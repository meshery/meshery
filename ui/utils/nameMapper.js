export const podNameMapper = (serviceMeshName, podName) => {
  switch (serviceMeshName) {
    case 'istio':
      if (podName.includes('istio-galley')) return 'Istio Galley';
      if (podName.includes('istio-ingressgateway')) return 'Istio Ingress Gateway';
      if (podName.includes('istio-egressgateway')) return 'Istio Egress Gateway';
      if (podName.includes('istio-policy')) return 'Istio Policy';
      if (podName.includes('istio-citadel')) return 'Istio Citadel';
      if (podName.includes('istio-telemetry')) return 'Istio Telemetry';
      if (podName.includes('istio-pilot')) return 'Istio Pilot';
      if (podName.includes('istio-tracing')) return 'Istio Tracing';
      if (podName.includes('istio-sidecar-injector')) return 'Istio Sidecar Injector';
      break;
    case 'kiali':
      return 'Kiali';
    case 'grafana':
      return 'Grafana';
    case 'prometheus':
      return 'Prometheus';
    default:
      break;
  }
  const podNameArr = podName.split('-');

  if (Array.isArray(podNameArr)) {
    const prettifiedPodName = podNameArr
      .slice(0, -2)
      .map((word) => word[0].toUpperCase() + word.substr(1))
      .join(' ');

    return prettifiedPodName;
  }
  return podName;
};

export const versionMapper = (versionName) =>
  versionName.charAt(0) === 'v' ? versionName : `v${versionName}`;

export const TelemetryComps = {
  GRAFANA: 'grafana',
  PROMETHEUS: 'prometheus',
};

// maps objects to telemetry comp if name is one of "TelemetryComps"
export function isTelemetryComponent(name) {
  const comps = Object.values(TelemetryComps);
  for (const comp in comps) {
    if (name?.toLowerCase().includes(comps[comp])) {
      return comps[comp];
    }
  }
  return '';
}
