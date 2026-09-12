import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any, _options?: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => {
      const Resolved = typeof Component === 'string' ? Component : Component?.displayName || 'div';
      // Render as a tag we can query so tests can verify it mounts.
      if (typeof Resolved === 'string') {
        return React.createElement(Resolved, props, children);
      }
      return <div {...props}>{children}</div>;
    };
    StyledComponent.displayName = 'StyledHeaderMock';
    return StyledComponent;
  };

  return {
    styled,
    AppBar: 'header',
    Toolbar: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Paper: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    MenuIcon: (props: any) => <svg data-testid="menu-icon" {...props} />,
    IconButton: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    darkTeal: { main: '#000' },
  };
});

import {
  HeaderAppBar,
  StyledToolbar,
  UserContainer,
  SettingsWrapper,
  PageTitleWrapper,
  MenuIconButton,
  UserSpan,
  HeaderIcons,
  CBadge,
  CBadgeContainer,
  CMenuContainer,
  IconButtonAvatar,
  IconButtonMenu,
  UserInfoContainer,
} from './Header.styles';

describe('Header.styles', () => {
  it('exports all styled components', () => {
    expect(HeaderAppBar).toBeDefined();
    expect(StyledToolbar).toBeDefined();
    expect(UserContainer).toBeDefined();
    expect(SettingsWrapper).toBeDefined();
    expect(PageTitleWrapper).toBeDefined();
    expect(MenuIconButton).toBeDefined();
    expect(UserSpan).toBeDefined();
    expect(HeaderIcons).toBeDefined();
    expect(CBadge).toBeDefined();
    expect(CBadgeContainer).toBeDefined();
    expect(CMenuContainer).toBeDefined();
    expect(IconButtonAvatar).toBeDefined();
    expect(IconButtonMenu).toBeDefined();
    expect(UserInfoContainer).toBeDefined();
  });

  it('renders styled components with children', () => {
    render(
      <UserContainer>
        <PageTitleWrapper>
          <span data-testid="title">Header Title</span>
        </PageTitleWrapper>
      </UserContainer>,
    );

    expect(screen.getByTestId('title')).toHaveTextContent('Header Title');
  });

  it('renders CBadge inside CBadgeContainer', () => {
    render(
      <CBadgeContainer>
        <CBadge>5</CBadge>
      </CBadgeContainer>,
    );

    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('renders MenuIconButton as a button-like element', () => {
    render(
      <MenuIconButton>
        <span>menu</span>
      </MenuIconButton>,
    );

    expect(screen.getByText('menu')).toBeInTheDocument();
  });

  it('renders UserInfoContainer with nested user span', () => {
    render(
      <UserInfoContainer>
        <UserSpan>user</UserSpan>
      </UserInfoContainer>,
    );

    expect(screen.getByText('user')).toBeInTheDocument();
  });
});
