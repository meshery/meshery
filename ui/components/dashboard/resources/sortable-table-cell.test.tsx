import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children, style }: { children?: React.ReactNode; style?: React.CSSProperties }) => (
    <div style={style}>{children}</div>
  ),
  Tooltip: ({
    title,
    children,
    placement,
  }: {
    title?: React.ReactNode;
    children?: React.ReactNode;
    placement?: string;
  }) => (
    <div data-testid="tooltip" data-title={String(title)} data-placement={placement}>
      {children}
    </div>
  ),
  Typography: ({ children, variant }: { children?: React.ReactNode; variant?: string }) => (
    <span data-variant={variant}>{children}</span>
  ),
  TableCell: ({
    children,
    onClick,
  }: {
    children?: React.ReactNode;
    onClick?: React.MouseEventHandler;
  }) => (
    <button onClick={onClick} type="button">
      {children}
    </button>
  ),
  TableSortLabel: ({ active, direction }: { active: boolean; direction?: string }) => (
    <span data-testid="sort-label" data-active={String(active)} data-direction={direction} />
  ),
}));

import { DefaultTableCell, SortableTableCell } from './sortable-table-cell';

describe('SortableTableCell', () => {
  it('renders label, sort indicator and calls onSort when clicked', async () => {
    const onSort = vi.fn();
    render(
      <SortableTableCell
        index={0}
        columnData={{ name: 'name', label: 'Name' }}
        columnMeta={{ name: 'name', direction: 'asc' }}
        onSort={onSort}
      />,
    );

    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.getByTestId('sort-label')).toHaveAttribute('data-active', 'true');
    expect(screen.getByTestId('sort-label')).toHaveAttribute('data-direction', 'asc');
    await userEvent.click(screen.getByRole('button'));
    expect(onSort).toHaveBeenCalledTimes(1);
  });

  it('defaults the sort direction to "asc" when not specified', () => {
    render(
      <SortableTableCell
        index={1}
        columnData={{ name: 'status', label: 'Status' }}
        columnMeta={{ name: 'status' }}
      />,
    );
    expect(screen.getByTestId('sort-label')).toHaveAttribute('data-direction', 'asc');
  });

  it('marks the sort label inactive when columnMeta name differs from columnData', () => {
    render(
      <SortableTableCell
        index={2}
        columnData={{ name: 'created', label: 'Created' }}
        columnMeta={{ name: 'other' }}
      />,
    );
    expect(screen.getByTestId('sort-label')).toHaveAttribute('data-active', 'false');
  });

  it('renders a tooltip wrapper around the icon when provided', () => {
    render(
      <SortableTableCell
        index={3}
        columnData={{ name: 'a', label: 'A' }}
        columnMeta={{ name: 'a' }}
        icon={<i data-testid="icon">i</i>}
        tooltip="hint"
      />,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'hint');
    expect(screen.getByTestId('icon')).toBeInTheDocument();
  });

  it('omits the tooltip wrapper when no icon is provided', () => {
    render(
      <SortableTableCell
        index={4}
        columnData={{ name: 'a', label: 'A' }}
        columnMeta={{ name: 'a' }}
      />,
    );
    expect(screen.queryByTestId('tooltip')).not.toBeInTheDocument();
  });
});

describe('DefaultTableCell', () => {
  it('renders just the label and no tooltip when no icon is provided', () => {
    render(<DefaultTableCell columnData={{ name: 'name', label: 'Name' }} />);
    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.queryByTestId('tooltip')).not.toBeInTheDocument();
  });

  it('renders an icon inside a Tooltip when icon is provided', () => {
    render(
      <DefaultTableCell
        columnData={{ name: 'a', label: 'A' }}
        icon={<i data-testid="icon">i</i>}
        tooltip="info"
      />,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'info');
    expect(screen.getByTestId('icon')).toBeInTheDocument();
  });

  it('defaults the tooltip title to "" when no tooltip prop is provided alongside an icon', () => {
    render(
      <DefaultTableCell columnData={{ name: 'a', label: 'A' }} icon={<i data-testid="icon" />} />,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', '');
  });
});
