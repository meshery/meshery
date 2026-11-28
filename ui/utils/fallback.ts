export function getFallbackImageBasedOnKind(kind?: string) {
  if (typeof kind !== 'string') {
    return undefined;
  }
  const fallbackComponent: Record<string, string> = {
    meshery: 'static/img/meshery-logo/meshery-logo.png',
    kubernetes: 'static/img/integrations/kubernetes.svg',
  };
  const normalized = kind.trim().toLowerCase();
  if (Object.prototype.hasOwnProperty.call(fallbackComponent, normalized)) {
    return fallbackComponent[normalized];
  }
  return undefined;
}

export function normalizeStaticImagePath(path?: unknown) {
  if (typeof path !== 'string') {
    return '';
  }

  const trimmed = path.trim();
  if (!trimmed || ['empty', 'none', 'null', 'undefined'].includes(trimmed.toLowerCase())) {
    return '';
  }

  // Inline SVG markup (e.g. a connection/component definition's styles.svgColor)
  // is encoded as a data URI so it can be used directly as an <img> src.
  const lower = trimmed.toLowerCase();
  if (lower.startsWith('<svg') || lower.startsWith('<?xml')) {
    return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(trimmed)}`;
  }

  if (trimmed.startsWith('http') || trimmed.startsWith('data:') || trimmed.startsWith('blob:')) {
    return trimmed;
  }

  const normalized = trimmed.replace(/^\/+/, '');

  // Model component icon SVGs are generated under ui/public at runtime and are served by Meshery Server
  // through the same ui/public-prefixed path.
  if (normalized.startsWith('ui/public/static/img/meshmodels/')) {
    return `/${normalized}`;
  }

  return `/${normalized.replace(/^ui\/public\//, '')}`;
}
