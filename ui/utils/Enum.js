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

export const MesheryPatternsCatalog = 'meshery-patterns-catalog';

export const MesheryFiltersCatalog = 'meshery-filters-catalog';

export const CONNECTION_KINDS = {
  MESHERY: 'meshery',
  KUBERNETES: 'kubernetes',
};
