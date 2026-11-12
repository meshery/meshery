import React from 'react';
import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import TipsCarousel from './TipsCarousel';

let intervalCallback: (() => void) | null = null;

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Typography: ({ children, ...props }: any) => <p {...props}>{children}</p>,
  useTheme: () => ({
    palette: {
      background: {
        brand: { default: 'brand-color' },
        constant: { white: '#fff' },
      },
    },
  }),
}));

vi.mock('../../assets/icons/Tipsicon', () => ({
  default: () => <svg data-testid="tips-icon" />,
}));

vi.mock('@/utils/hooks', () => ({
  useInterval: (callback: () => void) => {
    intervalCallback = callback;
  },
}));

describe('TipsCarousel', () => {
  beforeEach(() => {
    intervalCallback = null;
  });

  afterEach(() => {
    intervalCallback = null;
  });

  it('renders the tips icon and the first tip by default', () => {
    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);
    expect(screen.getByTestId('tips-icon')).toBeInTheDocument();
    expect(screen.getByText('Tip A')).toBeInTheDocument();
  });

  it('renders one bullet per tip', () => {
    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);
    const bullets = screen.getAllByText('•');
    expect(bullets).toHaveLength(3);
  });

  it('changes active tip when a bullet is clicked', async () => {
    const user = userEvent.setup();
    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);
    const bullets = screen.getAllByText('•');
    await user.click(bullets[2]);
    expect(screen.getByText('Tip C')).toBeInTheDocument();
  });

  it('advances the active tip when the interval fires and wraps to the first tip', () => {
    render(<TipsCarousel tips={['Tip A', 'Tip B']} />);
    expect(screen.getByText('Tip A')).toBeInTheDocument();
    expect(typeof intervalCallback).toBe('function');
    act(() => {
      intervalCallback?.();
    });
    expect(screen.getByText('Tip B')).toBeInTheDocument();
    act(() => {
      intervalCallback?.();
    });
    expect(screen.getByText('Tip A')).toBeInTheDocument();
  });
});
