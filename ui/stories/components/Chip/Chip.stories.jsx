import React from "react";
import Stack from "@mui/material/Stack";
import FaceIcon from "@mui/icons-material/Face";
import { styled, useTheme } from "@mui/material";
import { Chip } from "./Chip";

export default {
  title : "Components/Chip",
  component : Chip,
};

const CustomChip = styled(Chip)(({ theme }) => ({
  marginRight : theme.spacing(1),
  marginBottom : theme.spacing(1),
}));

export function Default(props) {
  return (
    <Chip {...props} />
  )
}

export function BasicChip() {
  return (
    <Stack direction="row" spacing={1}>
      <Chip label="Chip Filled" />
      <Chip label="Chip Outlined" variant="outlined" />
    </Stack>
  )
}

export function KubernetesChip() {
  const theme = useTheme();
  return (
    <Stack direction="row" spacing={1}>
      <Chip
        icon={<FaceIcon />}
        label="With Icon"
        variant="primary"
      />
      <CustomChip
        icon={<img src="/static/img/kubernetes.svg" sx={{ width : theme.spacing(2.5) }} style={{ width : 24, height : 24 }} />}
        label="0"
        variant="outlined"
      />
    </Stack>
  )
}

export function ColorChips() {
  return (
    <Stack spacing={1} alignItems="center">
      <Stack direction="row" spacing={1}>
        <Chip label="primary" color="primary" />
        <Chip label="success" color="success" />
      </Stack>
      <Stack direction="row" spacing={1}>
        <Chip label="primary" color="primary" variant="outlined" />
        <Chip label="success" color="success" variant="outlined" />
      </Stack>
    </Stack>
  )
}

export function ClickableChips() {
  const handleClick = () => {
    console.log("You clicked the Chip");
  };

  return (
    <Stack direction="row" spacing={1}>
      <Chip label="Clickable" onClick={handleClick} />
      <Chip label="Clickable" variant="outlined" onClick={handleClick} />
    </Stack>
  )
}

export function ClickableAndDeletableChips() {
  const handleClick = () => {
    console.info('You clicked the Chip.');
  };

  const handleDelete = () => {
    console.info('You clicked the delete icon.');
  };

  return (
    <Stack direction="row" spacing={1}>
      <Chip
        label="Clickable Deletable"
        onClick={handleClick}
        onDelete={handleDelete}
      />
      <Chip
        label="Clickable Deletable"
        variant="outlined"
        onClick={handleClick}
        onDelete={handleDelete}
      />
    </Stack>
  );
}

export function ClickableLinkChips() {
  return (
    <Stack direction="row" spacing={1}>
      <Chip label="Clickable Link" component="a" href="#basic-chip" clickable />
      <Chip
        label="Clickable Link"
        component="a"
        href="#basic-chip"
        variant="outlined"
        clickable
      />
    </Stack>
  );
}