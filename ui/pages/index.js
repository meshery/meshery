import React, {useState, useEffect} from "react";
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
import showServiceMesh from "@/components/Dashboard/showServiceMesh"

export default function Dashboard() {
  const theme = useTheme();
  const [meshScan, SetMeshScan] = useState([])
  const [meshScanNamespaces, SetMeshScanNamespaces]  = useState()
  const [activeMeshScanNamespace , SetActiveMeshScanNamespace] = useState()

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

  useEffect(() => {
    setMeshScanData
  });

 function setMeshScanData (controlPlanesData, dataPlanesData) {
    const namespaces = {};
    const activeNamespaces = {};
    const processedControlPlanesData = controlPlanesData?.controlPlanesState?.map((mesh) => {
      if (!mesh?.members?.length) {
        return;
      }
      let proxies = []

      if (Array.isArray(dataPlanesData?.dataPlanesState)) {
        const dataplane = dataPlanesData.dataPlanesState.find(mesh_ => mesh_.name === mesh.name)

        if (Array.isArray(dataplane?.proxies)) proxies = dataplane.proxies
      }
      const processedMember = mesh?.members?.map((member) => {
        if (namespaces[mesh.name]) {
          namespaces[mesh.name].add(member.namespace);
        } else {
          namespaces[mesh.name] = new Set([member.namespace]);
        }

        // retrieve data planes according to mesh name
        if (proxies.length > 0) {
          const controlPlaneMemberProxies = proxies.filter(proxy => proxy.controlPlaneMemberName === member.name)

          if (controlPlaneMemberProxies.length > 0) {
            member = {
              ...member,
              data_planes : controlPlaneMemberProxies
            }
          }
        }

        return member
      });
      namespaces[mesh.name] = [...namespaces[mesh.name]];
      activeNamespaces[mesh.name] = namespaces[mesh.name][0] || "";

      return {
        ...mesh,
        members : processedMember
      }
    });

    SetMeshScan (processedControlPlanesData?.filter(data => !!data).filter((data) => data.members?.length > 0)) 
    SetMeshScanNamespaces(namespaces)
    SetActiveMeshScanNamespace(activeNamespaces)
  };
  

  return (
    <Grid container spacing={2}>
      <Grid item xs={6}>
        <PaperWithTitle title="Service Mesh" titleVariant="h6">
          <showServiceMesh meshScan={meshScan} activeMeshScanNamespace={activeMeshScanNamespace} />
        </PaperWithTitle>
      </Grid>
      <Grid item xs={6}>
        <ConnectionStatus />
      </Grid>
    </Grid>
  );
}
