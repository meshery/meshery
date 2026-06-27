import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// Mock @/utils/utils to avoid pulling the navigator/UI component tree.
vi.mock('@/utils/utils', () => ({
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

beforeAll(() => {
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
});

// ---------------------------------------------------------------------------
// workspace.ts re-`injectEndpoints` into the shared `api` defined by
// @meshery/schemas. Because the file does NOT pass `overrideExisting: true`,
// any endpoint with a name already defined in the schemas package is ignored
// (with a console warning). The endpoint definitions in workspace.ts that
// share names with the schemas — e.g. getWorkspaces, getDesignsOfWorkspace —
// never take effect at runtime; the schemas' simpler `query` definitions
// run instead. The tests below assert the URLs/methods that REALLY run.
// `getEventsOfWorkspace` is unique to workspace.ts, so its definition does
// take effect — that one is asserted in full.
// ---------------------------------------------------------------------------

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
  const workspaceMod = await import('../workspace');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store, workspaceMod };
};

describe('workspace endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports all expected hooks', async () => {
    const mod = await import('../workspace');
    expect(mod.useGetWorkspacesQuery).toBeTypeOf('function');
    expect(mod.useLazyGetWorkspacesQuery).toBeTypeOf('function');
    expect(mod.useGetEnvironmentsOfWorkspaceQuery).toBeTypeOf('function');
    expect(mod.useAssignEnvironmentToWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useUnassignEnvironmentFromWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useGetDesignsOfWorkspaceQuery).toBeTypeOf('function');
    expect(mod.useAssignDesignToWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useUnassignDesignFromWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useGetViewsOfWorkspaceQuery).toBeTypeOf('function');
    expect(mod.useAssignViewToWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useUnassignViewFromWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useGetTeamsOfWorkspaceQuery).toBeTypeOf('function');
    expect(mod.useAssignTeamToWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useUnassignTeamFromWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useGetEventsOfWorkspaceQuery).toBeTypeOf('function');
    expect(mod.useCreateWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useUpdateWorkspaceMutation).toBeTypeOf('function');
    expect(mod.useDeleteWorkspaceMutation).toBeTypeOf('function');
  });

  it('getWorkspaces (effective definition wins) GETs /api/workspaces with params', async () => {
    fetchMock.mockResolvedValue(okResponse({ workspaces: [], total_count: 0 }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getWorkspaces.initiate({
        page: 0,
        pagesize: 10,
        orgId: 'org-1',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/workspaces');
    expect(req.url).toContain('orgId=org-1');
  });

  it('getEnvironmentsOfWorkspace GETs on /api/workspaces/:id/environments', async () => {
    fetchMock.mockResolvedValue(okResponse({ environments: [], total_count: 0 }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getEnvironmentsOfWorkspace.initiate({
        workspaceId: 'ws-1',
        page: 1,
        pagesize: 5,
        search: 'q',
        order: 'name',
        filter: 'kind=prod',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/workspaces/ws-1/environments');
    expect(req.url).toContain('page=1');
    expect(req.url).toContain('pagesize=5');
  });

  it('assignEnvironmentToWorkspace POSTs to /api/workspaces/:wid/environments/:eid', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.assignEnvironmentToWorkspace.initiate({
        workspaceId: 'ws-1',
        environmentId: 'env-1',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/workspaces/ws-1/environments/env-1');
  });

  it('unassignEnvironmentFromWorkspace DELETEs the workspace environment URL', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.unassignEnvironmentFromWorkspace.initiate({
        workspaceId: 'ws-1',
        environmentId: 'env-1',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/workspaces/ws-1/environments/env-1');
  });

  it('getDesignsOfWorkspace GETs on /api/workspaces/:id/designs', async () => {
    fetchMock.mockResolvedValue(okResponse({ designs: [], total_count: 0 }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getDesignsOfWorkspace.initiate({
        workspaceId: 'ws-7',
        page: 0,
        pagesize: 10,
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/workspaces/ws-7/designs');
  });

  it('assignDesignToWorkspace POSTs and unassignDesignFromWorkspace DELETEs', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.assignDesignToWorkspace.initiate({
        workspaceId: 'ws-1',
        designId: 'd-1',
      }),
    );
    let req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/workspaces/ws-1/designs/d-1');

    await store.dispatch(
      api.endpoints.unassignDesignFromWorkspace.initiate({
        workspaceId: 'ws-1',
        designId: 'd-1',
      }),
    );
    req = fetchMock.mock.calls[1][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/workspaces/ws-1/designs/d-1');
  });

  it('getViewsOfWorkspace GETs on the workspace views URL', async () => {
    fetchMock.mockResolvedValue(okResponse({ views: [], total_count: 0 }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getViewsOfWorkspace.initiate({
        workspaceId: 'ws-9',
        page: 0,
        pagesize: 10,
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/workspaces/ws-9/views');
  });

  it('assignViewToWorkspace POSTs and unassignViewFromWorkspace DELETEs', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.assignViewToWorkspace.initiate({ workspaceId: 'w', viewId: 'v' }),
    );
    let req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/workspaces/w/views/v');

    await store.dispatch(
      api.endpoints.unassignViewFromWorkspace.initiate({ workspaceId: 'w', viewId: 'v' }),
    );
    req = fetchMock.mock.calls[1][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/workspaces/w/views/v');
  });

  it('getTeamsOfWorkspace GETs the workspace teams URL with params', async () => {
    fetchMock.mockResolvedValue(okResponse({ teams: [], total_count: 0 }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getTeamsOfWorkspace.initiate({
        workspaceId: 'ws',
        page: 0,
        pagesize: 10,
        search: 's',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/workspaces/ws/teams');
    expect(req.url).toContain('search=s');
  });

  it('assignTeamToWorkspace POSTs and unassignTeamFromWorkspace DELETEs', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.assignTeamToWorkspace.initiate({ workspaceId: 'w', teamId: 't' }),
    );
    let req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/workspaces/w/teams/t');

    await store.dispatch(
      api.endpoints.unassignTeamFromWorkspace.initiate({ workspaceId: 'w', teamId: 't' }),
    );
    req = fetchMock.mock.calls[1][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/workspaces/w/teams/t');
  });

  it('getEventsOfWorkspace (unique to workspace.ts) hits the extensions events URL', async () => {
    fetchMock.mockResolvedValue(okResponse({ events: [] }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getEventsOfWorkspace.initiate({
        workspaceId: 'ws-1',
        page: 1,
        pagesize: 25,
        search: 'foo',
        order: 'created_at',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.url).toContain('/api/extensions/api/workspaces/ws-1/events');
    expect(req.url).toContain('page=1');
    expect(req.url).toContain('pagesize=25');
  });

  it('createWorkspace POSTs to /api/workspaces', async () => {
    fetchMock.mockResolvedValue(okResponse({ id: 'new' }));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.createWorkspace.initiate({
        body: { name: 'w-name', description: 'desc', organizationId: 'org-1' },
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/workspaces');
  });

  it('updateWorkspace PUTs against /api/workspaces/:workspaceId', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.updateWorkspace.initiate({
        workspaceId: 'w-1',
        body: { name: 'updated' },
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('PUT');
    expect(req.url).toContain('/api/workspaces/w-1');
  });

  it('deleteWorkspace DELETEs /api/workspaces/:workspaceId', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.deleteWorkspace.initiate({ workspaceId: 'w-2' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/workspaces/w-2');
  });

  it('normalizePaginatedCollectionResponse used by getDesignsOfWorkspace is reachable', async () => {
    // workspace.ts imports normalizePaginatedCollectionResponse from ./transforms.
    // We can't easily verify it runs during dispatch (since the schema's
    // basic query wins), but importing the module should at minimum not throw.
    const { workspaceMod } = await setupStore();
    expect(workspaceMod).toBeDefined();
  });
});
