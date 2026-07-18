import { CONNECTION_KINDS } from './Enum';

/**
 * Last-resort static icons by connection/registrant kind.
 *
 * Prefer dynamic sources first (callers already do this):
 * 1. `connection.kindLogo` / model styles / definition SVG
 * 2. `connectionMetadataState[kind].icon` (server connection definitions)
 * 3. only then `getFallbackImageBasedOnKind`
 *
 * Keys reuse CONNECTION_KINDS values so we do not invent a parallel kind vocabulary.
 * Paths are relative to the UI public root (`normalizeStaticImagePath` prefixes "/").
 */
const KIND_FALLBACK_ICONS = {
  [CONNECTION_KINDS.MESHERY]: 'static/img/meshery-logo/meshery-logo.png',
  [CONNECTION_KINDS.KUBERNETES]: 'static/img/integrations/kubernetes.svg',
  [CONNECTION_KINDS.GITHUB]: 'static/img/extensions/github.svg',
  [CONNECTION_KINDS.PROMETHEUS]: 'static/img/integrations/prometheus_logo_orange_circle.svg',
  [CONNECTION_KINDS.GRAFANA]: 'static/img/integrations/grafana_icon.svg',
  // Registry registrants (not CONNECTION_KINDS, but share this helper)
  artifacthub: 'static/img/integrations/artifacthub.svg',
  helm: 'static/img/extensions/helm_chart.svg',
};

/**
 * Resolve a well-known kind to a static fallback image path, or undefined.
 * Case-insensitive exact key match only — kinds on the wire already use
 * CONNECTION_KINDS values (e.g. "prometheus", not "Prometheus" / "prometheus_").
 */
export function getFallbackImageBasedOnKind(kind) {
  if (kind == null || kind === '') {
    return undefined;
  }
  return KIND_FALLBACK_ICONS[String(kind).trim().toLowerCase()];
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
