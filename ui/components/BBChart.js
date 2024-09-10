import { useRef, useEffect, memo } from 'react';
import { bb } from 'billboard.js';

const BBChart = ({ options }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    const chart = bb.generate({ ...options, bindto: chartRef.current });

    return () => {
      chart.destroy();
    };
  }, [options]);

  return <div ref={chartRef} onClick={(e) => e.stopPropagation()}></div>;
};

// Memoize the component to prevent unnecessary re-renders
export default memo(BBChart);
