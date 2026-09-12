import React from 'react';
import { render, screen, act } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div>{children}</div>,
  CircularProgress: ({ value }: any) => <div data-testid="progress" data-value={value} />,
  NoSsr: ({ children }: any) => <>{children}</>,
  Typography: ({ children }: any) => <span data-testid="time-typography">{children}</span>,
}));

import LoadTestTimerDialog from './load-test-timer-dialog';

describe('LoadTestTimerDialog', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns null when not open', () => {
    const { container } = render(
      <LoadTestTimerDialog open={false} t="30s" countDownComplete={vi.fn()} />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders the formatted duration in seconds when open', () => {
    render(<LoadTestTimerDialog open={true} t="30s" countDownComplete={vi.fn()} />);
    expect(screen.getByTestId('time-typography')).toHaveTextContent('30');
  });

  it('handles minutes via the m suffix', () => {
    render(<LoadTestTimerDialog open={true} t="2m" countDownComplete={vi.fn()} />);
    expect(screen.getByTestId('time-typography')).toHaveTextContent('02:00');
  });

  it('handles hours via the h suffix', () => {
    render(<LoadTestTimerDialog open={true} t="1h" countDownComplete={vi.fn()} />);
    expect(screen.getByTestId('time-typography')).toHaveTextContent('01:00:00');
  });

  it('decrements the displayed time every second', () => {
    render(<LoadTestTimerDialog open={true} t="5s" countDownComplete={vi.fn()} />);
    expect(screen.getByTestId('time-typography')).toHaveTextContent('5');
    act(() => {
      vi.advanceTimersByTime(1000);
    });
    expect(screen.getByTestId('time-typography')).toHaveTextContent('4');
  });

  it('invokes countDownComplete when the timer reaches zero', () => {
    const cb = vi.fn();
    render(<LoadTestTimerDialog open={true} t="2s" countDownComplete={cb} />);
    act(() => {
      vi.advanceTimersByTime(2000);
    });
    expect(cb).toHaveBeenCalled();
  });

  it('treats unparseable duration values as zero', () => {
    render(<LoadTestTimerDialog open={true} t="abcs" countDownComplete={vi.fn()} />);
    expect(screen.getByTestId('time-typography')).toHaveTextContent('0');
  });

  it('handles missing t prop', () => {
    render(<LoadTestTimerDialog open={true} t={undefined as any} countDownComplete={vi.fn()} />);
    expect(screen.getByTestId('time-typography')).toHaveTextContent('0');
  });

  it('resets timer when the open prop transitions to false', () => {
    const { rerender } = render(
      <LoadTestTimerDialog open={true} t="10s" countDownComplete={vi.fn()} />,
    );
    rerender(<LoadTestTimerDialog open={false} t="10s" countDownComplete={vi.fn()} />);
    expect(screen.queryByTestId('time-typography')).not.toBeInTheDocument();
  });
});
