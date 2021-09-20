import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchComponentsStatusThunk,
  initialiseOperatorStatusSubscriptionThunk,
  loadingSelector,
  mesheryComponentsSelector,
  updateConnectionStatus,
} from "../../mesheryComponentsSlice";
import { useEffect } from "react";

/**
 * React component that fetches the components status, initiates subscription and calls `render` method
 * with  components list and loading state
 * @param {{render, components}) => import("react").ReactElement}} props
 * @returns {import("react").ReactElement}
 */

const ComponentsStatus = (props) => {
  const dispatch = useDispatch();
  /** @type {boolean} */
  const loading = useSelector(loadingSelector);

  /**
   * @type
   * {{operator: import("../../mesheryComponentsSlice").MesheryComponent,
   * meshsync: import("../../mesheryComponentsSlice").MesheryComponent,
   *  broker: import("../../mesheryComponentsSlice").MesheryComponent,
   *  server: import("../../mesheryComponentsSlice").MesheryComponent}>}
   *  */
  const components = useSelector(mesheryComponentsSelector);
  const operatorError = useSelector((state) => state.mesheryComponents.operatorError);

  useEffect(() => {
    dispatch(fetchComponentsStatusThunk());
    dispatch(initialiseOperatorStatusSubscriptionThunk((res) => dispatch(updateConnectionStatus(res))));
  }, []);

  // eslint-disable-next-line react/prop-types
  return <div>{props.render({ loading, components, operatorError })}</div>;
};

export default ComponentsStatus;
