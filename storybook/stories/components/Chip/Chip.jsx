import React from "react";
import PropTypes from  "prop-types";
import { Chip as MuiChip } from "@mui/material";

export const Chip = ({  label = 'Chip Filled', ...rest }) => {
  return <MuiChip label={label} {...rest} />
}

Chip.propTypes = {
  label : PropTypes.string,
  clickable : PropTypes.bool,
  color : PropTypes.oneOf(['default', 'primary', 'secondary', 'error', 'info', 'success', 'warning']),
  disabled : PropTypes.bool,
  size : PropTypes.oneOf(['small', 'medium']),
  variant : PropTypes.oneOf(['filled', 'outlined'])
}

Chip.defaulProps = {
  label : "Chip Filled",
  variant : "filled",
  size : "small",
  color : "primary"
}