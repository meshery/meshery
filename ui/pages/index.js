import React, {useState} from "react";
import { Button, Grid, Link, Stack, Typography } from "@mui/material";
import { AdaptersChipList, AdaptersListContainer, MesheryServerVersionContainer } from "@/features/mesheryComponents";
import { PaperWithTitle } from "@/components/Paper";
import { nanoid } from "@reduxjs/toolkit";
import { useTheme } from "@mui/system";
import { KuberenetesClusterChip, KuberenetesClusterContainer } from "@/features/mesheryEnvironment";
import { GrafanaChip, MetricsContainer, PrometheusChip } from "@/features/mesheryEnvironment/components";
import SettingsIcon from "@mui/icons-material/Settings";
import { getMesheryVersionText } from "@/features/mesheryComponents/components/MesheryServer/helpers";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import { MetricsButton } from "@/components/Button";


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

        <AdaptersListContainer>{(props) => <AdaptersChipList {...props} />}</AdaptersListContainer>
        <MetricsContainer>
          {({ grafanas, prometheus, onGrafanaClick, onPrometheusClick }) => {
            return (
              <PaperWithTitle title="Metrics" containerProps={{ spacing: 2 }}>
                <Grid item xs={12} justifyContent="center" container>
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

                <Grid item xs={12} justifyContent="center" container>
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
              </PaperWithTitle>
            );
          }}
        </MetricsContainer>

        <MesheryServerVersionContainer>
          {({ serverVersion }) => (
            <PaperWithTitle title="Release">
              <Grid item xs={6}>
                <Typography variant="subtitle1" sx={{ fontWeight: "bold" }}>
                  Channel Subscription
                </Typography>
                {serverVersion.release_channel}
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle1" sx={{ fontWeight: "bold" }}>
                  Version
                </Typography>
                {getMesheryVersionText(serverVersion)}
                <Link
                  href={`https://docs.meshery.io/project/releases${
                    serverVersion.release_channel === "edge" ? "" : "/" + serverVersion.build
                  }`}
                  target="_blank"
                >
                  <OpenInNewIcon sx={{ fontSize: theme.spacing(2) }} />
                </Link>
              </Grid>
              <Grid item xs={12} sx={{ mt: theme.spacing(4) }}>
                {serverVersion.outdated ? (
                  <>
                    <Typography variant="subtitle1" sx={{ fontWeight: "bold" }}>
                      Newer version of Meshery available:
                    </Typography>
                    <Link
                      href={`https://docs.meshery.io/project/releases${
                        serverVersion.release_channel === "edge" ? "" : "/" + serverVersion.build
                      }`}
                      target="_blank"
                    >
                      {serverVersion.latest}
                    </Link>
                  </>
                ) : (
                  <Typography variant="body1" sx={{ fontStyle: "italic" }}>
                    Running latest version of Meshery
                  </Typography>
                )}
              </Grid>
            </PaperWithTitle>
          )}
        </MesheryServerVersionContainer>
      </Stack>
    </PaperWithTitle>
  );

  return (
    <Grid container spacing={2}>
      <Grid item xs={6}>
        <PaperWithTitle title="Service Mesh" titleVariant="h6"></PaperWithTitle>
      </Grid>
      <Grid item xs={6}>
        <ConnectionStatus />
      </Grid>
    </Grid>
  );
}
