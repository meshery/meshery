import React from "react";

const RightArrowIcon = ({ width, height, primaryFill = '#455a64', style={} }) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={width}
      height={height}
      style={style}
      viewBox="0 0 24 24"
      fill={primaryFill}
    >
      <polygon points="7.293 4.707 14.586 12 7.293 19.293 8.707 20.707 17.414 12 8.707 3.293 7.293 4.707" />
    </svg>
  );
};

export default RightArrowIcon;
