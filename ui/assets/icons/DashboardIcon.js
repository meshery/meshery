import React from "react";

const DashboardIcon = (props) => {
  return (
    <svg
      width={props.width ? props.width : "24px"}
      height={props.height ? props.height : "24px"}
      fill={props.fill ? props.fill : "currentColor"}
      onClick={props.onClick}
      className={props.className}
      color={props.color ? props.color : "unset"}
      fontSize={props.fontSize ? props.fontSize : "unset"}
      style={{ ...props.style }}
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <g>
        <path fill="none" d="M0 0h24v24H0z"/>
        <path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>
      </g>
    </svg>
  );
};

export default DashboardIcon;

