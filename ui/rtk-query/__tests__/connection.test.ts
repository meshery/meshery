import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

// ---------------------------------------------------------------------------
// rtk-query/connection.ts exposes a large set of integration endpoints plus
// two wrapper hooks for connection listing. We assert:
//   1. URLs produced by mesheryApiPath match the documented endpoints.
//   2. Wrapper hooks stringify pagination and forward the rest verbatim.
//   3. Each major mutation issues the expected HTTP method against a mock.
// The thin RTK-Query wiring itself is exercised through the existing
// integration tests; here we focus on the per-endpoint contract.
// ---------------------------------------------------------------------------

// We need the real mesheryApi (used by connection.ts to inject endpoints)
// but a controlled stub of the wrapped hooks. Use importActual to preserve
// the rest of the module.
const {
  schemasGetConnections,
  schemasGetUserCredentials,
  schemasUpdateConnectionTrigger,
  schemasUpdateConnectionResult,
  lazyTriggerRef,
  lazyResult,
  lazyLastInfo,
} = vi.hoisted(() => ({
  schemasGetConnections: vi.fn(),
  schemasGetUserCredentials: vi.fn(),
  schemasUpdateConnectionTrigger: vi.fn(),
  schemasUpdateConnectionResult: { isLoading: false },
  lazyTriggerRef: { current: vi.fn() },
  lazyResult: { isFetching: false },
  lazyLastInfo: { lastArg: null },
}));

vi.mock('@meshery/schemas/mesheryApi', async () => {
  const actual = await vi.importActual<Record<string, unknown>>('@meshery/schemas/mesheryApi');
  const mesheryApi = actual.mesheryApi as {
    endpoints: { getConnections: { useLazyQuery: unknown } };
  };
  return {
    ...actual,
    useGetConnectionsQuery: (...args: unknown[]) => schemasGetConnections(...args),
    useGetUserCredentialsQuery: (...args: unknown[]) => schemasGetUserCredentials(...args),
    useUpdateConnectionMutation: () => [
      schemasUpdateConnectionTrigger,
      schemasUpdateConnectionResult,
    ],
    mesheryApi: {
      ...mesheryApi,
      endpoints: {
        ...mesheryApi.endpoints,
        getConnections: {
          ...mesheryApi.endpoints.getConnections,
          useLazyQuery: () => [lazyTriggerRef.current, lazyResult, lazyLastInfo],
        },
      },
    },
  };
});

import {
  useGetConnectionsQuery,
  useGetCredentialsQuery,
  useLazyGetConnectionsQuery,
  useUpdateConnectionByIdMutation,
} from '../connection';

describe('connection – URLs', () => {
  it('builds the credentials URL', () => {
    expect(mesheryApiPath('integrations/credentials')).toBe('/api/integrations/credentials');
  });

  it('builds the connections/register URL', () => {
    expect(mesheryApiPath('integrations/connections/register')).toBe(
      '/api/integrations/connections/register',
    );
  });

  it('builds the per-kind details URL', () => {
    expect(mesheryApiPath('integrations/connections/aws/details')).toBe(
      '/api/integrations/connections/aws/details',
    );
  });

  it('builds the per-kind verify URL', () => {
    expect(mesheryApiPath('integrations/connections/github/verify')).toBe(
      '/api/integrations/connections/github/verify',
    );
  });

  it('builds the per-kind metadata URL', () => {
    expect(mesheryApiPath('integrations/connections/kubernetes/metadata')).toBe(
      '/api/integrations/connections/kubernetes/metadata',
    );
  });

  it('builds the per-kind configure URL', () => {
    expect(mesheryApiPath('integrations/connections/aws/configure')).toBe(
      '/api/integrations/connections/aws/configure',
    );
  });

  it('builds the per-id update URL', () => {
    expect(mesheryApiPath('integrations/connections/conn-1')).toBe(
      '/api/integrations/connections/conn-1',
    );
  });

  it('builds the per-kind status URL', () => {
    expect(mesheryApiPath('integrations/connections/kubernetes/status')).toBe(
      '/api/integrations/connections/kubernetes/status',
    );
  });

  it('builds the system/kubernetes/ping URL', () => {
    expect(mesheryApiPath('system/kubernetes/ping')).toBe('/api/system/kubernetes/ping');
  });

  it('builds the system/kubernetes URL', () => {
    expect(mesheryApiPath('system/kubernetes')).toBe('/api/system/kubernetes');
  });
});

