import { CONNECTION_KINDS } from '@/utils/Enum';
import { useNotification } from '@/utils/hooks/useNotification';
import dataFetch from 'lib/data-fetch';
import { EVENT_TYPES } from 'lib/event-types';
import { updateProgress } from 'lib/store';
import { useDispatch } from 'react-redux';

export function usePrometheusHook() {
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
        dispatch(updateProgress({ showProgress: false }));
        if (typeof result !== 'undefined') {
          notify({
            message: `Prometheus connection "${name}" pinged at ${server}`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      },
      self.handleError,
    );
  };
  return ping;
}

export function useGrafanaHook() {
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
        dispatch(updateProgress({ showProgress: false }));
        if (typeof result !== 'undefined') {
          notify({
            message: `Grafana connection "${name}" pinged at ${server}`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      },
      self.handleError,
    );
  };
  return ping;
}

export function withTelemetryHook(Component, telemetryConnType) {
  return function CompWrappedWithTelemetryHook (props) {
    let ping;
    if (telemetryConnType === CONNECTION_KINDS.PROMETHEUS) {
      ping = usePrometheusHook();
    } else {
      ping = useGrafanaHook();
    }
    return (<Component {...props} ping={ping} />);
  };
}
