import { describe, expect, it, vi } from 'vitest';
import { CONNECTION_KINDS } from '@/utils/Enum';
import {
  buildConnectionWizardKindConfigs,
  buildCredentialSecret,
  buildKubernetesImportSummary,
  connectionCreatedNotify,
  DEFAULT_CONNECTION_DOCS_URL,
  filterCredentialsForKind,
  getWizardStepLabels,
  isCreateConnectionQuery,
  kubernetesImportedNotify,
  normalizeCredentialPayload,
  resolveConnectionName,
} from './ConnectionWizard.helpers';

describe('ConnectionWizard.helpers', () => {
  it('filters credentials by the selected connection kind', () => {
    const credentials = [
      { id: '1', name: 'prom-a', type: CONNECTION_KINDS.PROMETHEUS },
      { id: '2', name: 'graf-a', type: CONNECTION_KINDS.GRAFANA },
      { id: '3', name: 'prom-b', kind: CONNECTION_KINDS.PROMETHEUS },
    ];

    expect(filterCredentialsForKind(credentials, CONNECTION_KINDS.PROMETHEUS)).toEqual([
      credentials[0],
      credentials[2],
    ]);
    expect(filterCredentialsForKind(credentials, CONNECTION_KINDS.KUBERNETES)).toEqual([]);
  });

  it('normalizes credentialName into the canonical name property', () => {
    expect(
      normalizeCredentialPayload({
        credentialName: 'grafana-token',
        secret: { grafanaAPIKey: 'abc' },
      }),
    ).toEqual({
      name: 'grafana-token',
      secret: { grafanaAPIKey: 'abc' },
    });
  });

  it('forwards the stored secret when associating an existing credential', () => {
    expect(
      buildCredentialSecret(
        { id: 'cred-1', name: 'prom-token', secret: { secret: 'super-secret' } },
        { name: 'unused' },
      ),
    ).toEqual({
      id: 'cred-1',
      name: 'prom-token',
      secret: 'super-secret',
    });
  });

  it('normalizes the new-credential form payload when no credential is selected', () => {
    expect(
      buildCredentialSecret(null, {
        credentialName: 'grafana-token',
        secret: { grafanaAPIKey: 'abc' },
      }),
    ).toEqual({
      name: 'grafana-token',
      secret: { grafanaAPIKey: 'abc' },
    });
  });

  it('builds wizard kind configs from registry connection definitions', () => {
    const promConnectionSchema = {
      type: 'object',
      properties: { url: { type: 'string' } },
    };
    const promCredentialSchema = {
      type: 'object',
      properties: { secret: { type: 'string' } },
    };

    const configs = buildConnectionWizardKindConfigs([
      {
        kind: CONNECTION_KINDS.PROMETHEUS,
        type: 'telemetry',
        subType: 'metrics',
        name: 'Prometheus',
        description: 'Register a Prometheus endpoint.',
        connectionSchema: promConnectionSchema,
        credentialSchema: promCredentialSchema,
      },
      {
        kind: CONNECTION_KINDS.KUBERNETES,
        name: 'Kubernetes',
        description: 'Import clusters from a kubeconfig.',
        metadata: { docsURL: 'https://docs.meshery.io/installation/kubernetes' },
      },
    ]);

    expect(configs).toEqual([
      {
        kind: CONNECTION_KINDS.PROMETHEUS,
        type: 'telemetry',
        subType: 'metrics',
        label: 'Prometheus',
        description: 'Register a Prometheus endpoint.',
        flow: 'generic',
        docsUrl: DEFAULT_CONNECTION_DOCS_URL,
        connectionSchema: promConnectionSchema,
        credentialSchema: promCredentialSchema,
        svgColor: null,
        svgWhite: null,
      },
      {
        kind: CONNECTION_KINDS.KUBERNETES,
        type: '',
        subType: '',
        label: 'Kubernetes',
        description: 'Import clusters from a kubeconfig.',
        flow: 'kubernetes',
        docsUrl: 'https://docs.meshery.io/installation/kubernetes',
        connectionSchema: null,
        credentialSchema: null,
        svgColor: null,
        svgWhite: null,
      },
    ]);
  });

  it('extracts svgColor/svgWhite from the definition styles', () => {
    const [config] = buildConnectionWizardKindConfigs([
      {
        kind: CONNECTION_KINDS.KUBERNETES,
        name: 'Kubernetes',
        styles: { svgColor: '<svg>color</svg>', svgWhite: '<svg>white</svg>' },
      },
    ]);

    expect(config.svgColor).toBe('<svg>color</svg>');
    expect(config.svgWhite).toBe('<svg>white</svg>');
  });

  it('treats empty connection/credential schema objects as absent', () => {
    const [config] = buildConnectionWizardKindConfigs([
      {
        kind: CONNECTION_KINDS.PROMETHEUS,
        name: 'Prometheus',
        connectionSchema: {},
        credentialSchema: undefined,
      },
    ]);

    expect(config.connectionSchema).toBeNull();
    expect(config.credentialSchema).toBeNull();
  });

  it('skips definitions without a kind and de-duplicates by kind', () => {
    const configs = buildConnectionWizardKindConfigs([
      { name: 'Missing kind' },
      { kind: CONNECTION_KINDS.GRAFANA, name: 'Grafana' },
      { kind: CONNECTION_KINDS.GRAFANA, name: 'Grafana duplicate' },
    ]);

    expect(configs).toHaveLength(1);
    expect(configs[0].kind).toBe(CONNECTION_KINDS.GRAFANA);
    expect(configs[0].label).toBe('Grafana');
  });

  it('returns an empty list when no definitions are provided', () => {
    expect(buildConnectionWizardKindConfigs(undefined)).toEqual([]);
    expect(buildConnectionWizardKindConfigs(null)).toEqual([]);
  });

  it('builds the credential step only for generic flows with credential schemas', () => {
    expect(
      getWizardStepLabels({
        kind: CONNECTION_KINDS.GRAFANA,
        flow: 'generic',
        hasCredentialSchema: true,
      }),
    ).toEqual(['Choose Kind', 'Configure Connection', 'Associate Credential', 'Review Connection']);

    expect(
      getWizardStepLabels({
        kind: CONNECTION_KINDS.KUBERNETES,
        flow: 'kubernetes',
        hasCredentialSchema: false,
      }),
    ).toEqual(['Choose Kind', 'Import Kubeconfig', 'Review Import']);
  });

  it('summarizes imported kubernetes context counts', () => {
    expect(
      buildKubernetesImportSummary({
        connected_contexts: [{ id: '1' }],
        registered_contexts: [{ id: '2' }, { id: '3' }],
        ignored_contexts: [],
        errored_contexts: [{ id: '4' }],
      }),
    ).toEqual({
      connectedCount: 1,
      registeredCount: 2,
      ignoredCount: 0,
      erroredCount: 1,
      importedCount: 3,
    });
  });

  it('resolves a stable connection name from form data', () => {
    expect(resolveConnectionName(CONNECTION_KINDS.PROMETHEUS, { name: 'prod-prom' })).toBe(
      'prod-prom',
    );
    expect(
      resolveConnectionName(CONNECTION_KINDS.PROMETHEUS, { url: 'https://prom.example' }),
    ).toBe('https://prom.example');
    expect(resolveConnectionName(CONNECTION_KINDS.PROMETHEUS, {})).toBe('prometheus-connection');
  });

  it('recognizes create-connection query flags', () => {
    expect(isCreateConnectionQuery('true')).toBe(true);
    expect(isCreateConnectionQuery('1')).toBe(true);
    expect(isCreateConnectionQuery('false')).toBe(false);
  });

  it('builds connection-created notify payloads with optional View connections action', () => {
    vi.stubGlobal('window', { location: { pathname: '/dashboard' } });
    const payload = connectionCreatedNotify('Prometheus');
    expect(payload.message).toBe('Prometheus connection created.');
    expect(payload.link?.href).toBe('/management/connections');
    vi.stubGlobal('window', { location: { pathname: '/management/connections' } });
    expect(connectionCreatedNotify('Grafana').link).toBeUndefined();
    vi.unstubAllGlobals();
  });

  it('formats kubernetes import notify summaries', () => {
    vi.stubGlobal('window', { location: { pathname: '/dashboard' } });
    expect(kubernetesImportedNotify(2).message).toBe('Imported 2 Kubernetes connections.');
    vi.unstubAllGlobals();
  });
});

