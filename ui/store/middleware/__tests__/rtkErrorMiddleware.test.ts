import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { rtkErrorMiddleware } from '../rtkErrorMiddleware';
import { pushEvent } from '../../slices/events';

const makeRejectedAction = (overrides: any = {}) => ({
  type: 'someApi/executeQuery/rejected',
  meta: {
    arg: { endpointName: 'getThing' },
    requestId: 'req-1',
    requestStatus: 'rejected',
    rejectedWithValue: true,
    ...overrides.meta,
  },
  payload: overrides.payload ?? { status: 500, data: { message: 'Server exploded' } },
  error: { message: 'Rejected' },
  ...overrides,
});

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
    expect(dispatched.payload.description).toBe('getThing: Boom');
    expect(dispatched.payload.action).toBe('api_error');
    expect(dispatched.payload.category).toBe('api');
    expect(dispatched.payload.createdAt).toBe(new Date().toISOString());
    expect(typeof dispatched.payload.id).toBe('string');
    expect(dispatched.payload.id).toMatch(/^rtk-error-\d+$/);
  });

  it('skips dispatch for 401 status (auth middleware territory)', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    const action = makeRejectedAction({ payload: { status: 401 } });
    rtkErrorMiddleware(storeApi as any)(next)(action);

    expect(storeApi.dispatch).not.toHaveBeenCalled();
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
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain('Specific message');
  });

  it('falls back to data.error when data.message is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({
        payload: { status: 500, data: { error: 'Error from data.error' } },
      }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Error from data.error',
    );
  });

  it('falls back to payload.error when data is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({ payload: { status: 500, error: 'Plain error' } }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain('Plain error');
  });

  it('falls back to "Request failed (status)" when no message/error is present', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(makeRejectedAction({ payload: { status: 503 } }));
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Request failed (503)',
    );
  });

  it('uses "unknown status" when status is missing', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(makeRejectedAction({ payload: {} }));
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Request failed (unknown status)',
    );
  });

  it('uses originalStatus when status is missing on payload', () => {
    const storeApi = makeStoreApi();
    const next = vi.fn();
    rtkErrorMiddleware(storeApi as any)(next)(
      makeRejectedAction({ payload: { originalStatus: 502 } }),
    );
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain(
      'Request failed (502)',
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
    expect(storeApi.dispatch.mock.calls[0][0].payload.description).toContain('API call:');
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
