import React from 'react';
import PropTypes from 'prop-types';
import { Button as MuiButton } from "@mui/material";

/**
 * Primary UI component for user interaction
 */
export const Button = ({
  variant = 'outlined',
  children = 'Button',
  ...props
}) => (
  <MuiButton variant={variant} {...props}>
    {children}
  </MuiButton>
);

Button.propTypes = {
  variant : PropTypes.oneOf(['text', 'outlined', 'contained']),
  color : PropTypes.oneOf([
    'inherit',
    'primary',
    'secondary',
    'success',
    'error',
    'info',
    'warning',
  ]),
  fullWidth : PropTypes.bool,
  disabled : PropTypes.bool,
  disabledElevation : PropTypes.bool,
  disableFocusRipple : PropTypes.bool,
  href : PropTypes.string,
  children : PropTypes.string,
};

Button.defaultProps = {
  variant : "contained",
  size : "medium",
  color : "primary",
};