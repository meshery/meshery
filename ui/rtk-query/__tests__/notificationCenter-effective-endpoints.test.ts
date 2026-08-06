import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

// notificationCenter.ts dispatches action creators from '../../store/slices/events'.
// Mock them so no real store is needed.
const { eventActions } = vi.hoisted(() => ({
  eventActions: {
    deleteEvent: vi.fn((p) => ({ type: 'events/deleteEvent', payload: p })),
    deleteEvents: vi.fn((p) => ({ type: 'events/deleteEvents', payload: p })),
    updateEventStatus: vi.fn((p) => ({ type: 'events/updateEventStatus', payload: p })),
    updateEvents: vi.fn((p) => ({ type: 'events/updateEvents', payload: p })),
  },
}));
vi.mock('../../store/slices/events', () => eventActions);
vi.mock('../store/slices/events', () => eventActions);

beforeAll(() => {
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
});

// ---------------------------------------------------------------------------
// Meshery Server serves the notification endpoints under `/api/system/events`.
// @meshery/schemas defines an overlapping `deleteEvent` operation pointed at
// `/api/events/{eventId}` - the path the server is *going* to serve once the
// UI is fully on the new events API and SSE is off (see the comment beside that
// route in server/router/server.go, and meshery/schemas#1134).
//
// notificationCenter.ts injected with `overrideExisting: false`, so RTK Query
// silently dropped the local `deleteEvent` and served every call from the
// schemas definition: `DELETE /api/events/undefined` - the wrong path, and
// `undefined` because callers pass `{ id }` while the schemas endpoint reads
// `eventId`. Deleting a notification could not work.
//
// These assertions dispatch through a real store, so they pin the definition
// that ACTUALLY serves the call rather than the one the module happens to
// declare. Drop this file's premise (and the `overrideExisting: true` it
// guards) when the server moves to `/api/events`.
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
  const notificationCenterMod = await import('../notificationCenter');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { store, notificationCenterApi: notificationCenterMod.notificationCenterApi };
};

describe('notificationCenter effective endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn();
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('deleteEvent DELETEs /api/system/events/:id, not the schemas /api/events path', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { store, notificationCenterApi } = await setupStore();

    await store.dispatch(notificationCenterApi.endpoints.deleteEvent.initiate({ id: 'e-1' }));

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/system/events/e-1');
    expect(req.url).not.toContain('undefined');
  });

  it('updateStatus PUTs /api/system/events/status/:id', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { store, notificationCenterApi } = await setupStore();

    await store.dispatch(
      notificationCenterApi.endpoints.updateStatus.initiate({ id: 'e-2', status: 'read' }),
    );

    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('PUT');
    expect(req.url).toContain('/api/system/events/status/e-2');
    expect(req.url).not.toContain('undefined');
  });

  it('deleteEvents DELETEs the bulk /api/system/events/bulk path', async () => {
    fetchMock.mockResolvedValue(okResponse({}));
    const { store, notificationCenterApi } = await setupStore();

    await store.dispatch(notificationCenterApi.endpoints.deleteEvents.initiate({ ids: ['e-3'] }));

    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('DELETE');
    expect(req.url).toContain('/api/system/events/bulk');
  });
});
