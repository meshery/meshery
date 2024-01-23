import { CONTROLLER_STATES } from '../../../utils/Enum';
import fetchMesheryOperatorStatus from '../../graphql/queries/OperatorStatusQuery';

export const isMesheryOperatorConnected = ({ operatorInstalled }) => operatorInstalled;

/**
 * Pings meshery operator
 * @param {() => Function} fetchMesheryOperatorStatus - function with which
 * we can query using graphql
 * @param  {(res) => void} successHandler
 * @param  {(err) => void} errorHandler
 */
export const pingMesheryOperator = (id, successcb, errorcb) => {
  const subscription = fetchMesheryOperatorStatus({
    k8scontextID: id,
  }).subscribe({
    next: (data) => {
      if (
        data === null ||
        data?.operator === null ||
        data?.operator?.status === CONTROLLER_STATES.UNKOWN
      ) {
        errorcb();
        subscription.unsubscribe();
        return;
      }
      successcb();
      subscription.unsubscribe();
    },
    error: () => {
      errorcb();
      subscription.unsubscribe();
    },
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
