import * as React from 'react';
import ResizeObserver from 'resize-observer-polyfill';

import Honeycomb, { type HoneycombProps } from './Honeycomb';
import { getGridColumnsCount } from './helpers';

type ResponsiveHoneycombProps<T> = Omit<HoneycombProps<T>, 'columns' | 'containerRef'> & {
  defaultWidth: number;
};

const ResponsiveHoneycomb = <T,>({
  size,
  defaultWidth,
  ...restProps
}: ResponsiveHoneycombProps<T>) => {
  const containerRef = React.useRef<HTMLDivElement>(null);
  const [columns, setColumns] = React.useState(() =>
    Math.max(1, getGridColumnsCount(size, defaultWidth)),
  );

  const updateColumns = React.useCallback(
    (width: number) => {
      if (width <= 0) {
        return;
      }

      const nextColumns = Math.max(1, getGridColumnsCount(size, width));
      setColumns((previousColumns) =>
        previousColumns === nextColumns ? previousColumns : nextColumns,
      );
    },
    [size],
  );

  React.useEffect(() => {
    const target = containerRef.current;

    if (!target) {
      return;
    }

    const observer = new ResizeObserver((entries) => {
      entries.forEach((entry) => {
        updateColumns(entry.contentRect.width);
      });
    });

    updateColumns(target.clientWidth);
    observer.observe(target);

    return () => observer.disconnect();
  }, [updateColumns]);

  return <Honeycomb containerRef={containerRef} size={size} {...restProps} columns={columns} />;
};

export default ResponsiveHoneycomb;
