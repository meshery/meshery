import React, { useState, useEffect } from "react";
import  "./SpinnerLoading.styles.css";
import classNames from "classnames";

// eslint-disable-next-line react/prop-types
 const SpinnerAnimated = ({ viewBox, width, height }) => {
  const [active, setActive] = useState(false);
  // const  = useStyles();

  useEffect(() => {
    const interval1 = setInterval(() => {
      if (active) setActive(false);
    }, 5000);
    const interval2 = setInterval(() => {
      if (!active) setActive(true);
    }, 300);
    return () => clearInterval(interval1, interval2);
  });

  return (

    <div>
      <div>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          data-name="Layer 1"
          viewBox={viewBox}
          width={width}
          height={height}
          className={active}
        >
          <title>color-logo-only</title>
          <polygon
            points="69.49 31.82 69.49 64.07 97.44 47.89 69.49 31.82"
            className={
              active ? classNames("cls1", "activeSvgElem1") : classNames("cls1", "svgElem1")
            }
          ></polygon>
          <polygon
            points="69.49 70.81 69.49 103.22 97.7 87.09 69.49 70.81"
            className={
              active ? classNames("cls1", "activeSvgElem2") : classNames("cls1", "svgElem2")
            }
          ></polygon>
          <polygon
            points="65.47 63.85 65.47 32.09 37.87 47.92 65.47 63.85"
            className={
              active ? classNames("cls2", "activeSvgElem3") : classNames("cls2", "svgElem3")
            }
          ></polygon>
          <path
            d="M10.1,103.1a67.79,67.79,0,0,0,21.41,21.55V90.71Z"
            className={
              active ? classNames("cls2", "activeSvgElem4") : classNames("cls2", "svgElem4")
            }
          ></path>
          <polygon
            points="65.47 103.06 65.47 71.05 37.8 87.07 65.47 103.06"
            className={
              active ? classNames("cls2", "activeSvgElem5") : classNames("cls2", "svgElem5")
            }
          ></polygon>
          <polygon
            points="35.54 122.63 63.56 106.61 35.54 90.41 35.54 122.63"
            className={
              active ? classNames("cls1", "activeSvgElem6") : classNames("cls1", "svgElem6")
            }
          ></polygon>
          <polygon
            points="99.61 122.8 99.61 90.63 71.63 106.63 99.61 122.8"
            className={
              active ? classNames("cls2", "activeSvgElem7") : classNames("cls2", "svgElem7")
            }
          ></polygon>
          <path
            d="M127,99.37a67.22,67.22,0,0,0,7.91-28.94L105.78,87.11Z"
            className={
              active ? classNames("cls2", "activeSvgElem8") : classNames("cls2", "svgElem8")
            }
          ></path>
          <polygon
            points="103.64 83.69 131.76 67.61 103.64 51.45 103.64 83.69"
            className={
              active ? classNames("cls1", "activeSvgElem9") : classNames("cls1", "svgElem9")
            }
          ></polygon>
          <polygon
            points="99.61 44.5 99.61 12.52 71.76 28.49 99.61 44.5"
            className={
              active ? classNames("cls2", "activeSvgElem10") : classNames("cls2", "svgElem10")
            }
          ></polygon>
          <polygon
            points="99.61 83.55 99.61 51.28 71.7 67.44 99.61 83.55"
            className={
              active ? classNames("cls2", "activeSvgElem11") : classNames("cls2", "svgElem11")
            }
          ></polygon>
          <polygon
            points="67.48 135.02 67.49 135.02 67.48 135.02 67.48 135.02"
            className={
              active ? classNames("cls2", "activeSvgElem12") : classNames("cls2", "svgElem12")
            }
          ></polygon>
          <polygon
            points="35.54 51.22 35.54 83.73 63.66 67.45 35.54 51.22"
            className={
              active ? classNames("cls1", "activeSvgElem13") : classNames("cls1", "svgElem13")
            }
          ></polygon>
          <path
            d="M65.47,0A67.2,67.2,0,0,0,35.83,7.83l29.64,17Z"
            className={
              active ? classNames("cls2", "activeSvgElem14") : classNames("cls2", "svgElem14")
            }
          ></path>
          <polygon
            points="35.54 12.3 35.54 44.62 63.68 28.48 35.54 12.3"
            className={
              active ? classNames("cls1", "activeSvgElem15") : classNames("cls1", "svgElem15")
            }
          ></polygon>
          <path
            d="M31.51,10.34A67.89,67.89,0,0,0,10.1,31.89L31.51,44.25Z"
            className={
              active ? classNames("cls2", "activeSvgElem16") : classNames("cls2", "svgElem16")
            }
          ></path>
          <path
            d="M99.43,8A67.23,67.23,0,0,0,69.49,0V25.15Z"
            className={
              active ? classNames("cls1", "activeSvgElem17") : classNames("cls1", "svgElem17")
            }
          ></path>
          <path
            d="M0,69.87A67.27,67.27,0,0,0,8.07,99.63L29.76,87.07Z"
            className={
              active ? classNames("cls1", "activeSvgElem18") : classNames("cls1", "svgElem18")
            }
          ></path>
          <path
            d="M8.07,35.37A67.16,67.16,0,0,0,0,65L29.79,47.91Z"
            className={
              active ? classNames("cls1", "activeSvgElem19") : classNames("cls1", "svgElem19")
            }
          ></path>
          <path
            d="M35.78,127.13A67.13,67.13,0,0,0,65.47,135V110.15Z"
            className={
              active ? classNames("cls2", "activeSvgElem20") : classNames("cls2", "svgElem20")
            }
          ></path>
          <path
            d="M124.92,32a67.9,67.9,0,0,0-21.28-21.52V44.3Z"
            className={
              active ? classNames("cls1", "activeSvgElem21") : classNames("cls1", "svgElem21")
            }
          ></path>
          <path
            d="M103.64,124.54A68,68,0,0,0,125,102.86L103.64,90.52Z"
            className={
              active ? classNames("cls1", "activeSvgElem22") : classNames("cls1", "svgElem22")
            }
          ></path>
          <path
            d="M135,64.81a67.06,67.06,0,0,0-8-29.35L105.49,47.88Z"
            className={
              active ? classNames("cls2", "activeSvgElem23") : classNames("cls2", "svgElem23")
            }
          ></path>
          <path
            d="M69.49,135a67.12,67.12,0,0,0,29.63-7.83L69.49,110Z"
            className={
              active ? classNames("cls1", "activeSvgElem24") : classNames("cls1", "svgElem24")
            }
          ></path>
          <polygon
            points="31.51 83.44 31.51 51.56 3.83 67.43 31.51 83.44"
            className={
              active ? classNames("cls2", "activeSvgElem25") : classNames("cls2", "svgElem25")
            }
          ></polygon>
        </svg>
      </div>
    </div>
  );
};


export default SpinnerAnimated