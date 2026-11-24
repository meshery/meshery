import { describe, expect, it, vi, beforeEach, afterEach, afterAll } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { api, mesheryApiPath } from '../index';

// The list endpoint is read by fetchBaseQuery as a (possibly relative) URL;
// pin an absolute base like the other rtk-query suites do (vi.hoisted runs
// ahead of the imports).
const { previousEndpointPrefix } = vi.hoisted(() => {
  const previous = process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
  return { previousEndpointPrefix: previous };
});

afterAll(() => {
  if (previousEndpointPrefix === undefined) {
    delete process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  } else {
    process.env.RTK_MESHERY_ENDPOINT_PREFIX = previousEndpointPrefix;
  }
});

// ---------------------------------------------------------------------------
// Unit tests for rtk-query/filter.ts. Endpoints managed:
//   GET    /api/filter                       getFilters
//   POST   /api/filter/clone/:id             cloneFilter
//   POST   /api/filter/catalog/publish       publishFilter
//   DELETE /api/filter/catalog/unpublish     unpublishFilter
//   DELETE /api/filter/:id                   deleteFilter
//   POST   /api/filter                       updateFilterFile
//   POST   /api/filter (octet-stream)        uploadFilterFile
// ---------------------------------------------------------------------------

describe('filter – URLs', () => {
  it('builds the base /filter URL', () => {
    expect(mesheryApiPath('filter')).toBe('/api/filter');
  });

  it('builds /filter/clone/:id', () => {
    expect(mesheryApiPath('filter/clone/abc-123')).toBe('/api/filter/clone/abc-123');
  });

  it('builds /filter/catalog/publish', () => {
    expect(mesheryApiPath('filter/catalog/publish')).toBe('/api/filter/catalog/publish');
  });

  it('builds /filter/catalog/unpublish', () => {
    expect(mesheryApiPath('filter/catalog/unpublish')).toBe('/api/filter/catalog/unpublish');
  });

  it('builds /filter/:id for delete', () => {
    expect(mesheryApiPath('filter/abc-123')).toBe('/api/filter/abc-123');
  });
});

describe('filter – module surface', () => {
  it('exposes all expected hooks', async () => {
    const mod = await import('../filter');
    expect(typeof mod.useGetFiltersQuery).toBe('function');
    expect(typeof mod.useCloneFilterMutation).toBe('function');
    expect(typeof mod.usePublishFilterMutation).toBe('function');
    expect(typeof mod.useUnpublishFilterMutation).toBe('function');
    expect(typeof mod.useDeleteFilterMutation).toBe('function');
    expect(typeof mod.useUpdateFilterFileMutation).toBe('function');
    expect(typeof mod.useUploadFilterFileMutation).toBe('function');
  });
});

describe('filter – HTTP contracts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('getFilters issues a GET with paging/search/order params', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve(JSON.stringify({ filters: [], total_count: 0 })),
    });

    const url = `${mesheryApiPath('filter')}?page=0&pagesize=10&order=asc&visibility=public&search=istio`;
    await fetch(url, { method: 'GET' });

    expect(global.fetch).toHaveBeenCalledWith(url, expect.objectContaining({ method: 'GET' }));
  });

  it('cloneFilter posts the body to the per-id clone endpoint', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    const body = { name: 'cloned-filter' };
    await fetch(mesheryApiPath('filter/clone/abc'), {
      method: 'POST',
      body: JSON.stringify(body),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/filter/clone/abc',
      expect.objectContaining({ method: 'POST', body: JSON.stringify(body) }),
    );
  });

  it('publishFilter posts publishBody to /filter/catalog/publish', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('filter/catalog/publish'), {
      method: 'POST',
      body: JSON.stringify({ id: 'f-1' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/filter/catalog/publish',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('unpublishFilter DELETEs /filter/catalog/unpublish with body', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('filter/catalog/unpublish'), {
      method: 'DELETE',
      body: JSON.stringify({ id: 'f-1' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/filter/catalog/unpublish',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('deleteFilter DELETEs the per-id URL', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 204,
      text: () => Promise.resolve(''),
    });

    await fetch(mesheryApiPath('filter/abc'), { method: 'DELETE' });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/filter/abc',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('uploadFilterFile posts with octet-stream content-type', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    const uploadBody = new ArrayBuffer(8);
    await fetch(mesheryApiPath('filter'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/octet-stream' },
      body: uploadBody,
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/filter',
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({ 'Content-Type': 'application/octet-stream' }),
      }),
    );
  });

  it('surfaces a 404 not-found error on getFilters', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 404,
      text: () => Promise.resolve('not found'),
    });

    const resp = await fetch(mesheryApiPath('filter'), { method: 'GET' });
    expect(resp.ok).toBe(false);
    expect(resp.status).toBe(404);
  });
});

describe('filter – cache invalidation', () => {
  const MUTATIONS = [
    { name: 'cloneFilter', arg: { filterID: 'f-1', body: { name: 'clone' } } },
    { name: 'publishFilter', arg: { publishBody: { id: 'f-1' } } },
    { name: 'unpublishFilter', arg: { unpublishBody: { id: 'f-1' } } },
    { name: 'deleteFilter', arg: { id: 'f-1' } },
    { name: 'updateFilterFile', arg: { updateBody: { id: 'f-1' } } },
    { name: 'uploadFilterFile', arg: { uploadBody: new ArrayBuffer(4) } },
  ] as const;

  const setup = () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(
        new Response(JSON.stringify({ filters: [], total_count: 0 }), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        }),
      ),
    );
    vi.stubGlobal('fetch', fetchMock);

    const store = configureStore({
      reducer: { [api.reducerPath]: api.reducer },
      middleware: (getDefault) => getDefault().concat(api.middleware),
    });

    return { fetchMock, store };
  };

  const urlOf = (call: unknown[]) =>
    typeof call[0] === 'string' ? call[0] : (call[0] as Request).url;

  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('refetches the filters list after every mutation', async () => {
    const { fetchMock, store } = setup();
    // filter.ts injects its endpoints into the shared api when imported.
    await import('../filter');

    for (const { name, arg } of MUTATIONS) {
      const listSub = store.dispatch(api.endpoints.getFilters.initiate({ page: 0, pagesize: 10 }));
      await listSub;
      const fetchesAfterList = fetchMock.mock.calls.length;

      await store.dispatch(api.endpoints[name].initiate(arg as never)).unwrap();

      // Let the invalidation-driven refetch settle.
      await new Promise((resolve) => setTimeout(resolve, 0));

      const listRefetches = fetchMock.mock.calls
        .slice(fetchesAfterList)
        .filter((call) => urlOf(call).includes('/api/filter?'));
      expect(
        listRefetches.length,
        `${name} should invalidate the filters tag so the list refetches`,
      ).toBeGreaterThan(0);

      fetchMock.mockClear();
      listSub.unsubscribe();
    }
  });
});
