import { describe, expect, it, vi, beforeEach } from 'vitest';

// pingMesheryOperator now dispatches the RTK Query endpoint
// `api.endpoints.getMesheryOperatorStatus.initiate(...)` against the store,
// then consumes the returned promise via .unwrap() and .unsubscribe().
// We mock the store and the RTK endpoint so we can drive the resolve/reject
// path and observe unsubscribe — same coverage shape the prior Relay-style
// test had, just against the post-migration plumbing.
const initiateMock = vi.fn();
const dispatchMock = vi.fn();
const unsubscribeMock = vi.fn();

vi.mock('@/store/index', () => ({
  store: {
    dispatch: (...args: unknown[]) => dispatchMock(...args),
  },
}));

vi.mock('@/rtk-query/index', () => ({
  api: {
    endpoints: {
      getMesheryOperatorStatus: {
        initiate: (...args: unknown[]) => initiateMock(...args),
      },
    },
  },
}));

// The source module imports '@/rtk-query/kubernetes' purely for its endpoint-
// injection side-effect. There is no exported value to mock; an empty module
// satisfies the import.
vi.mock('@/rtk-query/kubernetes', () => ({}));

import {
  getOperatorStatusFromQueryResult,
  isMesheryOperatorConnected,
  pingMesheryOperator,
} from '../mesheryOperator';

describe('isMesheryOperatorConnected', () => {
  it('returns the operatorInstalled flag', () => {
    expect(isMesheryOperatorConnected({ operatorInstalled: true })).toBe(true);
    expect(isMesheryOperatorConnected({ operatorInstalled: false })).toBe(false);
    expect(isMesheryOperatorConnected({ operatorInstalled: undefined })).toBe(undefined);
  });
});

describe('getOperatorStatusFromQueryResult', () => {
  it('returns disconnected information when the operator entry carries an error', () => {
    const result = getOperatorStatusFromQueryResult({ operator: { error: 'boom' } });
    expect(result[0]).toBe(false);
    expect(result[1].operatorInstalled).toBe(false);
    expect(result[1].NATSInstalled).toBe(false);
    expect(result[1].meshSyncInstalled).toBe(false);
    expect(result[1].operatorVersion).toBe('N/A');
  });

  it('returns disconnected when status is not ENABLED', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: { status: 'DISABLED', controllers: [] },
    });
    expect(result[0]).toBe(false);
    expect(result[1].operatorInstalled).toBe(false);
  });

  it('marks operator installed and aggregates broker/meshsync state when ENABLED', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: {
        status: 'ENABLED',
        version: 'v0.7.0',
        controllers: [
          { name: 'broker', status: 'ENABLED', version: '2.10.0' },
          { name: 'meshsync', status: 'ENABLED', version: '0.7.5' },
        ],
      },
    });
    expect(result[0]).toBe(true);
    const info = result[1];
    expect(info.operatorInstalled).toBe(true);
    expect(info.operatorVersion).toBe('v0.7.0');
    expect(info.NATSInstalled).toBe(true);
    expect(info.NATSVersion).toBe('2.10.0');
    expect(info.meshSyncInstalled).toBe(true);
    expect(info.meshSyncVersion).toBe('0.7.5');
  });

  it('marks individual controllers as not installed when their status is not ENABLED', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: {
        status: 'ENABLED',
        version: 'v0.7.0',
        controllers: [
          { name: 'broker', status: 'DISABLED', version: '2.10.0' },
          { name: 'meshsync', status: 'DISABLED', version: '0.7.5' },
        ],
      },
    });
    expect(result[0]).toBe(true);
    expect(result[1].operatorInstalled).toBe(true);
    expect(result[1].NATSInstalled).toBe(false);
    expect(result[1].NATSVersion).toBe('N/A');
    expect(result[1].meshSyncInstalled).toBe(false);
    expect(result[1].meshSyncVersion).toBe('N/A');
  });

  it('handles ENABLED status with no controllers array', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: { status: 'ENABLED', version: 'v0.7.0' },
    });
    expect(result[0]).toBe(true);
    expect(result[1].operatorInstalled).toBe(true);
    expect(result[1].operatorVersion).toBe('v0.7.0');
  });
});

