/** Add a leading slash to a non-empty URI, or return an empty string for invalid input. */
export default function normalizeURI(uri: unknown): string {
  if (typeof uri !== 'string' || uri.length === 0) return '';
  if (!uri.startsWith('/')) return '/' + uri;
  return uri;
}
