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
  await import('../schema');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('schema endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ schema: {} }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports useGetSchemaQuery', async () => {
    const mod = await import('../schema');
    expect(mod.useGetSchemaQuery).toBeTypeOf('function');
  });

  it('getSchema issues GET against the meshkit schema path', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getSchema.initiate({ schemaName: 'connection' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/schema/resource/connection');
  });

  it('getSchema URL-encodes the schema name in the path', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getSchema.initiate({ schemaName: 'environment' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/schema/resource/environment');
  });
});
