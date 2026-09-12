import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  GridOnIcon: (props: any) => <svg data-testid="grid-icon" {...props} />,
  TableChartIcon: (props: any) => <svg data-testid="table-icon" {...props} />,
  IconButton: ({ children, onClick, value, 'aria-label': ariaLabel, size }: any) => (
    <button
      data-testid="icon-button"
      data-size={size}
      data-value={value}
      aria-label={ariaLabel}
      onClick={onClick}
    >
      {children}
    </button>
  ),
  useTheme: () => ({ palette: { icon: { default: '#abcdef' } } }),
}));

import ViewSwitch from './ViewSwitch';

describe('ViewSwitch', () => {
  it('renders the table icon when in grid mode', () => {
    render(<ViewSwitch view="grid" changeView={vi.fn()} />);

    expect(screen.getByTestId('table-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('grid-icon')).not.toBeInTheDocument();
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'Grid View');
  });

  it('renders the grid icon when in table mode', () => {
    render(<ViewSwitch view="table" changeView={vi.fn()} />);

    expect(screen.getByTestId('grid-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('table-icon')).not.toBeInTheDocument();
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'Table View');
  });

  it('toggles from grid to table when clicked', async () => {
    const user = userEvent.setup();
    const changeView = vi.fn();
    render(<ViewSwitch view="grid" changeView={changeView} />);

    await user.click(screen.getByRole('button', { name: 'Switch View' }));
    expect(changeView).toHaveBeenCalledWith('table');
  });

  it('toggles from table to grid when clicked', async () => {
    const user = userEvent.setup();
    const changeView = vi.fn();
    render(<ViewSwitch view="table" changeView={changeView} />);

    await user.click(screen.getByRole('button', { name: 'Switch View' }));
    expect(changeView).toHaveBeenCalledWith('grid');
  });

  it('passes the current view as the icon button value', () => {
    render(<ViewSwitch view="grid" changeView={vi.fn()} />);
    expect(screen.getByTestId('icon-button')).toHaveAttribute('data-value', 'grid');
  });

  it('applies the theme icon color to the rendered icon', () => {
    render(<ViewSwitch view="grid" changeView={vi.fn()} />);
    // JSDOM normalises hex to rgb; assert by inspecting inline style.
    expect(screen.getByTestId('table-icon').getAttribute('style')).toMatch(
      /rgb\(171,\s*205,\s*239\)/,
    );
  });
});
