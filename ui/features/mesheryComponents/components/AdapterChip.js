/* eslint-disable react/prop-types */
import React from "react";
import { Chip } from "@material-ui/core";
import { Tooltip } from "@material-ui/core";

/**
 * React component that takes in adapter information among other things
 * as props and renders it in form of chip
 * @param {{adapter: import("../mesheryComponentsSlice").AdapterType}} props
 * @returns {import("react").ReactElement}
 */
export const AdapterChip = ({ adapter, handleClick, isDisabled }) => {
  const image = "/static/img/" + adapter.name.toLowerCase() + ".svg";
  const logoIcon = <img src={image} style={{ width: "1rem" }} />;
  return (
    <Tooltip title={isDisabled ? "This adapter is inactive" : `Active`}>
      <Chip
        label={adapter.adapter_location.split(":")[0]}
        onClick={handleClick}
        icon={logoIcon}
        variant={isDisabled ? "default" : "outlined"}
      />
    </Tooltip>
  );
};
