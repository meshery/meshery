import React from 'react';
import { iconLarge } from '../../../css/icons.styles';

const TwitterIcon = ({ width = 40, height = 40 }) => (
  <svg
    width={width}
    height={height}
    viewBox="0 0 29 29"
    style={iconLarge}
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g filter="url(#filter0_d_6961_29753)">
      <rect x="4.62891" y="3" width="20.2472" height="20.2472" rx="4.82076" fill="#252E31" />
    </g>
    <path
      d="M19.0914 10.5854C18.7722 10.7266 18.429 10.8226 18.0689 10.865C18.4361 10.6447 18.7186 10.2958 18.8513 9.88061C18.5081 10.084 18.1268 10.2323 17.7215 10.3128C17.3966 9.96677 16.9334 9.74927 16.4221 9.74927C15.4391 9.74927 14.6412 10.5472 14.6412 11.5302C14.6412 11.67 14.6567 11.8056 14.6878 11.9356C13.2077 11.8621 11.8956 11.1531 11.0185 10.0755C10.8646 10.3382 10.777 10.6447 10.777 10.9709C10.777 11.5881 11.0906 12.1333 11.5694 12.4525C11.277 12.4426 11.003 12.3635 10.7629 12.2293C10.7629 12.2378 10.7629 12.2449 10.7629 12.2519C10.7629 13.1149 11.3759 13.8351 12.1908 13.9976C12.0411 14.0385 11.8843 14.0611 11.7219 14.0611C11.6075 14.0611 11.4945 14.0498 11.3872 14.0286C11.6131 14.7362 12.2713 15.2503 13.0495 15.2658C12.4408 15.7432 11.6725 16.0271 10.8392 16.0271C10.6951 16.0271 10.5539 16.0186 10.4141 16.0017C11.2021 16.5073 12.1371 16.8025 13.1427 16.8025C16.4179 16.8025 18.2087 14.0894 18.2087 11.7364C18.2087 11.6602 18.2059 11.5825 18.2031 11.5062C18.5505 11.2548 18.8527 10.9413 19.0914 10.5854Z"
      fill="white"
    />
    <defs>
      <filter
        id="filter0_d_6961_29753"
        x="0.772299"
        y="0.107544"
        width="27.9603"
        height="27.9604"
        filterUnits="userSpaceOnUse"
        colorInterpolationFilters="sRGB"
      >
        <feFlood floodOpacity="0" result="BackgroundImageFix" />
        <feColorMatrix
          in="SourceAlpha"
          type="matrix"
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
          result="hardAlpha"
        />
        <feOffset dy="0.964152" />
        <feGaussianBlur stdDeviation="1.9283" />
        <feColorMatrix type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.08 0" />
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_6961_29753" />
        <feBlend
          mode="normal"
          in="SourceGraphic"
          in2="effect1_dropShadow_6961_29753"
          result="shape"
        />
      </filter>
    </defs>
  </svg>
);

export default TwitterIcon;
