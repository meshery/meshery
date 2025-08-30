import React from "react";

const DocumentColorIcon = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 256 256"
      width={props.width ? props.width : "24px"}
      height={props.height ? props.height : "24px"}
      onClick={props.onClick}
      className={props.className}
      style={{ ...props.style }}
    >
    <title>Document Text</title>
      <g fill="none" fillRule="nonzero" stroke="none" strokeWidth="1" strokeLinecap="butt" strokeLinejoin="miter" strokeMiterlimit="10" fontFamily="none" fontWeight="none" fontSize="none" textAnchor="none" style={{ mixBlendMode: "normal" }}>
        <g transform="scale(5.33333,5.33333)">
          <path d="M37,45h-26c-1.657,0 -3,-1.343 -3,-3v-36c0,-1.657 1.343,-3 3,-3h19l10,10v29c0,1.657 -1.343,3 -3,3z" fill="#00b39f"></path>
          <path d="M40,13h-10v-10z" fill="#bbdefb"></path>
          <path d="M30,13l10,10v-10z" fill="#1565c0"></path>
          <path d="M15,23h18v2h-18zM15,27h18v2h-18zM15,31h18v2h-18zM15,35h10v2h-10z" fill="#e3f2fd"></path>
        </g>
      </g>
    </svg>
  );
};

export default DocumentColorIcon;