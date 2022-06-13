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
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import PromptComponent from './PromptComponent';
import CloudUploadIcon from "@material-ui/icons/CloudUpload";
import AddIcon from '@material-ui/icons/Add';
import MeshsyncStatusQuery from './graphql/queries/MeshsyncStatusQuery';
import NatsStatusQuery from './graphql/queries/NatsStatusQuery';
import changeOperatorState from './graphql/mutations/OperatorStatusMutation';
import resetDatabase from "./graphql/queries/ResetDatabaseQuery";
import { updateProgress, actionTypes, setMeshsyncSubscription } from "../lib/store";
import fetchMesheryOperatorStatus from "./graphql/queries/OperatorStatusQuery";


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
  button : {
    marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1),
    padding : theme.spacing(1),
    borderRadius : 5
  },
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
  text : {
    width : "80%",
    wordWrap : "break-word"
  },
  add : {
    marginRight : theme.spacing(1)
  },
  FlushBtn : {
    margin : theme.spacing(0.5),
    padding : theme.spacing(1),
    borderRadius : 5
  },
  menu : {
    display : 'flex',
    alignItems : 'center'
  }
});

function MesherySettingsNew({ classes, enqueueSnackbar, closeSnackbar, updateProgress,
  operatorState, MeshSyncState, setMeshsyncSubscription }) {
  const [data, setData] = useState([])
  const [showMenu, setShowMenu] = useState([false])
  const [anchorEl, setAnchorEl] = useState(null);
  const [operatorInstalled, setOperatorInstalled] = useState([false]);
  const [NATSState, setNATSState] = useState(["UNKNOWN"]);
  const [NATSVersion, setNATSVersion] = useState(["N/A"]);
  const [operatorVersion, setOperatorVersion] = useState(["N/A"]);
  const [operatorProcessing, setOperatorProcessing] = useState([false]);
  const [operatorSwitch, setOperatorSwitch] = useState([false]);
  const [discover, setLastDiscover] = useState(['']);
  const [contexts, setContexts] = useState([]);
  const [k8sVersion, setK8sVersion] = useState(["N/A"]);


  const ref = useRef(null);
  const meshSyncResetRef = useRef(null);
  // const systemResetRef = useRef(null);

  const dateOptions = { weekday : 'long', year : 'numeric', month : 'long', day : 'numeric' };

  let k8sfileElementVal ="";
  let formData = new FormData();
  const stateUpdater = (state, updateFunc, updateValue, index) => {
    let newState = [...state];
    newState[index] = updateValue;
    updateFunc(newState);
  }

  // done
  useEffect(() => {
    // let meshSyncStatusEventsSubscription = subscribeMeshSyncStatusEvents((res) => {
    //   if (res.meshsync?.error) {
    //     handleError(res.meshsync?.error?.description || "MeshSync could not be reached");
    //     return;
    //   }
    // });


    let tableInfo = [];
    fetchAllContexts(25)
      .then(res => {
        handleContexts(res.contexts);
        res.contexts.forEach((ctx) => {
          let data = {
            context : ctx.name,
            location : ctx.server,
            deployment_type : "-",
            last_discovery : "",
            name : ctx.name,
            id : ctx.id
          };
          tableInfo.push(data);
        })
        setData(tableInfo);
      })
      .catch(handleError("failed to fetch contexts for the instance"))
      // console.log(contexts[0].id, "CTX");
    getKubernetesVersion();

    // let operatorStatusEventsSubscription = subscribeOperatorStatusEvents(setOperatorState);
    // setOperatorStatusSubscription(operatorStatusEventsSubscription);
    // fetchMesheryOperatorStatus().subscribe({
    //   next : (res) => {
    //     setOperatorState(res);
    //   },
    //   error : (err) => console.log("error at operator scan: " + err),
    // });

    // setMeshSyncStatusSubscription(meshSyncStatusEventsSubscription);

  }, [])

  useEffect(() => {
    let opSwitch = [];
    operatorState?.forEach((state, idx) => {
      opSwitch[idx] = state.operatorStatus.status !== 'ENABLED' ? false : true;
      setOperatorState(state.operatorStatus, idx);
    })
    setOperatorSwitch(opSwitch);
  }, [operatorState])
  const isMeshSyncActive = (ctx) => {
    return MeshSyncState.filter((state) => state.contextID === ctx ).length > 0;
  }

  // const handleResetDatabase = () => {
  //   return async () => {
  //     systemResetRef.current.show({
  //       title : "Reset Meshery Database?",
  //       subtitle : "Are you sure to reset all the data of Meshery?",
  //       options : ["Proceed", "Cancel"]
  //     });
  //     if (responseOfResetDatabase === "Continue Reset") {
  //       this.props.updateProgress({ showProgress : true });
  //       const self = this;
  //       resetDatabase({
  //         selector : {
  //           clearDB : "true",
  //           ReSync : "true",
  //           hardReset : "true",
  //         },
  //       }).subscribe({
  //         next : (res) => {
  //           self.props.updateProgress({ showProgress : false });
  //           if (res.resetStatus === "PROCESSING") {
  //             this.props.enqueueSnackbar(`Database reset successful.`, {
  //               variant : "success",
  //               action : (key) => (
  //                 <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
  //                   <CloseIcon />
  //                 </IconButton>
  //               ),
  //               autohideduration : 3000,
  //             })
  //           }
  //         },
  //         error : self.handleError("Database is not reachable, try restarting server.")
  //       });
  //     }
  //   }
  // }

  const handleFlushMeshSync = (index) => {
    return async () => { // add backend support to delete context realted data
      handleMenuClose(index);
      let response = await meshSyncResetRef.current.show({
        title : "Flush MeshSync data?",
        subtitle : "Are you sure to Flush MeshSync data?",
        options : ["PROCEED", "CANCEL"]
      });
      if (response === "PROCEED") {
        updateProgress({ showProgress : true });
        resetDatabase({
          selector : {
            clearDB : "true",
            ReSync : "false",
            hardReset : "false",
          },
          k8scontextID : contexts[index].id
        }).subscribe({
          next : (res) => {
            updateProgress({ showProgress : false });
            if (res.resetStatus === "PROCESSING") {
              enqueueSnackbar(`Database reset successful.`, {
                variant : "success",
                action : (key) => (
                  <IconButton key="close" aria-label="close" color="inherit" onClick={() => closeSnackbar(key)}>
                    <CloseIcon />
                  </IconButton>
                ),
                autohideduration : 2000,
              })
            }
          },
          error : handleError("Database is not reachable, try restarting server.")
        });
      }
    }
  }

  const handleContexts = (contexts) => {
    contexts.forEach((ctx) => {
      let cdt = new Date(ctx.created_at);
      let updt = new Date(ctx.updated_at);

      ctx.created_at = cdt.toLocaleDateString("en-US", options)
      + " " +  cdt.toLocaleTimeString("en-US");

      ctx.updated_at = updt.toLocaleDateString("en-US", options)
      + " " +  updt.toLocaleTimeString("en-US")
    })
    setContexts(contexts);
  }

  const handleMenuClose = (index) => {
    let menu = [...showMenu];
    menu[index] = false;
    setShowMenu(menu)
  }

  const setOperatorState = (res, index) => {

    if (res.operator?.error) {
      handleError("Operator could not be reached")(res.operator?.error?.description);
      setMeshsyncSubscription( { action : actionTypes.SET_MESHSYNC_SUBSCRIPTION, meshSyncState : null })
      stateUpdater(operatorProcessing, setOperatorProcessing, false, index);
      return false;
    }

    if (res.operator?.status === "ENABLED") {
      stateUpdater(operatorProcessing, setOperatorProcessing, false, index);
      res.operator?.controllers?.forEach((controller) => {
        if (controller.name === "broker" && controller.status.includes("CONNECTED")) {
          stateUpdater(NATSState, setNATSState, controller.status, index);
          stateUpdater(NATSVersion, setNATSVersion, controller.version, index);

        }
      });
      stateUpdater(operatorInstalled, setOperatorInstalled, true, index);
      stateUpdater(operatorSwitch, setOperatorSwitch, true, index);
      stateUpdater(operatorVersion, setOperatorVersion, res.operator?.version, index);
      return true;
    }
    if (res.operator?.status === "DISABLED") {
      stateUpdater(operatorProcessing, setOperatorProcessing, false, index);
    }

    if (res.operator?.status === "PROCESSING") {
      console.log("setting to processing");
      stateUpdater(operatorProcessing, setOperatorProcessing, true, index);
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
    let dt = new Date();
    const newDate = dt.toLocaleDateString("en-US", dateOptions) + "  " + dt.toLocaleTimeString("en-US");
    let newData = [...discover];
    newData[index] = newDate;
    setLastDiscover(newData);
  }

  const getKubernetesVersion = () => {
    dataFetch(
      "/api/oam/workload/APIService.K8s",
      { credentials : "same-origin" },
      (result) => {
        if (result) {
          let version = result[0]?.oam_definition?.spec?.metadata?.version;
          setK8sVersion(version);
        }
      }
    )
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

  const handleOperatorSwitch = (index) => {

    const variables = {
      status : `${!operatorSwitch[index] ? "ENABLED" : "DISABLED"}`,
      contextID : contexts[index].id
    };

    updateProgress({ showProgress : true });

    changeOperatorState((response, errors) => {
      updateProgress({ showProgress : false });
      if (errors !== undefined) {
        handleError(`Unable to ${operatorSwitch[index] === true ? "Uni" : "I"} nstall operator`);
      }
      enqueueSnackbar("Operator " + response.operatorStatus.toLowerCase(), { variant : "success",
        autoHideDuration : 2000,
        action : (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        ), });
      stateUpdater(operatorSwitch, setOperatorSwitch, !operatorSwitch[index], index);
    }, variables);
  };

  const handleConfigDelete = (id, index) => {
    updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/kubernetes/contexts/" + id,
      { credentials : "same-origin",
        method : "DELETE" },
      () => {
        updateProgress({ showProgress : false });
        if (index != undefined) {
          let newData = data.filter((dt, idx) => index != idx);
          setData(newData);
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
      const formdata = new FormData();
      formdata.append("k8sfile", field.files[0])
      textField.value=name;
      formData = formdata;

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
  }

  const columns = [
    {
      name : "contexts",
      label : "Contexts",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customBodyRender : (_, tableMeta, ) => {
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
      },
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
        customBodyRender : (value, tableMeta) => <p>{ discover[tableMeta.rowIndex] }</p>
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
                className={classes.menu}
                id="long-menu"
                MenuListProps={{
                  'aria-labelledby' : 'long-button',
                }}
                anchorEl={anchorEl}
                open={showMenu[tableMeta.rowIndex]}
                onClose={() => handleMenuClose(tableMeta.rowIndex)}
              >
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={handleFlushMeshSync(tableMeta.rowIndex)}
                  className={classes.FlushBtn}
                  data-cy="btnResetDatabase"
                >
                  <Typography> Flush MeshSync </Typography>
                </Button>
                <MenuItem>
                  <Switch
                    checked={operatorSwitch[tableMeta.rowIndex]}
                    onClick={() => handleOperatorSwitch(tableMeta.rowIndex)}
                    disabled={operatorProcessing[tableMeta.rowIndex]}
                    name="OperatorSwitch"
                    color="primary"
                  />
                    Operator
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
        handleConfigDelete(data[item.index].id)
      })
    },
    renderExpandableRow : (rowData, rowMetaData) => {
      return (
        <NoSsr>
          { operatorState &&
          <TableCell colSpan={6}>
            <TableContainer>
              <Table>
                {/* <TableRow> */}
                <TableCell className={classes.configBoxContainer}>
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
                              <ListItemText primary="Name" secondary={contexts[rowMetaData.rowIndex].name}/>
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="K8s Version" secondary={k8sVersion} />
                            </ListItem>
                          </List>
                        </Grid>
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <ListItemText primary="Created At" secondary={
                                contexts[rowMetaData.rowIndex].created_at
                              }/>
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="UpdatedAt" secondary={
                                contexts[rowMetaData.rowIndex].updated_at
                              } />
                            </ListItem>
                          </List>
                        </Grid>
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <ListItemText className={classes.text} primary="Server" secondary={
                                contexts[rowMetaData.rowIndex].server
                              }/>
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
                                title={operatorInstalled[rowData.rowIndex]
                                  ? `Version: ${operatorVersion[rowMetaData.rowIndex]}`
                                  : "Not Available"}
                                aria-label="meshSync"
                              >
                                <Chip
                                  // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
                                  label={"Operator"}
                                  // onDelete={handleReconfigure}
                                  onClick={() => handleOperatorClick(rowMetaData.rowIndex)}
                                  icon={<img src="/static/img/meshery-operator.svg" className={classes.icon} />}
                                  variant="outlined"
                                  data-cy="chipOperator"
                                />
                              </Tooltip>
                            </ListItem>
                          </List>
                        </Grid>
                        {operatorSwitch[rowMetaData.rowIndex] &&
                            <>
                              <Grid item xs={12} md={4}>
                                <List>
                                  <ListItem>
                                    <Tooltip
                                      title={isMeshSyncActive(data[rowMetaData.rowIndex].id) ? `Redeploy MeshSync` : "Not Available"}
                                      aria-label="meshSync"
                                    >
                                      <Chip
                                        label={"MeshSync"}
                                        onClick={() => handleMeshSyncClick(rowMetaData.rowIndex)}
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
                                      title={NATSState[rowMetaData.rowIndex] === "CONNECTED" ? `Reconnect NATS` : "Not Available"}
                                      aria-label="nats"
                                    >
                                      <Chip
                                        label={"NATS"}
                                        onClick={() => handleNATSClick(rowMetaData.rowIndex)}
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
                              <ListItemText primary="Operator State" secondary={operatorSwitch[rowMetaData?.rowIndex]
                                ? "Active"
                                : "Disabled"} />
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="Operator Version" secondary={operatorState[rowMetaData?.rowIndex]?.operatorStatus.version} />
                            </ListItem>
                          </List>
                        </Grid>
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <ListItemText primary="MeshSync State" secondary={isMeshSyncActive(data[rowMetaData?.rowIndex].id)
                                ? "Enabled"
                                : "Disabled"} />
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="MeshSync Version" secondary={MeshSyncState[rowMetaData?.rowIndex]?.OperatorControllerStatus.version} />
                            </ListItem>
                          </List>
                        </Grid>
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <ListItemText primary="NATS State" secondary={NATSState[rowMetaData?.rowIndex]} />
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="NATS Version" secondary={NATSVersion[rowMetaData?.rowIndex]} />
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
          }
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
                document.querySelector("#k8sfile")?.click();
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
      options : ["CANCEL", "UPLOAD"]
    })
    if (response === "UPLOAD") {
      uploadK8SConfig().then(() => {
        handleSuccess("successfully uploaded kubernetes config");
        fetchAllContexts(25)
          .then(res => {
            let newData = [...data];
            setData(newData);
            setContexts(res.contexts)
          })
          .catch(handleError("failed to get contexts"))
      }).
        catch(err => {
          handleError("failed to upload kubernetes config")(err)
        })
    }
  }

  const handleOperatorClick = (index) => {
    updateProgress({ showProgress : true });
    fetchMesheryOperatorStatus({ k8scontextID : contexts[index].id })
      .subscribe({ next : (res) => {
        let state = setOperatorState(res, index);
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

  // done
  const handleNATSClick = (index) => {
    updateProgress({ showProgress : true });
    NatsStatusQuery({ k8scontextID : contexts[index].id }).subscribe({
      next : (res) => {
        updateProgress({ showProgress : false });
        if (res.controller.name === "broker" && res.controller.status.includes("CONNECTED")) {
          let runningEndpoint = res.controller.status.substring("CONNECTED".length)
          enqueueSnackbar(`Broker was successfully pinged. Running at ${runningEndpoint}`, {
            variant : "success",
            action : (key) => (
              <IconButton key="close" aria-label="close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
            autohideduration : 2000,
          })
        } else {
          handleError("Meshery Broker could not be reached")("Meshery Server is not connected to Meshery Broker");
        }
        // make state changes for individual index
        stateUpdater(NATSState, setNATSState, res.controller.status.length !== 0 ? res.controller.status : "UNKNOWN", index)
        // setNATSState(res.controller.status.length !== 0 ? res.controller.status : "UNKNOWN");
        // setNATSVersion(res.controller.version);
        stateUpdater(NATSVersion, setNATSVersion, res.controller.version, index);
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

  const handleMeshSyncClick = (index) => {
    updateProgress({ showProgress : true });
    MeshsyncStatusQuery(({ k8scontextID : contexts[index].id })).subscribe({ next : (res) => {
      updateProgress({ showProgress : false });
      if (res.controller.name !== "meshsync" || !res.controller.status.includes("ENABLED")) {
        let newMeshSyncState = [...MeshSyncState]
        newMeshSyncState[index] = null;
        setMeshsyncSubscription({ type : actionTypes.SET_MESHSYNC_SUBSCRIPTION, meshSyncState : newMeshSyncState })
        handleError("MeshSync could not be reached")("MeshSync is unavailable");
      } else {
        let publishEndpoint = res.controller.status.substring("ENABLED".length)
        enqueueSnackbar(`MeshSync was successfully pinged. Publishing to ${publishEndpoint} `, {
          variant : "success",
          action : (key) => (
            <IconButton key="close" aria-label="close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ),
          autohideduration : 2000,
        })
      }
    },
    error : handleError("MeshSync status could not be retrieved"), });

    // connectToNats().subscribe({
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
    //     if (res.connectToNats === "CONNECTED") {
    //       this.props.updateProgress({ showProgress : false });
    //       this.props.enqueueSnackbar(`Successfully connected to NATS`, {
    //         variant : "success",
    //         action : (key) => (
    //           <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closesnackbar(key)}>
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
        <Typography className={classes.add}> Add Cluster </Typography>
        <AddIcon fontSize="small" />
      </Button>

      <DataTable
        columns = { columns }
        data = { data }
        options = { options }
      />
      <PromptComponent ref={ ref }/>
      <PromptComponent ref = {meshSyncResetRef} />
    </>
  )
}
const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const selectedK8sContexts = state.get('selectedK8sContexts')
  const operatorState = state.get('operatorState');
  const MeshSyncState = state.get('meshSyncState');
  return { k8sconfig, selectedK8sContexts, operatorState, MeshSyncState };
}
const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch),
  setMeshsyncSubscription : bindActionCreators(setMeshsyncSubscription, dispatch)
});

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesherySettingsNew)));

