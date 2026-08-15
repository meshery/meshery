import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Hoisted state we can swap between tests. `component` models what
// `useRemoteComponent` resolved the bundle to, including `undefined` -- which
// is what a bundle without a CommonJS default export yields, with no error.
const RemoteComponentStub = ({ injectProps }: any) => {
  remoteState.remoteProps = injectProps;
  return <div data-testid="remote-component">remote</div>;
};

const remoteState: { loading: boolean; err: any; remoteProps: any | null; component: unknown } = {
  loading: false,
  err: null,
  remoteProps: null,
  component: RemoteComponentStub,
};

vi.mock('@paciolan/remote-component', () => ({
  createUseRemoteComponent: () => (_url: string) => [
    remoteState.loading,
    remoteState.err,
    remoteState.component,
  ],
  getDependencies: () => ({}),
  createRequires: () => () => ({}),
}));

vi.mock('react-redux', () => ({
  useSelector: (sel: any) =>
    sel({
      ui: {
        providerCapabilities: { providerUrl: 'https://x' },
        selectedK8sContexts: ['ctx-1'],
        organization: { id: 'org-1' },
      },
      mesheryUi: {},
    }),
  Provider: ({ children }: any) => <>{children}</>,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  selectK8sConfig: () => [],
  selectSelectedK8sClusters: () => [],
}));

vi.mock('../../../store', () => ({ store: { getState: () => ({}) } }));

vi.mock('../../performance', () => ({
  default: () => <div data-testid="performance-component" />,
}));

vi.mock('../../meshery-mesh-interface/PatternServiceFormCore', () => ({
  default: () => <div data-testid="pattern-form-core" />,
}));

vi.mock('../../shared/Modal/Information/InfoModal', () => ({
  default: () => <div data-testid="info-modal" />,
}));

vi.mock('../../general/PromptComponent', () => ({
  default: React.forwardRef(() => <div data-testid="prompt" />),
}));

vi.mock('../../../utils/disabledComponents', () => ({
  ProviderUiAccessControl: class {
    isNavigatorComponentEnabled() {
      return true;
    }
  },
}));

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('../../shared/Modal/Modal', () => ({
  default: () => <div data-testid="modal" />,
}));

vi.mock('../../designs/export/ExportDesignModal', () => ({
  default: () => <div data-testid="export-modal" />,
}));

vi.mock('../../data-formatter', () => ({
  FormatStructuredData: () => null,
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  useFilterK8sContexts: () => () => [],
}));

vi.mock('@/utils/context/dynamicContext', () => ({
  useDynamicComponent: () => null,
}));

vi.mock('../../designs/lifecycle/ValidateDesign', () => ({
  ValidateDesign: () => null,
}));

vi.mock('../../designs/lifecycle/DryRun', () => ({
  DryRunDesign: () => null,
}));

vi.mock('../../designs/lifecycle/DeployStepper', () => ({
  DeployStepper: () => null,
  UnDeployStepper: () => null,
}));

vi.mock('machines/validator/designValidator', () => ({
  designValidationMachine: { id: 'design-validator' },
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/utils/eventBus', () => ({
  mesheryEventBus: { publish: vi.fn() },
}));

vi.mock('@/theme/hooks', () => ({
  ThemeTogglerCore: () => null,
}));

vi.mock('../../meshery-mesh-interface/PatternService/RJSF', () => ({
  default: () => null,
}));

vi.mock('../../shared/LoadingState/DynamicFullscreenLoader', () => ({
  // The real loader withholds its children entirely while loading; the stub has
  // to do the same or tests see a tree the app never renders.
  DynamicFullScreenLoader: ({ children, isLoading }: any) => (
    <div data-testid="dynamic-loader" data-loading={String(Boolean(isLoading))}>
      {isLoading ? null : children}
    </div>
  ),
}));

vi.mock('../../TroubleshootingComponent', () => ({
  default: () => <div data-testid="troubleshoot" />,
}));

vi.mock('@/components/shared/FormFields/typing-filter', () => ({
  default: () => null,
}));

vi.mock('../../registry/CreateModelModal', () => ({ default: () => null }));
vi.mock('../../registry/ImportModelModal', () => ({ default: () => null }));

vi.mock('../../workspaces/ViewInfoModal', () => ({
  ViewInfoModal: () => null,
}));

vi.mock('@/store/ProviderStoreWrapper', () => ({
  default: ({ children }: any) => <>{children}</>,
}));

vi.mock('@/utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({
    openModalWithDefault: vi.fn(),
    onLoadResource: vi.fn(),
  }),
}));

vi.mock('../NotificationCenter/formatters/relationship_evaluation', () => ({
  RelationshipEvaluationTraceFormatter: () => null,
}));

vi.mock('@/utils/hooks/useRegistryModal', () => ({
  useRegistryModal: () => ({ openModal: vi.fn() }),
}));

import {
  MESHERY_EXTENSION_CONTRACT_VERSION,
  describeInjectedCapabilityReport,
  isInjectedCapabilityReportSatisfied,
  reportInjectedCapabilities,
} from '@sistent/sistent';

