import { describe, expect, it, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import useEventListener from '../useEventListener';

describe('useEventListener', () => {
  it('subscribes to events on the supplied target', () => {
    const target = new EventTarget();
    const handler = vi.fn();
    renderHook(() => useEventListener('ping', handler, target));

    act(() => {
      target.dispatchEvent(new Event('ping'));
    });

    expect(handler).toHaveBeenCalledTimes(1);
  });

  it('defaults to window when no target is given', () => {
    const handler = vi.fn();
    renderHook(() => useEventListener('focus', handler));

    act(() => {
      window.dispatchEvent(new Event('focus'));
    });

    expect(handler).toHaveBeenCalled();
  });

  it('removes the listener on unmount', () => {
    const target = new EventTarget();
    const handler = vi.fn();
    const { unmount } = renderHook(() => useEventListener('ping', handler, target));

    unmount();

    act(() => {
      target.dispatchEvent(new Event('ping'));
    });

    expect(handler).not.toHaveBeenCalled();
  });

  it('uses the latest handler without re-subscribing', () => {
    const target = new EventTarget();
    const addSpy = vi.spyOn(target, 'addEventListener');
    const first = vi.fn();
    const second = vi.fn();
    const { rerender } = renderHook(({ h }) => useEventListener('ping', h, target), {
      initialProps: { h: first },
    });

    rerender({ h: second });

    act(() => {
      target.dispatchEvent(new Event('ping'));
    });

    expect(first).not.toHaveBeenCalled();
    expect(second).toHaveBeenCalledTimes(1);
    // Should have subscribed only once across the rerender
    expect(addSpy).toHaveBeenCalledTimes(1);
  });

  it('re-subscribes when the event name changes', () => {
    const target = new EventTarget();
    const handler = vi.fn();
    const { rerender } = renderHook(({ name }) => useEventListener(name, handler, target), {
      initialProps: { name: 'one' },
    });

    act(() => {
      target.dispatchEvent(new Event('one'));
    });
    expect(handler).toHaveBeenCalledTimes(1);

    rerender({ name: 'two' });

    // Old event should no longer fire the handler
    act(() => {
      target.dispatchEvent(new Event('one'));
    });
    expect(handler).toHaveBeenCalledTimes(1);

    act(() => {
      target.dispatchEvent(new Event('two'));
    });
    expect(handler).toHaveBeenCalledTimes(2);
  });

  it('gracefully no-ops when the target is null', () => {
    const handler = vi.fn();
    expect(() => {
      renderHook(() => useEventListener('ping', handler, null));
    }).not.toThrow();
    expect(handler).not.toHaveBeenCalled();
  });

  it('gracefully no-ops when the target lacks addEventListener', () => {
    const handler = vi.fn();
    const fakeTarget = {} as EventTarget;
    expect(() => {
      renderHook(() => useEventListener('ping', handler, fakeTarget));
    }).not.toThrow();
  });
});
