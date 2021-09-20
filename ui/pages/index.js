import React from "react";
import { Button, Grid, Link, Stack, Typography } from "@mui/material";
import { Paper } from "@/components/index";
import { AdaptersChipList, AdaptersListContainer, MesheryServerVersionContainer } from "@/features/mesheryComponents";
import { PaperWithTitle } from "@/components/Paper";
import { nanoid } from "@reduxjs/toolkit";
import { useTheme } from "@mui/system";
import { KuberenetesClusterChip, KuberenetesClusterContainer } from "@/features/mesheryEnvironment";
import { GrafanaChip, MetricsContainer, PrometheusChip } from "@/features/mesheryEnvironment/components";
import SettingsIcon from "@mui/icons-material/Settings";
import { getMesheryVersionText } from "@/features/mesheryComponents/components/MesheryServer/helpers";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";

export default function Dashboard() {
  const theme = useTheme();

  const ConnectionStatus = () => (
    <PaperWithTitle title="Connection Status" titleVariant="h6">
      <Stack spacing={2}>
        <KuberenetesClusterContainer
          render={({ clusters }) => {
            if (clusters.length > 0)
              return (
                <PaperWithTitle title="Kuberetnes">
                  {clusters.map((cluster) => (
                    <Grid item xs={6} key={`cluster+${nanoid()}`} sx={{ mb: theme.spacing(1) }}>
                      <KuberenetesClusterChip cluster={cluster} />
                    </Grid>
                  ))}
                </PaperWithTitle>
              );
            return <PaperWithTitle title="Kubernetes">Not connected to kubernetes</PaperWithTitle>;
          }}
        />

        <AdaptersListContainer render={(props) => <AdaptersChipList {...props} />} />
        <MetricsContainer
          render={({ grafanas, prometheus, onGrafanaClick, onPrometheusClick }) => {
            return (
              <PaperWithTitle title="Metrics" containerProps={{ spacing: 2 }}>
                <Grid item xs={12} justifyContent="center" container>
                  {grafanas !== null && grafanas?.length > 0 && grafanas[0] ? (
                    <GrafanaChip grafana={grafanas[0]} />
                  ) : (
                    <Button
                      variant="contained"
                      startIcon={<SettingsIcon />}
                      color="primary"
                      size="large"
                      onClick={onGrafanaClick}
                    >
                      Configure Grafana
                    </Button>
                  )}
                </Grid>

                <Grid item xs={12} justifyContent="center" container>
                  {prometheus !== null && prometheus?.length > 0 && prometheus[0] ? (
                    <PrometheusChip prometheus={prometheus[0]} />
                  ) : (
                    <Button
                      variant="contained"
                      startIcon={<SettingsIcon />}
                      color="primary"
                      size="large"
                      onClick={onPrometheusClick}
                    >
                      Configure Prometheus
                    </Button>
                  )}
                </Grid>
              </PaperWithTitle>
            );
          }}
        />

        <MesheryServerVersionContainer
          render={({ serverVersion }) => (
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
                  <OpenInNewIcon fontSize="small" />
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
        />
      </Stack>
    </PaperWithTitle>
  );

  return (
    <Grid container spacing={2}>
      <Grid item xs={6}>
        <Paper>asdf</Paper>
      </Grid>
      <Grid item xs={6}>
        <ConnectionStatus />
      </Grid>
    </Grid>
  );
}
