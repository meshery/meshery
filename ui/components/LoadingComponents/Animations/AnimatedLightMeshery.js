import React, { useState, useEffect } from 'react';
function getClassName(className, isActive) {
  if (!isActive) {
    return className;
  }

  return `${className} active`;
}

const AnimatedLightMeshery = (props) => {
  const [isActive, setIsActive] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setIsActive(false);
    }, 100);
  }, []);

  useEffect(() => {
    setTimeout(() => {
      setIsActive(!isActive);
    }, 4000);
  }, [isActive]);

  return (
    <div>
      <svg
        id="Layer_1"
        data-name="Layer 1"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 278.24 263.92"
        width="278.239990234375"
        height="263.9200134277344"
        {...props}
      >
        <title>meshery-logo-light-text</title>
        <polygon
          points="135.15 48.92 135.15 98.49 178.12 73.62 135.15 48.92"
          className={getClassName('svg-meshery-1', isActive)}
        ></polygon>
        <polygon
          points="135.15 108.86 135.15 158.67 178.51 133.88 135.15 108.86"
          className={getClassName('svg-meshery-2', isActive)}
        ></polygon>
        <polygon
          points="128.97 98.15 128.97 49.34 86.53 73.66 128.97 98.15"
          className={getClassName('svg-meshery-3', isActive)}
        ></polygon>
        <path
          d="M43.85,158.5a104.33,104.33,0,0,0,32.92,33.12V139.44Z"
          className={getClassName('svg-meshery-4', isActive)}
        ></path>
        <polygon
          points="128.97 158.44 128.97 109.22 86.43 133.85 128.97 158.44"
          className={getClassName('svg-meshery-5', isActive)}
        ></polygon>
        <polygon
          points="82.95 188.51 126.03 163.88 82.95 138.98 82.95 188.51"
          className={getClassName('svg-meshery-6', isActive)}
        ></polygon>
        <polygon
          points="181.46 188.78 181.46 139.32 138.44 163.91 181.46 188.78"
          className={getClassName('svg-meshery-7', isActive)}
        ></polygon>
        <path
          d="M223.6,152.75a103.12,103.12,0,0,0,12.16-44.48L190.93,133.9Z"
          className={getClassName('svg-meshery-8', isActive)}
        ></path>
        <polygon
          points="187.64 128.66 230.87 103.94 187.64 79.09 187.64 128.66"
          className={getClassName('svg-meshery-9', isActive)}
        ></polygon>
        <polygon
          points="181.46 68.4 181.46 19.25 138.64 43.79 181.46 68.4"
          className={getClassName('svg-meshery-10', isActive)}
        ></polygon>
        <polygon
          points="181.46 128.44 181.46 78.83 138.54 103.68 181.46 128.44"
          className={getClassName('svg-meshery-11', isActive)}
        ></polygon>
        <polygon
          points="132.06 207.57 132.07 207.57 132.05 207.57 132.06 207.57"
          className={getClassName('svg-meshery-12', isActive)}
        ></polygon>
        <polygon
          points="82.95 78.74 82.95 128.71 126.18 103.69 82.95 78.74"
          className={getClassName('svg-meshery-13', isActive)}
        ></polygon>
        <path
          d="M129,0A103.26,103.26,0,0,0,83.41,12L129,38.23Z"
          className={getClassName('svg-meshery-14', isActive)}
        ></path>
        <polygon
          points="82.95 18.91 82.95 68.58 126.22 43.78 82.95 18.91"
          className={getClassName('svg-meshery-15', isActive)}
        ></polygon>
        <path
          d="M76.77,15.89A104.42,104.42,0,0,0,43.84,49L76.77,68Z"
          className={getClassName('svg-meshery-16', isActive)}
        ></path>
        <path
          d="M181.17,12.28A103.4,103.4,0,0,0,135.15,0V38.66Z"
          className={getClassName('svg-meshery-17', isActive)}
        ></path>
        <path
          d="M28.32,107.4a103.21,103.21,0,0,0,12.42,45.75l33.34-19.3Z"
          className={getClassName('svg-meshery-18', isActive)}
        ></path>
        <path
          d="M40.73,54.37A103.3,103.3,0,0,0,28.33,99.9l45.8-26.25Z"
          className={getClassName('svg-meshery-19', isActive)}
        ></path>
        <path
          d="M83.32,195.43A103.35,103.35,0,0,0,129,207.52V169.33Z"
          className={getClassName('svg-meshery-20', isActive)}
        ></path>
        <path
          d="M220.36,49.16a104.43,104.43,0,0,0-32.72-33.09v52Z"
          className={getClassName('svg-meshery-21', isActive)}
        ></path>
        <path
          d="M187.64,191.44a104.29,104.29,0,0,0,32.87-33.32l-32.87-19Z"
          className={getClassName('svg-meshery-22', isActive)}
        ></path>
        <path
          d="M235.78,99.63a103.28,103.28,0,0,0-12.32-45.12l-33,19.09Z"
          className={getClassName('svg-meshery-23', isActive)}
        ></path>
        <path
          d="M135.15,207.52a103.14,103.14,0,0,0,45.54-12l-45.54-26.32Z"
          className={getClassName('svg-meshery-24', isActive)}
        ></path>
        <polygon
          points="76.77 128.26 76.77 79.26 34.2 103.66 76.77 128.26"
          className={getClassName('svg-meshery-25', isActive)}
        ></polygon>
        <polygon
          fill="#FFF"
          points="172.5 235.15 172.5 236.6 172.5 243.43 172.5 251.76 172.5 263.92 204.87 263.92 204.87 257.09 179.33 257.09 179.33 251.76 179.33 250.26 198.95 250.26 198.95 243.43 179.33 243.43 179.33 236.6 204.87 236.6 204.87 229.77 172.5 229.77 172.5 235.15"
        ></polygon>
        <polygon
          fill="#FFF"
          points="47.96 235.15 47.96 236.6 47.96 243.43 47.96 251.76 47.96 263.92 80.34 263.92 80.34 257.09 54.79 257.09 54.79 251.76 54.79 250.26 74.42 250.26 74.42 243.43 54.79 243.43 54.79 236.6 80.34 236.6 80.34 229.77 47.96 229.77 47.96 235.15"
        ></polygon>
        <path
          fill="#FFF"
          d="M125.15,255a8.53,8.53,0,0,1-.7,3.46,8.73,8.73,0,0,1-1.94,2.85,9.52,9.52,0,0,1-2.84,1.91,8.78,8.78,0,0,1-3.49.71h-32v-7.2h32a1.71,1.71,0,0,0,1.73-1.73v-2.78a1.72,1.72,0,0,0-1.73-1.73h-23a8.53,8.53,0,0,1-3.46-.7,9.15,9.15,0,0,1-4.81-4.8,8.61,8.61,0,0,1-.7-3.47v-2.77a8.77,8.77,0,0,1,.7-3.49,9.13,9.13,0,0,1,1.94-2.84,9,9,0,0,1,2.87-1.94,8.53,8.53,0,0,1,3.46-.7h32V237h-32a1.64,1.64,0,0,0-1.23.5,1.67,1.67,0,0,0-.5,1.23v2.77a1.67,1.67,0,0,0,.5,1.23,1.64,1.64,0,0,0,1.23.5h23a8.78,8.78,0,0,1,3.49.71,9.25,9.25,0,0,1,2.84,1.94,9.13,9.13,0,0,1,1.94,2.84,8.73,8.73,0,0,1,.7,3.48Z"
        ></path>
        <path
          fill="#FFF"
          d="M168.23,263.92h-6.51V250.48h-25.4v13.44h-6.51V243.24h31.91V229.77h6.51ZM136.32,240h-6.51V229.77h6.51Z"
        ></path>
        <polygon
          fill="#FFF"
          points="261.17 241.72 252.63 229.77 244.09 229.77 257.75 250.26 257.75 263.92 264.58 263.92 264.58 250.26 278.24 229.77 269.81 229.77 261.17 241.72"
        ></polygon>
        <path
          fill="#FFF"
          d="M216.49,236.28H234v7.1H219.9v6.51h3.23l12.39,14H244l-12.39-14h4.84a5.3,5.3,0,0,0,5.3-5.3v-9.52a5.3,5.3,0,0,0-5.3-5.3H209.8v34.15h6.83Z"
        ></path>
        <path
          fill="#FFF"
          d="M42.71,263.92H35.47V239a1.83,1.83,0,0,0-.16-.77,1.9,1.9,0,0,0-.43-.64,2,2,0,0,0-.64-.43,1.83,1.83,0,0,0-.77-.16H0v-7.24H33.47a8.81,8.81,0,0,1,3.57.73,9.49,9.49,0,0,1,3,2,9.38,9.38,0,0,1,2,3,8.81,8.81,0,0,1,.73,3.57Zm-35.47,0H0V237l7.24-2Zm17.71,0H17.71V240.8H25Z"
        ></path>
      </svg>
    </div>
  );
};

export default AnimatedLightMeshery;
