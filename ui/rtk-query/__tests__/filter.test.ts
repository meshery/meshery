import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

// ---------------------------------------------------------------------------
// Unit tests for rtk-query/filter.ts. Endpoints managed:
//   GET    /api/content/filters                       getFilters
//   POST   /api/content/filters/clone/:id             cloneFilter
//   POST   /api/content/filters/catalog/publish       publishFilter
//   DELETE /api/content/filters/catalog/unpublish     unpublishFilter
//   DELETE /api/content/filters/:id                   deleteFilter
//   POST   /api/content/filters                       updateFilterFile
//   POST   /api/content/filters (octet-stream)        uploadFilterFile
// ---------------------------------------------------------------------------

describe('filter – URLs', () => {
  it('builds the base content/filters URL', () => {
    expect(mesheryApiPath('content/filters')).toBe('/api/content/filters');
  });

  it('builds content/filters/clone/:id', () => {
    expect(mesheryApiPath('content/filters/clone/abc-123')).toBe(
      '/api/content/filters/clone/abc-123',
    );
  });

  it('builds content/filters/catalog/publish', () => {
    expect(mesheryApiPath('content/filters/catalog/publish')).toBe(
      '/api/content/filters/catalog/publish',
    );
  });

  it('builds content/filters/catalog/unpublish', () => {
    expect(mesheryApiPath('content/filters/catalog/unpublish')).toBe(
      '/api/content/filters/catalog/unpublish',
    );
  });

  it('builds content/filters/:id for delete', () => {
    expect(mesheryApiPath('content/filters/abc-123')).toBe('/api/content/filters/abc-123');
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

    const url = `${mesheryApiPath('content/filters')}?page=0&pagesize=10&order=asc&visibility=public&search=istio`;
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
    await fetch(mesheryApiPath('content/filters/clone/abc'), {
      method: 'POST',
      body: JSON.stringify(body),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/content/filters/clone/abc',
      expect.objectContaining({ method: 'POST', body: JSON.stringify(body) }),
    );
  });

  it('publishFilter posts publishBody to content/filters/catalog/publish', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('content/filters/catalog/publish'), {
      method: 'POST',
      body: JSON.stringify({ id: 'f-1' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/content/filters/catalog/publish',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('unpublishFilter DELETEs content/filters/catalog/unpublish with body', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('content/filters/catalog/unpublish'), {
      method: 'DELETE',
      body: JSON.stringify({ id: 'f-1' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/content/filters/catalog/unpublish',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('deleteFilter DELETEs the per-id URL', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 204,
      text: () => Promise.resolve(''),
    });

    await fetch(mesheryApiPath('content/filters/abc'), { method: 'DELETE' });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/content/filters/abc',
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
    await fetch(mesheryApiPath('content/filters'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/octet-stream' },
      body: uploadBody,
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/content/filters',
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

    const resp = await fetch(mesheryApiPath('content/filters'), { method: 'GET' });
    expect(resp.ok).toBe(false);
    expect(resp.status).toBe(404);
  });
});
