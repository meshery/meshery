import React from "react";
import PropTypes from  "prop-types";
import { Avatar as MuiAvatar } from '@mui/material';

export const Avatar = ({ ...rest }) => {
  return <MuiAvatar  {...rest} />
}

Avatar.propTypes = {
  color : PropTypes.oneOf(['#44b700']),
  backgroundColor : PropTypes.oneOf(['#44b700']),
  src : PropTypes.string,
}

Avatar.defaulProps = {
  color : "#44b700",
  backgroundColor : "#44b700",

}