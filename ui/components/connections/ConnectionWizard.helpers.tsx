import { CONNECTION_KINDS } from '@/utils/Enum';
import { EVENT_TYPES } from 'lib/event-types';

/*
 * Pure, presentation-free helpers for the Create Connection wizard: kind-config
 * derivation from the registry, credential shaping, kubeconfig-import parsing,
 * step labels, and success notifications. No JSX or React state lives here so it
 * stays trivially unit-testable (see ConnectionWizard.helpers.test.ts).
 */

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

/**
 * Connection-kind identifier as it appears on a registry connection definition
 * (e.g. "kubernetes", "prometheus"). Intentionally a widened `string`: the set
 * of registerable kinds is registry-driven, not a fixed enum.
 */
export type SupportedConnectionWizardKind = string;

/** A JSON Schema object carried on a connection/credential definition. */
export type JsonSchema = Record<string, unknown>;

/** Loosely-typed form/record payload exchanged with the wizard steps. */
type GenericRecord = Record<string, unknown>;

export type ConnectionWizardFlow = 'kubernetes' | 'generic';

// ---------------------------------------------------------------------------
// Connections page navigation + create-connection deep link
// ---------------------------------------------------------------------------

/** Lifecycle Connections page path (deep links + post-create navigation). */
export const CONNECTIONS_PATH = '/management/connections';

export const CREATE_CONNECTION_QUERY = {
  create: 'create',
  kind: 'kind',
} as const;

export const isCreateConnectionQuery = (value: string | string[] | undefined): boolean =>
  value === 'true' || value === '1';

// ---------------------------------------------------------------------------
// Success notifications
// ---------------------------------------------------------------------------

export type ConnectionCreatedNotifyPayload = {
  message: string;
  event_type: typeof EVENT_TYPES.SUCCESS | typeof EVENT_TYPES.WARNING;
  link?: { href: string; label: string };
};

const isOnConnectionsPage = (): boolean =>
  typeof window !== 'undefined' && window.location.pathname.startsWith(CONNECTIONS_PATH);

/**
 * Success snackbar after create/import. Plain string (BasicMarkdown-safe) plus an
 * optional same-tab action when not already on the Connections page.
 */
export const connectionCreatedNotify = (label: string): ConnectionCreatedNotifyPayload => {
  const name = (label && String(label).trim()) || '';
  const summary = name ? `${name} connection created.` : 'Connection created.';
  if (isOnConnectionsPage()) {
    return { message: summary, event_type: EVENT_TYPES.SUCCESS };
  }
  return {
    message: summary,
    event_type: EVENT_TYPES.SUCCESS,
    link: { href: CONNECTIONS_PATH, label: 'View connections' },
  };
};

export const kubernetesImportedNotify = (count: number): ConnectionCreatedNotifyPayload => {
  const noun = count === 1 ? 'connection' : 'connections';
  const summary = `Imported ${count} Kubernetes ${noun}.`;
  const event_type = count > 0 ? EVENT_TYPES.SUCCESS : EVENT_TYPES.WARNING;
  if (isOnConnectionsPage() || count === 0) {
    return { message: summary, event_type };
  }
  return {
    message: summary,
    event_type,
    link: { href: CONNECTIONS_PATH, label: 'View connections' },
  };
};

// ---------------------------------------------------------------------------
// Registry connection definitions -> wizard kind configs
// ---------------------------------------------------------------------------

export type ConnectionWizardKindConfig = {
  kind: SupportedConnectionWizardKind;
  // type/subType identify the connection together with kind; extensions can
  // target a specific kind-type-subType combination.
  type: string;
  subType: string;
  label: string;
  description: string;
  flow: ConnectionWizardFlow;
  docsUrl: string;
  // JSON Schemas carried on the connection definition; the wizard renders these
  // directly instead of fetching them from the registration state machine.
  connectionSchema: JsonSchema | null;
  credentialSchema: JsonSchema | null;
  // Inline SVG markup from the connection definition's styles, used for the
  // kind's icon (white variant for dark backgrounds).
  svgColor: string | null;
  svgWhite: string | null;
};

export type ConnectionDefinitionStyles = {
  svgColor?: string;
  svgWhite?: string;
  svgComplete?: string;
};

