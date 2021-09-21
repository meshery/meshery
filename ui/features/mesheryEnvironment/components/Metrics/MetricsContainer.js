import { useSelector } from "react-redux";

/**
 * React component that gets the metrics data from redux store and calls `render` method
 * with cluster data
 * @param {{render: () => import("react").ReactElement}} props
 * @returns {import("react").ReactElement}
 */

export const MetricsContainer = (props) => {
  /** @type {import("../../mesheryEnvironmentSlice").connectedGrafanas} */
  const grafanas = useSelector((state) => state.mesheryEnvironment.connectedGrafanas);
  /** @type {import("../../mesheryEnvironmentSlice").connectedPrometheus*/
  const prometheus = useSelector((state) => state.mesheryEnvironment.connectedPrometheus);

  // eslint-disable-next-line react/prop-types
  return props.children({ grafanas, prometheus });
};
