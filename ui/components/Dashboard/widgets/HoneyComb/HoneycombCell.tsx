import * as React from 'react';
import { HoneycombCell as StyledHoneycombCell } from '../../style';

type HoneycombCellProps = {
  children?: React.ReactNode;
  row?: number;
  column?: number;
};

const HoneycombCell = ({ children, row, column }: HoneycombCellProps) => {
  const safeRow = Number.isFinite(row) ? row : 1;
  const safeColumn = Number.isFinite(column) ? column : 1;

  return (
    <StyledHoneycombCell row={safeRow} column={safeColumn}>
      {children}
    </StyledHoneycombCell>
  );
};

export default HoneycombCell;
