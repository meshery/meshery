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
  createEnvironment,
  updateEnvironment,
  deleteEnvironment,
  ENVIRONMENTS_QUERY_RESULT,
  CONNECTIONS_QUERY_RESULT,
  NOOP_MUTATION,
  UI_STATE,
} = vi.hoisted(() => ({
  notify: vi.fn(),
  createEnvironment: vi.fn(),
  updateEnvironment: vi.fn(),
  deleteEnvironment: vi.fn(),
  ENVIRONMENTS_QUERY_RESULT: {
    data: { environments: [], totalCount: 0 },
    isLoading: false,
    isError: false,
    error: undefined,
  },
  CONNECTIONS_QUERY_RESULT: {
    data: { connections: [], totalCount: 0 },
    isError: false,
    error: undefined,
  },
  NOOP_MUTATION: [vi.fn()],
  UI_STATE: { ui: { organization: { id: 'org-1' } } },
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('../../rtk-query/environments', () => ({
  useCreateEnvironmentMutation: () => [createEnvironment],
  useUpdateEnvironmentMutation: () => [updateEnvironment],
  useDeleteEnvironmentMutation: () => [deleteEnvironment],
  useGetEnvironmentsQuery: () => ENVIRONMENTS_QUERY_RESULT,
  useGetEnvironmentConnectionsQuery: () => CONNECTIONS_QUERY_RESULT,
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
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

// The RJSF form is not under test; expose a button that submits the payload the
// real form would produce.
vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ handleSubmit }: any) => (
    <button
      data-testid="submit-environment"
      onClick={() =>
        handleSubmit({ organizationId: 'org-1', name: 'prod', description: 'production' })
      }
    >
      Save
    </button>
  ),
}));

vi.mock('./environment-card', () => ({ default: () => null }));
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
  Pagination: () => null,
  PaginationItem: () => null,
  PrimaryActionButtons: () => null,
  PROMPT_VARIANTS: { DANGER: 'danger' },
  SearchBar: () => null,
  TransferList: () => null,
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
// and `meshkit` is what the @meshery/schemas baseQuery wrapper attaches - which
// today is message/code/severity only, because that wrapper reads the
// snake_case spellings the server does not emit (meshery/schemas#1081).
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
