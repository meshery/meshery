import React from 'react';
import PropTypes from 'prop-types';
import { Tooltip } from '@material-ui/core';

const ReusableTooltip = ({ title, children, onClick, placement }) => {
  return (
    <Tooltip title={title} placement={placement} onClick={onClick} arrow interactive>
      {children}
    </Tooltip>
  );
};

ReusableTooltip.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  onClick: PropTypes.func,
  placement: PropTypes.string,
};

ReusableTooltip.defaultProps = {
  onClick: () => {},
  // placement: "top",
};

export default ReusableTooltip;
