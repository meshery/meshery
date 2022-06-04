import { 
  Chip, Grid, IconButton, List, ListItem, ListItemText, Menu, MenuItem, Switch, 
  Tooltip, Paper, NoSsr, TableCell, TableRow, TableContainer, Table
 } from '@material-ui/core';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import { useState } from 'react';
import DataTable from "mui-datatables";
import { withStyles } from '@material-ui/styles';

const styles = (theme) => ({
  operationButton : {
    [theme.breakpoints.down(1180)] : {
      marginRight : "25px",
    },
  },
  icon : { width : theme.spacing(2.5), },
  paper : { padding : theme.spacing(2), },
  heading : { textAlign : "center", },
  configBoxContainer : {
    [theme.breakpoints.down(1050)] : {
      flexGrow : 0,
      maxWidth : '100%',
      flexBasis : '100%',
    },
    [theme.breakpoints.down(1050)] : {
      flexDirection : "column",
    },
  },
  clusterConfiguratorWrapper : { padding : theme.spacing(5), display : "flex" },
  contentContainer : {
    [theme.breakpoints.down(1050)] : {
      flexDirection : "column",
    },
    flexWrap : "noWrap",
  },
  paper : { margin : theme.spacing(2), },
});

function MesherySettingsNew({ classes }) {
  const data = ["","","","","",""];
  const [showMenu, setShowMenu] = useState([false])
  const [anchorEl, setAnchorEl] = useState(null);
  const [operatorInstalled, setOperatorInstalled] = useState(false);
  const [meshSyncInstalled, setMeshSyncInstalled] = useState(false);
  const [meshSyncState, setMeshSyncState] = useState("N/A");
  const [NATSState, setNATSState] = useState("UNKNOWN");
  const [NATSVersion, setNATSVersion] = useState("N/A");
  const [meshSyncVersion, setMeshSyncVersion] = useState("N/A");
  const [operatorVersion, setOperatorVersion] = useState("N/A");
  const [operatorProcessing, setOperatorProcessing] = useState(false);
  const [operatorSwitch, setOperatorSwitch] = useState(false);

  const handleMenuClose = (index) => {
    let menu = [...showMenu];
    menu[index] = false;
    setShowMenu(menu)
  }

  const handleMenuOpen = (e, index) => {
    setAnchorEl(e.currentTarget)
    let menu = [...showMenu];
    menu[index] = true;
    setShowMenu(menu);
  }

  const columns = [
    {
      name : "Contexts",
      label : "Contexts",
      options : {
        filter : true,
        sort : true,
        searchable : true,
      }
    },
    {
      name : "deployment_type",
      label : "Type of Deployment",
      options : {
        filter : true,
        sort : true,
        searchable : true,
      }
    },
    {
      name : "location",
      label : "Location",
      options : {
        filter : true,
        sort : true,
        searchable : true,
      }
    },
    {
      name : "last_discovery",
      label : "Last Discovery",
      options : {
        filter : true,
        sort : true,
        searchable : true,
      }
    },
    {
      name : "Actions",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customBodyRender : (value, tableMeta, updateValue) => {
          return (
            <div>
              <IconButton
                aria-label="more"
                id="long-button"
                aria-controls={showMenu[tableMeta.rowIndex] ? 'long-menu' : undefined}
                aria-expanded={showMenu[tableMeta.rowIndex] ? 'true' : undefined}
                aria-haspopup="true"
                onClick={(e) => handleMenuOpen(e, tableMeta.rowIndex)}
              >
                <MoreVertIcon />
              </IconButton>
              <Menu
                id="long-menu"
                MenuListProps={{
                  'aria-labelledby' : 'long-button',
                }}
                anchorEl={anchorEl}
                open={showMenu[tableMeta.rowIndex]}
                onClose={() => handleMenuClose(tableMeta.rowIndex)}
                PaperProps={{
                  style : {
                    maxHeight : 48 * 4.5,
                    width : '20ch',
                  },
                }}
              >
                <MenuItem onClick={() => {
                  console.log("FMS"); handleMenuClose(tableMeta.rowIndex)
                }}>
                  FlushMeshSync
                </MenuItem>
                <MenuItem  onClick={() => {
                  console.log("RDO"); 
                  // handleMenuClose(tableMeta.rowIndex)
                }}>
                  <Switch
                    checked={operatorSwitch}
                    onClick={() => console.log("handle operator switch")}
                    disabled={operatorProcessing}
                    name="OperatorSwitch"
                    color="primary"
                  />
                </MenuItem>
              </Menu>
            </div>
          )
        },
      },

    }
  ]

  const options = {
    selectToolbarPlacement : 'none',
    print : false,
    download : false,
    selectableRows : false,
    expandableRows : true,
    expandableRowsOnClick : false,
    renderExpandableRow : (rowData, rowMetaData) => {
      console.log("RD", rowData, "RM", rowMetaData);
      return (
        <NoSsr>
            <TableCell colSpan={6}>
              <TableContainer>
                <Table>
                  <TableRow>
                    <TableCell className={classes.configBoxContainer}>
                      <Paper >
                        <div>
                          <Grid container spacing={1} >
                            <Grid item xs={12} md={5} className={classes.operationButton}>
                              <List>
                                <ListItem>
                                  <Tooltip
                                    title={operatorInstalled
                                      ? `Version: ${operatorVersion}`
                                      : "Not Available"}
                                    aria-label="meshSync"
                                  >
                                    <Chip
                                    // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
                                      label={"Operator"}
                                      // onDelete={self.handleReconfigure}
                                      onClick={self.handleOperatorClick}
                                      icon={<img src="/static/img/meshery-operator.svg" className={classes.icon} />}
                                      variant="outlined"
                                      data-cy="chipOperator"
                                    />
                                  </Tooltip>
                                </ListItem>
                              </List>
                            </Grid>
                            {operatorInstalled &&
                            <>
                              <Grid item xs={12} md={4}>
                                <List>
                                  <ListItem>
                                    <Tooltip
                                      title={meshSyncInstalled ? `Redeploy MeshSync` : "Not Available"}
                                      aria-label="meshSync"
                                    >
                                      <Chip
                                        label={"MeshSync"}
                                        onClick={() => console.log("HANDLE MESHSYNC CLICK")}
                                        icon={<img src="/static/img/meshsync.svg" className={classes.icon} />}
                                        variant="outlined"
                                        data-cy="chipMeshSync"
                                      />
                                    </Tooltip>
                                  </ListItem>
                                </List>
                              </Grid>
                              <Grid item xs={12} md={4}>
                                <List>
                                  <ListItem>
                                    <Tooltip
                                      title={NATSState === "CONNECTED" ? `Reconnect NATS` : "Not Available"}
                                      aria-label="nats"
                                    >
                                      <Chip
                                        label={"NATS"}
                                        onClick={() => console.log("HANDLE NATS CLICK")}
                                        icon={<img src="/static/img/nats-icon-color.svg" className={classes.icon} />}
                                        variant="outlined"
                                        data-cy="chipNATS"
                                      />
                                    </Tooltip>
                                  </ListItem>
                                </List>
                              </Grid>
                            </>
                            }
                          </Grid>
                          <Grid container spacing={1} className={classes.contentContainer}>
                            <Grid item xs={12} md={5}>
                              <List>
                                <ListItem>
                                  <ListItemText primary="Operator State" secondary={operatorInstalled
                                    ? "Active"
                                    : "Disabled"} />
                                </ListItem>
                                <ListItem>
                                  <ListItemText primary="Operator Version" secondary={operatorVersion} />
                                </ListItem>
                              </List>
                            </Grid>
                            <Grid item xs={12} md={5}>
                              <List>
                                <ListItem>
                                  <ListItemText primary="MeshSync State" secondary={meshSyncInstalled
                                    ? meshSyncState
                                    : "Disabled"} />
                                </ListItem>
                                <ListItem>
                                  <ListItemText primary="MeshSync Version" secondary={meshSyncVersion} />
                                </ListItem>
                              </List>
                            </Grid>
                            <Grid item xs={12} md={5}>
                              <List>
                                <ListItem>
                                  <ListItemText primary="NATS State" secondary={NATSState} />
                                </ListItem>
                                <ListItem>
                                  <ListItemText primary="NATS Version" secondary={NATSVersion} />
                                </ListItem>
                              </List>
                            </Grid>
                          </Grid>
                        </div>
                      </Paper>
                    </TableCell>
                    <TableCell>
                      <Paper >
                        <div>
                          <Grid container spacing={1} >
                            <Grid item xs={12} md={5} className={classes.operationButton}>
                              <List>
                                <ListItem>
                                  <Tooltip
                                    title={operatorInstalled
                                      ? `Version: ${operatorVersion}`
                                      : "Not Available"}
                                    aria-label="meshSync"
                                  >
                                    <Chip
                                    // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
                                      label={"Operator"}
                                      // onDelete={self.handleReconfigure}
                                      onClick={self.handleOperatorClick}
                                      icon={<img src="/static/img/meshery-operator.svg" className={classes.icon} />}
                                      variant="outlined"
                                      data-cy="chipOperator"
                                    />
                                  </Tooltip>
                                </ListItem>
                              </List>
                            </Grid>
                            {operatorInstalled &&
                            <>
                              <Grid item xs={12} md={4}>
                                <List>
                                  <ListItem>
                                    <Tooltip
                                      title={meshSyncInstalled ? `Redeploy MeshSync` : "Not Available"}
                                      aria-label="meshSync"
                                    >
                                      <Chip
                                        label={"MeshSync"}
                                        onClick={() => console.log("HANDLE MESHSYNC CLICK")}
                                        icon={<img src="/static/img/meshsync.svg" className={classes.icon} />}
                                        variant="outlined"
                                        data-cy="chipMeshSync"
                                      />
                                    </Tooltip>
                                  </ListItem>
                                </List>
                              </Grid>
                              <Grid item xs={12} md={4}>
                                <List>
                                  <ListItem>
                                    <Tooltip
                                      title={NATSState === "CONNECTED" ? `Reconnect NATS` : "Not Available"}
                                      aria-label="nats"
                                    >
                                      <Chip
                                        label={"NATS"}
                                        onClick={() => console.log("HANDLE NATS CLICK")}
                                        icon={<img src="/static/img/nats-icon-color.svg" className={classes.icon} />}
                                        variant="outlined"
                                        data-cy="chipNATS"
                                      />
                                    </Tooltip>
                                  </ListItem>
                                </List>
                              </Grid>
                            </>
                            }
                          </Grid>
                          <Grid container spacing={1} className={classes.contentContainer}>
                            <Grid item xs={12} md={5}>
                              <List>
                                <ListItem>
                                  <ListItemText primary="Operator State" secondary={operatorInstalled
                                    ? "Active"
                                    : "Disabled"} />
                                </ListItem>
                                <ListItem>
                                  <ListItemText primary="Operator Version" secondary={operatorVersion} />
                                </ListItem>
                              </List>
                            </Grid>
                            <Grid item xs={12} md={5}>
                              <List>
                                <ListItem>
                                  <ListItemText primary="MeshSync State" secondary={meshSyncInstalled
                                    ? meshSyncState
                                    : "Disabled"} />
                                </ListItem>
                                <ListItem>
                                  <ListItemText primary="MeshSync Version" secondary={meshSyncVersion} />
                                </ListItem>
                              </List>
                            </Grid>
                            <Grid item xs={12} md={5}>
                              <List>
                                <ListItem>
                                  <ListItemText primary="NATS State" secondary={NATSState} />
                                </ListItem>
                                <ListItem>
                                  <ListItemText primary="NATS Version" secondary={NATSVersion} />
                                </ListItem>
                              </List>
                            </Grid>
                          </Grid>
                        </div>
                      </Paper>
                    </TableCell>
                  </TableRow>
                </Table>
              </TableContainer>
            </TableCell>
        </NoSsr>
      )
    }
  }
  return (
    <DataTable
      columns = { columns }
      data = { data }
      options = { options }
    />
  )
}

export default withStyles(styles)(MesherySettingsNew);