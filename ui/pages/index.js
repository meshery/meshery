import React, {useState} from "react";
import { Button, Grid, Link, Stack, Typography, Box } from "@mui/material";
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { AdaptersChipList, AdaptersListContainer } from "@/features/mesheryComponents";
import { PaperWithTitle } from "@/components/Paper";
import { nanoid } from "@reduxjs/toolkit";
import { useTheme } from "@mui/system";
import { KuberenetesClusterChip, KuberenetesClusterContainer } from "@/features/mesheryEnvironment";
import { GrafanaChip, MetricsContainer, PrometheusChip } from "@/features/mesheryEnvironment/components";
import SettingsIcon from "@mui/icons-material/Settings";
import { MetricsButton } from "@/components/Button"

export default function Dashboard() {
  const theme = useTheme();

  const ConnectionStatus = () => (
    <PaperWithTitle title="Connection Status" titleVariant="h6">
      <Stack spacing={2}>
        <KuberenetesClusterContainer>
          {({ clusters }) => {
            if (clusters.length > 0 && clusters.filter((clstr) => !!clstr).length > 0)
              return (
                <PaperWithTitle title="Kubernetes">
                  {clusters.map((cluster) => (
                    <Grid item xs={6} key={`cluster+${nanoid()}`} sx={{ mb: theme.spacing(1) }}>
                      <KuberenetesClusterChip cluster={cluster} />
                    </Grid>
                  ))}
                </PaperWithTitle>
              );
            return <PaperWithTitle title="Kubernetes">Not connected to kubernetes</PaperWithTitle>;
          }}
        </KuberenetesClusterContainer>

        <AdaptersListContainer>{(props) => <PaperWithTitle title="Adapters">
        <Grid item xs={6}  sx={{ margin: theme.spacing(2.5)}}>
           <AdaptersChipList {...props} />
           </Grid> </PaperWithTitle> }</AdaptersListContainer>
        <MetricsContainer justifyContent="center" >
          {({ grafanas, prometheus, onGrafanaClick, onPrometheusClick }) => {
            return (
              <PaperWithTitle title="Metrics" containerProps={{ spacing: 2 }}>
               <Grid container  spacing={2} sx={{margin: "auto ", justifyContent:"center" }}>
                <Grid item  >
                  {grafanas !== null && grafanas?.length > 0 && grafanas[0] ? (
                    <GrafanaChip grafana={grafanas[0]} />
                  ) : (
                    <MetricsButton
                      variant="contained"
                      startIcon={<SettingsIcon />}
                      onClick={onGrafanaClick}
                      size="large"
                      
                    >
                      Configure Grafana
                    </MetricsButton>
                  )}
                </Grid>

                <Grid item  >
                  {prometheus !== null && prometheus?.length > 0 && prometheus[0] ? (
                    <PrometheusChip prometheus={prometheus[0]} />
                  ) : (
                    <MetricsButton
                      variant="contained"
                      startIcon={<SettingsIcon />}
                      size="large"
                      onClick={onPrometheusClick}
                    >
                      Configure Prometheus
                    </MetricsButton>
                  )}
                </Grid>
                </Grid>
              </PaperWithTitle>
            );
          }}
        </MetricsContainer>
      </Stack>
    </PaperWithTitle>
  );
  
const ShowServiceMesh = () => {
  return(
    <Box
    sx={{
      margin: "auto",
      display : "flex",
      justifyContent : "center",
      flexDirection : "column",
    }}
  >
    <Typography sx={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" color="textSecondary">
              No service meshes detected in the cluster.
              </Typography>
              <Button
                aria-label="Add Meshes"
                variant="contained"
                color="primary"
                size="large"
              >
                <AddCircleOutlineIcon  />
              Install Service Mesh
              </Button>
              </Box>
  )
}

  return (
    <Grid container spacing={2}>
      <Grid item xs={6}>
        <PaperWithTitle title="Service Mesh" titleVariant="h6">
          <ShowServiceMesh />
        </PaperWithTitle>
      </Grid>
      <Grid item xs={6}>
        <ConnectionStatus />
      </Grid>
    </Grid>
  );
}
