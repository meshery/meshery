import React, { createRef } from 'react';
import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// The real shipped connection definitions: the modal renders their
// transitionMap descriptions verbatim, so the tests below double as a drift
// guard on the authored copy.
import kubernetesDefinition from '../../../models/meshery-core/0.7.2/v1.0.0/connections/KubernetesConnection.json';
import grafanaDefinition from '../../../models/meshery-core/0.7.2/v1.0.0/connections/GrafanaConnection.json';
import prometheusDefinition from '../../../models/meshery-core/0.7.2/v1.0.0/connections/PrometheusConnection.json';
import artifactHubDefinition from '../../../models/meshery-core/0.7.2/v1.0.0/connections/ArtifactHubConnection.json';
import githubDefinition from '../../../models/meshery-core/0.7.2/v1.0.0/connections/GitHubConnection.json';

// Mutable store state so individual tests can swap in stale/absent definition
// data. Reset in beforeEach.
const mockStoreState: {
  ui: {
    connectionMetadataState: Record<string, { name?: string; transitionMap?: unknown }> | null;
  };
} = { ui: { connectionMetadataState: null } };

const seed = (definition: { name: string; kind: string; transitionMap: unknown }) => [
  definition.kind,
  { name: definition.name, transitionMap: definition.transitionMap },
];

const defaultConnectionMetadataState = () =>
  Object.fromEntries(
    [
      kubernetesDefinition,
      grafanaDefinition,
      prometheusDefinition,
      artifactHubDefinition,
      githubDefinition,
    ].map(seed),
  );

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) => selector(mockStoreState),
}));

vi.mock('@sistent/sistent', () => ({
  Modal: ({ open, title, children, headerIcon }) =>
    open ? (
      <div data-testid="modal">
        {headerIcon}
        <h2>{title}</h2>
        {children}
      </div>
    ) : null,
  ModalBody: ({ children }) => <div>{children}</div>,
  ModalFooter: ({ children }) => <div>{children}</div>,
  ModalButtonPrimary: ({ children, onClick, ...props }) => (
    <button data-variant="primary" onClick={onClick} type="button" {...props}>
      {children}
    </button>
  ),
  ModalButtonSecondary: ({ children, onClick, ...props }) => (
    <button data-variant="secondary" onClick={onClick} type="button" {...props}>
      {children}
    </button>
  ),
  ModalButtonDanger: ({ children, onClick, ...props }) => (
    <button data-variant="danger" onClick={onClick} type="button" {...props}>
      {children}
    </button>
  ),
  CustomTooltip: ({ children, title }) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  IconButton: ({ children, onClick, ...props }) => {
    // Drop MUI-only props so React does not warn in the unit mock.
    const { disableRipple: _disableRipple, sx: _sx, size: _size, ...domProps } = props;
    return (
      <button onClick={onClick} type="button" {...domProps}>
        {children}
      </button>
    );
  },
  InfoOutlinedIcon: () => <svg data-testid="info-outlined-icon" />,
  // Capture fill so we can assert amber/warning theming.
  WarningIcon: (props) => <svg data-testid="warning-icon" data-fill={props.fill} />,
  Box: ({ children, component, ...props }) => {
    const Tag = component || 'div';
    return <Tag {...props}>{children}</Tag>;
  },
  Typography: ({ children, component, ...props }) => {
    const Tag = component || 'span';
    return <Tag {...props}>{children}</Tag>;
  },
  useTheme: () => ({
    palette: {
      text: { secondary: 'gray', constant: { white: 'white' } },
      common: { white: 'white' },
      // Named tokens only in the mock (no hex/rgb literals — ui color-literal audit).
      status: { warning: 'warning-amber' },
      background: {
        error: { default: 'error-default', hover: 'error-hover' },
        warning: { default: 'warning-amber' },
      },
      warning: { main: 'warning-amber' },
    },
  }),
}));

vi.mock('../../css/icons.styles', () => ({
  iconLarge: {},
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../../utils/utils', () => ({
  formatToTitleCase: (value) => (value ? value.charAt(0).toUpperCase() + value.slice(1) : value),
}));

import ConnectionStateTransitionModal, {
  KUBERNETES_CONNECTION_LIFECYCLE_DOCS_URL,
  shouldShowTransitionDescription,
} from './ConnectionStateTransitionModal';
import type { ConnectionStateTransitionModalRef } from './ConnectionStateTransitionModal';

const setup = () => {
  const ref = createRef<ConnectionStateTransitionModalRef>();
  render(<ConnectionStateTransitionModal ref={ref} />);
  return ref;
};

