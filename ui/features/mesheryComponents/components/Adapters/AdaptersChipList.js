/* eslint-disable react/prop-types */
import React from "react";
import { AdapterChip } from "./AdapterChip";
import { CircularProgress, Grid } from "@mui/material";
import { PaperWithTitle } from "@/components/Paper";
import { useTheme } from "@mui/system";

/**
 * React component that renders the list of adapters
 * @param {{adapters: import("../../mesheryComponentsSlice").AdaptersListType, loading: {boolean}}} props
 * @returns {import("react").ReactElement}
 */
export const AdaptersChipList = ({ adapters, loading = false }) => {
  const theme = useTheme();
  return (
    <PaperWithTitle title="Adapters">
      {loading && <CircularProgress data-testid="circular-progress" />}
      {adapters.map((adapter) => (
        <Grid item xs={6} key={`${adapter?.adapter_location}`} sx={{ mb: theme.spacing(1) }}>
          <AdapterChip adapter={adapter} handleClick={() => null} />
        </Grid>
      ))}
    </PaperWithTitle>
  );
};
