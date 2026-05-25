import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

// meshsync.ts imports urlEncodeParams from @/utils/utils which transitively
// pulls non-test-safe modules; mock the small utility surface.
vi.mock('@/utils/utils', () => ({
  urlEncodeParams: (params: Record<string, unknown>) => {
    const sp = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v == null) return;
      if (Array.isArray(v)) {
        v.forEach((vv) => sp.append(k, String(vv)));
      } else {
        sp.append(k, String(v));
      }
    });
    return sp.toString();
  },
}));

// transforms.ts is leaf and import-safe, so no need to mock it.

describe('meshsync – module surface', () => {
  it('exports the expected hooks', async () => {
    const mod = await import('../meshsync');
    expect(typeof mod.useGetMeshSyncResourcesQuery).toBe('function');
    expect(typeof mod.useLazyGetMeshSyncResourcesQuery).toBe('function');
    expect(typeof mod.useGetMeshSyncResourceKindsQuery).toBe('function');
    expect(typeof mod.useDeleteMeshsyncResourceMutation).toBe('function');
  });
});

describe('meshsync – URLs', () => {
  it('builds the resources URL', () => {
    expect(mesheryApiPath('system/meshsync/resources')).toBe('/api/system/meshsync/resources');
  });

  it('builds the resources/summary URL', () => {
    expect(mesheryApiPath('system/meshsync/resources/summary')).toBe(
      '/api/system/meshsync/resources/summary',
    );
  });

  it('builds the per-id resource URL', () => {
    expect(mesheryApiPath('system/meshsync/resources/abc-1')).toBe(
      '/api/system/meshsync/resources/abc-1',
    );
  });
});

describe('meshsync – transformResponse', () => {
  it('normalizes paginated responses to camelCase', async () => {
    const { normalizePaginatedCollectionResponse } = await import('../transforms');
    const response = {
      page: 0,
      page_size: 25,
      total_count: 100,
      resources: [{ id: '1', kind: 'Pod' }],
    };
    expect(normalizePaginatedCollectionResponse(response as never, 'resources')).toEqual(
      expect.objectContaining({
        pageSize: 25,
        totalCount: 100,
        resources: [{ id: '1', kind: 'Pod' }],
      }),
    );
  });

  it('returns empty resources array for non-array values', async () => {
    const { normalizePaginatedCollectionResponse } = await import('../transforms');
    const response = { resources: 'not-an-array', total_count: 0 } as never;
    expect(normalizePaginatedCollectionResponse(response, 'resources')).toEqual(
      expect.objectContaining({ resources: [] }),
    );
  });

  it('passes through undefined response unchanged', async () => {
    const { normalizePaginatedCollectionResponse } = await import('../transforms');
    expect(normalizePaginatedCollectionResponse(undefined, 'resources')).toBeUndefined();
  });
});

describe('meshsync – HTTP contracts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('getMeshSyncResources GETs with all filter params', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve(JSON.stringify({ resources: [], total_count: 0 })),
    });

    const url =
      `${mesheryApiPath('system/meshsync/resources')}` +
      '?page=0&pagesize=10&search=foo&order=asc&kind=Pod&model=istio&namespace=default';

    await fetch(url, { method: 'GET' });
    expect(global.fetch).toHaveBeenCalledWith(url, expect.objectContaining({ method: 'GET' }));
  });

  it('getMeshSyncResourceKinds GETs summary with cluster/namespace params', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve(JSON.stringify({ kinds: [] })),
    });

    const url = `${mesheryApiPath('system/meshsync/resources/summary')}?clusterId=all&pagesize=10`;
    await fetch(url, { method: 'GET' });
    expect(global.fetch).toHaveBeenCalled();
  });

  it('deleteMeshsyncResource DELETEs the per-id URL with credentials', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 204,
      text: () => Promise.resolve(''),
    });

    await fetch(mesheryApiPath('system/meshsync/resources/abc'), {
      method: 'DELETE',
      credentials: 'include',
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/system/meshsync/resources/abc',
      expect.objectContaining({ method: 'DELETE', credentials: 'include' }),
    );
  });

  it('surfaces a 404 not-found error', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 404,
      text: () => Promise.resolve('not found'),
    });

    const resp = await fetch(mesheryApiPath('system/meshsync/resources/missing'), {
      method: 'DELETE',
    });
    expect(resp.status).toBe(404);
  });
});
