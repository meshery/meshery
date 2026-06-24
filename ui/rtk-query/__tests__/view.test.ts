import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// Mock @/utils/utils so we don't transitively pull React-component imports
// (SVG/JS-as-JSX) into Node-test environment. We only need the two encode
// helpers used by view.ts.
vi.mock('@/utils/utils', () => ({
  urlEncodeArrayParam: (key: string, array: unknown) => {
    if (typeof array === 'string') return `${key}=${array}`;
    if (!Array.isArray(array)) return '';
    return array.map((item) => `${key}=${item}`).join('&');
  },
  urlEncodeParams: (params: Record<string, unknown>) => {
    const usp = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v === undefined || v === null) return;
      if (Array.isArray(v)) {
        v.forEach((vv) => usp.append(k, String(vv)));
      } else {
        usp.append(k, String(v));
      }
    });
    return usp.toString();
  },
}));

// Mock the rtk-query helper `initiateQuery` so we don't transitively pull
// `../store` into the test environment.
vi.mock('../utils', () => ({
  shouldOverrideExisting: false,
  initiateQuery: vi.fn(
    async (query: { initiate: (variables: unknown) => unknown }, variables: unknown) => ({
      data: { stubbed: true, variables },
    }),
  ),
}));

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
  await import('../view');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('view endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ views: [] }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports the expected hooks and helpers', async () => {
    const mod = await import('../view');
    expect(mod.useGetViewQuery).toBeTypeOf('function');
    expect(mod.useUpdateViewVisibilityMutation).toBeTypeOf('function');
    expect(mod.useFetchViewsQuery).toBeTypeOf('function');
    expect(mod.useDeleteViewMutation).toBeTypeOf('function');
    expect(mod.getView).toBeTypeOf('function');
    expect(mod.viewsApi).toBeDefined();
  });

  it('getView GETs the extensions view URL by viewId', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getView.initiate({ viewId: 'v-1' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/extensions/api/content/views/v-1');
  });

  it('deleteView issues DELETE on the view URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.deleteView.initiate({ id: 'v-2' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/extensions/api/content/views/v-2');
  });

  it('updateViewVisibility issues PUT with the body', async () => {
    const { api, store } = await setupStore();
    const body = { visibility: 'public' };
    await store.dispatch(api.endpoints.updateViewVisibility.initiate({ id: 'v-3', body }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('PUT');
    expect(req.url).toContain('/api/extensions/api/content/views/v-3');
    expect(JSON.parse(await req.text())).toEqual(body);
  });

  it('fetchViews GETs the views URL with array+scalar params', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.fetchViews.initiate({
        visibility: ['public', 'private'],
        page: 1,
        pagesize: 25,
        search: 'foo',
        order: 'name',
        trim: true,
        shared: true,
        userId: 'user-1',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/extensions/api/content/views?');
    expect(req.url).toContain('visibility=public');
    expect(req.url).toContain('visibility=private');
    expect(req.url).toContain('page=1');
    expect(req.url).toContain('pagesize=25');
    expect(req.url).toContain('search=foo');
    expect(req.url).toContain('userId=user-1');
    expect(req.url).toContain('trim=true');
    expect(req.url).toContain('shared=true');
  });

  it('fetchViews falls back to empty/false for missing optional params', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.fetchViews.initiate({
        visibility: 'public',
        page: 0,
        pagesize: 10,
        userId: 'u',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('search=');
    expect(req.url).toContain('order=');
    expect(req.url).toContain('trim=false');
    expect(req.url).toContain('shared=false');
  });

  it('getView helper delegates to initiateQuery with the getView endpoint', async () => {
    const utilsMod = await import('../utils');
    const viewMod = await import('../view');
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    const result = await viewMod.getView({ viewId: 'helper-1' });
    expect(utilsMod.initiateQuery).toHaveBeenCalledTimes(1);
    const [endpointArg, variablesArg] = (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(endpointArg).toBe(viewMod.viewsApi.endpoints.getView);
    expect(variablesArg).toEqual({ viewId: 'helper-1' });
    expect(result).toEqual({ data: { stubbed: true, variables: { viewId: 'helper-1' } } });
  });
});
