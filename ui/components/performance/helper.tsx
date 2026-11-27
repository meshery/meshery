import { v4 as uuid } from 'uuid';
import { ctxUrl } from '../../utils/multi-ctx';

/**
 * generateTestName takes in test name and service mesh name
 * and generates a random name (if test name is an empty string or is falsy) or
 * will return the given name
 */
export const generateTestName = (name: string, meshName: string): string => {
  if (!name || name.trim() === '') {
    const mesh = meshName === '' || meshName === 'None' ? 'No mesh' : meshName;
    return `${mesh}_${new Date().getTime()}`;
  }

  return name;
};

export function generateUUID(): string {
  return uuid();
}

/**
 * Constructs the SSE run URL for a performance profile, properly joining
 * query parameters whether or not Kubernetes contexts are selected.
 */
export function buildPerformanceRunUrl(
  id: string,
  selectedK8sContexts?: string[] | null,
  params?: string,
): string {
  const baseURL = ctxUrl(`/api/user/performance/profiles/${id}/run`, selectedK8sContexts);
  const queryParts = ['cert=true'];
  if (params) {
    queryParts.push(params);
  }
  const queryString = queryParts.join('&');
  return `${baseURL}${baseURL.includes('?') ? '&' : '?'}${queryString}`;
}
