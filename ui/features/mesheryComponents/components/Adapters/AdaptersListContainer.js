import { useDispatch, useSelector } from "react-redux";
import { adaptersSelector, fetchAvailableAdaptersThunk, loadingSelector } from "../../mesheryComponentsSlice";
import { useEffect } from "react";

/**
 * React component that fetches the available adapters and calls `render` method
 * with adapters list and loading state
 * @param {{render: () => import("react").ReactElement}} props
 * @returns {import("react").ReactElement}
 */

export const AdaptersListContainer = (props) => {
  const dispatch = useDispatch();

  /**
   * @type {true | false}
   */
  const loading = useSelector(loadingSelector);
  /**
   * @type {import("../../mesheryComponentsSlice.js").AdaptersListType}
   */
  const adapters = useSelector(adaptersSelector);

  useEffect(() => {
    dispatch(fetchAvailableAdaptersThunk());
  }, []);

  // eslint-disable-next-line react/prop-types
  return props.children({ loading, adapters });
};
