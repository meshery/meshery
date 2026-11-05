import React from "react";

const LeftArrowIcon = ({ width, height, primaryFill = "#455a64", style={} }) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={width}
      height={height}
      style={style}
      viewBox="0 0 24 24"
      fill={primaryFill}
    >
      <polygon points="15.293 3.293 6.586 12 15.293 20.707 16.707 19.293 9.414 12 16.707 4.707 15.293 3.293" />
    </svg>
  );
};

export default LeftArrowIcon;
