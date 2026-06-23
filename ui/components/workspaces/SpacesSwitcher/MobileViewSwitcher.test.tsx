import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/assets/icons/DashboardSwitcherIcon', () => ({
  default: () => <svg data-testid="dashboard-icon" />,
}));

vi.mock('@/assets/icons/OrgOutlinedIcon', () => ({
  default: () => <svg data-testid="org-icon" />,
}));

vi.mock('css/icons.styles', () => ({
  iconLarge: {},
  iconXLarge: {},
}));

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Box: ({ children }: any) => <div>{children}</div>,
    Button: ({ children, onClick, ...rest }: any) => (
      <button onClick={onClick} {...rest}>
        {children}
      </button>
    ),
    ClickAwayListener: ({ children }: any) => <div>{children}</div>,
    Grid2: ({ children }: any) => <div>{children}</div>,
    Slide: ({ children, in: open }: any) =>
      open ? <div data-testid="slide">{children}</div> : null,
    styled,
    useMediaQuery: () => false,
    useTheme: () => ({
      palette: { icon: { default: 'icon-default' } },
    }),
    WorkspaceIcon: () => <svg data-testid="ws-icon" />,
  };
});

vi.mock('../../layout/Header/Header.styles', () => ({
  CMenuContainer: ({ children }: any) => <div data-testid="menu-container">{children}</div>,
}));

vi.mock('./SpaceSwitcher', () => ({
  OrgMenu: ({ open }: any) => (open ? <div data-testid="org-menu">org menu</div> : null),
}));

vi.mock('./WorkspaceSwitcher', () => ({
  default: ({ open }: any) =>
    open ? <div data-testid="workspace-switcher">workspace switcher</div> : null,
}));

import { MobileOrgWksSwither } from './MobileViewSwitcher';

describe('MobileOrgWksSwither', () => {
  const router = { push: vi.fn() };
  const organization = { id: 'org-1', name: 'Acme' };

  it('renders the dashboard switcher button', () => {
    render(<MobileOrgWksSwither organization={organization} router={router as any} />);
    expect(screen.getByTestId('dashboard-icon')).toBeInTheDocument();
  });

  it('opens the slide-in menu when the trigger is clicked', async () => {
    const user = userEvent.setup();
    render(<MobileOrgWksSwither organization={organization} router={router as any} />);

    expect(screen.queryByTestId('slide')).not.toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /contexts/i }));
    expect(screen.getByTestId('slide')).toBeInTheDocument();
    // The slide-in contains OrgMenu and WorkspaceSwitcher with open=true
    expect(screen.getByTestId('org-menu')).toBeInTheDocument();
    expect(screen.getByTestId('workspace-switcher')).toBeInTheDocument();
  });

  it('renders OrgMenu and WorkspaceSwitcher inside the menu container once opened', async () => {
    const user = userEvent.setup();
    render(<MobileOrgWksSwither organization={organization} router={router as any} />);
    const trigger = screen.getByRole('button', { name: /contexts/i });

    await user.click(trigger);
    expect(screen.getByTestId('menu-container')).toBeInTheDocument();
    expect(screen.getByTestId('org-menu')).toBeInTheDocument();
    expect(screen.getByTestId('workspace-switcher')).toBeInTheDocument();
  });
});
