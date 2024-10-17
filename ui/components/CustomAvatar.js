import React from 'react';
import { Badge, styled, useTheme } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

const StyledBadge = styled(Badge)(({ theme, color }) => ({
  '& .MuiBadge-badge': {
    backgroundColor: color || theme.palette.success.main,
    color: theme.palette.background.paper,
    boxShadow: `0 0 0 2px ${theme.palette.background.paper}`,
    '&::after': {
      position: 'absolute',
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
      transform: 'scale(2.4)',
      opacity: 0,
    },
  },
}));

const Root = styled('div')(({ theme }) => ({
  display: 'flex',
  '& > *': {
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(-0.75),
  },
}));

function CustomAvatar({ children, color }) {
  const theme = useTheme();

  return (
    <UsesSistent>
      <Root theme={theme}>
        <StyledBadge
          overlap="circular"
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}
          variant="dot"
          color={color}
          theme={theme}
        >
          {children}
        </StyledBadge>
      </Root>
    </UsesSistent>
  );
}

export default CustomAvatar;
