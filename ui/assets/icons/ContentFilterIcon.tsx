import React from "react";

const ContentFilterIcon = (props) => {
  return (
    <svg
      id="content-filter-icon-svg"
      version="1.1"
      viewBox="0 0 24 24"
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

      <path fill-rule="evenodd" clip-rule="evenodd" d="M10 18H14V16H10V18ZM3 6V8H21V6H3ZM6 13H18V11H6V13Z" fill={props.fill}/>
    </svg>
  );
};

export default ContentFilterIcon;

