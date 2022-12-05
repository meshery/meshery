import React from "react";

const SquarePollVerticalIcon = (props) => {
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
      <path d="M384 31.1H64c-35.35 0-64 28.65-64 63.1v320c0 35.35 28.65 64 64 64h320c35.35 0 64-28.65 64-64v-320C448 60.65 419.3 31.1 384 31.1zM160 368C160 376.9 152.9 384 144 384h-32C103.1 384 96 376.9 96 368v-128C96 231.1 103.1 224 112 224h32C152.9 224 160 231.1 160 240V368zM256 368c0 8.875-7.125 16-16 16h-32C199.1 384 192 376.9 192 368v-224c0-8.875 7.125-16 16-16h32c8.875 0 16 7.125 16 16V368zM352 368c0 8.875-7.125 16-16 16h-32c-8.875 0-16-7.125-16-16v-64C288 295.1 295.1 288 304 288h32C344.9 288 352 295.1 352 304V368z" />
    </svg>
  );
};

export default SquarePollVerticalIcon;

