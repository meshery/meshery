import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ConnectionManagementPageWithErrorBoundary from './index';

const router = {
  query: {} as Record<string, unknown>,
  pathname: '/connections',
  push: vi.fn(),
  isReady: true,
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('@sistent/sistent', () => ({
  NoSsr: ({ children }) => <>{children}</>,
  ErrorBoundary: ({ children }) => <>{children}</>,
  AppBar: ({ children }) => <div>{children}</div>,
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

vi.mock('./ConnectionTable', () => ({
  default: ({ selectedConnectionId, updateUrlWithConnectionId }) => (
    <div>
      <div data-testid="connection-table">connection:{selectedConnectionId ?? 'none'}</div>
      <button onClick={() => updateUrlWithConnectionId?.('cluster-2')} type="button">
        Update Connection Id
      </button>
    </div>
  ),
}));

vi.mock('./meshSync', () => ({
  default: ({ selectedResourceId, updateUrlWithResourceId }) => (
    <div>
      <div data-testid="meshsync-table">resource:{selectedResourceId ?? 'none'}</div>
      <button onClick={() => updateUrlWithResourceId?.('resource-2')} type="button">
        Update Resource Id
      </button>
    </div>
  ),
}));

vi.mock('../General/Modals/Modal', () => ({
  default: () => <div data-testid="create-connection-modal" />,
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/rtk-query/schema', () => ({
  useGetSchemaQuery: () => ({ data: undefined }),
}));

vi.mock('../General/error-404/index', () => ({
  default: () => <div data-testid="default-error" />,
}));

vi.mock('../General/ErrorBoundary', () => ({
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
    router.isReady = true;
  });

  it('defaults to the connections tab and ignores non-string connection ids', () => {
    router.query = {
      tab: ['meshsync'],
      connectionId: ['cluster-1'],
    };

    render(<ConnectionManagementPageWithErrorBoundary />);

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
});
