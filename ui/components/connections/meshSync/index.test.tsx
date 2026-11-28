import React from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import MeshSyncTable from './index';

const notify = vi.fn();
const getK8sClusterIdsFromCtxId = vi.fn();
const useGetMeshSyncResourcesQuery = vi.fn();
const useGetMeshSyncResourceKindsQuery = vi.fn();
const getResponsiveColumnVisibility = vi.fn();
const getColumnValue = vi.fn();
let dataTableProps: any;
let searchBarProps: any;
let filterProps: any;
let visibilityControlProps: any;
let registerConnectionModalProps: any;
let statusSelectProps: any;
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
  CustomColumnVisibilityControl: (props) => {
    visibilityControlProps = props;
    return <div />;
  },
  ResponsiveDataTable: (props) => {
    dataTableProps = props;
    return <div data-testid="mesh-sync-table" />;
  },
  SearchBar: (props) => {
    searchBarProps = props;
    return <div />;
  },
  UniversalFilter: (props) => {
    filterProps = props;
    return <div />;
  },
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
  DataTableToolbar: ({ search, filter, columnVisibility }) => (
    <div data-testid="datatable-toolbar">
      {search}
      {filter}
      {columnVisibility}
    </div>
  ),
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
  JsonParse: (item) => (typeof item === 'string' ? JSON.parse(item || '{}') : item),
  camelcaseToSnakecase: (value) => value,
  getColumnValue: (...args) => getColumnValue(...args),
  getVisibilityColums: (columns) => columns.filter((col) => col.options?.display !== false),
}));

vi.mock('../../../utils/responsive-column', () => ({
  getResponsiveColumnVisibility: (...args) => getResponsiveColumnVisibility(...args),
}));

vi.mock('../RegisterConnectionModal', () => ({
  default: (props) => {
    registerConnectionModalProps = props;
    return <div />;
  },
}));

vi.mock('../ConnectionChip', () => ({
  ConnectionStateChip: ({ status }) => (
    <div data-testid="connection-state-chip" data-status={status}>
      {status}
    </div>
  ),
}));

vi.mock('../styles', () => ({
  ContentContainer: ({ children }) => <div>{children}</div>,
  ConnectionStyledSelect: ({ children, ...props }) => {
    statusSelectProps = props;
    return <div>{children}</div>;
  },
  MeshSyncNonInteractiveStatus: ({ children, ...props }) => (
    <span data-testid="meshsync-status-non-interactive" {...props}>
      {children}
    </span>
  ),
  InnerTableContainer: ({ children }) => <div>{children}</div>,
}));

vi.mock('./MeshSyncEmptyState', () => ({
  default: () => <div data-testid="mesh-sync-empty-state" />,
}));

