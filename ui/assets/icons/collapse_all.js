import React from 'react';
import { useTheme } from '@layer5/sistent';

const CollapseAllIcon = (props) => {
  const theme = useTheme();

  return (
    <svg
      width="24"
      height="24"
      viewBox="0 0 20 20"
      fill={props.fill || theme.palette.icon.default}
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M6.16667 18.3334L5 17.1667L10 12.1667L15 17.1667L13.8333 18.3334L10 14.5L6.16667 18.3334ZM10 7.83335L5 2.83335L6.16667 1.66669L10 5.50002L13.8333 1.66669L15 2.83335L10 7.83335Z" />
    </svg>
  );
};

export default CollapseAllIcon;
