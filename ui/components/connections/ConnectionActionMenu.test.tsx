import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

let isMobile = false;

vi.mock('@sistent/sistent', () => ({
  SettingsIcon: () => <svg data-testid="settings-icon" />,
  CopyLinkIcon: () => <svg data-testid="copy-link-icon" />,
  ListItemIcon: ({ children }: any) => <span>{children}</span>,
  MenuItem: ({ children, onClick, permissionKey }: any) => (
    <button
      type="button"
      onClick={onClick}
      data-permission-key={permissionKey?.id ?? ''}
      data-testid="menu-item"
    >
      {children}
    </button>
  ),
  MenuList: ({ children }: any) => <div data-testid="menu-list">{children}</div>,
  Popover: ({ open, children }: any) => (open ? <div data-testid="popover">{children}</div> : null),
  BottomSheet: ({ children, open, title }: any) =>
    open ? (
      <div data-testid="bottom-sheet" data-title={title}>
        {children}
      </div>
    ) : null,
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({
    breakpoints: { down: (key: string) => `(max-width:${key})` },
  }),
  useMediaQuery: (query: string) => query === '(max-width:sm)' && isMobile,
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: {
    LifecycleManagementEditConnection: { id: 'edit-connection', function: 'connection' },
  },
}));

vi.mock('./styles', () => ({
  ActionListItem: ({ children }: any) => <div data-testid="action-list-item">{children}</div>,
  ActionButton: ({ children, onClick, permissionKey, ...props }: any) => (
    <button
      type="button"
      onClick={onClick}
      data-permission-key={permissionKey?.id ?? ''}
      {...props}
    >
      {children}
    </button>
  ),
}));

vi.mock('../../css/icons.styles', () => ({ iconMedium: {} }));

import { ConnectionActionMenu } from './ConnectionActionMenu';

describe('ConnectionActionMenu', () => {
  const onClose = vi.fn();
  const onConfigure = vi.fn();
  const onConfigureControllers = vi.fn();
  const onCopyLink = vi.fn();

  beforeEach(() => {
    isMobile = false;
    onClose.mockReset();
    onConfigure.mockReset();
    onConfigureControllers.mockReset();
    onCopyLink.mockReset();
  });

  const renderMenu = (overrides = {}) =>
    render(
      <ConnectionActionMenu
        anchorEl={document.createElement('button')}
        open
        onClose={onClose}
        onConfigure={onConfigure}
        onConfigureControllers={onConfigureControllers}
        onCopyLink={onCopyLink}
        {...overrides}
      />,
    );

  it('renders a popover of actions on desktop viewports', () => {
    renderMenu();
    expect(screen.getByTestId('popover')).toBeInTheDocument();
    expect(screen.queryByTestId('bottom-sheet')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: /configure controllers/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /copy link/i })).toBeInTheDocument();
  });

  it('omits actions whose handlers are not provided', () => {
    renderMenu({ onConfigureControllers: undefined, onCopyLink: undefined });
    expect(screen.getByRole('button', { name: /^configure$/i })).toBeInTheDocument();
    expect(
      screen.queryByRole('button', { name: /configure controllers/i }),
    ).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /copy link/i })).not.toBeInTheDocument();
  });

  it('invokes the action and closes the popover', async () => {
    const user = userEvent.setup();
    renderMenu();
    await user.click(screen.getByRole('button', { name: /^configure$/i }));
    expect(onConfigure).toHaveBeenCalledTimes(1);
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('renders a BottomSheet of actions on mobile viewports', async () => {
    isMobile = true;
    const user = userEvent.setup();
    renderMenu();

    expect(screen.queryByTestId('popover')).not.toBeInTheDocument();
    expect(screen.getByTestId('bottom-sheet')).toHaveAttribute('data-title', 'Connection Actions');
    expect(screen.getByTestId('menu-list')).toBeInTheDocument();
    expect(screen.getAllByTestId('menu-item')).toHaveLength(3);

    await user.click(screen.getByRole('button', { name: /copy link/i }));
    expect(onCopyLink).toHaveBeenCalledTimes(1);
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('does not render either surface when closed', () => {
    isMobile = true;
    renderMenu({ open: false });
    expect(screen.queryByTestId('bottom-sheet')).not.toBeInTheDocument();
    expect(screen.queryByTestId('popover')).not.toBeInTheDocument();
  });
});
