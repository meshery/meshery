import React from 'react';

export const OCIImageIcon = ({ height = 24, width = 24 }) => {
  const st0 = {
    fill: 'none',
  };

  const st1 = {
    // fill: '#262261', // orginal branding color has very low contrast
    fill: '#326CE5',
    // kubernetes blue
  };

  return (
    <svg
      version="1.1"
      id="Layer_1"
      xmlns="http://www.w3.org/2000/svg"
      x="0px"
      y="0px"
      width={width}
      height={height}
      viewBox="0 0 94.6 94.6"
    >
      <g>
        <polygon style={st0} points="82.8,57.8 83.3,77.3 88.8,76.5 88.2,57 	" />
        <path
          style={st0}
          d="M47.3,29.4c4.5,0,8.1-3.6,8.1-8.1s-3.6-8.1-8.1-8.1c-4.5,0-8.1,3.6-8.1,8.1S42.8,29.4,47.3,29.4z M47.3,15.5
		c3.2,0,5.9,2.6,5.9,5.9s-2.6,5.9-5.9,5.9c-3.2,0-5.9-2.6-5.9-5.9S44.1,15.5,47.3,15.5z"
        />
        <polygon style={st0} points="69.1,59.7 69.6,79.2 75.1,78.4 74.5,59 	" />

        <polygon style={st0} points="55.4,61.6 55.9,81.1 61.4,80.4 60.8,60.9 	" />
        <polygon style={st0} points="41.7,63.6 42.3,83.1 47.7,82.3 47.1,62.8 	" />
        <polygon style={st0} points="28,65.5 28.6,85 34,84.2 33.4,64.7 	" />
        <path
          style={st1}
          d="M19.1,60.4l0.8,31.6l74-10.4L93,50.3L19.1,60.4z M28.6,85L28,65.5l5.4-0.8L34,84.2L28.6,85z M42.3,83.1
		l-0.6-19.5l5.4-0.8l0.6,19.5L42.3,83.1z M55.9,81.1l-0.6-19.5l5.4-0.8l0.6,19.5L55.9,81.1z M69.6,79.2l-0.6-19.5l5.4-0.8l0.6,19.5
		L69.6,79.2z M83.3,77.3l-0.6-19.5l5.4-0.8l0.6,19.5L83.3,77.3z"
        />
        <polygon style={st1} points="0.5,50.2 1.4,79.4 14,88.4 13.3,59.4 	" />
        <polygon
          style={st1}
          points="73.6,33.9 73.6,35.1 73.6,41.1 67.6,41.1 27,41.1 22.4,41.1 1.8,44 16.9,54.8 90.7,44.7 75.1,33.7
		"
        />
        <path
          style={st1}
          d="M47.3,27.2c3.2,0,5.9-2.6,5.9-5.9s-2.6-5.9-5.9-5.9c-3.2,0-5.9,2.6-5.9,5.9S44.1,27.2,47.3,27.2z"
        />
        <path
          style={st1}
          d="M27,35.1h40.6V7.6H54.9l-1.6-6h-12l-1.6,6h-3.3v-3h-5.6v3H27V35.1z M47.3,13.3c4.5,0,8.1,3.6,8.1,8.1
		s-3.6,8.1-8.1,8.1c-4.5,0-8.1-3.6-8.1-8.1S42.8,13.3,47.3,13.3z"
        />
      </g>
    </svg>
  );
};

export default OCIImageIcon;
