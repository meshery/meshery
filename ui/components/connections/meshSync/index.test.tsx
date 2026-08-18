import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import MeshSyncTable from './index';

const notify = vi.fn();
const getK8sClusterIdsFromCtxId = vi.fn();
const useGetMeshSyncResourcesQuery = vi.fn();
const useGetMeshSyncResourceKindsQuery = vi.fn();
const getResponsiveColumnVisibility = vi.fn();
let dataTableProps: any;
let windowWidth = 1200;

vi.mock('@sistent/sistent', () => ({
  Tooltip: ({ children }) => <div>{children}</div>,
  Grid2: ({ children }) => <div>{children}</div>,
  Box: ({ children }) => <div>{children}</div>,
  Typography: ({ children }) => <span>{children}</span>,
  FormControl: ({ children }) => <div>{children}</div>,
  MenuItem: ({ children }) => <div>{children}</div>,
  Table: ({ children }) => <div>{children}</div>,
  FormattedTime: ({ date }) => <span>{String(date)}</span>,
  CustomColumnVisibilityControl: () => <div />,
  ResponsiveDataTable: (props) => {
    dataTableProps = props;
    return <div data-testid="mesh-sync-table" />;
  },
  SearchBar: () => <div />,
  UniversalFilter: ({
    selectedFilters,
    setSelectedFilters,
    handleApplyFilter,
  }: {
    selectedFilters: Record<string, string>;
    setSelectedFilters: (filters: Record<string, string>) => void;
    handleApplyFilter: (filters?: Record<string, string>) => void;
  }) => (
    <div data-testid="universal-filter">
      <button
        type="button"
        data-testid="apply-kind-pod"
        onClick={() => {
          const next = { ...selectedFilters, kind: 'Pod' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply Pod
      </button>
      <button
        type="button"
        data-testid="apply-kind-all"
        onClick={() => {
          const next = { ...selectedFilters, kind: 'All' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply kind All
      </button>
      <button
        type="button"
        data-testid="apply-model-core"
        onClick={() => {
          const next = { ...selectedFilters, model: 'core' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply model core
      </button>
      <button
        type="button"
        data-testid="apply-model-all"
        onClick={() => {
          const next = { ...selectedFilters, model: 'All' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply model All
      </button>
      <button
        type="button"
        data-testid="apply-namespace-default"
        onClick={() => {
          const next = { ...selectedFilters, namespace: 'default' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply namespace default
      </button>
      <button
        type="button"
        data-testid="apply-namespace-all"
        onClick={() => {
          const next = { ...selectedFilters, namespace: 'All' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply namespace All
      </button>
    </div>
  ),
  TableCell: ({ children }) => <div>{children}</div>,
  TableRow: ({ children }) => <div>{children}</div>,
  styled: (Component) => () => {
    const StyledComponent = ({ children, ...props }) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  },
  accentGrey: 'gray',
}));

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('../../../utils/multi-ctx', () => ({
  getK8sClusterIdsFromCtxId: (...args) => getK8sClusterIdsFromCtxId(...args),
}));

vi.mock('@/rtk-query/meshsync', () => ({
  useGetMeshSyncResourceKindsQuery: (...args) => useGetMeshSyncResourceKindsQuery(...args),
  useGetMeshSyncResourcesQuery: (...args) => useGetMeshSyncResourcesQuery(...args),
}));

vi.mock('../../../utils/dimension', () => ({
  useWindowDimensions: () => ({ width: windowWidth }),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector) =>
    selector({
      ui: {
        k8sConfig: { currentContext: 'dev' },
        selectedK8sContexts: ['all'],
      },
    }),
}));

vi.mock('../metadata', () => ({
  MeshSyncDataFormatter: () => <div />,
}));

vi.mock('../common', () => ({
  DefaultTableCell: () => <div />,
  SortableTableCell: () => <div />,
}));

vi.mock('../../../utils/utils', () => ({
  JsonParse: JSON.parse,
  camelcaseToSnakecase: (value) => value,
  getColumnValue: () => null,
  getVisibilityColums: () => [],
}));

vi.mock('../../../utils/responsive-column', () => ({
  getResponsiveColumnVisibility: (...args) => getResponsiveColumnVisibility(...args),
}));

vi.mock('../RegisterConnectionModal', () => ({
  default: () => <div />,
}));

vi.mock('../ConnectionChip', () => ({
  ConnectionStateChip: () => <div />,
}));

vi.mock('../styles', () => ({
  ContentContainer: ({ children }) => <div>{children}</div>,
  ConnectionStyledSelect: ({ children }) => <div>{children}</div>,
  InnerTableContainer: ({ children }) => <div>{children}</div>,
}));

vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }) => <div>{children}</div>,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

vi.mock('./MeshSyncEmptyState', () => ({
  default: () => <div data-testid="mesh-sync-empty-state" />,
}));

describe('MeshSyncTable', () => {
  beforeEach(() => {
    dataTableProps = undefined;
    notify.mockReset();
    getK8sClusterIdsFromCtxId.mockReset();
    useGetMeshSyncResourcesQuery.mockReset();
    useGetMeshSyncResourceKindsQuery.mockReset();
    getResponsiveColumnVisibility.mockReset();
    windowWidth = 1200;

    getK8sClusterIdsFromCtxId.mockReturnValue(['cluster-a', 'cluster-b']);
    useGetMeshSyncResourcesQuery.mockReturnValue({
      data: {
        resources: [
          {
            id: 'resource-1',
            metadata: { name: 'pod-a', namespace: 'default', creationTimestamp: '2026-05-08' },
            apiVersion: 'v1',
            kind: 'Pod',
            model: 'core',
            cluster_id: 'cluster-a',
            pattern_resources: '',
            status: 'discovered',
          },
        ],
        totalCount: 1,
      },
      isError: false,
      error: undefined,
    });
    useGetMeshSyncResourceKindsQuery.mockReturnValue({
      data: { kinds: [], namespaces: [] },
    });
    getResponsiveColumnVisibility.mockImplementation((columnNames, _colViews, width) =>
      Object.fromEntries(
        columnNames.map((columnName) => [
          columnName,
          columnName === 'model' ? width >= 1000 : true,
        ]),
      ),
    );
  });

  it('passes stable cluster ids to both queries', () => {
    render(<MeshSyncTable />);

    expect(useGetMeshSyncResourcesQuery).toHaveBeenCalledWith(
      expect.objectContaining({
        clusterIds: JSON.stringify(['cluster-a', 'cluster-b']),
      }),
    );
    expect(useGetMeshSyncResourceKindsQuery).toHaveBeenCalledWith(
      expect.objectContaining({
        clusterIds: ['cluster-a', 'cluster-b'],
      }),
    );
  });

  it('clears kind from the MeshSync query when kind is reset to All', async () => {
    const user = userEvent.setup();
    const { rerender } = render(<MeshSyncTable />);

    await user.click(screen.getByTestId('apply-kind-pod'));
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ kind: 'Pod' }),
      );
    });

    await user.click(screen.getByTestId('apply-kind-all'));
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ kind: null }),
      );
    });
  });

  it('clears model from the MeshSync query when model is reset to All', async () => {
    const user = userEvent.setup();
    const { rerender } = render(<MeshSyncTable />);

    await user.click(screen.getByTestId('apply-model-core'));
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ model: 'core' }),
      );
    });

    await user.click(screen.getByTestId('apply-model-all'));
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ model: null }),
      );
    });
  });

  it('clears namespace from the MeshSync query when namespace is reset to All', async () => {
    const user = userEvent.setup();
    const { rerender } = render(<MeshSyncTable />);

    await user.click(screen.getByTestId('apply-namespace-default'));
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ namespace: 'default' }),
      );
    });

    await user.click(screen.getByTestId('apply-namespace-all'));
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ namespace: null }),
      );
    });
  });

  it('notifies when fetching mesh sync resources fails', () => {
    useGetMeshSyncResourcesQuery.mockReturnValue({
      data: { resources: [], totalCount: 0 },
      isError: true,
      error: { data: 'MeshSync unavailable' },
    });

    render(<MeshSyncTable />);

    expect(notify).toHaveBeenNthCalledWith(
      1,
      expect.objectContaining({
        message: 'Error fetching MeshSync Resources',
        event_type: expect.objectContaining({ type: 'error' }),
        details: 'MeshSync unavailable',
      }),
    );
  });

  it('recomputes responsive column visibility when the window width changes', async () => {
    const { rerender } = render(<MeshSyncTable />);

    expect(dataTableProps.columnVisibility.model).toBe(true);

    windowWidth = 900;
    rerender(<MeshSyncTable />);

    await waitFor(() => {
      expect(dataTableProps.columnVisibility.model).toBe(false);
    });
  });
});
