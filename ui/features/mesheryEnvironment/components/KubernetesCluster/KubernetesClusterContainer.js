import { useDispatch, useSelector } from "react-redux";
import { useEffect } from "react";
import { mesherySystemSyncThunk } from "@/features/mesheryEnvironment/mesheryEnvironmentSlice";

/**
 * React component that gets the cluster data from redux store  and calls `render` method
 * with cluster data
 * @param {{render: () => import("react").ReactElement}} props
 * @returns {import("react").ReactElement}
 */

export const KuberenetesClusterContainer = (props) => {
  const dispatch = useDispatch();
  /** @type {import("@/features/mesheryEnvironment/mesheryEnvironmentSlice").KubernetesClusters} */
  const clusters = useSelector((state) => state.mesheryEnvironment.kubernetesClusters);

  useEffect(() => {
    dispatch(mesherySystemSyncThunk());
  }, []);

  // eslint-disable-next-line react/prop-types
  // return <div>{props.render({ clusters })}</div>;
  return props.children({ clusters });
};
