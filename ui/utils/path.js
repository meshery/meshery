/**
 * @returns {string} path - Returns current location pathname - Not the full path but the first part of the path string separated by "/"
 */
export function getPath() {
  let path = typeof window !== "undefined" ? window.location.pathname : "";
  if (path.lastIndexOf("/") > 0) {
    path = path.substring(0, path.lastIndexOf("/"));
  }
  path += typeof window !== "undefined" ? window.location.search : "";
  return path;
}

export const pathToPageTitleMapper = {
  "/": "Dashboard",
  "/configuration/applications": "Applications",
  "/configuration/patterns": "Patterns",
  "/configuration/filters": "Filters",
  "/management/consul": "Consul",
  "/management/istio": "Istio",
  "/management/citrix": "Citrix Service Mesh",
  "/management/kuma": "Kuma",
  "/management/linkerd": "Linkerd",
  "/management/nsm": "Network Service Mesh",
  "/management/octarine": "Octarine",
  "/management/nginx": "NGINX Service Mesh",
  "/management/osm": "Open Service Mesh",
  "/management/traefik-mesh": "Traefik Mesh",
  "/performance": "Performance",
  "/smi_results": "Conformance",
  "/system/connections": "Connection Wizard",
  "/settings": "Settings",
};
