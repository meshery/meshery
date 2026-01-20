import React from 'react';
import { styled, Badge } from '@sistent/sistent';
import type { Theme } from '@sistent/sistent';

type StyledBadgeProps = {
  color?: string;
};

const StyledBadge = styled(Badge)<StyledBadgeProps>(({ theme, color }: { theme: Theme; color?: string }) => ({
  '& .MuiBadge-badge': {
    backgroundColor: color || theme.palette.background?.brand?.default || theme.palette.primary.main,
    color: color || theme.palette.background?.brand?.default || theme.palette.primary.main,
    '&::after': {
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      borderRadius: '50%',
      animation: 'ripple 1.2s infinite ease-in-out',
      border: '1px solid currentColor',
      content: '""',
    },
  },
  '@keyframes ripple': {
    '0%': {
      transform: 'scale(.8)',
      opacity: 1,
    },
    '100%': {
      transform: 'scale(1.8)',
      opacity: 0,
    },
  },
}));

const Container = styled('div')(({ theme }: { theme: Theme }) => ({
  display: 'flex',
  '& > *': {
    marginLeft: theme.spacing(0.5),
    marginRight: -theme.spacing(0.75),
  },
}));

type BadgeAvatarsProps = {
  children: React.ReactNode;
  color?: string;
};

export default function BadgeAvatars({ children, color }: BadgeAvatarsProps) {
  return (
    <Container>
      <StyledBadge
        overlap="circular"
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        variant="dot"
        {...(color && { color: color as 'primary' | 'secondary' | 'default' | 'error' | 'info' | 'success' | 'warning' })}
      >
        {children}
      </StyledBadge>
    </Container>
  );
}
