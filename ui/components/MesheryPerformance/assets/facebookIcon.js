import React from "react";
import { iconLarge } from "../../../css/icons.styles"

const FacebookIcon = ({ width = 40, height = 40 }) => (
  <svg width={width} height={height} viewBox="0 0 28 29" style={iconLarge} fill="none" xmlns="http://www.w3.org/2000/svg">
    <g filter="url(#filter0_d_6961_29756)">
      <rect x="3.87598" y="3" width="20.2472" height="20.2472" rx="4.82076" fill="#252E31"/>
    </g>
    <path d="M16.5443 12.8065H15.1772V17.4103H13.5977V12.8065H12.7277V12.0453L13.5977 11.6207V11.196C13.5977 10.5366 13.76 10.055 14.0845 9.7512C14.409 9.44739 14.9286 9.29549 15.6432 9.29549C16.1887 9.29549 16.6738 9.37662 17.0984 9.53888L16.6945 10.6989C16.3769 10.5988 16.0834 10.5487 15.8141 10.5487C15.5897 10.5487 15.4275 10.616 15.3274 10.7507C15.2272 10.8819 15.1772 11.051 15.1772 11.2582V11.6207H16.5443V12.8065Z" fill="white"/>
    <defs>
      <filter id="filter0_d_6961_29756" x="0.0193691" y="0.107544" width="27.9603" height="27.9604" filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
        <feFlood floodOpacity="0" result="BackgroundImageFix"/>
        <feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
        <feOffset dy="0.964152"/>
        <feGaussianBlur stdDeviation="1.9283"/>
        <feColorMatrix type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.08 0"/>
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_6961_29756"/>
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_6961_29756" result="shape"/>
      </filter>
    </defs>
  </svg>
);

export default FacebookIcon;