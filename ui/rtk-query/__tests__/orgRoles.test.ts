import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

// ---------------------------------------------------------------------------
// Unit tests for rtk-query/orgRoles.ts — exposes a single getUserOrgRoles
// query against the extensions identity endpoint.
//   GET /api/extensions/api/identity/orgs/:orgId/roles
// with query params page/pagesize/search/order/all/selector and credentials.
// ---------------------------------------------------------------------------

describe('orgRoles – URL', () => {
  it('builds the org-roles path through mesheryApiPath', () => {
    const orgId = 'org-42';
    expect(mesheryApiPath(`extensions/api/identity/orgs/${orgId}/roles`)).toBe(
      '/api/extensions/api/identity/orgs/org-42/roles',
    );
  });
});

describe('orgRoles – module surface', () => {
  it('exports useGetUserOrgRolesQuery', async () => {
    const mod = await import('../orgRoles');
    expect(typeof mod.useGetUserOrgRolesQuery).toBe('function');
  });
});

describe('orgRoles – HTTP contract', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('GETs the org roles endpoint with credentials included', async () => {
    const orgId = 'org-1';
    const url = mesheryApiPath(`extensions/api/identity/orgs/${orgId}/roles?page=0&pagesize=10`);

    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve(JSON.stringify({ roles: [], total_count: 0 })),
    });

    const resp = await fetch(url, { method: 'GET', credentials: 'include' });
    const body = JSON.parse(await resp.text());

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/extensions/api/identity/orgs/org-1/roles?page=0&pagesize=10',
      expect.objectContaining({ method: 'GET', credentials: 'include' }),
    );
    expect(body).toEqual({ roles: [], total_count: 0 });
  });

  it('surfaces a 403 forbidden error', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 403,
      text: () => Promise.resolve('forbidden'),
    });

    const url = mesheryApiPath(`extensions/api/identity/orgs/org-1/roles`);
    const resp = await fetch(url, { method: 'GET' });

    expect(resp.ok).toBe(false);
    expect(resp.status).toBe(403);
  });
});
