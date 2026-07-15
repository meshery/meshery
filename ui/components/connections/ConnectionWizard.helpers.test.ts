import { describe, expect, it } from 'vitest';
import { CONNECTION_KINDS } from '@/utils/Enum';
import {
  buildConnectionWizardKindConfigs,
  buildCredentialSecret,
  buildKubernetesImportSummary,
  DEFAULT_CONNECTION_DOCS_URL,
  filterCredentialsForKind,
  getWizardStepLabels,
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
});
