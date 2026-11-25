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

  it('allows a zero delay without falling back to default', () => {
    const fn = vi.fn();
    const { result } = renderHook(() => useDebouncedCallback(fn, 0));

    act(() => {
      result.current('instant');
    });
    expect(fn).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(0);
    });
    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('instant');
  });

  it('falls back to default 300ms when delay is null, negative, NaN, Infinity, or non-numeric', () => {
    const numericOrNullBadDelays: Array<number | null> = [
      null,
      -100,
      Number.NaN,
      Number.POSITIVE_INFINITY,
    ];

    for (const badDelay of numericOrNullBadDelays) {
      const fn = vi.fn();
      const { result, unmount } = renderHook(() => useDebouncedCallback(fn, badDelay));

      act(() => {
        result.current('payload');
      });

      act(() => {
        vi.advanceTimersByTime(299);
      });
      expect(fn).not.toHaveBeenCalled();

      act(() => {
        vi.advanceTimersByTime(1);
      });
      expect(fn).toHaveBeenCalledTimes(1);
      expect(fn).toHaveBeenCalledWith('payload');
      unmount();
    }

    const fnForStringDelay = vi.fn();
    const { result: strResult, unmount: unmountStr } = renderHook(() =>
      // @ts-expect-error testing invalid string delay input
      useDebouncedCallback(fnForStringDelay, '500'),
    );

    act(() => {
      strResult.current('payload');
      vi.advanceTimersByTime(299);
    });
    expect(fnForStringDelay).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(fnForStringDelay).toHaveBeenCalledTimes(1);
    expect(fnForStringDelay).toHaveBeenCalledWith('payload');
    unmountStr();
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

  it('invokes the latest callback without changing debounced function identity', () => {
    const first = vi.fn();
    const second = vi.fn();
    const { result, rerender } = renderHook(({ fn }) => useDebouncedCallback(fn, 100), {
      initialProps: { fn: first },
    });

    const initialDebounced = result.current;

    act(() => {
      result.current('latest');
    });

    rerender({ fn: second });
    expect(result.current).toBe(initialDebounced);

    act(() => {
      vi.advanceTimersByTime(100);
    });

    expect(first).not.toHaveBeenCalled();
    expect(second).toHaveBeenCalledTimes(1);
    expect(second).toHaveBeenCalledWith('latest');
  });

  it('does not throw when callback is null, undefined, or non-callable', () => {
    const { result: nullRes } = renderHook(() => useDebouncedCallback(null, 100));
    const { result: undefRes } = renderHook(() => useDebouncedCallback(undefined, 100));
    // @ts-expect-error testing non-callable input
    const { result: badRes } = renderHook(() => useDebouncedCallback('not-a-fn', 100));

    expect(() => {
      act(() => {
        nullRes.current();
        undefRes.current();
        badRes.current();
        vi.advanceTimersByTime(200);
      });
    }).not.toThrow();
  });

  it('cancels pending invocation when callback transitions to null before timer fires', () => {
    const fn = vi.fn();
    const { result, rerender } = renderHook(
      ({ cb }: { cb: (() => void) | null }) => useDebouncedCallback(cb, 100),
      { initialProps: { cb: fn as (() => void) | null } },
    );

    act(() => {
      result.current();
    });

    rerender({ cb: null });

    act(() => {
      vi.advanceTimersByTime(200);
    });

    expect(fn).not.toHaveBeenCalled();
  });

  it('cancels pending invocation when delay changes before timer fires', () => {
    const fn = vi.fn();
    const { result, rerender } = renderHook(({ delay }) => useDebouncedCallback(fn, delay), {
      initialProps: { delay: 100 },
    });

    act(() => {
      result.current();
    });

    rerender({ delay: 500 });

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(fn).not.toHaveBeenCalled();
  });
});
