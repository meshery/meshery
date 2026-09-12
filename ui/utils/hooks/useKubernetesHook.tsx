import { useNotification } from '../../utils/hooks/useNotification';
import { getErrorMessage } from '../../components/connections/ConnectionTable.constants';
import {
  useLazyPingKubernetesQuery,
  useLazyGetOperatorStatusQuery,
  useLazyGetMeshsyncStatusQuery,
  useLazyGetBrokerStatusQuery,
} from '@/rtk-query/connection';
import { EVENT_TYPES } from '../../lib/event-types';
import { useCallback } from 'react';
import { CONTROLLERS, CONTROLLER_STATES } from '../../utils/Enum';
import _ from 'lodash';
import { useDispatch, useSelector } from 'react-redux';
import { updateProgressAction } from '@/store/slices/mesheryUi';

export default function useKubernetesHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerPing] = useLazyPingKubernetesQuery();

  // Memoized so consumers can list `ping` in hook dep arrays without
  // invalidating their memos every render.
  const ping = useCallback(
    async (name, server, connectionID) => {
      dispatch(updateProgressAction({ showProgress: true }));
      try {
        const result = await triggerPing(connectionID).unwrap();

        const serverVersion = result?.server_version ?? result?.serverVersion ?? null;
        const message = serverVersion
          ? `Connected successfully to ${name} (${server}) - Kubernetes ${serverVersion}`
          : `Connected successfully to ${name} (${server})`;

        notify({
          message,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (err) {
        notify({
          message: `Connection failed for ${name} (${server}) - unable to reach cluster`,
          details: getErrorMessage(err, 'Unable to reach cluster'),
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
const handleErrorGenerator = (dispatch, notify) => (message, error) => {
  dispatch(updateProgressAction({ showProgress: false }));
  console.error(message, error);
  notify({
    message: message,
    details: error,
    event_type: EVENT_TYPES.ERROR,
  });
};

const handleSuccessGenerator =
  (dispatch, notify) =>
  (message, variant = 'success') => {
    dispatch(updateProgressAction({ showProgress: false }));
    const variantMap = {
      success: EVENT_TYPES.SUCCESS,
      info: EVENT_TYPES.INFO,
      warning: EVENT_TYPES.WARNING,
      error: EVENT_TYPES.ERROR,
    };
    notify({
      message,
      event_type: variantMap[variant] ?? EVENT_TYPES.SUCCESS,
    });
  };

const handleInfoGenerator = (notify) => (message) => {
  notify({
    message: message,
    event_type: EVENT_TYPES.INFO,
  });
};

export function useMesheryOperator() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerOperatorStatus] = useLazyGetOperatorStatusQuery();
  const handleError = handleErrorGenerator(dispatch, notify);
  const handleSuccess = handleSuccessGenerator(dispatch, notify);

  const ping = async ({ connectionID }) => {
    dispatch(updateProgressAction({ showProgress: true }));
    try {
      // REST returns the operator status object directly: { status, controller, ... }.
      const res = await triggerOperatorStatus(connectionID).unwrap();
      const status = String(res?.status ?? CONTROLLER_STATES.UNKNOWN)
        .trim()
        .toUpperCase();

      const statusToVariantMap = {
        [CONTROLLER_STATES.DEPLOYED]: 'success',
        [CONTROLLER_STATES.DEPLOYING]: 'info',
        [CONTROLLER_STATES.NOTDEPLOYED]: 'error',
        [CONTROLLER_STATES.UNKNOWN]: 'error',
      };
      const variant = statusToVariantMap[status] || 'warning';

      handleSuccess(`Meshery Operator status: ${status}`, variant);
    } catch (err) {
      handleError(`Meshery Operator not reachable`, err);
    }
  };

  return { ping };
}
export function useMeshsSyncController() {
  const { notify } = useNotification();
  const dispatch = useDispatch();

  const handleError = handleErrorGenerator(dispatch, notify);
  const handleSuccess = handleSuccessGenerator(dispatch, notify);
  const handleInfo = handleInfoGenerator(notify);

  const [triggerMeshsyncStatus] = useLazyGetMeshsyncStatusQuery();

  const ping = async ({ connectionID, onSuccess, onError }) => {
    dispatch(updateProgressAction({ showProgress: true }));
    try {
      // REST returns the MeshSync status object directly: { name, version, status, ... }.
      const res = await triggerMeshsyncStatus(connectionID).unwrap();
      dispatch(updateProgressAction({ showProgress: false }));

      if (res.name === 'MeshSync' && res.status.includes('Connected')) {
        let publishEndpoint = res.status.substring('Connected'.length);
        handleSuccess(
          `MeshSync was pinged. ${publishEndpoint != '' ? `Publishing to ${publishEndpoint}` : ''}`,
        );
      } else if (
        res.name === 'MeshSync' &&
        (res.status === 'Running' || res.status.includes('Running'))
      ) {
        handleInfo(
          `MeshSync is running${res.version ? ` (${res.version})` : ''}, but Meshery could not confirm a connection to the Meshery Broker (the broker may be unreachable from Meshery, or MeshSync is still connecting).`,
        );
      } else if (res.name === 'MeshSync' && res.status === 'Deployed') {
        handleInfo('MeshSync is deployed but connection status unclear');
      } else if (
        res.name === 'MeshSync' &&
        !res.status.includes('Unknown') &&
        !res.status.includes('UNKNOWN')
      ) {
        handleInfo('MeshSync is not publishing to Meshery Broker');
      } else {
        handleError('MeshSync could not be reached');
      }
      onSuccess && onSuccess(res);
    } catch (err) {
      dispatch(updateProgressAction({ showProgress: false }));
      handleError('MeshSync status could not be retrieved', err);
      onError && onError(err);
    }
  };

  return { ping };
}

export const useNatsController = () => {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const [triggerBrokerStatus] = useLazyGetBrokerStatusQuery();

  const handleError = handleErrorGenerator(dispatch, notify);
  const handleSuccess = handleSuccessGenerator(dispatch, notify);
  const handleInfo = handleInfoGenerator(notify);

  const ping = async ({ connectionID, onSuccess, onError }) => {
    dispatch(updateProgressAction({ showProgress: true }));
    try {
      // REST returns the broker status object directly: { name, version, status, ... }.
      const res = await triggerBrokerStatus(connectionID).unwrap();
      dispatch(updateProgressAction({ showProgress: false }));

      if (res.name === 'MesheryBroker' && res.status.includes('Connected')) {
        let runningEndpoint = res.status.substring('Connected'.length);
        handleSuccess(
          `Broker was pinged. ${runningEndpoint != '' ? `Running at ${runningEndpoint}` : ''}`,
        );
      } else if (
        res.name === 'MesheryBroker' &&
        (res.status === 'Deployed' || res.status === 'DEPLOYED')
      ) {
        handleInfo(
          `Meshery Broker is deployed${res.version ? ` (${res.version})` : ''} but not connected to Meshery Server`,
        );
      } else {
        handleError(
          'Meshery Broker could not be reached',
          'Meshery Server is not connected to Meshery Broker',
        );
      }

      onSuccess &&
        onSuccess({
          rawResponse: res,
          isReachable: true,
          natsStatus: res.status,
          NATSVersion: res.version,
        });
    } catch (err) {
      onError && onError(err);
      handleError('NATS status could not be retrieved', err);
    }
  };

  return {
    ping,
  };
};

export const useControllerStatus = (controllerState) => {
  const getContextStatus = (connectionID) => {
    const defaultState = {
      operatorState: CONTROLLER_STATES.DISABLED,
      operatorVersion: 'Not Available',
      meshSyncState: CONTROLLER_STATES.DISABLED,
      meshSyncVersion: 'Not Available',
      natsState: CONTROLLER_STATES.DISABLED,
      natsVersion: 'Not Available',
    };

    const controller = controllerState?.filter((op) => op.connectionId === connectionID);
    if (!controller) {
      return defaultState;
    }

    function getMeshSyncStats() {
      if (!controller) return defaultState;
      const meshsyncController = controller?.find(
        (ctlr) => ctlr?.controller === CONTROLLERS.MESHSYNC,
      );
      // meshsync is at 1st idx
      if (meshsyncController) {
        return {
          meshSyncState: meshsyncController?.status,
          meshSyncVersion: meshsyncController?.version,
        };
      }
    }

    function getBrokerStats() {
      if (!controller) return defaultState;
      const brokerController = controller?.find((ctlr) => ctlr?.controller === CONTROLLERS.BROKER);
      if (brokerController) {
        return {
          natsState: brokerController?.status,
          natsVersion: brokerController?.version,
        };
      }
    }

    function getOperatorStatus(connectionID) {
      const operator = controllerState?.find(
        (op) => op.connectionId === connectionID && op.controller === CONTROLLERS.OPERATOR,
      );
      if (!operator) {
        return defaultState;
      }

      return {
        operatorState: operator.status,
        operatorVersion: operator?.version,
      };
    }

    const actualOperatorState = {
      ...getOperatorStatus(connectionID),
      ...getMeshSyncStats(),
      ...getBrokerStats(),
    };

    return _.merge(defaultState, actualOperatorState);
  };

  return {
    getControllerStatesByConnectionID: getContextStatus,
  };
};

export const useFilterK8sContexts = (k8sContexts, predicate) => {
  const { controllerState: meshsyncControllerState } = useSelector((state) => state.ui);
  const { getControllerStatesByConnectionID } = useControllerStatus(meshsyncControllerState);

  return k8sContexts.filter((ctx) => {
    const operatorsStatus = getControllerStatesByConnectionID(ctx.connectionId);
    return predicate({ ...operatorsStatus, context: ctx });
  });
};
