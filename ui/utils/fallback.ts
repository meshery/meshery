export function getFallbackImageBasedOnKind(kind) {
  // Last-resort icons when connection.kindLogo / definition icon is missing.
  // Connections table: kindLogo || getFallbackImageBasedOnKind(kind).
  const fallbackComponent = {
    meshery: 'static/img/meshery-logo/meshery-logo.png',
    kubernetes: 'static/img/integrations/kubernetes.svg',
    // Artifact Hub is a registrant/connection kind on Lifecycle → Connections;
    // without this entry the name chip falls through to a generic icon (#20864).
    artifacthub: 'static/img/integrations/artifacthub.svg',
  };
  if (kind == null || kind === '') {
    return undefined;
  }
  return fallbackComponent[String(kind).trim().toLowerCase()];
}

export function normalizeStaticImagePath(path) {
  if (!path) {
    return '';
  }

  const trimmed = path.trim();
  if (!trimmed || ['empty', 'none', 'null', 'undefined'].includes(trimmed.toLowerCase())) {
    return '';
  }

  // Inline SVG markup (e.g. a connection/component definition's styles.svgColor)
  // is encoded as a data URI so it can be used directly as an <img> src.
  if (trimmed.startsWith('<svg') || trimmed.startsWith('<?xml')) {
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
