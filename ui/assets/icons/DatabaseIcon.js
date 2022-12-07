import React from "react";

const DatabaseIcon = (props) => {
  return (
    <svg
      viewBox="0 0 448 512"
      xmlns="http://www.w3.org/2000/svg"
      width={props.width ? props.width : "24px"}
      height={props.height ? props.height : "24px"}
      onClick={props.onClick}
      className={props.className}
      color={props.color ? props.color : "unset"}
      fontSize={props.fontSize ? props.fontSize : "unset"}
      style={{ ...props.style }}
      fill={props.fill ? props.fill : "currentColor"}
    >
      <path d="M448 73.12v45.75C448 159.1 347.6 192 224 192S0 159.1 0 118.9V73.12C0 32.88 100.4 0 224 0S448 32.88 448 73.12zM448 176v102.9C448 319.1 347.6 352 224 352S0 319.1 0 278.9V176c48.12 33.12 136.2 48.62 224 48.62S399.9 209.1 448 176zM448 336v102.9C448 479.1 347.6 512 224 512s-224-32.88-224-73.13V336c48.12 33.13 136.2 48.63 224 48.63S399.9 369.1 448 336z" />
    </svg>
  );
};

export default DatabaseIcon;

