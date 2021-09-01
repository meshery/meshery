import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchComponentsStatusThunk,
  initialiseOperatorStatusSubscriptionThunk,
  loadingSelector,
  updateConnectionStatus,
} from "../mesheryComponentsSlice";
import { CircularProgress } from "@material-ui/core";
import { useEffect } from "react";
import { nanoid } from "@reduxjs/toolkit";

const ComponentsStatus = () => {
  const dispatch = useDispatch();
  const loading = useSelector(loadingSelector);
  const components = useSelector((state) => state.mesheryComponents);

  useEffect(() => {
    dispatch(fetchComponentsStatusThunk());
    dispatch(
      initialiseOperatorStatusSubscriptionThunk((res) =>
        dispatch(updateConnectionStatus(res))
      )
    );
  }, []);

  return (
    <>
      <div>{loading && <CircularProgress />}</div>
      <ul>
        {Object.keys(components)?.map((comp) => (
          <li key={nanoid()}>{comp} </li>
        ))}
      </ul>
    </>
  );
};

export default ComponentsStatus;
