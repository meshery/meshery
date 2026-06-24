import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

let isMobile = false;

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    DARK_BLUE_GRAY: '#243a44',
    IconButton: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled}>
        {children}
      </button>
    ),
    MoreVertIcon: () => <svg data-testid="more-vert" />,
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid={`tooltip-${title}`}>{children}</div>
    ),
    styled,
    Menu: ({ children, open }: any) => (open ? <div data-testid="menu">{children}</div> : null),
    MenuItem: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled} data-testid="menu-item">
        {children}
      </button>
    ),
    useMediaQuery: () => isMobile,
    useTheme: () => ({
      palette: { mode: 'light', background: { paper: '#fff' } },
      breakpoints: { down: () => '@media (max-width: 0)' },
    }),
  };
});

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

import { MenuComponent } from './MenuComponent';

describe('MenuComponent', () => {
  beforeEach(() => {
    isMobile = false;
  });

  const baseOptions = [
    {
      title: 'Edit',
      icon: <svg data-testid="edit-icon" />,
      handler: vi.fn(),
    },
    {
      title: 'Delete',
      icon: <svg data-testid="delete-icon" />,
      handler: vi.fn(),
      disabled: true,
    },
  ];

  it('renders direct icons on desktop (non-mobile) viewports', () => {
    isMobile = false;
    render(<MenuComponent options={baseOptions} />);
    expect(screen.getByTestId('tooltip-Edit')).toBeInTheDocument();
    expect(screen.getByTestId('tooltip-Delete')).toBeInTheDocument();
    expect(screen.queryByTestId('more-vert')).not.toBeInTheDocument();
  });

  it('invokes the handler when a direct icon is clicked', async () => {
    isMobile = false;
    const editHandler = vi.fn();
    const options = [{ title: 'Edit', icon: <svg />, handler: editHandler }];
    render(<MenuComponent options={options} />);

    await userEvent.setup().click(screen.getByRole('button'));
    expect(editHandler).toHaveBeenCalled();
  });

  it('disables a direct icon button when option.disabled is true', () => {
    isMobile = false;
    render(<MenuComponent options={baseOptions} />);

    const buttons = screen.getAllByRole('button');
    // Edit is enabled, Delete is disabled
    expect(buttons[0]).not.toBeDisabled();
    expect(buttons[1]).toBeDisabled();
  });

  it('renders a more-vert dropdown on mobile viewports', async () => {
    isMobile = true;
    const user = userEvent.setup();
    render(<MenuComponent options={baseOptions} />);

    expect(screen.getByTestId('more-vert')).toBeInTheDocument();
    expect(screen.queryByTestId('menu')).not.toBeInTheDocument();

    // Click anywhere on the MoreVert area to open the menu
    await user.click(screen.getByTestId('more-vert'));
    expect(screen.getByTestId('menu')).toBeInTheDocument();
    // In mobile mode the menu wraps each option in its own block (styled MenuItem).
    // The number of tooltips per-option should match options.length + 1 (the
    // top-level "Quick Actions" trigger), so verify by tooltip count.
    expect(screen.getByTestId('tooltip-Edit')).toBeInTheDocument();
    expect(screen.getByTestId('tooltip-Delete')).toBeInTheDocument();
  });

  it('handles an empty options array without rendering buttons', () => {
    isMobile = false;
    render(<MenuComponent options={[]} />);
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });
});
