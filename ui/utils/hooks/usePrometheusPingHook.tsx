import { useCallback } from 'react';
import { useDispatch } from 'react-redux';
import { useNotification } from './useNotification';
import { getErrorMessage } from '../../components/connections/ConnectionTable.constants';
import { useLazyPingPrometheusConnectionQuery } from '@/rtk-query/telemetryPrometheus';
import { EVENT_TYPES } from '../../lib/event-types';
import { updateProgressAction } from '@/store/slices/mesheryUi';

/**
 * Imperative `ping(connectionID, name)` for a registered Prometheus connection.
 * Mirrors useGrafanaPingHook; hits GET /api/telemetry/prometheus/{id}/ping and
 * surfaces reachability (+ version when returned) as a notification.
 */
export default function usePrometheusPingHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerPing] = useLazyPingPrometheusConnectionQuery();

  const ping = useCallback(
    async (connectionID: string, name?: string) => {
      const label = name || 'Prometheus';
      dispatch(updateProgressAction({ showProgress: true }));
      try {
        const result = await triggerPing({ connectionID }).unwrap();
        const version = result?.version;
        notify({
          message: version
            ? `Connected successfully to ${label} (Prometheus ${version})`
            : `Connected successfully to ${label}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (err) {
        notify({
          message: `Connection failed for ${label} — unable to reach Prometheus`,
          details: getErrorMessage(err, 'Unable to reach Prometheus'),
          event_type: EVENT_TYPES.ERROR,
        });
      } finally {
        dispatch(updateProgressAction({ showProgress: false }));
      }
    },
    [dispatch, notify, triggerPing],
  );

  return ping;
}
