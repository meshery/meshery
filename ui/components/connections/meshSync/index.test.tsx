import React from 'react';
import { render, waitFor } from '@testing-library/react';
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
  UniversalFilter: () => <div />,
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

vi.mock('./RegisterConnectionModal', () => ({
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