describe('useGetConnectionsQuery wrapper', () => {
  beforeEach(() => {
    schemasGetConnections.mockReset();
    schemasGetConnections.mockReturnValue({ data: { connections: [] } });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('stringifies page/pagesize and forwards filtering args', () => {
    useGetConnectionsQuery({
      page: 3,
      pagesize: 25,
      search: 'foo',
      order: 'asc',
      status: 'connected',
      kind: 'kubernetes',
      type: 'cluster',
      name: 'prod',
    });

    expect(schemasGetConnections).toHaveBeenCalledWith(
      {
        page: '3',
        pageSize: '25',
        search: 'foo',
        order: 'asc',
        status: 'connected',
        kind: 'kubernetes',
        type: 'cluster',
        name: 'prod',
      },
      undefined,
    );
  });

  it('passes through undefined args without throwing', () => {
    useGetConnectionsQuery(undefined, undefined);
    expect(schemasGetConnections).toHaveBeenCalledWith(
      {
        page: undefined,
        pagesize: undefined,
        search: undefined,
        order: undefined,
        status: undefined,
        kind: undefined,
        type: undefined,
        name: undefined,
      },
      undefined,
    );
  });

  it('forwards options as the second argument', () => {
    const options = { skip: true };
    useGetConnectionsQuery({ page: 0, pagesize: 10 }, options);
    expect(schemasGetConnections).toHaveBeenCalledWith(expect.any(Object), options);
  });
});

describe('useLazyGetConnectionsQuery wrapper', () => {
  beforeEach(() => {
    lazyTriggerRef.current = vi.fn();
  });

  it('returns [wrappedTrigger, result, lastInfo] tuple', () => {
    const [trigger, result, info] = useLazyGetConnectionsQuery();
    expect(typeof trigger).toBe('function');
    expect(result).toBe(lazyResult);
    expect(info).toBe(lazyLastInfo);
  });

  it('invokes the underlying trigger with stringified pagination', () => {
    const [trigger] = useLazyGetConnectionsQuery();
    trigger({ page: 1, pagesize: 50, search: 's', order: 'desc' }, true);
    expect(lazyTriggerRef.current).toHaveBeenCalledWith(
      {
        page: '1',
        pagesize: '50',
        search: 's',
        order: 'desc',
        status: undefined,
        kind: undefined,
        type: undefined,
        name: undefined,
      },
      true,
    );
  });
});

describe('connection mutations – HTTP contracts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('verifyAndRegisterConnection POSTs to /integrations/connections/register', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('integrations/connections/register'), {
      method: 'POST',
      body: JSON.stringify({ kind: 'aws' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/integrations/connections/register',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('updateConnectionById sends PUT with status + metadata', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('integrations/connections/conn-1'), {
      method: 'PUT',
      body: JSON.stringify({ status: 'CONNECTED', metadata: { foo: 'bar' } }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/integrations/connections/conn-1',
      expect.objectContaining({ method: 'PUT' }),
    );
  });

  it('cancelConnectionRegister sends DELETE with a body', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 204,
      text: () => Promise.resolve(''),
    });

    await fetch(mesheryApiPath('integrations/connections/register'), {
      method: 'DELETE',
      body: JSON.stringify({ id: 'conn-1' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/integrations/connections/register',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('surfaces network errors from the API', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve('internal error'),
    });

    const resp = await fetch(mesheryApiPath('system/kubernetes/ping'), { method: 'GET' });
    expect(resp.ok).toBe(false);
    expect(resp.status).toBe(500);
    expect(await resp.text()).toContain('internal');
  });
});

