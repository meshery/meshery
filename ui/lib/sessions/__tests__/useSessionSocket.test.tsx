import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useSessionSocket } from '../useSessionSocket';

/**
 * Captures the handlers each SessionSocket was constructed with, so a test can
 * fire a discarded socket's callbacks *after* its replacement exists — which is
 * the ordering React's StrictMode double-mount actually produces.
 */
const constructed: Array<{
  url: string;
  handlers: {
    onReady?: (caps: unknown) => void;
    onData?: (chunk: Uint8Array) => void;
    onClose?: (info: unknown) => void;
  };
  close: ReturnType<typeof vi.fn>;
}> = [];

vi.mock('../sessionSocket', () => ({
  SessionSocket: class {
    close = vi.fn();
    constructor(url: string, handlers: Record<string, unknown>) {
      constructed.push({ url, handlers, close: this.close });
    }
  },
}));

beforeEach(() => {
  constructed.length = 0;
});

describe('useSessionSocket', () => {
  it('drops data from a socket it has already torn down', () => {
    const onData = vi.fn();
    const { rerender } = renderHook(({ url }) => useSessionSocket({ url, onData }), {
      initialProps: { url: 'ws://host/a' },
    });

    rerender({ url: 'ws://host/b' });
    const [first, second] = constructed;

    act(() => {
      first.handlers.onData?.(new Uint8Array([1, 2, 3]));
    });
    expect(onData).not.toHaveBeenCalled();

    act(() => {
      second.handlers.onData?.(new Uint8Array([4, 5, 6]));
    });
    expect(onData).toHaveBeenCalledTimes(1);
  });

  it('keeps the live socket’s status after a superseded socket closes', () => {
    const { result, rerender } = renderHook(({ url }) => useSessionSocket({ url }), {
      initialProps: { url: 'ws://host/a' },
    });

    rerender({ url: 'ws://host/b' });
    const [first, second] = constructed;

    act(() => {
      second.handlers.onReady?.({ terminal: true, logs: false });
    });
    expect(result.current.status).toBe('open');

    act(() => {
      first.handlers.onClose?.({ graceful: false, message: 'The connection was lost.' });
    });

    expect(result.current.status).toBe('open');
    expect(result.current.closed).toBeNull();
  });

  it('holds the socket shut when there is no url', () => {
    renderHook(() => useSessionSocket({ url: null }));
    expect(constructed).toHaveLength(0);
  });
});
