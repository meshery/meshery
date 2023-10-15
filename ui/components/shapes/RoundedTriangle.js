export function RoundedTriangleShape({ color }) {
  return (
    <svg width="30" height="30" viewBox="0 0 30 30" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M19.547 2.625c-2.021-3.5-7.073-3.5-9.094 0L.797 19.351c-2.021 3.5.505 7.875 4.546 7.875h19.314c4.041 0 6.567-4.375 4.546-7.875L19.547 2.625Z"
        fill={color}
        stroke="white"
        stroke-width="1.5"
      />
    </svg>
  );
}
