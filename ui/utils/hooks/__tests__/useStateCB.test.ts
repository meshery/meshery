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

  it('always invokes the latest changeTrackCB after rerenders', () => {
    const firstTrack = vi.fn();
    const secondTrack = vi.fn();
    const { result, rerender } = renderHook(({ cb }) => useStateCB(0, cb), {
      initialProps: { cb: firstTrack },
    });

    rerender({ cb: secondTrack });

    act(() => {
      result.current[1](10);
    });

    expect(firstTrack).not.toHaveBeenCalled();
    expect(secondTrack).toHaveBeenCalledTimes(1);
    expect(secondTrack).toHaveBeenCalledWith(10);
  });

  it('supports lazy initializer functions and keeps getStateRefValue accurate', () => {
    const init = vi.fn(() => 99);
    const { result } = renderHook(() => useStateCB(init));

    expect(init).toHaveBeenCalledTimes(1);
    expect(result.current[0]).toBe(99);
    expect(result.current[2]()).toBe(99);
  });

  it('supports functional state updaters and keeps stateRef in sync across sequential updates', () => {
    const changeTrack = vi.fn();
    const after = vi.fn();
    const { result } = renderHook(() => useStateCB(0, changeTrack));

    act(() => {
      result.current[1]((prev) => prev + 1);
      result.current[1]((prev) => prev + 2, after);
      expect(result.current[2]()).toBe(3);
    });

    expect(result.current[0]).toBe(3);
    expect(result.current[2]()).toBe(3);
    expect(after).toHaveBeenCalledTimes(1);
    expect(after).toHaveBeenCalledWith(3);
    expect(changeTrack).toHaveBeenCalledTimes(1);
    expect(changeTrack).toHaveBeenCalledWith(3);
  });

  it('ignores non-callable callbacks and preserves setState and getter identities', () => {
    const invalidTrack = 'not-a-function' as unknown as (value: number) => void;
    const { result, rerender } = renderHook(() => useStateCB(0, invalidTrack));
    const [, initialSetState, initialGetRef] = result.current;

    expect(() => {
      act(() => {
        initialSetState(5, 'invalid' as unknown as (value: number) => void);
      });
    }).not.toThrow();

    rerender();
    expect(result.current[0]).toBe(5);
    expect(result.current[1]).toBe(initialSetState);
    expect(result.current[2]).toBe(initialGetRef);
  });

  it('preserves function-valued state when resolved from an updater', () => {
    const initialFn = () => 1;
    const nextFn = () => 2;
    const after = vi.fn();
    const { result } = renderHook(() => useStateCB<() => number>(() => initialFn));

    act(() => {
      result.current[1](() => nextFn, after);
    });

    expect(result.current[0]).toBe(nextFn);
    expect(result.current[0]()).toBe(2);
    expect(result.current[2]()).toBe(nextFn);
    expect(after).toHaveBeenCalledWith(nextFn);
  });
});
