import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import FlipCard from './FlipCard';

vi.mock('./FlipCard.styles', () => ({
  FlipCardWrapper: ({ children, onClick }: any) => (
    <div data-testid="flip-card-wrapper" onClick={onClick}>
      {children}
    </div>
  ),
  InnerCard: ({ children, style }: any) => (
    <div data-testid="inner-card" style={style}>
      {children}
    </div>
  ),
}));

vi.mock('@/utils/hooks', () => ({
  useTimeout: (callback: () => void, delay: number) => {
    // Synchronously schedule via real timers to mimic the actual hook behavior.
    React.useEffect(() => {
      if (delay === null) return;
      const id = setTimeout(() => callback(), delay);
      return () => clearTimeout(id);
    }, [delay, callback]);
  },
}));

describe('FlipCard', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders the front child by default', () => {
    render(
      <FlipCard>
        <div data-testid="front">FRONT</div>
        <div data-testid="back">BACK</div>
      </FlipCard>,
    );

    expect(screen.getByTestId('front')).toBeInTheDocument();
    expect(screen.queryByTestId('back')).not.toBeInTheDocument();
  });

  it('throws when not given exactly two children', () => {
    // Suppress expected React error logs while asserting the throw.
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
    expect(() =>
      render(
        <FlipCard>
          <div>only one child</div>
        </FlipCard>,
      ),
    ).toThrow('FlipCard requires exactly two child components');
    spy.mockRestore();
  });

  it('flips and reveals the back face after the delayed swap', () => {
    const onClick = vi.fn();
    const onShow = vi.fn();

    render(
      <FlipCard duration={600} onClick={onClick} onShow={onShow}>
        <div data-testid="front">FRONT</div>
        <div data-testid="back">BACK</div>
      </FlipCard>,
    );

    expect(screen.getByTestId('front')).toBeInTheDocument();

    act(() => {
      fireEvent.click(screen.getByTestId('flip-card-wrapper'));
    });

    expect(onClick).toHaveBeenCalledTimes(1);
    expect(onShow).toHaveBeenCalledTimes(1);

    // Advance past duration/6 (100ms) so the activeBack timeout fires.
    act(() => {
      vi.advanceTimersByTime(150);
    });

    expect(screen.getByTestId('back')).toBeInTheDocument();
    expect(screen.queryByTestId('front')).not.toBeInTheDocument();
  });

  it('uses the default 500ms duration when none is provided', () => {
    render(
      <FlipCard>
        <div data-testid="front">FRONT</div>
        <div data-testid="back">BACK</div>
      </FlipCard>,
    );

    const inner = screen.getByTestId('inner-card');
    expect(inner).toHaveStyle({ transition: 'transform 500ms' });
  });

  it('applies the scale transform after a flip click', () => {
    render(
      <FlipCard duration={300}>
        <div data-testid="front">FRONT</div>
        <div data-testid="back">BACK</div>
      </FlipCard>,
    );

    expect(screen.getByTestId('inner-card')).not.toHaveStyle({ transform: 'scale(-1,1)' });

    act(() => {
      fireEvent.click(screen.getByTestId('flip-card-wrapper'));
    });

    expect(screen.getByTestId('inner-card')).toHaveStyle({ transform: 'scale(-1,1)' });
  });

  it('does not call onClick/onShow when omitted', () => {
    render(
      <FlipCard>
        <div data-testid="front">FRONT</div>
        <div data-testid="back">BACK</div>
      </FlipCard>,
    );

    expect(() => {
      act(() => {
        fireEvent.click(screen.getByTestId('flip-card-wrapper'));
      });
    }).not.toThrow();
  });
});
