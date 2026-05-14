import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { DefaultTableCell, SortableTableCell } from './index';

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children }) => <div>{children}</div>,
  Typography: ({ children }) => <span>{children}</span>,
  TableCell: ({ children, onClick }) => (
    <button onClick={onClick} type="button">
      {children}
    </button>
  ),
  TableSortLabel: ({ active, direction }) => (
    <span data-testid="sort-label" data-active={String(active)} data-direction={direction} />
  ),
}));

vi.mock('@/components/meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ title, children }) => (
    <div data-testid="custom-text-tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

describe('SortableTableCell', () => {
  it('renders the label, tooltip icon, and triggers sorting on click', async () => {
    const user = userEvent.setup();
    const onSort = vi.fn();

    render(
      <SortableTableCell
        index={0}
        columnData={{ name: 'name', label: 'Name' }}
        columnMeta={{ name: 'name', direction: 'desc' }}
        onSort={onSort}
        icon={<span>i</span>}
        tooltip="Sort by name"
      />,
    );

    await user.click(screen.getByRole('button'));

    expect(onSort).toHaveBeenCalledTimes(1);
    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.getByTestId('custom-text-tooltip')).toHaveAttribute('data-title', 'Sort by name');
    expect(screen.getByTestId('sort-label')).toHaveAttribute('data-active', 'true');
    expect(screen.getByTestId('sort-label')).toHaveAttribute('data-direction', 'desc');
  });
});

describe('DefaultTableCell', () => {
  it('renders without a tooltip wrapper when no icon is provided', () => {
    render(<DefaultTableCell columnData={{ label: 'Status' }} />);

    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.queryByTestId('custom-text-tooltip')).not.toBeInTheDocument();
  });
});
