import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import useDebouncedCallback from '../useDebounce';

describe('useDebouncedCallback', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('delays callback invocation by the default 300ms', () => {
    const fn = vi.fn();
    const { result } = renderHook(() => useDebouncedCallback(fn));

    act(() => {
      result.current();
    });
    expect(fn).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(299);
    });
    expect(fn).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(fn).toHaveBeenCalledTimes(1);
  });

  it('honours a custom delay', () => {
    const fn = vi.fn();
    const { result } = renderHook(() => useDebouncedCallback(fn, 100));

    act(() => {
      result.current('arg');
    });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(fn).toHaveBeenCalledWith('arg');
  });

  it('resets the timer when called rapidly and only fires once with the latest args', () => {
    const fn = vi.fn();
    const { result } = renderHook(() => useDebouncedCallback(fn, 200));

    act(() => {
      result.current('a');
      vi.advanceTimersByTime(100);
      result.current('b');
      vi.advanceTimersByTime(100);
      result.current('c');
      vi.advanceTimersByTime(200);
    });

    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('c');
  });

  it('forwards all arguments', () => {
    const fn = vi.fn();
    const { result } = renderHook(() => useDebouncedCallback(fn, 50));

    act(() => {
      result.current(1, 'two', { three: 3 });
      vi.advanceTimersByTime(50);
    });

    expect(fn).toHaveBeenCalledWith(1, 'two', { three: 3 });
  });

  it('cancels the pending callback on unmount', () => {
    const fn = vi.fn();
    const { result, unmount } = renderHook(() => useDebouncedCallback(fn, 100));

    act(() => {
      result.current();
    });

    unmount();

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(fn).not.toHaveBeenCalled();
  });

  it('rebuilds the debounced function when the delay changes', () => {
    const fn = vi.fn();
    const { result, rerender } = renderHook(({ delay }) => useDebouncedCallback(fn, delay), {
      initialProps: { delay: 100 },
    });

    const first = result.current;
    rerender({ delay: 500 });
    expect(result.current).not.toBe(first);
  });
});
