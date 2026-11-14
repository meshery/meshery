/**
 * normalizeURI takes in a uri and adds "/"
 * to the start of it if it doesn't have
 * one already.
 *
 * Returns an empty string if uri is null, undefined, not a string, or empty.
 *
 * The goal is to safely prepend a leading slash to paths
 * without throwing errors on invalid or absent values.
 *
 * @param {unknown} [uri]
 * @returns {string}
 */
export default function normalizeURI(uri?: unknown): string {
  if (typeof uri !== 'string' || uri.length === 0) {
    return '';
  }

  if (uri.startsWith('/')) {
    return uri;
  }

  return '/' + uri;
}
