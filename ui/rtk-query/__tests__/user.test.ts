import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// Mock @/utils/utils to avoid pulling the navigator/UI component tree
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
  initiateQuery: vi.fn(
    async (
      query: { initiate: (variables: unknown, options?: unknown) => unknown },
      variables?: unknown,
      options?: unknown,
    ) => ({
      data: { stubbed: true, variables, options },
    }),
  ),
}));

// ExtensionPointSchemaValidator is loaded via dynamic `require` inside
// getExtensionsByType.transformResponse — its real module is not necessary
// for our tests. Provide a no-op default that returns the array passed in.
// The path here matches what user.ts uses: `require('../utils/ExtensionPointSchemaValidator')`.
vi.mock('../../utils/ExtensionPointSchemaValidator', () => ({
  default: (_type: string) => (items: unknown) => items,
}));
vi.mock('../utils/ExtensionPointSchemaValidator', () => ({
  default: (_type: string) => (items: unknown) => items,
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

const textResponse = (text: string, status = 500) => ({
  ok: status < 400,
  status,
  redirected: false,
  headers: new Headers({ 'content-type': 'text/plain' }),
  url: '',
  text: () => Promise.resolve(text),
  json: () => Promise.reject(new SyntaxError('non-json')),
  clone() {
    return this;
  },
});

const setupStore = async () => {
  vi.resetModules();
  const apiMod = await import('../index');
  const userMod = await import('../user');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store, userMod };
};

describe('user endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports the expected hooks and helpers', async () => {
    const mod = await import('../user');
    expect(mod.useGetExtensionsByTypeQuery).toBeTypeOf('function');
    expect(mod.useLazyGetExtensionsByTypeQuery).toBeTypeOf('function');
    expect(mod.useGetFullPageExtensionsQuery).toBeTypeOf('function');
    expect(mod.useLazyGetFullPageExtensionsQuery).toBeTypeOf('function');
    expect(mod.useGetLoadTestPrefsQuery).toBeTypeOf('function');
    expect(mod.useUpdateLoadTestPrefsMutation).toBeTypeOf('function');
    expect(mod.useHandleUserInviteMutation).toBeTypeOf('function');
    expect(mod.useLazyGetTokenQuery).toBeTypeOf('function');
    expect(mod.useGetUserPrefQuery).toBeTypeOf('function');
    expect(mod.useUpdateUserPrefMutation).toBeTypeOf('function');
    expect(mod.useGetUserPrefWithContextQuery).toBeTypeOf('function');
    expect(mod.useUpdateUserPrefWithContextMutation).toBeTypeOf('function');
    expect(mod.useGetLoggedInUserQuery).toBeTypeOf('function');
    expect(mod.useGetProviderCapabilitiesQuery).toBeTypeOf('function');
    expect(mod.useHandleFeedbackFormSubmissionMutation).toBeTypeOf('function');
    expect(mod.useGetAllUsersQuery).toBeTypeOf('function');
    expect(mod.useRemoveUserFromTeamMutation).toBeTypeOf('function');
    expect(mod.useGetSystemVersionQuery).toBeTypeOf('function');
    expect(mod.useInstallProviderExtensionMutation).toBeTypeOf('function');
    expect(mod.useRemoveProviderExtensionMutation).toBeTypeOf('function');
    expect(mod.useGetUserProfileSummaryByIdQuery).toBeTypeOf('function');
    expect(mod.useGetUserByIdQuery).toBeTypeOf('function');
    expect(mod.useGetUsersForOrgQuery).toBeTypeOf('function');
    expect(mod.useGetTeamsQuery).toBeTypeOf('function');
    expect(mod.useLazyGetTeamsQuery).toBeTypeOf('function');
    expect(mod.userApi).toBeDefined();
    expect(mod.getProviderCapabilities).toBeTypeOf('function');
    expect(mod.getUserAccessToken).toBeTypeOf('function');
    expect(mod.getUserProfile).toBeTypeOf('function');
    expect(mod.getSystemVersion).toBeTypeOf('function');
    expect(mod.getAllUsers).toBeTypeOf('function');
  });

  it('getLoadTestPrefs GETs /api/user/prefs with context query and coerces legacy wrk2 to fortio', async () => {
    fetchMock.mockResolvedValue(
      okResponse({
        loadTestPrefs: { c: 5, qps: 10, t: '60s', gen: 'wrk2' },
      }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getLoadTestPrefs.initiate(['ctx-a', 'ctx-b']));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/user/prefs?contexts=ctx-a&contexts=ctx-b');
    expect(req.credentials).toBe('include');
    expect(res.data).toEqual({ c: 5, qps: 10, t: '60s', gen: 'fortio' });
  });

  it('getLoadTestPrefs normalizes missing fields to defaults', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getLoadTestPrefs.initiate([]));
    expect(res.data).toEqual({ c: 0, qps: 0, t: '30s', gen: 'fortio' });
  });

  it('updateLoadTestPrefs POSTs JSON body with context and content-type', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.updateLoadTestPrefs.initiate({
        selectedK8sContexts: ['ctx-a'],
        loadTestPrefs: { c: 3, qps: 5, t: '20s', gen: 'fortio' },
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/user/prefs?contexts=ctx-a');
    expect(req.headers.get('content-type')).toContain('application/json');
    const body = JSON.parse(await req.text());
    expect(body).toEqual({ loadTestPrefs: { c: 3, qps: 5, t: '20s', gen: 'fortio' } });
  });

  it('getToken GETs /api/token', async () => {
    fetchMock.mockResolvedValue(okResponse({ token: 'abc' }));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getToken.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/token');
    expect(req.credentials).toBe('include');
  });

  it('getUserPref GETs /api/user/prefs (string url shorthand)', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getUserPref.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/user/prefs');
    expect(req.method).toBe('GET');
  });

  it('updateUserPref POSTs the arg body directly', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.updateUserPref.initiate({ selectedOrganizationID: 'org-9' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/user/prefs');
    expect(JSON.parse(await req.text())).toEqual({ selectedOrganizationID: 'org-9' });
  });

  it('getUserPrefWithContext uses same-origin credentials with contexts in URL', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getUserPrefWithContext.initiate(['ctx-x']));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/user/prefs?contexts=ctx-x');
    expect(req.credentials).toBe('same-origin');
  });

  it('updateUserPrefWithContext POSTs body with JSON content-type', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.updateUserPrefWithContext.initiate({
        selectedK8sContexts: ['ctx-x'],
        body: { x: 1 },
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/user/prefs?contexts=ctx-x');
    expect(req.headers.get('content-type')).toContain('application/json');
    expect(JSON.parse(await req.text())).toEqual({ x: 1 });
  });

  it('getLoggedInUser GETs /api/user', async () => {
    fetchMock.mockResolvedValue(okResponse({ id: 'u1' }));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getLoggedInUser.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/user');
  });

  it('getProviderCapabilities GETs /api/provider/capabilities and normalizes the response', async () => {
    fetchMock.mockResolvedValue(
      okResponse({
        provider_name: 'Meshery',
        provider_type: 'local',
        provider_url: 'http://x',
        provider_description: ['Local provider'],
        capabilities: [{ feature: 'a' }],
      }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getProviderCapabilities.initiate({}));
    expect(res.data).toEqual(
      expect.objectContaining({
        providerName: 'Meshery',
        providerType: 'local',
        providerUrl: 'http://x',
        providerDescription: ['Local provider'],
        capabilities: [{ feature: 'a' }],
      }),
    );
  });

  it('getUserProfileSummaryById GETs /api/user/profile/:id and normalizes the response', async () => {
    fetchMock.mockResolvedValue(
      okResponse({
        id: 'u-1',
        email: 'a@b',
        first_name: 'A',
        last_name: 'B',
        avatar_url: 'http://avatar',
      }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getUserProfileSummaryById.initiate({ id: 'u-1' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/user/profile/u-1');
    expect(res.data).toEqual({
      id: 'u-1',
      email: 'a@b',
      userId: 'u-1',
      avatarUrl: 'http://avatar',
      firstName: 'A',
      lastName: 'B',
    });
  });

  it('getUserProfileSummaryById falls back to raw text when body is not JSON', async () => {
    fetchMock.mockResolvedValue(textResponse('Internal Server Error'));
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getUserProfileSummaryById.initiate({ id: 'u-2' }),
    );
    // Non-2xx status produces an error.data being the raw text body
    expect(res.error).toBeDefined();
    expect((res.error as { data: unknown }).data).toBe('Internal Server Error');
  });

  it('getUserProfileSummaryById returns undefined when body is empty', async () => {
    fetchMock.mockResolvedValue(textResponse('', 200));
    const { api, store } = await setupStore();
    const res = await store.dispatch(
      api.endpoints.getUserProfileSummaryById.initiate({ id: 'u-3' }),
    );
    // normalizeUserProfileSummary(undefined) returns undefined
    expect(res.data).toBeUndefined();
  });

  it('getExtensionsByType GETs /api/provider/capabilities and returns an array', async () => {
    // ExtensionPointSchemaValidator is loaded via dynamic CommonJS require
    // inside transformResponse, which under Vitest may not resolve in the
    // ESM-only test environment. The transform's catch arm returns [] in
    // that case — both branches must yield an array, so we assert just that.
    const extensions = { userPrefs: [{ component: 'CompA' }] };
    fetchMock.mockResolvedValue(okResponse({ extensions }));
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getExtensionsByType.initiate('userPrefs'));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/provider/capabilities');
    expect(req.credentials).toBe('include');
    expect(Array.isArray(res.data)).toBe(true);
  });

  it('getExtensionsByType returns [] when the response has no matching extension type', async () => {
    fetchMock.mockResolvedValue(okResponse({ extensions: {} }));
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getExtensionsByType.initiate('foo'));
    expect(res.data).toEqual([]);
  });

  it('getFullPageExtensions returns name/uri tuples for full_page entries', async () => {
    fetchMock.mockResolvedValue(
      okResponse({
        extensions: {
          dashboard: [
            { type: 'full_page', href: { uri: '/dashboard' } },
            { type: 'iframe', href: { uri: '/x' } },
          ],
          another: [{ type: 'full_page', href: { uri: '/another' } }],
        },
      }),
    );
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getFullPageExtensions.initiate({}));
    expect(res.data).toEqual(
      expect.arrayContaining([
        { name: 'dashboard', uri: '/dashboard' },
        { name: 'another', uri: '/another' },
      ]),
    );
  });

  it('getFullPageExtensions returns [] when extensions is missing', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getFullPageExtensions.initiate({}));
    expect(res.data).toEqual([]);
  });

  it('getSystemVersion GETs /api/system/version', async () => {
    fetchMock.mockResolvedValue(okResponse({ build: '0.0.1' }));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getSystemVersion.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/system/version');
  });

  it('installProviderExtension POSTs /api/provider/extension/install with the extension payload', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    const body = {
      extType: 'navigator',
      packageUrl: 'https://example.com/provider-meshery.tar.gz',
      extensionMetadata: { title: 'Kanvas' },
    };
    await store.dispatch(api.endpoints.installProviderExtension.initiate(body));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/provider/extension/install');
    expect(JSON.parse(await req.text())).toEqual(body);
  });

  it('removeProviderExtension POSTs /api/provider/extension/remove with the extension identity', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    const body = {
      extType: 'navigator',
      title: 'Kanvas',
    };
    await store.dispatch(api.endpoints.removeProviderExtension.initiate(body));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/provider/extension/remove');
    expect(JSON.parse(await req.text())).toEqual(body);
  });

  it('handleFeedbackFormSubmission POSTs the extensions feedback URL', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    const userFeedbackRequestBody = { rating: 5, comment: 'great' };
    await store.dispatch(
      api.endpoints.handleFeedbackFormSubmission.initiate({ userFeedbackRequestBody }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/extensions/api/identity/users/notify/feedback');
    expect(JSON.parse(await req.text())).toEqual(userFeedbackRequestBody);
  });

  it('getAllUsers GETs /api/identity/users with params', async () => {
    fetchMock.mockResolvedValue(okResponse({ users: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getAllUsers.initiate({
        page: 1,
        pagesize: 10,
        search: 'x',
        order: 'name',
        filter: 'role=admin',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/identity/users');
    expect(req.url).toContain('page=1');
    expect(req.url).toContain('pagesize=10');
    expect(req.url).toContain('search=x');
    expect(req.url).toContain('order=name');
    expect(req.url).toContain('filter=role%3Dadmin');
  });

  it('removeUserFromTeam DELETEs the extensions team-user URL', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.removeUserFromTeam.initiate({
        orgId: 'org',
        teamId: 'team',
        userId: 'user',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/extensions/api/identity/orgs/org/teams/team/users/user');
  });

  it('handleUserInvite POSTs to the extensions invite URL', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    const userInvite = { emails: ['a@b'] };
    await store.dispatch(api.endpoints.handleUserInvite.initiate({ orgId: 'org-1', userInvite }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/extensions/api/identity/orgs/org-1/users/invite');
    expect(JSON.parse(await req.text())).toEqual(userInvite);
  });

  it('getAccessToken GETs /api/token and returns just the token string', async () => {
    fetchMock.mockResolvedValue(okResponse({ token: 'tok-123' }));
    const { api, store } = await setupStore();
    const res = await store.dispatch(api.endpoints.getAccessToken.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/token');
    expect(res.data).toBe('tok-123');
  });

  it('useGetUserByIdQuery skips when id is empty (no fetch fires)', async () => {
    const { renderHook } = await import('@testing-library/react');
    const React = await import('react');
    const { Provider } = await import('react-redux');
    const { api, store, userMod } = await setupStore();
    fetchMock.mockResolvedValue(okResponse({}));

    const wrapper = ({ children }: { children: React.ReactNode }) =>
      React.createElement(Provider, { store, children });

    renderHook(() => userMod.useGetUserByIdQuery(''), { wrapper });
    // No fetch should occur because skip is forced true for falsy ids.
    expect(fetchMock).not.toHaveBeenCalled();
    // Use the api to satisfy the unused-binding warning.
    expect(api.endpoints.getUserProfileSummaryById).toBeDefined();
  });

  it('useGetUserByIdQuery passes id through and fetches when id is truthy', async () => {
    const { renderHook, waitFor } = await import('@testing-library/react');
    const React = await import('react');
    const { Provider } = await import('react-redux');
    const { store, userMod } = await setupStore();
    fetchMock.mockResolvedValue(okResponse({ id: 'user-1', email: 'a@b' }));

    const wrapper = ({ children }: { children: React.ReactNode }) =>
      React.createElement(Provider, { store, children });

    renderHook(() => userMod.useGetUserByIdQuery('user-1'), { wrapper });
    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/user/profile/user-1');
  });

  it('useGetUserByIdQuery options.skip=true tightens the inferred skip', async () => {
    const { renderHook } = await import('@testing-library/react');
    const React = await import('react');
    const { Provider } = await import('react-redux');
    const { store, userMod } = await setupStore();
    fetchMock.mockResolvedValue(okResponse({}));

    const wrapper = ({ children }: { children: React.ReactNode }) =>
      React.createElement(Provider, { store, children });
    renderHook(() => userMod.useGetUserByIdQuery('id-1', { skip: true }), { wrapper });
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it('getProviderCapabilities helper delegates to initiateQuery against userApi.endpoints.getProviderCapabilities', async () => {
    const utilsMod = await import('../utils');
    const userMod = await import('../user');
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    await userMod.getProviderCapabilities();
    expect(utilsMod.initiateQuery).toHaveBeenCalledWith(
      userMod.userApi.endpoints.getProviderCapabilities,
    );
  });

  it('getUserAccessToken helper delegates to initiateQuery with the getAccessToken endpoint', async () => {
    const utilsMod = await import('../utils');
    const userMod = await import('../user');
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    await userMod.getUserAccessToken();
    expect(utilsMod.initiateQuery).toHaveBeenCalledWith(
      userMod.userApi.endpoints.getAccessToken,
      {},
      {},
    );
  });

  it('getUserProfile helper delegates to initiateQuery with the getLoggedInUser endpoint', async () => {
    const utilsMod = await import('../utils');
    const userMod = await import('../user');
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    await userMod.getUserProfile();
    expect(utilsMod.initiateQuery).toHaveBeenCalledWith(userMod.userApi.endpoints.getLoggedInUser);
  });

  it('getSystemVersion helper delegates to initiateQuery with the getSystemVersion endpoint', async () => {
    const utilsMod = await import('../utils');
    const userMod = await import('../user');
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    await userMod.getSystemVersion();
    expect(utilsMod.initiateQuery).toHaveBeenCalledWith(userMod.userApi.endpoints.getSystemVersion);
  });

  it('getAllUsers helper delegates to initiateQuery with skip when search is empty', async () => {
    const utilsMod = await import('../utils');
    const userMod = await import('../user');
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    await userMod.getAllUsers({ page: 0, pagesize: 10, search: '' });
    expect(utilsMod.initiateQuery).toHaveBeenCalledWith(
      userMod.userApi.endpoints.getAllUsers,
      { page: 0, pagesize: 10, search: '' },
      { skip: true },
    );
    (utilsMod.initiateQuery as ReturnType<typeof vi.fn>).mockClear();
    await userMod.getAllUsers({ page: 0, pagesize: 10, search: 'a' });
    expect(utilsMod.initiateQuery).toHaveBeenCalledWith(
      userMod.userApi.endpoints.getAllUsers,
      { page: 0, pagesize: 10, search: 'a' },
      { skip: false },
    );
  });
});
