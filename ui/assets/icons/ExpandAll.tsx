import React from 'react';

const ExpandAllIcon = ({ width = 24, height = 24, fill = "currentColor", style = {} }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 36 35"
    style={style}
    height={height}
    width={width}
    fill={fill}
  >
    <g id="Layer_2" data-name="Layer 2">
      <g id="Layer_1-2" data-name="Layer 1">
        <path d="M34,0H12a2,2,0,0,0-2,2V9h2V2H34V24H26v2h8a2,2,0,0,0,2-2V2A2,2,0,0,0,34,0Z" />
        <path d="M26,24V11a2,2,0,0,0-2-2H2a2,2,0,0,0-2,2V33a2,2,0,0,0,2,2H24a2,2,0,0,0,2-2V24Zm-2,9H2V11H24Z" />
        <polygon points="12 21 6 21 6 23 12 23 12 29 14 29 14 23 20 23 20 21 14 21 14 15 12 15 12 21" />
      </g>
    </g>
  </svg>
);

export default ExpandAllIcon;
