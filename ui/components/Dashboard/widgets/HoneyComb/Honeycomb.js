import * as React from 'react';
import { getColumnSize, getRowSize } from './helpers';
import HoneycombCell from './HoneycombCell.js';
import { HoneycombContainer } from '../../style';

const Honeycomb = React.forwardRef(({ items, renderItem, size, columns, className }, ref) => {
  const rowSize = getRowSize(size);
  const columnSize = getColumnSize(size);

  return (
    <HoneycombContainer
      ref={ref}
      className={className}
      columnSize={columnSize}
      columns={columns}
      rowSize={rowSize}
    >
      {items.map((item, index) => {
        const row = 1 + Math.floor(index / columns) * 3;
        const column = 1 + (index % columns) * 4;
        const renderedItem = renderItem(item, index);

        return (
          <HoneycombCell key={index} row={row} column={column}>
            {renderedItem}
          </HoneycombCell>
        );
      })}
    </HoneycombContainer>
  );
});

Honeycomb.displayName = 'Honeycomb';

export default Honeycomb;
