import { useNotification } from '../../utils/hooks/useNotification';
import { errorHandlerGenerator, successHandlerGenerator } from '../../utils/helpers/common';
import { pingKubernetes } from '../../utils/helpers/kubernetesHelpers';
import { pingMesheryOperator } from '../../utils/helpers/mesheryOperator';
import { EVENT_TYPES } from '../../lib/event-types';
import MeshsyncStatusQuery from '../graphql/queries/MeshsyncStatusQuery';
import { useEffect, useState } from 'react';
import fetchMesheryOperatorStatus from '../graphql/queries/OperatorStatusQuery';
import NatsStatusQuery from '../graphql/queries/NatsStatusQuery';
import { CONTROLLERS, CONTROLLER_STATES } from '../../utils/Enum';
import _ from 'lodash';
import { useDispatch, useSelector } from 'react-redux';
import { updateProgressAction } from '@/store/slices/mesheryUi';

export default function useKubernetesHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const ping = (name, server, connectionID) => {
    dispatch(updateProgressAction({ showProgress: true }));
    pingKubernetes(
      successHandlerGenerator(notify, `Kubernetes context ${name} at ${server} pinged`, () =>
        dispatch(updateProgressAction({ showProgress: false })),
      ),
      errorHandlerGenerator(notify, `Kubernetes context ${name} at ${server} not reachable`, () =>
        dispatch(updateProgressAction({ showProgress: false })),
      ),
      connectionID,
    );
  };

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
  const handleError = handleErrorGenerator(dispatch, notify);
  const handleSuccess = handleSuccessGenerator(dispatch, notify);

  const ping = ({ connectionID }) => {
    dispatch(updateProgressAction({ showProgress: true }));
    pingMesheryOperator(
      connectionID,
      (res) => {
        const status = String(res?.operator?.status ?? CONTROLLER_STATES.UNKNOWN)
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
      },
      (err) => handleError(`Meshery Operator not reachable`, err),
    );
  };

  return { ping };
}
export function useMeshsSyncController() {
  const { notify } = useNotification();
  const dispatch = useDispatch();

  const handleError = handleErrorGenerator(dispatch, notify);
  const handleSuccess = handleSuccessGenerator(dispatch, notify);
  const handleInfo = handleInfoGenerator(notify);

  const ping = ({ connectionID, subscribe = false, onSuccess, onError }) => {
    dispatch(updateProgressAction({ showProgress: true }));

    const subscription = MeshsyncStatusQuery({ connectionID: connectionID }).subscribe({
      next: (res) => {
        dispatch(updateProgressAction({ showProgress: false }));

        if (res.controller.name === 'MeshSync' && res.controller.status.includes('Connected')) {
          let publishEndpoint = res.controller.status.substring('Connected'.length);
          handleSuccess(
            `MeshSync was pinged. ${
              publishEndpoint != '' ? `Publishing to ${publishEndpoint}` : ''
            }`,
          );
        } else if (
          res.controller.name === 'MeshSync' &&
          (res.controller.status === 'Running' || res.controller.status.includes('Running'))
        ) {
          handleInfo(
            `MeshSync is running (${res.controller.version}), but is not connected to Meshery Broker.`,
          );
        } else if (res.controller.name === 'MeshSync' && res.controller.status === 'Deployed') {
          handleInfo('MeshSync is deployed but connection status unclear');
        } else if (
          res.controller.name === 'MeshSync' &&
          !res.controller.status.includes('Unknown') &&
          !res.controller.status.includes('UNKNOWN')
        ) {
          handleInfo('MeshSync is not publishing to Meshery Broker');
        } else {
          handleError('MeshSync could not be reached');
        }
        onSuccess && onSuccess(res);

        !subscribe && subscription && subscription?.unsubscribe();
      },
      error: (err) => {
        dispatch(updateProgressAction({ showProgress: false }));
        handleError('MeshSync status could not be retrieved', err);
        onError && onError(err);
        !subscribe && subscription && subscription?.unsubscribe();
      },
    });

    return subscription;
  };

  return { ping };
}

