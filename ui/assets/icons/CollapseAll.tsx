import React from 'react';

const CollapseAllIcon = ({ width = 24, height = 24, fill = "currentColor", style = {} }) => (
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
        <path fill="none" d="M26,11V24h8V2H12V9H24A2,2,0,0,1,26,11Z" />
        <path d="M34,0H12a2,2,0,0,0-2,2V9h2V2H34V24H26v2h8a2,2,0,0,0,2-2V2A2,2,0,0,0,34,0Z" />
        <path d="M26,11a2,2,0,0,0-2-2H2a2,2,0,0,0-2,2V33a2,2,0,0,0,2,2H24a2,2,0,0,0,2-2V11ZM2,11H24V33H2Z" />
        <rect x="6" y="21.5" width="14" height="2" />
      </g>
    </g>
  </svg>
);

export default CollapseAllIcon;
