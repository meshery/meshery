import React, {useState} from "react"
import {Grid, Typography, Box, Paper, MenuItem, Button, TableContainer, Table, TableHead, TableRow, TableCell, Tooltip,} from "@mui/material"
import {  getK8sClusterNamesFromCtxId } from "@/utils//multi-ctx";
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { styled } from "@mui/material/styles";
import {useRouter} from "next/router";

export const TableWrapper = styled(Box)(({theme})=>({
  fontSize : "15px", 
  color : '#fff', 
  paddingBottom : '10px', 
  padding : '1vh'
}))

function ShowServiceMesh ({meshScan, activeMeshScanNamespace, meshScanNamespaces}) {

  const router = useRouter();

    function emptyStateMessageForServiceMeshesInfo () {
        const clusters = getSelectedK8sContextsNames();
        if (clusters.length === 0) {
          return "No Cluster is selected to show the Service Mesh Information"
        }
        if (clusters.includes("all")) {
          return `No service meshes detected in any of the cluster.`
        }
        return `No service meshes detected in the ${clusters.join(", ")} cluster(s).`
      }

     function getSelectedK8sContextsNames ()  {
        return getK8sClusterNamesFromCtxId()
      } 
    
      function Meshcard (mesh, components = []) {
  
        if (Array.isArray(components) && components.length)
          return (
            <Paper elevation={1} sx={{ padding : "2rem", marginTop : "1rem" }}>
              <Grid container justify="space-between" spacing={1}>
                <Grid item>
                  <div style={{ display : "flex", alignItems : "center", marginBottom : "1rem" }}>
                    <img src={mesh.icon} style={{ marginRight : "0.75rem" }} />
                    <Typography variant="h6">{mesh.tag}</Typography>
                  </div>
                </Grid>
                <Grid item>
                  {activeMeshScanNamespace[mesh.name] && (
                    <Select
                      value={activeMeshScanNamespace[mesh.name]}
                    >
                      {meshScanNamespaces[mesh.name] &&
                        meshScanNamespaces[mesh.name].map((ns) => <MenuItem value={ns}>{ns}</MenuItem>)}
                    </Select>
                  )}
                </Grid>
              </Grid>
              <TableContainer>
                <Table aria-label="Deployed service mesh details">
                  <TableHead>
                    <TableRow>
                      <TableCell align="center">Component</TableCell>
                      <TableCell align="center">Version</TableCell>
                      <TableCell align="center">Proxy</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {components
                      .filter((comp) => comp.namespace === activeMeshScanNamespace[mesh.name])
                      .map((component) => {
                        return (
                          <TableRow key={component.name.full}>
                            <TableCell align="center">{component.component}</TableCell>
                            <Tooltip
                              key={`component-${component.name}`}
                              title={
                                Array.isArray(component?.data_planes) && component.data_planes.length > 0 ? (
                                  component.data_planes.map((cont) => {
                                    return (
                                      <TableWrapper key={cont.name}>
                                        <p>Name: {cont?.containerName ? cont.containerName : 'Unspecified'}</p>
                                        <p>Status: {cont?.status?.ready ? 'ready' : 'not ready'}</p>
                                        {!cont?.status?.ready && (
                                          typeof cont?.status?.lastState === 'object' && cont?.status?.lastState !== null && Object.keys(cont.status.lastState).length > 0 && (
                                            <div>
                                              <p>Last state: {Object.keys(cont?.status?.lastState)[0]} <br /> Error: {Object.values(cont?.status?.lastState)[0]?.exitCode} <br /> Finished at: {Object.values(cont?.status?.lastState)[0]?.finishedAt}</p>
                                            </div>
                                          )
                                        )}
                                        {typeof cont?.status?.state === 'object' && cont?.status?.state !== null && Object.keys(cont.status.state).length > 0 && (
                                          <p>State: {Object.keys(cont.status.state)[0]}</p>
                                        )}
                                        {cont?.status?.restartCount && (
                                          <p>Restart count: {cont?.status.restartCount}</p>
                                        )}
                                        <p>Image: {cont.image}</p>
                                        <p>Ports: <br /> {cont?.ports && cont.ports.map(port => `[ ${port?.name ? port.name : 'Unknown'}, ${port?.containerPort ? port.containerPort : 'Unknown'}, ${port?.protocol ? port.protocol : 'Unknown'} ]`).join(', ')}</p>
                                        {cont?.resources && (
                                          <div>
                                            Resources used: <br />
      
                                            <div style={{ paddingLeft : '2vh' }}>
                                              {cont?.resources?.limits && (
                                                <div>
                                                  <p>Limits: <br />
                                                    CPU: {cont?.resources?.limits?.cpu} - Memory: {cont?.resources?.limits?.memory}</p>
                                                </div>
                                              )}
                                              {cont?.resources?.requests && (
                                                <div>
                                                  <p>Requests: <br />
                                                    CPU: {cont?.resources?.requests?.cpu} - Memory: {cont?.resources?.requests?.memory}</p>
                                                </div>
                                              )}
                                            </div>
                                          </div>
                                        )}
                                      </TableWrapper>
                                    )
                                  })
                                ) : "No proxy attached"}
                            >
                              <TableCell align="center">{component?.data_planes?.length || 0}</TableCell>
                            </Tooltip>
                          </TableRow>
                        )
                      })}
                  </TableBody>
                </Table>
              </TableContainer>
            </Paper>
          );
      
        return null;
      };
       

    return(
  <>
          {meshScan && Object.keys(meshScan).length
            ? (
              <>
                {meshScan.map((mesh) => {
                  let tag = "";
                  mesh.name
                    .split("_")
                    .forEach((element) => {
                      tag = tag + " " + element[0].toUpperCase() + element.slice(1, element.length);
                    });
                  return Meshcard(
                    { name : mesh.name, tag : tag, icon : "/static/img/" + mesh.name + ".svg" },
                    mesh.members
                  );
                })}
              </>
            )
            : (
              <Box
              sx={{
                margin: "auto",
                display : "flex",
                justifyContent : "center",
                flexDirection : "column",
              }}
            >
                <Typography style={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" color="textSecondary">
                  {emptyStateMessageForServiceMeshesInfo}
                </Typography>
                <Button
                  aria-label="Add Meshes"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={() => router.push("/management")}
                >
                  <AddCircleOutlineIcon />
                  Install Service Mesh
                </Button>
              </Box>
            )}
        </>
    )
  }

  export default ShowServiceMesh