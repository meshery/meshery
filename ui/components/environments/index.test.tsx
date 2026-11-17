import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

/**
 * Regression coverage for the "Create Environment succeeded" lie.
 *
 * An Org Admin whose create was rejected by the remote provider (HTTP 403) saw
 * a green "Environment created" toast and no environment. Two independent
 * defects produced that:
 *
 *   1. `.then(handleSuccess('...'))` - note the immediate call - evaluated
 *      handleSuccess while the callback was merely being constructed, so the
 *      success notification fired synchronously, before the request settled
 *      and regardless of the outcome.
 *   2. `handleError` was curried, `handleError(msg) => (error) => ...`, but
 *      every call site invoked it as `handleError(msg)`, producing a function
 *      that was thrown away. No failure ever reached the user.
 *
 * These tests pin both halves: success only on resolve, failure always on
 * reject, and the MeshKit metadata the server sends carried into the toast.
 */

// vi.mock factories are hoisted above these declarations, so the fixtures they
// close over must be hoisted too (vi.hoisted) or they are still in the temporal
// dead zone when a factory runs. Query results MUST also be referentially
// stable across renders: the component appends query results into state inside
// effects keyed on the query object identity, so a freshly-built object per
// render is an infinite render loop that exhausts the heap, not a test failure.
const {
  notify,
  getEnvironmentsQuery,
  getEnvironmentConnectionsQuery,
  createEnvironment,
  updateEnvironment,
  deleteEnvironment,
  canSpy,
  permissionOverrides,
  ENVIRONMENTS_QUERY_RESULT,
  CONNECTIONS_QUERY_RESULT,
  NOOP_MUTATION,
  UI_STATE,
} = vi.hoisted(() => {
  const overrides = new Map<string, boolean>();
  return {
    notify: vi.fn(),
    getEnvironmentsQuery: vi.fn(),
    getEnvironmentConnectionsQuery: vi.fn(),
    createEnvironment: vi.fn(),
    updateEnvironment: vi.fn(),
    deleteEnvironment: vi.fn(),
    permissionOverrides: overrides,
    canSpy: vi.fn((key?: { id?: string }) => overrides.get(key?.id ?? '') ?? true),
    ENVIRONMENTS_QUERY_RESULT: {
      data: { environments: [], totalCount: 0 },
      isLoading: false,
      isError: false,
      error: undefined,
    },
    CONNECTIONS_QUERY_RESULT: {
      data: { connections: [], totalCount: 0 },
      currentData: undefined as any,
      isError: false,
      error: undefined,
    },
    NOOP_MUTATION: [vi.fn()],
    UI_STATE: { ui: { organization: { id: 'org-1' } } },
  };
});

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('../../rtk-query/environments', () => ({
  useCreateEnvironmentMutation: () => [createEnvironment, { isLoading: false }],
  useUpdateEnvironmentMutation: () => [updateEnvironment, { isLoading: false }],
  useDeleteEnvironmentMutation: () => [deleteEnvironment],
  useGetEnvironmentsQuery: (args: unknown, options: unknown) => {
    getEnvironmentsQuery(args, options);
    return ENVIRONMENTS_QUERY_RESULT;
  },
  useGetEnvironmentConnectionsQuery: (args: unknown, options: unknown) => {
    getEnvironmentConnectionsQuery(args, options);
    const isSkipped = Boolean((options as any)?.skip);
    const resolvedCurrentData =
      CONNECTIONS_QUERY_RESULT.currentData !== undefined
        ? CONNECTIONS_QUERY_RESULT.currentData
        : isSkipped
          ? undefined
          : CONNECTIONS_QUERY_RESULT.data;
    return {
      data: CONNECTIONS_QUERY_RESULT.data,
      currentData: resolvedCurrentData,
      isError: CONNECTIONS_QUERY_RESULT.isError,
      error: CONNECTIONS_QUERY_RESULT.error,
    };
  },
  useAddConnectionToEnvironmentMutation: () => NOOP_MUTATION,
  useRemoveConnectionFromEnvironmentMutation: () => NOOP_MUTATION,
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) => selector(UI_STATE),
}));

vi.mock('next/router', () => ({
  withRouter: (Component: React.ComponentType) => Component,
}));

vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));

vi.mock('@/utils/can', () => ({ default: () => true }));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: (_, prop) => ({ id: String(prop), function: 'fn' }) }),
}));

