import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

const schemaMocks = vi.hoisted(() => {
  const result = { isLoading: false };
  return {
    result,
    upsertProfileTrigger: vi.fn(),
    deleteProfileTrigger: vi.fn(),
    useGetPerformanceProfilesQuery: vi.fn((queryArg, options) => ({
      hook: 'getPerformanceProfiles',
      queryArg,
      options,
    })),
    useUpsertPerformanceProfileMutation: vi.fn(() => [vi.fn(), result] as const),
    useGetPerformanceProfileQuery: vi.fn((queryArg, options) => ({
      hook: 'getPerformanceProfile',
      queryArg,
      options,
    })),
    useDeletePerformanceProfileMutation: vi.fn(() => [vi.fn(), result] as const),
  };
});

vi.mock('@meshery/schemas/mesheryApi', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@meshery/schemas/mesheryApi')>()),
  useGetPerformanceProfilesQuery: schemaMocks.useGetPerformanceProfilesQuery,
  useUpsertPerformanceProfileMutation: schemaMocks.useUpsertPerformanceProfileMutation,
  useGetPerformanceProfileQuery: schemaMocks.useGetPerformanceProfileQuery,
  useDeletePerformanceProfileMutation: schemaMocks.useDeletePerformanceProfileMutation,
}));

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

const loadModule = async () => {
  const mod = await import('../performance-profile');
  return mod;
};

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
    schemaMocks.useGetPerformanceProfilesQuery.mockClear();
    schemaMocks.useUpsertPerformanceProfileMutation.mockClear();
    schemaMocks.useGetPerformanceProfileQuery.mockClear();
    schemaMocks.useDeletePerformanceProfileMutation.mockClear();
    schemaMocks.upsertProfileTrigger = vi.fn();
    schemaMocks.deleteProfileTrigger = vi.fn();
    schemaMocks.useUpsertPerformanceProfileMutation.mockReturnValue([
      schemaMocks.upsertProfileTrigger,
      schemaMocks.result,
    ] as const);
    schemaMocks.useDeletePerformanceProfileMutation.mockReturnValue([
      schemaMocks.deleteProfileTrigger,
      schemaMocks.result,
    ] as const);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports all expected hooks', async () => {
    const mod = await loadModule();
    expect(mod.useGetPerformanceProfilesQuery).toBeTypeOf('function');
    expect(mod.useSavePerformanceProfileMutation).toBeTypeOf('function');
    expect(mod.useGetProfileResultsQuery).toBeTypeOf('function');
    expect(mod.useGetPerformanceProfileByIdQuery).toBeTypeOf('function');
    expect(mod.useDeletePerformanceProfileMutation).toBeTypeOf('function');
    expect(mod.useGetProfileResultsByIdQuery).toBeTypeOf('function');
  });

  it('getPerformanceProfiles delegates to the schema-generated hook with normalized params', async () => {
    const { useGetPerformanceProfilesQuery } = await loadModule();
    const result = useGetPerformanceProfilesQuery(
      {
        page: 2,
        pagesize: 25,
        search: 'foo',
        order: 'name desc',
      },
      { skip: true },
    );

    expect(result).toEqual({
      hook: 'getPerformanceProfiles',
      queryArg: {
        page: '2',
        pagesize: '25',
        search: 'foo',
        order: 'name desc',
      },
      options: { skip: true },
    });
    expect(schemaMocks.useGetPerformanceProfilesQuery).toHaveBeenCalledTimes(1);
  });

  it('savePerformanceProfile delegates to schema-generated upsertPerformanceProfile', async () => {
    const { useSavePerformanceProfileMutation } = await loadModule();
    const body = { name: 'profile-a', endpoints: ['/x'] };
    const [trigger, result] = useSavePerformanceProfileMutation();

    trigger({ body });

    expect(result).toBe(schemaMocks.result);
    expect(schemaMocks.useUpsertPerformanceProfileMutation).toHaveBeenCalledTimes(1);
    expect(schemaMocks.upsertProfileTrigger).toHaveBeenCalledWith({ body });
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
    expect(url).toContain('/api/performance/profiles/results');
    expect(url).toContain('from=2024-01-01');
    expect(url).toContain('to=2024-02-01');
  });

  it('getPerformanceProfileById delegates to schema-generated getPerformanceProfile', async () => {
    const { useGetPerformanceProfileByIdQuery } = await loadModule();
    const result = useGetPerformanceProfileByIdQuery({ id: 'abc' }, { skip: true });

    expect(result).toEqual({
      hook: 'getPerformanceProfile',
      queryArg: { performanceProfileId: 'abc' },
      options: { skip: true },
    });
    expect(schemaMocks.useGetPerformanceProfileQuery).toHaveBeenCalledTimes(1);
  });

  it('deletePerformanceProfile delegates to schema-generated deletePerformanceProfile', async () => {
    const { useDeletePerformanceProfileMutation } = await loadModule();
    const [trigger, result] = useDeletePerformanceProfileMutation();

    trigger({ id: 'xyz' });

    expect(result).toBe(schemaMocks.result);
    expect(schemaMocks.useDeletePerformanceProfileMutation).toHaveBeenCalledTimes(1);
    expect(schemaMocks.deleteProfileTrigger).toHaveBeenCalledWith({ performanceProfileId: 'xyz' });
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
    expect(url).toContain('/api/performance/profiles/pp-2/results');
    expect(url).toContain('page=1');
    expect(url).toContain('pagesize=5');
  });
});
