import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);
let width = 1280;

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: {
    WorkspaceManagementEditWorkspace: { id: 'edit', function: 'workspace' },
    WorkspaceManagementDeleteWorkspace: { id: 'delete', function: 'workspace' },
  },
}));

vi.mock('@sistent/sistent', () => ({
  AccessTimeFilledIcon: () => <svg data-testid="time-icon" />,
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid={`tooltip-${title}`}>{children}</div>
  ),
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  EditIcon: () => <svg data-testid="edit-icon" />,
  GroupAddIcon: () => <svg data-testid="group-add-icon" />,
  IconButton: ({ children, onClick, disabled, ...props }: any) => (
    <button onClick={onClick} disabled={disabled} {...props}>
      {children}
    </button>
  ),
  ListItemIcon: ({ children }: any) => <span>{children}</span>,
  Menu: ({ children, open }: any) => (open ? <div data-testid="menu">{children}</div> : null),
  MenuItem: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={disabled} data-testid="menu-item">
      {children}
    </button>
  ),
  MoreVertIcon: () => <svg data-testid="more-vert-icon" />,
  useTheme: () => ({ palette: { icon: { default: 'icon-default' } } }),
  useWindowDimensions: () => ({ width }),
}));

vi.mock('./styles', () => ({
  TableIconsContainer: ({ children }: any) => <span data-testid="icons-container">{children}</span>,
  IconWrapper: ({ children }: any) => <div data-testid="icon-wrapper">{children}</div>,
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('.', () => ({
  WORKSPACE_ACTION_TYPES: { CREATE: 'create', EDIT: 'edit' },
}));

import WorkspaceActionList from './WorkspaceActionList';

describe('WorkspaceActionList', () => {
  const handleTeamsModalOpen = vi.fn();
  const handleActivityModalOpen = vi.fn();
  const handleWorkspaceModalOpen = vi.fn();
  const handleDeleteWorkspaceConfirm = vi.fn();
  const workspace = { id: 'ws-1', name: 'devspace' };

  beforeEach(() => {
    handleTeamsModalOpen.mockReset();
    handleActivityModalOpen.mockReset();
    handleWorkspaceModalOpen.mockReset();
    handleDeleteWorkspaceConfirm.mockReset();
    can.mockReset();
    can.mockReturnValue(true);
    width = 1280;
  });

  const renderComponent = () =>
    render(
      <WorkspaceActionList
        handleTeamsModalOpen={handleTeamsModalOpen}
        handleActivityModalOpen={handleActivityModalOpen}
        handleWorkspaceModalOpen={handleWorkspaceModalOpen}
        handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
        workspaceId={workspace.id}
        workspaceName={workspace.name}
        selectedWorkspace={workspace}
      />,
    );

  it('renders four direct-action buttons on desktop widths', () => {
    renderComponent();
    // 4 actions on desktop: teams, activity, edit, delete
    const buttons = screen.getAllByRole('button');
    expect(buttons.length).toBe(4);
    expect(screen.getByTestId('tooltip-Assign Teams')).toBeInTheDocument();
    expect(screen.getByTestId('tooltip-Recent Activity')).toBeInTheDocument();
    expect(screen.getByTestId('tooltip-Edit Workspace')).toBeInTheDocument();
    expect(screen.getByTestId('tooltip-Delete Workspace')).toBeInTheDocument();
  });

  it('calls handleTeamsModalOpen with the workspace id and name', async () => {
    const user = userEvent.setup();
    renderComponent();

    const teamsButton = screen.getByRole('button', { name: /assign-teams/i });
    await user.click(teamsButton);
    expect(handleTeamsModalOpen).toHaveBeenCalledWith(expect.anything(), 'ws-1', 'devspace');
  });

  it('calls handleActivityModalOpen when activity button is clicked', async () => {
    const user = userEvent.setup();
    renderComponent();
    await user.click(screen.getByRole('button', { name: /recent-activity/i }));
    expect(handleActivityModalOpen).toHaveBeenCalledWith(expect.anything(), 'ws-1', 'devspace');
  });

  it('calls handleWorkspaceModalOpen with the edit action and selected workspace', async () => {
    const user = userEvent.setup();
    renderComponent();
    await user.click(screen.getByRole('button', { name: /edit-workspace/i }));
    expect(handleWorkspaceModalOpen).toHaveBeenCalledWith(expect.anything(), 'edit', workspace);
  });

  it('calls handleDeleteWorkspaceConfirm when delete button is clicked', async () => {
    const user = userEvent.setup();
    renderComponent();
    await user.click(screen.getByRole('button', { name: /delete-workspace/i }));
    expect(handleDeleteWorkspaceConfirm).toHaveBeenCalledWith(expect.anything(), workspace);
  });

  it('disables edit/delete when permission CAN returns false', () => {
    can.mockImplementation((action: string) => {
      if (action === 'edit' || action === 'delete') return false;
      return true;
    });
    renderComponent();
    expect(screen.getByRole('button', { name: /edit-workspace/i })).toBeDisabled();
    expect(screen.getByRole('button', { name: /delete-workspace/i })).toBeDisabled();
  });

  it('switches to a mobile menu when the viewport is narrow', async () => {
    width = 800;
    const user = userEvent.setup();
    renderComponent();

    expect(screen.queryByTestId('menu')).not.toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /more/i }));
    expect(screen.getByTestId('menu')).toBeInTheDocument();
    // Each of the 4 actions should be rendered as a menu item
    expect(screen.getAllByTestId('menu-item')).toHaveLength(4);
  });
});
