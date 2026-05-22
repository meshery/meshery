import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// Tests for the performance-profile RTK-query endpoints.
//
// Strategy: spin up a real Redux store containing the api reducer/middleware,
// mock global.fetch, and dispatch each endpoint's initiate(). This verifies
// the URL, HTTP method, params and request body that the endpoint produces.
// ---------------------------------------------------------------------------

beforeAll(() => {
  // Make URLs absolute so fetchBaseQuery's `new Request` doesn't reject.
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
  await import('../performance-profile');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

const getFetchCall = (fetchMock: FetchMock, index = 0) => {
  const req = fetchMock.mock.calls[index][0] as Request;
  return {
    url: req.url,
    method: req.method,
    request: req,
  };
};

describe('performance-profile endpoints', () => {
  let fetchMock: FetchMock;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports all expected hooks', async () => {
    const mod = await import('../performance-profile');
    expect(mod.useGetPerformanceProfilesQuery).toBeTypeOf('function');
    expect(mod.useSavePerformanceProfileMutation).toBeTypeOf('function');
    expect(mod.useGetProfileResultsQuery).toBeTypeOf('function');
    expect(mod.useGetPerformanceProfileByIdQuery).toBeTypeOf('function');
    expect(mod.useDeletePerformanceProfileMutation).toBeTypeOf('function');
    expect(mod.useGetProfileResultsByIdQuery).toBeTypeOf('function');
  });

  it('getPerformanceProfiles issues GET with pagination params', async () => {
    fetchMock.mockResolvedValue(okResponse({ profiles: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getPerformanceProfiles.initiate({
        page: 2,
        pagesize: 25,
        search: 'foo',
        order: 'name desc',
      }),
    );
    const { url, method } = getFetchCall(fetchMock);
    expect(method).toBe('GET');
    expect(url).toContain('/api/user/performance/profiles');
    expect(url).toContain('page=2');
    expect(url).toContain('pagesize=25');
    expect(url).toContain('search=foo');
    expect(url).toContain('order=name+desc');
  });

  it('savePerformanceProfile issues POST with the body', async () => {
    fetchMock.mockResolvedValue(okResponse({ id: 'pp-1' }));
    const { api, store } = await setupStore();
    const body = { name: 'profile-a', endpoints: ['/x'] };
    await store.dispatch(api.endpoints.savePerformanceProfile.initiate({ body }));
    const { url, method, request } = getFetchCall(fetchMock);
    expect(method).toBe('POST');
    expect(url).toContain('/api/user/performance/profiles');
    const sent = await request.text();
    expect(JSON.parse(sent)).toEqual(body);
  });

  it('getProfileResults issues GET with date range params', async () => {
    fetchMock.mockResolvedValue(okResponse({ results: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getProfileResults.initiate({
        page: 0,
        pagesize: 10,
        search: '',
        order: '',
        from: '2024-01-01',
        to: '2024-02-01',
      }),
    );
    const { url, method } = getFetchCall(fetchMock);
    expect(method).toBe('GET');
    expect(url).toContain('/api/user/performance/profiles/results');
    expect(url).toContain('from=2024-01-01');
    expect(url).toContain('to=2024-02-01');
  });

  it('getPerformanceProfileById issues GET with profile id in path', async () => {
    fetchMock.mockResolvedValue(okResponse({ id: 'abc' }));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getPerformanceProfileById.initiate({ id: 'abc' }));
    const { url, method } = getFetchCall(fetchMock);
    expect(method).toBe('GET');
    expect(url).toContain('/api/user/performance/profiles/abc');
  });

  it('deletePerformanceProfile issues DELETE', async () => {
    fetchMock.mockResolvedValue(okResponse({ ok: true }));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.deletePerformanceProfile.initiate({ id: 'xyz' }));
    const { url, method } = getFetchCall(fetchMock);
    expect(method).toBe('DELETE');
    expect(url).toContain('/api/user/performance/profiles/xyz');
  });

  it('getProfileResultsById issues GET for a given profile id with params', async () => {
    fetchMock.mockResolvedValue(okResponse({ results: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getProfileResultsById.initiate({
        id: 'pp-2',
        page: 1,
        pagesize: 5,
        search: 'q',
        order: 'created_at',
      }),
    );
    const { url, method } = getFetchCall(fetchMock);
    expect(method).toBe('GET');
    expect(url).toContain('/api/user/performance/profiles/pp-2/results');
    expect(url).toContain('page=1');
    expect(url).toContain('pagesize=5');
  });
});
