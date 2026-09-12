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
  const userKeysMod = await import('../userKeys');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store, userKeysMod };
};

describe('userKeys exports', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({ keys: [] }));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports useGetUserKeysQuery and useLazyGetUserKeysQuery', async () => {
    const mod = await import('../userKeys');
    expect(mod.useGetUserKeysQuery).toBeTypeOf('function');
    expect(mod.useLazyGetUserKeysQuery).toBeTypeOf('function');
  });

  it('underlying getUserKeys (schemas) GETs the org keys URL', async () => {
    // userKeys.ts simply wraps the schemas-defined getUserKeys hook. We exercise
    // the actual endpoint via the store dispatch path to confirm the URL/method.
    const { api, store } = await setupStore();
    const endpoint = api.endpoints.getUserKeys;
    expect(endpoint).toBeDefined();
    await store.dispatch(endpoint.initiate({ orgId: 'org-1' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/identity/orgs/org-1/users/keys');
  });

  it('useGetUserKeysQuery passes orgId through and fetches the keys URL', async () => {
    const { renderHook, waitFor } = await import('@testing-library/react');
    const React = await import('react');
    const { Provider } = await import('react-redux');
    const { store, userKeysMod } = await setupStore();

    const wrapper = ({ children }: { children: React.ReactNode }) =>
      React.createElement(Provider, { store, children });

    renderHook(() => userKeysMod.useGetUserKeysQuery({ orgId: 'org-42' }, { skip: false }), {
      wrapper,
    });

    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/identity/orgs/org-42/users/keys');
  });

  it('useGetUserKeysQuery with skip:true does not fetch', async () => {
    const { renderHook } = await import('@testing-library/react');
    const React = await import('react');
    const { Provider } = await import('react-redux');
    const { store, userKeysMod } = await setupStore();

    const wrapper = ({ children }: { children: React.ReactNode }) =>
      React.createElement(Provider, { store, children });

    renderHook(() => userKeysMod.useGetUserKeysQuery({ orgId: 'org-42' }, { skip: true }), {
      wrapper,
    });
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it('useLazyGetUserKeysQuery wraps the schemas lazy hook trigger with orgId mapping', async () => {
    // We need to spy on the schemas mesheryApi.endpoints.getUserKeys.useLazyQuery,
    // which is what useLazyGetUserKeysQuery internally calls.
    vi.resetModules();
    const schemasModule = await import('@meshery/schemas/mesheryApi');
    const trigger = vi.fn();
    const triggerResult = { data: undefined } as never;
    const lastPromiseInfo = {} as never;
    const useLazyQuerySpy = vi
      .spyOn(schemasModule.mesheryApi.endpoints.getUserKeys, 'useLazyQuery')
      .mockReturnValue([trigger, triggerResult, lastPromiseInfo] as never);

    const userKeysMod = await import('../userKeys');
    const [wrappedTrigger, res, info] = userKeysMod.useLazyGetUserKeysQuery();
    expect(useLazyQuerySpy).toHaveBeenCalled();
    expect(res).toBe(triggerResult);
    expect(info).toBe(lastPromiseInfo);

    wrappedTrigger({ orgId: 'org-7' }, true);
    expect(trigger).toHaveBeenCalledWith({ orgId: 'org-7' }, true);

    wrappedTrigger(undefined, undefined);
    expect(trigger).toHaveBeenLastCalledWith({ orgId: undefined }, undefined);

    useLazyQuerySpy.mockRestore();
  });
});
