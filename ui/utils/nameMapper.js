export const podNameMapper = (serviceMeshName, podName) => {
  switch (serviceMeshName) {
    case "istio":
      if (podName.includes("istio-galley")) return "Istio Galley";
      else if (podName.includes("istio-ingressgateway")) return "Istio Ingress Gateway";
      else if (podName.includes("istio-egressgateway")) return "Istio Egress Gateway";
      else if (podName.includes("istio-policy")) return "Istio Policy";
      else if (podName.includes("istio-citadel")) return "Istio Citadel";
      else if (podName.includes("istio-telemetry")) return "Istio Telemetry";
      else if (podName.includes("istio-pilot")) return "Istio Pilot";
      else if (podName.includes("istio-tracing")) return "Istio Tracing";
      else if (podName.includes("istio-sidecar-injector")) return "Istio Sidecar Injector";
      break;
    case "kiali":
      return "Kiali";
    case "grafana":
      return "Grafana";
    case "prometheus":
      return "Prometheus";
    case "osm":
      if (podName.includes("osm-injector")) return "OSM Injector";
      else if (podName.includes("osm-controller")) return "OSM Controller";
      break;
    default:
      break;
  }
  return podName;
};

export const versionMapper = (versionName) => (versionName.charAt(0) === "v" ? versionName : `v${versionName}`);
