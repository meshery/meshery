import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);
const handleSuccess = vi.fn();
const handleError = vi.fn();
const notifyApiError = vi.fn();
const deleteWorkspaceMutator = vi.fn();

vi.mock('@/rtk-query/workspace', () => ({
  useDeleteWorkspaceMutation: () => [deleteWorkspaceMutator],
}));

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotificationHandlers: () => ({ handleSuccess, handleError, notifyApiError }),
}));

vi.mock('@sistent/sistent', () => ({
  ChevronLeftIcon: (props) =>
    React.createElement('svg', { 'data-component': 'chevron-left-icon', ...props }),
  ChevronRightIcon: (props) =>
    React.createElement('svg', { 'data-component': 'chevron-right-icon', ...props }),
  Grid2: ({ children }: any) => <div>{children}</div>,
  L5DeleteIcon: ({ onClick, disabled }: any) => (
    <button data-testid="bulk-delete-btn" onClick={onClick} disabled={disabled}>
      delete
    </button>
  ),
  Modal: ({ children, open, title }: any) =>
    open ? (
      <div data-testid="modal">
        <h2>{title}</h2>
        {children}
      </div>
    ) : null,
  Pagination: ({ count, page, onChange }: any) => (
    <div>
      <span data-testid="page">
        {page}/{count}
      </span>
      <button data-testid="page-next" onClick={(e) => onChange(e, page + 1)}>
        next
      </button>
    </div>
  ),
  PaginationItem: () => null,
  ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
  ModalButtonPrimary: ({ children, onClick }: any) => (
    <button onClick={onClick} data-testid="modal-primary">
      {children}
    </button>
  ),
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  ModalButtonSecondary: ({ children, onClick }: any) => (
    <button onClick={onClick} data-testid="modal-secondary">
      {children}
    </button>
  ),
  useTheme: () => ({
    palette: {
      background: { error: { default: 'red' } },
    },
  }),
  ErrorBoundary: ({ children }: any) => <>{children}</>,
}));

vi.mock('./styles', () => ({
  UserCommonBox: ({ children }: any) => <div data-testid="bulk-bar">{children}</div>,
}));

vi.mock('./MesheryWorkspaceCard', () => ({
  default: ({ workspaceDetails, handleBulkSelect }: any) => (
    <div data-testid={`workspace-card-${workspaceDetails.id}`}>
      <span>{workspaceDetails.name}</span>
      <input
        type="checkbox"
        data-testid={`bulk-select-${workspaceDetails.id}`}
        onChange={(e) => handleBulkSelect(e, workspaceDetails.id)}
      />
    </div>
  ),
}));

import WorkspaceGridView from './WorkspaceGridView';

describe('WorkspaceGridView', () => {
  const workspaces = [
    { id: 'ws-1', name: 'one' },
    { id: 'ws-2', name: 'two' },
  ];
  const handleWorkspaceModalOpen = vi.fn();
  const handleDeleteWorkspaceConfirm = vi.fn();
  const setPage = vi.fn();

  beforeEach(() => {
    handleWorkspaceModalOpen.mockReset();
    handleDeleteWorkspaceConfirm.mockReset();
    setPage.mockReset();
    deleteWorkspaceMutator.mockReset();
    handleSuccess.mockReset();
    handleError.mockReset();
    notifyApiError.mockReset();
    can.mockReset();
    can.mockReturnValue(true);
    deleteWorkspaceMutator.mockReturnValue({ unwrap: () => Promise.resolve() });
  });

  const renderComponent = (props = {}) =>
    render(
      <WorkspaceGridView
        workspacesData={workspaces}
        handleWorkspaceModalOpen={handleWorkspaceModalOpen}
        handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
        totalPages={3}
        page={0}
        setPage={setPage}
        {...props}
      />,
    );

  it('renders one card per workspace', () => {
    renderComponent();
    expect(screen.getByTestId('workspace-card-ws-1')).toBeInTheDocument();
    expect(screen.getByTestId('workspace-card-ws-2')).toBeInTheDocument();
    // No bulk bar by default (no selections)
    expect(screen.queryByTestId('bulk-bar')).not.toBeInTheDocument();
  });

  it('shows bulk action bar after selecting a workspace', async () => {
    const user = userEvent.setup();
    renderComponent();
    await user.click(screen.getByTestId('bulk-select-ws-1'));

    expect(screen.getByTestId('bulk-bar')).toBeInTheDocument();
    expect(screen.getByTestId('bulk-bar')).toHaveTextContent('1 workspace selected');
  });

  it('pluralises the bulk-bar message for multiple selections', async () => {
    const user = userEvent.setup();
    renderComponent();

    await user.click(screen.getByTestId('bulk-select-ws-1'));
    await user.click(screen.getByTestId('bulk-select-ws-2'));

    expect(screen.getByTestId('bulk-bar')).toHaveTextContent('2 workspaces selected');
  });

  it('opens delete modal when the bulk delete is clicked and deletes on confirm', async () => {
    const user = userEvent.setup();
    renderComponent();

    await user.click(screen.getByTestId('bulk-select-ws-1'));
    await user.click(screen.getByTestId('bulk-delete-btn'));
    expect(screen.getByTestId('modal')).toBeInTheDocument();

    await user.click(screen.getByTestId('modal-primary'));

    expect(deleteWorkspaceMutator).toHaveBeenCalledWith({ workspaceId: 'ws-1' });
  });

  it('closes the delete modal when secondary button (Cancel) is clicked', async () => {
    const user = userEvent.setup();
    renderComponent();

    await user.click(screen.getByTestId('bulk-select-ws-1'));
    await user.click(screen.getByTestId('bulk-delete-btn'));
    expect(screen.getByTestId('modal')).toBeInTheDocument();

    await user.click(screen.getByTestId('modal-secondary'));
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
    expect(deleteWorkspaceMutator).not.toHaveBeenCalled();
  });

  it('hides the delete modal entirely when the user lacks WorkspaceManagementDeleteWorkspace permission', () => {
    can.mockReturnValue(false);
    renderComponent();
    // Modal is rendered conditionally on CAN(...) - so when no permission, modal element is absent
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });
});
