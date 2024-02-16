import { CONNECTION_KINDS } from '@/utils/Enum';
import { useNotification } from '@/utils/hooks/useNotification';
import dataFetch from 'lib/data-fetch';
import { EVENT_TYPES } from 'lib/event-types';
import { updateProgress } from 'lib/store';
import { useDispatch } from 'react-redux';

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
  const ping = (name, server, connectionID) => {
    dispatch(updateProgress({ showProgress: true }));
    dataFetch(
      `/api/telemetry/metrics/ping/${connectionID}`,
      {
        credentials: 'include',
      },
      (result) => {
        if (typeof result !== 'undefined') {
          notify({
            message: `Prometheus connection "${name}" pinged at ${server}`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      },
      self.handleError,
    );
    dispatch(updateProgress({ showProgress: false }));
  };
  return ping;
}

function PingGrafana() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const ping = (name, server, connectionID) => {
    dispatch(updateProgress({ showProgress: true }));
    dataFetch(
      `/api/telemetry/metrics/grafana/ping/${connectionID}`,
      {
        credentials: 'include',
      },
      (result) => {
        if (typeof result !== 'undefined') {
          notify({
            message: `Grafana connection "${name}" pinged at ${server}`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      },
      self.handleError,
    );
    dispatch(updateProgress({ showProgress: false }));
  };
  return ping;
}

export function withTelemetryHook(Component, telemetryConnType) {
  return function CompWrappedWithTelemetryHook(props) {
    const ping = useTelemetryHook(telemetryConnType);
    return <Component {...props} ping={ping} />;
  };
}
