import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// Tests for the getCatalogFilters RTK-query endpoint introduced as part of
// the GraphQL → REST migration. Lives in its own file so that the env-prefix
// setup happens before `../index` (and therefore the schemas package's
// baseUrl) are loaded.
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
  await import('../filter');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('filter – getCatalogFilters endpoint', () => {
  let fetchMock: FetchMock;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('unwraps the {filters: [...]} envelope and returns just the array', async () => {
    fetchMock.mockResolvedValue(
      okResponse({
        page: 0,
        page_size: 10,
        total_count: 2,
        filters: [{ id: 'f1' }, { id: 'f2' }],
      }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getCatalogFilters.initiate({ search: 'envelope' }),
    );
    expect(res.data).toEqual([{ id: 'f1' }, { id: 'f2' }]);
  });

  it('returns the array as-is when the server already returns an unwrapped array', async () => {
    fetchMock.mockResolvedValue(okResponse([{ id: 'f1' }]));
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getCatalogFilters.initiate({ search: 'raw-array' }),
    );
    expect(res.data).toEqual([{ id: 'f1' }]);
  });

  it('returns empty array when neither envelope nor array is provided', async () => {
    fetchMock.mockResolvedValue(okResponse({ message: 'oops' }));
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getCatalogFilters.initiate({ search: 'unknown-shape' }),
    );
    expect(res.data).toEqual([]);
  });

  it('issues GET on /api/filter/catalog with pagination params', async () => {
    fetchMock.mockResolvedValue(okResponse({ filters: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getCatalogFilters.initiate({
        page: 1,
        pagesize: 25,
        search: 'x',
        order: 'name desc',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/filter/catalog');
    expect(req.url).toContain('page=1');
    expect(req.url).toContain('pagesize=25');
    expect(req.url).toContain('search=x');
    expect(req.url).toContain('order=name+desc');
  });
});
