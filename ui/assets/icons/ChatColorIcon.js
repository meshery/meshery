import React from "react";

const ChatColorIcon = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 -1 104 106"
      width={props.width ? props.width : "24px"}
      height={props.height ? props.height : "24px"}
      onClick={props.onClick}
      className={props.className}
      style={{ ...props.style }}
    >
      <defs></defs>
      <path fill="#231f20" d="M51.87,0C23.71,0,0,22.83,0,51c0,.91,0,52.81,0,52.81l51.86-.05c28.16,0,51-23.71,51-51.87S80,0,51.87,0Z"/>
      <path fill="#fff9ae" d="M52.37,19.74A31.62,31.62,0,0,0,24.58,66.41l-5.72,18.4L39.4,80.17a31.61,31.61,0,1,0,13-60.43Z"/>
      <path fill="#00aeef" d="M77.45,32.12a31.6,31.6,0,0,1-38.05,48L18.86,84.82l20.91-2.47A31.6,31.6,0,0,0,77.45,32.12Z"/>
      <path fill="#00a94f" d="M71.63,26.29A31.6,31.6,0,0,1,38.8,78L18.86,84.82,39.4,80.17A31.6,31.6,0,0,0,71.63,26.29Z"/>
      <path fill="#f15d22" d="M26.47,67.11a31.61,31.61,0,0,1,51-35A31.61,31.61,0,0,0,24.58,66.41l-5.72,18.4Z"/>
      <path fill="#e31b23" d="M24.58,66.41A31.61,31.61,0,0,1,71.63,26.29a31.61,31.61,0,0,0-49,39.63l-3.76,18.9Z"/>
    </svg>
  );
};

export default ChatColorIcon;
