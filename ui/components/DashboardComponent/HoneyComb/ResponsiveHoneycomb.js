import * as React from 'react';
import ResizeObserver from 'resize-observer-polyfill';

import Honeycomb from './Honeycomb';
import { getGridColumnsCount } from './helpers';

const ResponsiveHoneycomb = ({ size, defaultWidth, ...restProps }) => {
  const containerRef = React.useRef(null);
  const [columns, setColumns] = React.useState(() => {
    const calculatedColumns = getGridColumnsCount(size, defaultWidth);
    return Math.max(1, calculatedColumns); // Ensure at least 1 column
  });

  React.useEffect(() => {
    if (!containerRef.current) return;

    const target = containerRef.current;
    const observer = new ResizeObserver((entries) => {
      entries.forEach((entry) => {
        const width = entry.contentRect.width;
        if (width > 0) {
          const calculatedColumns = getGridColumnsCount(size, width);
          setColumns(Math.max(1, calculatedColumns));
        }
      });
    });

    // Initial calculation
    const initialWidth = target.clientWidth;
    if (initialWidth > 0) {
      const calculatedColumns = getGridColumnsCount(size, initialWidth);
      setColumns(Math.max(1, calculatedColumns));
    }

    observer.observe(target);
    return () => observer.disconnect();
  }, [size]);

  return <Honeycomb ref={containerRef} size={size} {...restProps} columns={columns} />;
};

export default ResponsiveHoneycomb;
