import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../style', () => ({
  HoneycombCell: ({
    children,
    row,
    column,
  }: {
    children?: React.ReactNode;
    row: number;
    column: number;
  }) => (
    <div data-testid="cell" data-row={row} data-column={column}>
      {children}
    </div>
  ),
}));

import HoneycombCell from './HoneycombCell';

describe('HoneycombCell', () => {
  it('forwards row/column when finite values are provided', () => {
    render(
      <HoneycombCell row={5} column={3}>
        text
      </HoneycombCell>,
    );
    const cell = screen.getByTestId('cell');
    expect(cell).toHaveAttribute('data-row', '5');
    expect(cell).toHaveAttribute('data-column', '3');
    expect(cell).toHaveTextContent('text');
  });

  it('defaults row/column to 1 when non-finite values are passed', () => {
    render(<HoneycombCell row={NaN} column={Infinity} />);
    const cell = screen.getByTestId('cell');
    expect(cell).toHaveAttribute('data-row', '1');
    expect(cell).toHaveAttribute('data-column', '1');
  });

  it('defaults row/column to 1 when neither prop is provided', () => {
    render(<HoneycombCell />);
    const cell = screen.getByTestId('cell');
    expect(cell).toHaveAttribute('data-row', '1');
    expect(cell).toHaveAttribute('data-column', '1');
  });
});
