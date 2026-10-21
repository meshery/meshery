import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { rtkErrorMiddleware } from '../rtkErrorMiddleware';
import { pushEvent } from '../../slices/events';

const makeRejectedAction = (overrides: any = {}) => {
  const { meta: metaOverrides, payload, error, ...rest } = overrides;

  return {
    type: 'someApi/executeQuery/rejected',
    meta: {
      arg: { endpointName: 'getThing' },
      requestId: 'req-1',
      requestStatus: 'rejected',
      rejectedWithValue: true,
      ...metaOverrides,
    },
    payload: payload ?? { status: 500, data: { message: 'Server exploded' } },
    error: { message: 'Rejected', ...error },
    ...rest,
  };
};

const makeStoreApi = () => ({
  dispatch: vi.fn(),
  getState: vi.fn(),
});

describe('rtkErrorMiddleware', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2024-06-01T00:00:00Z'));
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('passes the action through unchanged via next()', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn((action) => action);
    const action = makeRejectedAction();
    const result = rtkErrorMiddleware(storeApi as any)(next)(action);
    expect(next).toHaveBeenCalledWith(action);
    expect(result).toBe(action);
  });

  it('dispatches a pushEvent for rejected actions with non-401 status', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    const action = makeRejectedAction({
      payload: { status: 500, data: { message: 'Boom' } },
    });
    rtkErrorMiddleware(storeApi as any)(next)(action);

    expect(storeApi.dispatch).toHaveBeenCalledTimes(1);
    const dispatched = storeApi.dispatch.mock.calls[0][0];
    expect(dispatched.type).toBe(pushEvent({} as any).type);
    expect(dispatched.payload.severity).toBe('error');
    expect(dispatched.payload.description).toBe('Unable to load thing: Boom');
    expect(dispatched.payload.action).toBe('api_error');
    expect(dispatched.payload.category).toBe('api');
    // Avoid a race against wall-clock millisecond drift on slow CI runners:
    // assert createdAt is a parseable ISO string within a second of "now".
    expect(typeof dispatched.payload.createdAt).toBe('string');
    const createdAtMs = Date.parse(dispatched.payload.createdAt);
    expect(Number.isFinite(createdAtMs)).toBe(true);
    expect(Math.abs(createdAtMs - Date.now())).toBeLessThan(1000);
    expect(typeof dispatched.payload.id).toBe('string');
    expect(dispatched.payload.id).toMatch(/^rtk-error-\d+$/);
  });

  it('dispatches a visible error event for 401 status instead of swallowing it', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    const action = makeRejectedAction({
      payload: { status: 401, data: { error: 'Unauthorized' } },
    });
    rtkErrorMiddleware(storeApi as any)(next)(action);

    expect(storeApi.dispatch).toHaveBeenCalledTimes(1);
    const dispatched = storeApi.dispatch.mock.calls[0][0];
    expect(dispatched.type).toBe(pushEvent({} as any).type);
    expect(dispatched.payload.severity).toBe('error');
    expect(dispatched.payload.description).toBe('Unable to load thing: Unauthorized');
    expect(next).toHaveBeenCalledWith(action);
  });

  it('does not dispatch for non-rejected actions', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    const action = { type: 'someApi/executeQuery/fulfilled', payload: { ok: true } };
    rtkErrorMiddleware(storeApi as any)(next)(action);
    expect(storeApi.dispatch).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalledWith(action);
  });

  it('prefers data.message in the error payload', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({
        payload: { status: 500, data: { message: 'Specific message' } },
      }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toBe(
      'Unable to load thing: Specific message',
    );
  });

  it('falls back to data.error when data.message is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({
        payload: { status: 500, data: { error: 'Error from data.error' } },
      }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toBe(
      'Unable to load thing: Error from data.error',
    );
  });

  it('falls back to payload.error when data is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({ payload: { status: 500, error: 'Plain error' } }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toBe(
      'Unable to load thing: Plain error',
    );
  });

  it('falls back to "Request failed (status)" when no message/error is present', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(makeRejectedAction({ payload: { status: 503 } }));
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Unable to load thing. Request failed (503).',
    );
  });

  it('uses "unknown status" when status is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(makeRejectedAction({ payload: {} }));
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Unable to load thing. Request failed (unknown status).',
    );
  });

  it('uses originalStatus when status is missing on payload', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({ payload: { originalStatus: 502 } }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Unable to load thing. Request failed (502).',
    );
  });

  it('uses "API call" when endpoint name is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    const action = {
      type: 'rejected',
      payload: { status: 500, data: { message: 'err' } },
      meta: { requestId: 'req-1', requestStatus: 'rejected', rejectedWithValue: true },
    };
    rtkErrorMiddleware(storeApi as any)(next)(action);
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Unable to complete api call: err',
    );
  });

  it('turns fetch errors into user-friendly notification copy for notification center events', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({
        meta: { arg: { endpointName: 'getEvents' } },
        payload: { status: 'FETCH_ERROR', error: 'TypeError: Failed to fetch' },
      }),
    );

    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toBe(
      'Unable to load notifications. Check your connection and try again.',
    );
    expect(console.error).toHaveBeenCalledWith(
      expect.stringContaining('[RTK Query] getEvents failed:'),
      'TypeError: Failed to fetch',
    );
  });

  it('formats timeout errors with retry-oriented copy', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({
        payload: { status: 'TIMEOUT_ERROR', error: 'Request timed out' },
      }),
    );

    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toBe(
      'Timed out while trying to load thing. Please try again.',
    );
  });

  it('logs the error via console.error', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({ payload: { status: 500, data: { message: 'foo' } } }),
    );
    expect(console.error).toHaveBeenCalledWith(
      expect.stringContaining('[RTK Query] getThing failed:'),
      'foo',
    );
  });

  it('returns the result of next(action) for rejected actions', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn((a) => ({ pass: a.type }));
    const action = makeRejectedAction();
    const result = rtkErrorMiddleware(storeApi as any)(next)(action);
    expect(result).toEqual({ pass: action.type });
  });
});
