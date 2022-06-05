import {
  Chip, Grid, IconButton, List, ListItem, ListItemText, Menu, MenuItem, Switch,
  Tooltip, Paper, NoSsr, TableCell, TableContainer, Table, Button, Typography,
  TextField, FormGroup, InputAdornment
} from '@material-ui/core';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import CloseIcon from "@material-ui/icons/Close";
import { withSnackbar } from "notistack";
import { useState, useEffect, useRef } from 'react';
import DataTable from "mui-datatables";
import { withStyles } from '@material-ui/styles';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../lib/store';
import subscribeOperatorStatusEvents from './graphql/subscriptions/OperatorStatusSubscription';
import subscribeMeshSyncStatusEvents from './graphql/subscriptions/MeshSyncStatusSubscription';
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import fetchMesheryOperatorStatus from './graphql/queries/OperatorStatusQuery';
import PromptComponent from './PromptComponent';
import CloudUploadIcon from "@material-ui/icons/CloudUpload";
import MeshsyncStatusQuery from './graphql/queries/MeshsyncStatusQuery';
import NatsStatusQuery from './graphql/queries/NatsStatusQuery';

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
  fileInputStyle : { display : "none", },
  button : { marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1), },
  grey : { background : "WhiteSmoke",
    padding : theme.spacing(2),
    borderRadius : "inherit", },
  fileLabelText : { cursor : "pointer",
    "& *" : { cursor : "pointer", }, },
  subtitle : {
    minWidth : 400,
    overflowWrap : 'anywhere',
    textAlign : 'left',
    padding : '5px'
  },
});

