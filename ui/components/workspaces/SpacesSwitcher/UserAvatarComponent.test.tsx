import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://meshery.io',
}));

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={title}>
      {children}
    </div>
  ),
  Link: ({ children, href }: any) => (
    <a href={href} data-testid="link">
      {children}
    </a>
  ),
  Avatar: ({ alt, src }: any) => <div data-testid="avatar" data-alt={alt} data-src={src} />,
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({
    palette: {},
    breakpoints: { down: () => '@media (max-width: 0)' },
    spacing: (n: number) => `${n * 8}px`,
  }),
}));

vi.mock('./styles', () => ({
  StyledAvatarContainer: ({ children }: any) => (
    <div data-testid="avatar-container">{children}</div>
  ),
  StyledUserDetailsContainer: ({ children }: any) => <div>{children}</div>,
  StyledUpdatedText: ({ children }: any) => <span data-testid="email-text">{children}</span>,
}));

import UserAvatarComponent from './UserAvatarComponent';

describe('UserAvatarComponent', () => {
  it('renders the avatar with name, src and a link to the user profile', () => {
    render(
      <UserAvatarComponent
        userData={{
          id: 'user-1',
          firstName: 'Alice',
          lastName: 'Smith',
          avatarUrl: 'https://avatar.example/a.png',
          email: 'alice@example.com',
        }}
      />,
    );
    expect(screen.getByTestId('avatar')).toHaveAttribute('data-alt', 'Alice');
    expect(screen.getByTestId('avatar')).toHaveAttribute(
      'data-src',
      'https://avatar.example/a.png',
    );
    expect(screen.getByTestId('link')).toHaveAttribute('href', 'https://meshery.io/user/user-1');
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Smith', { exact: false })).toBeInTheDocument();
    expect(screen.getByTestId('email-text')).toHaveTextContent('alice@example.com');
  });

  it('uses the tooltip to display the full name', () => {
    render(
      <UserAvatarComponent
        userData={{
          id: 'user-2',
          firstName: 'Bob',
          lastName: 'Jones',
          avatarUrl: '',
          email: 'bob@example.com',
        }}
      />,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'Bob Jones');
  });

  it('handles missing lastName gracefully', () => {
    render(
      <UserAvatarComponent
        userData={{
          id: 'user-3',
          firstName: 'Solo',
          avatarUrl: '',
          email: 's@e.com',
        }}
      />,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'Solo ');
  });

  it('handles fully missing userData (renders fallbacks)', () => {
    render(<UserAvatarComponent userData={undefined as any} />);
    expect(screen.getByTestId('avatar-container')).toBeInTheDocument();
  });
});