beforeEach(() => {
  mockStoreState.ui.connectionMetadataState = defaultConnectionMetadataState();
});

describe('shouldShowTransitionDescription', () => {
  it('hides prompt-style "Are you sure" leftovers from older definition versions', () => {
    expect(
      shouldShowTransitionDescription(
        'Are you sure you want to transition from discovered to registered?',
        { currentStatus: 'discovered', targetStatus: 'registered' },
      ),
    ).toBe(false);
  });

  it('hides copy that merely restates the from→to transition', () => {
    expect(
      shouldShowTransitionDescription('transition from discovered to registered', {
        currentStatus: 'discovered',
        targetStatus: 'registered',
      }),
    ).toBe(false);
  });

  it('keeps authored descriptions that explain the consequences', () => {
    expect(
      shouldShowTransitionDescription('Registration description from the connection definition.', {
        currentStatus: 'discovered',
        targetStatus: 'registered',
      }),
    ).toBe(true);
  });

  it.each([
    ['KubernetesConnection.json', kubernetesDefinition],
    ['GrafanaConnection.json', grafanaDefinition],
    ['PrometheusConnection.json', prometheusDefinition],
    ['ArtifactHubConnection.json', artifactHubDefinition],
    ['GitHubConnection.json', githubDefinition],
  ])('every authored description in %s is displayable (no prompt-style copy)', (_name, def) => {
    Object.entries(def.transitionMap).forEach(([currentStatus, transitions]) => {
      (transitions as { nextState: string; description?: string }[]).forEach((transition) => {
        expect(
          shouldShowTransitionDescription(transition.description, {
            currentStatus,
            targetStatus: transition.nextState,
          }),
        ).toBe(true);
      });
    });
  });
});

