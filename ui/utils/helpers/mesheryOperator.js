import { CONTROLLER_STATES } from '../Enum';
import fetchMesheryOperatorStatus from '../../components/graphql/queries/OperatorStatusQuery';

export const isMesheryOperatorConnected = ({ operatorInstalled }) => operatorInstalled;

/**
 * Pings meshery operator
 * @param {() => Function} fetchMesheryOperatorStatus - function with which
 * we can query using graphql
 * @param  {(res: object) => void} successHandler - called when operator is reachable,
 *         receives the GraphQL result (e.g. { operator: { status, controller, connectionID } }).
 * @param  {(err: any) => void} errorHandler - called when operator is unreachable or status is UNKNOWN,
 *         receives the error object or response for context.
 */
export const pingMesheryOperator = (id, successcb, errorcb) => {
  const subscription = fetchMesheryOperatorStatus({
    connectionID: id,
  }).subscribe({
    next: (data) => {
      if (
        data === null ||
        data?.operator === null ||
        data?.operator?.status === CONTROLLER_STATES.UNKOWN
      ) {
        errorcb && errorcb(data);
        subscription.unsubscribe();
        return;
      }
      successcb && successcb(data);
      subscription.unsubscribe();
    },
    error: (err) => {
      errorcb && errorcb(err ?? new Error('Unknown error from pingMesheryOperator'));
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
