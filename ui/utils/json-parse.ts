/**
 * Safely parses a JSON string into an object or array. If the input is already
 * an object, it is returned as-is. If parsing fails or the input is nullish,
 * the provided fallback value is returned instead.
 *
 * This utility prevents unhandled `SyntaxError` exceptions inside React render
 * functions (e.g. MUI-Datatables `customBodyRender`) where a thrown error would
 * crash the entire component tree.
 */
export function safeJsonParse<T extends object = Record<string, unknown>>(
  value: unknown,
  fallback: T = {} as T,
): T {
  if (value === null || value === undefined) {
    return fallback;
  }

  if (typeof value === 'object') {
    return value as T;
  }

  if (typeof value !== 'string') {
    return fallback;
  }

  try {
    const parsed = JSON.parse(value);
    return parsed !== null && typeof parsed === 'object' ? (parsed as T) : fallback;
  } catch {
    return fallback;
  }
}
