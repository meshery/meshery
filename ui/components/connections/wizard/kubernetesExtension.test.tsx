import React, { createRef } from 'react';
import { configureStore } from '@reduxjs/toolkit';
import { act, render, renderHook, screen } from '@testing-library/react';
import { Provider } from 'react-redux';
import { SistentThemeProvider } from '@sistent/sistent';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { WizardContext, WizardFormRefs } from './types';
import type { ConnectionWizardKindConfig } from '../ConnectionWizard.helpers';

const discoverUnwrap = vi.fn();
const uploadUnwrap = vi.fn();
const updateUnwrap = vi.fn();
const discoverTrigger = vi.fn(() => ({ unwrap: discoverUnwrap }));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    useHasPermission: () => true,
  };
});

vi.mock('@/rtk-query/connection', () => ({
  useAddKubernetesConfigMutation: () => [() => ({ unwrap: uploadUnwrap })],
  useDiscoverKubernetesContextsMutation: () => [discoverTrigger],
  useGetCredentialsQuery: () => ({ data: { credentials: [] } }),
  useGetControllerDiagnosticsQuery: () => ({ data: undefined }),
  usePerformConnectionActionMutation: () => [() => ({ unwrap: () => Promise.resolve({}) })],
  useProcessConnectionRegistrationMutation: () => [() => ({ unwrap: () => Promise.resolve({}) })],
  useUpdateConnectionByIdMutation: () => [() => ({ unwrap: updateUnwrap })],
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

// Master loads credentials through a Redux-backed hook. This file already
// stubs the connection RTK hooks and does not exercise the credential step.
vi.mock('@/utils/hooks/useAllUserCredentials', () => ({
  useAllUserCredentials: () => ({ credentials: [], isLoading: false, isError: false }),
}));

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

// The step chrome pulls RJSF, and the configure receipt pulls connection
// metadata. Both drag in JS assets this unit test does not render.
vi.mock('../ConnectionWizardStepContent', () => ({
  StepHeader: ({ title, subtitle }: { title: string; subtitle?: string }) => (
    <div>
      <h2>{title}</h2>
      {subtitle ? <p>{subtitle}</p> : null}
    </div>
  ),
  KubernetesImportStep: () => <div>import kubeconfig</div>,
}));

vi.mock('../metadata', () => ({
  default: () => null,
}));

vi.mock('../ConnectionStateTransitionModal', () => ({
  default: () => null,
}));

import { kubernetesExtension } from './kubernetesExtension';
import { useConnectionWizard } from './useConnectionWizard';

const store = configureStore({
  reducer: {
    ui: (state = { controllerState: {}, connectionMetadataState: {} }) => state,
  },
});

const kubernetesKind = {
  kind: 'kubernetes',
  type: 'platform',
  subType: 'orchestrator',
  label: 'Kubernetes',
  description: 'Kubernetes cluster',
  flow: 'kubernetes',
  docsUrl: 'https://docs.meshery.io',
  connectionSchema: null,
  credentialSchema: null,
  svgColor: null,
  svgWhite: null,
} as unknown as ConnectionWizardKindConfig;

const renderWithProviders = (ui: React.ReactElement) =>
  render(
    <Provider store={store}>
      <SistentThemeProvider initialMode="light">{ui}</SistentThemeProvider>
    </Provider>,
  );

const makeCtx = (postConfig: Record<string, unknown>): WizardContext => {
  const data = {
    availableKinds: [],
    isLoadingKinds: false,
    connectionIconMap: {},
    kindConfig: kubernetesKind,
    connectionFormData: {},
    credentialMode: 'existing' as const,
    selectedCredentialId: '',
    credentialName: '',
    credentialFormData: {},
    skipCredentialVerification: false,
    kubeconfigFile: new File(['kubeconfig'], 'clusters.yaml', { lastModified: 1_700_000_000_000 }),
    registrationId: null,
    registrationResult: null,
    registrationError: null,
    postConfig,
  };
  const formRefs: WizardFormRefs = {
    connection: createRef(),
    credential: createRef(),
  };
  return {
    mode: 'create',
    data,
    patch: (partial) => Object.assign(data, partial),
    patchPostConfig: (partial) => {
      data.postConfig = { ...data.postConfig, ...partial };
    },
    services: {
      notify: vi.fn(),
      registerConnection: vi.fn(),
      connectConnection: vi.fn(),
      discoverKubeContexts: vi.fn(),
      uploadKubeconfig: vi.fn(),
      updateConnectionById: vi.fn(),
      setMeshsyncMode: vi.fn(),
      flushMeshsync: vi.fn(),
      credentials: [],
    },
    formRefs,
    advance: () => undefined,
  };
};

const Review = kubernetesExtension.registerStep.Component;
const Receipt = kubernetesExtension.receiptStep.Component;

describe('kubernetes connection wizard', () => {
  beforeEach(() => {
    discoverUnwrap.mockReset();
    uploadUnwrap.mockReset();
    updateUnwrap.mockReset();
    discoverTrigger.mockClear();
  });

  it('marks review contexts Reachable or Unreachable before they are saved', () => {
    renderWithProviders(
      <Review
        ctx={makeCtx({
          discoveredContexts: [
            { id: 'up', name: 'prod', server: 'https://up.example', reachable: true },
            {
              id: 'down',
              name: 'dev',
              server: 'https://down.example',
              reachable: false,
              authNote: 'uses an exec credential plugin',
            },
          ],
          contextChoices: {
            up: { selected: true, name: 'prod', meshsyncDeploymentMode: 'embedded' },
            down: { selected: true, name: 'dev', meshsyncDeploymentMode: 'embedded' },
          },
        })}
      />,
    );

    expect(screen.getByText('Reachable')).toBeInTheDocument();
    expect(screen.getByText('Unreachable')).toBeInTheDocument();
    expect(screen.getByText(/uses an exec credential plugin/)).toBeInTheDocument();
    expect(screen.queryByText('discovered')).not.toBeInTheDocument();
    expect(screen.queryByText('not found')).not.toBeInTheDocument();
  });

  it('renders a failure on Done when nothing was imported', async () => {
    const ctx = makeCtx({
      discoveredContexts: [
        { id: 'bad', name: 'bad-cluster', server: 'https://bad', reachable: true },
      ],
      contextChoices: {
        bad: { selected: true, name: 'bad-cluster', meshsyncDeploymentMode: 'embedded' },
      },
      connectOnImport: true,
    });
    ctx.services.uploadKubeconfig = vi.fn().mockResolvedValue({
      erroredContexts: [{ name: 'bad-cluster', error: 'save failed' }],
    });

    await kubernetesExtension.registerStep.onNext(ctx);

    renderWithProviders(<Receipt ctx={ctx} />);

    expect(screen.getByRole('heading', { name: 'Kubernetes import failed' })).toBeInTheDocument();
    expect(
      screen.queryByRole('heading', { name: 'Kubernetes import complete' }),
    ).not.toBeInTheDocument();
    expect(screen.getByTestId('kubernetes-import-outcome')).toHaveAttribute(
      'data-outcome',
      'failure',
    );
    expect(screen.getByText('bad-cluster')).toBeInTheDocument();
    expect(screen.getByText(/save failed/)).toBeInTheDocument();
  });

  it('lists errored, ignored, and failed-connect contexts with reasons on a partial import', async () => {
    const ctx = makeCtx({
      discoveredContexts: [
        { id: 'up', name: 'prod', server: 'https://up.example', reachable: true },
      ],
      contextChoices: {
        up: { selected: true, name: 'prod', meshsyncDeploymentMode: 'embedded' },
      },
      connectOnImport: true,
    });
    ctx.services.uploadKubeconfig = vi.fn().mockResolvedValue({
      registeredContexts: [{ connectionId: 'conn-prod', name: 'prod', reachable: true }],
      erroredContexts: [{ name: 'bad-cluster', error: 'credential rejected' }],
      ignoredContexts: [{ name: 'old-cluster', error: 'already ignored' }],
    });
    ctx.services.updateConnectionById = vi.fn().mockRejectedValue(new Error('connection refused'));

    expect(await kubernetesExtension.registerStep.onNext(ctx)).toBe(true);

    renderWithProviders(<Receipt ctx={ctx} />);

    expect(
      screen.getByRole('heading', { name: 'Kubernetes import completed with issues' }),
    ).toBeInTheDocument();
    expect(screen.getByTestId('kubernetes-import-outcome')).toHaveAttribute(
      'data-outcome',
      'partial',
    );
    expect(screen.getByText('bad-cluster')).toBeInTheDocument();
    expect(screen.getByText(/Errored: credential rejected/)).toBeInTheDocument();
    expect(screen.getByText('old-cluster')).toBeInTheDocument();
    expect(screen.getByText(/Ignored: already ignored/)).toBeInTheDocument();
    expect(screen.getAllByText('prod').length).toBeGreaterThan(0);
    expect(screen.getByText(/Failed to connect: connection refused/)).toBeInTheDocument();
  });

  it('keeps the success presentation when every context imports', () => {
    renderWithProviders(
      <Receipt
        ctx={makeCtx({
          createdCount: 1,
          connectedCount: 1,
          unreachableCount: 0,
          importedContexts: [{ name: 'prod', status: 'connected' }],
          importIssues: [],
        })}
      />,
    );

    expect(screen.getByRole('heading', { name: 'Kubernetes import complete' })).toBeInTheDocument();
    expect(screen.getByTestId('kubernetes-import-outcome')).toHaveAttribute(
      'data-outcome',
      'success',
    );
  });

  it('tells an empty kubeconfig apart from contexts skipped for unsupported auth', async () => {
    const empty = makeCtx({});
    empty.services.discoverKubeContexts = vi
      .fn()
      .mockResolvedValue({ contexts: [], skippedContexts: [] });
    expect(await kubernetesExtension.detailsStep.onNext(empty)).toBe(false);
    expect(empty.services.notify).toHaveBeenCalledWith(
      expect.objectContaining({ message: 'The kubeconfig does not contain any contexts.' }),
    );

    const skipped = makeCtx({});
    skipped.services.discoverKubeContexts = vi.fn().mockResolvedValue({
      contexts: [],
      skippedContexts: [
        {
          name: 'eks',
          reason:
            'The context could not be read from the kubeconfig (uses an exec credential plugin).',
        },
      ],
    });
    expect(await kubernetesExtension.detailsStep.onNext(skipped)).toBe(false);
    expect(skipped.services.notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: expect.stringContaining(
          'eks: The context could not be read from the kubeconfig (uses an exec credential plugin).',
        ),
      }),
    );
    expect(skipped.services.notify).not.toHaveBeenCalledWith(
      expect.objectContaining({ message: 'The kubeconfig does not contain any contexts.' }),
    );
  });

  it('keeps review edits when Back then Continue reuses the same file, and Done does not offer Back', async () => {
    const file = new File(['apiVersion: v1'], 'clusters.yaml', { lastModified: 1_700_000_000_000 });
    discoverUnwrap.mockResolvedValue({
      contexts: [
        { id: 'ctx-1', name: 'prod', server: 'https://up.example', reachable: true },
        { id: 'ctx-2', name: 'dev', server: 'https://down.example', reachable: false },
      ],
      skippedContexts: [],
    });
    uploadUnwrap.mockResolvedValue({
      registeredContexts: [{ connectionId: 'conn-1', name: 'renamed-prod', reachable: false }],
    });

    const { result } = renderHook(() =>
      useConnectionWizard({
        mode: 'create',
        isOpen: true,
        availableKinds: [kubernetesKind],
        initialKindConfig: kubernetesKind,
      }),
    );

    await act(async () => {
      await result.current.next();
    });
    act(() => {
      result.current.ctx.patch({ kubeconfigFile: file });
    });
    await act(async () => {
      await result.current.next();
    });

    expect(result.current.activeStep?.label).toBe('Review Contexts');
    act(() => {
      result.current.ctx.patchPostConfig({
        contextChoices: {
          'ctx-1': { selected: true, name: 'renamed-prod', meshsyncDeploymentMode: 'operator' },
          'ctx-2': { selected: false, name: 'dev', meshsyncDeploymentMode: 'embedded' },
        },
        connectOnImport: false,
      });
    });

    act(() => {
      result.current.back();
    });
    await act(async () => {
      await result.current.next();
    });

    expect(discoverUnwrap).toHaveBeenCalledTimes(1);
    const choices = result.current.ctx.data.postConfig.contextChoices as Record<
      string,
      { name: string; meshsyncDeploymentMode: string; selected: boolean }
    >;
    expect(choices['ctx-1']).toMatchObject({
      name: 'renamed-prod',
      meshsyncDeploymentMode: 'operator',
      selected: true,
    });
    expect(result.current.ctx.data.postConfig.connectOnImport).toBe(false);

    await act(async () => {
      await result.current.next();
    });

    expect(result.current.activeStep?.label).toBe('Done');
    expect(result.current.canGoBack).toBe(false);

    const sent = discoverTrigger.mock.calls[0]?.[0] as { body: FormData } | undefined;
    expect(sent?.body.get('discoveryDetail')).toBe('true');
  });
});
