import React, { useState, useEffect } from "react";

function getClassName(className, isActive) {
  if (!isActive) {
    return className;
  }

  return `${className} active`;
}

const AnimatedFilter = (props) => {
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
        width="100"
        height="100"
        viewBox="0 0 100 100"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        {...props}
      >
        <g clipPath="url(#clip0_4_2)">
          <path
            d="M61.4379 0C61.4379 0.176471 61.4379 0.352941 61.4379 0.539216C61.4379 6.87255 56.3039 12.0049 49.9722 12.0049C43.6389 12.0049 38.5065 6.87092 38.5065 0.539216C38.5065 0.352941 38.5065 0.176471 38.5065 0L0 0V100H100V0H61.4379Z"
            fill="#00D3A9"
            className={getClassName("svg-wa-1", isActive)}
          ></path>
          <path
            d="M23.2287 53.8905H29.8562L34.3807 77.9869H34.4624L39.9003 53.8905H46.0996L51.0114 78.2826H51.1078L56.2647 53.8905H62.7647L54.3186 89.2973H47.7418L42.8709 65.2009H42.7434L37.5294 89.2973H30.83L23.2287 53.8905ZM70.2369 53.8905H80.6846L91.0604 89.2973H84.2238L81.9673 81.4182H70.0653L68.3235 89.2973H61.665L70.2369 53.8905ZM74.214 62.6176L71.3235 75.6078H80.3202L77.0016 62.6176H74.214Z"
            fill="white"
            className={getClassName("svg-wa-2", isActive)}
          ></path>
        </g>
        <defs>
          <clipPath id="clip0_4_2">
            <rect
              width="100"
              height="100"
              fill="white"
              className={getClassName("svg-wa-3", isActive)}
            ></rect>
          </clipPath>
        </defs>
      </svg>
    </div>
  );
};

export default AnimatedFilter;