// ---------------------------------------------------------------------------
// Source connections (Artifact Hub + GitHub)
//
// These kinds ship purely as registry connection definitions (added in
// meshery/meshery#20926) with no bespoke wizard code: the wizard surfaces and
// drives them entirely through the generic, schema-driven path. The definition
// shapes below mirror the on-disk models
// (models/meshery-core/0.7.2/v1.0.0/connections/{ArtifactHubConnection,GitHubConnection}.json)
// so the tests fail if the wire contract the wizard depends on drifts.
// ---------------------------------------------------------------------------

const ARTIFACT_HUB_DEFINITION = {
  kind: 'artifacthub',
  type: 'source',
  subType: 'registry',
  name: 'Artifact Hub',
  description: 'An Artifact Hub instance that Meshery sources packages from.',
  connectionSchema: {
    type: 'object',
    title: 'Artifact Hub Connection',
    required: ['url'],
    properties: {
      url: { type: 'string', format: 'uri', title: 'Artifact Hub Endpoint' },
      name: { type: 'string', title: 'Connection Name' },
    },
  },
  credentialSchema: {
    type: 'object',
    title: 'Artifact Hub Credential',
    properties: {
      apiKeyId: { type: 'string', title: 'API Key ID' },
      apiKeySecret: { type: 'string', format: 'password', title: 'API Key Secret' },
    },
    dependencies: {
      apiKeyId: ['apiKeySecret'],
      apiKeySecret: ['apiKeyId'],
    },
  },
  styles: { svgColor: '<svg>ah-color</svg>', svgWhite: '<svg>ah-white</svg>' },
};

