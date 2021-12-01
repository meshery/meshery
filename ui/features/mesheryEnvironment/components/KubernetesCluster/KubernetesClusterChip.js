/* eslint-disable react/prop-types */
import React from "react";
import { Tooltip } from "@mui/material";
import Image from "next/image";
import { useTheme } from "@mui/system";
import Chip from "@/components/Chip";
import { useSnackbar } from 'notistack';

/**
 * React component that takes in kubernetes clusters information among other things
 * as props and renders it in form of chip
 * @param {{cluster: import("@/features/mesheryEnvironment/mesheryEnvironmentSlice").KubernetesCluster}} props
 * @returns {import("react").ReactElement}
 */

export const KuberenetesClusterChip = ({ cluster, handleClick }) => {

  const theme = useTheme();
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const handleKubernetesClick = () => {
     fetch(
      "/api/system/kubernetes/ping",
      { credentials : "same-origin",
        credentials : "include", },
     )
     .then(function(response) {
       if(!response.ok){
         enqueueSnackbar("Kubernetes cannot be pinged",{
           variant : "error"
         });
       }
       else{
         enqueueSnackbar("Kubernetes pinged successfully",{
           variant : "success"
         });
       }
     })
  }
  let image = "/static/img/kubernetes.svg";
  const logoIcon = (
    <Image
      src={image}
      style={{ width: "1rem" }}
      alt={cluster?.configuredServer}
      width={theme.spacing(2.5)}
      height={theme.spacing(2.5)}
    />
  );
  return (
    <Tooltip title={`Server: ${cluster?.configuredServer}`}>
      <Chip
        label={cluster.inClusterConfig ? "Using In Cluster Config" : cluster.contextName}
        onClick={handleKubernetesClick}
        icon={logoIcon}
        variant="outlined"
      />
    </Tooltip>
  );
};
