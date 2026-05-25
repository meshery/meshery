import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../style', () => ({
  HoneycombContainer: ({
    children,
    columnSize,
    columns,
    rowSize,
    className,
  }: {
    children?: React.ReactNode;
    columnSize: number;
    columns: number;
    rowSize: number;
    className?: string;
  }) => (
    <div
      data-testid="container"
      data-column-size={columnSize}
      data-columns={columns}
      data-row-size={rowSize}
      className={className}
    >
      {children}
    </div>
  ),
}));

vi.mock('./HoneycombCell', () => ({
  default: ({
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

import Honeycomb from './Honeycomb';

describe('Honeycomb', () => {
  it('renders one cell per item and passes computed sizes to the container', () => {
    render(
      <Honeycomb
        items={['a', 'b', 'c']}
        renderItem={(item, index) => (
          <span data-testid={`item-${index}`} key={index}>
            {item as string}
          </span>
        )}
        size={20}
        columns={2}
      />,
    );

    const container = screen.getByTestId('container');
    expect(container).toHaveAttribute('data-columns', '2');
    // rowSize = max(1, 20/2) = 10; columnSize = max(1, sqrt(3) * 20 / 4)
    expect(container).toHaveAttribute('data-row-size', '10');
    expect(Number(container.getAttribute('data-column-size'))).toBeGreaterThan(8);

    expect(screen.getAllByTestId('cell')).toHaveLength(3);
    expect(screen.getByTestId('item-0')).toHaveTextContent('a');
    expect(screen.getByTestId('item-2')).toHaveTextContent('c');
  });

  it('assigns row/column based on index and columns count', () => {
    // For 5 items in 2 columns, item indexes 0..4 should produce:
    // index 0 -> row 1 + 0 = 1, col 1 + 0 = 1
    // index 1 -> row 1 + 0 = 1, col 1 + 4 = 5
    // index 2 -> row 1 + 3 = 4, col 1 + 0 = 1
    render(
      <Honeycomb
        items={[0, 1, 2, 3, 4]}
        renderItem={(item, index) => <span key={index}>{item as number}</span>}
        size={10}
        columns={2}
      />,
    );

    const cells = screen.getAllByTestId('cell');
    expect(cells[0]).toHaveAttribute('data-row', '1');
    expect(cells[0]).toHaveAttribute('data-column', '1');
    expect(cells[1]).toHaveAttribute('data-row', '1');
    expect(cells[1]).toHaveAttribute('data-column', '5');
    expect(cells[2]).toHaveAttribute('data-row', '4');
    expect(cells[2]).toHaveAttribute('data-column', '1');
    expect(cells[3]).toHaveAttribute('data-row', '4');
    expect(cells[3]).toHaveAttribute('data-column', '5');
    expect(cells[4]).toHaveAttribute('data-row', '7');
    expect(cells[4]).toHaveAttribute('data-column', '1');
  });

  it('forwards className to the container', () => {
    render(
      <Honeycomb items={[]} renderItem={() => null} size={10} columns={2} className="hex-grid" />,
    );
    expect(screen.getByTestId('container')).toHaveClass('hex-grid');
  });
});
