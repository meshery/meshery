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
  await import('../resource');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('resource endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ actors: [] }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports both hooks', async () => {
    const mod = await import('../resource');
    expect(mod.useGetAccessActorsInfoOfResourceQuery).toBeTypeOf('function');
    expect(mod.useCreateAndRevokeResourceAccessRecordMutation).toBeTypeOf('function');
  });

  it('getAccessActorsInfoOfResource issues GET on the extensions sharing path', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getAccessActorsInfoOfResource.initiate({
        resourceType: 'design',
        resourceId: 'res-1',
        actorType: 'user',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/extensions/api/resource/design/share/res-1/user');
  });

  it('createAndRevokeResourceAccessRecord issues POST with payload', async () => {
    const { api, store } = await setupStore();
    const payload = { actorIds: ['u1'], role: 'reader' };
    await store.dispatch(
      api.endpoints.createAndRevokeResourceAccessRecord.initiate({
        resourceType: 'view',
        resourceId: 'res-2',
        resourceAccessMappingPayload: payload,
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/extensions/api/resource/view/share/res-2');
    const sent = JSON.parse(await req.text());
    expect(sent).toEqual(payload);
  });
});
