import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = ({ children, ...props }: any) => {
      if (typeof Component === 'string') {
        return React.createElement(Component, props, children);
      }
      return <div {...props}>{children}</div>;
    };
    Styled.displayName = 'StyledNotificationCenterMock';
    return Styled;
  };

  return {
    styled,
    alpha: (color: string) => color,
    Badge: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    Drawer: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Grid: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    IconButton: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
    InfoIcon: () => null,
    createTheme: () => ({ breakpoints: { up: () => '', down: () => '' } }),
    notificationColors: {
      info: { main: '#info' },
      error: { main: '#error', dark: '#error_dark' },
      warning: { main: '#warning', light: '#warning_light' },
      success: { main: '#success' },
    },
  };
});

vi.mock('../../../assets/icons/AlertIcon', () => ({ default: () => null }));
vi.mock('../../../assets/icons/ErrorIcon', () => ({ default: () => null }));
vi.mock('../../../assets/icons/ReadIcon', () => ({ default: () => null }));

import {
  DarkBackdrop,
  SideList,
  StyledSubtitle,
  NotificationButton,
  NotificationDrawer,
  StyledNotificationDrawer,
  FullView,
  PeekView,
  Container,
  Header,
  NotificationContainer,
  Title,
  TitleBellIcon,
  SeverityChip,
  SeverityChips,
  Notification,
  DetailsContainer,
  ListDetails,
  ListItem,
  SoicialListItem,
  ListButton,
  Expanded,
  StyledAvatarStack,
  ActorAvatar,
  Message,
  GridItem,
  MenuPaper,
  OptionList,
  OptionListItem,
  SocialListItem,
  StyledBadge,
  Root,
  Summary,
} from './notificationCenter.style';

describe('notificationCenter styled components', () => {
  it('exports all named styled components', () => {
    [
      DarkBackdrop,
      SideList,
      StyledSubtitle,
      NotificationButton,
      NotificationDrawer,
      StyledNotificationDrawer,
      FullView,
      PeekView,
      Container,
      Header,
      NotificationContainer,
      Title,
      TitleBellIcon,
      SeverityChip,
      SeverityChips,
      Notification,
      DetailsContainer,
      ListDetails,
      ListItem,
      SoicialListItem,
      ListButton,
      Expanded,
      StyledAvatarStack,
      ActorAvatar,
      Message,
      GridItem,
      MenuPaper,
      OptionList,
      OptionListItem,
      SocialListItem,
      StyledBadge,
      Root,
      Summary,
    ].forEach((component) => {
      expect(component).toBeDefined();
    });
  });

  it('renders Container with children', () => {
    render(
      <Container>
        <span data-testid="content">Content</span>
      </Container>,
    );
    expect(screen.getByTestId('content')).toHaveTextContent('Content');
  });

  it('renders title content', () => {
    render(
      <Title>
        <span data-testid="title">A title</span>
      </Title>,
    );
    expect(screen.getByTestId('title')).toBeInTheDocument();
  });

  it('renders SeverityChip / SeverityChips composition', () => {
    render(
      <SeverityChips>
        <SeverityChip>
          <span data-testid="chip-label">warning</span>
        </SeverityChip>
      </SeverityChips>,
    );
    expect(screen.getByTestId('chip-label')).toHaveTextContent('warning');
  });

  it('renders DarkBackdrop component', () => {
    render(<DarkBackdrop data-testid="backdrop" />);
    expect(screen.getByTestId('backdrop')).toBeInTheDocument();
  });

  it('renders SideList component', () => {
    render(
      <SideList data-testid="side-list">
        <span>inside</span>
      </SideList>,
    );
    expect(screen.getByTestId('side-list')).toBeInTheDocument();
  });
});