// The RJSF form is not under test; expose a button that submits the payload the
// real form would produce.
vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ handleSubmit, isSubmitting }: any) => (
    <button
      data-testid="submit-environment"
      disabled={isSubmitting}
      onClick={() =>
        handleSubmit({ organizationId: 'org-1', name: 'prod', description: 'production' })
      }
    >
      Save
    </button>
  ),
}));

vi.mock('./environment-card', () => ({
  default: ({ onAssignConnection, environmentDetails }: any) => (
    <div data-testid={`env-card-${environmentDetails?.id}`}>
      <button
        data-testid={`assign-btn-${environmentDetails?.id}`}
        onClick={(e) => onAssignConnection(e, environmentDetails)}
      >
        Assign
      </button>
    </div>
  ),
}));
vi.mock('./styles', () => ({
  CreateButtonWrapper: ({ children }: any) => <div>{children}</div>,
  BulkActionWrapper: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('../general/PromptComponent', () => ({
  default: React.forwardRef(() => null),
}));
vi.mock('../lifecycle/general', () => ({ EmptyState: () => <div>empty</div> }));
vi.mock('../general/error-404/index', () => ({ default: () => <div>no access</div> }));
vi.mock('../../assets/icons/AddIconCircleBorder', () => ({ default: () => null }));
vi.mock('../../assets/icons/Environment', () => ({ default: () => null }));
vi.mock('../../assets/icons/Connection', () => ({ default: () => null }));

vi.mock('@sistent/sistent', () => ({
  // This suite exercises the create flow, not authorization: grant every
  // capability so the permission gates never mask the behaviour under test.
  useHasPermission: (key: any) => canSpy(key),
  Button: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  ChevronLeftIcon: () => null,
  ChevronRightIcon: () => null,
  DeleteIcon: () => null,
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  Modal: ({ children, open }: any) => (open ? <div>{children}</div> : null),
  ModalBody: ({ children }: any) => <div>{children}</div>,
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  NoSsr: ({ children }: any) => <>{children}</>,
  Pagination: ({ onChange }: any) => (
    <button data-testid="go-to-page-3" onClick={(e) => onChange(e, 3)}>
      page 3
    </button>
  ),
  PaginationItem: () => null,
  PrimaryActionButtons: ({ secondaryButtonProps }: any) => (
    <button data-testid="cancel-connection-modal" onClick={secondaryButtonProps?.onClick}>
      Cancel
    </button>
  ),
  PROMPT_VARIANTS: { DANGER: 'danger' },
  SearchBar: ({ onSearch }: any) => (
    <button data-testid="search-prod" onClick={() => onSearch('prod')}>
      search
    </button>
  ),
  TransferList: ({
    assignableData,
    originalAssignedData,
    assignablePage,
    originalLeftCount,
    originalRightCount,
  }: any) => (
    <div
      data-testid="mock-transfer-list"
      data-assignable-ids={(assignableData || []).map((c: any) => c.id).join(',')}
      data-assigned-ids={(originalAssignedData || []).map((c: any) => c.id).join(',')}
      data-original-left-count={originalLeftCount}
      data-original-right-count={originalRightCount}
      data-move-all-left-disabled={
        originalLeftCount !== undefined && (assignableData?.length ?? 0) < originalLeftCount
      }
      data-move-all-right-disabled={
        originalRightCount !== undefined && (originalAssignedData?.length ?? 0) < originalRightCount
      }
    >
      <button data-testid="paginate-assignable-btn" onClick={assignablePage}>
        Next Page
      </button>
    </div>
  ),
  Typography: ({ children }: any) => <span>{children}</span>,
  createAndEditEnvironmentSchema: {},
  createAndEditEnvironmentUiSchema: {},
  useTheme: () => ({
    palette: {
      icon: { default: '#000', secondary: '#111' },
      background: { constant: { table: '#fff' } },
      text: { default: '#000' },
    },
  }),
}));

import Environments from './index';

const EVENT_ERROR = 'error';
const EVENT_SUCCESS = 'success';

// Shape the real chain produces for a provider-rejected create: `data` is the
// verbatim server envelope (camelCase, per server/models/httputil/httputil.go)
// and `meshkit` is what the @meshery/schemas baseQuery wrapper attaches -
// since v1.3.37 (meshery/schemas#1081) it carries the full envelope, reading
// the server's camelCase detail arrays with a snake_case fallback, so the
// probable cause and remediation list arrive populated.
// `utils/helpers/__tests__/meshkitErrorChain.test.ts` pins that transform
// against the real client; this fixture mirrors its output.
const REJECTED_CREATE = {
  status: 403,
  data: {
    error: 'Unable to create the environment',
    code: 'meshery-server-1448',
    severity: 'ALERT',
    probableCause: ['Your account does not have permission.'],
    suggestedRemediation: ['Ask an organization owner to grant the Environment role.'],
  },
  meshkit: {
    message: 'Unable to create the environment',
    code: 'meshery-server-1448',
    severity: 'ALERT',
    probableCause: ['Your account does not have permission.'],
    suggestedRemediation: ['Ask an organization owner to grant the Environment role.'],
  },
};

const openCreateModalAndSubmit = async () => {
  const user = userEvent.setup();
  render(<Environments />);
  await user.click(screen.getByText('Create'));
  await user.click(await screen.findByTestId('submit-environment'));
};

const notifiedEventTypes = () =>
  notify.mock.calls.map(([arg]) => arg?.event_type?.type ?? arg?.event_type);

describe('Environments create flow notifications', () => {
  beforeEach(() => {
    permissionOverrides.clear();
    notify.mockReset();
    createEnvironment.mockReset();
    updateEnvironment.mockReset();
    deleteEnvironment.mockReset();
  });

  it('does not report success when the provider rejects the create', async () => {
    createEnvironment.mockReturnValue({ unwrap: () => Promise.reject(REJECTED_CREATE) });

    await openCreateModalAndSubmit();

    await waitFor(() => expect(notify).toHaveBeenCalled());

    const types = notifiedEventTypes();
    expect(types).toContain(EVENT_ERROR);
    expect(types).not.toContain(EVENT_SUCCESS);
  });

  it('carries the MeshKit code and remediation into the failure notification', async () => {
    createEnvironment.mockReturnValue({ unwrap: () => Promise.reject(REJECTED_CREATE) });

    await openCreateModalAndSubmit();

    await waitFor(() => expect(notify).toHaveBeenCalled());

    const [payload] = notify.mock.calls[0];
    expect(payload.message).toContain('Unable to create the environment');
    expect(payload.message).toContain('meshery-server-1448');
    expect(payload.message).toContain('Ask an organization owner to grant the Environment role.');
  });

  it('reports success only once the create actually resolves', async () => {
    createEnvironment.mockReturnValue({ unwrap: () => Promise.resolve({ id: 'env-1' }) });

    await openCreateModalAndSubmit();

    await waitFor(() => expect(notify).toHaveBeenCalled());

    const types = notifiedEventTypes();
    expect(types).toContain(EVENT_SUCCESS);
    expect(types).not.toContain(EVENT_ERROR);
    expect(notify.mock.calls[0][0].message).toBe('Environment "prod" created');
  });

  // Pins the #20854 root cause: the submit handler once read `organization`
  // off the form (always undefined) and sent an empty org, so the request
  // either failed or created nothing. The payload must carry a populated
  // `organizationId` in the canonical camelCase wire spelling.
  it('sends a populated organizationId in the create payload', async () => {
    createEnvironment.mockReturnValue({ unwrap: () => Promise.resolve({ id: 'env-1' }) });

    await openCreateModalAndSubmit();

    await waitFor(() => expect(createEnvironment).toHaveBeenCalled());

    expect(createEnvironment).toHaveBeenCalledWith({
      environmentPayload: expect.objectContaining({ organizationId: 'org-1' }),
    });
  });
});

describe('Environments search paging', () => {
  const lastQueryArgs = () => getEnvironmentsQuery.mock.calls.at(-1)?.[0];

  it('goes back to the first page when the search changes', async () => {
    const originalData = ENVIRONMENTS_QUERY_RESULT.data;
    // Pagination only renders when the current page has environments.
    ENVIRONMENTS_QUERY_RESULT.data = {
      environments: [{ id: 'env-1', name: 'staging' }],
      totalCount: 40,
    } as any;
    try {
      const user = userEvent.setup();
      render(<Environments />);

      await user.click(screen.getByTestId('go-to-page-3'));
      // The pagination handler is debounced.
      await waitFor(() => expect(lastQueryArgs()).toMatchObject({ page: 2, search: '' }));

      await user.click(screen.getByTestId('search-prod'));
      await waitFor(() => expect(lastQueryArgs()).toMatchObject({ page: 0, search: 'prod' }));
    } finally {
      ENVIRONMENTS_QUERY_RESULT.data = originalData;
    }
  });
});

describe('Environments permissions and wire query skipping', () => {
  beforeEach(() => {
    permissionOverrides.clear();
    canSpy.mockReset();
    canSpy.mockImplementation(
      (key?: { id?: string }) => permissionOverrides.get(key?.id ?? '') ?? true,
    );
    getEnvironmentsQuery.mockClear();
    getEnvironmentConnectionsQuery.mockClear();
  });

  it('skips useGetEnvironmentsQuery and renders 404 fallback when WorkspaceManagementViewEnvironment is denied', () => {
    permissionOverrides.set('WorkspaceManagementViewEnvironment', false);
    render(<Environments />);

    expect(screen.getByText('no access')).toBeInTheDocument();
    expect(getEnvironmentsQuery).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: true }),
    );
    expect(getEnvironmentConnectionsQuery).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: true }),
    );
  });

  it('runs useGetEnvironmentsQuery without skip when view permission and orgId are present', () => {
    render(<Environments />);

    expect(getEnvironmentsQuery).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: false }),
    );
  });
});

