import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// Tests for the kubernetes RTK-query endpoints introduced as part of the
// GraphQL → REST migration. We mirror the strategy used by
// performance-profile.test.ts: a fresh Redux store per test, mock fetch,
// then dispatch endpoint initiate() and assert on the produced Request.
// ---------------------------------------------------------------------------

beforeAll(() => {
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
});

type FetchMock = ReturnType<typeof vi.fn>;

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
  await import('../kubernetes');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('kubernetes endpoints – module surface', () => {
  it('exports the expected hooks', async () => {
    const mod = await import('../kubernetes');
    expect(mod.useGetKubernetesNamespacesQuery).toBeTypeOf('function');
    expect(mod.useLazyGetKubernetesNamespacesQuery).toBeTypeOf('function');
    expect(mod.useGetMesheryOperatorStatusQuery).toBeTypeOf('function');
    expect(mod.useLazyGetMesheryOperatorStatusQuery).toBeTypeOf('function');
    expect(mod.useGetMeshsyncStatusQuery).toBeTypeOf('function');
    expect(mod.useLazyGetMeshsyncStatusQuery).toBeTypeOf('function');
    expect(mod.useGetNatsStatusQuery).toBeTypeOf('function');
    expect(mod.useLazyGetNatsStatusQuery).toBeTypeOf('function');
  });
});

describe('kubernetes endpoints – HTTP', () => {
  let fetchMock: FetchMock;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('getKubernetesNamespaces issues GET /api/system/kubernetes/namespaces with clusterIds', async () => {
    fetchMock.mockResolvedValue(okResponse({ namespaces: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getKubernetesNamespaces.initiate({ clusterIds: ['c1', 'c2'] }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/system/kubernetes/namespaces');
    expect(req.url).toContain('clusterIds=c1');
  });

  it('getKubernetesNamespaces normalizes a string[] to {namespaces: [{namespace}]}', async () => {
    fetchMock.mockResolvedValue(okResponse(['default', 'kube-system']));
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getKubernetesNamespaces.initiate({}));
    expect(res.data).toEqual({
      namespaces: [{ namespace: 'default' }, { namespace: 'kube-system' }],
    });
  });

  it('getKubernetesNamespaces preserves existing {namespaces} envelope', async () => {
    fetchMock.mockResolvedValue(
      okResponse({ namespaces: [{ namespace: 'default' }, { namespace: 'meshery' }] }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getKubernetesNamespaces.initiate({}));
    expect(res.data).toEqual({
      namespaces: [{ namespace: 'default' }, { namespace: 'meshery' }],
    });
  });

  it('getMesheryOperatorStatus hits the per-connection operator status URL', async () => {
    fetchMock.mockResolvedValue(okResponse({ status: 'ENABLED' }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getMesheryOperatorStatus.initiate({ connectionID: 'abc-123' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/system/kubernetes/connections/abc-123/operator/status');
  });

  it('getMesheryOperatorStatus wraps a bare status response in {operator}', async () => {
    fetchMock.mockResolvedValue(okResponse({ status: 'ENABLED', version: '0.7.1' }));
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getMesheryOperatorStatus.initiate({ connectionID: 'abc' }),
    );
    expect(res.data).toEqual({ operator: { status: 'ENABLED', version: '0.7.1' } });
  });

  it('getMesheryOperatorStatus preserves an already-wrapped {operator} envelope', async () => {
    fetchMock.mockResolvedValue(
      okResponse({ operator: { status: 'DISABLED', controller: 'mesheryctl' } }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getMesheryOperatorStatus.initiate({ connectionID: 'abc' }),
    );
    expect(res.data).toEqual({
      operator: { status: 'DISABLED', controller: 'mesheryctl' },
    });
  });

  it('getMeshsyncStatus hits the per-connection meshsync URL', async () => {
    fetchMock.mockResolvedValue(okResponse({ name: 'MeshSync', status: 'Running', version: '1' }));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getMeshsyncStatus.initiate({ connectionID: 'k8s-1' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/kubernetes/connections/k8s-1/meshsync/status');
  });

  it('getMeshsyncStatus wraps bare controller in {controller}', async () => {
    fetchMock.mockResolvedValue(okResponse({ name: 'MeshSync', status: 'Running', version: '1' }));
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getMeshsyncStatus.initiate({ connectionID: 'k8s-1' }),
    );
    expect(res.data).toEqual({
      controller: { name: 'MeshSync', status: 'Running', version: '1' },
    });
  });

  it('getNatsStatus hits the per-connection nats URL', async () => {
    fetchMock.mockResolvedValue(
      okResponse({ name: 'MesheryBroker', status: 'Connected nats://x', version: '2' }),
    );
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getNatsStatus.initiate({ connectionID: 'k8s-1' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/kubernetes/connections/k8s-1/nats/status');
  });

  it('getNatsStatus wraps bare controller in {controller}', async () => {
    fetchMock.mockResolvedValue(
      okResponse({ name: 'MesheryBroker', status: 'Connected', version: '2' }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getNatsStatus.initiate({ connectionID: 'k8s-1' }),
    );
    expect(res.data).toEqual({
      controller: { name: 'MesheryBroker', status: 'Connected', version: '2' },
    });
  });

  it('connection ID is URL-encoded in the path', async () => {
    fetchMock.mockResolvedValue(okResponse({ status: 'ENABLED' }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getMesheryOperatorStatus.initiate({ connectionID: 'a/b/c' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/kubernetes/connections/a%2Fb%2Fc/operator/status');
  });
});
