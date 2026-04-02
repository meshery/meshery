import { describe, it, expect, vi } from 'vitest';
import { authMiddleware } from '../authMiddleware';

describe('authMiddleware', () => {
  const createFakeStore = () => ({
    dispatch: vi.fn(),
    getState: vi.fn(),
  });

  it('dispatches SESSION_EXPIRED on 401 rejected action', () => {
    const store = createFakeStore();
    const next = vi.fn();

    const action = {
      type: 'api/executeQuery/rejected',
      meta: { rejectedWithValue: true },
      payload: { status: 401 },
    };

    // Simulate RTK's isRejectedWithValue check
    Object.defineProperty(action, 'meta', {
      value: { ...action.meta, rejectedWithValue: true },
    });

    authMiddleware(store)(next)(action);

    expect(next).toHaveBeenCalledWith(action);
  });

  it('passes through non-401 rejected actions without dispatching', () => {
    const store = createFakeStore();
    const next = vi.fn();

    const action = {
      type: 'api/executeQuery/rejected',
      payload: { status: 500 },
      meta: { rejectedWithValue: true },
    };

    authMiddleware(store)(next)(action);

    expect(next).toHaveBeenCalledWith(action);
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
