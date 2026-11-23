import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import useInterval from '../useInterval';

describe('useInterval', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('invokes the callback every `delay` ms', () => {
    const cb = vi.fn();
    renderHook(() => useInterval(cb, 100));

    expect(cb).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(cb).toHaveBeenCalledTimes(1);

    act(() => {
      vi.advanceTimersByTime(250);
    });
    expect(cb).toHaveBeenCalledTimes(3);
  });

  it('does nothing when delay is null', () => {
    const cb = vi.fn();
    renderHook(() => useInterval(cb, null));

    act(() => {
      vi.advanceTimersByTime(10_000);
    });

    expect(cb).not.toHaveBeenCalled();
  });

  it('does nothing when delay is undefined', () => {
    const cb = vi.fn();
    renderHook(() => useInterval(cb, undefined));

    act(() => {
      vi.advanceTimersByTime(10_000);
    });

    expect(cb).not.toHaveBeenCalled();
  });

  it('does nothing when delay is negative', () => {
    const cb = vi.fn();
    renderHook(() => useInterval(cb, -100));

    act(() => {
      vi.advanceTimersByTime(10_000);
    });

    expect(cb).not.toHaveBeenCalled();
  });

  it('does nothing when delay is NaN or non-numeric', () => {
    const cb = vi.fn();
    renderHook(() => useInterval(cb, Number.NaN));
    // @ts-expect-error testing invalid delay input
    renderHook(() => useInterval(cb, '500'));

    act(() => {
      vi.advanceTimersByTime(10_000);
    });

    expect(cb).not.toHaveBeenCalled();
  });

  it('does not throw when callback is not a function', () => {
    // @ts-expect-error testing invalid callback input
    expect(() => renderHook(() => useInterval(null, 100))).not.toThrow();
    // @ts-expect-error testing invalid callback input
    expect(() => renderHook(() => useInterval(undefined, 100))).not.toThrow();

    act(() => {
      vi.advanceTimersByTime(500);
    });
  });

  it(
    'starts interval when callback transitions from non-callable to callable with unchanged delay',
    () => {
    const cb = vi.fn();
    const { rerender } = renderHook(({ fn }) => useInterval(fn, 100), {
      initialProps: { fn: null as unknown as () => void },
    });

    act(() => {
      vi.advanceTimersByTime(200);
    });
    expect(cb).not.toHaveBeenCalled();

    rerender({ fn: cb });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(cb).toHaveBeenCalledTimes(1);
  });

  it('always invokes the most recent callback without resetting the timer', () => {
    const first = vi.fn();
    const second = vi.fn();
    const { rerender } = renderHook(({ cb }) => useInterval(cb, 100), {
      initialProps: { cb: first },
    });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(first).toHaveBeenCalledTimes(1);

    rerender({ cb: second });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(first).toHaveBeenCalledTimes(1);
    expect(second).toHaveBeenCalledTimes(1);
  });

  it('switches to the new delay when it changes', () => {
    const cb = vi.fn();
    const { rerender } = renderHook(({ delay }) => useInterval(cb, delay), {
      initialProps: { delay: 100 as number | null },
    });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(cb).toHaveBeenCalledTimes(1);

    rerender({ delay: 500 });

    act(() => {
      vi.advanceTimersByTime(500);
    });
    expect(cb).toHaveBeenCalledTimes(2);
  });

  it('stops invoking the callback after unmount', () => {
    const cb = vi.fn();
    const { unmount } = renderHook(() => useInterval(cb, 50));

    act(() => {
      vi.advanceTimersByTime(50);
    });
    expect(cb).toHaveBeenCalledTimes(1);

    unmount();
    act(() => {
      vi.advanceTimersByTime(500);
    });
    expect(cb).toHaveBeenCalledTimes(1);
  });
});
