import React from "react";
import { Grid, Stack } from "@mui/material";
import { Paper } from "@/components/index";
import { AdaptersChipList, AdaptersListContainer } from "@/features/mesheryComponents";
import { KuberenetesClusterContainer } from "@/features/mesheryComponents/components/KubernetesCluster/KubernetesClusterContainer";
import { KuberenetesClusterChip } from "@/features/mesheryComponents/components/KubernetesCluster/KubernetesClusterChip";
import { PaperWithTitle } from "@/components/Paper";
import { nanoid } from "@reduxjs/toolkit";
import { useTheme } from "@mui/system";

export default function Dashboard() {
  const theme = useTheme();
  return (
    <Grid container spacing={2}>
      <Grid item xs={6}>
        <Paper>asdf</Paper>
      </Grid>
      <Grid item xs={6}>
        <Paper>
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
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
}
