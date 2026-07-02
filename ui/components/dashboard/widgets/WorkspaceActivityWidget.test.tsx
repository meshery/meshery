import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let workspacesQueryReturn: {
  data?: { workspaces?: Array<{ id: string; name: string }> };
  isError?: boolean;
  isLoading?: boolean;
} = { data: { workspaces: [] } };

let eventsQueryReturn: {
  data?: { data?: unknown[] };
  isLoading?: boolean;
  isError?: boolean;
} = { data: { data: [] }, isLoading: false, isError: false };

let currentOrg: { id?: string } | undefined = { id: 'org-1' };

const activityCardSpy = vi.fn();

vi.mock('@/rtk-query/workspace', () => ({
  useGetWorkspacesQuery: () => workspacesQueryReturn,
  useGetEventsOfWorkspaceQuery: () => eventsQueryReturn,
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { organization: currentOrg } }),
}));

vi.mock('@sistent/sistent', () => ({
  WorkspaceActivityCard: (props: {
    selectedWorkspace?: string;
    workspaces: Array<{ id: string; name: string }>;
    activities: unknown[];
    isEventsLoading?: boolean;
    workspacePagePath?: string;
    handleWorkspaceChange?: (event: { target: { value: string } }) => void;
  }) => {
    activityCardSpy(props);
    return (
      <div
        data-testid="activity-card"
        data-selected={String(props.selectedWorkspace)}
        data-activities={String(props.activities.length)}
        data-loading={String(props.isEventsLoading)}
      >
        <button
          type="button"
          onClick={() =>
            props.handleWorkspaceChange?.({
              target: { value: 'ws-2' },
            } as never)
          }
        >
          changeWs
        </button>
      </div>
    );
  },
}));

const widgetErrorFallbackSpy = vi.fn();

vi.mock('./WidgetErrorFallback', () => ({
  default: (props: { widgetTitle: string; message?: string }) => {
    widgetErrorFallbackSpy(props);
    return (
      <div data-testid="widget-error-fallback" data-title={props.widgetTitle}>
        {props.message}
      </div>
    );
  },
}));

import WorkspaceActivityWidget from './WorkspaceActivityWidget';

describe('WorkspaceActivityWidget', () => {
  beforeEach(() => {
    activityCardSpy.mockReset();
    widgetErrorFallbackSpy.mockReset();
    workspacesQueryReturn = {
      data: {
        workspaces: [
          { id: 'ws-1', name: 'WS One' },
          { id: 'ws-2', name: 'WS Two' },
        ],
      },
      isError: false,
      isLoading: false,
    };
    eventsQueryReturn = {
      data: { data: [{ id: 'evt-1' }, { id: 'evt-2' }] },
      isLoading: false,
      isError: false,
    };
    currentOrg = { id: 'org-1' };
  });

  it('selects the first workspace by default and forwards events', () => {
    render(<WorkspaceActivityWidget />);
    const card = screen.getByTestId('activity-card');
    expect(card).toHaveAttribute('data-selected', 'ws-1');
    expect(card).toHaveAttribute('data-activities', '2');
  });

  it('updates selected workspace when the change handler is invoked', async () => {
    render(<WorkspaceActivityWidget />);
    expect(screen.getByTestId('activity-card')).toHaveAttribute('data-selected', 'ws-1');
    await userEvent.click(screen.getByRole('button', { name: 'changeWs' }));
    expect(screen.getByTestId('activity-card')).toHaveAttribute('data-selected', 'ws-2');
  });

  it('returns empty activities when events query errors', () => {
    eventsQueryReturn = { isError: true };
    render(<WorkspaceActivityWidget />);
    expect(screen.getByTestId('activity-card')).toHaveAttribute('data-activities', '0');
  });

  it('reflects loading state via card prop', () => {
    eventsQueryReturn = { isLoading: true, data: { data: [] } };
    render(<WorkspaceActivityWidget />);
    expect(screen.getByTestId('activity-card')).toHaveAttribute('data-loading', 'true');
  });

  it('falls back to empty selected workspace when there are no workspaces', () => {
    workspacesQueryReturn = { data: { workspaces: [] } };
    render(<WorkspaceActivityWidget />);
    expect(screen.getByTestId('activity-card')).toHaveAttribute('data-selected', '');
  });

  it('still renders even when the current org has no id', () => {
    currentOrg = {};
    workspacesQueryReturn = { data: { workspaces: [] } };
    render(<WorkspaceActivityWidget />);
    expect(screen.getByTestId('activity-card')).toBeInTheDocument();
  });

  it('shows the error fallback when the workspaces query fails', () => {
    workspacesQueryReturn = { isError: true };
    render(<WorkspaceActivityWidget />);
    expect(screen.getByTestId('widget-error-fallback')).toHaveAttribute(
      'data-title',
      'Workspace Activity',
    );
    expect(screen.queryByTestId('activity-card')).not.toBeInTheDocument();
  });
});
