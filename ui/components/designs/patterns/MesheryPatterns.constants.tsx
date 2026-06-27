export const ACTION_TYPES = {
  FETCH_PATTERNS: {
    name: 'FETCH_PATTERNS',
    error_msg: 'Failed to fetch designs',
  },
  UPDATE_PATTERN: {
    name: 'UPDATE_PATTERN',
    error_msg: 'Failed to update design file',
  },
  DELETE_PATTERN: {
    name: 'DELETE_PATTERN',
    error_msg: 'Failed to delete design file',
  },
  DEPLOY_PATTERN: {
    name: 'DEPLOY_PATTERN',
    error_msg: 'Failed to deploy design file',
  },
  UNDEPLOY_PATTERN: {
    name: 'UNDEPLOY_PATTERN',
    error_msg: 'Failed to undeploy design file',
  },
  UPLOAD_PATTERN: {
    name: 'UPLOAD_PATTERN',
    error_msg: 'Failed to upload design file',
  },
  CLONE_PATTERN: {
    name: 'CLONE_PATTERN',
    error_msg: 'Failed to clone design file',
  },
  PUBLISH_CATALOG: {
    name: 'PUBLISH_CATALOG',
    error_msg: 'Failed to publish catalog',
  },
  UNPUBLISH_CATALOG: {
    name: 'UNPUBLISH_CATALOG',
    error_msg: 'Failed to unpublish catalog',
  },
  SCHEMA_FETCH: {
    name: 'SCHEMA_FETCH',
    error_msg: 'failed to fetch import schema',
  },
  EVALUATE_RELATIONSHIP: {
    name: 'EVALUATE_RELATIONSHIP',
    error_msg: 'Failed to evaluate design relationships',
  },
};

export const genericClickHandler = (ev, fn) => {
  ev.stopPropagation();
  fn(ev);
};

export function resetSelectedPattern() {
  return { show: false, pattern: null };
}
