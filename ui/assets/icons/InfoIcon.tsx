import React from "react";
const InfoIcon = ({ fill = 'currentColor', height = 24, width = 24, style = {} }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width={width} height={height} style={style} viewBox="0 0 20 20" fill="none">
      <g clip-path="url(#clip0_19460_1750)">
        <path
          d="M10.0003 1.66666C5.40033 1.66666 1.66699 5.39999 1.66699 9.99999C1.66699 14.6 5.40033 18.3333 10.0003 18.3333C14.6003 18.3333 18.3337 14.6 18.3337 9.99999C18.3337 5.39999 14.6003 1.66666 10.0003 1.66666ZM10.8337 14.1667H9.16699V9.16666H10.8337V14.1667ZM10.8337 7.49999H9.16699V5.83332H10.8337V7.49999Z"
          fill={fill}
        />
      </g>
      <defs>
        <clipPath id="clip0_19460_1750">
          <rect width={width} height={height} fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
};

export default InfoIcon