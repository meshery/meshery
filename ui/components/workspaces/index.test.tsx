import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

/**
 * Regression coverage for silent workspace failures.
 *
 * `handleError` was curried - `handleError(action) => (error) => ...` - but
 * every call site invoked it as `handleError('some message')`, which merely
 * built a function and discarded it. No failed workspace create, update or
 * delete ever reached the user. The same defect in the environments component
 * additionally produced a false success toast (see ../environments/index.test.tsx).
 *
 * The contract pinned here: a rejected create emits an ERROR notification
 * carrying the server's MeshKit metadata, and never a SUCCESS one.
 */

const notify = vi.fn();
const createWorkspace = vi.fn();
const updateWorkspace = vi.fn();
const deleteWorkspace = vi.fn();

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
  useNotificationHandlers: () => ({
    handleSuccess: vi.fn(),
    handleError: vi.fn(),
    handleInfo: vi.fn(),
    notifyApiError: vi.fn(),
  }),
}));

// Mocked hook results must be referentially stable: the component keys effects
// and memos on these objects, so a fresh object per render is a render loop.
const WORKSPACES_QUERY_RESULT = { data: { workspaces: [], total_count: 0 } };
const EMPTY_QUERY_RESULT = { data: undefined, isLoading: false };
const NOOP_MUTATION = [vi.fn()];

vi.mock('../../rtk-query/workspace', () => ({
  useCreateWorkspaceMutation: () => [createWorkspace],
  useUpdateWorkspaceMutation: () => [updateWorkspace],
  useDeleteWorkspaceMutation: () => [deleteWorkspace],
  useGetWorkspacesQuery: () => WORKSPACES_QUERY_RESULT,
  useGetTeamsOfWorkspaceQuery: () => EMPTY_QUERY_RESULT,
  useGetEventsOfWorkspaceQuery: () => EMPTY_QUERY_RESULT,
  useAssignTeamToWorkspaceMutation: () => NOOP_MUTATION,
  useUnassignTeamFromWorkspaceMutation: () => NOOP_MUTATION,
}));

vi.mock('@/rtk-query/user', () => ({
  useGetUsersForOrgQuery: () => EMPTY_QUERY_RESULT,
  useRemoveUserFromTeamMutation: () => NOOP_MUTATION,
}));

const UI_STATE = { ui: { organization: { id: 'org-1' } } };

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) => selector(UI_STATE),
}));

vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));
vi.mock('@/utils/can', () => ({ default: () => true }));
vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

// The factory is hoisted above module scope, so it must not close over any
// top-level const - build the default value inline.
vi.mock('@/utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({
    createNewWorkspaceModalOpen: false,
    setCreateNewWorkspaceModalOpen: () => {},
  }),
}));

vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ handleSubmit }: any) => (
    <button
      data-testid="submit-workspace"
      onClick={() =>
        handleSubmit({ organizationId: 'org-1', name: 'team-space', description: 'shared' })
      }
    >
      Save
    </button>
  ),
}));

vi.mock('../general/PromptComponent', () => ({ default: React.forwardRef(() => null) }));
vi.mock('@/components/lifecycle/general', () => ({ EmptyState: () => <div>empty</div> }));
vi.mock('@/components/general/ViewSwitch', () => ({ default: () => null }));
vi.mock('./WorkspaceGridView', () => ({ default: () => null }));
vi.mock('./WorkspaceDataTable', () => ({ default: () => null }));
vi.mock('./styles', () => ({ CreateButtonWrapper: ({ children }: any) => <div>{children}</div> }));
vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('@/assets/icons/AddIconCircleBorder', () => ({ default: () => null }));
vi.mock('@/assets/icons/RightArrowIcon', () => ({ default: () => null }));
vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div>{children}</div>,
  Breadcrumbs: ({ children }: any) => <nav>{children}</nav>,
  Button: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  CustomColumnVisibilityControl: () => null,
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  Modal: ({ children, open }: any) => (open ? <div>{children}</div> : null),
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  NoSsr: ({ children }: any) => <>{children}</>,
  PROMPT_VARIANTS: { DANGER: 'danger' },
  SearchBar: () => null,
  TeamsIcon: () => null,
  Typography: ({ children }: any) => <span>{children}</span>,
  WorkspaceIcon: () => null,
  WorkspaceRecentActivityModal: () => null,
  WorkspaceTeamsTable: () => null,
  createAndEditWorkspaceSchema: {},
  createAndEditWorkspaceUiSchema: {},
  editWorkspaceSchema: {},
  useTheme: () => ({
    palette: {
      icon: { default: '#000', secondary: '#111' },
      text: { default: '#000' },
      common: { white: '#fff' },
      background: { brand: { default: '#000' }, constant: { table: '#fff' } },
    },
  }),
}));

import Workspaces from './index';

const notifiedEventTypes = () =>
  notify.mock.calls.map(([arg]) => arg?.event_type?.type ?? arg?.event_type);

const openCreateModalAndSubmit = async () => {
  const user = userEvent.setup();
  render(<Workspaces onSelectWorkspace={undefined} />);
  await user.click(screen.getByText('Create'));
  await user.click(await screen.findByTestId('submit-workspace'));
};

describe('Workspaces create flow notifications', () => {
  beforeEach(() => {
    notify.mockReset();
    createWorkspace.mockReset();
    updateWorkspace.mockReset();
    deleteWorkspace.mockReset();
  });

  it('surfaces the failure when the provider rejects the create', async () => {
    // `data` is the verbatim server envelope (camelCase, per
    // server/models/httputil/httputil.go); `meshkit` is what the
    // @meshery/schemas baseQuery wrapper actually attaches - message/code/
    // severity only, because it reads snake_case spellings the server does not
    // emit (meshery/schemas#1081). The real transform is pinned in
    // `utils/helpers/__tests__/meshkitErrorChain.test.ts`.
    createWorkspace.mockReturnValue({
      unwrap: () =>
        Promise.reject({
          status: 403,
          data: {
            error: 'Unable to create the workspace',
            code: 'meshery-server-1454',
            severity: 'ALERT',
            probableCause: ['Your account does not have permission to create workspaces.'],
            suggestedRemediation: ['Ask an organization owner to grant the Workspace role.'],
          },
          meshkit: {
            message: 'Unable to create the workspace',
            code: 'meshery-server-1454',
            severity: 'ALERT',
          },
        }),
    });

    await openCreateModalAndSubmit();

    await waitFor(() => expect(notify).toHaveBeenCalled());

    const types = notifiedEventTypes();
    expect(types).toContain('error');
    expect(types).not.toContain('success');

    const [payload] = notify.mock.calls[0];
    expect(payload.message).toContain('Unable to create the workspace');
    expect(payload.message).toContain('meshery-server-1454');
    expect(payload.message).toContain('Ask an organization owner to grant the Workspace role.');
  });

  it('reports success only once the create actually resolves', async () => {
    createWorkspace.mockReturnValue({ unwrap: () => Promise.resolve({ id: 'ws-1' }) });

    await openCreateModalAndSubmit();

    await waitFor(() => expect(notify).toHaveBeenCalled());

    const types = notifiedEventTypes();
    expect(types).toContain('success');
    expect(types).not.toContain('error');
    expect(notify.mock.calls[0][0].message).toBe('Workspace "team-space" created');
  });
});
