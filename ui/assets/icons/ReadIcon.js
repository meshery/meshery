import React from "react";

const ReadIcon = ({ height, width, fill, style = {} }) => {
  return (
    <svg version="1.1" xmlns="http://www.w3.org/2000/svg"
      x="0px"
      y="0px"
      style={style}
      height={height}
      width={width}
      fill="#fff"
      viewBox="0 0 20 16"
    >
      <path d="M18,0H2C0.9,0,0,0.9,0,2l0,12c0,1.1,0.9,2,2,2h16c1.1,0,2-0.9,2-2V2C20,0.9,19.1,0,18,0z M18,4l-8,5L2,4V2l8,5l8-5V4z" fill={fill} />
    </svg>


  );
};

export default ReadIcon;