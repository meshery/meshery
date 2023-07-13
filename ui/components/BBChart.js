import { useRef , useEffect } from "react"
import { bb } from "billboard.js"

const BBChart = ({ options }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    const chart = bb.generate({ ...options,bindto : chartRef.current });

    return () => {
      chart.destroy();
    };
  }, [options]);

  return <div ref={chartRef} ></div>;
};

export default BBChart