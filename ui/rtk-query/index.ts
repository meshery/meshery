export { mesheryApi as api } from '@meshery/schemas/mesheryApi';

const ABSOLUTE_URL_PATTERN = /^[a-z]+:\/\//i;
const MESHERY_API_PREFIX = '/api';
const MESHERY_EXTENSION_PREFIX = '/api/extensions';

export const mesheryApiPath = (path = '') => {
  if (!path) {
    return MESHERY_API_PREFIX;
  }

  if (ABSOLUTE_URL_PATTERN.test(path) || path.startsWith('//')) {
    return path;
  }

  if (path === MESHERY_EXTENSION_PREFIX || path.startsWith(`${MESHERY_EXTENSION_PREFIX}/`)) {
    return path;
  }

  if (path.startsWith('/extensions/')) {
    return `${MESHERY_API_PREFIX}${path}`;
  }

  if (path.startsWith('extensions/')) {
    return `${MESHERY_API_PREFIX}/${path}`;
  }

  if (path === MESHERY_API_PREFIX || path.startsWith(`${MESHERY_API_PREFIX}/`)) {
    return path;
  }

  if (path.startsWith('/evaluate')) {
    return path;
  }

  return `${MESHERY_API_PREFIX}/${path.replace(/^\/+/, '')}`;
};
