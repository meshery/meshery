import { describe, it, expect, vi } from 'vitest';
import { authMiddleware } from '../authMiddleware';

function makeRejectedWithValueAction(payload) {
  return {
    type: 'api/executeQuery/rejected',
    payload,
    meta: {
      rejectedWithValue: true,
      requestId: 'test-req',
      requestStatus: 'rejected',
      aborted: false,
      condition: false,
    },
    error: { message: 'Rejected' },
  };
}

describe('authMiddleware', () => {
  const createFakeStore = () => ({
    dispatch: vi.fn(),
    getState: vi.fn(),
  });

  it('dispatches SESSION_EXPIRED on 401 rejected action', () => {
    const store = createFakeStore();
    const next = vi.fn();
    const action = makeRejectedWithValueAction({ status: 401 });

    authMiddleware(store)(next)(action);

    expect(next).toHaveBeenCalledWith(action);
    expect(store.dispatch).toHaveBeenCalledWith({ type: 'SESSION_EXPIRED' });
    expect(store.dispatch).toHaveBeenCalledTimes(1);
  });

  it('passes through non-401 rejected actions without dispatching SESSION_EXPIRED', () => {
    const store = createFakeStore();
    const next = vi.fn();
    const action = makeRejectedWithValueAction({ status: 500 });

    authMiddleware(store)(next)(action);

    expect(next).toHaveBeenCalledWith(action);
    expect(store.dispatch).not.toHaveBeenCalled();
  });

  it('passes through fulfilled actions without dispatching', () => {
    const store = createFakeStore();
    const next = vi.fn();
    const action = {
      type: 'api/executeQuery/fulfilled',
      payload: { data: 'ok' },
    };

    authMiddleware(store)(next)(action);

    expect(next).toHaveBeenCalledWith(action);
    expect(store.dispatch).not.toHaveBeenCalled();
  });
});
