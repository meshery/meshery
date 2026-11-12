import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useIsTextTruncated } from '../useIsTextTruncated';

describe('useIsTextTruncated', () => {
  let originalResizeObserver: typeof globalThis.ResizeObserver;
  let observerCallback: ResizeObserverCallback | null = null;
  let observeMock: ReturnType<typeof vi.fn>;
  let disconnectMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    observerCallback = null;
    observeMock = vi.fn();
    disconnectMock = vi.fn();

    originalResizeObserver = globalThis.ResizeObserver;
    globalThis.ResizeObserver = class MockResizeObserver {
      constructor(callback: ResizeObserverCallback) {
        observerCallback = callback;
      }
      observe = observeMock;
      unobserve = vi.fn();
      disconnect = disconnectMock;
    } as unknown as typeof ResizeObserver;
  });

  afterEach(() => {
    globalThis.ResizeObserver = originalResizeObserver;
  });

  const createMockElement = (scrollWidth: number, clientWidth: number): HTMLElement => {
    const el = document.createElement('div');
    Object.defineProperty(el, 'scrollWidth', {
      configurable: true,
      get: () => scrollWidth,
    });
    Object.defineProperty(el, 'clientWidth', {
      configurable: true,
      get: () => clientWidth,
    });
    return el;
  };

  it('returns false when ref.current is null', () => {
    const ref = { current: null };
    const { result } = renderHook(() => useIsTextTruncated(ref));

    expect(result.current).toBe(false);
    expect(observeMock).not.toHaveBeenCalled();
  });

  it('returns false when text is not overflowing (scrollWidth <= clientWidth)', () => {
    const element = createMockElement(100, 100);
    const ref = { current: element };

    const { result } = renderHook(() => useIsTextTruncated(ref));

    expect(result.current).toBe(false);
    expect(observeMock).toHaveBeenCalledWith(element);
  });

  it('returns true when text is overflowing (scrollWidth > clientWidth)', () => {
    const element = createMockElement(150, 100);
    const ref = { current: element };

    const { result } = renderHook(() => useIsTextTruncated(ref));

    expect(result.current).toBe(true);
    expect(observeMock).toHaveBeenCalledWith(element);
  });

  it('updates state to true when element resizes and overflows', () => {
    let scrollWidth = 100;
    let clientWidth = 100;
    const element = document.createElement('div');
    Object.defineProperty(element, 'scrollWidth', {
      configurable: true,
      get: () => scrollWidth,
    });
    Object.defineProperty(element, 'clientWidth', {
      configurable: true,
      get: () => clientWidth,
    });
    const ref = { current: element };

    const { result } = renderHook(() => useIsTextTruncated(ref));
    expect(result.current).toBe(false);

    // Simulate element resize making text overflow
    scrollWidth = 200;
    clientWidth = 100;

    act(() => {
      observerCallback?.([], {} as ResizeObserver);
    });

    expect(result.current).toBe(true);
  });

  it('updates state to false when element resizes and no longer overflows', () => {
    let scrollWidth = 200;
    let clientWidth = 100;
    const element = document.createElement('div');
    Object.defineProperty(element, 'scrollWidth', {
      configurable: true,
      get: () => scrollWidth,
    });
    Object.defineProperty(element, 'clientWidth', {
      configurable: true,
      get: () => clientWidth,
    });
    const ref = { current: element };

    const { result } = renderHook(() => useIsTextTruncated(ref));
    expect(result.current).toBe(true);

    // Simulate drawer expanding / container widening
    scrollWidth = 200;
    clientWidth = 300;

    act(() => {
      observerCallback?.([], {} as ResizeObserver);
    });

    expect(result.current).toBe(false);
  });

  it('disconnects ResizeObserver on unmount', () => {
    const element = createMockElement(100, 100);
    const ref = { current: element };

    const { unmount } = renderHook(() => useIsTextTruncated(ref));

    expect(disconnectMock).not.toHaveBeenCalled();
    unmount();
    expect(disconnectMock).toHaveBeenCalledTimes(1);
  });

  it('re-evaluates truncation when dependencies change', () => {
    let scrollWidth = 100;
    let clientWidth = 100;
    const element = document.createElement('div');
    Object.defineProperty(element, 'scrollWidth', {
      configurable: true,
      get: () => scrollWidth,
    });
    Object.defineProperty(element, 'clientWidth', {
      configurable: true,
      get: () => clientWidth,
    });
    const ref = { current: element };

    const { result, rerender } = renderHook(
      ({ open, text }) => useIsTextTruncated(ref, [open, text]),
      { initialProps: { open: false, text: 'Short' } },
    );

    expect(result.current).toBe(false);

    // Change text and simulate larger content
    scrollWidth = 250;
    rerender({ open: true, text: 'Very Long Workspace Name That Overflows' });

    expect(result.current).toBe(true);
  });
});
