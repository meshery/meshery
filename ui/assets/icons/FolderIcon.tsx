import React from "react";

const FolderIcon = ({ width="1.5rem", height="1.5rem", fill, style = {} }) => (
  <svg
    width={width}
    height={height}
    viewBox="0 0 31 30"
    fill="none"
    style={style}
    xmlns="http://www.w3.org/2000/svg"
  >
    <g clipPath="url(#a)">
      <path
        d="M28.625 3.75h-11.25V1.875A1.874 1.874 0 0 0 15.5 0H2.375C1.339 0 .5.84.5 1.875v26.25c0 .473.18.9.469 1.23v.176h.177c.33.289.756.469 1.229.469h26.25c.473 0 .9-.18 1.23-.469h.176v-.177a1.86 1.86 0 0 0 .469-1.229v-22.5a1.874 1.874 0 0 0-1.875-1.875Zm-5.156 17.813H7.53a.468.468 0 1 1 0-.938h15.94a.468.468 0 1 1 0 .938Zm0-2.813H7.53a.468.468 0 1 1 0-.938h15.94a.468.468 0 1 1 0 .938ZM7.062 15.469c0-.26.21-.469.47-.469h7.5a.468.468 0 1 1 0 .938h-7.5a.468.468 0 0 1-.47-.47Zm22.5-8.906H1.438V1.875c0-.518.42-.938.938-.938H15.5c.517 0 .938.42.938.938v2.813h12.187c.517 0 .938.42.938.937v.938Z"
        fill={fill ? fill : '#3C494F'}
      />
    </g>
    <defs>
      <clipPath id="a">
        <path fill={fill ? fill : '#fff'} d="M.5 0h30v30H.5z" />
      </clipPath>
    </defs>
  </svg>
);

export default FolderIcon;
