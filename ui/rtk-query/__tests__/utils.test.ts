import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// initiateQuery dispatches an RTK-Query initiate action via the application
// store. We mock the store module so the helper can be exercised in isolation
// without instantiating the real RTK-Query middleware (which would require
// configuring an entire reducer tree).
const { dispatch } = vi.hoisted(() => ({ dispatch: vi.fn() }));
vi.mock('../../store', () => ({ store: { dispatch } }));

import { initiateQuery } from '../utils';

describe('initiateQuery', () => {
  beforeEach(() => {
    dispatch.mockReset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('dispatches query.initiate(variables) via the store', async () => {
    dispatch.mockResolvedValue({ data: { ok: true } });
    const initiate = vi.fn().mockReturnValue({ type: 'thunk-action' });
    const variables = { id: 'abc' };

    const result = await initiateQuery({ initiate }, variables);

    expect(initiate).toHaveBeenCalledWith(variables);
    expect(dispatch).toHaveBeenCalledWith({ type: 'thunk-action' });
    expect(result).toEqual({ data: { ok: true } });
  });

  it('releases the cache subscription that initiate() opened', async () => {
    // The real dispatch returns a QueryActionCreatorResult, which is a promise
    // that also carries unsubscribe(). Leaving it unsubscribed pins the cache
    // entry for the life of the session.
    const unsubscribe = vi.fn();
    const subscription = Object.assign(Promise.resolve({ data: { ok: true } }), { unsubscribe });
    dispatch.mockReturnValue(subscription);
    const initiate = vi.fn().mockReturnValue({ type: 'thunk-action' });

    const result = await initiateQuery({ initiate }, { id: 'abc' });

    expect(result).toEqual({ data: { ok: true } });
    expect(unsubscribe).toHaveBeenCalledTimes(1);
  });

  it('releases the subscription even when the query rejects', async () => {
    const unsubscribe = vi.fn();
    const subscription = Object.assign(Promise.reject(new Error('boom')), { unsubscribe });
    dispatch.mockReturnValue(subscription);
    const initiate = vi.fn().mockReturnValue({ type: 'thunk-action' });

    const result = await initiateQuery({ initiate }, undefined);

    expect(result.isError).toBe(true);
    expect(unsubscribe).toHaveBeenCalledTimes(1);
  });

  it('returns the awaited dispatch result on success', async () => {
    const fulfilled = {
      data: { foo: 'bar' },
      isSuccess: true,
      isError: false,
      isLoading: false,
      isFetching: false,
    };
    dispatch.mockResolvedValue(fulfilled);
    const initiate = vi.fn().mockReturnValue({ type: 'thunk' });

    const result = await initiateQuery({ initiate }, undefined);

    expect(result).toBe(fulfilled);
  });

  it('returns an error-shaped object when dispatch throws synchronously', async () => {
    const err = new Error('initiate threw');
    const initiate = vi.fn(() => {
      throw err;
    });

    const result = await initiateQuery({ initiate }, { id: 1 });

    expect(result).toEqual({
      data: null,
      error: err,
      isFetching: false,
      isSuccess: false,
      isLoading: false,
      isError: true,
    });
  });

  it('returns an error-shaped object when dispatch rejects', async () => {
    const err = new Error('dispatch rejected');
    dispatch.mockRejectedValue(err);
    const initiate = vi.fn().mockReturnValue({ type: 'thunk' });

    const result = await initiateQuery({ initiate }, {});

    expect(result.isError).toBe(true);
    expect(result.error).toBe(err);
    expect(result.data).toBeNull();
    expect(result.isSuccess).toBe(false);
  });
});
