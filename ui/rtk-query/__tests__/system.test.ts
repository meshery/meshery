import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

beforeAll(() => {
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
});

const okResponse = (body: unknown = {}) => ({
  ok: true,
  status: 200,
  redirected: false,
  headers: new Headers({ 'content-type': 'application/json' }),
  url: '',
  text: () => Promise.resolve(JSON.stringify(body)),
  json: () => Promise.resolve(body),
  clone() {
    return this;
  },
});

const setupStore = async () => {
  vi.resetModules();
  const apiMod = await import('../index');
  await import('../system');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('system endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ contexts: [], total_count: 0 }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports the expected hooks', async () => {
    const mod = await import('../system');
    expect(mod.useGetDatabaseSummaryQuery).toBeTypeOf('function');
    expect(mod.useGetAdaptersQuery).toBeTypeOf('function');
    expect(mod.useGetAvailableAdaptersQuery).toBeTypeOf('function');
    expect(mod.useLazyPingAdapterQuery).toBeTypeOf('function');
    expect(mod.useManageAdapterMutation).toBeTypeOf('function');
    expect(mod.useGetSystemSyncQuery).toBeTypeOf('function');
    expect(mod.useLazyGetSystemSyncQuery).toBeTypeOf('function');
    expect(mod.useGetKubernetesContextsQuery).toBeTypeOf('function');
    expect(mod.useLazyGetKubernetesContextsQuery).toBeTypeOf('function');
    expect(mod.useAdapterOperationMutation).toBeTypeOf('function');
    expect(mod.useLazyGetSmiResultsQuery).toBeTypeOf('function');
  });

  it('getDatabaseSummary issues GET with pagination params', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getDatabaseSummary.initiate({
        page: 1,
        pagesize: 20,
        search: 'foo',
        order: 'name',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/system/database');
    expect(req.url).toContain('page=1');
    expect(req.url).toContain('pagesize=20');
    expect(req.url).toContain('search=foo');
    expect(req.url).toContain('order=name');
  });

  it('getAdapters issues GET /api/system/adapters with credentials include', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getAdapters.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/system/adapters');
    expect(req.credentials).toBe('include');
  });

  it('getAvailableAdapters issues GET /api/system/availableAdapters', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getAvailableAdapters.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/availableAdapters');
    expect(req.method).toBe('GET');
    expect(req.credentials).toBe('include');
  });

  it('pingAdapter passes the adapter argument as a query param', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.pingAdapter.initiate('istio.local:10000'));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/adapters');
    expect(req.url).toContain('adapter=istio.local%3A10000');
    expect(req.credentials).toBe('include');
  });

  it('getSystemSync issues GET /api/system/sync', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getSystemSync.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/sync');
    expect(req.method).toBe('GET');
  });

  it('getKubernetesContexts normalizes the response via transformResponse', async () => {
    fetchMock.mockResolvedValue(
      okResponse({
        total_count: 1,
        contexts: [
          {
            id: 'ctx-1',
            connection_id: 'conn-1',
            is_current_context: true,
            created_by: 'me',
            deployment_type: 'in_cluster',
          },
        ],
      }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getKubernetesContexts.initiate({ pagesize: 50, search: 'prod' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/kubernetes/contexts');
    expect(req.url).toContain('pagesize=50');
    expect(req.url).toContain('search=prod');
    // The transformResponse normalizes snake_case to camelCase
    expect(res.data).toEqual(
      expect.objectContaining({
        totalCount: 1,
        contexts: [
          expect.objectContaining({
            connectionId: 'conn-1',
            isCurrentContext: true,
            createdBy: 'me',
            deploymentType: 'in_cluster',
          }),
        ],
      }),
    );
  });

  it('getKubernetesContexts uses default page size and empty search when args undefined', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getKubernetesContexts.initiate(undefined));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('pagesize=10');
    // Empty search string should still be encoded in the params
    expect(req.url).toMatch(/search=&?/);
  });

  it('adapterOperation issues POST with form-encoded body and custom url override', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.adapterOperation.initiate({
        url: 'system/adapter/operation',
        body: 'foo=bar&baz=qux',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/system/adapter/operation');
    expect(req.headers.get('content-type')).toContain('application/x-www-form-urlencoded');
    expect(await req.text()).toBe('foo=bar&baz=qux');
  });

  it('getSmiResults issues GET on /api/smi/results with optional params', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getSmiResults.initiate({
        page: 0,
        pagesize: 5,
        search: 'x',
        order: 'created_at desc',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/smi/results');
    expect(req.url).toContain('page=0');
    expect(req.url).toContain('pagesize=5');
    expect(req.credentials).toBe('include');
  });

  it('manageAdapter DELETE branch issues DELETE with adapter param', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.manageAdapter.initiate({
        method: 'DELETE',
        adapter: 'istio.local:10000',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/system/adapter/manage');
    expect(req.url).toContain('adapter=istio.local%3A10000');
  });

  it('manageAdapter default branch issues POST with form-encoded meshLocationURL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.manageAdapter.initiate({
        meshLocationURL: 'tcp://meshery-istio:10000',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/system/adapter/manage');
    expect(req.headers.get('content-type')).toContain('application/x-www-form-urlencoded');
    const body = await req.text();
    expect(body).toContain('meshLocationURL=');
    expect(body).toContain(encodeURIComponent('tcp://meshery-istio:10000'));
  });
});
