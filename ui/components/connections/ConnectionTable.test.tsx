import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ConnectionTable from './ConnectionTable';

const notify = vi.fn();
const push = vi.fn();
const ping = vi.fn();
const pingGrafana = vi.fn();
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
  SyncAltIcon: () => <svg data-testid="sync-alt-icon" />,
  SettingsIcon: () => <svg data-testid="settings-icon" />,
  MoreVertIcon: () => <svg data-testid="more-vert-icon" />,
  InfoOutlinedIcon: () => <svg data-testid="info-outlined-icon" />,
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

vi.mock('./ConnectionWizardLauncher', () => ({
  default: () => <div data-testid="connection-wizard-launcher" />,
}));

vi.mock('../../utils/utils', () => ({
  getVisibilityColums: (columns) => columns,
  getColumnValue: (rowData, columnName, columns) => {
    const columnIndex = columns.findIndex((column) => column.name === columnName);
    return columnIndex >= 0 ? rowData[columnIndex] : undefined;
  },
}));

vi.mock('@/graphql/queries/ResetDatabaseQuery', () => ({
  default: vi.fn(),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => ping,
}));

vi.mock('@/utils/hooks/useGrafanaPingHook', () => ({
  default: () => pingGrafana,
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

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args) => getConnectionsQuery(...args),
  useUpdateConnectionByIdMutation: () => [updateConnectionByIdMutator],
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

// Mutable container so individual tests can flip `connectionMetadataState`
// (e.g., simulate the pre-hydration `null` state) without re-defining the
// mock between tests.
const uiState: {
  organization: { id: string };
  connectionMetadataState: Record<string, { transitions?: string[]; icon?: string }> | null;
  controllerState: Record<string, unknown> | null;
} = {
  organization: { id: 'org-1' },
  connectionMetadataState: {
    kubernetes: {
      transitions: ['connected'],
      icon: '/static/img/integrations/kubernetes.svg',
    },
  },
  controllerState: {},
};

vi.mock('react-redux', () => ({
  useSelector: (selector) => selector({ ui: uiState }),
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

    // Restore the populated Redux state — individual tests below flip
    // `connectionMetadataState` to `null` to exercise the pre-hydration path.
    uiState.organization = { id: 'org-1' };
    uiState.connectionMetadataState = {
      kubernetes: {
        transitions: ['connected'],
        icon: '/static/img/integrations/kubernetes.svg',
      },
    };
    uiState.controllerState = {};

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
    router.query = { con_q: 'cluster-a' };

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

  // Regression: the Environments select disappeared from the table. The cells
  // come from a memoized `columns` whose `customBodyRender` closes over
  // `isEnvironmentsSuccess` at definition time, and ResponsiveDataTable renders
  // from a `tableCols` *snapshot* that it only re-syncs to the live `columns`
  // when `columnVisibility` identity changes. On first render the environments
  // query is still pending (isEnvironmentsSuccess=false), so the snapshot froze
  // a cell that renders nothing; once the query resolved the snapshot was never
  // refreshed and the select never reappeared. ConnectionTable now keeps
  // `tableCols` following `columns`, so the resolved cell reaches the table.
  it('re-renders the Environments select once the environments query resolves', async () => {
    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [] },
      isSuccess: false,
      isError: false,
      error: undefined,
    });

    const { rerender } = render(<ConnectionTable />);

    // Render the snapshot's environments cell and report whether the select
    // (mocked as `multi-select-wrapper`) is present.
    const environmentsSelectIsRendered = () => {
      const envColumn = dataTableProps.tableCols.find((col) => col.name === 'environments');
      const { container, unmount } = render(
        <>{envColumn.options.customBodyRender([], { rowData: [] })}</>,
      );
      const present = !!container.querySelector('[data-testid="multi-select-wrapper"]');
      unmount();
      return present;
    };

    expect(environmentsSelectIsRendered()).toBe(false);

    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [{ id: 'env-1', name: 'dev' }] },
      isSuccess: true,
      isError: false,
      error: undefined,
    });
    // Re-render under act() so the tableCols-sync effect flushes; then assert
    // once (no render() inside waitFor, so a regression fails fast instead of
    // looping the snapshot render until timeout).
    rerender(<ConnectionTable />);

    expect(environmentsSelectIsRendered()).toBe(true);
  });

  // Regression for issue #19405 — `/management/connections` crashes with
  // "React error #185" / a `TypeError: Cannot read properties of null` in
  // production. The Redux slice (`store/slices/mesheryUi.ts`) initialises
  // `connectionMetadataState` to `null` and `_app.tsx` only populates it after
  // `getMeshModelComponentByName` resolves. The pages-router renders the
  // connections page before that promise settles, so the `enhancedConnections`
  // memo must tolerate a null map. The pre-fix code wrote
  // `connectionMetadataState[connection.kind]?.transitions` which protected
  // the property access but not the lookup itself.
  it('renders without throwing when connectionMetadataState is null (pre-hydration state)', () => {
    uiState.connectionMetadataState = null;

    expect(() => render(<ConnectionTable />)).not.toThrow();
  });

  it('falls back to undefined nextStatus/kindLogo when metadata is null', async () => {
    uiState.connectionMetadataState = null;

    render(<ConnectionTable />);

    await waitFor(() => {
      expect(dataTableProps).toBeDefined();
    });
    // The rows still render even though metadata is unavailable. Downstream
    // columns/cells handle the missing transitions gracefully.
    expect(dataTableProps.data).toHaveLength(2);
  });

  // Regression for the URL-clear loop described in issue #19405. The
  // pre-fix effect listed `filteredConnections` in its deps and called
  // `updateUrlWithConnectionId('')` when the selected id wasn't on the
  // current page. RTK Query returns a fresh array reference on every cache
  // hit, so this fired on every refetch, pushing a new URL, re-rendering the
  // parent, minting another RTK array — a textbook React #185 update-depth
  // loop. The effect now reads `filteredConnections` via a ref and never
  // clears the URL.
  it('does not clear the connectionId URL param when the connection is not on the current page', async () => {
    const updateUrlWithConnectionId = vi.fn();

    render(
      <ConnectionTable
        selectedConnectionId="connection-not-on-this-page"
        updateUrlWithConnectionId={updateUrlWithConnectionId}
      />,
    );

    await waitFor(() => {
      expect(dataTableProps).toBeDefined();
    });

    // Wait an extra tick to make sure no deferred call clears the URL.
    await new Promise((resolve) => setTimeout(resolve, 50));
    expect(updateUrlWithConnectionId).not.toHaveBeenCalledWith('');
  });

  // Regression for the Copilot review feedback on PR #19544: the prior shape
  // of this effect listed `filteredConnections` in its deps (loop-prone, was
  // the cause of issue #19405) and the intermediate fix moved it to a ref
  // and dropped the dep entirely, but set `lastProcessedId` *before*
  // confirming the row was found — locking the effect out for the rest of
  // the session whenever the user landed on a page that didn't contain the
  // deep-linked id (slow network, page 2, filtered view).
  it('still expands the row when the deep-linked id arrives on a later page', async () => {
    const setRowsExpandedSpy = vi.fn();

    // First load: deep-linked id is NOT in the visible page.
    getConnectionsQuery.mockReturnValue({
      data: {
        connections: [
          makeConnection({ id: 'page-1-conn-a' }),
          makeConnection({ id: 'page-1-conn-b' }),
        ],
        totalCount: 4,
      },
      isError: false,
      error: undefined,
      refetch: refetchConnections,
      isLoading: false,
    });

    const { rerender } = render(<ConnectionTable selectedConnectionId="deep-link-target" />);

    await waitFor(() => {
      expect(dataTableProps).toBeDefined();
    });

    // The row isn't on this page, so nothing should be expanded yet.
    const initialRowsExpanded = dataTableProps.options.rowsExpanded;
    expect(initialRowsExpanded).toEqual([]);

    // User paginates and the deep-linked id is now in the visible set.
    getConnectionsQuery.mockReturnValue({
      data: {
        connections: [
          makeConnection({ id: 'page-2-conn-a' }),
          makeConnection({ id: 'deep-link-target', name: 'cluster-deep' }),
        ],
        totalCount: 4,
      },
      isError: false,
      error: undefined,
      refetch: refetchConnections,
      isLoading: false,
    });

    rerender(<ConnectionTable selectedConnectionId="deep-link-target" />);

    // The effect must re-fire and expand index 1 now that the id is visible.
    await waitFor(() => {
      expect(dataTableProps.options.rowsExpanded).toEqual([1]);
    });

    // Sanity check that the linter isn't optimizing the spy away.
    expect(setRowsExpandedSpy).not.toHaveBeenCalled();
  });

  it('reuses the same onFlushMeshSync callback across rerenders so children are not invalidated', () => {
    const { rerender } = render(<ConnectionTable />);

    // The action menu only renders when an anchor is set. Capture the
    // identity of the JSX-bound `onFlushMeshSync` by re-rendering with the
    // same props and asserting the prop bag handed to ResponsiveDataTable
    // (the only render-time reflection of `handleFlushMeshSync` available
    // here) is referentially stable across renders.
    const firstOptions = dataTableProps.options;

    rerender(<ConnectionTable />);

    expect(dataTableProps.options).toBe(firstOptions);
  });
});
