export const isMesheryOperatorConnected = (params?: { operatorInstalled?: boolean } | null) =>
  params?.operatorInstalled;

/**
 * Returns the connection status of Operator, Meshsync, and Broker (NATS)
 * using the result of graphql `operatorStatusQuery` query
 *
 * @param {object} res - Result of the graphql query
 * @returns {[boolean, object]} result - array with final states and
 * reachability of operator
 */

export const getOperatorStatusFromQueryResult = (res?: unknown) => {
  let operatorInformation = {
    operatorInstalled: false,
    NATSInstalled: false,
    meshSyncInstalled: false,
    operatorSwitch: false,
    operatorVersion: 'N/A',
    meshSyncVersion: 'N/A',
    NATSVersion: 'N/A',
  };

  const queryResult = res as {
    operator?: {
      error?: unknown;
      status?: string;
      version?: string;
      controllers?: Array<{
        name?: string;
        status?: string;
        version?: string;
        [key: string]: unknown;
      }>;
    };
  } | null;

  if (!queryResult || typeof queryResult !== 'object' || queryResult.operator?.error) {
    return [false, operatorInformation];
  }

  if (queryResult.operator?.status === 'ENABLED') {
    if (Array.isArray(queryResult.operator?.controllers)) {
      queryResult.operator.controllers.forEach((controller) => {
        if (
          !controller ||
          typeof controller !== 'object' ||
          typeof controller.name !== 'string' ||
          controller.name === ''
        ) {
          return;
        }

        operatorInformation = {
          ...operatorInformation,
          [controller.name]: controller,
        };
        if (controller.name === 'broker' && controller.status === 'ENABLED') {
          operatorInformation = {
            ...operatorInformation,
            NATSInstalled: true,
            NATSVersion: controller.version || 'N/A',
          };
        } else if (controller.name === 'meshsync' && controller.status === 'ENABLED') {
          operatorInformation = {
            ...operatorInformation,
            meshSyncInstalled: true,
            meshSyncVersion: controller.version || 'N/A',
          };
        }
      });
    }

    operatorInformation = {
      ...operatorInformation,
      operatorInstalled: true,
      operatorVersion: queryResult.operator?.version || 'N/A',
    };

    return [true, operatorInformation];
  }

  return [false, operatorInformation];
};
