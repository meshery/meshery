import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/assets/icons/OrgOutlinedIcon', () => ({
  default: () => <svg data-testid="org-icon" />,
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: { height: 20, width: 20 },
  iconMedium: { height: 24, width: 24 },
  iconLarge: { height: 32, width: 32 },
  iconXLarge: {},
}));

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Box: ({ children, className }: any) => <div className={className}>{children}</div>,
    BottomSheet: ({ children, open, title }: any) =>
      open ? (
        <div data-testid="bottom-sheet" role="dialog">
          {title ? <h2>{title}</h2> : null}
          {children}
        </div>
      ) : null,
    Button: ({ children, onClick, variant, fullWidth, permissionKey, sx, ...rest }: any) => (
      <button onClick={onClick} {...rest}>
        {children}
      </button>
    ),
    Typography: ({ children }: any) => <span>{children}</span>,
    Grid2: ({ children }: any) => <div>{children}</div>,
    styled,
    useTheme: () => ({
      palette: {
        common: { white: '#ffffff' },
        icon: { default: 'icon-default' },
        text: { secondary: 'text-secondary', default: 'text-default' },
        background: { brand: { default: 'brand-default' } },
        divider: 'divider-color',
      },
      spacing: (value: number) => `${value * 8}px`,
    }),
    WorkspaceIcon: () => <svg data-testid="ws-icon" />,
  };
});

vi.mock('@/utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({
    openModal: vi.fn(),
    setCreateNewWorkspaceModalOpen: vi.fn(),
    setSelectedWorkspace: vi.fn(),
  }),
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

  it('renders the org and workspace switcher buttons', () => {
    render(<MobileOrgWksSwither organization={organization} router={router as any} />);
    expect(screen.getByTestId('org-icon')).toBeInTheDocument();
    expect(screen.getByTestId('ws-icon')).toBeInTheDocument();
  });

  it('opens the bottom sheet when the trigger is clicked', async () => {
    const user = userEvent.setup();
    render(<MobileOrgWksSwither organization={organization} router={router as any} />);

    expect(screen.queryByTestId('bottom-sheet')).not.toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /contexts/i }));
    expect(screen.getByTestId('bottom-sheet')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /organization & workspace/i })).toBeInTheDocument();
    expect(screen.getByTestId('org-menu')).toBeInTheDocument();
    expect(screen.getByTestId('workspace-switcher')).toBeInTheDocument();
  });

  it('renders OrgMenu and WorkspaceSwitcher inside the bottom sheet once opened', async () => {
    const user = userEvent.setup();
    render(<MobileOrgWksSwither organization={organization} router={router as any} />);
    const trigger = screen.getByRole('button', { name: /contexts/i });

    await user.click(trigger);
    expect(screen.getByTestId('bottom-sheet')).toBeInTheDocument();
    expect(screen.getByTestId('org-menu')).toBeInTheDocument();
    expect(screen.getByTestId('workspace-switcher')).toBeInTheDocument();
  });
});
