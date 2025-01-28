import React from 'react';
import { styled, Badge } from '@layer5/sistent';

const StyledBadge = styled(Badge)(({ theme, color }) => ({
  '& .MuiBadge-badge': {
    backgroundColor: color || '#44b700',
    color: color || '#44b700',
    boxShadow: `0 0 0 2px ${theme.palette.background.paper}`,
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

const Container = styled('div')(({ theme }) => ({
  display: 'flex',
  '& > *': {
    marginLeft: theme.spacing(0.5),
    marginRight: -theme.spacing(0.75),
  },
}));

export default function BadgeAvatars({ children, color }) {
  return (
    <Container>
      <StyledBadge
        overlap="circular"
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        variant="dot"
        color={color}
      >
        {children}
      </StyledBadge>
    </Container>
  );
}
