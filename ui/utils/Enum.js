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

export const EXTENSIONS = {
  MESHMAP: 'meshmap',
};

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

// Fetch from GraphQL/REST API remove this
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
