/* eslint-disable react/prop-types */
import React from "react";
import { Grid } from "@mui/material";
import { AdapterChip } from "./AdapterChip";
import { CircularProgress } from "@mui/material";

/**
 * React component that renders the list of adapters
 * @param {{adapters: import("../mesheryComponentsSlice").AdaptersListType, loading: {boolean}}} props
 * @returns {import("react").ReactElement}
 */
export const AdaptersChipList = ({ adapters, loading }) => (
  <Grid>
    {loading && <CircularProgress />}
    {adapters.map((adapter) => (
      <AdapterChip key={`${adapter?.adapter_location}`} adapter={adapter} isDisabled={false} />
    ))}
  </Grid>
);
