import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import MeshSyncTable from './index';

const notify = vi.fn();
const getK8sClusterIdsFromCtxId = vi.fn();
const useGetMeshSyncResourcesQuery = vi.fn();
const useGetMeshSyncResourceKindsQuery = vi.fn();

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
  ResponsiveDataTable: () => <div data-testid="mesh-sync-table" />,
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
  useWindowDimensions: () => ({ width: 1200 }),
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
    notify.mockReset();
    getK8sClusterIdsFromCtxId.mockReset();
    useGetMeshSyncResourcesQuery.mockReset();
    useGetMeshSyncResourceKindsQuery.mockReset();

    getK8sClusterIdsFromCtxId.mockReturnValue(['cluster-a', 'cluster-b']);
    useGetMeshSyncResourcesQuery.mockReturnValue({
      data: { resources: [], totalCount: 0 },
      isError: false,
      error: undefined,
    });
    useGetMeshSyncResourceKindsQuery.mockReturnValue({
      data: { kinds: [], namespaces: [] },
    });
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
});
