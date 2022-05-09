import React, { useState, useEffect } from "react";
import "./mesheryAnimation.css";


function getClassName(className, isActive) {
  if (!isActive) {
    return className;
  }

  return `${className} active`;
}

export const MesheryAnimation = (props) => {
  const [isActive, setIsActive] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setIsActive(false);
    }, 100);
  }, [])

  useEffect(() => {
    setTimeout(() => {
      setIsActive(!isActive);
    }, 2000);
  }, [isActive])
  return (
    <div>
      <div>
      <svg width="135" height="136" viewBox="0 0 135 136" fill="none" xmlns="http://www.w3.org/2000/svg" {...props}>
<g clip-path="url(#clip0_5308_25299)">
<path d="M69.49 31.82V64.07L97.44 47.89L69.49 31.82Z" className={getClassName("svg-elem-1", isActive)} fill="#00D3A9"/>
<path d="M69.49 70.81V103.22L97.7 87.09L69.49 70.81Z" className={getClassName("svg-elem-2", isActive)} fill="#00D3A9"/>
<path d="M65.47 63.85V32.09L37.87 47.92L65.47 63.85Z" className={getClassName("svg-elem-3", isActive)} fill="#00B39F"/>
<path d="M10.1 103.1C15.5146 111.81 22.835 119.179 31.51 124.65V90.71L10.1 103.1Z" className={getClassName("svg-elem-4", isActive)} fill="#00B39F"/>
<path d="M65.4701 103.06V71.05L37.8 87.07L65.4701 103.06Z"className={getClassName("svg-elem-5", isActive)} fill="#00B39F"/>
<path d="M35.54 122.63L63.56 106.61L35.54 90.41V122.63Z"className={getClassName("svg-elem-6", isActive)} fill="#00D3A9"/>
<path d="M99.61 122.8V90.63L71.63 106.63L99.61 122.8Z"className={getClassName("svg-elem-7", isActive)} fill="#00B39F"/>
<path d="M127 99.37C131.783 90.4424 134.487 80.5494 134.91 70.43L105.78 87.11L127 99.37Z"className={getClassName("svg-elem-8", isActive)} fill="#00B39F"/>
<path d="M103.64 83.69L131.76 67.61L103.64 51.45V83.69Z"className={getClassName("svg-elem-9", isActive)} fill="#00D3A9"/>
<path d="M99.61 44.5V12.52L71.76 28.49L99.61 44.5Z"className={getClassName("svg-elem-10", isActive)} fill="#00B39F"/>
<path d="M99.61 83.55V51.28L71.7 67.44L99.61 83.55Z" className={getClassName("svg-elem-11", isActive)}fill="#00B39F"/>
<path d="M35.54 51.22V83.73L63.66 67.45L35.54 51.22Z"className={getClassName("svg-elem-12", isActive)} fill="#00D3A9"/>
<path d="M65.47 0C55.1186 0.291396 44.9747 2.97112 35.83 7.83L65.47 24.83V0Z"className={getClassName("svg-elem-13", isActive)} fill="#00B39F"/>
<path d="M35.54 12.3V44.62L63.68 28.48L35.54 12.3Z"className={getClassName("svg-elem-14", isActive)} fill="#00D3A9"/>
<path d="M31.51 10.34C22.8372 15.814 15.5173 23.1818 10.1 31.89L31.51 44.25V10.34Z"className={getClassName("svg-elem-15", isActive)} fill="#00B39F"/>
<path d="M99.43 8C90.2092 3.03814 79.9568 0.2987 69.49 0V25.15L99.43 8Z" className={getClassName("svg-elem-16", isActive)}fill="#00D3A9"/>
<path d="M0 69.87C0.34952 80.281 3.11209 90.4686 8.07 99.63L29.76 87.07L0 69.87Z"className={getClassName("svg-elem-17", isActive)} fill="#00D3A9"/>
<path d="M8.07 35.37C3.12771 44.4908 0.36551 54.6326 0 65L29.79 47.91L8.07 35.37Z"className={getClassName("svg-elem-18", isActive)} fill="#00D3A9"/>
<path d="M35.78 127.13C44.9355 132.013 55.0981 134.706 65.47 135V110.15L35.78 127.13Z"className={getClassName("svg-elem-19", isActive)} fill="#00B39F"/>
<path d="M124.92 32C119.536 23.317 112.262 15.961 103.64 10.48V44.3L124.92 32Z"className={getClassName("svg-elem-20", isActive)} fill="#00D3A9"/>
<path d="M103.64 124.54C112.307 119.02 119.609 111.608 125 102.86L103.64 90.52V124.54Z"className={getClassName("svg-elem-21", isActive)} fill="#00D3A9"/>
<path d="M135 64.81C134.617 54.543 131.88 44.5013 127 35.46L105.49 47.88L135 64.81Z"className={getClassName("svg-elem-22", isActive)} fill="#00B39F"/>
<path d="M69.49 135C79.8383 134.71 89.9793 132.03 99.12 127.17L69.49 110V135Z"className={getClassName("svg-elem-23", isActive)} fill="#00D3A9"/>
<path d="M31.51 83.44V51.56L3.82996 67.43L31.51 83.44Z"className={getClassName("svg-elem-24", isActive)} fill="#00B39F"/>
</g>
<defs>
<clipPath id="clip0_5308_25299">
<rect width="134.95" height="135.02" className={getClassName("svg-elem-25", isActive)} fill="white"/>
</clipPath>
</defs>
</svg>

      {/* <svg
    id="Layer_1"
    data-name="Layer 1"
    xmlns="http://www.w3.org/2000/svg"
    width={width}
    height={height}
    {...props}
  >
    <defs>
      <style>{".cls-1{fill:#00d3a9}.cls-2{fill:#00b39f}"}</style>
    </defs>
    <title>{"meshery-logo-light"}</title>
    <path
      className={getClassName("svg-elem-1", isActive)}
      d="M69.49 31.82v32.25l27.95-16.18-27.95-16.07z"
    />
    <path
      className={getClassName("svg-elem-2", isActive)}
      d="M69.49 70.81v32.41L97.7 87.09 69.49 70.81z"
    />
    <path
      className={getClassName("svg-elem-3", isActive)}
      d="M65.47 63.85V32.09l-27.6 15.83 27.6 15.93z"
    />
    <path
      className={getClassName("svg-elem-4", isActive)}
      d="M10.1 103.1a67.79 67.79 0 0 0 21.41 21.55V90.71Z"
    />
    <path
      className={getClassName(" svg-elem-5", isActive)}
      d="M65.47 103.06V71.05L37.8 87.07l27.67 15.99z"
    />
    <path
      className={getClassName(" svg-elem-6", isActive)}
      d="m35.54 122.63 28.02-16.02-28.02-16.2v32.22z"
    />
    <path
      className={getClassName(" svg-elem-7", isActive)}
      d="M99.61 122.8V90.63l-27.98 16 27.98 16.17z"
    />
    <path
      className={getClassName(" svg-elem-8", isActive)}
      d="M127 99.37a67.22 67.22 0 0 0 7.91-28.94l-29.13 16.68Z"
    />
    <path
      className={getClassName(" svg-elem-9", isActive)}
      d="m103.64 83.69 28.12-16.08-28.12-16.16v32.24z"
    />
    <path
      className={getClassName(" svg-elem-10", isActive)}
      d="M99.61 44.5V12.52L71.76 28.49 99.61 44.5z"
    />
    <path
      className={getClassName(" svg-elem-11", isActive)}
      d="M99.61 83.55V51.28L71.7 67.44l27.91 16.11z"
    />
    <path className={getClassName(" svg-elem-12", isActive)} d="M67.48 135.02h.01-.01z" />
    <path
      className={getClassName(" svg-elem-13", isActive)}
      d="M35.54 51.22v32.51l28.12-16.28-28.12-16.23z"
    />
    <path
      className={getClassName(" svg-elem-14", isActive)}
      d="M65.47 0a67.2 67.2 0 0 0-29.64 7.83l29.64 17Z"
    />
    <path
      className={getClassName(" svg-elem-15", isActive)}
      d="M35.54 12.3v32.32l28.14-16.14L35.54 12.3z"
    />
    <path
      className={getClassName(" svg-elem-16", isActive)}
      d="M31.51 10.34A67.89 67.89 0 0 0 10.1 31.89l21.41 12.36Z"
    />
    <path
      className={getClassName(" svg-elem-17", isActive)}
      d="M99.43 8a67.23 67.23 0 0 0-29.94-8v25.15Z"
    />
    <path
      className={getClassName(" svg-elem-18", isActive)}
      d="M0 69.87a67.27 67.27 0 0 0 8.07 29.76l21.69-12.56Z"
    />
    <path
      className={getClassName(" svg-elem-19", isActive)}
      d="M8.07 35.37A67.16 67.16 0 0 0 0 65l29.79-17.09Z"
    />
    <path
      className={getClassName(" svg-elem-20", isActive)}
      d="M35.78 127.13A67.13 67.13 0 0 0 65.47 135v-24.85Z"
    />
    <path
      className={getClassName(" svg-elem-21", isActive)}
      d="M124.92 32a67.9 67.9 0 0 0-21.28-21.52V44.3Z"
    />
    <path
      className={getClassName(" svg-elem-22", isActive)}
      d="M103.64 124.54A68 68 0 0 0 125 102.86l-21.36-12.34Z"
    />
    <path
      className={getClassName(" svg-elem-23", isActive)}
      d="M135 64.81a67.06 67.06 0 0 0-8-29.35l-21.51 12.42Z"
    />
    <path
      className={getClassName("svg-elem-24", isActive)}
      d="M69.49 135a67.12 67.12 0 0 0 29.63-7.83L69.49 110Z"
    />
    <path
      className={getClassName(" svg-elem-25", isActive)}
      d="M31.51 83.44V51.56L3.83 67.43l27.68 16.01z"
    />
  </svg> */}
      </div>
    </div>
  );
};