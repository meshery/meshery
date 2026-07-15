import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const pingKubernetes = vi.fn();
const getControllerStatesByConnectionID = vi.fn();
const dispatchMock = vi.fn();
const notifyMock = vi.fn();
const fetchSystemSyncMock = vi.fn(() => ({ unwrap: () => Promise.resolve({ k8sConfig: [] }) }));
const updateConnectionStatusMock = vi.fn(() => ({ unwrap: () => Promise.resolve({}) }));
const useGetConnectionsQueryMock = vi.fn(() => ({ data: { connections: [] } }));
const useGetProviderCapabilitiesQueryMock = vi.fn(() => ({
  data: { providerUrl: 'https://x', extensions: {} },
  isError: false,
  error: null,
}));

const useMediaQueryMock = vi.fn(() => true);

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (sel: any) =>
    sel({
      ui: {
        controllerState: {},
        connectionMetadataState: { kubernetes: { icon: 'icons/k8s.svg' } },
      },
    }),
}));

vi.mock('../NotificationCenter/index', () => ({
  NotificationDrawerButton: () => <button data-testid="notification-drawer">bell</button>,
}));

vi.mock('../../User', () => ({
  default: () => <div data-testid="user" />,
}));

vi.mock('../../../utils/helpers/common', () => ({
  successHandlerGenerator: () => () => {},
  errorHandlerGenerator: () => () => {},
}));

vi.mock('../../connections/ConnectionChip', () => ({
  ConnectionChip: ({ title, onDelete, handlePing, status }: any) => (
    <div data-testid="connection-chip" data-status={status}>
      <button onClick={handlePing} type="button">
        {title}
      </button>
      {onDelete ? (
        <button data-testid="chip-delete" onClick={onDelete} type="button">
          delete
        </button>
      ) : null}
    </div>
  ),
}));

vi.mock('../../../rtk-query/system', () => ({
  useLazyGetSystemSyncQuery: () => [fetchSystemSyncMock],
}));

vi.mock('../../../rtk-query/connection', () => ({
  useUpdateConnectionStatusMutation: () => [updateConnectionStatusMock],
  useGetConnectionsQuery: () => useGetConnectionsQueryMock(),
}));

vi.mock('@/rtk-query/connection', () => ({
  useUpdateConnectionStatusMutation: () => [updateConnectionStatusMock],
  useGetConnectionsQuery: () => useGetConnectionsQueryMock(),
}));

vi.mock('../../../utils/Enum', () => ({
  CONNECTION_KINDS: { KUBERNETES: 'kubernetes' },
  CONNECTION_STATES: { DELETED: 'deleted' },
}));

vi.mock('../../PromptComponent', () => ({
  default: React.forwardRef((_props: any, ref: any) => {
    if (ref) {
      ref.current = { show: vi.fn(() => Promise.resolve('CANCEL')) };
    }
    return <div data-testid="prompt-component" />;
  }),
}));

vi.mock('../../../css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../../ExtensionSandbox', () => ({
  createPathForRemoteComponent: (uri: string) => `/remote/${uri}`,
}));

vi.mock('../../RemoteComponent', () => ({
  default: () => <div data-testid="remote-component" />,
}));

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: notifyMock }),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => pingKubernetes,
  useControllerStatus: () => ({
    getControllerStatesByConnectionID: (id?: string) => {
      getControllerStatesByConnectionID(id);
      return { operatorState: 'running', meshSyncState: 'connected', natsState: 'ok' };
    },
  }),
}));

vi.mock('../../../utils/utils', () => ({
  formatToTitleCase: (v: string) => v,
}));

vi.mock('../../registry/RegistryModal', () => ({
  default: () => <div data-testid="registry-modal" />,
}));

