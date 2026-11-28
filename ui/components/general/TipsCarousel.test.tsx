import React from 'react';
import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import TipsCarousel from './TipsCarousel';

let intervalCallback: (() => void) | null = null;

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Typography: ({ children, ...props }: any) => <p {...props}>{children}</p>,
  styled: (Component: any) => (_styles: any) => {
    const StyledComponent = ({ children, ...props }: any) =>
      React.createElement(Component, props, children);

    StyledComponent.displayName = 'StyledComponentMock';
    return StyledComponent;
  },
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

  it('renders one accessible button per tip', () => {
    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);

    const indicators = screen.getAllByRole('button');

    expect(indicators).toHaveLength(3);
    expect(screen.getByRole('button', { name: 'Show tip 1 of 3' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Show tip 2 of 3' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Show tip 3 of 3' })).toBeInTheDocument();
  });

  it('marks only the active indicator with aria-current', () => {
    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);

    expect(screen.getByRole('button', { name: 'Show tip 1 of 3' })).toHaveAttribute(
      'aria-current',
      'step',
    );
    expect(screen.getByRole('button', { name: 'Show tip 2 of 3' })).not.toHaveAttribute(
      'aria-current',
    );
    expect(screen.getByRole('button', { name: 'Show tip 3 of 3' })).not.toHaveAttribute(
      'aria-current',
    );
  });

  it('allows indicators to receive keyboard focus using Tab', async () => {
    const user = userEvent.setup();

    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);

    await user.tab();

    expect(screen.getByRole('button', { name: 'Show tip 1 of 3' })).toHaveFocus();

    await user.tab();

    expect(screen.getByRole('button', { name: 'Show tip 2 of 3' })).toHaveFocus();
  });

  it('activates an indicator with Enter', async () => {
    const user = userEvent.setup();

    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);

    const secondIndicator = screen.getByRole('button', { name: 'Show tip 2 of 3' });

    secondIndicator.focus();
    await user.keyboard('{Enter}');

    expect(screen.getByText('Tip B')).toBeInTheDocument();
    expect(secondIndicator).toHaveAttribute('aria-current', 'step');
    expect(screen.getByRole('button', { name: 'Show tip 1 of 3' })).not.toHaveAttribute(
      'aria-current',
    );
  });

  it('activates an indicator with Space', async () => {
    const user = userEvent.setup();

    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);

    const thirdIndicator = screen.getByRole('button', { name: 'Show tip 3 of 3' });

    thirdIndicator.focus();
    await user.keyboard(' ');

    expect(screen.getByText('Tip C')).toBeInTheDocument();
    expect(thirdIndicator).toHaveAttribute('aria-current', 'step');
    expect(screen.getByRole('button', { name: 'Show tip 1 of 3' })).not.toHaveAttribute(
      'aria-current',
    );
  });

  it('changes active tip when an indicator is clicked', async () => {
    const user = userEvent.setup();

    render(<TipsCarousel tips={['Tip A', 'Tip B', 'Tip C']} />);

    await user.click(screen.getByRole('button', { name: 'Show tip 3 of 3' }));

    expect(screen.getByText('Tip C')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Show tip 3 of 3' })).toHaveAttribute(
      'aria-current',
      'step',
    );
  });

  it('advances the active tip when the interval fires and wraps to the first tip', () => {
    render(<TipsCarousel tips={['Tip A', 'Tip B']} />);

    expect(screen.getByText('Tip A')).toBeInTheDocument();
    expect(typeof intervalCallback).toBe('function');

    act(() => {
      intervalCallback?.();
    });

    expect(screen.getByText('Tip B')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Show tip 2 of 2' })).toHaveAttribute(
      'aria-current',
      'step',
    );

    act(() => {
      intervalCallback?.();
    });

    expect(screen.getByText('Tip A')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Show tip 1 of 2' })).toHaveAttribute(
      'aria-current',
      'step',
    );
  });
});