const GITHUB_DEFINITION = {
  kind: 'github',
  type: 'source',
  subType: 'git',
  name: 'GitHub',
  description: 'A GitHub repository that Meshery sources from.',
  connectionSchema: {
    type: 'object',
    title: 'GitHub Connection',
    required: ['url'],
    properties: {
      url: { type: 'string', format: 'uri', title: 'Repository URL' },
      branch: { type: 'string', title: 'Branch' },
      name: { type: 'string', title: 'Connection Name' },
    },
  },
  credentialSchema: {
    type: 'object',
    title: 'GitHub Credential',
    properties: {
      token: { type: 'string', format: 'password', title: 'Access Token' },
    },
  },
  styles: { svgColor: '<svg>gh-color</svg>', svgWhite: '<svg>gh-white</svg>' },
};

describe('ConnectionWizard.helpers — source connections (Artifact Hub + GitHub)', () => {
  it('surfaces Artifact Hub + GitHub from registry definitions via the generic flow', () => {
    const configs = buildConnectionWizardKindConfigs([ARTIFACT_HUB_DEFINITION, GITHUB_DEFINITION]);

    expect(configs).toEqual([
      {
        kind: 'artifacthub',
        type: 'source',
        subType: 'registry',
        label: 'Artifact Hub',
        description: 'An Artifact Hub instance that Meshery sources packages from.',
        // Non-kubernetes kinds drive the generic (form-based) flow, not the
        // kubeconfig-import flow.
        flow: 'generic',
        docsUrl: DEFAULT_CONNECTION_DOCS_URL,
        connectionSchema: ARTIFACT_HUB_DEFINITION.connectionSchema,
        credentialSchema: ARTIFACT_HUB_DEFINITION.credentialSchema,
        svgColor: '<svg>ah-color</svg>',
        svgWhite: '<svg>ah-white</svg>',
      },
      {
        kind: 'github',
        type: 'source',
        subType: 'git',
        label: 'GitHub',
        description: 'A GitHub repository that Meshery sources from.',
        flow: 'generic',
        docsUrl: DEFAULT_CONNECTION_DOCS_URL,
        connectionSchema: GITHUB_DEFINITION.connectionSchema,
        credentialSchema: GITHUB_DEFINITION.credentialSchema,
        svgColor: '<svg>gh-color</svg>',
        svgWhite: '<svg>gh-white</svg>',
      },
    ]);
  });

  // The wizard consumes these schemas as opaque JSON Schema (Record<string,
  // unknown>); this narrows to the fields the assertions read.
  const asSchemaShape = (schema: Record<string, unknown> | null) =>
    (schema ?? {}) as {
      properties?: Record<string, unknown>;
      required?: string[];
      dependencies?: Record<string, string[]>;
    };

  it('carries the endpoint/repo connectionSchema so the Configure step renders', () => {
    const [artifactHub, github] = buildConnectionWizardKindConfigs([
      ARTIFACT_HUB_DEFINITION,
      GITHUB_DEFINITION,
    ]);

    // The generic details step only advances when a connectionSchema is present
    // (genericDetailsStep.canProceed), so a non-null schema is what makes the
    // Configure Connection step usable for these kinds.
    const artifactHubSchema = asSchemaShape(artifactHub.connectionSchema);
    const githubSchema = asSchemaShape(github.connectionSchema);
    expect(Object.keys(artifactHubSchema.properties ?? {})).toContain('url');
    expect(artifactHubSchema.required).toEqual(['url']);
    expect(Object.keys(githubSchema.properties ?? {})).toEqual(['url', 'branch', 'name']);
    expect(githubSchema.required).toEqual(['url']);
  });

  it('carries the credentialSchema so the Associate Credential step is shown, not hidden', () => {
    const [artifactHub, github] = buildConnectionWizardKindConfigs([
      ARTIFACT_HUB_DEFINITION,
      GITHUB_DEFINITION,
    ]);

    // genericCredentialStep.hidden === !credentialSchema, so a non-null
    // credentialSchema keeps the credential step in the flow for both kinds.
    expect(artifactHub.credentialSchema).not.toBeNull();
    expect(github.credentialSchema).not.toBeNull();

    // Artifact Hub pairs the API key id + secret (both-or-neither); GitHub takes
    // a single optional access token.
    expect(asSchemaShape(artifactHub.credentialSchema).dependencies).toEqual({
      apiKeyId: ['apiKeySecret'],
      apiKeySecret: ['apiKeyId'],
    });
    expect(Object.keys(asSchemaShape(github.credentialSchema).properties ?? {})).toEqual(['token']);
  });

  it('lists the credential + review steps for a generic source connection', () => {
    // GitHub carries a credentialSchema, so the generic flow inserts the
    // Associate Credential step between Configure and Review.
    expect(
      getWizardStepLabels({ kind: 'github', flow: 'generic', hasCredentialSchema: true }),
    ).toEqual(['Choose Kind', 'Configure Connection', 'Associate Credential', 'Review Connection']);
  });

  it('names a source connection from its endpoint/repo URL when unnamed', () => {
    expect(resolveConnectionName('artifacthub', { url: 'https://artifacthub.io' })).toBe(
      'https://artifacthub.io',
    );
    expect(
      resolveConnectionName('github', {
        url: 'https://github.com/meshery/meshery',
        branch: 'master',
      }),
    ).toBe('https://github.com/meshery/meshery');
    // An explicit name always wins over the URL.
    expect(
      resolveConnectionName('github', { name: 'meshery-repo', url: 'https://github.com/x/y' }),
    ).toBe('meshery-repo');
    expect(resolveConnectionName('artifacthub', {})).toBe('artifacthub-connection');
  });

  it('filters existing credentials to the selected source connection kind', () => {
    const credentials = [
      { id: '1', name: 'ah-key', type: 'artifacthub' },
      { id: '2', name: 'gh-token', kind: 'github' },
      { id: '3', name: 'prom', type: 'prometheus' },
    ];

    expect(filterCredentialsForKind(credentials, 'artifacthub')).toEqual([credentials[0]]);
    expect(filterCredentialsForKind(credentials, 'github')).toEqual([credentials[1]]);
  });
});
