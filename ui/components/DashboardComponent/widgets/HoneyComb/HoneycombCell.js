import * as React from 'react';
import { HoneycombCell as StyledHoneycombCell } from '../../style';

const HoneycombCell = ({ children, row, column }) => {
  const safeRow = Number.isFinite(row) ? row : 1;
  const safeColumn = Number.isFinite(column) ? column : 1;

  return (
    <StyledHoneycombCell row={safeRow} column={safeColumn}>
      {children}
    </StyledHoneycombCell>
  );
};

export default HoneycombCell;
