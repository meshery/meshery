import React from 'react';
import { act, render, screen } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let lastObserver:
  | {
      callback: ResizeObserverCallback;
      disconnect: () => void;
      observe: () => void;
    }
  | undefined;

vi.mock('resize-observer-polyfill', () => ({
  default: class MockResizeObserver {
    callback: ResizeObserverCallback;
    constructor(cb: ResizeObserverCallback) {
      this.callback = cb;
      lastObserver = {
        callback: cb,
        disconnect: vi.fn(),
        observe: vi.fn(),
      };
    }
    observe() {
      lastObserver?.observe();
    }
    disconnect() {
      lastObserver?.disconnect();
    }
    unobserve() {}
  },
}));

const honeycombSpy = vi.fn();

vi.mock('./Honeycomb', () => ({
  default: (props: {
    columns: number;
    size: number;
    items: unknown[];
    renderItem: (item: unknown, index: number) => React.ReactNode;
    containerRef: React.RefObject<HTMLDivElement>;
  }) => {
    honeycombSpy(props);
    // Forward containerRef so the useEffect resize observer in
    // ResponsiveHoneycomb can attach to a real DOM node.
    return (
      <div data-testid="honeycomb-mock" data-columns={props.columns} ref={props.containerRef}>
        items={props.items.length}
      </div>
    );
  },
}));

import ResponsiveHoneycomb from './ResponsiveHoneycomb';

describe('ResponsiveHoneycomb', () => {
  beforeEach(() => {
    honeycombSpy.mockReset();
    lastObserver = undefined;
  });

  it('renders the Honeycomb child with an initial columns count derived from defaultWidth', () => {
    render(
      <ResponsiveHoneycomb
        items={['a', 'b']}
        renderItem={() => null}
        size={10}
        defaultWidth={200}
      />,
    );
    const honeycomb = screen.getByTestId('honeycomb-mock');
    // For size=10, defaultWidth=200 -> floor(200/(sqrt(3)*10)) ~ 11
    expect(Number(honeycomb.getAttribute('data-columns'))).toBeGreaterThan(0);
  });

  it('updates the column count when ResizeObserver reports a new width', () => {
    render(
      <ResponsiveHoneycomb
        items={['a', 'b']}
        renderItem={() => null}
        size={10}
        defaultWidth={200}
      />,
    );

    const initialColumns = Number(
      screen.getByTestId('honeycomb-mock').getAttribute('data-columns'),
    );

    act(() => {
      lastObserver?.callback(
        [
          {
            contentRect: { width: 800 } as DOMRectReadOnly,
          } as ResizeObserverEntry,
        ],
        {} as ResizeObserver,
      );
    });

    const updatedColumns = Number(
      screen.getByTestId('honeycomb-mock').getAttribute('data-columns'),
    );
    expect(updatedColumns).toBeGreaterThan(initialColumns);
  });

  it('ignores ResizeObserver entries with non-positive widths', () => {
    render(
      <ResponsiveHoneycomb items={['a']} renderItem={() => null} size={10} defaultWidth={200} />,
    );

    const initialColumns = screen.getByTestId('honeycomb-mock').getAttribute('data-columns');

    act(() => {
      lastObserver?.callback(
        [{ contentRect: { width: 0 } as DOMRectReadOnly } as ResizeObserverEntry],
        {} as ResizeObserver,
      );
    });

    const columnsAfter = screen.getByTestId('honeycomb-mock').getAttribute('data-columns');
    expect(columnsAfter).toBe(initialColumns);
  });
});
