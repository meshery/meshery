import React from 'react';
import SvgIcon from '@material-ui/core/SvgIcon';

export function DotSquare(props) {
  return (
    <SvgIcon
      className="close"
      fontSize="inherit"
      style={{ width: 14, height: 14, fillOpacity: 0.5 }}
      fill={props.fill}
      {...props}
    >
      {/* tslint:disable-next-line: max-line-length */}
      <g clip-path="url(#clip0_23103_14908)">
        <path
          d="M22.047 22.074V1.927H1.927V22.074H22.047ZM22.047 24H1.927C1.39167 24 0.936667 23.8127 0.562 23.438C0.187333 23.0633 0 22.6083 0 22.073V1.926C0 1.40867 0.187333 0.958333 0.562 0.575C0.936667 0.191667 1.39167 0 1.927 0H22.074C22.5913 0 23.0417 0.191667 23.425 0.575C23.8083 0.958333 24 1.40867 24 1.926V22.073C24 22.6083 23.8083 23.0633 23.425 23.438C23.0417 23.8127 22.5823 24 22.047 24Z"
        />
        <circle cx="12" cy="12" r="4"/>
      </g>
      <defs>
        <clipPath id="clip0_23103_14908">
          <rect width="24" height="24" />
        </clipPath>
      </defs>
    </SvgIcon>
  );
}

export default DotSquare;