vi.mock('@sistent/sistent', () => ({
  Checkbox: ({ checked, onChange }: any) => (
    <input
      type="checkbox"
      checked={Boolean(checked)}
      onChange={(e) => onChange?.(e)}
      data-testid="header-checkbox"
      readOnly
    />
  ),
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={title}>
      {children}
    </div>
  ),
  Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  styled: (Component: any) => () => {
    const Styled = ({ children, ...props }: any) =>
      typeof Component === 'string' ? (
        React.createElement(Component, props, children)
      ) : (
        <div {...props}>{children}</div>
      );
    return Styled;
  },
  PROMPT_VARIANTS: { DANGER: 'danger' },
  TextField: (props: any) => <input data-testid="search-field" {...props} />,
  ClickAwayListener: ({ children }: any) => <>{children}</>,
  IconButton: ({ children, onClick, ...props }: any) => (
    <button type="button" onClick={onClick} {...props}>
      {children}
    </button>
  ),
  Slide: ({ children, in: open }: any) =>
    open ? <div data-testid="slide-open">{children}</div> : null,
  Grid2: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Hidden: ({ children }: any) => <>{children}</>,
  NoSsr: ({ children }: any) => <>{children}</>,
  useTheme: () => ({
    palette: {
      background: {
        secondary: '#222',
        card: '#fff',
        constant: { white: '#fff' },
        brand: { default: '#brand' },
      },
      common: { white: '#fff' },
      mode: 'light',
    },
    breakpoints: {
      up: () => '',
      down: () => '',
      between: () => '',
    },
    spacing: (val: number) => `${val * 8}px`,
    shadows: [],
  }),
  useMediaQuery: () => useMediaQueryMock(),
  SearchIcon: () => <svg data-testid="search-icon" />,
  SettingsIcon: () => <svg data-testid="settings-icon" />,
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  darkTeal: { main: '#000' },
  AppBar: ({ children }: any) => <header>{children}</header>,
  Toolbar: ({ children }: any) => <div>{children}</div>,
  Paper: ({ children }: any) => <div>{children}</div>,
  MenuIcon: () => <svg data-testid="menu-icon" />,
}));

vi.mock('@/utils/can', () => ({
  CanShow: ({ children }: any) => <>{children}</>,
}));

vi.mock('../../workspaces/SpacesSwitcher/SpaceSwitcher', () => ({
  default: () => <div data-testid="space-switcher" />,
}));

vi.mock('./HeaderMenu', () => ({
  default: () => <div data-testid="header-menu" />,
}));

vi.mock('../../connections/ConnectionFormModal', () => ({
  default: ({ isOpenModal }: any) => (isOpenModal ? <div data-testid="connection-modal" /> : null),
}));

vi.mock('../../MesherySettingsEnvButtons', () => ({
  default: () => <div data-testid="meshery-env-buttons" />,
}));

vi.mock('./Header.styles', () => {
  const make =
    (testId?: string) =>
    ({ children, ...props }: any) => (
      <div {...(testId ? { 'data-testid': testId } : {})} {...props}>
        {children}
      </div>
    );
  return {
    HeaderAppBar: ({ children, ...props }: any) => (
      <header data-testid="header-app-bar" {...props}>
        {children}
      </header>
    ),
    UserContainer: make('user-container'),
    PageTitleWrapper: make(),
    CBadgeContainer: make(),
    CMenuContainer: make(),
    HeaderIcons: () => <svg data-testid="header-icons-svg" />,
    MenuIconButton: ({ children, onClick }: any) => (
      <button type="button" onClick={onClick} data-testid="menu-icon-btn">
        {children}
      </button>
    ),
    UserSpan: make('user-span'),
    CBadge: ({ children, onClick }: any) => (
      <span data-testid="cbadge" onClick={onClick}>
        {children}
      </span>
    ),
    StyledToolbar: make('styled-toolbar'),
    UserInfoContainer: make('user-info-container'),
  };
});

vi.mock('@/rtk-query/user', () => ({
  getUserAccessToken: vi.fn(),
  getUserProfile: vi.fn(),
  useGetProviderCapabilitiesQuery: () => useGetProviderCapabilitiesQueryMock(),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateK8SConfig: (payload: any) => ({ type: 'ui/updateK8SConfig', payload }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'error' },
}));

vi.mock('../../../utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({ openModal: vi.fn() }),
}));

import Header, { K8sContextConnectionChip } from './Header';

