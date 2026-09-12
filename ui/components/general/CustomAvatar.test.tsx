import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  // styled(Component) returns a function which returns a wrapped Component
  // that just forwards all props through. This way StyledBadge -> Badge.
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledMock';
    return StyledComponent;
  };

  return {
    styled,
    Badge: ({ children, color, variant, overlap, anchorOrigin }: any) => (
      <div
        data-testid="badge"
        data-color={color}
        data-variant={variant}
        data-overlap={overlap}
        data-anchor={JSON.stringify(anchorOrigin)}
      >
        {children}
      </div>
    ),
  };
});

import BadgeAvatars from './CustomAvatar';

describe('BadgeAvatars (CustomAvatar)', () => {
  it('renders children inside the styled badge', () => {
    render(
      <BadgeAvatars>
        <span data-testid="avatar-child">A</span>
      </BadgeAvatars>,
    );

    expect(screen.getByTestId('avatar-child')).toBeInTheDocument();
    expect(screen.getByTestId('badge')).toBeInTheDocument();
  });

  it('forwards the color prop to the Badge', () => {
    render(
      <BadgeAvatars color="#ff00ff">
        <span>child</span>
      </BadgeAvatars>,
    );

    expect(screen.getByTestId('badge')).toHaveAttribute('data-color', '#ff00ff');
  });

  it('uses the default Badge props (variant=dot, overlap=circular, anchorOrigin bottom-right)', () => {
    render(
      <BadgeAvatars>
        <span>child</span>
      </BadgeAvatars>,
    );

    const badge = screen.getByTestId('badge');
    expect(badge).toHaveAttribute('data-variant', 'dot');
    expect(badge).toHaveAttribute('data-overlap', 'circular');
    expect(badge.getAttribute('data-anchor')).toBe(
      JSON.stringify({ vertical: 'bottom', horizontal: 'right' }),
    );
  });

  it('renders multiple children', () => {
    render(
      <BadgeAvatars color="red">
        <span data-testid="child-a">A</span>
        <span data-testid="child-b">B</span>
      </BadgeAvatars>,
    );

    expect(screen.getByTestId('child-a')).toBeInTheDocument();
    expect(screen.getByTestId('child-b')).toBeInTheDocument();
  });
});
