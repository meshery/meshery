import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  adaptersSelector,
  fetchAvailableAdaptersThunk,
  loadingSelector,
} from "../mesheryComponentsSlice";
import { CircularProgress } from "@material-ui/core";
import { useEffect } from "react";
import { nanoid } from "@reduxjs/toolkit";

const AdaptersList = () => {
  const dispatch = useDispatch();
  const loading = useSelector(loadingSelector);
  const adapters = useSelector(adaptersSelector);

  useEffect(() => {
    dispatch(fetchAvailableAdaptersThunk());
  }, []);

  return (
    <>
      <div>{loading && <CircularProgress />}</div>
      {adapters?.length && (
        <ul>
          {adapters.map((ad) => (
            <li key={nanoid()}>{ad.adapter_location}</li>
          ))}
        </ul>
      )}
    </>
  );
};

export default AdaptersList;
