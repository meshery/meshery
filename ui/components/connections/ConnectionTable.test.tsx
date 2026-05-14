import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ConnectionTable from './ConnectionTable';

const notify = vi.fn();
const push = vi.fn();
const ping = vi.fn();
const modalShow = vi.fn();
const updateConnectionByIdMutator = vi.fn();
const addConnectionToEnvironmentMutator = vi.fn();
const removeConnectionFromEnvironmentMutator = vi.fn();
const saveEnvironmentMutator = vi.fn();
const refetchConnections = vi.fn();
const getConnectionsQuery = vi.fn();
const getEnvironmentsQuery = vi.fn();
const updateVisibleColumns = vi.fn();
let dataTableProps: any;
let windowWidth = 1280;

const router = {
  query: {} as Record<string, unknown>,
  push,
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children }) => <div>{children}</div>,
  CustomColumnVisibilityControl: () => <div data-testid="column-visibility-control" />,
  SearchBar: () => <div data-testid="search-bar" />,
  UniversalFilter: () => <div data-testid="universal-filter" />,
  ResponsiveDataTable: (props) => {
    dataTableProps = props;
    return <div data-testid="responsive-data-table" />;
  },
  PROMPT_VARIANTS: {
    DANGER: 'danger',
    WARNING: 'warning',
  },
  MenuItem: ({ children }) => <div>{children}</div>,
  Box: ({ children }) => <div>{children}</div>,
  IconButton: ({ children, onClick, ...props }) => (
    <button onClick={onClick} type="button" {...props}>
      {children}
    </button>
  ),
  Typography: ({ children }) => <span>{children}</span>,
  Table: ({ children }) => <div>{children}</div>,
  Grid2: ({ children }) => <div>{children}</div>,
  Button: ({ children, onClick, disabled, ...props }) => (
    <button onClick={onClick} disabled={disabled} type="button" {...props}>
      {children}
    </button>
  ),
  ListItem: ({ children }) => <div>{children}</div>,
  FormControl: ({ children }) => <div>{children}</div>,
  styled: (Component) => () => {
    const StyledComponent = ({ children, ...props }) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  },
  accentGrey: 'gray',
  createTheme: () => ({ breakpoints: {} }),
  useTheme: () => ({
    palette: {
      error: { dark: 'darkred' },
      common: { white: 'white' },
    },
  }),
  TableCell: ({ children }) => <div>{children}</div>,
  TableRow: ({ children }) => <div>{children}</div>,
  Popover: ({ open, children }) => (open ? <div>{children}</div> : null),
  DeleteIcon: () => <svg data-testid="delete-icon" />,
}));

vi.mock('./styles', () => ({
  ContentContainer: ({ children }) => <div>{children}</div>,
  CreateButton: ({ children }) => <div>{children}</div>,
  InnerTableContainer: ({ children }) => <div>{children}</div>,
  ActionListItem: ({ children }) => <div>{children}</div>,
  ConnectionStyledSelect: ({ children }) => <div>{children}</div>,
}));

vi.mock('../data-formatter', () => ({
  FormatId: ({ id }) => <span>{id}</span>,
  formatDate: (value) => value,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../shared/LoadingState/LoadingComponent', () => ({
  default: () => <div data-testid="loading-screen" />,
}));

vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }) => <div>{children}</div>,
}));

vi.mock('../MesherySettingsEnvButtons', () => ({
  default: () => <div data-testid="settings-buttons" />,
}));

vi.mock('../../utils/utils', () => ({
  getVisibilityColums: (columns) => columns,
  getColumnValue: (rowData, columnName, columns) => {
    const columnIndex = columns.findIndex((column) => column.name === columnName);
    return columnIndex >= 0 ? rowData[columnIndex] : undefined;
  },
}));

vi.mock('../graphql/queries/ResetDatabaseQuery', () => ({
  default: vi.fn(),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => ping,
}));

vi.mock('./ConnectionChip', () => ({
  ConnectionStateChip: () => <div />,
  TooltipWrappedConnectionChip: () => <div />,
}));

vi.mock('./common', () => ({
  DefaultTableCell: () => <div />,
  SortableTableCell: () => <div />,
}));

vi.mock('./metadata', () => ({
  default: () => <div data-testid="connection-metadata" />,
}));

vi.mock('../../utils/responsive-column', () => ({
  getResponsiveColumnVisibility: (...args) => updateVisibleColumns(...args),
}));

vi.mock('../../utils/dimension', () => ({
  useWindowDimensions: () => ({ width: windowWidth }),
}));

vi.mock('../multi-select-wrapper', () => ({
  default: () => <div data-testid="multi-select-wrapper" />,
}));

vi.mock('../../rtk-query/environments', () => ({
  useAddConnectionToEnvironmentMutation: () => [addConnectionToEnvironmentMutator],
  useGetEnvironmentsQuery: (...args) => getEnvironmentsQuery(...args),
  useRemoveConnectionFromEnvironmentMutation: () => [removeConnectionFromEnvironmentMutator],
  useSaveEnvironmentMutation: () => [saveEnvironmentMutator],
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/utils/permission_constants', () => ({
  keys: {
    ASSIGN_CONNECTIONS_TO_ENVIRONMENT: { action: 'assign', subject: 'environment' },
    CHANGE_CONNECTION_STATE: { action: 'change', subject: 'connection' },
    DELETE_A_CONNECTION: { action: 'delete', subject: 'connection' },
    FLUSH_MESHSYNC_DATA: { action: 'flush', subject: 'meshsync' },
  },
}));

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args) => getConnectionsQuery(...args),
  useUpdateConnectionByIdMutation: () => [updateConnectionByIdMutator],
}));

