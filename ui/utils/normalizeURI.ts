/**
 * normalizeURI takes in a uri and adds "/"
 * to the start of it if they they if it doesn't
 * them already
 *
 * The goal is to be able to easily append the uris
 * without concerning about the slashes
 * @param {string} uri
 * @returns {string}
 */
export default function normalizeURI(uri) {
  if (!uri.startsWith('/')) return '/' + uri;
  return uri;
}
