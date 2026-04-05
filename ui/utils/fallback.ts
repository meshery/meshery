export function getFallbackImageBasedOnKind(kind) {
  const fallbackComponent = {
    meshery: 'static/img/meshery-logo.png',
    kubernetes: 'static/img/kubernetes.svg',
  };
  return fallbackComponent[kind];
}

export function normalizeStaticImagePath(path) {
  if (!path) {
    return '';
  }

  if (path.startsWith('http')) {
    return path;
  }

  const trimmed = path.replace(/^\/?ui\/public\//, '').replace(/^\/+/, '');
  return `/${trimmed}`;
}
