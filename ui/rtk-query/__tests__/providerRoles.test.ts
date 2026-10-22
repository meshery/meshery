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
  await import('../providerRoles');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('providerRoles endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ user: 'me' }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports useGetUserProviderRolesQuery', async () => {
    const mod = await import('../providerRoles');
    expect(mod.useGetUserProviderRolesQuery).toBeTypeOf('function');
  });

  it('getUserProviderRoles issues GET /api/user with credentials', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getUserProviderRoles.initiate({}));
    expect(fetchMock).toHaveBeenCalled();
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/user');
    expect(req.method).toBe('GET');
    // credentials: 'include' produces request.credentials === 'include'
    expect(req.credentials).toBe('include');
  });
});
