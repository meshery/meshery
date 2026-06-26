import { useCallback } from 'react';
import { useDispatch } from 'react-redux';
import { useNotification } from './useNotification';
import { getErrorMessage } from '../../components/connections/ConnectionTable.constants';
import { useLazyPingGrafanaConnectionQuery } from '@/rtk-query/telemetryGrafana';
import { EVENT_TYPES } from '../../lib/event-types';
import { updateProgressAction } from '@/store/slices/mesheryUi';

/**
 * useGrafanaPingHook returns an imperative `ping(connectionID, name)` that checks
 * a registered Grafana connection's reachability via the telemetry ping endpoint
 * and surfaces the outcome as a notification. The server additionally emits a
 * persisted connection event for the result.
 */
export default function useGrafanaPingHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerPing] = useLazyPingGrafanaConnectionQuery();

  // Memoized so consumers can list `ping` in hook dep arrays without
  // invalidating their memos every render.
  const ping = useCallback(
    async (connectionID: string, name?: string) => {
      const label = name || 'Grafana';
      dispatch(updateProgressAction({ showProgress: true }));
      try {
        const result = await triggerPing({ connectionID }).unwrap();
        const version = result?.version;
        notify({
          message: version
            ? `Connected successfully to ${label} (Grafana ${version})`
            : `Connected successfully to ${label}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (err) {
        notify({
          message: `Connection failed for ${label} — unable to reach Grafana`,
          details: getErrorMessage(err, 'Unable to reach Grafana'),
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
