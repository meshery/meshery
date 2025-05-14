import React from 'react';
import { useTheme } from '@layer5/sistent';

const MoveFileIcon = (props) => {
  const theme = useTheme();

  return (
    <svg
      width={props.width || '24'}
      height={props.height || '24'}
      viewBox="0 0 24 24"
      fill={props.fill || theme.palette.icon.default}
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M14 2H6C4.9 2 4 2.9 4 4V20C4 20.41 4.12 20.8 4.34 21.12C4.41 21.23 4.5 21.33 4.59 21.41C4.95 21.78 5.45 22 6 22H13.53C13 21.42 12.61 20.75 12.35 20H6V4H13V9H18V12C18.7 12 19.37 12.12 20 12.34V8L14 2M18 23L23 18.5L20 15.8L18 14V17H14V20H18V23Z" />
    </svg>
  );
};

export default MoveFileIcon;
