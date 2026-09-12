import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { render } from '@testing-library/react';
import { act } from 'react';
import React from 'react';
import { getWindowDimensions, useWindowDimensions } from '../dimension';

describe('getWindowDimensions', () => {
  let originalInnerWidth: number;
  let originalInnerHeight: number;

  beforeEach(() => {
    originalInnerWidth = window.innerWidth;
    originalInnerHeight = window.innerHeight;
  });

  afterEach(() => {
    Object.defineProperty(window, 'innerWidth', {
      configurable: true,
      writable: true,
      value: originalInnerWidth,
    });
    Object.defineProperty(window, 'innerHeight', {
      configurable: true,
      writable: true,
      value: originalInnerHeight,
    });
  });

  it('returns the current window inner dimensions', () => {
    Object.defineProperty(window, 'innerWidth', {
      configurable: true,
      writable: true,
      value: 1024,
    });
    Object.defineProperty(window, 'innerHeight', {
      configurable: true,
      writable: true,
      value: 768,
    });

    expect(getWindowDimensions()).toEqual({ width: 1024, height: 768 });
  });

  it('reflects changes to window.innerWidth/innerHeight', () => {
    Object.defineProperty(window, 'innerWidth', { configurable: true, writable: true, value: 320 });
    Object.defineProperty(window, 'innerHeight', {
      configurable: true,
      writable: true,
      value: 480,
    });

    expect(getWindowDimensions()).toEqual({ width: 320, height: 480 });
  });
});

describe('useWindowDimensions', () => {
  let originalInnerWidth: number;
  let originalInnerHeight: number;

  beforeEach(() => {
    vi.useFakeTimers();
    originalInnerWidth = window.innerWidth;
    originalInnerHeight = window.innerHeight;
    Object.defineProperty(window, 'innerWidth', { configurable: true, writable: true, value: 800 });
    Object.defineProperty(window, 'innerHeight', {
      configurable: true,
      writable: true,
      value: 600,
    });
  });

  afterEach(() => {
    vi.useRealTimers();
    Object.defineProperty(window, 'innerWidth', {
      configurable: true,
      writable: true,
      value: originalInnerWidth,
    });
    Object.defineProperty(window, 'innerHeight', {
      configurable: true,
      writable: true,
      value: originalInnerHeight,
    });
  });

  function Probe({ onRender }: { onRender: (dims: { width: number; height: number }) => void }) {
    const dims = useWindowDimensions();
    onRender(dims);
    return null;
  }

  it('returns the initial window dimensions', () => {
    const onRender = vi.fn();
    render(<Probe onRender={onRender} />);
    expect(onRender).toHaveBeenCalledWith({ width: 800, height: 600 });
  });

  it('updates the dimensions on resize (debounced by 500ms)', () => {
    const onRender = vi.fn();
    render(<Probe onRender={onRender} />);
    onRender.mockClear();

    Object.defineProperty(window, 'innerWidth', {
      configurable: true,
      writable: true,
      value: 1200,
    });
    Object.defineProperty(window, 'innerHeight', {
      configurable: true,
      writable: true,
      value: 900,
    });

    act(() => {
      window.dispatchEvent(new Event('resize'));
    });

    // Within the debounce window: no update yet
    expect(onRender).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(onRender).toHaveBeenCalledWith({ width: 1200, height: 900 });
  });

  it('cleans up the resize listener on unmount', () => {
    const removeListenerSpy = vi.spyOn(window, 'removeEventListener');
    const { unmount } = render(<Probe onRender={() => {}} />);
    unmount();
    expect(removeListenerSpy.mock.calls.some(([eventName]) => eventName === 'resize')).toBe(true);
    removeListenerSpy.mockRestore();
  });
});
