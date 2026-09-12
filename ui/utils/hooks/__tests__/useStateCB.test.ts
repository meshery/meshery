import { describe, expect, it, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import useStateCB from '../useStateCB';

describe('useStateCB', () => {
  it('returns the initial state', () => {
    const { result } = renderHook(() => useStateCB(42));
    expect(result.current[0]).toBe(42);
  });

  it('updates state when setState is called', () => {
    const { result } = renderHook(() => useStateCB('a'));

    act(() => {
      result.current[1]('b');
    });

    expect(result.current[0]).toBe('b');
  });

  it('does NOT fire the change-track callback on initial mount, only on subsequent updates', () => {
    const changeTrack = vi.fn();
    const { result } = renderHook(() => useStateCB(0, changeTrack));

    expect(changeTrack).not.toHaveBeenCalled();

    act(() => {
      result.current[1](1);
    });

    expect(changeTrack).toHaveBeenCalledTimes(1);
    expect(changeTrack).toHaveBeenCalledWith(1);
  });

  it('invokes the per-call callback with the new state after the render commits', () => {
    const { result } = renderHook(() => useStateCB(0));
    const after = vi.fn();

    act(() => {
      result.current[1](7, after);
    });

    expect(after).toHaveBeenCalledTimes(1);
    expect(after).toHaveBeenCalledWith(7);
  });

  it('exposes a getter that returns the latest synchronously-recorded value', () => {
    const { result } = renderHook(() => useStateCB('init'));
    const [, setState, getRef] = result.current;

    act(() => {
      setState('next');
    });

    expect(getRef()).toBe('next');
  });

  it('drops the previous per-call callback when setState is called twice in a row', () => {
    const { result } = renderHook(() => useStateCB(0));
    const first = vi.fn();
    const second = vi.fn();

    act(() => {
      result.current[1](1, first);
      result.current[1](2, second);
    });

    expect(first).not.toHaveBeenCalled();
    expect(second).toHaveBeenCalledTimes(1);
    expect(second).toHaveBeenCalledWith(2);
  });

  it('invokes the change-track callback for every change after first mount', () => {
    const changeTrack = vi.fn();
    const { result } = renderHook(() => useStateCB(0, changeTrack));

    act(() => {
      result.current[1](1);
    });
    act(() => {
      result.current[1](2);
    });
    act(() => {
      result.current[1](3);
    });

    expect(changeTrack).toHaveBeenCalledTimes(3);
    expect(changeTrack.mock.calls.map((c) => c[0])).toEqual([1, 2, 3]);
  });
});
