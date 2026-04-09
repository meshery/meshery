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

  const trimmed = path.trim();
  if (!trimmed || ['empty', 'none', 'null', 'undefined'].includes(trimmed.toLowerCase())) {
    return '';
  }

  if (trimmed.startsWith('http') || trimmed.startsWith('data:') || trimmed.startsWith('blob:')) {
    return trimmed;
  }

  const normalized = trimmed.replace(/^\/+/, '');

  // Meshmodel SVGs are generated under ui/public at runtime and are served by Meshery Server
  // through the same ui/public-prefixed path.
  if (normalized.startsWith('ui/public/static/img/meshmodels/')) {
    return `/${normalized}`;
  }

  return `/${normalized.replace(/^ui\/public\//, '')}`;
}
