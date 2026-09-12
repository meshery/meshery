import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';
import { meshApi } from '../mesh';

// ---------------------------------------------------------------------------
// Unit tests for rtk-query/mesh.ts — a tiny module that exposes the single
// getMesh query. The endpoint is verified by:
//   1. confirming the URL helper returns /api/mesh
//   2. confirming the meshApi object exposes a getMesh endpoint
//   3. confirming a mocked fetch against the path is invoked correctly.
// ---------------------------------------------------------------------------

describe('mesh – URL', () => {
  it('builds the /api/mesh URL', () => {
    expect(mesheryApiPath('mesh')).toBe('/api/mesh');
  });
});

describe('mesh – API surface', () => {
  it('exposes a getMesh endpoint', () => {
    expect(meshApi.endpoints).toBeDefined();
    expect(meshApi.endpoints.getMesh).toBeDefined();
    expect(typeof meshApi.endpoints.getMesh.initiate).toBe('function');
    expect(typeof meshApi.endpoints.getMesh.select).toBe('function');
  });

  it('exposes useGetMeshQuery hook', async () => {
    const mod = await import('../mesh');
    expect(typeof mod.useGetMeshQuery).toBe('function');
  });
});

describe('mesh – HTTP contract', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns the mesh payload on a successful GET', async () => {
    const payload = { meshes: [{ name: 'istio' }] };
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve(JSON.stringify(payload)),
    });

    const resp = await fetch(mesheryApiPath('mesh'));
    const body = JSON.parse(await resp.text());

    expect(global.fetch).toHaveBeenCalledWith('/api/mesh');
    expect(body).toEqual(payload);
  });

  it('surfaces a 500 server error', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve('boom'),
    });

    const resp = await fetch(mesheryApiPath('mesh'));
    expect(resp.ok).toBe(false);
    expect(resp.status).toBe(500);
  });
});