import NavigatorExtension, { buildExtensionInjectProps } from './NavigatorExtension';

describe('NavigatorExtension', () => {
  beforeEach(() => {
    remoteState.loading = false;
    remoteState.err = null;
    remoteState.remoteProps = null;
    remoteState.component = RemoteComponentStub;
  });

  it('renders the remote component wrapped in a fullscreen loader', () => {
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.getByTestId('dynamic-loader')).toHaveAttribute('data-loading', 'false');
    expect(screen.getByTestId('remote-component')).toBeInTheDocument();
  });

  it('renders an error state when the remote component fails to load', () => {
    remoteState.err = new Error('boom');
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText(/please pardon our mesh/i)).toBeInTheDocument();
    expect(screen.getByText(/boom/)).toBeInTheDocument();
    expect(screen.getByTestId('troubleshoot')).toBeInTheDocument();
  });

  it('passes the loading flag through to DynamicFullScreenLoader', () => {
    remoteState.loading = true;
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.getByTestId('dynamic-loader')).toHaveAttribute('data-loading', 'true');
  });

  it('injects the expected props onto the remote component', () => {
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(remoteState.remoteProps).toBeTruthy();
    expect(remoteState.remoteProps.providerCapabilities).toEqual({ providerUrl: 'https://x' });
    expect(remoteState.remoteProps.selectedK8sContexts).toEqual(['ctx-1']);
    expect(remoteState.remoteProps.currentOrganization).toEqual({ id: 'org-1' });
    expect(typeof remoteState.remoteProps.PatternServiceFormCore).toBe('function');
    // The mesheryStore exposes getters for k8s config / selected clusters.
    expect(typeof remoteState.remoteProps.mesheryStore.selectedK8sClusters.get).toBe('function');
    expect(typeof remoteState.remoteProps.mesheryStore.k8sConfig.get).toBe('function');
    // Backwards-compatible alias remains in place.
    expect(remoteState.remoteProps.CapabilitiesRegistryClass).toBe(
      remoteState.remoteProps.ProviderUiAccessControlClass,
    );
    // Extensions read this to detect that they were built against a different
    // contract revision than the host that loaded them.
    expect(remoteState.remoteProps.contractVersion).toBe(MESHERY_EXTENSION_CONTRACT_VERSION);
  });

  it('surfaces the real cause when the bundle exposes no default export', () => {
    // `useRemoteComponent` yields `undefined` with no error in this case, so
    // without an explicit guard React fails with an opaque "Element type is
    // invalid" far from the actual problem.
    remoteState.component = undefined;
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText(/no CommonJS default export/i)).toBeInTheDocument();
    // This is the bundler-config case, so the commonjs2 remedy is the right one.
    expect(screen.getByText(/commonjs2/i)).toBeInTheDocument();
    expect(screen.getByText(/https:\/\/remote\.example\/component\.js/)).toBeInTheDocument();
  });

  it('waits for loading to finish before reporting a missing default export', () => {
    remoteState.loading = true;
    remoteState.component = undefined;
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it.each([
    ['memo', () => React.memo(RemoteComponentStub)],
    ['forwardRef', () => React.forwardRef((props: any) => <RemoteComponentStub {...props} />)],
  ])('renders a %s-wrapped default export', (_label, makeComponent) => {
    remoteState.component = makeComponent();
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.getByTestId('remote-component')).toBeInTheDocument();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('rejects a default export that is an element rather than a component', () => {
    // A bundle that exports `<Foo />` instead of `Foo` still carries a
    // `$$typeof` tag, so a bare "has $$typeof" check would wave it through and
    // React would then fail with the opaque "Element type is invalid" this
    // guard exists to replace. Only forwardRef/memo/lazy tags are components.
    remoteState.component = <RemoteComponentStub />;
    render(<NavigatorExtension url="https://remote.example/component.js" />);
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText(/is not a React component/i)).toBeInTheDocument();
    // The remedies differ: this bundle exported *something*, so pointing the
    // author at the commonjs2 bundler setting would send them after the wrong fix.
    expect(screen.queryByText(/commonjs2/i)).not.toBeInTheDocument();
  });
});

describe('buildExtensionInjectProps', () => {
  it('satisfies every capability the extension contract declares', () => {
    // The host <-> extension boundary has no compile-time link: extensions read
    // these keys off an untyped bag, so a rename here is invisible until the
    // feature that reads it is dead in production. This assertion is the gate.
    const report = reportInjectedCapabilities(
      buildExtensionInjectProps({
        providerCapabilities: {},
        selectedK8sContexts: [],
        currentOrganization: { id: 'org-1' },
        openWorkspaceModal: vi.fn(),
        openRegistryModal: vi.fn(),
        setCurrentLoadedResourceInOrgWorkspaceSession: vi.fn(),
      }),
    );

    expect(
      isInjectedCapabilityReportSatisfied(report),
      String(describeInjectedCapabilityReport(report)),
    ).toBe(true);
  });
});