describe('Environments connection assignment modal', () => {
  const originalEnvData = ENVIRONMENTS_QUERY_RESULT.data;
  const originalConnData = CONNECTIONS_QUERY_RESULT.data;

  beforeEach(() => {
    ENVIRONMENTS_QUERY_RESULT.data = {
      environments: [{ id: 'env-1', name: 'production' }],
      totalCount: 1,
    } as any;
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [
        { id: 'conn-1', name: 'Cluster A' },
        { id: 'conn-2', name: 'Cluster B' },
      ],
      totalCount: 2,
    } as any;
  });

  afterEach(() => {
    ENVIRONMENTS_QUERY_RESULT.data = originalEnvData;
    CONNECTIONS_QUERY_RESULT.data = originalConnData;
    CONNECTIONS_QUERY_RESULT.currentData = undefined;
  });

  it('cleans up connection state on close and does not duplicate connections when reopening', async () => {
    const user = userEvent.setup();
    render(<Environments />);

    // Open connection modal
    await user.click(screen.getByTestId('assign-btn-env-1'));

    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    const transferList = screen.getByTestId('mock-transfer-list');
    expect(transferList.getAttribute('data-assignable-ids')).toBe('conn-1,conn-2');

    // Close connection modal
    await user.click(screen.getByTestId('cancel-connection-modal'));
    await waitFor(() => {
      expect(screen.queryByTestId('mock-transfer-list')).not.toBeInTheDocument();
    });

    // Reopen connection modal for the same environment
    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    // Should still have exactly conn-1,conn-2 without duplicates
    const reopenedTransferList = screen.getByTestId('mock-transfer-list');
    expect(reopenedTransferList.getAttribute('data-assignable-ids')).toBe('conn-1,conn-2');
  });

  it('deduplicates incoming connections with identical ids', async () => {
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [
        { id: 'conn-1', name: 'Cluster A' },
        { id: 'conn-1', name: 'Cluster A Duplicate' },
        { id: 'conn-2', name: 'Cluster B' },
      ],
      totalCount: 3,
    } as any;

    const user = userEvent.setup();
    render(<Environments />);

    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    const transferList = screen.getByTestId('mock-transfer-list');
    expect(transferList.getAttribute('data-assignable-ids')).toBe('conn-1,conn-2');
  });

  it('resets pagination offsets back to 0 when modal is closed and reopened', async () => {
    // 50 total connections: 2 pages of 25 items each
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [{ id: 'conn-1', name: 'Cluster A' }],
      totalCount: 50,
    } as any;

    const user = userEvent.setup();
    render(<Environments />);

    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    // Advance to next page (page 1)
    await user.click(screen.getByTestId('paginate-assignable-btn'));
    await waitFor(() => {
      expect(getEnvironmentConnectionsQuery).toHaveBeenCalledWith(
        expect.objectContaining({ page: 1 }),
        expect.anything(),
      );
    });

    // Close modal
    await user.click(screen.getByTestId('cancel-connection-modal'));
    await waitFor(() => {
      expect(screen.queryByTestId('mock-transfer-list')).not.toBeInTheDocument();
    });

    // Reopen modal: should query starting from page 0 again
    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    const assignableCalls = getEnvironmentConnectionsQuery.mock.calls.filter(
      ([args, options]: [any, any]) => args?.filter === '{"assigned":false}' && !options?.skip,
    );
    expect(assignableCalls.length).toBeGreaterThan(0);
    const lastAssignableCall = assignableCalls.at(-1);
    expect(lastAssignableCall?.[0]).toMatchObject({ page: 0 });
  });

  it('enables move-all when deduplicated list is complete even if server totalCount includes duplicates', async () => {
    // 1 page of connections on server, but totalCount is 3 due to duplicate IDs:
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [
        { id: 'conn-1', name: 'Cluster A' },
        { id: 'conn-1', name: 'Cluster A Duplicate' },
        { id: 'conn-2', name: 'Cluster B' },
      ],
      totalCount: 3,
    } as any;

    const user = userEvent.setup();
    render(<Environments />);

    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    const transferList = screen.getByTestId('mock-transfer-list');
    // Deduplicated list has 2 items:
    expect(transferList.getAttribute('data-assignable-ids')).toBe('conn-1,conn-2');
    // Once all pages are loaded (page 0 of 1), the count given to TransferList uses the deduplicated length:
    expect(transferList.getAttribute('data-original-left-count')).toBe('2');
    // Move-all is enabled because assignableData.length (2) is not less than originalLeftCount (2):
    expect(transferList.getAttribute('data-move-all-left-disabled')).toBe('false');
  });

  it('keeps move-all disabled while final-page request is pending or if it fails', async () => {
    CONNECTIONS_QUERY_RESULT.data = {
      connections: Array.from({ length: 25 }, (_, i) => ({
        id: `conn-${i}`,
        name: `Cluster ${i}`,
      })),
      totalCount: 50,
    } as any;

    const user = userEvent.setup();
    render(<Environments />);

    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    const transferList = screen.getByTestId('mock-transfer-list');
    expect(transferList.getAttribute('data-original-left-count')).toBe('50');
    expect(transferList.getAttribute('data-move-all-left-disabled')).toBe('true');

    // Simulate advance to page 1 while the response is pending (currentData is null/pending)
    CONNECTIONS_QUERY_RESULT.currentData = null;

    await user.click(screen.getByTestId('paginate-assignable-btn'));

    // Move-all must remain disabled because page 1 has not successfully accumulated
    expect(transferList.getAttribute('data-original-left-count')).toBe('50');
    expect(transferList.getAttribute('data-move-all-left-disabled')).toBe('true');
  });

  it('does not advance to subsequent page until current page response has accumulated', async () => {
    // 75 total connections: 3 pages of 25 items each
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [{ id: 'conn-0', name: 'Cluster 0' }],
      totalCount: 75,
    } as any;

    const user = userEvent.setup();
    render(<Environments />);

    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    // Advance to page 1
    await user.click(screen.getByTestId('paginate-assignable-btn'));
    await waitFor(() => {
      expect(getEnvironmentConnectionsQuery).toHaveBeenCalledWith(
        expect.objectContaining({ page: 1 }),
        expect.anything(),
      );
    });

    // While page 1 has not accumulated yet, click paginate again
    await user.click(screen.getByTestId('paginate-assignable-btn'));

    // Should NOT have requested page 2 because page 1 has not accumulated
    expect(getEnvironmentConnectionsQuery).not.toHaveBeenCalledWith(
      expect.objectContaining({ page: 2 }),
      expect.anything(),
    );
  });

  it('marks empty page response as loaded so pagination and completion do not stall', async () => {
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [{ id: 'conn-0', name: 'Cluster 0' }],
      totalCount: 50,
    } as any;

    const user = userEvent.setup();
    const { rerender } = render(<Environments />);

    await user.click(screen.getByTestId('assign-btn-env-1'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-transfer-list')).toBeInTheDocument();
    });

    // Advance to page 1
    await user.click(screen.getByTestId('paginate-assignable-btn'));

    // Simulate page 1 returning an empty list
    CONNECTIONS_QUERY_RESULT.data = {
      connections: [],
      totalCount: 50,
    } as any;
    rerender(<Environments />);

    await waitFor(() => {
      const transferList = screen.getByTestId('mock-transfer-list');
      // Because page 1 was accumulated (even though empty), list is marked fully loaded
      expect(transferList.getAttribute('data-original-left-count')).toBe('1');
    });
  });
});
