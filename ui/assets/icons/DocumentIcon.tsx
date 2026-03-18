import React from "react";

const DocumentIcon = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={props.width ? props.width : "24px"}
      height={props.height ? props.height : "24px"}
      fill={props.fill ? props.fill : "currentColor"}
      onClick={props.onClick}
      className={props.className}
      color={props.color ? props.color : "unset"}
      fontSize={props.fontSize ? props.fontSize : "unset"}
      style={{ ...props.style }}
      viewBox="0 0 512 512"
    >
      <title>Document Text</title>
      <path
        d="M416 221.25V416a48 48 0 01-48 48H144a48 48 0 01-48-48V96a48 48 0 0148-48h98.75a32 32 0 0122.62 9.37l141.26 141.26a32 32 0 019.37 22.62z"
        fill="none"
        stroke="currentColor"
        strokeLinejoin="round"
        strokeWidth="32"
      />
      <path
        d="M256 56v120a32 32 0 0032 32h120M176 288h160M176 368h160"
        fill="none"
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="32"
      />
    </svg>
  );
};

export default DocumentIcon;

