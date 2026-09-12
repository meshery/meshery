import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children, style }: any) => (
    <div data-testid="grid" data-style={JSON.stringify(style || {})}>
      {children}
    </div>
  ),
  Typography: ({ children, style }: any) => (
    <span data-style={JSON.stringify(style || {})}>{children}</span>
  ),
}));

vi.mock('./curvedArrowIcon', () => ({
  default: () => <svg data-testid="curved-arrow" />,
}));

import EmptyState from './index';

describe('EmptyState', () => {
  it('renders message, pointer label and icon when supplied', () => {
    render(
      <EmptyState
        icon={<svg data-testid="state-icon" />}
        message="Nothing here yet"
        pointerLabel="Click here"
      />,
    );

    expect(screen.getByTestId('curved-arrow')).toBeInTheDocument();
    expect(screen.getByTestId('state-icon')).toBeInTheDocument();
    expect(screen.getByText('Nothing here yet')).toBeInTheDocument();
    expect(screen.getByText('Click here')).toBeInTheDocument();
  });

  it('still renders structural grids when only message is provided', () => {
    render(<EmptyState message="Empty" />);
    expect(screen.getAllByTestId('grid').length).toBe(2);
    expect(screen.getByText('Empty')).toBeInTheDocument();
  });
});
