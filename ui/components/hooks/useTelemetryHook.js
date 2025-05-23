import { useDispatch } from 'react-redux';
import { updateProgressAction } from '@/store/slices/mesheryUi';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useLazyPingPrometheusQuery, useLazyPingGrafanaQuery } from '@/rtk-query/telemetry';

export function useTelemetryHook(connectionType) {
  switch (connectionType) {
    case CONNECTION_KINDS.PROMETHEUS:
      return PingPrometheus();
    case CONNECTION_KINDS.GRAFANA:
      return PingGrafana();
  }
}

function PingPrometheus() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerPingPrometheus] = useLazyPingPrometheusQuery();

  const ping = async (name, server, connectionID) => {
    dispatch(updateProgressAction({ showProgress: true }));

    try {
      const result = await triggerPingPrometheus({ connectionId: connectionID }).unwrap();
      if (typeof result !== 'undefined') {
        notify({
          message: `Prometheus connection "${name}" pinged at ${server}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      }
    } catch (error) {
      let cleanErrorMessage = 'There was an error communicating with Prometheus';
      let serverError;
      if (error && typeof error.data === 'string') {
        serverError = error.data.replace(/^Status Code: \d+\.\s*/, '').trim();
      } else if (error.error) {
        serverError = error.error;
      }
      notify({
        message: `${cleanErrorMessage}: ${serverError || ''}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }

    dispatch(updateProgressAction({ showProgress: false }));
  };

  return ping;
}

function PingGrafana() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerPingGrafana] = useLazyPingGrafanaQuery();

  const ping = async (name, server, connectionID) => {
    dispatch(updateProgressAction({ showProgress: true }));

    try {
      const result = await triggerPingGrafana({ connectionId: connectionID }).unwrap();
      if (typeof result !== 'undefined') {
        notify({
          message: `Grafana connection "${name}" pinged at ${server}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      }
    } catch {
      notify({
        message: 'There was an error communicating with Grafana',
        event_type: EVENT_TYPES.ERROR,
      });
    }

    dispatch(updateProgressAction({ showProgress: false }));
  };

  return ping;
}

export function withTelemetryHook(Component, telemetryConnType) {
  return function CompWrappedWithTelemetryHook(props) {
    const ping = useTelemetryHook(telemetryConnType);
    return <Component {...props} ping={ping} />;
  };
}
