/* eslint-disable react/prop-types */
import React from "react";
import { Tooltip } from "@mui/material";
import Image from "next/image";
import { useTheme } from "@mui/system";
import Chip from "@/components/Chip";

/**
 * React component that takes in adapter information among other things
 * as props and renders it in form of chip
 * @param {{adapter: import("../../mesheryComponentsSlice").AdapterType}} props
 * @returns {import("react").ReactElement}
 */
export const AdapterChip = ({ adapter, handleClick }) => {
  const theme = useTheme();
  let image = "/static/img/meshery-logo/meshery-logo.svg";
  if (adapter?.name) image = "/static/img/" + adapter.name.toLowerCase() + ".svg";
  const logoIcon = (
    <Image
      src={image}
      alt={adapter?.name?.toLowerCase() || "adapter"}
      width={theme.spacing(2.5)}
      height={theme.spacing(2.5)}
    />
  );
  return (
    <Tooltip title={!adapter.isActive ? "This adapter is inactive" : `Active`}>
      <Chip
        label={adapter.adapter_location.split(":")[0]}
        onClick={handleClick}
        icon={logoIcon}
        sx={{m: theme.spacing(1)}}
        variant={!adapter.isActive ? "default" : "outlined"}
        data-testid="adapter-chip"
      />
    </Tooltip>
  );
};
