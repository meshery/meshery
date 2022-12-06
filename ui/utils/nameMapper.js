export const podNameMapper = (serviceMeshName, podName) => {
  switch (serviceMeshName) {
    case "istio":
      if (podName.includes("istio-galley")) return "Istio Galley";
      if (podName.includes("istio-ingressgateway")) return "Istio Ingress Gateway";
      if (podName.includes("istio-egressgateway")) return "Istio Egress Gateway";
      if (podName.includes("istio-policy")) return "Istio Policy";
      if (podName.includes("istio-citadel")) return "Istio Citadel";
      if (podName.includes("istio-telemetry")) return "Istio Telemetry";
      if (podName.includes("istio-pilot")) return "Istio Pilot";
      if (podName.includes("istio-tracing")) return "Istio Tracing";
      if (podName.includes("istio-sidecar-injector")) return "Istio Sidecar Injector";
      break;
    case "kiali":
      return "Kiali";
    case "grafana":
      return "Grafana";
    case "prometheus":
      return "Prometheus";
    case "osm":
      if (podName.includes("osm-injector")) return "OSM Injector";
      if (podName.includes("osm-controller")) return "OSM Controller";
      break;
    default:
      break;
  }
  const podNameArr = podName.split("-");

  if (Array.isArray(podNameArr)) {
    const prettifiedPodName = podNameArr
      .slice(0, -2)
      .map((word) => word[0].toUpperCase() + word.substr(1))
      .join(" ");

    return prettifiedPodName;
  }
  return podName;
};

export const versionMapper = (versionName) => (versionName.charAt(0) === "v" ? versionName : `v${versionName}`);



export function getMeshProperties(name) {
  switch (name) {
    case "istio": return { name, img : "/static/img/istio.svg", color : "#466BB0" }
    case "linkerd": return { name, img : "/static/img/linkerd.svg", color : "#2beda7" }
    case "kuma": return { name, img : "/static/img/kuma.svg", color : "#291953" }
    case "nginx_service_mesh": return { name, img : "/static/img/nginx.svg", color : "#009639" }
    case "open_service_mesh": return { name, img : "/static/img/openservicemesh.svg", color : "#6ED3B9" }
    case "cilium_service_mesh": return { name, img : "/static/img/cilium_service_mesh.svg", color : "#6B91C7" }
    case "citrix": return { name, img : "/static/img/citrix_service_mesh.svg", color : "#466BB0" }
    case "traefik_mesh": return { name, img : "/static/img/traefik_mesh.svg", color : "#9D0FB0" }
    case "consul": return { name, img : "/static/img/consul.svg", color : "#D62783" }
    case "app_mesh": return { name, img : "/static/img/app_mesh.svg", color : "#F49322" }
    case "core": return { name, img : "/static/img/meshery-logo.png", color : "#00B39F" }
    case "kubernetes": return { name, img : "/static/img/kubernetes.svg", color : "#293CDA" }
    default: return {}
  }
}