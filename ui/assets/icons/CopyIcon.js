import React from 'react';

export const CopyIcon = ({
  width = 24,
  height = 24,
  fill = 'currentColor',
  ...props
}) => (
  <svg xmlns="http://www.w3.org/2000/svg" width={width} height={height} viewBox="0 0 18 18" fill="none"  {...props}>
    <path fillRule="evenodd" clipRule="evenodd" d="M12 0.75H3C2.175 0.75 1.5 1.425 1.5 2.25V12.75H3V2.25H12V0.75ZM14.25 3.75H6C5.175 3.75 4.5 4.425 4.5 5.25V15.75C4.5 16.575 5.175 17.25 6 17.25H14.25C15.075 17.25 15.75 16.575 15.75 15.75V5.25C15.75 4.425 15.075 3.75 14.25 3.75ZM14.25 15.75H6V5.25H14.25V15.75Z" fill={fill} />
  </svg>
);

export default CopyIcon;