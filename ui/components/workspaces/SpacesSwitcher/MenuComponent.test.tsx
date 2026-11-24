import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Independent flags so tablet (overflow Menu) is not collapsed into mobile (BottomSheet).
let isOverflow = false;
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
    useMediaQuery: (query: string) => {
      if (query === '(max-width:xl)') return isOverflow;
      if (query === '(max-width:sm)') return isMobile;
      return false;
    },
    useTheme: () => ({
      palette: { mode: 'light', background: { paper: '#fff' } },
      breakpoints: { down: (key: string) => `(max-width:${key})` },
    }),
    BottomSheet: ({ children, open }: any) =>
      open ? <div data-testid="bottom-sheet">{children}</div> : null,
    MenuList: ({ children }: any) => <div data-testid="menu-list">{children}</div>,
    ListItemIcon: ({ children }: any) => <span>{children}</span>,
    Typography: ({ children }: any) => <span>{children}</span>,
  };
});

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

import { MenuComponent } from './MenuComponent';

describe('MenuComponent', () => {
  beforeEach(() => {
    isOverflow = false;
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

  it('renders an overflow Menu on tablet viewports (collapsed, not mobile)', async () => {
    isOverflow = true;
    isMobile = false;
    const user = userEvent.setup();
    render(<MenuComponent options={baseOptions} />);

    expect(screen.getByTestId('more-vert')).toBeInTheDocument();
    expect(screen.queryByTestId('menu')).not.toBeInTheDocument();
    expect(screen.queryByTestId('bottom-sheet')).not.toBeInTheDocument();

    await user.click(screen.getByTestId('more-vert'));
    expect(screen.getByTestId('menu')).toBeInTheDocument();
    expect(screen.queryByTestId('bottom-sheet')).not.toBeInTheDocument();
  });

  it('renders a bottom sheet on mobile viewports', async () => {
    isOverflow = true;
    isMobile = true;
    const user = userEvent.setup();
    render(<MenuComponent options={baseOptions} />);

    expect(screen.getByTestId('more-vert')).toBeInTheDocument();
    expect(screen.queryByTestId('bottom-sheet')).not.toBeInTheDocument();

    await user.click(screen.getByTestId('more-vert'));
    expect(screen.getByTestId('bottom-sheet')).toBeInTheDocument();
    expect(screen.queryByTestId('menu')).not.toBeInTheDocument();
    expect(screen.getByTestId('menu-list')).toBeInTheDocument();
  });

  it('handles an empty options array without rendering buttons', () => {
    isMobile = false;
    render(<MenuComponent options={[]} />);
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });
});
