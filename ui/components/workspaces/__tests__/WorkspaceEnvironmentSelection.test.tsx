import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import WorkspaceEnvironmentSelection from '../WorkspaceEnvironmentSelection';

describe('WorkspaceEnvironmentSelection', () => {
  const mockWorkspaceId = 'test-workspace-123';
  const mockAssignMutation = vi.fn().mockReturnValue([vi.fn()]);
  const mockUnassignMutation = vi.fn().mockReturnValue([vi.fn()]);
  const mockNotificationHandlers = vi.fn().mockReturnValue({
    handleSuccess: vi.fn(),
    handleError: vi.fn(),
    handleInfo: vi.fn(),
    handleWarn: vi.fn(),
    notifyApiError: vi.fn(),
  });

  it('renders with placeholder and popup caret icon present', () => {
    const mockGetQuery = vi.fn().mockImplementation(({ filter }) => {
      if (filter) {
        return {
          data: {
            environments: [
              { id: 'env-1', name: 'Production' },
              { id: 'env-2', name: 'Staging' },
            ],
          },
          isLoading: false,
        };
      }
      return {
        data: {
          environments: [],
        },
        isLoading: false,
      };
    });

    render(
      <WorkspaceEnvironmentSelection
        workspaceId={mockWorkspaceId}
        useAssignEnvironmentToWorkspaceMutation={mockAssignMutation}
        useGetEnvironmentsOfWorkspaceQuery={mockGetQuery}
        useUnassignEnvironmentFromWorkspaceMutation={mockUnassignMutation}
        useNotificationHandlers={mockNotificationHandlers}
        isAssignedEnvironmentAllowed={true}
      />,
    );

    expect(screen.getByPlaceholderText('Assigned Environment')).toBeInTheDocument();
    // Verify the popup indicator caret button is present
    const openButton = screen.getByRole('button', { name: /open/i });
    expect(openButton).toBeInTheDocument();
  });

  it('renders assigned environments as chips with dropdown caret intact', () => {
    const mockGetQuery = vi.fn().mockImplementation(({ filter }) => {
      if (filter) {
        return {
          data: {
            environments: [{ id: 'env-2', name: 'Staging' }],
          },
          isLoading: false,
        };
      }
      return {
        data: {
          environments: [{ id: 'env-1', name: 'Production' }],
        },
        isLoading: false,
      };
    });

    render(
      <WorkspaceEnvironmentSelection
        workspaceId={mockWorkspaceId}
        useAssignEnvironmentToWorkspaceMutation={mockAssignMutation}
        useGetEnvironmentsOfWorkspaceQuery={mockGetQuery}
        useUnassignEnvironmentFromWorkspaceMutation={mockUnassignMutation}
        useNotificationHandlers={mockNotificationHandlers}
        isAssignedEnvironmentAllowed={true}
      />,
    );

    expect(screen.getByText('Production')).toBeInTheDocument();
    const openButton = screen.getByRole('button', { name: /open/i });
    expect(openButton).toBeInTheDocument();
  });
});
