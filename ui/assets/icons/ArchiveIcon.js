
import React from "react";
const ArchiveIcon = ({ height, width, fill, innerFill = "#fff", style = {} }) => {
  return (
    <svg
      style={style}
      height={height}
      width={width}
      fill={fill}
      viewBox="0 0 20 20"
    >
      <g clip-path="url(#clip0_19418_1412)">
        <path d="M17.1167 4.35833L15.9583 2.95833C15.7333 2.675 15.3917 2.5 15 2.5H5C4.60833 2.5 4.26667 2.675 4.03333 2.95833L2.88333 4.35833C2.64167 4.64167 2.5 5.01667 2.5 5.41667V15.8333C2.5 16.75 3.25 17.5 4.16667 17.5H15.8333C16.75 17.5 17.5 16.75 17.5 15.8333V5.41667C17.5 5.01667 17.3583 4.64167 17.1167 4.35833ZM10 14.5833L5.41667 10H8.33333V8.33333H11.6667V10H14.5833L10 14.5833ZM4.26667 4.16667L4.94167 3.33333H14.9417L15.725 4.16667H4.26667Z" fill={fill} />
      </g>
      <defs>
        <clipPath id="clip0_19418_1412">
          <rect width="20" height="20" fill={innerFill} />
        </clipPath>
      </defs>
    </svg>
  );
};

export default ArchiveIcon;