describe('pingMesheryOperator', () => {
  // Each test installs a fresh mock promise that we can resolve/reject
  // synchronously and then await microtasks to drain. The shape mirrors
  // what RTK Query's initiate() returns: an object with .unwrap() and
  // .unsubscribe(). We control resolve/reject from the test body.
  let resolveFn: (value: unknown) => void;
  let rejectFn: (reason: unknown) => void;
  let dispatchedPromise: { unwrap: () => Promise<unknown>; unsubscribe: () => void };

  beforeEach(() => {
    initiateMock.mockReset();
    dispatchMock.mockReset();
    unsubscribeMock.mockReset();

    const inner = new Promise<unknown>((res, rej) => {
      resolveFn = res;
      rejectFn = rej;
    });
    dispatchedPromise = {
      unwrap: () => inner,
      unsubscribe: unsubscribeMock,
    };

    initiateMock.mockReturnValue('initiate-action');
    dispatchMock.mockReturnValue(dispatchedPromise);
  });

  // Microtask-flush helper. The .then/.catch/.finally chain runs across
  // multiple ticks, so a single await isn't enough to drive the whole chain.
  const flushAll = async () => {
    for (let i = 0; i < 5; i++) {
      await Promise.resolve();
    }
  };

  it('dispatches getMesheryOperatorStatus.initiate with the connection id and forceRefetch', () => {
    pingMesheryOperator('conn-1', vi.fn(), vi.fn());
    expect(initiateMock).toHaveBeenCalledWith({ connectionID: 'conn-1' }, { forceRefetch: true });
    expect(dispatchMock).toHaveBeenCalledWith('initiate-action');
  });

  it('invokes the success handler and unsubscribes on a healthy response', async () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const data = { operator: { status: 'ENABLED' } };
    resolveFn(data);
    await flushAll();

    expect(success).toHaveBeenCalledWith(data);
    expect(error).not.toHaveBeenCalled();
    expect(unsubscribeMock).toHaveBeenCalledTimes(1);
  });

  it('invokes the error handler when the response is null', async () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    resolveFn(null);
    await flushAll();

    expect(error).toHaveBeenCalledWith(null);
    expect(success).not.toHaveBeenCalled();
    expect(unsubscribeMock).toHaveBeenCalledTimes(1);
  });

  it('invokes the error handler when operator is null', async () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const data = { operator: null };
    resolveFn(data);
    await flushAll();

    expect(error).toHaveBeenCalledWith(data);
    expect(success).not.toHaveBeenCalled();
  });

  it('invokes the error handler when operator status is UNKOWN', async () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const data = { operator: { status: 'UNKOWN' } };
    resolveFn(data);
    await flushAll();

    expect(error).toHaveBeenCalledWith(data);
    expect(success).not.toHaveBeenCalled();
  });

  it('routes rejected promises through the error handler and unsubscribes', async () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const e = new Error('socket closed');
    rejectFn(e);
    await flushAll();

    expect(error).toHaveBeenCalledWith(e);
    expect(unsubscribeMock).toHaveBeenCalledTimes(1);
  });

  it('substitutes a default error when the rejection is nullish', async () => {
    const error = vi.fn();
    pingMesheryOperator('conn-1', vi.fn(), error);

    rejectFn(null);
    await flushAll();

    expect(error).toHaveBeenCalledTimes(1);
    const arg = error.mock.calls[0][0];
    expect(arg).toBeInstanceOf(Error);
    expect((arg as Error).message).toBe('Unknown error from pingMesheryOperator');
  });

  it('does not throw when the success callback is undefined', async () => {
    pingMesheryOperator('conn-1', undefined as never, vi.fn());
    resolveFn({ operator: { status: 'ENABLED' } });
    await flushAll();
    expect(unsubscribeMock).toHaveBeenCalledTimes(1);
  });

  it('does not throw when the error callback is undefined on rejection', async () => {
    pingMesheryOperator('conn-1', vi.fn(), undefined as never);
    rejectFn(new Error('boom'));
    await flushAll();
    expect(unsubscribeMock).toHaveBeenCalledTimes(1);
  });

  it('does not throw when the error callback is undefined on null response', async () => {
    pingMesheryOperator('conn-1', vi.fn(), undefined as never);
    resolveFn(null);
    await flushAll();
    expect(unsubscribeMock).toHaveBeenCalledTimes(1);
  });
});
