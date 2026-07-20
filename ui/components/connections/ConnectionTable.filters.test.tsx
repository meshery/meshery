import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ConnectionTable from './ConnectionTable';

const notify = vi.fn();
const getConnectionsQuery = vi.fn();
const getEnvironmentsQuery = vi.fn();
const updateVisibleColumns = vi.fn();
let dataTableProps: { data?: unknown[] } | undefined;

const routerReplace = vi.fn();
const router = {
  query: {} as Record<string, unknown>,
  pathname: '/management/connections',
  push: vi.fn(),
  replace: (...args: unknown[]) => routerReplace(...args),
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children }) => <div>{children}</div>,
  getRelativeTime: (date: string) => date,
  getFullFormattedTime: (date: string) => date,
  CustomColumnVisibilityControl: () => null,
  SearchBar: () => null,
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
        data-testid="apply-status-connected"
        onClick={() => {
          const next = { ...selectedFilters, status: 'connected' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply connected
      </button>
      <button
        type="button"
        data-testid="apply-status-all"
        onClick={() => {
          const next = { ...selectedFilters, status: 'All' };
          setSelectedFilters(next);
          handleApplyFilter(next);
        }}
      >
        Apply all
      </button>
    </div>
  ),
  DataTableToolbar: ({ filter }: { filter?: React.ReactNode }) => <div>{filter}</div>,
  ResponsiveDataTable: (props) => {
    dataTableProps = props;
    return <div data-testid="responsive-data-table" />;
  },
  MenuItem: ({ children }) => <div>{children}</div>,
  Box: ({ children }) => <div>{children}</div>,
  SyncAltIcon: () => null,
  SettingsIcon: () => null,
  MoreVertIcon: () => null,
  InfoOutlinedIcon: () => null,
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
    return StyledComponent;
  },
  accentGrey: 'gray',
  createTheme: () => ({ breakpoints: {} }),
  useTheme: () => ({
    palette: { error: { dark: 'darkred' }, common: { white: 'white' } },
  }),
  TableCell: ({ children }) => <div>{children}</div>,
  TableRow: ({ children }) => <div>{children}</div>,
  Popover: ({ open, children }) => (open ? <div>{children}</div> : null),
  DeleteIcon: () => null,
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
  default: () => null,
}));

vi.mock('../../utils/utils', () => ({
  getVisibilityColums: (columns) => columns,
  getColumnValue: () => undefined,
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => vi.fn(),
}));

vi.mock('@/utils/hooks/useGrafanaPingHook', () => ({
  default: () => vi.fn(),
}));

vi.mock('@/utils/hooks/usePrometheusPingHook', () => ({
  default: () => vi.fn(),
}));

vi.mock('./ConnectionChip', () => ({
  ConnectionStateChip: () => null,
  TooltipWrappedConnectionChip: () => null,
}));

vi.mock('./common', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));

vi.mock('./metadata', () => ({
  default: () => null,
}));

vi.mock('../../utils/responsive-column', () => ({
  getResponsiveColumnVisibility: (...args) => updateVisibleColumns(...args),
}));

vi.mock('../../utils/dimension', () => ({
  useWindowDimensions: () => ({ width: 1280 }),
}));

vi.mock('../multi-select-wrapper', () => ({
  default: () => null,
}));

vi.mock('../../rtk-query/environments', () => ({
  useAddConnectionToEnvironmentMutation: () => [vi.fn()],
  useGetEnvironmentsQuery: (...args) => getEnvironmentsQuery(...args),
  useRemoveConnectionFromEnvironmentMutation: () => [vi.fn()],
  useSaveEnvironmentMutation: () => [vi.fn()],
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
  useUpdateConnectionByIdMutation: () => [vi.fn()],
  usePerformConnectionActionMutation: () => [vi.fn(() => ({ unwrap: () => Promise.resolve({}) }))],
}));

vi.mock('../../assets/icons/disconnect', () => ({
  default: () => null,
}));

vi.mock('./ConnectionStateTransitionModal', () => ({
  default: React.forwardRef(function ConnectionStateTransitionModalMock(_, ref) {
    React.useImperativeHandle(ref, () => ({ show: vi.fn() }));
    return null;
  }),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector) =>
    selector({
      ui: {
        organization: { id: 'org-1' },
        connectionMetadataState: {
          kubernetes: {
            transitions: ['connected'],
            icon: '/static/img/integrations/kubernetes.svg',
          },
        },
        controllerState: {},
      },
    }),
}));

describe('ConnectionTable filters', () => {
  beforeEach(() => {
    dataTableProps = undefined;
    notify.mockReset();
    routerReplace.mockReset();
    getConnectionsQuery.mockReset();
    getEnvironmentsQuery.mockReset();
    updateVisibleColumns.mockReset();
    router.query = {};
    routerReplace.mockImplementation((url: { query?: Record<string, unknown> }) => {
      router.query = { ...(url?.query ?? {}) };
      return Promise.resolve(true);
    });

    getConnectionsQuery.mockReturnValue({
      data: {
        connections: [
          {
            id: 'connection-1',
            name: 'cluster-a',
            kind: 'kubernetes',
            status: 'connected',
            type: 'cluster',
            subType: 'managed',
            metadata: { name: 'cluster-a' },
            environments: [],
            createdAt: '2026-05-08T12:00:00Z',
            updatedAt: '2026-05-09T12:00:00Z',
          },
          {
            id: 'connection-2',
            name: 'cluster-b',
            kind: 'kubernetes',
            status: 'disconnected',
            type: 'cluster',
            subType: 'managed',
            metadata: { name: 'cluster-b' },
            environments: [],
            createdAt: '2026-05-08T12:00:00Z',
            updatedAt: '2026-05-09T12:00:00Z',
          },
        ],
        totalCount: 2,
      },
      isError: false,
      error: undefined,
      refetch: vi.fn(),
      isLoading: false,
    });

    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [] },
      isSuccess: true,
      isError: false,
      error: undefined,
    });
    updateVisibleColumns.mockImplementation((columnNames) =>
      Object.fromEntries(columnNames.map((name) => [name, true])),
    );
  });

  it('clears status/kind from the connections query when filters are reset to All', async () => {
    const user = userEvent.setup();
    const { rerender } = render(<ConnectionTable />);

    await user.click(screen.getByTestId('apply-status-connected'));
    rerender(<ConnectionTable />);

    await waitFor(() => {
      expect(getConnectionsQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ status: 'connected' }),
        undefined,
      );
    });

    await user.click(screen.getByTestId('apply-status-all'));
    rerender(<ConnectionTable />);

    await waitFor(() => {
      expect(getConnectionsQuery).toHaveBeenLastCalledWith(
        expect.objectContaining({ status: undefined, kind: undefined }),
        undefined,
      );
    });
    expect(dataTableProps?.data).toHaveLength(2);
  });
});
