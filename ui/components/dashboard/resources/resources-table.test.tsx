import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

// Mock the leaf dependencies that the SUT pulls in. The SUT itself is NOT
// mocked — we import it and exercise its real exports below.
vi.mock('@/rtk-query/meshsync', () => ({
  useLazyGetMeshSyncResourcesQuery: () => [
    vi.fn().mockResolvedValue({ resources: [], total_count: 0 }),
    { data: undefined, isFetching: false },
  ],
  useGetMeshSyncResourceKindsQuery: () => ({ data: { kinds: [] } }),
}));

vi.mock('react-redux', () => ({
  useSelector: () => ({}),
}));

const routerState = { isReady: false, query: {}, push: vi.fn(), replace: vi.fn() };
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
  getK8sClusterIdsFromCtxId: () => [],
}));

vi.mock('../../../utils/responsive-column', () => ({
  updateVisibleColumns: vi.fn(),
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
  UniversalFilter: () => <div data-testid="universal-filter" />,
}));

// Now import the real SUT.
import ResourcesTable, { ACTION_TYPES } from './resources-table';

describe('resources-table module', () => {
  beforeEach(() => {
    routerState.isReady = false;
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
    const { container } = render(
      <ResourcesTable
        updateProgress={vi.fn()}
        k8sConfig={{}}
        resourceConfig={() => ({ columns: [], options: {} })}
        submenu={false}
        workloadType="pods"
        selectedK8sContexts={[]}
      />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('mounts the inner table once next/router becomes ready', () => {
    routerState.isReady = true;
    render(
      <ResourcesTable
        updateProgress={vi.fn()}
        k8sConfig={{}}
        resourceConfig={() => ({ columns: [], options: {} })}
        submenu={false}
        workloadType="pods"
        selectedK8sContexts={[]}
      />,
    );
    expect(screen.getByTestId('responsive-data-table')).toBeInTheDocument();
  });
});
