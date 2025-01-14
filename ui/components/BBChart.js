import { useRef, useEffect, memo } from 'react';
import { bb } from 'billboard.js';
import { ErrorBoundary } from '@layer5/sistent';

const BBChart = ({ options }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    const chart = bb.generate({ ...options, bindto: chartRef.current });

    return () => {
      chart.destroy();
    };
  }, [options]);

  return (
    <ErrorBoundary>
      <div ref={chartRef} onClick={(e) => e.stopPropagation()}></div>
    </ErrorBoundary>
  );
};

// Memoize the component to prevent unnecessary re-renders
export default memo(BBChart);
