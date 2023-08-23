const ServiceMeshIcon = (props) => {
  const fill= props.fill ? props.fill : "currentColor"
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 172.63 184.19"
      width={props.width ? props.width : "24px"}
      height={props.height ? props.height : "24px"}
      fill={props.fill ? props.fill : "currentColor"}
      onClick={props.onClick}
      className={props.className}
      color={props.color ? props.color : "unset"}
      fontSize={props.fontSize ? props.fontSize : "unset"}
      style={{ ...props.style }}

    >
      <defs>
        <style>
          {
          `.cls-1{fill:${fill}}.cls-3{opacity:0.8;fill:none;stroke:${fill};stroke-miterlimit:10;stroke-width:5px}`
          }
        </style>
      </defs>
      <g id="Layer_2" data-name="Layer 2">
        <g id="Layer_1-2" data-name="Layer 1">
          <path
            d="M146.32 184.19A26.29 26.29 0 0 1 120 157.94a26.3 26.3 0 0 1 26.28-26.28 26.28 26.28 0 0 1 26.27 26.28 26.27 26.27 0 0 1-26.23 26.25ZM74.55 118.37a26.3 26.3 0 0 1-26.28-26.28 26.29 26.29 0 0 1 26.28-26.27 26.27 26.27 0 0 1 26.27 26.27 26.28 26.28 0 0 1-26.27 26.28ZM146.32 118.37A26.3 26.3 0 0 1 120 92.09a26.29 26.29 0 0 1 26.28-26.27 26.27 26.27 0 0 1 26.27 26.27 26.28 26.28 0 0 1-26.23 26.28ZM146.35 52.54a26.27 26.27 0 0 1-26.27-26.27A26.27 26.27 0 0 1 146.35 0a26.28 26.28 0 0 1 26.28 26.27 26.28 26.28 0 0 1-26.28 26.27Z"
            className="cls-1"
          />
          <g
            style={{
              opacity : 0.8,
            }}
          >
            <path
              d="M48.27 92.09C32.33 92.09 30 43.94 1.51 22.42M120 157.94c-4.78 0-7.17-16.46-9.56-32.91s-4.75-32.92-9.58-32.92"
              className="cls-3"
            />
            <path
              d="M120 157.94c-4.78 0-7.17-16.46-9.56-32.91s-4.75-32.92-9.58-32.92"
              className="cls-3"
            />
            <path
              d="M100.82 92.09c4.79 0 7.18.09 9.57.17s4.78.17 9.57.17"
              className="cls-3"
            />
            <path
              d="M100.82 92.09c4.79 0 7.18.09 9.57.17s4.78.17 9.57.17M146.32 131.64c0-3.31.08-5 .16-6.62s.17-3.32.17-6.64"
              className="cls-3"
            />
            <path
              d="M146.32 131.64c0-3.31.08-5 .16-6.62s.17-3.32.17-6.64M146.32 65.82c0-3.31.08-5 .16-6.62s.17-3.32.17-6.64"
              className="cls-3"
            />
            <path
              d="M146.32 65.82c0-3.31.08-5 .16-6.62s.17-3.32.17-6.64M100.82 92.09c4.83 0 7.24-16.45 9.65-32.9s4.83-32.91 9.66-32.91"
              className="cls-3"
            />
            <path
              d="M100.82 92.09c4.83 0 7.24-16.45 9.65-32.9s4.83-32.91 9.66-32.91"
              className="cls-3"
            />
          </g>
        </g>
      </g>
    </svg>
  )
}

export default ServiceMeshIcon