import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import debounceWidthProvider from './debounceWidthProvider';

class MockResizeObserver {
  static instances: MockResizeObserver[] = [];
  callback: ResizeObserverCallback;
  observed: Element[] = [];

  constructor(cb: ResizeObserverCallback) {
    this.callback = cb;
    MockResizeObserver.instances.push(this);
  }
  observe(target: Element) {
    this.observed.push(target);
  }
  disconnect() {
    this.observed = [];
  }
  unobserve() {}
  trigger(width: number) {
    this.callback(
      [
        {
          contentRect: { width } as DOMRectReadOnly,
        } as ResizeObserverEntry,
      ],
      this as unknown as ResizeObserver,
    );
  }
}

describe('debounceWidthProvider', () => {
  const originalRO = global.ResizeObserver;
  const originalRAF = global.requestAnimationFrame;
  const originalCAF = global.cancelAnimationFrame;

  beforeEach(() => {
    MockResizeObserver.instances = [];
    (global as any).ResizeObserver = MockResizeObserver;
    // Synchronous rAF for simpler assertions.
    (global as any).requestAnimationFrame = (cb: FrameRequestCallback) => {
      cb(0);
      return 1;
    };
    (global as any).cancelAnimationFrame = vi.fn();
  });

  afterEach(() => {
    (global as any).ResizeObserver = originalRO;
    (global as any).requestAnimationFrame = originalRAF;
    (global as any).cancelAnimationFrame = originalCAF;
  });

  it('renders the wrapped component immediately at the default width', () => {
    const Inner = vi.fn(({ width }: { width: number }) => <div data-testid="inner">w={width}</div>);
    const Wrapped = debounceWidthProvider(Inner);

    render(<Wrapped />);
    expect(screen.getByTestId('inner')).toHaveTextContent('w=1280');
  });

  it('updates the width immediately on the first resize observation', async () => {
    const Inner = ({ width }: { width: number }) => <div data-testid="inner">w={width}</div>;
    const Wrapped = debounceWidthProvider(Inner);

    render(<Wrapped />);
    expect(MockResizeObserver.instances).toHaveLength(1);
    MockResizeObserver.instances[0].trigger(840);

    await waitFor(() => {
      expect(screen.getByTestId('inner')).toHaveTextContent('w=840');
    });
  });

  it('waits to render until measurement when measureBeforeMount is set', async () => {
    const Inner = ({ width }: { width: number }) => <div data-testid="inner">w={width}</div>;
    const Wrapped = debounceWidthProvider(Inner);

    render(<Wrapped measureBeforeMount={true} />);

    expect(screen.queryByTestId('inner')).not.toBeInTheDocument();

    MockResizeObserver.instances[0].trigger(640);

    await waitFor(() => {
      expect(screen.getByTestId('inner')).toHaveTextContent('w=640');
    });
  });

  it('sets a displayName based on the inner component', () => {
    const Named = ({ width }: { width: number }) => <div>w={width}</div>;
    Named.displayName = 'CustomInner';
    const Wrapped = debounceWidthProvider(Named);
    expect(Wrapped.displayName).toBe('DebouncedWidthProvider(CustomInner)');
  });

  it('falls back to "Component" when the wrapped element has no name', () => {
    const Wrapped = debounceWidthProvider(() => <div />);
    expect(Wrapped.displayName).toBe('DebouncedWidthProvider(Component)');
  });

  it('forwards className and style to the wrapper div', () => {
    const Inner = ({ width }: { width: number }) => <div>w={width}</div>;
    const Wrapped = debounceWidthProvider(Inner);
    const { container } = render(<Wrapped className="cool" style={{ padding: '12px' }} />);
    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper).toHaveClass('cool');
    expect(wrapper.style.padding).toBe('12px');
  });
});
