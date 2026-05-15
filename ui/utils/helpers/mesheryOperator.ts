import { CONTROLLER_STATES } from '../Enum';
import { store } from '@/store/index';
import { api } from '@/rtk-query/index';
// Importing kubernetes here triggers the endpoint injection side-effect,
// so api.endpoints.getMesheryOperatorStatus is registered when this helper
// is used outside of a React render tree.
import '@/rtk-query/kubernetes';

export const isMesheryOperatorConnected = ({ operatorInstalled }) => operatorInstalled;

/**
 * Pings meshery operator (one-shot). Replaces the legacy Relay fetchQuery
 * with a direct dispatch of the RTK Query endpoint. We bypass the React
 * hook because this helper is called from imperative call sites that need
 * to fire-and-forget a single request.
 *
 * @param id connection ID for the kubernetes context
 * @param successcb invoked when operator is reachable; receives
 *        `{ operator: { status, controller, connectionID, ... } }`.
 * @param errorcb invoked when operator is unreachable or status is UNKNOWN.
 */
export const pingMesheryOperator = (id, successcb, errorcb) => {
  // forceRefetch: true matches the previous one-shot semantics — we do
  // not want a cached value here, the caller is explicitly pinging.
  const promise = store.dispatch(
    (api.endpoints as any).getMesheryOperatorStatus.initiate(
      { connectionID: id },
      { forceRefetch: true },
    ),
  );

  promise
    .unwrap()
    .then((data) => {
      if (
        data === null ||
        data?.operator === null ||
        data?.operator?.status === CONTROLLER_STATES.UNKOWN
      ) {
        errorcb && errorcb(data);
        return;
      }
      successcb && successcb(data);
    })
    .catch((err) => {
      errorcb && errorcb(err ?? new Error('Unknown error from pingMesheryOperator'));
    })
    .finally(() => {
      // Mirror the legacy behaviour of letting the cache entry get GC'd.
      promise.unsubscribe();
    });
};

/**
 * Returns the connection status of Operator, Meshsync, and Broker (NATS)
 * using the result of graphql `operatorStatusQuery` query
 *
 * @param {object} res - Result of the graphql query
 * @returns {[boolean, object]} result - array with final states and
 * reachability of operator
 */

export const getOperatorStatusFromQueryResult = (res) => {
  var operatorInformation = {
    operatorInstalled: false,
    NATSInstalled: false,
    meshSyncInstalled: false,
    operatorSwitch: false,
    operatorVersion: 'N/A',
    meshSyncVersion: 'N/A',
    NATSVersion: 'N/A',
  };

  if (res.operator?.error) {
    return [false, operatorInformation];
  }

  if (res.operator?.status === 'ENABLED') {
    res.operator?.controllers?.forEach((controller) => {
      operatorInformation = {
        ...operatorInformation,
        [controller.name]: controller,
      };
      if (controller.name === 'broker' && controller.status === 'ENABLED') {
        operatorInformation = {
          ...operatorInformation,
          NATSInstalled: true,
          NATSVersion: controller.version,
        };
      } else if (controller.name === 'meshsync' && controller.status === 'ENABLED') {
        operatorInformation = {
          ...operatorInformation,
          meshSyncInstalled: true,
          meshSyncVersion: controller.version,
        };
      }
    });

    operatorInformation = {
      ...operatorInformation,
      operatorInstalled: true,
      operatorVersion: res.operator?.version,
    };

    return [true, operatorInformation];
  }

  return [false, operatorInformation];
};
