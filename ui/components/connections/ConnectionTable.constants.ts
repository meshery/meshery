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
  WorkspaceManagementCreateEnvironment: {
    name: 'WorkspaceManagementCreateEnvironment',
    error_msg: 'Failed to create environment',
  },
};

// A single permissible state transition for a connection, mirroring the
// `ConnectionStateTransition` schema (meshery/schemas v1beta3/connection).
export type ConnectionStateTransition = {
  nextState: string;
  description?: string;
};

// `transitionMap` from a connection definition: keyed by current status, each
// value is the list of states reachable from that status. Authored per-kind in
// meshery core (models/.../connections/*.json) and surfaced to the UI
// via the connection definitions, replacing the previously hardcoded map.
export type ConnectionTransitionMap = Record<string, ConnectionStateTransition[]>;

// The states a connection may transition to from its current status.
export const getNextStates = (
  transitionMap: ConnectionTransitionMap | undefined,
  currentStatus: string,
): string[] => (transitionMap?.[currentStatus] ?? []).map((transition) => transition.nextState);

// The human-readable description for a specific transition, falling back to a
// generic prompt when the definition does not describe it.
export const getStatusTransition = (
  transitionMap: ConnectionTransitionMap | undefined,
  connectionState: string,
  transitionState: string,
) => {
  const transition = transitionMap?.[connectionState]?.find((t) => t.nextState === transitionState);

  return (
    transition?.description ||
    `Are you sure you want to transition from ${connectionState.toUpperCase()} to ${transitionState.toUpperCase()}?`
  );
};

export const CONNECTION_DOCS_URL = `https://docs.meshery.io/concepts/logical/connections#states-and-the-lifecycle-of-connections`;
export const ENVIRONMENT_DOCS_URL = `https://docs.meshery.io/concepts/logical/environments`;
