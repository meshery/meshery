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

// The definition-authored description for a specific transition, or undefined
// when the definition does not describe it. Connection definitions
// (models/.../connections/*.json) are the source of truth for this copy; the
// transition modal supplies its own generic fallback, so no prompt is
// synthesized here.
export const getTransitionDescription = (
  transitionMap: ConnectionTransitionMap | undefined,
  currentStatus: string | undefined,
  targetStatus: string,
): string | undefined => {
  if (!currentStatus) {
    return undefined;
  }
  return transitionMap?.[currentStatus.toLowerCase()]?.find(
    (transition) => transition.nextState === targetStatus.toLowerCase(),
  )?.description;
};

export const CONNECTION_DOCS_URL = `https://docs.meshery.io/concepts/logical/connections#states-and-the-lifecycle-of-connections`;
export const ENVIRONMENT_DOCS_URL = `https://docs.meshery.io/concepts/logical/environments`;

// The table's column names follow the v1beta3 camelCase wire shape
// (createdAt, updatedAt), but the server's `order` query param addresses DB
// columns (created_at, updated_at - see SanitizeOrderInput in
// server/models/connection_persister.go). Translate a UI sort order like
// "createdAt desc" into its server form; unknown fields pass through
// unchanged, which also keeps older bookmarked URLs with snake_case sort
// params working.
const UI_TO_SERVER_SORT_COLUMN: Record<string, string> = {
  createdAt: 'created_at',
  updatedAt: 'updated_at',
};

const SERVER_TO_UI_SORT_COLUMN: Record<string, string> = Object.fromEntries(
  Object.entries(UI_TO_SERVER_SORT_COLUMN).map(([uiColumn, serverColumn]) => [
    serverColumn,
    uiColumn,
  ]),
);

export const toServerSortOrder = (sortOrder: string): string => {
  const trimmed = sortOrder.trim();
  // Guard against empty / whitespace-only input so we never emit " desc".
  if (!trimmed) return 'created_at desc';
  const [field, direction] = trimmed.split(/\s+/);
  const serverField = UI_TO_SERVER_SORT_COLUMN[field] ?? field;
  // SanitizeOrderInput accepts exactly "<column> <asc|desc>"; a bare column
  // would be silently dropped server-side, so default the direction.
  return `${serverField} ${direction || 'desc'}`;
};

// The inverse of toServerSortOrder, for the table's active-sort indicator:
// mui-datatables matches `options.sortOrder.name` against a column name, so a
// bookmarked snake_case param (created_at desc) has to be mapped back to the
// camelCase column before it reaches the table, or the indicator matches
// nothing and silently disappears.
export const toUiSortOrder = (sortOrder: string): string => {
  const trimmed = sortOrder.trim();
  if (!trimmed) return 'createdAt desc';
  const [field, direction] = trimmed.split(/\s+/);
  const uiField = SERVER_TO_UI_SORT_COLUMN[field] ?? field;
  return `${uiField} ${direction || 'desc'}`;
};