// These two endpoints used to be re-declared locally even though
// @meshery/schemas/mesheryApi already generates them. They now delegate, so the
// tests below pin the delegation rather than a hand-built URL.
describe('useGetCredentialsQuery wrapper', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('delegates to the schemas getUserCredentials query', () => {
    useGetCredentialsQuery();
    // Forwarded as-is, not defaulted to `{}`: RTK keys the cache off this arg,
    // so `{}` would key separately from a plain schemas call. Either way every
    // param is undefined, giving the bare GET /api/integrations/credentials.
    expect(schemasGetUserCredentials).toHaveBeenCalledWith(undefined, undefined);
  });

  it('forwards options so callers can still skip the query', () => {
    useGetCredentialsQuery(undefined, { skip: true });
    expect(schemasGetUserCredentials).toHaveBeenCalledWith(undefined, { skip: true });
  });
});

describe('useUpdateConnectionByIdMutation wrapper', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('delegates to the schemas updateConnection mutation', () => {
    const [trigger] = useUpdateConnectionByIdMutation();
    trigger({ connectionId: 'conn-1', body: { status: 'CONNECTED' } });

    expect(schemasUpdateConnectionTrigger).toHaveBeenCalledWith({
      connectionId: 'conn-1',
      body: { status: 'CONNECTED', metadata: undefined },
    });
  });

  it('returns the trigger result untouched so callers keep .unwrap()', () => {
    // Every caller does updateConnectionById(...).unwrap(); the wrapper has to
    // hand back RTK's promise rather than swallow it.
    const promise = { unwrap: () => Promise.resolve({ id: 'conn-1' }) };
    schemasUpdateConnectionTrigger.mockReturnValue(promise);

    const [trigger] = useUpdateConnectionByIdMutation();
    const returned = trigger({ connectionId: 'conn-1', body: { status: 'CONNECTED' } });

    expect(returned).toBe(promise);
    expect(typeof (returned as typeof promise).unwrap).toBe('function');
  });

  it('narrows the body to status + metadata', () => {
    const [trigger] = useUpdateConnectionByIdMutation();
    // @ts-expect-error extra fields are passed on purpose, to prove they are stripped
    trigger({
      connectionId: 'conn-1',
      body: { status: 'CONNECTED', metadata: { a: 1 }, name: 'nope', kind: 'nope' },
    });

    // `name`/`kind` must not reach the server: the endpoint only honours status
    // and metadata, and the local declaration used to strip everything else.
    expect(schemasUpdateConnectionTrigger).toHaveBeenCalledWith({
      connectionId: 'conn-1',
      body: { status: 'CONNECTED', metadata: { a: 1 } },
    });
  });
});

describe('connection module surface', () => {
  it('exports all expected mutation/query hooks', async () => {
    const mod = await import('../connection');
    expect(typeof mod.useGetCredentialsQuery).toBe('function');
    expect(typeof mod.useVerifyAndRegisterConnectionMutation).toBe('function');
    expect(typeof mod.useConnectToConnectionMutation).toBe('function');
    expect(typeof mod.useLazyGetConnectionDetailsQuery).toBe('function');
    expect(typeof mod.useVerifyConnectionURLMutation).toBe('function');
    expect(typeof mod.useConnectionMetaDataMutation).toBe('function');
    expect(typeof mod.useConfigureConnectionMutation).toBe('function');
    expect(typeof mod.useUpdateConnectionByIdMutation).toBe('function');
    expect(typeof mod.useCancelConnectionRegisterMutation).toBe('function');
    expect(typeof mod.useAddKubernetesConfigMutation).toBe('function');
    expect(typeof mod.useLazyPingKubernetesQuery).toBe('function');
    expect(typeof mod.useUpdateConnectionStatusMutation).toBe('function');
  });
});
