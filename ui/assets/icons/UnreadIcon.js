import React from "react";

const UnreadIcon = ({ height, width, fill, style = {} }) => {
  return (
    <svg
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      xml:space="preserve"
      x="0px"
      y="0px"
      style={style}
      height={height}
      width={width}
      fill="#fff"
      viewBox="0 0 20 19"
    >
      <path
        d="M20,7c0-0.7-0.4-1.3-0.9-1.7L10,0L1,5.3C0.4,5.7,0,6.3,0,7v10c0,1.1,0.9,2,2,2h16c1.1,0,2-0.9,2-2L20,7z M10,12L1.7,6.8
		L10,2l8.3,4.8L10,12z"
        fill={fill}
      />
    </svg>
  );
};

export default UnreadIcon;
