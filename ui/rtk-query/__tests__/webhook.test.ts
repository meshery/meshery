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
  await import('../webhook');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('webhook endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ ok: true }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports useSupportWebHookMutation', async () => {
    const mod = await import('../webhook');
    expect(mod.useSupportWebHookMutation).toBeTypeOf('function');
  });

  it('supportWebHook issues POST on the extensions webhook URL with the body', async () => {
    const { api, store } = await setupStore();
    const body = { message: 'hello', email: 'me@x.io' };
    await store.dispatch(api.endpoints.supportWebHook.initiate({ type: 'support', body }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/extensions/api/webhook/support');
    expect(JSON.parse(await req.text())).toEqual(body);
  });

  it('supportWebHook URL respects the type segment', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.supportWebHook.initiate({ type: 'feedback', body: {} }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/extensions/api/webhook/feedback');
  });
});
