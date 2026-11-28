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

// Only plain digit spellings pass. The run request sends the typed text as is,
// and the server reads concurrent requests with strconv.Atoi, which fails on
// an exponent such as "1e3" and falls back to 1 request instead of 1000.
// A digit string too long to hold exactly (or one that converts to Infinity)
// fails too, so the value saved and sent is the one the user typed.
const isWholeNumberAtLeast = (value: string | number | null | undefined, min: number): boolean => {
  const text = String(value).trim();
  if (!/^[+-]?\d+$/.test(text)) return false;
  const number = Number(text);
  return Number.isSafeInteger(number) && number >= min;
};

/**
 * validateLoadParams checks the load settings of the performance form.
 * Concurrent requests must be a whole number of at least 1; QPS must be a
 * whole number of at least 0, where 0 requests the maximum rate. Each field
 * gets an error message, or an empty string when the value is valid.
 */
export function validateLoadParams(
  c: string | number | null | undefined,
  qps: string | number | null | undefined,
): { cError: string; qpsError: string } {
  return {
    cError: isWholeNumberAtLeast(c, 1) ? '' : 'Enter a whole number of 1 or more',
    qpsError: isWholeNumberAtLeast(qps, 0) ? '' : 'Enter a whole number of 0 or more (0 = max QPS)',
  };
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
