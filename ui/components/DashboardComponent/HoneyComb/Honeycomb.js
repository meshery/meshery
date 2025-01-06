import * as React from 'react';
import { getColumnSize, getRowSize } from './helpers';
import HoneycombCell from './HoneycombCell.js';

const Honeycomb = React.forwardRef(({ items, renderItem, size, columns, className }, ref) => {
  const rowSize = getRowSize(size);
  const columnSize = getColumnSize(size);

  return (
    <ul
      ref={ref}
      className={className}
      style={{
        display: 'grid',
        gridTemplateColumns: `repeat(${columns * 4}, ${columnSize}px)`,
        justifyContent: 'center',
        gridAutoRows: `${rowSize}px`,
        padding: `0 ${columnSize}px`,
        listStyle: 'none',
      }}
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
    </ul>
  );
});

Honeycomb.displayName = 'Honeycomb';

export default Honeycomb;
