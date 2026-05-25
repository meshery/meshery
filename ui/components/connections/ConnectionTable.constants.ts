export const getErrorMessage = (error: any, fallback = 'Unknown error') => {
  if (!error) {
    return fallback;
  }

  if (typeof error === 'string') {
    return error;
  }

  if (typeof error === 'object') {
    if ('error' in error && typeof error.error === 'string') {
      return error.error;
    }

    if ('message' in error && typeof error.message === 'string') {
      return error.message;
    }

    if ('data' in error && typeof error.data === 'string') {
      return error.data;
    }
  }

  return fallback;
};

export const ACTION_TYPES = {
  FETCH_CONNECTIONS: {
    name: 'FETCH_CONNECTIONS',
    error_msg: 'Failed to fetch connections',
  },
  UPDATE_CONNECTION: {
    name: 'UPDATE_CONNECTION',
    error_msg: 'Failed to update connection',
  },
  DELETE_CONNECTION: {
    name: 'DELETE_CONNECTION',
    error_msg: 'Failed to delete connection',
  },
  FETCH_CONNECTION_STATUS_TRANSITIONS: {
    name: 'FETCH_CONNECTION_STATUS_TRANSITIONS',
    error_msg: 'Failed to fetch connection transitions',
  },
  FETCH_ENVIRONMENT: {
    name: 'FETCH_ENVIRONMENT',
    error_msg: 'Failed to fetch environment',
  },
  CREATE_ENVIRONMENT: {
    name: 'CREATE_ENVIRONMENT',
    error_msg: 'Failed to create environment',
  },
};

const kubernetesConnectionTransitions = {
  connected: {
    disconnected:
      'Are you sure you want to transition from CONNECTED to DISCONNECTED? This will perform planned maintenance by removing the operator but keeping the cluster registered.',
    ignored:
      'Are you sure you want to transition from CONNECTED to IGNORED? This will mark the connection as ignored due to unplanned maintenance, without deleting the registration.',
    deleted:
      'Are you sure you want to transition from CONNECTED to DELETED? This will undeploy the operator and unregister the cluster completely.',
    'not found':
      'Are you sure you want to transition from CONNECTED to NOT FOUND? Meshery could not connect to the cluster or it is currently unavailable. You can either delete the connection or try re-registering.',
  },
  disconnected: {
    connected:
      'Are you sure you want to transition from DISCONNECTED to CONNECTED? This will reconnect the cluster and redeploy the operator after maintenance.',
    deleted:
      'Are you sure you want to transition from DISCONNECTED to DELETED? This will remove the cluster completely by undeploying the operator and unregistering.',
  },
  ignored: {
    deleted:
      'Are you sure you want to transition from IGNORED to DELETED? This will completely remove the ignored cluster by undeploying the operator and unregistering.',
    registered:
      'Are you sure you want to transition from IGNORED to REGISTER? This will reinitiate the registration process for the ignored connection and attempt to connect it again.',
  },
  'not found': {
    discovered:
      'Are you sure you want to transition from NOT FOUND to DISCOVERED? You are trying to re-register the cluster. Meshery will attempt to reconnect to the cluster.',
    deleted:
      'Are you sure you want to transition from NOT FOUND to DELETED? This will remove the unreachable connection completely by unregistering it.',
  },
};

// <connection-kind>: TransitionMap (status:{status:description})
export const CONNECTION_STATE_TRANSITIONS = {
  kubernetes: kubernetesConnectionTransitions,
};

export const getStatusTransition = (
  connectionKind: string,
  connectionState: string,
  transitionState: string,
) => {
  // This is for one connection kind that is kubernetes, and adding other connection kinds
  // here will make it more complex.
  // This issue can be resolved if we add the transition messages in the connection schemas
  // and use the same schema to get the transition messages.
  // Github issue: https://github.com/meshery/schemas/issues/303

  switch (connectionKind) {
    case 'kubernetes':
      return kubernetesConnectionTransitions[connectionState][transitionState];
    default:
      return `Are you sure you want to transition from ${connectionState} to ${transitionState}?`;
  }
};

export const CONNECTION_DOCS_URL = `https://docs.meshery.io/concepts/logical/connections#states-and-the-lifecycle-of-connections`;
export const ENVIRONMENT_DOCS_URL = `https://docs.meshery.io/concepts/logical/environments`;
