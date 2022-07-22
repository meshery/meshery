/* eslint-disable react/prop-types */
import React from "react";
import { AdapterChip } from "./AdapterChip";
import { CircularProgress, Grid } from "@mui/material";
import { useTheme } from "@mui/system";

/**
 * React component that renders the list of adapters
 * @param {{adapters: import("../../mesheryComponentsSlice").AdaptersListType, loading: {boolean}}} props
 * @returns {import("react").ReactElement}
 */
export const AdaptersChipList = ({ adapters, loading = false }) => {
  const theme = useTheme();
  return (
    <>
      {loading && <CircularProgress data-testid="circular-progress" />}
      {adapters.map((adapter) => (
          <AdapterChip adapter={adapter}  handleClick={() => null} />
      ))}
    </>
  );
};
