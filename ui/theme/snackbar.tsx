import {
  BasicMarkdown,
  Box,
  CheckCircleIcon as CheckCircle,
  CircularProgress,
  ErrorIcon as Error,
  InfoIcon as Info,
  lighten,
  styled,
  WarningIcon as Warning,
} from '@sistent/sistent';
import { SnackbarContent } from 'notistack';
import { forwardRef } from 'react';

const StyledSnackbarContent = styled(SnackbarContent, {
  shouldForwardProp: (prop) => prop !== 'variant',
})(({ theme, variant }) => {
  const variantTextColors = {
    success: theme.palette.text.success,
    info: theme.palette.text.info,
    warning: theme.palette.text.warning,
    error: theme.palette.text.error,
  };

  const baseColor = variantTextColors[variant] || variantTextColors.info;

  const backgroundColor = theme.palette.mode === 'light' ? lighten(baseColor, 0.95) : '#323232';

  return {
    backgroundColor,
    color: baseColor,
    pointerEvents: 'auto',
    borderRadius: '0.3rem',
    boxShadow: `0 0px 4px ${theme.palette.background.tabs}`,
  };
});

const StyledCircularProgress = styled(CircularProgress)(({ theme }) => ({
  color: theme.palette.text.info,
  marginRight: '0.75rem',
  height: '24px !important',
  width: '24px !important',
}));

const iconProps = { style: { marginRight: '0.5rem' } };

const contentStyle = {
  display: 'flex',
  alignItems: 'center',
  padding: '0.5rem 1rem',
  width: '100%',
};

// `BasicMarkdown` emits sibling block elements (<p>, <ul>, ...) with no wrapper
// of its own. Dropped straight into the row above, those blocks become sibling
// flex items and lay out horizontally, so a multi-block MeshKit error renders
// as "Unable to create the environmentmeshery-server-1448". Its own column
// keeps them stacked. `flexBasis: auto` preserves the intrinsic width a
// single-line toast has today, and `minWidth: 0` lets a long unbreakable string
// wrap instead of overflowing the toast.
const messageStyle = {
  display: 'flex',
  flexDirection: 'column',
  flex: '1 1 auto',
  minWidth: 0,
  overflowWrap: 'anywhere',
};

export const ThemeResponsiveSnackbar = forwardRef((props, forwardedRef) => {
  // notistack v3 custom content props: the snackbar identifier arrives as
  // `id` (React strips `key`, so it must never be read from props). `action`
  // callbacks receive that id so closeSnackbar(id) dismisses only this toast.
  const { id, variant, message, action } = props;

  // Function to determine the icon based on variant
  const getIcon = () => {
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
        return <StyledCircularProgress />;
      default:
        return null;
    }
  };

  return (
    <StyledSnackbarContent ref={forwardedRef} variant={variant}>
      <div data-testid={`SnackbarContent-${variant}`} style={contentStyle}>
        {getIcon()}
        <div data-testid="SnackbarContent-message" style={messageStyle}>
          <BasicMarkdown content={message} />
        </div>
        <Box marginLeft={'auto'} paddingLeft={'0.5rem'}>
          {action && action(id)}
        </Box>
      </div>
    </StyledSnackbarContent>
  );
});

ThemeResponsiveSnackbar.displayName = 'ThemeResponsiveSnackbar';
