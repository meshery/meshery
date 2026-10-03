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
  submittedOrgId,
} = vi.hoisted(() => ({
  notify: vi.fn(),
  createWorkspace: vi.fn(),
  updateWorkspace: vi.fn(),
  deleteWorkspace: vi.fn(),
  WORKSPACES_QUERY_RESULT: { data: { workspaces: [], totalCount: 0 }, isLoading: false },
  EMPTY_QUERY_RESULT: { data: undefined, isLoading: false },
  NOOP_MUTATION: [vi.fn()],
  UI_STATE: { ui: { organization: { id: 'org-1' } as { id: string } | null } },
  // Boxed so tests can override the id the mocked form "submits" without
  // reassigning the hoisted binding itself (vi.hoisted values are const).
  submittedOrgId: { current: 'org-1' as string | undefined },
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
  useCreateWorkspaceMutation: () => [createWorkspace, { isLoading: false }],
  useUpdateWorkspaceMutation: () => [updateWorkspace, { isLoading: false }],
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
  RJSFModalWrapper: ({ handleSubmit, isSubmitting }: any) => (
    <button
      data-testid="submit-workspace"
      disabled={isSubmitting}
      onClick={() =>
        handleSubmit({
          organizationId: submittedOrgId.current,
          name: 'team-space',
          description: 'shared',
        })
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
  CircularProgress: (props: any) => <div {...props} />,
  CustomColumnVisibilityControl: () => null,
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid={`tooltip-${title}`}>{children}</div>
  ),
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  Modal: ({ children, open }: any) => (open ? <div>{children}</div> : null),
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  NoSsr: ({ children }: any) => <>{children}</>,
  PROMPT_VARIANTS: { DANGER: 'danger' },
  Select: ({ children, ...rest }: any) => <select {...rest}>{children}</select>,
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
  // `DefaultError` (rendered by the new permission guard) imports `styled`
  // from error-404/styles.tsx to build its `ErrorMain` wrapper. Emulate the
  // real styled-components-style API closely enough that a tag-name call
  // (`styled('main')`) or component call (`styled(SomeComponent)`) both
  // return a renderable element, ignoring the style function/object passed
  // to the second call.
  styled: (Component: any) => () => {
    if (typeof Component === 'string') {
      return ({ children, ...rest }: any) => <Component {...rest}>{children}</Component>;
    }
    return Component;
  },
}));

vi.mock('../general/error-404/index', () => ({
  default: ({ permissionKey }: any) => <div data-testid="default-error">{permissionKey}</div>,
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
    // @meshery/schemas baseQuery wrapper attaches - since v1.3.37
    // (meshery/schemas#1081) it carries the full envelope, reading the
    // server's camelCase detail arrays with a snake_case fallback, so the
    // probable cause and remediation list arrive populated. The real
    // transform is pinned in
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
            probableCause: ['Your account does not have permission to create workspaces.'],
            suggestedRemediation: ['Ask an organization owner to grant the Workspace role.'],
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

/**
 * Regression coverage for meshery/meshery#21263.
 *
 * organizationId is a hidden field in the create-workspace form, seeded from
 * `organization?.id` when the modal opens. If the org context hasn't
 * hydrated yet, the user has no way to see or fix it, and the request used to
 * reach the server with a missing/empty organizationId - which only surfaced
 * as an opaque "Unable to unmarshal the : workspace" error. The Create button
 * is disabled (with a tooltip) so that doomed click can't happen at all; the
 * submit-time guard below is the remaining defense-in-depth for the other
 * entry point into handleWorkspaceModalOpen (WorkspaceSwitcher's "+ Create
 * Workspace", which does not go through this button).
 */
describe('Workspaces create flow — organization guard', () => {
  beforeEach(() => {
    notify.mockReset();
    createWorkspace.mockReset();
  });

  afterEach(() => {
    UI_STATE.ui.organization = { id: 'org-1' };
    submittedOrgId.current = 'org-1';
  });

  it('disables the Create button with an explanatory tooltip while the organization has not loaded yet', async () => {
    UI_STATE.ui.organization = null;
    render(<Workspaces onSelectWorkspace={undefined} />);

    const createButton = screen.getByText('Create').closest('button');
    expect(createButton).toBeDisabled();
    expect(screen.getByTestId(/^tooltip-Organization is still loading/)).toBeInTheDocument();

    // A disabled native button does not dispatch click at all - the doomed
    // action must be unreachable, not merely toasted after the fact.
    await userEvent.setup().click(screen.getByText('Create'));

    expect(screen.queryByTestId('submit-workspace')).not.toBeInTheDocument();
    expect(createWorkspace).not.toHaveBeenCalled();
    expect(notify).not.toHaveBeenCalled();
  });

  it('submit-time guard rejects creation when organizationId is missing, even with the modal already open', async () => {
    submittedOrgId.current = undefined;

    await openCreateModalAndSubmit();

    expect(createWorkspace).not.toHaveBeenCalled();
    await waitFor(() => expect(notify).toHaveBeenCalled());
    const types = notifiedEventTypes();
    expect(types).toContain('error');
    expect(types).not.toContain('success');
    expect(notify.mock.calls[0][0].message).toMatch(/organization/i);
  });
});

/**
 * Flagged in review on meshery/meshery#21335: the workspaces query is
 * skipped until organization?.id hydrates, and a skipped RTK Query hook
 * never reports isLoading - so `workspacesData` stayed undefined and
 * `workspaces` defaulted to [], rendering EmptyState's "Click Create" prompt
 * at the exact moment the Create button was disabled for "Organization is
 * still loading…". Neither EmptyState nor the data views should render
 * before organization and the workspaces query have both resolved.
 */
describe('Workspaces list — loading state', () => {
  afterEach(() => {
    UI_STATE.ui.organization = { id: 'org-1' };
    WORKSPACES_QUERY_RESULT.data = { workspaces: [], totalCount: 0 };
    WORKSPACES_QUERY_RESULT.isLoading = false;
  });

  it('shows a loading indicator instead of the empty state while the organization has not loaded yet', () => {
    UI_STATE.ui.organization = null;
    render(<Workspaces onSelectWorkspace={undefined} />);

    expect(screen.getByTestId('workspaces-loading')).toBeInTheDocument();
    expect(screen.queryByText('empty')).not.toBeInTheDocument();
  });

  it('shows a loading indicator instead of the empty state while the workspaces query is still in flight', () => {
    WORKSPACES_QUERY_RESULT.data = undefined;
    WORKSPACES_QUERY_RESULT.isLoading = true;
    render(<Workspaces onSelectWorkspace={undefined} />);

    expect(screen.getByTestId('workspaces-loading')).toBeInTheDocument();
    expect(screen.queryByText('empty')).not.toBeInTheDocument();
  });

  it('shows the empty state once organization and workspace data have both resolved with zero workspaces', () => {
    render(<Workspaces onSelectWorkspace={undefined} />);

    expect(screen.queryByTestId('workspaces-loading')).not.toBeInTheDocument();
    expect(screen.getByText('empty')).toBeInTheDocument();
  });
});
