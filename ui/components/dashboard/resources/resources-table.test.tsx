import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const fetchMeshSyncResources = vi.fn();
const getClusterIds = vi.fn();

// Mock the leaf dependencies that the SUT pulls in. The SUT itself is NOT
// mocked — we import it and exercise its real exports below.
vi.mock('@/rtk-query/meshsync', () => ({
  useLazyGetMeshSyncResourcesQuery: () => [
    (...args: unknown[]) => fetchMeshSyncResources(...args),
    { data: undefined, isFetching: false },
  ],
  useGetMeshSyncResourceKindsQuery: () => ({
    data: { kinds: [], namespaces: ['default', 'kube-system'] },
  }),
}));

vi.mock('react-redux', () => ({
  useSelector: () => ({}),
}));

const routerState = {
  isReady: false,
  query: {} as Record<string, unknown>,
  pathname: '/dashboard',
  push: vi.fn(),
  replace: vi.fn(),
};
vi.mock('next/router', () => ({
  useRouter: () => routerState,
}));

vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }: any) => <div data-testid="tool-wrapper">{children}</div>,
}));

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('../../../utils/multi-ctx', () => ({
  getK8sClusterIdsFromCtxId: (...args: unknown[]) => getClusterIds(...args),
}));

vi.mock('../../../utils/responsive-column', () => ({
  updateVisibleColumns: () => ({}),
}));

vi.mock('../../../utils/dimension', () => ({
  useWindowDimensions: () => ({ width: 1024, height: 768 }),
}));

vi.mock('../../../utils/utils', () => ({
  camelcaseToSnakecase: (s: string) => s.replace(/([A-Z])/g, '_$1').toLowerCase(),
}));

vi.mock('../view', () => ({
  default: () => <div data-testid="view-component" />,
}));

vi.mock('./config', () => ({
  ALL_VIEW: 'all',
  SINGLE_VIEW: 'single',
}));

vi.mock('@sistent/sistent', () => ({
  CustomColumnVisibilityControl: () => <div data-testid="column-visibility" />,
  ResponsiveDataTable: () => <div data-testid="responsive-data-table" />,
  SearchBar: () => <div data-testid="search-bar" />,
  Slide: ({ children, in: visible }: any) =>
    visible ? <div data-testid="slide">{children}</div> : null,
  UniversalFilter: ({ handleApplyFilter, setSelectedFilters }: any) => (
    <div data-testid="universal-filter">
      <button
        type="button"
        data-testid="apply-namespace"
        onClick={() => {
          const next = { namespace: 'kube-system' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply Namespace
      </button>
      <button
        type="button"
        data-testid="apply-all"
        onClick={() => {
          const next = { namespace: 'All' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply All
      </button>
    </div>
  ),
}));

import ResourcesTable, { ACTION_TYPES } from './resources-table';

const baseProps = {
  updateProgress: vi.fn(),
  k8sConfig: {},
  useResourceConfig: () => ({
    name: 'Pod',
    columns: [{ name: 'name' }],
    colViews: [],
    options: {},
  }),
  submenu: false,
  workloadType: 'pods',
  selectedK8sContexts: ['all'],
};

describe('resources-table module', () => {
  beforeEach(() => {
    routerState.isReady = false;
    routerState.query = {};
    fetchMeshSyncResources.mockReset();
    getClusterIds.mockReset();
    getClusterIds.mockReturnValue(['cluster-a']);
    fetchMeshSyncResources.mockReturnValue({
      unwrap: () => Promise.resolve({ resources: [], page: 0, totalCount: 0, pageSize: 10 }),
    });
  });

  it('exports the FETCH_MESHSYNC_RESOURCES action type', () => {
    expect(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES).toBeDefined();
    expect(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES.name).toBe('FETCH_MESHSYNC_RESOURCES');
    expect(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES.error_msg).toBe(
      'Failed to fetch meshsync resources',
    );
  });

  it('renders nothing while next/router is not yet ready', () => {
    routerState.isReady = false;
    const { container } = render(<ResourcesTable {...baseProps} />);
    expect(container.firstChild).toBeNull();
  });

  it('mounts the inner table once next/router becomes ready', () => {
    routerState.isReady = true;
    render(<ResourcesTable {...baseProps} />);
    expect(screen.getByTestId('responsive-data-table')).toBeInTheDocument();
  });

  it('applies the namespace from the Apply payload on first click', async () => {
    const user = userEvent.setup();
    routerState.isReady = true;
    render(<ResourcesTable {...baseProps} />);

    await user.click(screen.getByTestId('apply-namespace'));

    await waitFor(() => {
      expect(fetchMeshSyncResources).toHaveBeenCalledWith(
        expect.objectContaining({ namespace: 'kube-system' }),
      );
    });
  });

  it('clears the namespace filter when Apply resets to All', async () => {
    const user = userEvent.setup();
    routerState.isReady = true;
    render(<ResourcesTable {...baseProps} />);

    await user.click(screen.getByTestId('apply-namespace'));
    await waitFor(() => {
      expect(fetchMeshSyncResources).toHaveBeenCalledWith(
        expect.objectContaining({ namespace: 'kube-system' }),
      );
    });

    fetchMeshSyncResources.mockClear();
    await user.click(screen.getByTestId('apply-all'));

    await waitFor(() => {
      expect(fetchMeshSyncResources).toHaveBeenCalled();
      const lastCall = fetchMeshSyncResources.mock.calls.at(-1)?.[0];
      expect(lastCall).not.toHaveProperty('namespace');
    });
  });
});
