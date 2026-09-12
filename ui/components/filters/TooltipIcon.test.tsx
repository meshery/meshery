import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={title}>
      {children}
    </div>
  ),
  IconButton: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
}));

import TooltipIcon from './TooltipIcon';

describe('TooltipIcon', () => {
  it('renders the tooltip with the provided title and wraps children in an IconButton', () => {
    render(
      <TooltipIcon onClick={() => {}} title="An action">
        <span data-testid="child">child</span>
      </TooltipIcon>,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'An action');
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('forwards click events to onClick', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();
    render(
      <TooltipIcon onClick={onClick} title="Click me">
        <span>X</span>
      </TooltipIcon>,
    );
    await user.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