describe('ConnectionStateTransitionModal', () => {
  it('renders nothing until show() is called', () => {
    setup();
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('explains a Kubernetes delete with the definition-authored copy and a red confirm button', async () => {
    const ref = setup();

    let resolved: boolean | undefined;
    act(() => {
      ref.current
        .show({
          targetStatus: 'deleted',
          currentStatus: 'connected',
          kind: 'kubernetes',
          connections: [{ id: 'c1', name: 'prod-cluster' }],
        })
        .then((value) => (resolved = value));
    });

    expect(await screen.findByText('Delete Kubernetes connection?')).toBeInTheDocument();
    expect(screen.getByTestId('warning-icon')).toHaveAttribute('data-fill', 'warning-amber');
    expect(screen.getByTestId('connection-transition-lead')).toHaveTextContent('prod-cluster');

    // Body copy comes from KubernetesConnection.json (connected → deleted),
    // not from anything hardcoded in the UI.
    const description = screen.getByTestId('connection-transition-description');
    expect(description).toHaveTextContent('associated credential');
    expect(description).toHaveTextContent('Meshery Operator');
    expect(description).toHaveTextContent('when present');
    expect(description).toHaveTextContent('The Kubernetes cluster itself is not deleted');
    expect(description).toHaveTextContent('auto-reconnect');
    expect(screen.queryByTestId('connection-transition-fallback')).not.toBeInTheDocument();

    const confirm = screen.getByTestId('connection-transition-confirm');
    expect(confirm).toHaveAttribute('data-variant', 'danger');
    expect(confirm).toHaveTextContent('Delete');

    await userEvent.click(confirm);
    expect(resolved).toBe(true);
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  // Title casing the kind slug yields "Github"/"Artifacthub". The definition's
  // authored `name` is the display label, so the dialog must use it.
  it.each([
    ['github', 'Delete GitHub connection?', 'repository itself is not modified'],
    [
      'artifacthub',
      'Delete Artifact Hub connection?',
      'Artifact Hub instance itself is not modified',
    ],
  ])(
    'titles a %s transition with the definition-authored display name',
    async (kind, expectedTitle, expectedCopy) => {
      const ref = setup();

      act(() => {
        ref.current.show({
          targetStatus: 'deleted',
          currentStatus: 'connected',
          kind,
          connections: [{ id: 'c1', name: 'meshery/meshery' }],
        });
      });

      expect(await screen.findByText(expectedTitle)).toBeInTheDocument();
      expect(screen.getByTestId('connection-transition-description')).toHaveTextContent(
        expectedCopy,
      );
      expect(screen.queryByTestId('connection-transition-fallback')).not.toBeInTheDocument();
    },
  );

  it('falls back to title-cased kind when the definition carries no display name', async () => {
    mockStoreState.ui.connectionMetadataState = {
      github: { transitionMap: githubDefinition.transitionMap },
    };
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        currentStatus: 'connected',
        kind: 'github',
        connections: [{ id: 'c1', name: 'meshery/meshery' }],
      });
    });

    expect(await screen.findByText('Delete Github connection?')).toBeInTheDocument();
  });

  it('falls back to the generic line when the current state is unknown', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        kind: 'kubernetes',
        connections: [{ id: 'c1', name: 'prod-cluster' }],
      });
    });

    await screen.findByText('Delete Kubernetes connection?');
    expect(screen.queryByTestId('connection-transition-description')).not.toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-fallback')).toHaveTextContent(
      'This will transition the connection to the Deleted state.',
    );
  });

  it('falls back when the definition still carries prompt-style copy (stale registry rows)', async () => {
    mockStoreState.ui.connectionMetadataState = {
      kubernetes: {
        transitionMap: {
          'not found': [
            {
              nextState: 'deleted',
              description:
                'Are you sure you want to transition from not found to deleted? This will remove the unreachable connection completely by unregistering it.',
            },
          ],
        },
      },
    };
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        currentStatus: 'not found',
        kind: 'kubernetes',
        connections: [{ id: 'c1', name: 'unreachable-cluster' }],
      });
    });

    await screen.findByText('Delete Kubernetes connection?');
    expect(screen.queryByTestId('connection-transition-description')).not.toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-fallback')).toBeInTheDocument();
  });

  it('falls back for kinds without a connection definition', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        currentStatus: 'connected',
        kind: 'meshery',
        connections: [{ id: 'c1', name: 'Meshery Cloud' }],
      });
    });

    await screen.findByText('Delete Meshery connection?');
    expect(screen.queryByTestId('connection-transition-description')).not.toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-fallback')).toBeInTheDocument();
  });

  it('normalizes mixed-case targetStatus so delete still uses danger styling', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'DELETED',
        kind: 'kubernetes',
        connections: [{ id: 'c1', name: 'prod-cluster' }],
      });
    });

    expect(await screen.findByText('Delete Kubernetes connection?')).toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-confirm')).toHaveAttribute(
      'data-variant',
      'danger',
    );
  });

  it('reconciles bulk title count with "and N more" when some names are empty', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        kind: 'kubernetes',
        connections: [
          { id: 'c1', name: 'docker-desktop' },
          { id: 'c2', name: 'meshery' },
          { id: 'c3', name: 'Artifact Hub' },
          { id: 'c4', name: '' },
          { id: 'c5', name: 'delta' },
        ],
      });
    });

    // 5 total; 3 listed names; remainder uses connections.length so 5 - 3 = 2.
    expect(await screen.findByText('Delete 5 Kubernetes connections?')).toBeInTheDocument();
    const lead = screen.getByTestId('connection-transition-lead');
    expect(lead).toHaveTextContent('docker-desktop, meshery, Artifact Hub');
    expect(lead).toHaveTextContent('and 2 more');
  });

  it('resolves the definition copy for a bulk selection whose current states agree', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        kind: 'kubernetes',
        connections: [
          { id: 'c1', name: 'alpha', status: 'connected' },
          { id: 'c2', name: 'beta', status: 'connected' },
        ],
      });
    });

    expect(await screen.findByText('Delete 2 Kubernetes connections?')).toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-description')).toHaveTextContent(
      'Meshery Operator',
    );
    expect(screen.getByTestId('connection-transition-bulk-scope')).toHaveTextContent(
      'This applies to each selected connection.',
    );
  });

  it('uses the generic line for a bulk selection with mixed current states', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        kind: 'kubernetes',
        connections: [
          { id: 'c1', name: 'alpha', status: 'connected' },
          { id: 'c2', name: 'beta', status: 'discovered' },
        ],
      });
    });

    expect(await screen.findByText('Delete 2 Kubernetes connections?')).toBeInTheDocument();
    expect(screen.queryByTestId('connection-transition-description')).not.toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-fallback')).toHaveTextContent(
      'This will transition each selected connection to the Deleted state.',
    );
  });

  it('settles an in-flight confirmation as cancelled when show() is re-entered', async () => {
    const ref = setup();

    let firstResolved: boolean | undefined;
    act(() => {
      ref.current
        .show({
          targetStatus: 'deleted',
          connections: [{ id: 'c1', name: 'first' }],
        })
        .then((value) => (firstResolved = value));
    });

    let secondResolved: boolean | undefined;
    act(() => {
      ref.current
        .show({
          targetStatus: 'deleted',
          connections: [{ id: 'c2', name: 'second' }],
        })
        .then((value) => (secondResolved = value));
    });

    await screen.findByTestId('connection-transition-confirm');
    expect(firstResolved).toBe(false);

    await userEvent.click(screen.getByTestId('connection-transition-confirm'));
    expect(secondResolved).toBe(true);
  });

  it('resolves false on cancel', async () => {
    const ref = setup();

    let resolved: boolean | undefined;
    act(() => {
      ref.current
        .show({
          targetStatus: 'deleted',
          connections: [{ id: 'c1', name: 'some-conn' }],
        })
        .then((value) => (resolved = value));
    });

    await userEvent.click(await screen.findByTestId('connection-transition-cancel'));
    expect(resolved).toBe(false);
  });

  it('uses action-style titles and primary confirm for forward transitions', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'registered',
        currentStatus: 'discovered',
        kind: 'kubernetes',
        connections: [{ id: 'c1', name: 'prod-cluster' }],
      });
    });

    // Reuses CONNECTION_STATE_TO_TRANSITION_MAP verbs (Register), same shape as Delete.
    expect(await screen.findByText('Register Kubernetes connection?')).toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-lead')).toHaveTextContent(
      'You are about to register the connection',
    );
    expect(screen.getByTestId('connection-transition-lead')).toHaveTextContent('prod-cluster');
    // Do not restate from→to in the lead; title/verb already carry the target.
    expect(screen.getByTestId('connection-transition-lead')).not.toHaveTextContent(
      'from DISCOVERED to REGISTERED',
    );
    const confirm = screen.getByTestId('connection-transition-confirm');
    expect(confirm).toHaveAttribute('data-variant', 'primary');
    expect(confirm).toHaveAttribute('data-severity', 'primary');
    expect(confirm).toHaveTextContent('Register');
    // Definition copy for discovered → registered (KubernetesConnection.json).
    expect(screen.getByTestId('connection-transition-description')).toHaveTextContent(
      'Registers the discovered cluster with Meshery',
    );
  });

  it('frames disconnect with a caution (warning) confirm, not delete-red or bare Confirm', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'disconnected',
        currentStatus: 'connected',
        kind: 'kubernetes',
        connections: [{ id: 'c1', name: 'prod-cluster' }],
      });
    });

    expect(await screen.findByText('Disconnect Kubernetes connection?')).toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-lead')).toHaveTextContent(
      'You are about to disconnect the connection',
    );
    // Definition copy for connected → disconnected (KubernetesConnection.json).
    expect(screen.getByTestId('connection-transition-description')).toHaveTextContent(
      'MeshSync data already collected is kept',
    );
    const confirm = screen.getByTestId('connection-transition-confirm');
    expect(confirm).toHaveAttribute('data-severity', 'caution');
    expect(confirm).toHaveTextContent('Disconnect');
    // Still the primary component (Sistent has no warning button); severity is via sx.
    expect(confirm).toHaveAttribute('data-variant', 'primary');
  });

  it('uses Discover verb for discovered transitions (consistent with Delete/Disconnect)', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'discovered',
        currentStatus: 'not found',
        connections: [{ id: 'c1', name: 'prom' }],
      });
    });

    expect(await screen.findByText('Discover connection?')).toBeInTheDocument();
    expect(screen.getByTestId('connection-transition-confirm')).toHaveTextContent('Discover');
  });

  it('offers an info tooltip linking the relevant docs', async () => {
    const ref = setup();

    act(() => {
      ref.current.show({
        targetStatus: 'deleted',
        kind: 'kubernetes',
        connections: [{ id: 'c1', name: 'prod-cluster' }],
      });
    });

    await screen.findByTestId('connection-transition-info');
    const tooltips = screen.getAllByTestId('tooltip');
    const docsTooltip = tooltips.find((tooltip) =>
      (tooltip.getAttribute('data-title') || '').includes(KUBERNETES_CONNECTION_LIFECYCLE_DOCS_URL),
    );
    expect(docsTooltip).toBeDefined();
    expect(docsTooltip.getAttribute('data-title')).toContain(
      KUBERNETES_CONNECTION_LIFECYCLE_DOCS_URL,
    );
    expect(docsTooltip.getAttribute('data-title')).toContain(
      'docs.meshery.io/concepts/logical/connections',
    );
  });
});