export const useGetOperatorInfoQuery = ({ connectionID }) => {
  const { notify } = useNotification();
  const dispatch = useDispatch();

  const handleError = handleErrorGenerator(dispatch, notify);
  // const handleSuccess = handleSuccessGenerator(dispatch, notify);
  const handleInfo = handleInfoGenerator(notify);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setIsLoading(true);
    dispatch(updateProgressAction({ showProgress: true }));
    handleInfo('Fetching Meshery Operator status');
    // react-realy fetchQuery function returns a "Observable". To start a request subscribe needs to be called.
    // The data is stored into the react-relay store, the data is retrieved by subscribing to the relay store.
    // This subscription only subscribes to the fetching of the query and not to any subsequent changes to data in the relay store.
    const tempSubscription = fetchMesheryOperatorStatus({ connectionID: connectionID }).subscribe({
      next: () => {
        setIsLoading(false);

        dispatch(updateProgressAction({ showProgress: false }));
        // const [isReachable, operatorInfo] = getOperatorStatusFromQueryResult(res);
        // setOperatorInfo({
        //   isReachable,
        //   ...operatorInfo,
        //   meshSyncStatus: operatorInfo.MeshSync ? operatorInfo.MeshSync.status : '',
        //   meshSyncVersion: operatorInfo.MeshSync ? operatorInfo.MeshSync.version : '',
        //   NATSVersion: operatorInfo.MesheryBroker ? operatorInfo.MesheryBroker.version : '',
        //   natsStatus: operatorInfo.MesheryBroker ? operatorInfo.MesheryBroker.status : '',
        // });
      },
      error: () => {
        setIsLoading(false);
        handleError('Meshery Operator status could not be retrieved');
      },
    });
    return () => {
      setIsLoading(false);
      dispatch(updateProgressAction({ showProgress: false }));
      tempSubscription?.unsubscribe();
    };
  }, []);

  return {
    isLoading,
  };
};

export const useNatsController = () => {
  const { notify } = useNotification();
  const dispatch = useDispatch();

  const handleError = handleErrorGenerator(dispatch, notify);
  const handleSuccess = handleSuccessGenerator(dispatch, notify);
  const handleInfo = handleInfoGenerator(notify);

  const ping = ({ connectionID, subscribe = false, onSuccess, onError }) => {
    dispatch(updateProgressAction({ showProgress: true }));
    const subscription = NatsStatusQuery({ connectionID }).subscribe({
      next: (res) => {
        dispatch(updateProgressAction({ showProgress: false }));

        if (
          res.controller.name === 'MesheryBroker' &&
          res.controller.status.includes('Connected')
        ) {
          let runningEndpoint = res.controller.status.substring('Connected'.length);
          handleSuccess(
            `Broker was pinged. ${runningEndpoint != '' ? `Running at ${runningEndpoint}` : ''}`,
          );
        } else if (
          res.controller.name === 'MesheryBroker' &&
          (res.controller.status === 'Deployed' || res.controller.status === 'DEPLOYED')
        ) {
          handleInfo(
            `Meshery Broker is deployed (${res.controller.version}) but not connected to Meshery Server`,
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
            natsStatus: res.controller.status,
            NATSVersion: res.controller.version,
          });

        !subscribe && subscription?.unsubscribe();
      },
      error: (err) => {
        onError && onError(err);
        handleError('NATS status could not be retrieved', err);
        !subscribe && subscription?.unsubscribe();
      },
    });

    return subscription;
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

    const controller = controllerState?.filter((op) => op.connectionID === connectionID);
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
        (op) => op.connectionID === connectionID && op.controller === CONTROLLERS.OPERATOR,
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
    const operatorsStatus = getControllerStatesByConnectionID(ctx.connection_id);
    return predicate({ ...operatorsStatus, context: ctx });
  });
};
