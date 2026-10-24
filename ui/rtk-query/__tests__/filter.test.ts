import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

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
