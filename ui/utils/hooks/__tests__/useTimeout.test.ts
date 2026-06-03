import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import useTimeout from '../useTimeout';

describe('useTimeout', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('invokes the callback once after `delay` ms', () => {
    const cb = vi.fn();
    renderHook(() => useTimeout(cb, 100));

    act(() => {
      vi.advanceTimersByTime(99);
    });
    expect(cb).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(cb).toHaveBeenCalledTimes(1);

    act(() => {
      vi.advanceTimersByTime(500);
    });
    expect(cb).toHaveBeenCalledTimes(1);
  });

  it('does nothing when delay is null', () => {
    const cb = vi.fn();
    renderHook(() => useTimeout(cb, null));

    act(() => {
      vi.advanceTimersByTime(10_000);
    });
    expect(cb).not.toHaveBeenCalled();
  });

  it('always fires the latest callback', () => {
    const first = vi.fn();
    const second = vi.fn();
    const { rerender } = renderHook(({ cb }) => useTimeout(cb, 100), {
      initialProps: { cb: first },
    });

    rerender({ cb: second });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(first).not.toHaveBeenCalled();
    expect(second).toHaveBeenCalledTimes(1);
  });

  it('cancels the timer when unmounted before firing', () => {
    const cb = vi.fn();
    const { unmount } = renderHook(() => useTimeout(cb, 100));

    unmount();

    act(() => {
      vi.advanceTimersByTime(200);
    });
    expect(cb).not.toHaveBeenCalled();
  });

  it('re-runs the timer when delay changes', () => {
    const cb = vi.fn();
    const { rerender } = renderHook(({ delay }) => useTimeout(cb, delay), {
      initialProps: { delay: 100 as number | null },
    });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(cb).toHaveBeenCalledTimes(1);

    rerender({ delay: 50 });

    act(() => {
      vi.advanceTimersByTime(50);
    });
    expect(cb).toHaveBeenCalledTimes(2);
  });

  it('re-runs the timer when extra deps change', () => {
    const cb = vi.fn();
    const { rerender } = renderHook(({ dep }) => useTimeout(cb, 100, [dep]), {
      initialProps: { dep: 'a' },
    });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(cb).toHaveBeenCalledTimes(1);

    rerender({ dep: 'b' });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(cb).toHaveBeenCalledTimes(2);
  });
});
