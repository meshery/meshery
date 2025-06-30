export function getFallbackImageBasedOnKind(kind) {
  const fallbackComponent = {
    meshery: 'static/img/meshery-logo.png',
    kubernetes: 'static/img/kubernetes.svg',
  };
  return fallbackComponent[kind];
}
