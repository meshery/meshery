import * as React from 'react';

const HoneycombCell = ({ children, row, column }) => {
  // Ensure row and column are valid numbers
  const safeRow = Number.isFinite(row) ? row : 1;
  const safeColumn = Number.isFinite(column) ? column : 1;

  return (
    <li
      style={{
        gridRow: `${safeRow} / span 4`,
        gridColumn: `${safeColumn} / span 4`,
        position: 'relative',
        pointerEvents: 'none',
        transform: safeRow % 2 ? 'translateX(25%)' : 'translateX(-25%)',
      }}
    >
      {children}
    </li>
  );
};

export default HoneycombCell;
