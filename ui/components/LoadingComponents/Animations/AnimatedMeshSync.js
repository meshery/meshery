import React, { useState, useEffect } from 'react';

function getClassName(className, isActive) {
  if (!isActive) {
    return className;
  }

  return `${className} active`;
}

export default function AnimatedMeshSync(props) {
  const [isActive, setIsActive] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setIsActive(false);
    }, 100);
  }, []);

  useEffect(() => {
    setTimeout(() => {
      setIsActive(!isActive);
    }, 2000);
  }, [isActive]);

  return (
    <div>
      <svg
        id="Layer_1"
        data-name="Layer 1"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 361 359"
        width="361"
        height="359"
        {...props}
      >
        <title>meshsync</title>
        <path
          d="M63.5,113.5a133.13,133.13,0,0,1,94-67.5l-8-46A181.53,181.53,0,0,0,22,91Z"
          className={getClassName('svg-sync-5', isActive)}
        ></path>
        <path
          d="M65.5,247a133.42,133.42,0,0,1-19-69A130.14,130.14,0,0,1,55,131.5l-44-16A184.54,184.54,0,0,0,0,178a182.54,182.54,0,0,0,25.5,93Z"
          className={getClassName('svg-sync-4', isActive)}
        ></path>
        <path
          d="M182,312h-1.5C139,312,102,293,77,263L41,293a181.18,181.18,0,0,0,139.5,66h2Z"
          className={getClassName('svg-sync-3', isActive)}
        ></path>
        <path
          d="M297,244a134.07,134.07,0,0,1-95,66l7.5,46.5a180.18,180.18,0,0,0,128-89.5Z"
          className={getClassName('svg-sync-2', isActive)}
        ></path>
        <path
          d="M296,110.5A133.46,133.46,0,0,1,314,178a128,128,0,0,1-9,48l44,17a184.66,184.66,0,0,0,12-64.5,180.8,180.8,0,0,0-24.5-91Z"
          className={getClassName('svg-sync-1', isActive)}
        ></path>
        <rect
          x="251.3"
          y="17.85"
          width="44"
          height="67.8"
          style={{ fill: 'none' }}
          className={getClassName('svg-sync-8', isActive)}
        ></rect>
        <polygon
          points="261.4 85.65 251.3 75.55 275.1 51.75 251.3 27.95 261.4 17.85 295.3 51.75 261.4 85.65"
          className={getClassName('svg-sync-9', isActive)}
        ></polygon>
        <rect
          x="183.5"
          y="17.85"
          width="44"
          height="67.8"
          style={{ fill: 'none' }}
          className={getClassName('svg-sync-6', isActive)}
        ></rect>
        <polygon
          points="193.6 85.65 183.5 75.55 207.3 51.75 183.5 27.95 193.6 17.85 227.5 51.75 193.6 85.65"
          className={getClassName('svg-sync-7', isActive)}
        ></polygon>
      </svg>
    </div>
  );
}
