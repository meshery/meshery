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
// @meshery/schemas. Because this file now passes `overrideExisting: true`,
// the local endpoint implementations in ui/rtk-query/workspace.ts are the
// effective runtime definitions for overlapping workspace endpoints.
//
// The tests below assert the custom behavior that now runs, including
// getWorkspaces, expandInfo enrichment, and infinite-scroll
// serializeQueryArgs / merge behavior.
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

  it('getWorkspaces GETs /api/workspaces with params through the active workspace.ts implementation', async () => {
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
    expect(req.url).not.toContain('expandInfo');
  });

  it('getWorkspaces with expandInfo enriches workspaces with related resource counts', async () => {
    fetchMock.mockImplementation((input) => {
      const url = (input as Request).url || String(input);

      if (url.includes('/api/workspaces?')) {
        return Promise.resolve(
          okResponse({
            workspaces: [
              { id: 'ws-1', name: 'Workspace 1' },
              { id: 'ws-2', name: 'Workspace 2' },
            ],
            total_count: 2,
          }),
        );
      }

      if (url.includes('/api/workspaces/ws-1/designs')) {
        return Promise.resolve(okResponse({ designs: [], total_count: 11 }));
      }

      if (url.includes('/api/workspaces/ws-1/environments')) {
        return Promise.resolve(okResponse({ environments: [], total_count: 12 }));
      }

      if (url.includes('/api/extensions/api/workspaces/ws-1/views')) {
        return Promise.resolve(okResponse({ views: [], total_count: 13 }));
      }

      if (url.includes('/api/extensions/api/workspaces/ws-1/teams')) {
        return Promise.resolve(okResponse({ teams: [], total_count: 14 }));
      }

      if (url.includes('/api/workspaces/ws-2/designs')) {
        return Promise.resolve(okResponse({ designs: [], total_count: 21 }));
      }

      if (url.includes('/api/workspaces/ws-2/environments')) {
        return Promise.resolve(okResponse({ environments: [], total_count: 22 }));
      }

      if (url.includes('/api/extensions/api/workspaces/ws-2/views')) {
        return Promise.resolve(okResponse({ views: [], total_count: 23 }));
      }

      if (url.includes('/api/extensions/api/workspaces/ws-2/teams')) {
        return Promise.resolve(okResponse({ teams: [], total_count: 24 }));
      }

      return Promise.resolve(okResponse({}));
    });

    const { api, store } = await setupStore();

    const result = await store.dispatch(
      api.endpoints.getWorkspaces.initiate({
        page: 0,
        pagesize: 10,
        orgId: 'org-1',
        expandInfo: true,
      }),
    );

    expect(fetchMock).toHaveBeenCalledTimes(9);
    expect(result.data.workspaces).toEqual([
      {
        id: 'ws-1',
        name: 'Workspace 1',
        designCount: 11,
        environmentCount: 12,
        viewCount: 13,
        teamCount: 14,
      },
      {
        id: 'ws-2',
        name: 'Workspace 2',
        designCount: 21,
        environmentCount: 22,
        viewCount: 23,
        teamCount: 24,
      },
    ]);

    const requestedUrls = fetchMock.mock.calls.map((call) => (call[0] as Request).url);
    expect(requestedUrls.some((url) => url.includes('/api/workspaces/ws-1/designs'))).toBe(true);
    expect(requestedUrls.some((url) => url.includes('/api/workspaces/ws-1/environments'))).toBe(
      true,
    );
    expect(
      requestedUrls.some((url) => url.includes('/api/extensions/api/workspaces/ws-1/views')),
    ).toBe(true);
    expect(
      requestedUrls.some((url) => url.includes('/api/extensions/api/workspaces/ws-1/teams')),
    ).toBe(true);
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

  it('getDesignsOfWorkspace GETs on /api/workspaces/:id/designs and normalizes the response', async () => {
    fetchMock.mockResolvedValue(okResponse({ designs: [], total_count: 0 }));
    const { api, store } = await setupStore();

    const result = await store.dispatch(
      api.endpoints.getDesignsOfWorkspace.initiate({
        workspaceId: 'ws-7',
        page: 0,
        pagesize: 10,
      }),
    );

    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/workspaces/ws-7/designs');
    expect(result.data.designs).toEqual([]);
  });

  it('getDesignsOfWorkspace merges pages when infiniteScroll is enabled', async () => {
    fetchMock
      .mockResolvedValueOnce(okResponse({ designs: [{ id: 'design-1' }], total_count: 2 }))
      .mockResolvedValueOnce(okResponse({ designs: [{ id: 'design-2' }], total_count: 2 }));

    const { api, store } = await setupStore();

    const firstPage = await store.dispatch(
      api.endpoints.getDesignsOfWorkspace.initiate({
        workspaceId: 'ws-1',
        infiniteScroll: true,
        page: 0,
        pagesize: 1,
      }),
    );

    const secondPage = await store.dispatch(
      api.endpoints.getDesignsOfWorkspace.initiate({
        workspaceId: 'ws-1',
        infiniteScroll: true,
        page: 1,
        pagesize: 1,
      }),
    );

    expect(firstPage.data.designs).toEqual([{ id: 'design-1' }]);
    expect(secondPage.data.designs).toEqual([{ id: 'design-1' }, { id: 'design-2' }]);
    expect(fetchMock).toHaveBeenCalledTimes(2);
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

  it('getViewsOfWorkspace GETs on the extensions workspace views URL', async () => {
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
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/extensions/api/workspaces/ws-9/views');
  });

  it('getViewsOfWorkspace merges pages when infiniteScroll is enabled', async () => {
    fetchMock
      .mockResolvedValueOnce(okResponse({ views: [{ id: 'view-1' }], total_count: 2 }))
      .mockResolvedValueOnce(okResponse({ views: [{ id: 'view-2' }], total_count: 2 }));

    const { api, store } = await setupStore();

    const firstPage = await store.dispatch(
      api.endpoints.getViewsOfWorkspace.initiate({
        workspaceId: 'ws-1',
        infiniteScroll: true,
        page: 0,
        pagesize: 1,
      }),
    );

    const secondPage = await store.dispatch(
      api.endpoints.getViewsOfWorkspace.initiate({
        workspaceId: 'ws-1',
        infiniteScroll: true,
        page: 1,
        pagesize: 1,
      }),
    );

    expect(firstPage.data.views).toEqual([{ id: 'view-1' }]);
    expect(secondPage.data.views).toEqual([{ id: 'view-1' }, { id: 'view-2' }]);
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it('assignViewToWorkspace POSTs and unassignViewFromWorkspace DELETEs', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.assignViewToWorkspace.initiate({ workspaceId: 'w', viewId: 'v' }),
    );
    let req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/extensions/api/workspaces/w/views/v');

    await store.dispatch(
      api.endpoints.unassignViewFromWorkspace.initiate({ workspaceId: 'w', viewId: 'v' }),
    );
    req = fetchMock.mock.calls[1][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/extensions/api/workspaces/w/views/v');
  });

  it('getTeamsOfWorkspace GETs the extensions workspace teams URL with params', async () => {
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
    expect(req.url).toContain('/api/extensions/api/workspaces/ws/teams');
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
    expect(req.url).toContain('/api/extensions/api/workspaces/w/teams/t');

    await store.dispatch(
      api.endpoints.unassignTeamFromWorkspace.initiate({ workspaceId: 'w', teamId: 't' }),
    );
    req = fetchMock.mock.calls[1][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/extensions/api/workspaces/w/teams/t');
  });

  it('getEventsOfWorkspace hits the extensions events URL', async () => {
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
        name: 'w-name',
        description: 'desc',
        organizationId: 'org-1',
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

  it('imports the active workspace module successfully', async () => {
    const { workspaceMod } = await setupStore();
    expect(workspaceMod).toBeDefined();
  });
});