describe('K8sContextConnectionChip', () => {
  beforeEach(() => {
    pingKubernetes.mockReset();
    getControllerStatesByConnectionID.mockReset();
  });

  it('renders a ConnectionChip with the supplied context information', () => {
    render(
      <K8sContextConnectionChip
        ctx={{ id: 'ctx-1', name: 'cluster-a', server: 'https://a', connectionId: 'conn-1' }}
        connectionMetadataState={{ kubernetes: { icon: 'icons/k8s.svg' } }}
        meshsyncControllerState={{}}
        connections={[{ id: 'conn-1', status: 'connected' }]}
      />,
    );

    const chip = screen.getByTestId('connection-chip');
    expect(chip).toHaveAttribute('data-status', 'connected');
    expect(screen.getByRole('button', { name: 'cluster-a' })).toBeInTheDocument();
  });

  it('does not render a delete control when onDelete is not provided', () => {
    render(
      <K8sContextConnectionChip
        ctx={{ id: 'ctx-1', name: 'cluster-a', server: 'https://a', connectionId: 'conn-1' }}
        connectionMetadataState={{}}
        meshsyncControllerState={{}}
        connections={[]}
      />,
    );

    expect(screen.queryByTestId('chip-delete')).not.toBeInTheDocument();
  });

  it('calls handlePing with the context details when clicked', () => {
    render(
      <K8sContextConnectionChip
        ctx={{ id: 'ctx-1', name: 'cluster-a', server: 'https://a', connectionId: 'conn-1' }}
        connectionMetadataState={{}}
        meshsyncControllerState={{}}
        connections={[]}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'cluster-a' }));
    expect(pingKubernetes).toHaveBeenCalledWith('cluster-a', 'https://a', 'conn-1');
  });

  it('renders a checkbox when selectable is true', () => {
    render(
      <K8sContextConnectionChip
        ctx={{ id: 'ctx-1', name: 'cluster-a', server: 'https://a', connectionId: 'conn-1' }}
        selectable
        selected={false}
        onSelectChange={vi.fn()}
        connectionMetadataState={{}}
        meshsyncControllerState={{}}
        connections={[]}
      />,
    );

    expect(screen.getByTestId('header-checkbox')).toBeInTheDocument();
  });

  it('renders a delete button when onDelete is provided', () => {
    const onDelete = vi.fn();
    render(
      <K8sContextConnectionChip
        ctx={{ id: 'ctx-1', name: 'cluster-a', server: 'https://a', connectionId: 'conn-1' }}
        onDelete={onDelete}
        connectionMetadataState={{}}
        meshsyncControllerState={{}}
        connections={[]}
      />,
    );

    fireEvent.click(screen.getByTestId('chip-delete'));
    expect(onDelete).toHaveBeenCalledWith('cluster-a', 'conn-1');
  });
});

describe('Header (default export)', () => {
  beforeEach(() => {
    dispatchMock.mockReset();
    notifyMock.mockReset();
    useGetProviderCapabilitiesQueryMock.mockReturnValue({
      data: { providerUrl: 'https://x', extensions: {} },
      isError: false,
      error: null,
    });
  });

  it('renders the header app bar with notification, user, and menu controls', () => {
    render(
      <Header
        onDrawerToggle={vi.fn()}
        contexts={{ totalCount: 0, contexts: [] }}
        activeContexts={[]}
        setActiveContexts={vi.fn()}
        searchContexts={vi.fn()}
      />,
    );

    expect(screen.getByTestId('header-app-bar')).toBeInTheDocument();
    expect(screen.getByTestId('notification-drawer')).toBeInTheDocument();
    expect(screen.getByTestId('user')).toBeInTheDocument();
    // The Header component wraps HeaderMenu in a span that the source code
    // tags with `data-testid="header-menu"`, and the test mock itself uses
    // the same id — so we expect at least one match.
    expect(screen.getAllByTestId('header-menu').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByTestId('registry-modal')).toBeInTheDocument();
  });

  it('renders the space switcher and respects user info container', () => {
    render(
      <Header
        onDrawerToggle={vi.fn()}
        contexts={{ totalCount: 0, contexts: [] }}
        activeContexts={[]}
        setActiveContexts={vi.fn()}
        searchContexts={vi.fn()}
      />,
    );
    expect(screen.getByTestId('space-switcher')).toBeInTheDocument();
    expect(screen.getByTestId('user-info-container')).toBeInTheDocument();
  });

  it('renders the k8s-cluster badge with the supplied total count', () => {
    render(
      <Header
        onDrawerToggle={vi.fn()}
        contexts={{ totalCount: 3, contexts: [] }}
        activeContexts={[]}
        setActiveContexts={vi.fn()}
        searchContexts={vi.fn()}
      />,
    );
    expect(screen.getByTestId('cbadge')).toHaveTextContent('3');
  });
});
