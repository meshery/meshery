import { describe, expect, it, vi, beforeEach } from 'vitest';

// The module under test calls `fetchMesheryOperatorStatus(...).subscribe(...)`
// (a Relay-style fetchQuery wrapper). We mock the import so we can drive the
// next/error callbacks directly without standing up a Relay environment.
const fetchMock = vi.fn();
vi.mock('@/graphql/queries/OperatorStatusQuery', () => ({
  default: (...args: unknown[]) => fetchMock(...args),
}));

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
  let unsubscribe: ReturnType<typeof vi.fn>;
  let savedHandlers: { next?: (v: unknown) => void; error?: (e: unknown) => void };

  beforeEach(() => {
    fetchMock.mockReset();
    unsubscribe = vi.fn();
    savedHandlers = {};
    fetchMock.mockImplementation(() => ({
      subscribe: ({ next, error }: { next: (v: unknown) => void; error: (e: unknown) => void }) => {
        savedHandlers.next = next;
        savedHandlers.error = error;
        return { unsubscribe };
      },
    }));
  });

  it('passes the connection id to fetchMesheryOperatorStatus', () => {
    pingMesheryOperator('conn-1', vi.fn(), vi.fn());
    expect(fetchMock).toHaveBeenCalledWith({ connectionID: 'conn-1' });
  });

  it('invokes the success handler and unsubscribes on a healthy response', () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const data = { operator: { status: 'ENABLED' } };
    savedHandlers.next?.(data);

    expect(success).toHaveBeenCalledWith(data);
    expect(error).not.toHaveBeenCalled();
    expect(unsubscribe).toHaveBeenCalledTimes(1);
  });

  it('invokes the error handler when the response is null', () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    savedHandlers.next?.(null);

    expect(error).toHaveBeenCalledWith(null);
    expect(success).not.toHaveBeenCalled();
    expect(unsubscribe).toHaveBeenCalledTimes(1);
  });

  it('invokes the error handler when operator is null', () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const data = { operator: null };
    savedHandlers.next?.(data);

    expect(error).toHaveBeenCalledWith(data);
    expect(success).not.toHaveBeenCalled();
  });

  it('invokes the error handler when operator status is UNKOWN', () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const data = { operator: { status: 'UNKOWN' } };
    savedHandlers.next?.(data);

    expect(error).toHaveBeenCalledWith(data);
    expect(success).not.toHaveBeenCalled();
  });

  it('routes raw errors through the error handler and unsubscribes', () => {
    const success = vi.fn();
    const error = vi.fn();
    pingMesheryOperator('conn-1', success, error);

    const e = new Error('socket closed');
    savedHandlers.error?.(e);

    expect(error).toHaveBeenCalledWith(e);
    expect(unsubscribe).toHaveBeenCalledTimes(1);
  });

  it('substitutes a default error when the subscription emits a nullish error', () => {
    const error = vi.fn();
    pingMesheryOperator('conn-1', vi.fn(), error);

    savedHandlers.error?.(null);

    expect(error).toHaveBeenCalledTimes(1);
    const arg = error.mock.calls[0][0];
    expect(arg).toBeInstanceOf(Error);
    expect((arg as Error).message).toBe('Unknown error from pingMesheryOperator');
  });

  it('does not throw when the success callback is undefined', () => {
    pingMesheryOperator('conn-1', undefined as never, vi.fn());
    expect(() => savedHandlers.next?.({ operator: { status: 'ENABLED' } })).not.toThrow();
    expect(unsubscribe).toHaveBeenCalledTimes(1);
  });

  it('does not throw when the error callback is undefined', () => {
    pingMesheryOperator('conn-1', vi.fn(), undefined as never);
    expect(() => savedHandlers.next?.(null)).not.toThrow();
    expect(unsubscribe).toHaveBeenCalledTimes(1);
    expect(() => savedHandlers.error?.(new Error('boom'))).not.toThrow();
  });
});
