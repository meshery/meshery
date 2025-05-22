import React from 'react';
import { useTheme } from '@layer5/sistent';

const MoveFileIcon = (props) => {
  const theme = useTheme();

  return (
    <>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width={props.width || '24'}
        height={props.height || '24'}
        viewBox="0 0 24 24"
      >
        <g>
          <path
            fill={props.fill || theme.palette.icon.default}
            fill-rule="evenodd"
            d="M2 6a3 3 0 0 1 3-3h4.172a3 3 0 0 1 2.12.879L12.415 5H19a3 3 0 0 1 3 3v10a3 3 0 0 1-3 3H5a3 3 0 0 1-3-3zm11.08 10.664 2.667-3a1 1 0 0 0 0-1.328l-2.666-3a1 1 0 1 0-1.495 1.328L12.773 12H9a1 1 0 1 0 0 2h3.773l-1.187 1.336a1 1 0 0 0 1.495 1.328z"
            clip-rule="evenodd"
            opacity="1"
            data-original="#000000"
            class=""
          ></path>
        </g>
      </svg>
    </>
  );
};

export default MoveFileIcon;
