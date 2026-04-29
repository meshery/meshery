import React from "react";

const AlertIcon = ({ height, width, fill, style = {} }) => {
  return (
    <svg
      style={style}
      height={height}
      width={width}
      fill="#fff"
      viewBox="0 0 20 20"
    >
      <path d="M11.6042 11.6667H9.89591V7.49999H11.6042V11.6667ZM11.6042 15H9.89591V13.3333H11.6042V15ZM1.35425 17.5H20.1459L10.7501 1.66666L1.35425 17.5Z" fill={fill} />

    </svg>
  );
};

export default AlertIcon;