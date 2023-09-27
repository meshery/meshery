import React from "react";

const CrossCircleIcon = (props) => {
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
      <g clip-path="url(#clip0_19418_2922)">
        <path d="M12 2C6.47 2 2 6.47 2 12C2 17.53 6.47 22 12 22C17.53 22 22 17.53 22 12C22 6.47 17.53 2 12 2ZM17 15.59L15.59 17L12 13.41L8.41 17L7 15.59L10.59 12L7 8.41L8.41 7L12 10.59L15.59 7L17 8.41L13.41 12L17 15.59Z" fill={props.fill} fill-opacity="0.54"/>
      </g>
      <defs>
        <clipPath id="clip0_19418_2922">
          <rect width="24" height="24" fill={props.stroke || "#fff"}/>
        </clipPath>
      </defs>
    </svg>
  );
};

export default CrossCircleIcon;
