const Filter = ({ fill, height, width, ...rest }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    height={height}
    width={width}
    viewBox="0 0 50 50"
    // className={className}
    {...rest}>
    <g>
      <path d="M30.7 0C30.7 0.1 30.7 0.2 30.7 0.3C30.7 3.4 28.1 6 25 6C21.8 6 19.3 3.5 19.3 0.3C19.3 0.2 19.3 0.1 19.3 0H0V49.6H50V0H30.7ZM23.9 44.3L21.5 32.4H21.4L18.8 44.3H15.5L11.7 26.8H15L17.3 38.7L20 26.7H23L25.5 38.8L28.1 26.7H31.3L27.1 44.2L23.9 44.3ZM42.1 44.3L41 40.4H35L34.1 44.3H30.8L35.1 26.8H40.3L45.5 44.3H42.1Z" fill={fill} />
      <path d="M35.6992 37.5L37.0992 31H38.4992L40.1992 37.5" fill={fill} />
    </g>
    <defs>
      <clipPath>
        <rect width={width} height={height} />
      </clipPath>
    </defs>

  </svg>
)

export default Filter;