describe('MeshSyncTable', () => {
  beforeEach(() => {
    dataTableProps = undefined;
    searchBarProps = undefined;
    filterProps = undefined;
    visibilityControlProps = undefined;
    registerConnectionModalProps = undefined;
    statusSelectProps = undefined;
    notify.mockReset();
    getColumnValue.mockReset();
    getColumnValue.mockReturnValue(null);
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
            pattern_resource: 'pattern-resource-1',
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

  it('reads and sorts the pattern resources column by the key the API sends', () => {
    render(<MeshSyncTable />);

    const { columns, data, options } = dataTableProps;
    const index = columns.findIndex(({ label }) => label === 'Pattern resources');
    expect(data[0][columns[index].name]).toBe('pattern-resource-1');
    // The visibility map is keyed by the same names.
    expect(dataTableProps.columnVisibility).toHaveProperty(columns[index].name);

    act(() => {
      options.onTableChange('sort', {
        activeColumn: index,
        announceText: 'Table now sorted by Pattern resources : descending',
      });
    });

    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ order: 'pattern_resource desc' }),
    );
  });

  it('offers the pattern resources column in the column visibility menu', () => {
    render(<MeshSyncTable />);

    expect(visibilityControlProps.columns.map(({ label }) => label)).toContain('Pattern resources');
  });

  it('keeps the pattern resources column hidden until it is turned on', async () => {
    const { getResponsiveColumnVisibility: actual } = await vi.importActual<any>(
      '../../../utils/responsive-column',
    );
    getResponsiveColumnVisibility.mockImplementation(actual);

    render(<MeshSyncTable />);

    expect(dataTableProps.columnVisibility.pattern_resource).toBe(false);
  });

  it('resets to the first page when the search changes', () => {
    render(<MeshSyncTable />);

    act(() => {
      dataTableProps.options.onTableChange('changePage', { page: 2 });
    });
    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ page: '2' }),
    );

    act(() => {
      searchBarProps.onSearch('cor');
    });

    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ page: 0, search: 'cor' }),
    );
  });

  it('resets to the first page when rows per page changes', () => {
    render(<MeshSyncTable />);

    act(() => {
      dataTableProps.options.onTableChange('changePage', { page: 2 });
    });
    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ page: '2' }),
    );

    // The datatable clamps its own page here; keeping page 2 with 25 rows
    // would ask for a page past the end and show an empty table.
    act(() => {
      dataTableProps.options.onTableChange('changeRowsPerPage', { rowsPerPage: 25, page: 0 });
    });

    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ page: 0, pagesize: '25' }),
    );
  });

  it('resets to the first page when filters are applied', () => {
    render(<MeshSyncTable />);

    act(() => {
      dataTableProps.options.onTableChange('changePage', { page: 2 });
    });

    act(() => {
      filterProps.setSelectedFilters({ kind: 'Pod', model: 'All', namespace: 'All' });
    });
    act(() => {
      filterProps.handleApplyFilter();
    });

    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ page: 0, kind: 'Pod' }),
    );
  });

  it('resets to the first page when the selected clusters change', () => {
    const { rerender } = render(<MeshSyncTable />);

    act(() => {
      dataTableProps.options.onTableChange('changePage', { page: 2 });
    });
    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({ page: '2' }),
    );

    getK8sClusterIdsFromCtxId.mockReturnValue(['cluster-c']);
    rerender(<MeshSyncTable />);

    expect(useGetMeshSyncResourcesQuery).toHaveBeenLastCalledWith(
      expect.objectContaining({
        page: 0,
        clusterIds: JSON.stringify(['cluster-c']),
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
    expect(notify).toHaveBeenCalledTimes(1);
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

  it('renders non-connection-capable discovered rows non-interactively with cursor: default and does not open the modal on click', () => {
    getColumnValue.mockImplementation((rowData, columnName, cols) => {
      const idx = cols.findIndex((col) => col?.name === columnName);
      return idx === -1 ? undefined : rowData[idx];
    });

    render(<MeshSyncTable />);

    const statusCol = dataTableProps.columns.find((col) => col.name === 'status');
    const rowData = dataTableProps.columns.map((col) => {
      if (col.name === 'id') return 'resource-uuid-non-conn';
      if (col.name === 'kind') return 'Pod';
      if (col.name === 'component_metadata') return { capabilities: [] };
      return null;
    });

    render(statusCol.options.customBodyRender('discovered', { rowData, rowIndex: 0 }));

    const wrapper = screen.getByTestId('meshsync-status-non-interactive');
    const chip = screen.getByTestId('connection-state-chip');
    expect(wrapper).toBeInTheDocument();
    expect(wrapper).toContainElement(chip);
    expect(chip).toHaveAttribute('data-status', 'discovered');
    expect(statusSelectProps).toBeUndefined();

    const parentClick = vi.fn();
    const { container } = render(
      <div onClick={parentClick}>
        {statusCol.options.customBodyRender('discovered', { rowData, rowIndex: 0 })}
      </div>,
    );
    fireEvent.click(
      container.querySelector('[data-testid="meshsync-status-non-interactive"]') as Element,
    );

    expect(parentClick).not.toHaveBeenCalled();
    expect(registerConnectionModalProps.openRegistrationModal).toBe(false);
  });

  it('stops click propagation on connection-capable rows and opens the registration modal on register selection', () => {
    getColumnValue.mockImplementation((rowData, columnName, cols) => {
      const idx = cols.findIndex((col) => col?.name === columnName);
      return idx === -1 ? undefined : rowData[idx];
    });

    render(<MeshSyncTable />);

    const statusCol = dataTableProps.columns.find((col) => col.name === 'status');
    const capabilities = { connection: true, urls: ['http://10.0.0.1:9090'] };
    const rowData = dataTableProps.columns.map((col) => {
      if (col.name === 'id') return 'resource-uuid-123';
      if (col.name === 'kind') return 'Service';
      if (col.name === 'metadata') return { name: 'prom-svc', uid: 'k8s-uid' };
      if (col.name === 'component_metadata') {
        return { capabilities, metadata: '{"host":"10.0.0.1"}' };
      }
      return null;
    });

    render(statusCol.options.customBodyRender('discovered', { rowData, rowIndex: 0 }));

    expect(statusSelectProps.value).toBe('discovered');

    const stopPropagation = vi.fn();
    act(() => {
      statusSelectProps.onClick({ stopPropagation, target: {} });
    });
    expect(stopPropagation).toHaveBeenCalledTimes(1);
    expect(registerConnectionModalProps.openRegistrationModal).toBe(false);

    act(() => {
      statusSelectProps.onChange({ target: { value: 'register' } });
    });

    expect(registerConnectionModalProps.openRegistrationModal).toBe(true);
    expect(registerConnectionModalProps.connectionData).toEqual({
      capabilities,
      metadata: { host: '10.0.0.1' },
      resourceID: 'resource-uuid-123',
    });
  });

  it('opens registration dialog for connection-capable Service with authoritative discovered status', () => {
    getColumnValue.mockImplementation((rowData, columnName, cols) => {
      const idx = cols.findIndex((col) => col?.name === columnName);
      return idx === -1 ? undefined : rowData[idx];
    });

    useGetMeshSyncResourcesQuery.mockReturnValue({
      data: {
        resources: [
          {
            id: 'svc-uuid-1',
            kind: 'Service',
            model: 'core',
            cluster_id: 'cluster-a',
            status: 'discovered',
            component_metadata: {
              capabilities: { connection: true, urls: ['http://10.0.0.1:9090'] },
              metadata: '{"host":"10.0.0.1"}',
            },
          },
        ],
        totalCount: 1,
      },
      isError: false,
    });

    render(<MeshSyncTable />);

    const { columns, data } = dataTableProps;
    expect(data[0].status).toBe('discovered');

    const statusCol = columns.find((col) => col.name === 'status');
    const rowData = columns.map((col) => {
      if (col.name === 'id') return 'svc-uuid-1';
      if (col.name === 'kind') return 'Service';
      if (col.name === 'component_metadata') {
        return {
          capabilities: { connection: true, urls: ['http://10.0.0.1:9090'] },
          metadata: '{"host":"10.0.0.1"}',
        };
      }
      return null;
    });

    render(statusCol.options.customBodyRender(data[0].status, { rowData, rowIndex: 0 }));

    expect(statusSelectProps).toBeDefined();
    expect(statusSelectProps.value).toBe('discovered');

    act(() => {
      statusSelectProps.onChange({ target: { value: 'register' } });
    });

    expect(registerConnectionModalProps.openRegistrationModal).toBe(true);
    expect(registerConnectionModalProps.connectionData).toEqual({
      capabilities: { connection: true, urls: ['http://10.0.0.1:9090'] },
      metadata: { host: '10.0.0.1' },
      resourceID: 'svc-uuid-1',
    });
  });

  it.each([
    ['missing', undefined],
    ['raw Kubernetes status', { phase: 'Running' }],
  ])(
    'does not project discovered status or offer registration for a connection-capable inventory row with %s',
    (_label, status) => {
      getColumnValue.mockImplementation((rowData, columnName, cols) => {
        const idx = cols.findIndex((col) => col?.name === columnName);
        return idx === -1 ? undefined : rowData[idx];
      });

      useGetMeshSyncResourcesQuery.mockReturnValue({
        data: {
          resources: [
            {
              id: 'svc-uuid-no-status',
              kind: 'Service',
              model: 'core',
              cluster_id: 'cluster-a',
              status,
              component_metadata: {
                capabilities: { connection: true, urls: ['http://10.0.0.1:9090'] },
                metadata: '{"host":"10.0.0.1"}',
              },
            },
          ],
          totalCount: 1,
        },
        isError: false,
      });

      render(<MeshSyncTable />);

      const { columns, data } = dataTableProps;
      expect(data[0].status).toBe(status);

      const statusCol = columns.find((col) => col.name === 'status');
      const rowData = columns.map((col) => {
        if (col.name === 'id') return 'svc-uuid-no-status';
        if (col.name === 'component_metadata') {
          return {
            capabilities: { connection: true, urls: ['http://10.0.0.1:9090'] },
            metadata: '{"host":"10.0.0.1"}',
          };
        }
        return null;
      });

      render(statusCol.options.customBodyRender(data[0].status, { rowData, rowIndex: 0 }));

      expect(statusSelectProps).toBeUndefined();
      expect(screen.getByTestId('meshsync-status-non-interactive')).toBeInTheDocument();
      expect(registerConnectionModalProps.openRegistrationModal).toBe(false);
    },
  );

  it.each(['ignored', 'connected', 'registered'])(
    'renders connection-capable %s rows without a registration dropdown',
    (status) => {
      getColumnValue.mockImplementation((rowData, columnName, cols) => {
        const idx = cols.findIndex((col) => col?.name === columnName);
        return idx === -1 ? undefined : rowData[idx];
      });

      render(<MeshSyncTable />);

      const statusCol = dataTableProps.columns.find((col) => col.name === 'status');
      const rowData = dataTableProps.columns.map((col) =>
        col.name === 'component_metadata' ? { capabilities: { connection: true } } : null,
      );

      render(statusCol.options.customBodyRender(status, { rowData, rowIndex: 0 }));

      expect(statusSelectProps).toBeUndefined();
      expect(screen.getByTestId('meshsync-status-non-interactive')).toBeInTheDocument();
      expect(screen.getByTestId('connection-state-chip')).toHaveAttribute('data-status', status);
    },
  );

  it.each([undefined, null, '', 0, 'unknown', { phase: 'Running' }])(
    'keeps registration unavailable for inventory row with missing or invalid status %s',
    (status) => {
      getColumnValue.mockImplementation((rowData, columnName, cols) => {
        const idx = cols.findIndex((col) => col?.name === columnName);
        return idx === -1 ? undefined : rowData[idx];
      });

      render(<MeshSyncTable />);

      const statusCol = dataTableProps.columns.find((col) => col.name === 'status');
      const rowData = dataTableProps.columns.map((col) =>
        col.name === 'component_metadata' ? { capabilities: { connection: true } } : null,
      );

      render(statusCol.options.customBodyRender(status, { rowData, rowIndex: 0 }));

      expect(statusSelectProps).toBeUndefined();
      expect(screen.getByTestId('meshsync-status-non-interactive')).toBeInTheDocument();
    },
  );

  it('resolves resourceID from the id column instead of the trailing metadata column when registering', () => {
    getColumnValue.mockImplementation((rowData, columnName, cols) => {
      const idx = cols.findIndex((col) => col?.name === columnName);
      return idx === -1 ? undefined : rowData[idx];
    });

    render(<MeshSyncTable />);

    const statusCol = dataTableProps.columns.find((col) => col.name === 'status');
    const rowData = dataTableProps.columns.map((col) => {
      if (col.name === 'id') return 'resource-uuid-123';
      if (col.name === 'kind') return 'Service';
      if (col.name === 'metadata') return { name: 'pod-a', uid: 'k8s-uid' };
      if (col.name === 'component_metadata') {
        return { capabilities: { connection: true }, metadata: '{"host":"10.0.0.1"}' };
      }
      return null;
    });

    render(statusCol.options.customBodyRender('discovered', { rowData, rowIndex: 0 }));

    act(() => {
      statusSelectProps.onChange({ target: { value: 'register' } });
    });

    expect(registerConnectionModalProps.connectionData).toEqual({
      capabilities: { connection: true },
      metadata: { host: '10.0.0.1' },
      resourceID: 'resource-uuid-123',
    });
  });

  it('does not throw when componentMetadata.metadata is undefined during status onChange', () => {
    getColumnValue.mockImplementation((rowData, columnName, cols) => {
      const idx = cols.findIndex((col) => col?.name === columnName);
      return idx === -1 ? undefined : rowData[idx];
    });

    render(<MeshSyncTable />);

    const statusCol = dataTableProps.columns.find((col) => col.name === 'status');
    const rowData = dataTableProps.columns.map((col) => {
      if (col.name === 'id') return 'resource-uuid-456';
      if (col.name === 'kind') return 'Service';
      if (col.name === 'component_metadata') return { capabilities: { connection: true } };
      return undefined;
    });

    render(statusCol.options.customBodyRender('discovered', { rowData, rowIndex: 0 }));

    expect(() => {
      act(() => {
        statusSelectProps.onChange({ target: { value: 'register' } });
      });
    }).not.toThrow();

    expect(registerConnectionModalProps.connectionData).toEqual({
      capabilities: { connection: true },
      metadata: undefined,
      resourceID: 'resource-uuid-456',
    });
  });
});