function MesherySettingsNew({ classes, enqueueSnackbar, closeSnackbar, updateProgress, k8sconfig }) {
  const [data, setData] = useState([])
  let k8sfileElementVal ="";
  const [showMenu, setShowMenu] = useState([false])
  const [anchorEl, setAnchorEl] = useState(null);
  // const [contexts, setContexts] = useState([]);vd` v
  const [operatorInstalled, setOperatorInstalled] = useState(false);
  const [meshSyncInstalled, setMeshSyncInstalled] = useState(false);
  const [meshSyncState, setMeshSyncState] = useState("N/A");
  const [NATSState, setNATSState] = useState("UNKNOWN");
  const [NATSVersion, setNATSVersion] = useState("N/A");
  const [meshSyncVersion, setMeshSyncVersion] = useState("N/A");
  const [operatorVersion, setOperatorVersion] = useState("N/A");
  const [operatorProcessing, setOperatorProcessing] = useState(false);
  const [operatorSwitch, setOperatorSwitch] = useState(false);
  const [meshSyncStatusSubscription, setMeshSyncStatusSubscription] = useState(null);
  const [operatorStatusSubscription, setOperatorStatusSubscription] = useState(null);
  const [lastDiscover, setLastDiscover] = useState([""]);
  const [contexts, setContexts] = useState([]);
  // const [toUploadContexts, setToUploadContexts] = useState([]);
  // const [open, setOpen] = useState(false);
  const [formData, setFormData] = useState(null);
  const ref = useRef(null);

  const dateOptions = { weekday : 'long', year : 'numeric', month : 'long', day : 'numeric' };

  useEffect(() => {
    console.log("ICG", k8sconfig);
    let meshSyncStatusEventsSubscription = subscribeMeshSyncStatusEvents((res) => {
      if (res.meshsync?.error) {
        handleError(res.meshsync?.error?.description || "MeshSync could not be reached");
        return;
      }
    });

    let tableInfo = [];
    fetchAllContexts(25)
      .then(res => {
        console.log(res, "CTX");
        setContexts(res.contexts);
        res.contexts.forEach((ctx) => {
          let updatedAt = new Date(ctx.updated_at);
          let data = {
            context : ctx.name,
            location : ctx.server,
            deployment_type : "fix",
            // last_discovery : lastDiscover[idx],
            last_discovery : updatedAt.toLocaleDateString("en-US", dateOptions) +
                             + " " + updatedAt.toLocaleTimeString("en-US"),
            name : ctx.name,
            id : ctx.id
          };
          tableInfo.push(data);
        })
        setData(tableInfo);
      })
      .catch(handleError("failed to fetch contexts for the instance"))

    let operatorStatusEventsSubscription = subscribeOperatorStatusEvents(setOperatorState);
    setOperatorStatusSubscription(operatorStatusEventsSubscription);
    fetchMesheryOperatorStatus().subscribe({
      next : (res) => {
        setOperatorState(res);
      },
      error : (err) => console.log("error at operator scan: " + err),
    });

    setMeshSyncStatusSubscription(meshSyncStatusEventsSubscription);

    return () => {
      if (meshSyncStatusSubscription) {
        meshSyncStatusSubscription.dispose();
      }
      if (operatorStatusSubscription) {
        operatorStatusSubscription.dispose();
      }
    }
  }, [])

  const handleMenuClose = (index) => {
    let menu = [...showMenu];
    menu[index] = false;
    setShowMenu(menu)
  }

  const setOperatorState = (res) => {
    console.log("incoming change")
    if (res.operator?.error) {
      handleError("Operator could not be reached")(res.operator?.error?.description);
      setOperatorProcessing(false);
      return false;
    }

    if (res.operator?.status === "ENABLED") {
      setOperatorProcessing(false);
      res.operator?.controllers?.forEach((controller) => {
        if (controller.name === "broker" && controller.status.includes("CONNECTED")) {
          setNATSState(controller.status);
          setNATSVersion(controller.version);
        } else if (controller.name === "meshsync" && controller.status.includes("ENABLED")) {
          setMeshSyncInstalled(true);
          setMeshSyncVersion(controller.version);
          setMeshSyncState(controller.status);
        }
      });
      setOperatorInstalled(true);
      setOperatorSwitch(true);
      setOperatorVersion(res.operator?.version);
      return true;
    }
    if (res.operator?.status === "DISABLED") {
      setOperatorProcessing(false);
    }

    if (res.operator?.status === "PROCESSING") {
      console.log("setting to processing");
      setOperatorProcessing(true);
    }
  }

  async function fetchAllContexts(number) {
    return await promisifiedDataFetch("/api/system/kubernetes/contexts?pageSize=" + number)
  }

  const handleError = (msg) => (error) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar(`${msg}: ${error}`, { variant : "error", preventDuplicate : true,
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000, });
  };

  const handleMenuOpen = (e, index) => {
    setAnchorEl(e.currentTarget)
    let menu = [...showMenu];
    menu[index] = true;
    setShowMenu(menu);
    console.log(contexts, "CTX");
  }

  const handleSuccess = msg => {
    updateProgress({ showProgress : false });
    enqueueSnackbar(msg, {
      variant : "success",
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000,
    });
  }

  const handleLastDiscover = (index) => {
    let ld = lastDiscover;
    let dt = new Date();

    ld[index] = dt.toLocaleDateString("en-US", dateOptions) + "  " + dt.toLocaleTimeString("en-US");
    setLastDiscover(ld);
  }

  const handleKubernetesClick = (context, index) => {
    updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/kubernetes/ping?context=" + context,
      { credentials : "same-origin" },
      (result) => {
        updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          handleLastDiscover(index);
          console.log("ld", lastDiscover);
          enqueueSnackbar("Kubernetes was successfully pinged!", {
            variant : "success",
            "data-cy" : "k8sSuccessSnackbar",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
          });
        }
      },
      handleError("Kubernetes config could not be validated")
    );
  };

  const handleConfigDelete = (id, index) => {
    console.log("test");
    updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/kubernetes/contexts/" + id,
      { credentials : "same-origin",
        method : "DELETE" },
      (result) => {
        updateProgress({ showProgress : false });
        if (index != undefined) {
          let newData = data.filter((dt, idx) => index != idx);
          console.log(index, "LL");
          setData(newData);
          console.log("result", result);
        }
      },
      handleError("failed to delete kubernetes context")
    );
  }

  const handleChange = () => {
    const field = document.getElementById("k8sfile");
    const textField = document.getElementById("k8sfileLabelText");
    if (field instanceof HTMLInputElement) {
      if (field.files.length < 1) return;
      const name = field.files[0].name;
      const formData = new FormData();
      formData.append("k8sfile", field.files[0])
      textField.value=name;
      setFormData(formData);
    }
  }

  const uploadK8SConfig = async () => {
    await promisifiedDataFetch(
      "/api/system/kubernetes",
      {
        method : "POST",
        body : formData,
      }
    )
    console.log("changed");
  }
  const columns = [
    {
      name : "contexts",
      label : "Contexts",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customBodyRender : (value, tableMeta, updateValue) => {
          console.log("v", value, "tm", tableMeta, "uv", updateValue);
          return (
            <Tooltip title={`Server: ${tableMeta.rowData[2]}`}>
              <Chip
                label={data[tableMeta.rowIndex].name}
                onDelete={() => handleConfigDelete(data[tableMeta.rowIndex].id, tableMeta.rowIndex)}
                onClick={() => handleKubernetesClick(data[tableMeta.rowIndex].id, tableMeta.rowIndex)}
                icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
                variant="outlined"
                data-cy="chipContextName"
              />
            </Tooltip>
          )
        }
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
        customBodyRender : (value, tableMeta) => {
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
    print : false,
    download : false,
    expandableRows : true,
    expandableRowsOnClick : false,
    onRowsDelete : (td) => {
      td.data.forEach((item) => {
        console.log("D", item);
        handleConfigDelete(data[item.index].id)
      })
    },
    renderExpandableRow : (rowData, rowMetaData) => {
      console.log("RDASA", rowData, "RM", rowMetaData);
      return (
        <NoSsr>
          <TableCell colSpan={6}>
            <TableContainer>
              <Table>
                {/* <TableRow> */}
                <TableCell>
                  <Paper >
                    <div>
                      <Grid container spacing={1} >
                        <Grid item xs={12} md={5} className={classes.operationButton}>
                          <List>
                            <ListItem>
                              <Tooltip title={`Server: ${rowData[2].server}`}
                              >
                                <Chip
                                  label={data[rowMetaData.rowIndex].name}
                                  onClick={() => handleKubernetesClick(data[rowMetaData.rowIndex].id, rowMetaData.rowIndex)}
                                  icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
                                  variant="outlined"
                                  data-cy="chipContextName"
                                />
                              </Tooltip>
                            </ListItem>
                          </List>
                        </Grid>
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
                <TableCell className={classes.configBoxContainer}>
                  <Paper >
                    <div>
                      <Grid container spacing={1} >
                        <Grid item xs={12} md={4} className={classes.operationButton}>
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
                                  // onDelete={handleReconfigure}
                                  onClick={handleOperatorClick}
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
                                        onClick={handleMeshSyncClick}
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
                                        onClick={handleNATSClick}
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
                {/* </TableRow> */}
              </Table>
            </TableContainer>
          </TableCell>
        </NoSsr>
      )
    }
  }

  const handleClick = async () => {
    const modal = ref.current;
    let response = await modal.show({
      title : "Add Kuberneted Cluster(s)",
      subtitle :
      <>
        <div>
          <Typography variant="h6">
          Upload your kubeconfig
          </Typography>
          <Typography variant="body2">
          commonly found at ~/.kube/config
          </Typography>
          <FormGroup>
            <input
              id="k8sfile"
              type="file"
              value={k8sfileElementVal}
              onChange={handleChange}
              className={classes.fileInputStyle}
            />
            <TextField
              id="k8sfileLabelText"
              name="k8sfileLabelText"
              className={classes.fileLabelText}
              label="Upload kubeconfig"
              variant="outlined"
              fullWidth
              onClick={() => {
                document.querySelector("#k8sfile")?.click(); setOpen(true);
              }}
              margin="normal"
              InputProps={{ readOnly : true,
                endAdornment : (
                  <InputAdornment position="end">
                    <CloudUploadIcon />
                  </InputAdornment>
                ), }}
            />
          </FormGroup>
          {/* <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
            className={classes.dialogBox}
          >
            <DialogContent>
              <DialogContentText className={ classes.subtitle }>
                <Typography>
                    Available Contexts
                </Typography>
                {
                  toUploadContexts.map((ctx) => (
                    <Chip
                      label={ctx.name}
                      // onDelete={() => handleReconfigure(ctx.id)}
                      onClick={() => handleKubernetesClick(data[tableMeta.rowIndex].id, tableMeta.rowIndex)}
                      icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
                      variant="outlined"
                      data-cy="chipContextName"
                    />
                  ))
                }
              </DialogContentText>
            </DialogContent>
          </Dialog> */}
        </div>
      </>,
      options : ["Cancel", "Upload"]
    })
    if (response === "Upload") {
      uploadK8SConfig().then(() => {
        handleSuccess("successfully uploaded kubernetes config");
        fetchAllContexts(25)
          .then(res => setContexts(res.contexts))
          .catch(handleError("failed to get contexts"))
      }).
        catch(err => {
          handleError("failed to upload kubernetes config")(err)
        })
    }
  }

  const handleOperatorClick = () => {
    updateProgress({ showProgress : true });
    fetchMesheryOperatorStatus().subscribe({ next : (res) => {
      console.log(res);
      let state = setOperatorState(res);
      updateProgress({ showProgress : false });
      if (state == true) {
        enqueueSnackbar("Operator was successfully pinged!", { variant : "success",
          autoHideDuration : 2000,
          action : (key) => (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ), });
      } else {
        handleError("Operator could not be reached")("Operator is disabled");
      }
    },
    error : handleError("Operator could not be pinged"), });
  };

  const handleNATSClick = () => {
    updateProgress({ showProgress : true });
    NatsStatusQuery().subscribe({
      next : (res) => {
        updateProgress({ showProgress : false });
        if (res.controller.name === "broker" && res.controller.status.includes("CONNECTED")) {
          let runningEndpoint = res.controller.status.substring("CONNECTED".length)
          enqueueSnackbar(`Broker was successfully pinged. Running at ${runningEndpoint}`, {
            variant : "success",
            action : (key) => (
              <IconButton key="close" aria-label="close" color="inherit" onClick={() => closesnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
            autohideduration : 2000,
          })
        } else {
          handleError("Meshery Broker could not be reached")("Meshery Server is not connected to Meshery Broker");
        }
        setNATSState(res.controller.status.length !== 0 ? res.controller.status : "UNKNOWN");
        setNATSVersion(res.controller.version);
      },
      error : handleError("NATS status could not be retrieved"), });

    // connectToNats().subscribe({
    //   next : (res) => {
    //     if (res.connectToNats === "PROCESSING") {
    //       updateProgress({ showProgress : false });
    //       enqueueSnackbar(`Reconnecting to NATS...`, {
    //         variant : "info",
    //         action : (key) => (
    //           <IconButton key="close" aria-label="close" color="inherit" onClick={() => closesnackbar(key)}>
    //             <CloseIcon />
    //           </IconButton>
    //         ),
    //         autohideduration : 7000,
    //       })
    //     }
    //     if (res.connectToNats === "CONNECTED") {
    //       updateProgress({ showProgress : false });
    //       enqueueSnackbar(`Successfully connected to NATS`, {
    //         variant : "success",
    //         action : (key) => (
    //           <IconButton key="close" aria-label="close" color="inherit" onClick={() => closesnackbar(key)}>
    //             <CloseIcon />
    //           </IconButton>
    //         ),
    //         autohideduration : 7000,
    //       })
    //     }

    //   },
    //   error : handleError("Failed to request reconnection with NATS"),
    // });

  };

  const handleMeshSyncClick = () => {
    updateProgress({ showProgress : true });
    MeshsyncStatusQuery().subscribe({ next : (res) => {
      updateProgress({ showProgress : false });
      if (res.controller.name === "meshsync" && res.controller.status.includes("ENABLED")) {
        setMeshSyncInstalled(true);
        setMeshSyncVersion(res.controller.version);
        setMeshSyncState(res.controller.status);
        let publishEndpoint = res.controller.status.substring("ENABLED".length)
        enqueueSnackbar(`MeshSync was successfully pinged. Publishing to ${publishEndpoint} `, {
          variant : "success",
          action : (key) => (
            <IconButton key="close" aria-label="close" color="inherit" onClick={() => closesnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ),
          autohideduration : 2000,
        })
      } else {
        handleError("MeshSync could not be reached")("MeshSync is unavailable");
        setMeshSyncInstalled(false);
        setMeshSyncVersion("");
        setMeshSyncState(res.controller.status);
      }
    },
    error : handleError("MeshSync status could not be retrieved"), });

    // deployMeshSync().subscribe({
    //   next : (res) => {
    //     if (res.deployMeshsync === "PROCESSING") {
    //       updateProgress({ showProgress : false });
    //       enqueueSnackbar(`MeshSync deployment in progress`, {
    //         variant : "info",
    //         action : (key) => (
    //           <IconButton key="close" aria-label="close" color="inherit" onClick={() => closesnackbar(key)}>
    //             <CloseIcon />
    //           </IconButton>
    //         ),
    //         autohideduration : 7000,
    //       })
    //     }

    //   },
    //   error : handleError("Failed to request Meshsync redeployment"),
    // });
  };

  return (
    <>
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        onClick={handleClick}
        className={classes.button}
        data-cy="btnResetDatabase"
      >
      Add Cluster
      </Button>

      <DataTable
        columns = { columns }
        data = { data }
        options = { options }
      />
      <PromptComponent ref={ ref }/>
    </>
  )
}
const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig').toJS();
  return k8sconfig;
}
const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesherySettingsNew)));

