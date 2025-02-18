import { BasicMarkdown, CircularProgress, styled } from '@layer5/sistent';
import { SnackbarContent } from 'notistack';
import { forwardRef } from 'react';
import { CheckCircle, Error, Info, Warning } from '@mui/icons-material';
import { lighten } from '@mui/material';

const drawerWidth = 256;

export const StyledFooterText = styled('span')({
  cursor: 'pointer',
  display: 'inline',
  verticalAlign: 'middle',
});

export const StyledRoot = styled('div')({
  display: 'flex',
  minHeight: '100vh',
  height: '100vh',
});

export const StyledFooterBody = styled('footer')(({ theme }) => ({
  padding: theme.spacing(2),
  backgroundColor:
    theme.palette.mode === 'dark' ? theme.palette.background.card : theme.palette.common.white,
}));
export const StyledMainContent = styled('main')(({ theme }) => ({
  backgroundColor:
    theme.palette.mode === 'dark'
      ? theme.palette.background.elevatedComponents
      : theme.palette.background.hover,
  flex: 1,
  padding: '48px 36px 24px',
  [theme.breakpoints.down('sm')]: {
    padding: '24px 16px 16px',
  },
}));

export const StyledAppContent = styled('div')({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  position: 'relative',
  overflow: 'visible',
});

export const StyledContentWrapper = styled('div')({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  overflow: 'auto',
  minHeight: 0,
});

export const StyledDrawer = styled('nav', {
  shouldForwardProp: (prop) => prop !== 'isDrawerCollapsed',
})(({ theme, isDrawerCollapsed }) => ({
  [theme.breakpoints.up('sm')]: {
    width: isDrawerCollapsed ? theme.spacing(8.4) + 1 : drawerWidth,
    flexShrink: 0,
  },
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: isDrawerCollapsed
      ? theme.transitions.duration.leavingScreen
      : theme.transitions.duration.enteringScreen,
  }),
  '& > div:first-child': {
    height: 'inherit',
    width: 'inherit',
  },
  height: '100%',
  overflow: 'visible',
  paddingRight: '4rem',
  [theme.breakpoints.up('xs')]: {
    paddingRight: '0',
  },
}));

const StyledSnackbarContent = styled(SnackbarContent)(({ theme, variant }) => {
  const notificationColors = {
    success: theme.palette.text.success,
    info: theme.palette.text.info,
    warning: theme.palette.text.warning,
    error: theme.palette.text.error,
  };

  const baseColor = notificationColors[variant] || notificationColors.info;

  const backgroundColor = theme.palette.mode === 'light' ? lighten(baseColor, 0.95) : '#323232';

  return {
    backgroundColor,
    color: baseColor,
    pointerEvents: 'auto',
    borderRadius: '0.3rem',
    boxShadow: `0 0px 4px ${theme.palette.background.tabs}`,
  };
});

export const ThemeResponsiveSnackbar = forwardRef((props, forwardedRef) => {
  const { variant, message, action, key, theme } = props;

  // Function to determine the icon based on variant
  const getIcon = () => {
    const iconProps = { style: { marginRight: '0.5rem' } };
    switch (variant) {
      case 'error':
        return <Error {...iconProps} />;
      case 'success':
        return <CheckCircle {...iconProps} />;
      case 'warning':
        return <Warning {...iconProps} />;
      case 'info':
        return <Info {...iconProps} />;
      case 'loading':
        return <CircularProgress size={24} {...iconProps} />;
      default:
        return null;
    }
  };

  return (
    <StyledSnackbarContent ref={forwardedRef} variant={variant} theme={theme}>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          padding: '0.5rem',
          width: '100%',
        }}
      >
        {getIcon()}
        <BasicMarkdown content={message} />
        <div style={{ marginLeft: 'auto' }}>{action && action(key)}</div>
      </div>
    </StyledSnackbarContent>
  );
});

ThemeResponsiveSnackbar.displayName = 'ThemeResponsiveSnackbar';
