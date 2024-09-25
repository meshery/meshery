import React from 'react';
import Badge from '@mui/material/Badge';
import { styled, keyframes } from '@mui/material/styles';

const rippleAnimation = keyframes`
  0% {
    transform: scale(.8);
    opacity: 1;
  }
  100% {
    transform: scale(1.8);
    opacity: 0;
  }
`;

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
      animation: `${rippleAnimation} 1.2s infinite ease-in-out`,
      border: '1px solid currentColor',
      content: '""',
    },
  },
}));

const BadgeAvatars = ({ children, color }) => {
  return (
    <StyledBadge
      overlap="circular"
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
      }}
      variant="dot"
      color={color}
      sx={{
        display: 'flex',
        '& > *': {
          ml: 0.5,
          mr: -0.75,
        },
      }}
    >
      {children}
    </StyledBadge>
  );
};

export default BadgeAvatars;
