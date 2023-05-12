import React from "react";
import PropTypes from "prop-types";
import { AppBar as MuiAppBar } from "@mui/material";

export const AppBar = ({ props, ...rest }) => {
  return <MuiAppBar {...props} {...rest} />
};

AppBar.propTypes = {
  classes : PropTypes.string,
  color : PropTypes.oneOf(['transparent']),
  enableColorOnDark : PropTypes.bool,
  position : PropTypes.oneOf(['fixed', 'absolute', 'sticky', 'static', 'relative'])
}

AppBar.defaultProps = {
  position : "fixed",
  color : "primary"
}