/**
 * Shape of a connection definition as returned by the registry endpoint
 * `GET /api/meshmodels/connections`. Only the fields the wizard consumes are
 * declared; the payload carries more (schemaVersion, status, model, ...).
 */
export type ConnectionDefinition = {
  kind?: string;
  type?: string;
  subType?: string;
  name?: string;
  description?: string;
  connectionSchema?: JsonSchema | null;
  credentialSchema?: JsonSchema | null;
  styles?: ConnectionDefinitionStyles | null;
  metadata?: {
    flow?: ConnectionWizardFlow;
    docsURL?: string;
  } & Record<string, unknown>;
};

export const DEFAULT_CONNECTION_DOCS_URL = 'https://docs.meshery.io/concepts/logical/connections';

/** Treat empty/non-object schema payloads as absent so steps can branch on null. */
const asSchema = (value?: JsonSchema | null): JsonSchema | null =>
  value && typeof value === 'object' && Object.keys(value).length > 0 ? value : null;

/**
 * Determines which wizard flow a connection kind drives. Kubernetes is unique:
 * it imports clusters from a kubeconfig file rather than rendering a generic
 * registration form. A definition can override this via `metadata.flow`.
 */
export const resolveConnectionWizardFlow = (
  definition: ConnectionDefinition,
): ConnectionWizardFlow => {
  if (definition.metadata?.flow) {
    return definition.metadata.flow;
  }

  return definition.kind === CONNECTION_KINDS.KUBERNETES ? 'kubernetes' : 'generic';
};

/**
 * Maps the connection definitions returned by the registry endpoint into the
 * config the wizard renders. Replaces the previously hardcoded kind list so the
 * set of registerable connections is driven by what is registered in the
 * registry. Definitions without a `kind` are skipped and duplicates (by
 * kind|type|subType) are collapsed to the first occurrence.
 */
export const buildConnectionWizardKindConfigs = (
  definitions?: ConnectionDefinition[] | null,
): ConnectionWizardKindConfig[] => {
  if (!Array.isArray(definitions)) {
    return [];
  }

  const seen = new Set<string>();

  return definitions.reduce<ConnectionWizardKindConfig[]>((configs, definition) => {
    const kind = definition?.kind;
    if (!kind) {
      return configs;
    }

    const dedupeKey = [kind, definition.type || '', definition.subType || ''].join('|');
    if (seen.has(dedupeKey)) {
      return configs;
    }
    seen.add(dedupeKey);
    configs.push({
      kind,
      type: definition.type || '',
      subType: definition.subType || '',
      label: definition.name || kind,
      description: definition.description || '',
      flow: resolveConnectionWizardFlow(definition),
      docsUrl:
        (typeof definition.metadata?.docsURL === 'string' && definition.metadata.docsURL) ||
        DEFAULT_CONNECTION_DOCS_URL,
      connectionSchema: asSchema(definition.connectionSchema),
      credentialSchema: asSchema(definition.credentialSchema),
      svgColor: definition.styles?.svgColor || null,
      svgWhite: definition.styles?.svgWhite || null,
    });

    return configs;
  }, []);
};

// ---------------------------------------------------------------------------
// Credential shaping
// ---------------------------------------------------------------------------

export type CredentialRecord = {
  id?: string;
  name?: string;
  type?: string;
  kind?: string;
  secret?: Record<string, unknown>;
};

/**
 * The `credentialSecret` payload handed to the registration state machine. It is
 * either the forwarded fields of an existing credential or an arbitrary
 * normalized form record, so it keeps the open `GenericRecord` index while
 * naming the fields callers commonly read/stamp.
 */
export type CredentialSecretPayload = GenericRecord & {
  id?: string;
  name?: string;
  secret?: unknown;
};

export const filterCredentialsForKind = (
  credentials: CredentialRecord[],
  kind?: string | null,
): CredentialRecord[] => {
  if (!kind) {
    return [];
  }

  const normalizedKind = kind.toLowerCase();

  return credentials.filter((credential) => {
    const credentialKind = credential.type || credential.kind;
    return typeof credentialKind === 'string' && credentialKind.toLowerCase() === normalizedKind;
  });
};

export const normalizeCredentialPayload = (formData?: GenericRecord | null): GenericRecord => {
  if (!formData) {
    return {};
  }

  const credentialName = formData.name || formData.credentialName;
  const normalized = { ...formData };

  delete normalized.credentialName;

  if (credentialName && typeof credentialName === 'string') {
    normalized.name = credentialName;
  }

  return normalized;
};

