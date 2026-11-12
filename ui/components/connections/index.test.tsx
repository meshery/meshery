import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Keys } from '@meshery/schemas/permissions';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ConnectionManagementPageWithErrorBoundary from './index';

const { mockUseHasPermission, mockUseGetSchemaQuery, mockOpenCreateConnection, mockReplace } =
  vi.hoisted(() => ({
    mockUseHasPermission: vi.fn(),
    mockUseGetSchemaQuery: vi.fn(),
    mockOpenCreateConnection: vi.fn(),
    mockReplace: vi.fn(),
  }));

const router = {
  query: {} as Record<string, unknown>,
  pathname: '/connections',
  push: vi.fn(),
  replace: mockReplace,
  isReady: true,
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('@sistent/sistent', () => ({
  NoSsr: ({ children }) => <>{children}</>,
  ErrorBoundary: ({ children }) => <>{children}</>,
  AppBar: ({ children }) => <div>{children}</div>,
  InfoTooltip: () => <svg data-testid="info-icon" />,
  useHasPermission: (permissionKey) => mockUseHasPermission(permissionKey),
  useTheme: () => ({ palette: { icon: { default: '#000' } } }),
}));

vi.mock('./styles', () => ({
  ConnectionIconText: ({ children }) => <span>{children}</span>,
  ConnectionTab: ({ label }) => <div>{label}</div>,
  ConnectionTabs: ({ children, onChange }) => (
    <div>
      <button onClick={(event) => onChange?.(event, 0)} type="button">
        Connections Tab
      </button>
      <button onClick={(event) => onChange?.(event, 1)} type="button">
        MeshSync Tab
      </button>
      {children}
    </div>
  ),
}));

// Records every `updateUrlWithConnectionId` reference the parent passes down.
// Tests assert against this list to catch callback identity churn (which
// previously cascaded into ConnectionTable's `options` memo and caused
// React error #185 on /management/connections).
const connectionTableCallbackRefs: Array<unknown> = [];

vi.mock('./ConnectionTable', () => ({
  default: ({ selectedConnectionId, updateUrlWithConnectionId, tabs }) => {
    connectionTableCallbackRefs.push(updateUrlWithConnectionId);
    return (
      <div>
        {tabs}
        <div data-testid="connection-table">connection:{selectedConnectionId ?? 'none'}</div>
        <button onClick={() => updateUrlWithConnectionId?.('cluster-2')} type="button">
          Update Connection Id
        </button>
      </div>
    );
  },
}));

vi.mock('@/utils/context/ConnectionWizardContextProvider', () => ({
  useConnectionWizardModal: () => ({
    openCreateConnection: mockOpenCreateConnection,
    closeCreateConnection: vi.fn(),
    open: false,
  }),
}));

vi.mock('./meshSync', () => ({
  default: ({ selectedResourceId, updateUrlWithResourceId, tabs }) => (
    <div>
      {tabs}
      <div data-testid="meshsync-table">resource:{selectedResourceId ?? 'none'}</div>
      <button onClick={() => updateUrlWithResourceId?.('resource-2')} type="button">
        Update Resource Id
      </button>
    </div>
  ),
}));

vi.mock('../shared/Modal/Modal', () => ({
  default: () => <div data-testid="create-connection-modal" />,
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/rtk-query/schema', () => ({
  useGetSchemaQuery: (args, options) => mockUseGetSchemaQuery(args, options),
}));

vi.mock('../general/error-404/index', () => ({
  default: ({ permissionKey }) => (
    <div data-testid="default-error" data-permission-key={JSON.stringify(permissionKey)} />
  ),
}));

vi.mock('../shared/ErrorBoundary/ErrorBoundary', () => ({
  default: () => <div data-testid="error-fallback" />,
}));

vi.mock('../../assets/icons/Connection', () => ({
  default: () => <svg data-testid="connection-icon" />,
}));

vi.mock('../../assets/icons/Meshsync', () => ({
  default: () => <svg data-testid="meshsync-icon" />,
}));

describe('connections index page', () => {
  beforeEach(() => {
    router.query = {};
    router.push.mockReset();
    router.replace.mockReset();
    router.isReady = true;
    connectionTableCallbackRefs.length = 0;
    mockOpenCreateConnection.mockReset();
    mockUseHasPermission.mockReset();
    mockUseHasPermission.mockReturnValue(true);
    mockUseGetSchemaQuery.mockReset();
    mockUseGetSchemaQuery.mockReturnValue({ data: undefined });
  });

  it('defaults to the connections tab and ignores non-string connection ids', () => {
    router.query = {
      tab: ['meshsync'],
      connectionId: ['cluster-1'],
    };

    render(<ConnectionManagementPageWithErrorBoundary />);

    expect(mockUseGetSchemaQuery).toHaveBeenCalledWith({ schemaName: 'helmRepo' }, { skip: false });
    expect(screen.getByTestId('connection-table')).toHaveTextContent('connection:none');
    expect(screen.queryByTestId('meshsync-table')).not.toBeInTheDocument();
  });

  it('renders the MeshSync tab when the query param is a string', () => {
    router.query = {
      tab: 'meshsync',
      connectionId: 'resource-1',
    };

    render(<ConnectionManagementPageWithErrorBoundary />);

    expect(screen.getByTestId('meshsync-table')).toHaveTextContent('resource:resource-1');
    expect(screen.queryByTestId('connection-table')).not.toBeInTheDocument();
  });

  it('updates the router query when the active tab changes', async () => {
    const user = userEvent.setup();

    render(<ConnectionManagementPageWithErrorBoundary />);

    await user.click(screen.getByRole('button', { name: 'MeshSync Tab' }));

    expect(router.push).toHaveBeenCalledWith(
      {
        pathname: '/connections',
        query: { tab: 'meshsync' },
      },
      undefined,
      { shallow: true },
    );
  });

  it('updates the router query when a child table selects a connection', async () => {
    const user = userEvent.setup();

    render(<ConnectionManagementPageWithErrorBoundary />);

    await user.click(screen.getByRole('button', { name: 'Update Connection Id' }));

    expect(router.push).toHaveBeenCalledWith(
      {
        pathname: '/connections',
        query: { connectionId: 'cluster-2' },
      },
      undefined,
      { shallow: true },
    );
  });

  // Regression test for React error #185 ("Maximum update depth exceeded") on
  // /management/connections. The parent used to recreate
  // `updateUrlWithConnectionId` on every render because its `useCallback` deps
  // listed `query` and `push` from `useRouter()`, both of which Next.js's
  // pages-router reassigns each render. That fresh callback identity
  // invalidated ConnectionTable's `options` memo and the in-table
  // expansion-sync effect, cascading into setState loops. The parent now
  // mirrors router state into refs so the callback identity is preserved.
  it('hands ConnectionTable a stable updateUrlWithConnectionId across re-renders', () => {
    router.query = {};

    const { rerender } = render(<ConnectionManagementPageWithErrorBoundary />);

    // Simulate Next.js handing us a brand-new `query` object reference on the
    // next render with the same logical contents (no URL change, just a
    // different object identity from `useRouter()`).
    router.query = {};
    rerender(<ConnectionManagementPageWithErrorBoundary />);

    // And again — but this time also include the `connectionId` the URL would
    // grow once a row is expanded. The dedupe guard in
    // `updateUrlWithConnectionId` previously closed over `connectionId` in its
    // dep list, so this is the render that historically minted a new callback.
    router.query = { connectionId: 'cluster-1' };
    rerender(<ConnectionManagementPageWithErrorBoundary />);

    expect(connectionTableCallbackRefs.length).toBeGreaterThanOrEqual(3);
    const [first, ...rest] = connectionTableCallbackRefs;
    expect(typeof first).toBe('function');
    for (const ref of rest) {
      expect(ref).toBe(first);
    }
  });

  it('skips schema query, ignores cached schema data, blocks deep links, and renders DefaultError when permission is denied', () => {
    mockUseHasPermission.mockReturnValue(false);
    mockUseGetSchemaQuery.mockReturnValue({
      data: { rjsfSchema: { title: 'Cached Helm Repo' }, uiSchema: {} },
    });
    router.query = {
      create: 'true',
      kind: 'kubernetes',
    };

    render(<ConnectionManagementPageWithErrorBoundary />);

    expect(mockUseHasPermission).toHaveBeenCalledWith(Keys.WorkspaceManagementViewConnections);
    expect(mockUseGetSchemaQuery).toHaveBeenCalledWith({ schemaName: 'helmRepo' }, { skip: true });
    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.WorkspaceManagementViewConnections),
    );
    expect(screen.queryByTestId('connection-table')).not.toBeInTheDocument();
    expect(screen.queryByTestId('meshsync-table')).not.toBeInTheDocument();
    expect(screen.queryByTestId('create-connection-modal')).not.toBeInTheDocument();
    expect(mockOpenCreateConnection).not.toHaveBeenCalled();
    expect(router.replace).not.toHaveBeenCalled();
  });

  it('hides tables and skips schema query when permission is revoked mid-session', () => {
    const { rerender } = render(<ConnectionManagementPageWithErrorBoundary />);

    expect(screen.getByTestId('connection-table')).toBeInTheDocument();
    expect(mockUseGetSchemaQuery).toHaveBeenLastCalledWith(
      { schemaName: 'helmRepo' },
      { skip: false },
    );

    mockUseHasPermission.mockReturnValue(false);
    rerender(<ConnectionManagementPageWithErrorBoundary />);

    expect(mockUseGetSchemaQuery).toHaveBeenLastCalledWith(
      { schemaName: 'helmRepo' },
      { skip: true },
    );
    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.WorkspaceManagementViewConnections),
    );
    expect(screen.queryByTestId('connection-table')).not.toBeInTheDocument();
  });
});
