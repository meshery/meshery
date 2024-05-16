import React from 'react';

const DetailsIcon = ({ width = '24px', height = '24px', fill = '#fff', ...props }) => (
  <svg
    width={width}
    height={height}
    fill={fill}
    {...props}
    viewBox="0 0 24 24"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g clip-path="url(#clip0_31178_5015)">
      <path d="M3 13H5V11H3V13ZM3 17H5V15H3V17ZM3 9H5V7H3V9ZM7 13H21V11H7V13ZM7 17H21V15H7V17ZM7 7V9H21V7H7Z" />
    </g>
    <defs>
      <clipPath id="clip0_31178_5015">
        <rect width="24" height="24" fill="white" />
      </clipPath>
    </defs>
  </svg>
);

export default DetailsIcon;
