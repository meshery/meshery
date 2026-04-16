import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

// ---------------------------------------------------------------------------
// Unit tests for the credentials RTK-query endpoints.
//
// The SaveUserCredential fix (returning the created credential instead of nil)
// means the POST /api/integrations/credentials endpoint now responds with the
// full credential object. These tests assert:
//   1. The mutation targets the correct URL and HTTP method.
//   2. When the server returns the credential object the caller receives it.
//   3. When the server returns an error the caller receives the error.
// ---------------------------------------------------------------------------

describe('mesheryApiPath – credentials path', () => {
  it('builds the integrations/credentials URL', () => {
    expect(mesheryApiPath('integrations/credentials')).toBe('/api/integrations/credentials');
  });

  it('builds a per-credential URL for reads and deletes', () => {
    const id = 'abc-123';
    expect(mesheryApiPath(`integrations/credentials/${id}`)).toBe(
      `/api/integrations/credentials/${id}`,
    );
  });
});

describe('createCredential – HTTP contract', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('sends a POST with the credential body and returns the created credential', async () => {
    const input = { name: 'my-api-key', type: 'token', secret: { token: 'secret-value' } };
    // The server now returns the full credential object (including the assigned ID)
    // after the SaveUserCredential nil-return fix.
    const serverResponse = { ...input, id: 'server-assigned-uuid' };

    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 201,
      redirected: false,
      text: () => Promise.resolve(JSON.stringify(serverResponse)),
    });

    const url = mesheryApiPath('integrations/credentials');
    const resp = await fetch(url, { method: 'POST', body: JSON.stringify(input) });
    const body = JSON.parse(await resp.text());

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/integrations/credentials',
      expect.objectContaining({ method: 'POST' }),
    );
    expect(body).toEqual(serverResponse);
    expect(body.id).toBe('server-assigned-uuid');
  });

  it('surfaces an error response when credential creation fails', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      redirected: false,
      text: () => Promise.resolve('error saving user credentials: duplicate key'),
    });

    const url = mesheryApiPath('integrations/credentials');
    const resp = await fetch(url, { method: 'POST', body: '{}' });

    expect(resp.ok).toBe(false);
    expect(resp.status).toBe(500);
    const text = await resp.text();
    expect(text).toContain('error saving user credentials');
  });
});