vi.mock('../../assets/icons/InfoOutlined', () => ({
  default: () => <svg />,
}));

vi.mock('../../assets/icons/disconnect', () => ({
  default: () => <svg />,
}));

vi.mock('../PromptComponent', () => ({
  default: React.forwardRef(function PromptComponentMock(_, ref) {
    React.useImperativeHandle(ref, () => ({
      show: modalShow,
    }));
    return <div data-testid="prompt-component" />;
  }),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector) =>
    selector({
      ui: {
        organization: { id: 'org-1' },
        connectionMetadataState: {
          kubernetes: { transitions: ['connected'], icon: '/static/img/kubernetes.svg' },
        },
        controllerState: {},
      },
    }),
}));

const makeConnection = (overrides = {}) => ({
  id: 'connection-1',
  name: 'cluster-a',
  kind: 'kubernetes',
  status: 'connected',
  type: 'cluster',
  subType: 'managed',
  metadata: {
    name: 'cluster-a',
    server: 'https://cluster-a.local',
  },
  environments: [],
  created_at: '2026-05-08',
  ...overrides,
});

describe('ConnectionTable', () => {
  beforeEach(() => {
    dataTableProps = undefined;
    notify.mockReset();
    push.mockReset();
    ping.mockReset();
    modalShow.mockReset();
    refetchConnections.mockReset();

    updateConnectionByIdMutator.mockReset();
    addConnectionToEnvironmentMutator.mockReset();
    removeConnectionFromEnvironmentMutator.mockReset();
    saveEnvironmentMutator.mockReset();
    getConnectionsQuery.mockReset();
    getEnvironmentsQuery.mockReset();
    updateVisibleColumns.mockReset();
    windowWidth = 1280;

    updateConnectionByIdMutator.mockImplementation(({ connectionId, body }) => ({
      unwrap: () => Promise.resolve({ connectionId, body }),
    }));
    addConnectionToEnvironmentMutator.mockImplementation(() => ({
      unwrap: () => Promise.resolve({}),
    }));
    removeConnectionFromEnvironmentMutator.mockImplementation(() => ({
      unwrap: () => Promise.resolve({}),
    }));
    saveEnvironmentMutator.mockImplementation(() => ({
      unwrap: () => Promise.resolve({ id: 'env-1', name: 'dev' }),
    }));
    modalShow.mockResolvedValue('DELETE');

    router.query = {};

    getConnectionsQuery.mockReturnValue({
      data: {
        connections: [
          makeConnection(),
          makeConnection({
            id: 'connection-2',
            name: 'cluster-b',
            metadata: { name: 'cluster-b', server: 'https://cluster-b.local' },
          }),
        ],
        totalCount: 2,
      },
      isError: false,
      error: undefined,
      refetch: refetchConnections,
      isLoading: false,
    });

    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [{ id: 'env-1', name: 'dev' }] },
      isSuccess: true,
      isError: false,
      error: undefined,
    });
    updateVisibleColumns.mockImplementation((columnNames, _colViews, width) =>
      Object.fromEntries(
        columnNames.map((columnName) => [columnName, columnName === 'kind' ? width >= 1000 : true]),
      ),
    );
  });

  it('hydrates search from a string router query and passes it to the connections query', async () => {
    router.query = { searchText: 'cluster-a' };

    render(<ConnectionTable />);

    await waitFor(() => {
      expect(getConnectionsQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ search: 'cluster-a' }),
        undefined,
      );
    });
  });

  it('surfaces query failures through notifications', async () => {
    getConnectionsQuery.mockReturnValue({
      data: { connections: [], totalCount: 0 },
      isError: true,
      error: { data: 'connections unavailable' },
      refetch: refetchConnections,
      isLoading: false,
    });
    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [] },
      isSuccess: false,
      isError: true,
      error: { message: 'environments unavailable' },
    });

    render(<ConnectionTable />);

    await waitFor(() => {
      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Failed to fetch environment: environments unavailable',
        }),
      );
      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Failed to fetch connections: connections unavailable',
        }),
      );
    });
  });

  it('deletes the selected connections through the toolbar action', async () => {
    const user = userEvent.setup();

    render(<ConnectionTable />);

    const toolbar = dataTableProps.options.customToolbarSelect({
      data: [{ index: 0 }, { index: 1 }],
    });
    render(toolbar);

    await user.click(screen.getByRole('button', { name: /delete/i }));

    await waitFor(() => {
      expect(modalShow).toHaveBeenCalled();
      expect(updateConnectionByIdMutator).toHaveBeenCalledTimes(2);
    });

    expect(updateConnectionByIdMutator).toHaveBeenNthCalledWith(1, {
      connectionId: 'connection-1',
      body: { status: 'deleted' },
    });
    expect(updateConnectionByIdMutator).toHaveBeenNthCalledWith(2, {
      connectionId: 'connection-2',
      body: { status: 'deleted' },
    });
  });

  it('recomputes responsive column visibility when the window width changes', async () => {
    const { rerender } = render(<ConnectionTable />);

    expect(dataTableProps.columnVisibility.kind).toBe(true);

    windowWidth = 800;
    rerender(<ConnectionTable />);

    await waitFor(() => {
      expect(dataTableProps.columnVisibility.kind).toBe(false);
    });
  });
});
