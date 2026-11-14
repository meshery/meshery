import React from 'react';
import '@testing-library/jest-dom';
import { configureStore } from '@reduxjs/toolkit';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { Provider } from 'react-redux';
import { SistentThemeProvider } from '@sistent/sistent';
import WorkspaceDataTable from './WorkspaceDataTable';
import WorkspaceActionList from './WorkspaceActionList';
import { ResponsiveDataTable } from '@sistent/sistent';
import { describe, it, expect, vi, type Mock } from 'vitest';
import { useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { useGetSelectedOrganization } from '@/rtk-query/user';

const store = configureStore({
  reducer: { ui: (state = { organization: { id: 'org-1' } }) => state },
});

// Mock RTK Query hooks
vi.mock('@/rtk-query/workspace', () => ({
  useGetWorkspacesQuery: vi.fn(),
  useGetEnvironmentsOfWorkspaceQuery: vi.fn(() => ({
    data: { environments: [] },
    isLoading: false,
  })),
  useAssignEnvironmentToWorkspaceMutation: vi.fn(() => [vi.fn()]),
  useUnassignEnvironmentFromWorkspaceMutation: vi.fn(() => [vi.fn()]),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetSelectedOrganization: vi.fn(),
}));

// =========================================================================
// VITE STATIC ANALYSIS MOCKS
// =========================================================================
// The following mocks are genuinely required because Vite statically crawls
// the entire import graph of WorkspaceDataTable before executing vi.mock.
// WorkspaceDataTable imports WorkSpaceContentDataTable, which imports a deep
// tree of Next.js components (OrgSwitcher, etc.). Vite crashes on Next.js
// absolute imports (utils/*, assets/*) and JSX inside .js files (drawer-icons).
vi.mock('core-js-pure/actual/disposable-stack', () => ({}));
vi.mock('@mui/x-tree-view', () => ({ TreeItem: () => null, TreeView: () => null }));
vi.mock('@mui/x-date-pickers', () => ({}));
vi.mock('@mui/x-data-grid', () => ({ DataGrid: () => null }));
vi.mock('assets/icons/OrgIcon', () => ({ default: () => null }));
vi.mock('rtk-query/organization', () => ({ useGetOrgsQuery: () => ({}) }));
vi.mock('lib/event-types', () => ({ EVENT_TYPES: {} }));
vi.mock('@/utils/hooks/useNotification', () => ({
  useNotificationHandlers: () => ({
    notifyApiError: vi.fn(),
    handleSuccess: vi.fn(),
    handleError: vi.fn(),
  }),
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('utils/hooks/useNotification', () => ({
  useNotificationHandlers: () => ({
    notifyApiError: vi.fn(),
    handleSuccess: vi.fn(),
    handleError: vi.fn(),
  }),
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('@/components/lifecycle', () => ({ WorkspacesComponent: () => null }));
vi.mock('../../pages/_app', () => ({ mesheryExtensionRoute: '/' }));
vi.mock('../../public/static/img/drawer-icons/application_svg.js', () => ({ default: () => null }));
vi.mock('../../public/static/img/drawer-icons/configuration_hover_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/conformance_svg.js', () => ({ default: () => null }));
vi.mock('../../public/static/img/drawer-icons/discuss_forum_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/extensions_svg.js', () => ({ default: () => null }));
vi.mock('../../public/static/img/drawer-icons/filter_svg.js', () => ({ default: () => null }));
vi.mock('../../public/static/img/drawer-icons/lifecycle_hover_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/lifecycle_mgmt_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/pattern_svg.js', () => ({ default: () => null }));
vi.mock('../../public/static/img/drawer-icons/performance-icon_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/performance_hover_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/performance_svg.js', () => ({ default: () => null }));
vi.mock('../../public/static/img/drawer-icons/servicemeshinterface-icon-white_svg.js', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/smp-white-text_svg.js', () => ({
  default: () => null,
}));
// =========================================================================

// Stop the deep dependency chain at runtime to keep rendering fast
vi.mock('./WorkSpaceContentDataTable', () => ({ default: () => null }));
vi.mock('./WorkspaceActionList', () => ({ default: vi.fn(() => null) }));

// Intercept ResponsiveDataTable to test customBodyRender directly
vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    useHasPermission: () => true,
    useWindowDimensions: () => ({ width: 1024, height: 768 }),
    ResponsiveDataTable: vi.fn((props) => {
      let renderedActions = null;
      const actionsCol = props.columns.find((c: any) => c.name === 'actions');
      if (actionsCol && actionsCol.options.customBodyRender) {
        // Render it with a simulated row index and row data
        renderedActions = actionsCol.options.customBodyRender(null, {
          rowIndex: 0,
          rowData: ['1', 'Alpha', 'a@b.com', 'desc', 0, 0, 0, 0, 'a', null, null, null, null],
        });
      }
      return <div data-testid="mock-responsive-data-table">{renderedActions}</div>;
    }),
  };
});

describe('WorkspaceDataTable customBodyRender', () => {
  it('resolves the correct workspace by ID rather than visual row index', () => {
    const mockWorkspaces = [
      { id: '2', name: 'Beta', ownerId: 'b' }, // index 0 visually
      { id: '1', name: 'Alpha', ownerId: 'a' }, // index 1 visually
    ];

    (useGetSelectedOrganization as Mock).mockReturnValue({
      selectedOrganization: { id: 'org-1' },
    });

    (useGetWorkspacesQuery as Mock).mockReturnValue({
      data: { workspaces: mockWorkspaces, totalCount: 2 },
      isLoading: false,
    });

    render(
      <Provider store={store}>
        <SistentThemeProvider initialMode="light">
          <WorkspaceDataTable
            handleWorkspaceModalOpen={vi.fn()}
            handleTeamsModalOpen={vi.fn()}
            handleActivityModalOpen={vi.fn()}
            handleDeleteWorkspaceConfirm={vi.fn()}
            columnVisibility={{ actions: true }}
            selectedWorkspace={{ id: '', name: '' }}
            handleRowClick={vi.fn()}
            setColumnVisibility={vi.fn()}
            search=""
            viewType="table"
          />
        </SistentThemeProvider>
      </Provider>,
    );

    // Asserts that customBodyRender properly used the ID from rowData ('1' Alpha)
    // rather than the visual rowIndex (0 which corresponds to '2' Beta in the data array)
    expect(WorkspaceActionList).toHaveBeenCalledWith(
      expect.objectContaining({
        workspaceId: '1',
        selectedWorkspace: expect.objectContaining({
          id: '1',
          name: 'Alpha',
        }),
      }),
      undefined,
    );
  });
});

describe('WorkspaceDataTable search paging', () => {
  const renderTable = (search: string) => (
    <Provider store={store}>
      <SistentThemeProvider initialMode="light">
        <WorkspaceDataTable
          handleWorkspaceModalOpen={vi.fn()}
          handleTeamsModalOpen={vi.fn()}
          handleActivityModalOpen={vi.fn()}
          handleDeleteWorkspaceConfirm={vi.fn()}
          columnVisibility={{ actions: true }}
          selectedWorkspace={{ id: '', name: '' }}
          handleRowClick={vi.fn()}
          setColumnVisibility={vi.fn()}
          search={search}
          viewType="table"
        />
      </SistentThemeProvider>
    </Provider>
  );

  const lastQueryArgs = () => (useGetWorkspacesQuery as Mock).mock.calls.at(-1)?.[0];

  it('goes back to the first page when the search changes', () => {
    (useGetSelectedOrganization as Mock).mockReturnValue({
      selectedOrganization: { id: 'org-1' },
    });
    (useGetWorkspacesQuery as Mock).mockReturnValue({
      data: { workspaces: [{ id: '1', name: 'Alpha', ownerId: 'a' }], totalCount: 30 },
      isLoading: false,
    });

    const { rerender } = render(renderTable(''));

    const { onTableChange } = (ResponsiveDataTable as unknown as Mock).mock.calls.at(-1)[0].options;
    act(() => {
      onTableChange('changePage', { page: 2, selectedRows: { data: [] } });
    });
    expect(lastQueryArgs()).toMatchObject({ page: 2, search: '' });

    rerender(renderTable('alpha'));

    expect(lastQueryArgs()).toMatchObject({ page: 0, search: 'alpha' });
  });
});

// The workspaces list is the schemas v1beta3 contract: every field is camelCase.
const WORKSPACE_ROW = {
  id: 'ws-1',
  name: 'Alpha',
  description: 'desc',
  orgName: 'Acme',
  owner: 'Jane',
  ownerId: 'owner-1',
  ownerEmail: 'jane@example.com',
  ownerAvatar: 'https://avatars.example.com/jane.png',
  createdAt: '2026-01-01T00:00:00Z',
  updatedAt: '2026-01-02T00:00:00Z',
};

describe('WorkspaceDataTable column wiring', () => {
  const renderWithRow = (row: Record<string, unknown>, setColumnVisibility = vi.fn()) => {
    (useGetSelectedOrganization as Mock).mockReturnValue({
      selectedOrganization: { id: 'org-1' },
    });
    (useGetWorkspacesQuery as Mock).mockReturnValue({
      data: { workspaces: [row], totalCount: 1 },
      isLoading: false,
    });
    render(
      <Provider store={store}>
        <SistentThemeProvider initialMode="light">
          <WorkspaceDataTable
            handleWorkspaceModalOpen={vi.fn()}
            handleTeamsModalOpen={vi.fn()}
            handleActivityModalOpen={vi.fn()}
            handleDeleteWorkspaceConfirm={vi.fn()}
            columnVisibility={{}}
            selectedWorkspace={{ id: '', name: '' }}
            handleRowClick={vi.fn()}
            setColumnVisibility={setColumnVisibility}
            search=""
            viewType="table"
          />
        </SistentThemeProvider>
      </Provider>,
    );
    const { columns } = (ResponsiveDataTable as unknown as Mock).mock.calls.at(-1)[0];
    // mui-datatables reads each cell from the row by its column name.
    const rowData = columns.map((c: any) => row[c.name]);
    const ownerCell = columns
      .find((c: any) => c.name === 'owner')
      .options.customBodyRender(row.owner, { rowIndex: 0, rowData });
    return { columns, ownerCell };
  };

  it('names every column after a field of the camelCase workspace row', () => {
    const { columns } = renderWithRow(WORKSPACE_ROW);
    expect(columns.map((c: any) => c.name)).toEqual([
      'id',
      'ownerEmail',
      'ownerId',
      'orgName',
      'name',
      'description',
      'ownerAvatar',
      'owner',
      'environments',
      'designCount',
      'viewCount',
      'teamCount',
      'environmentCount',
      'createdAt',
      'updatedAt',
      'actions',
    ]);
  });

  it("shows the owner's avatar and links to the owner's profile", () => {
    const open = vi.spyOn(window, 'open').mockImplementation(() => null);
    const { ownerCell } = renderWithRow(WORKSPACE_ROW);
    const { container } = render(
      <SistentThemeProvider initialMode="light">{ownerCell}</SistentThemeProvider>,
    );

    fireEvent.click(container.querySelector('.MuiAvatar-root')!);
    expect(open).toHaveBeenCalledWith(expect.stringMatching(/\/user\/owner-1$/), '_blank');
    expect(screen.getByRole('img', { name: 'Jane' })).toHaveAttribute(
      'src',
      WORKSPACE_ROW.ownerAvatar,
    );
    open.mockRestore();
  });

  it('falls back to the org name for a workspace without an owner', () => {
    const { ownerCell } = renderWithRow({ ...WORKSPACE_ROW, owner: undefined });
    render(ownerCell);
    expect(screen.getByText('Acme')).toBeInTheDocument();
  });

  it('seeds the responsive visibility under the same column names', () => {
    const setColumnVisibility = vi.fn();
    renderWithRow(WORKSPACE_ROW, setColumnVisibility);
    // width is mocked to 1024: 'na' columns start hidden, 'l' columns shown.
    expect(setColumnVisibility).toHaveBeenCalledWith(
      expect.objectContaining({ ownerEmail: false, createdAt: false, updatedAt: true }),
    );
  });
});
