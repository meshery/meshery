import { useDispatch, useSelector } from "react-redux";
import { useEffect } from "react";
import { fetchMesheryServerDetailsThunk } from "@/features/mesheryComponents/mesheryComponentsSlice";

/**
 * React component that gets the cluster data from redux store  and calls `render` method
 * with cluster data
 * @param {{render: () => import("react").ReactElement}} props
 * @returns {import("react").ReactElement}
 */

export const MesheryServerVersionContainer = (props) => {
  const dispatch = useDispatch();
  /** @type {import("@/features/mesheryComponents/mesheryComponentsSlice").MesheryServerVersionType} */
  const serverVersion = useSelector((state) => state.mesheryComponents.server.version);

  useEffect(() => {
    dispatch(fetchMesheryServerDetailsThunk());
  }, []);

  // eslint-disable-next-line react/prop-types
  return props.children({ serverVersion });
};
