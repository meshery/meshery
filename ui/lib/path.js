export function getPath(){
  let path = (typeof window !== 'undefined' ? window.location.pathname : '');
  path += (typeof window !== 'undefined' ? window.location.search : '');
  return path;
}