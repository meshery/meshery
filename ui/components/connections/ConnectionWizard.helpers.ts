import { CONNECTION_KINDS } from '@/utils/Enum';

export type SupportedConnectionWizardKind = string;

export type ConnectionWizardFlow = 'kubernetes' | 'generic';

export type JsonSchema = Record<string, unknown>;

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

const asSchema = (value?: JsonSchema | null): JsonSchema | null =>
  value && typeof value === 'object' && Object.keys(value).length > 0 ? value : null;

export const DEFAULT_CONNECTION_DOCS_URL = 'https://docs.meshery.io/concepts/logical/connections';

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
 * registry. Definitions without a `kind` are skipped.
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

export type CredentialRecord = {
  id?: string;
  name?: string;
  type?: string;
  kind?: string;
  secret?: Record<string, unknown>;
};

type GenericRecord = Record<string, unknown>;

export const filterCredentialsForKind = (credentials: CredentialRecord[], kind?: string | null) => {
  if (!kind) {
    return [];
  }

  const normalizedKind = kind.toLowerCase();

  return credentials.filter((credential) => {
    const credentialKind = credential.type || credential.kind;
    return typeof credentialKind === 'string' && credentialKind.toLowerCase() === normalizedKind;
  });
};

export const normalizeCredentialPayload = (formData?: GenericRecord | null) => {
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
) => {
  if (selectedCredential) {
    return {
      id: selectedCredential.id,
      name: selectedCredential.name,
      secret: selectedCredential.secret?.secret,
    };
  }

  return normalizeCredentialPayload(credentialFormData);
};

export const getWizardStepLabels = ({
  kind,
  flow,
  hasCredentialSchema,
}: {
  kind?: SupportedConnectionWizardKind | null;
  flow?: ConnectionWizardKindConfig['flow'];
  hasCredentialSchema?: boolean;
}) => {
  const configureLabel = flow === 'kubernetes' ? 'Import Kubeconfig' : 'Configure Connection';

  const steps = ['Choose Kind', configureLabel];

  if (flow === 'generic' && hasCredentialSchema) {
    steps.push('Associate Credential');
  }

  steps.push(kind === CONNECTION_KINDS.KUBERNETES ? 'Review Import' : 'Review Connection');

  return steps;
};

// The kubeconfig upload endpoint (POST /system/kubernetes) returns its context
// buckets in camelCase (registeredContexts, ...); older/mocked payloads used
// snake_case, so accept either form.
export const getKubernetesContexts = (
  response: GenericRecord | null | undefined,
  bucket: 'connected' | 'registered' | 'ignored' | 'errored',
): GenericRecord[] => {
  const camel = {
    connected: 'connectedContexts',
    registered: 'registeredContexts',
    ignored: 'ignoredContexts',
    errored: 'erroredContexts',
  }[bucket];
  const snake = {
    connected: 'connected_contexts',
    registered: 'registered_contexts',
    ignored: 'ignored_contexts',
    errored: 'errored_contexts',
  }[bucket];
  const value = response?.[camel] ?? response?.[snake];
  return Array.isArray(value) ? (value as GenericRecord[]) : [];
};

export const buildKubernetesImportSummary = (response?: GenericRecord | null) => {
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

export const resolveConnectionName = (kind: string, formData?: GenericRecord | null) => {
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
