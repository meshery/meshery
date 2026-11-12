import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../style', () => ({
  HexagonWrapper: ({
    children,
    className,
    style,
    onClick,
  }: {
    children?: React.ReactNode;
    className?: string;
    style?: React.CSSProperties;
    onClick?: React.MouseEventHandler<HTMLDivElement>;
  }) => (
    <div
      data-testid="hexagon-wrapper"
      className={className}
      style={style}
      onClick={onClick}
      role="button"
    >
      {children}
    </div>
  ),
}));

import Hexagon from './Hexagon';

describe('Hexagon', () => {
  it('forwards className, style, children, and onClick to the wrapper', async () => {
    const onClick = vi.fn();
    render(
      <Hexagon className="custom" style={{ color: 'red' }} onClick={onClick}>
        inside
      </Hexagon>,
    );
    const wrapper = screen.getByTestId('hexagon-wrapper');
    expect(wrapper).toHaveTextContent('inside');
    expect(wrapper).toHaveClass('custom');
    expect(wrapper).toHaveStyle({ color: 'rgb(255, 0, 0)' });
    await userEvent.click(wrapper);
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('renders without optional props', () => {
    render(<Hexagon />);
    expect(screen.getByTestId('hexagon-wrapper')).toBeInTheDocument();
  });
});