/**
 * Builds the `credentialSecret` payload the registration state machine expects.
 *
 * For an existing credential we must forward the stored secret (nested under
 * `secret.secret`) alongside the id and name, otherwise the backend `register`
 * (verify) step rehydrates an empty `PromCred`/`GrafanaCred` and verification
 * fails for any auth-protected endpoint. For a new credential we pass the
 * normalized form payload, which the backend persists verbatim.
 */
export const buildCredentialSecret = (
  selectedCredential?: CredentialRecord | null,
  credentialFormData?: GenericRecord | null,
): CredentialSecretPayload => {
  if (selectedCredential) {
    return {
      id: selectedCredential.id,
      name: selectedCredential.name,
      secret: selectedCredential.secret?.secret,
    };
  }

  return normalizeCredentialPayload(credentialFormData);
};

// ---------------------------------------------------------------------------
// Wizard step labels
// ---------------------------------------------------------------------------

export const getWizardStepLabels = ({
  kind,
  flow,
  hasCredentialSchema,
}: {
  kind?: SupportedConnectionWizardKind | null;
  flow?: ConnectionWizardKindConfig['flow'];
  hasCredentialSchema?: boolean;
}): string[] => {
  const configureLabel = flow === 'kubernetes' ? 'Import Kubeconfig' : 'Configure Connection';

  const steps = ['Choose Kind', configureLabel];

  if (flow === 'generic' && hasCredentialSchema) {
    steps.push('Associate Credential');
  }

  steps.push(kind === CONNECTION_KINDS.KUBERNETES ? 'Review Import' : 'Review Connection');

  return steps;
};

// ---------------------------------------------------------------------------
// Kubernetes kubeconfig-import result parsing
// ---------------------------------------------------------------------------

export type KubernetesContextBucket = 'connected' | 'registered' | 'ignored' | 'errored';

// The kubeconfig upload endpoint (POST /system/kubernetes) returns its context
// buckets in camelCase (registeredContexts, ...); older/mocked payloads used
// snake_case, so accept either form. Keyed once at module scope rather than
// rebuilt per call.
const KUBERNETES_CONTEXT_KEYS: Record<KubernetesContextBucket, { camel: string; snake: string }> = {
  connected: { camel: 'connectedContexts', snake: 'connected_contexts' },
  registered: { camel: 'registeredContexts', snake: 'registered_contexts' },
  ignored: { camel: 'ignoredContexts', snake: 'ignored_contexts' },
  errored: { camel: 'erroredContexts', snake: 'errored_contexts' },
};

export const getKubernetesContexts = (
  response: GenericRecord | null | undefined,
  bucket: KubernetesContextBucket,
): GenericRecord[] => {
  const { camel, snake } = KUBERNETES_CONTEXT_KEYS[bucket];
  const value = response?.[camel] ?? response?.[snake];
  return Array.isArray(value) ? (value as GenericRecord[]) : [];
};

export type KubernetesImportSummary = {
  connectedCount: number;
  registeredCount: number;
  ignoredCount: number;
  erroredCount: number;
  importedCount: number;
};

export const buildKubernetesImportSummary = (
  response?: GenericRecord | null,
): KubernetesImportSummary => {
  const connectedCount = getKubernetesContexts(response, 'connected').length;
  const registeredCount = getKubernetesContexts(response, 'registered').length;
  const ignoredCount = getKubernetesContexts(response, 'ignored').length;
  const erroredCount = getKubernetesContexts(response, 'errored').length;

  return {
    connectedCount,
    registeredCount,
    ignoredCount,
    erroredCount,
    importedCount: connectedCount + registeredCount + ignoredCount,
  };
};

// ---------------------------------------------------------------------------
// Connection naming
// ---------------------------------------------------------------------------

export const resolveConnectionName = (kind: string, formData?: GenericRecord | null): string => {
  const explicitName = formData?.name;
  if (typeof explicitName === 'string' && explicitName.trim().length > 0) {
    return explicitName.trim();
  }

  const url = formData?.url;
  if (typeof url === 'string' && url.trim().length > 0) {
    return url.trim();
  }

  return `${kind}-connection`;
};
