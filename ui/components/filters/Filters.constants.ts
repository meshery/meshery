export const ACTION_TYPES = {
  FETCH_FILTERS: {
    name: 'FETCH_FILTERS',
    error_msg: 'Failed to fetch filter',
  },
  DELETE_FILTERS: {
    name: 'DELETE_FILTERS',
    error_msg: 'Failed to delete filter file',
  },
  DEPLOY_FILTERS: {
    name: 'DEPLOY_FILTERS',
    error_msg: 'Failed to deploy filter file',
  },
  UNDEPLOY_FILTERS: {
    name: 'UNDEPLOY_FILTERS',
    error_msg: 'Failed to undeploy filter file',
  },
  UPLOAD_FILTERS: {
    name: 'UPLOAD_FILTERS',
    error_msg: 'Failed to upload filter file',
  },
  CLONE_FILTERS: {
    name: 'CLONE_FILTER',
    error_msg: 'Failed to clone filter file',
  },
  PUBLISH_CATALOG: {
    name: 'PUBLISH_CATALOG',
    error_msg: 'Failed to publish catalog',
  },
  UNPUBLISH_CATALOG: {
    name: 'PUBLISH_CATALOG',
    error_msg: 'Failed to publish catalog',
  },
  SCHEMA_FETCH: {
    name: 'SCHEMA_FETCH',
    error_msg: 'failed to fetch import schema',
  },
};

export const COLUMN_VIEWS: [string, string][] = [
  ['name', 'xs'],
  ['created_at', 'm'],
  ['updated_at', 'l'],
  ['visibility', 's'],
  ['Actions', 'xs'],
];
