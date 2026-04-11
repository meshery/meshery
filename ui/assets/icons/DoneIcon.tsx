import React from 'react';

export const DoneIcon = ({ width, height, fill, style = {} }) => {

  return (
    <svg xmlns="http://www.w3.org/2000/svg" width={width} height={height} style={style} viewBox="0 0 88 88" fill="none">
      <g clip-path="url(#clip0_19514_2274)">
        <path d="M33.9997 58L19.9997 44L15.333 48.6666L33.9997 67.3333L73.9997 27.3333L69.333 22.6666L33.9997 58Z" fill={fill} />
      </g>
      <rect x="2" y="2" width="84" height="84" rx="42" stroke={fill} stroke-width="4" />
      <defs>
        <clipPath id="clip0_19514_2274">
          <rect x="4" y="4" width="80" height="80" rx="40" fill="white" />
        </clipPath>
      </defs>
    </svg>
  )
}


export default DoneIcon