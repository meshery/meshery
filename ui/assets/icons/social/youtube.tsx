export default function Youtube({ height, width, className }) {
  return (
    <svg
      height={height}
      width={width}
      className={className}
      style={{
        enableBackground: "new 0 0 400 400"
      }}
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 400 400"
    >
      <g id="Dark_Blue" style={{ display: "none" }}>
        <path
          fill="#FFFFFF"
          style={{ display: "inline" }}
          d="M350,400H50c-27.6,0-50-22.4-50-50V50C0,22.4,22.4,0,50,0h300c27.6,0,50,22.4,50,50v300
        C400,377.6,377.6,400,350,400z"
        />
      </g>
      <path
        fill="currentColor"
        d="M358.3,119.9c-3.8-14.3-15-25.6-29.2-29.4C303.3,83.6,200,83.6,200,83.6s-103.3,0-129.1,6.9
      c-14.2,3.8-25.4,15.1-29.2,29.4c-6.9,25.9-6.9,80.1-6.9,80.1s0,54.1,6.9,80.1c3.8,14.3,15,25.6,29.2,29.4c25.8,7,129.1,7,129.1,7
      s103.3,0,129.1-7c14.2-3.8,25.4-15.1,29.2-29.4c6.9-25.9,6.9-80.1,6.9-80.1S365.2,145.9,358.3,119.9z M166.2,249.1v-98.2l86.4,49.1
      L166.2,249.1z"
      />
    </svg>
  );
}
