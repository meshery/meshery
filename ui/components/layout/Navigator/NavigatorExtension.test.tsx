import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Hoisted state we can swap between tests.
const remoteState: { loading: boolean; err: any; remoteProps: any | null } = {
  loading: false,
  err: null,
  remoteProps: null,
};

vi.mock('@paciolan/remote-component', () => ({
  createUseRemoteComponent: () => (_url: string) => {
    const RemoteComponent = ({ injectProps }: any) => {
      remoteState.remoteProps = injectProps;
      return <div data-testid="remote-component">remote</div>;
    };
    return [remoteState.loading, remoteState.err, RemoteComponent];
  },
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

vi.mock('../../PromptComponent', () => ({
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
  DynamicFullScreenLoader: ({ children, isLoading }: any) => (
    <div data-testid="dynamic-loader" data-loading={String(Boolean(isLoading))}>
      {children}
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

import NavigatorExtension from './NavigatorExtension';

describe('NavigatorExtension', () => {
  beforeEach(() => {
    remoteState.loading = false;
    remoteState.err = null;
    remoteState.remoteProps = null;
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
  });
});
