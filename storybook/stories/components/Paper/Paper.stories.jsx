import React from "react";
import Stack from "@mui/material/Stack";
import { Paper } from "./Paper";

export default {
  title : "Components/Paper",
  component : Paper,
  tags: ['autodocs'],
};

export function BasicPaper() {
  return (
    <Stack direction="row" spacing={1}>
      <Paper variant="outlined" />
    </Stack>
  )
}

export function ElevationPaper() {
  return (
    <Stack direction="row" spacing={1}>
      <Paper elevation={12} />
    </Stack>
  )
}