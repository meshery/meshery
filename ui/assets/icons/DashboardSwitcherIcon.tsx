import React from 'react';

function DashboardSwitcherIcon({ width = '24px', height = '24px', ...props }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 32 32"
      width={width}
      height={height}
      {...props}
    >
      <circle
        cx="25.6"
        cy="23.6"
        r="4.4"
        fill="#fff"
        stroke="#fff"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path fill="#fff" d="M11.8,3.6h-6.8v1.2h6.8v-1.2Z" />
      <path fill="#fff" d="M15.2,5.1H1.7v1.2h13.5v-1.2Z" />
      <path
        fill="#fff"
        d="M2.7,29.2h2.3v-4c0-1.9,1.5-3.4,3.4-3.4s3.4,1.5,3.4,3.4v4h2.3V6.6H2.7v22.5ZM9.3,8
          c0-.1,0-.2.2-.2h2.5c.1,0,.2,0,.2.2v2.1c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1ZM9.3,11.6c0-.1,0-.2.2-.2h2.5
          c.1,0,.2,0,.2.2v2.1c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1ZM9.3,14.8c0-.1,0-.2.2-.2h2.5c.1,0,.2,0,.2.2v2.1
          c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1ZM9.3,18.3c0-.1,0-.2.2-.2h2.5c.1,0,.2,0,.2.2v2.1c0,.1,0,.2-.2.2h-2.5
          c-.1,0-.2,0-.2-.2v-2.1ZM4.6,8c0-.1,0-.2.2-.2h2.5c.1,0,.2,0,.2.2v2.1c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1ZM4.6,11.6
          c0-.1,0-.2.2-.2h2.5c.1,0,.2,0,.2.2v2.1c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1ZM4.6,14.8c0-.1,0-.2.2-.2h2.5
          c.1,0,.2,0,.2.2v2.1c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1ZM4.6,18.3c0-.1,0-.2.2-.2h2.5c.1,0,.2,0,.2.2v2.1
          c0,.1,0,.2-.2.2h-2.5c-.1,0-.2,0-.2-.2v-2.1Z"
      />
      <path
        fill="#fff"
        stroke="#fff"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M24.1,12.8l-4.6-7.5c0-.1-.2-.2-.3-.3-.4-.2-1,0-1.2.3l-1.3,2.1v6.7h6.8c.2,0,.4,0,.5-.2.4-.3.5-.8.2-1.2Z"
      />
    </svg>
  );
}

export default DashboardSwitcherIcon;
