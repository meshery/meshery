export const FILE_OPS = {
  FILE_UPLOAD: 'upload',
  URL_UPLOAD: 'url_upload',
  UPDATE: 'update',
  DELETE: 'delete',
  DOWNLOAD: 'download',
  CLONE: 'clone',
};

export const CON_OPS = {
  DELETE: 'delete',
  UPDATE: 'update',
  CREATE: 'create',
};
export const ACTIONS = {
  DEPLOY: 2,
  UNDEPLOY: 1,
  VERIFY: 0,
};

export const DEPLOYMENT_TYPE = {
  IN_CLUSTER: 'in_cluster',
  OUT_CLUSTER: 'out_of_cluster',
};

export const VISIBILITY = {
  PRIVATE: 'private',
  PUBLIC: 'public',
  PUBLISHED: 'published',
};

export const EVENT_TYPES = {
  ADDED: 'ADDED',
  DELETED: 'DELETED',
  MODIFIED: 'MODIFIED',
};

export const REGISTRY_ITEM_STATES = {
  ENABLED: 'enabled',
  IGNORED: 'ignored',
};

export const REGISTRY_ITEM_STATES_TO_TRANSITION_MAP = {
  [REGISTRY_ITEM_STATES.ENABLED]: 'Enable',
  [REGISTRY_ITEM_STATES.IGNORED]: 'Ignore',
};

// Mirrors the ConnectionStatusValue enum in meshery/schemas
// (schemas/constructs/v1beta1/connection/api.yml), which is the source of
// truth for these wire values. TODO: type-check this object's values against
// the schemas-exported `v1beta1.ConnectionStatusValue` TS type once
// @meshery/schemas publishes it (it isn't in the published package yet).
export const CONNECTION_STATES = {
  DISCOVERED: 'discovered',
  REGISTERED: 'registered',
  CONNECTED: 'connected',
  IGNORED: 'ignored',
  MAINTENANCE: 'maintenance',
  DISCONNECTED: 'disconnected',
  DELETED: 'deleted',
  NOTFOUND: 'not found',
};

export const CONTROLLERS = {
  BROKER: 'BROKER',
  OPERATOR: 'OPERATOR',
  MESHSYNC: 'MESHSYNC',
};

// DEPLOYED..CONNECTED mirror the ControllerStatusValue enum newly formalized
// in meshery/schemas (schemas/constructs/v1beta1/system/api.yml,
// meshery/schemas#1064), itself mirroring the MesheryControllerStatus
// GraphQL enum (server/internal/graphql/schema/schema.graphql) - these are
// real wire values, including the published UNKOWN spelling. DISABLED and
// UNKNOWN (correct spelling) below are UI-only sentinel states with no wire
// equivalent (see ui/utils/hooks/useKubernetesHook.tsx) and do not belong in
// the schema. TODO: once @meshery/schemas publishes the new enum, type-check
// the wire-backed subset here against the schemas-exported
// `v1beta1.ControllerStatusValue` TS type.
export const CONTROLLER_STATES = {
  DEPLOYED: 'DEPLOYED',
  NOTDEPLOYED: 'NOTDEPLOYED',
  UNDEPLOYED: 'UNDEPLOYED',
  DEPLOYING: 'DEPLOYING',
  ENABLED: 'ENABLED',
  UNKOWN: 'UNKOWN',
  RUNNING: 'RUNNING',
  CONNECTED: 'CONNECTED',
  DISABLED: 'DISABLED',
  UNKNOWN: 'UNKNOWN',
};

export const MesheryPatternsCatalog = 'meshery-patterns-catalog';

export const MesheryFiltersCatalog = 'meshery-filters-catalog';

// Remove this fetch all connections initially
export const CONNECTION_KINDS_DEF = ['MESHERY', 'KUBERNETES', 'PROMETHEUS', 'GRAFANA', 'GITHUB'];

export const CONNECTION_KINDS = {
  MESHERY: 'meshery',
  KUBERNETES: 'kubernetes',
  PROMETHEUS: 'prometheus',
  GRAFANA: 'grafana',
  GITHUB: 'github',
};

export const MESHSYNC_DEPLOYMENT_TYPE = {
  OPERATOR: 'operator',
  EMBEDDED: 'embedded',
};

export const MESHSYNC_STATES = {
  DISCOVERED: 'discovered',
  REGISTER: 'register',
};

export const TRANSFER_COMPONENT = {
  CHIP: 'chip',
  OTHER: 'other',
};

export const CONNECTION_STATE_TO_TRANSITION_MAP = {
  [CONNECTION_STATES.IGNORED]: 'Ignore',
  [CONNECTION_STATES.CONNECTED]: 'Connect',
  [CONNECTION_STATES.REGISTERED]: 'Register',
  [CONNECTION_STATES.DISCOVERED]: 'Discover',
  [CONNECTION_STATES.DELETED]: 'Delete',
  [CONNECTION_STATES.MAINTENANCE]: 'Maintenance',
  [CONNECTION_STATES.DISCONNECTED]: 'Disconnect',
  [CONNECTION_STATES.NOTFOUND]: 'Not Found',
};

export const FILTER = 'filter';
export const PATTERN = 'pattern';

// Meshery Extension Point
// Add your UI plugin into this extension point.
// Learn more: https://docs.meshery.io/extensibility/ui
export const EXTENSION_NAMES = {
  EXTENSION: 'kanvas',
};

export const EXTENSIONS = {
  [EXTENSION_NAMES.EXTENSION]: {
    name: 'Extension',
    signup_header: 'Visual design and operation',
    signup_button: 'Open Extension',
    signup_url: '/extension/meshmap',
    show_popup: true,
  },
  Catalog: {
    name: 'Meshery Catalog',
  },
  Snapshot: {
    name: 'Snapshot',
    signup_url: 'https://meshery.io/extensions/github-action-meshery-snapshot',
  },
  PerformanceAnalysis: {
    name: 'Performance Analysis',
  },
};
export const RESOURCE_TYPE = {
  FILTER: 'filter',
  DESIGN: 'design',
  CATALOG: 'catalog',
  VIEW: 'view',
};

export const APP_MODE = {
  DESIGN: 'design',
  OPERATOR: 'operator',
};

export const VIEW_VISIBILITY = {
  PUBLIC: 'public',
  PRIVATE: 'private',
};

// Meshery Extension Point
// ---
// Purpose: Notify Remote Providers of changes in Golang dependencies
// Learn more: See https://docs.meshery.io/extensibility
// Add your repository to the list: https://github.com/meshery/meshery/issues/new/choose
// ---
