import React from "react";
import Stack from "@mui/material/Stack";
import { Tab } from "./Tab"

export default {
  title : "Components/Tab",
  component : Tab,
  tags: ['autodocs'],
};

export function BasicTab() {
  return (
    <Stack direction="row" spacing={1}>
      <Tab />
    </Stack>
  )
}

export function ColoredTab() {
  return (
    <Stack direction="row" spacing={1}>
      <Tab textColor="secondary" indicatorColor="secondary" />
    </Stack>
  )
}

export function DisableddTab() {
  return (
    <Stack direction="row" spacing={1}>
      <Tab disabled={true} />
    </Stack>
  )
}

export function verticalTab() {
  return (
    <Stack direction="row" spacing={1}>
      <Tab orientation="vertical"/>
    </Stack>
  )
}