import '@testing-library/jest-dom';
import { configureStore } from '@reduxjs/toolkit';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { ThemeProvider, createTheme } from '@sistent/sistent';
import WorkspaceDataTable from './WorkspaceDataTable';
import { describe, it, expect, vi } from 'vitest';

const store = configureStore({ reducer: {} });
import { useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { useGetSelectedOrganization } from '@/rtk-query/user';

vi.mock('@/rtk-query/workspace', () => ({
  ...vi.importActual('@/rtk-query/workspace'),
  useGetWorkspacesQuery: vi.fn(),
}));

vi.mock('@/rtk-query/user', () => ({
  ...vi.importActual('@/rtk-query/user'),
  useGetSelectedOrganization: vi.fn(),
}));

vi.mock('@sistent/sistent', () => {
  const actual = vi.importActual('@sistent/sistent');
  return {
    ...actual,
    useHasPermission: () => true,
    useWindowDimensions: () => ({ width: 1024, height: 768 }),
  };
});
// Mock performance SVG icon to avoid transform errors
vi.mock('../../public/static/img/drawer-icons/performance_svg', () => ({
  __esModule: true,
  default: () => null,
}));
// Mock WorkspaceActionList to expose props for verification
vi.mock('./WorkspaceActionList', () => (props) => {
  const { workspaceId, workspaceName, selectedWorkspace } = props;
  return (
    <div
      data-testid="workspace-action-list"
      data-workspace-id={workspaceId}
      data-workspace-name={workspaceName}
    >
      {JSON.stringify(selectedWorkspace)}
    </div>
  );
});

describe('WorkspaceDataTable crash evidence', () => {
  it('should demonstrate the rowIndex bug', async () => {
    const mockWorkspaces = [
      { id: '3', name: 'Gamma', owner_id: 'c', description: 'desc 3', updated_at: '2023-01-03' },
      { id: '1', name: 'Alpha', owner_id: 'a', description: 'desc 1', updated_at: '2023-01-01' },
      { id: '2', name: 'Beta', owner_id: 'b', description: 'desc 2', updated_at: '2023-01-02' },
    ];

    (useGetSelectedOrganization as jest.Mock).mockReturnValue({
      selectedOrganization: { id: 'org-1' },
    });

    (useGetWorkspacesQuery as jest.Mock).mockReturnValue({
      data: { workspaces: mockWorkspaces, totalCount: 3 },
      isLoading: false,
    });

    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    render(
      <Provider store={store}>
        <ThemeProvider theme={createTheme()}>
          <WorkspaceDataTable
            handleWorkspaceModalOpen={jest.fn()}
            handleTeamsModalOpen={jest.fn()}
            handleActivityModalOpen={jest.fn()}
            handleDeleteWorkspaceConfirm={jest.fn()}
            columnVisibility={{ actions: true }}
            selectedWorkspace={{ id: '', name: '' }}
            handleRowClick={jest.fn()}
            setColumnVisibility={jest.fn()}
            search=""
            viewType="table"
          />
        </ThemeProvider>
      </Provider>,
    );

    await waitFor(() => {
      expect(screen.getByText('Alpha')).toBeInTheDocument();
    });

    // In a default render, MUIDataTable renders the original array, so rowIndex matches dataIndex.
    // Let's sort the table visually by clicking "Name" header
    const nameHeader = screen.getByText('Name');
    fireEvent.click(nameHeader);

    // MUIDataTable will sort the visible page.
    await new Promise((r) => setTimeout(r, 100)); // wait for sort to apply

    // Verify that each rendered WorkspaceActionList receives the correct workspace via ID
    const actionLists = screen.getAllByTestId('workspace-action-list');
    expect(actionLists).toHaveLength(3);
    // After sorting by name ascending, rows should be Alpha, Beta, Gamma
    expect(actionLists[1]).toHaveAttribute('data-workspace-id', '2');
    const selectedWorkspace = JSON.parse(actionLists[1].textContent);
    expect(selectedWorkspace.id).toBe('2');

    consoleSpy.mockRestore();
  });
});
