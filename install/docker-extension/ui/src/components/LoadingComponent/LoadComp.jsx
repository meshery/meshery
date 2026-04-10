import React, { useState, useEffect } from "react";
import "./loadComp.css";


function getClassName(className, isActive) {
  if (!isActive) {
    return className;
  }

  return `${className} active`;
}

export const LoadComp = (props) => {
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
      <svg
    id="Layer_1"
    data-name="Layer 1"
    xmlns="http://www.w3.org/2000/svg"
    width={200.11}
    height={195.78}
    {...props}
  >
    <defs>
      <style>{".cls-1{fill:#00b39f}.cls-2{fill:#00d3a9}"}</style>
    </defs>
    <title>{"5-light-no-trim"}</title>
    <path
      id="_Path_"
      data-name="&lt;Path&gt;"
      className={getClassName("svg-el-1", isActive)}
      d="M57.83 0h36.93v17.79H57.83z"
    />
    <path
      id="_Path_2"
      data-name="&lt;Path&gt;"
      className={getClassName("svg-el-2", isActive)}
      d="M98.11 59.93v-7.38c0-8.48-7.14-15.36-15.94-15.36h-59.6v-4.85H.75v21h75.54v6.57Z"
    />
    <path
      id="_Path_3"
      data-name="&lt;Path&gt;"
      className={getClassName("svg-el-3", isActive)}
      d="M22.57 26.52v-8.73h29.22V0H.75v26.52h21.82z"
    />
    <path
      id="_Path_4"
      data-name="&lt;Path&gt;"
      className={getClassName("svg-el-4", isActive)}
      d="M76.29 65.64V76H47.07v17.78h36c8.31 0 15-6.48 15-14.48V65.64Z"
    />
    <path
      id="_Path_5"
      data-name="&lt;Path&gt;"
      className={getClassName("svg-el-5", isActive)}
      d="M20.82 66.61V76H41v17.78H15c-8.27 0-15-6.48-15-14.48V66.61Z"
    />
  </svg>
  );
};