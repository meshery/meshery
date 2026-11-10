import * as React from 'react';
import { getColumnSize, getRowSize } from './helpers';
import HoneycombCell from './HoneycombCell';
import { HoneycombContainer } from '../../style';

export type HoneycombProps<T> = {
  items: T[];
  renderItem: (item: T, index: number) => React.ReactNode;
  size: number;
  columns: number;
  className?: string;
  containerRef?: React.Ref<HTMLDivElement>;
};

const Honeycomb = <T,>({
  items,
  renderItem,
  size,
  columns,
  className,
  containerRef,
}: HoneycombProps<T>) => {
  const rowSize = getRowSize(size);
  const columnSize = getColumnSize(size);

  return (
    <HoneycombContainer
      ref={containerRef}
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
};

export default Honeycomb;
