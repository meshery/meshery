export function getPath() {
  let path = typeof window !== 'undefined' ? window.location.pathname : '';
  if (path.lastIndexOf('/') > 0) {
    path = path.substring(0, path.lastIndexOf('/'));
  }
  path += typeof window !== 'undefined' ? window.location.search : '';
  return path;
}
