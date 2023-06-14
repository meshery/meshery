/**
 * Finds the host on which Meshery is running, it could be
 * localhost:{port} or an IP adress or a webadress
 * @returns {string} application host
 */
function getHost() {
  return window.location.host;
}

/**
 * The protocol used by the Meshery Application.
 * It could be , "http:" or "https:"
 * @returns {string}
 */
export function getProtocol() {
  return window.location.protocol;
}

/**
 *
 * @returns {string} base url on which Meshery is runnning,
 * @example http://localhost:9081
 */
export function getWebAdress() {
  return getProtocol() + "//" + getHost();
}

/**
 * get value from query param
 * @example for "http://localhost:9081/extension/meshmap?componentInterface=meshmodel" ,
 * if called the function like getQueryParam("componentInterface"), then it will return "meshmodel"
 * @param {string} queryKey
 * @returns {string} queryVal
 */
export function getQueryParam(queryKey) {
  let queryParamString = window.location.search;
  queryParamString = queryParamString.replace("?", "");

  let queryVal = "";

  queryParamString.split("&").forEach(query => {
    if (query.split("=")[0] === queryKey) {
      if (!queryVal) {
        queryVal = query.split("=")[1];
      }
    }
  });

  return queryVal;
}

export function getRawUrlFromCssUrlString(url) {
  if (!url) return;

  // turns url(http://localhost:9081/path/to/svg) to http://localhost:9081/path/to/svg
  if (url.startsWith("url")) {
    url = url.slice(4).slice(0, -1);
  }

  // turns http://localhost:9081/path/to/svg to "/path/to/svg"
  if (url.startsWith("http")) {
    return url.split("//")?.[1]?.split("/")?.slice(1)?.join("/");
  }

  return url;
}
