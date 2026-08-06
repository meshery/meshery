import React from 'react';
import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

// `../utils` (appendInvalidatesTags) imports the app store at module scope, so
// the store has to be stubbed before workspace.ts is loaded.
vi.mock('../../store', () => ({ store: { dispatch: vi.fn() } }));

vi.mock('@/utils/utils', () => ({
  urlEncodeParams: (params: Record<string, unknown>) => {
    const usp = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v === undefined || v === null) return;
      usp.append(k, String(v));
    });
    return usp.toString();
  },
}));

beforeAll(() => {
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
});

// ---------------------------------------------------------------------------
// workspace.ts's create/update/delete hooks are wrappers: callers pass
// `{ workspaceId, workspacePayload }`, the schemas endpoints take
// `{ workspaceId, body }`. These tests drive the real hooks so the mapping
// itself is exercised - dispatching the endpoint directly would assert the
// schemas contract and never touch the wrapper.
//
// The wrappers used to emit the shape of workspace.ts's own local endpoint
// declarations (`{name, description, organizationId}` and `{id}`). Those
// declarations never ran: `injectEndpoints` without `overrideExisting: true`
// silently ignores a name the schemas package already defines, so the schemas
// endpoints served every call and read `body` and `workspaceId` as `undefined`.
// Create and update sent an empty body; update and delete addressed
// `/api/workspaces/undefined`. Every assertion below fails against that.
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

const setup = async () => {
  vi.resetModules();
  const apiMod = await import('../index');
  const workspaceMod = await import('../workspace');
  const { configureStore } = await import('@reduxjs/toolkit');
  const { Provider } = await import('react-redux');

  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });

  const wrapper = ({ children }: { children: React.ReactNode }) =>
    React.createElement(Provider, { store }, children);

  return { workspaceMod, wrapper };
};

describe('workspace mutation wrappers', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('useCreateWorkspaceMutation sends the payload as the request body', async () => {
    fetchMock.mockResolvedValue(okResponse({ id: 'new' }));
    const { workspaceMod, wrapper } = await setup();

    const { result } = renderHook(() => workspaceMod.useCreateWorkspaceMutation(), { wrapper });
    await act(async () => {
      await result.current[0]({
        workspacePayload: { name: 'w-name', description: 'desc', organization_id: 'org-1' },
      });
    });

    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/workspaces');
    expect(await req.clone().json()).toEqual({
      name: 'w-name',
      description: 'desc',
      organizationId: 'org-1',
    });
  });

  it('useUpdateWorkspaceMutation addresses the workspace by id and sends a body', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { workspaceMod, wrapper } = await setup();

    const { result } = renderHook(() => workspaceMod.useUpdateWorkspaceMutation(), { wrapper });
    await act(async () => {
      await result.current[0]({
        workspaceId: 'w-1',
        workspacePayload: { name: 'updated', description: 'd', organizationId: 'org-1' },
      });
    });

    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('PUT');
    expect(req.url).toContain('/api/workspaces/w-1');
    expect(req.url).not.toContain('undefined');
    expect(await req.clone().json()).toEqual({
      name: 'updated',
      description: 'd',
      organizationId: 'org-1',
    });
  });

  it('useDeleteWorkspaceMutation addresses the workspace by id', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { workspaceMod, wrapper } = await setup();

    const { result } = renderHook(() => workspaceMod.useDeleteWorkspaceMutation(), { wrapper });
    await act(async () => {
      await result.current[0]({ workspaceId: 'w-2' });
    });

    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/workspaces/w-2');
    expect(req.url).not.toContain('undefined');
  });
});
