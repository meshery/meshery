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

// vi.mock factories are hoisted above these declarations, so the fixtures they
// close over must be hoisted too (vi.hoisted) or they are still in the temporal
// dead zone when a factory runs. Mocked hook results must also be referentially
// stable: the component keys effects and memos on these objects, so a fresh
// object per render is a render loop.
const {
  notify,
  createWorkspace,
  updateWorkspace,
  deleteWorkspace,
  WORKSPACES_QUERY_RESULT,
  EMPTY_QUERY_RESULT,
  NOOP_MUTATION,
  UI_STATE,
} = vi.hoisted(() => ({
  notify: vi.fn(),
  createWorkspace: vi.fn(),
  updateWorkspace: vi.fn(),
  deleteWorkspace: vi.fn(),
  WORKSPACES_QUERY_RESULT: { data: { workspaces: [], totalCount: 0 } },
  EMPTY_QUERY_RESULT: { data: undefined, isLoading: false },
  NOOP_MUTATION: [vi.fn()],
  UI_STATE: { ui: { organization: { id: 'org-1' } } },
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
  useNotificationHandlers: () => ({
    handleSuccess: vi.fn(),
    handleError: vi.fn(),
    handleInfo: vi.fn(),
    notifyApiError: vi.fn(),
  }),
}));

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
  // This suite exercises the create flow, not authorization: grant every
  // capability so the permission gates never mask the behaviour under test.
  useHasPermission: () => true,
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

  // The create modal (its Save button) must survive a failed create so the
  // user keeps their typed input, and disappear only once the create resolves.
  // A deferred promise lets us observe the modal between submit and settlement.
  it('keeps the create modal open on failure and closes it only on success', async () => {
    let rejectCreate: (reason?: unknown) => void;
    createWorkspace.mockReturnValue({
      unwrap: () =>
        new Promise((_resolve, reject) => {
          rejectCreate = reject;
        }),
    });

    await openCreateModalAndSubmit();
    // Still pending: modal open.
    expect(screen.getByTestId('submit-workspace')).toBeInTheDocument();

    rejectCreate({ status: 403, data: { error: 'nope' }, meshkit: { message: 'nope' } });
    await waitFor(() => expect(notify).toHaveBeenCalled());
    // Rejected: modal must remain open.
    expect(screen.getByTestId('submit-workspace')).toBeInTheDocument();
    expect(notifiedEventTypes()).not.toContain('success');

    // A subsequent create that resolves must close the modal.
    createWorkspace.mockReturnValue({ unwrap: () => Promise.resolve({ id: 'ws-1' }) });
    await userEvent.setup().click(screen.getByTestId('submit-workspace'));
    await waitFor(() => expect(screen.queryByTestId('submit-workspace')).not.toBeInTheDocument());
  });
});
