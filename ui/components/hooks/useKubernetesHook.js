import { useDispatch } from 'react-redux';
import { updateProgress } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import { errorHandlerGenerator, successHandlerGenerator } from '../ConnectionWizard/helpers/common';
import { pingKubernetes } from '../ConnectionWizard/helpers/kubernetesHelpers';
import { pingMesheryOperator } from '../ConnectionWizard/helpers/mesheryOperator';
import { EVENT_TYPES } from '../../lib/event-types';
import MeshsyncStatusQuery from '../graphql/queries/MeshsyncStatusQuery';
import { useEffect, useState } from 'react';
import fetchMesheryOperatorStatus from '../graphql/queries/OperatorStatusQuery';
import NatsStatusQuery from '../graphql/queries/NatsStatusQuery';
import { CONTROLLERS, CONTROLLER_STATES } from '../../utils/Enum';
import _ from 'lodash';

export default function useKubernetesHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const ping = (name, server, connectionID) => {
    dispatch(updateProgress({ showProgress: true }));
    pingKubernetes(
      successHandlerGenerator(notify, `Kubernetes context ${name} at ${server} pinged`, () =>
        dispatch(updateProgress({ showProgress: false })),
      ),
      errorHandlerGenerator(notify, `Kubernetes context ${name} at ${server} not reachable`, () =>
        dispatch(updateProgress({ showProgress: false })),
      ),
      connectionID,
    );
  };

  return ping;
}
const handleErrorGenerator = (dispatch, notify) => (message, error) => {
  dispatch(updateProgress({ showProgress: false }));
  console.error(message, error);
  notify({
    message: message,
    details: error,
    event_type: EVENT_TYPES.ERROR,
  });
};

const handleSuccessGenerator = (dispatch, notify) => (message) => {
  dispatch(updateProgress({ showProgress: false }));
  notify({
    message: message,
    event_type: EVENT_TYPES.SUCCESS,
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

  const ping = (contextID) => {
    dispatch(updateProgress({ showProgress: true }));
    pingMesheryOperator(
      contextID,
      () => handleSuccess(`Meshery Operator  pinged`),
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

  // takes connectionID as input not the contextID
  const ping = ({ connectionID, subscribe = false, onSuccess, onError }) => {
    dispatch(updateProgress({ showProgress: true }));

    const subscription = MeshsyncStatusQuery({ connectionID: connectionID }).subscribe({
      next: (res) => {
        dispatch(updateProgress({ showProgress: false }));

        if (res.controller.name === 'MeshSync' && res.controller.status.includes('Connected')) {
          let publishEndpoint = res.controller.status.substring('Connected'.length);
          handleSuccess(
            `MeshSync was pinged. ${
              publishEndpoint != '' ? `Publishing to ${publishEndpoint}` : ''
            }`,
          );
        } else if (
          res.controller.name === 'MeshSync' &&
          !res.controller.status.includes('Unknown')
        ) {
          handleError('MeshSync is not publishing to Meshery Broker');
        } else {
          handleError('MeshSync could not be reached');
        }
        onSuccess && onSuccess(res);

        !subscribe && subscription && subscription?.unsubscribe();
      },
      error: (err) => {
        dispatch(updateProgress({ showProgress: false }));
        handleError('MeshSync status could not be retrieved', err);
        onError && onError(err);
        !subscribe && subscription && subscription?.unsubscribe();
      },
    });

    return subscription;
  };

  return { ping };
}

export const useGetOperatorInfoQuery = ({ contextID }) => {
  const { notify } = useNotification();
  const dispatch = useDispatch();

  const handleError = handleErrorGenerator(dispatch, notify);
  // const handleSuccess = handleSuccessGenerator(dispatch, notify);
  const handleInfo = handleInfoGenerator(notify);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setIsLoading(true);
    dispatch(updateProgress({ showProgress: true }));
    handleInfo('Fetching Meshery Operator status');
    const tempSubscription = fetchMesheryOperatorStatus({ k8scontextID: contextID }).subscribe({
      next: () => {
        setIsLoading(false);

        dispatch(updateProgress({ showProgress: false }));
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
      dispatch(updateProgress({ showProgress: false }));
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

  const ping = ({ connectionID, subscribe = false, onSuccess, onError }) => {
    dispatch(updateProgress({ showProgress: true }));
    const subscription = NatsStatusQuery({ connectionID }).subscribe({
      next: (res) => {
        dispatch(updateProgress({ showProgress: false }));

        if (
          res.controller.name === 'MesheryBroker' &&
          res.controller.status.includes('Connected')
        ) {
          let runningEndpoint = res.controller.status.substring('Connected'.length);
          handleSuccess(
            `Broker was pinged. ${runningEndpoint != '' ? `Running at ${runningEndpoint}` : ''}`,
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
  const getContextStatus = (ctxId) => {
    const defaultState = {
      operatorState: CONTROLLER_STATES.DISABLED,
      operatorVersion: 'Not Available',
      meshSyncState: CONTROLLER_STATES.DISABLED,
      meshSyncVersion: 'Not Available',
      natsState: CONTROLLER_STATES.DISABLED,
      natsVersion: 'Not Available',
    };

    const controller = controllerState?.filter((op) => op.contextId === ctxId);
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

    function getOperatorStatus(ctxId) {
      const operator = controllerState?.find(
        (op) => op.contextId === ctxId && op.controller === CONTROLLERS.OPERATOR,
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
      ...getOperatorStatus(ctxId),
      ...getMeshSyncStats(),
      ...getBrokerStats(),
    };

    return _.merge(defaultState, actualOperatorState);
  };

  return {
    getControllerStatesByContexID: getContextStatus,
  };
};
