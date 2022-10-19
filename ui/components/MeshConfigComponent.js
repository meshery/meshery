import {
  Chip, Grid, IconButton, List, ListItem, ListItemText, Menu, MenuItem, Switch,
  Tooltip, Paper, NoSsr, TableCell, TableContainer, Table, Button, Typography,
  TextField, FormGroup, InputAdornment
} from '@material-ui/core';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import TableSortLabel from "@material-ui/core/TableSortLabel";
import CloseIcon from "@material-ui/icons/Close";
import { withSnackbar } from "notistack";
import { useState, useEffect, useRef } from 'react';
import DataTable from "mui-datatables";
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import PromptComponent from './PromptComponent';
import CloudUploadIcon from "@material-ui/icons/CloudUpload";
import AddIcon from "@material-ui/icons/AddCircleOutline";
import MeshsyncStatusQuery from './graphql/queries/MeshsyncStatusQuery';
import NatsStatusQuery from './graphql/queries/NatsStatusQuery';
import changeOperatorState from './graphql/mutations/OperatorStatusMutation';
import resetDatabase from "./graphql/queries/ResetDatabaseQuery";
import { updateProgress } from "../lib/store";
import fetchMesheryOperatorStatus from "./graphql/queries/OperatorStatusQuery";
import _ from "lodash";

const styles = (theme) => ({
  operationButton : {
    [theme.breakpoints.down(1180)] : {
      marginRight : "25px",
    },
  },
  icon : { width : theme.spacing(2.5), },
  paper : { margin : theme.spacing(2), padding : theme.spacing(2), },
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
  fileInputStyle : { display : "none", },
  button : {
    padding : theme.spacing(1),
    borderRadius : 5
  },
  grey : {
    background : "WhiteSmoke",
    padding : theme.spacing(2),
    borderRadius : "inherit",
  },
  fileLabelText : {
    cursor : "pointer",
    "& *" : { cursor : "pointer", },
  },
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
  addIcon : {
    paddingLeft : theme.spacing(.5),
    marginRight : theme.spacing(.5),
  },
  FlushBtn : {
    margin : theme.spacing(0.5),
    padding : theme.spacing(1),
    borderRadius : 5
  },
  menu : {
    display : 'flex',
    alignItems : 'center'
  },
  table : {
    marginTop : theme.spacing(1.5)
  },
  uploadCluster : {
    overflow : "hidden"
  }
});

const ENABLED = "ENABLED"
const DISABLED = "DISABLED"

function MesherySettingsNew({ classes, enqueueSnackbar, closeSnackbar, updateProgress,
  operatorState, k8sconfig }) {
  const [data, setData] = useState([])
  const [showMenu, setShowMenu] = useState([false])
  const [anchorEl, setAnchorEl] = useState(null);
  const [NATSState, setNATSState] = useState(["UNKNOWN"]);
  const [NATSVersion, setNATSVersion] = useState(["N/A"]);
  const [contexts, setContexts] = useState([]);
  const [k8sVersion, setK8sVersion] = useState(["N/A"]);
  const [discover, setLastDiscover] = useState(['']);
  const [_operatorState, _setOperatorState] = useState(operatorState || []);
  const ref = useRef(null);
  const meshSyncResetRef = useRef(null);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;

  const dateOptions = { weekday : 'long', year : 'numeric', month : 'long', day : 'numeric' };

  let k8sfileElementVal = "";
  let formData = new FormData();

  const stateUpdater = (state, updateFunc, updateValue, index) => {
    let newState = [...state];
    newState[index] = updateValue;
    updateFunc(newState);
  }

  const setTableData = () => {
    let tableInfo = [];
    handleContexts(k8sconfig);
    k8sconfig.forEach(ctx => {
      let data = {
        context : ctx.name,
        location : ctx.server,
        deployment_type : ctx.inClusterConfig ? "In Cluster" : "Out of Cluster",
        last_discovery : setDateTime(new Date()),
        id : ctx.id
      };
      tableInfo.push(data);
    })
    setData(tableInfo);
  }

  useEffect(() => {
    setTableData();
  },[k8sconfig])
  useEffect(() => {
    setTableData();
    k8sconfig.forEach(ctx => {
      const tempSubscription = fetchMesheryOperatorStatus({ k8scontextID : ctx.id }).
        subscribe({
          next : (res) => {
            if (!_operatorState?.find(opSt => opSt.contextID === ctx.id)) {
              const x = updateCtxInfo(ctx.id, res)
              _setOperatorState(x)
            }
            tempSubscription.unsubscribe();
          },
          error : (err) => console.log("error at operator scan: " + err),
        })
    })
    getKubernetesVersion(); // change to per context and also ping return server version. no need to check workload
  }, [])

  useEffect(() => {
    if (operatorState) {
      _setOperatorState(operatorState);
    }
  }, [operatorState])

  const handleFlushMeshSync = (index) => {
    return async () => {
      handleMenuClose(index);
      let response = await meshSyncResetRef.current.show({
        title : `Flush MeshSync data for ${data[index].context} ?`,
        subtitle : `Are you sure to Flush MeshSync data for “${data[index].context}”? Fresh MeshSync data will be repopulated for this context, if MeshSync is actively running on this cluster.`,
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

  const setDateTime = (dt) => {
    return dt.toLocaleDateString("en-US", dateOptions)
      + " " + dt.toLocaleTimeString("en-US");
  }

  const handleContexts = (contexts) => {
    let ctxs = []
    contexts.forEach((ctx) => {
      let tempCtx = { ...ctx }
      tempCtx.created_at = setDateTime(new Date(ctx.created_at));
      tempCtx.updated_at = setDateTime(new Date(ctx.updated_at));
      ctxs.push(tempCtx);
    })
    setContexts(ctxs);
  }

  const handleMenuClose = (index) => {
    let menu = [...showMenu];
    menu[index] = false;
    setShowMenu(menu)
  }

  const handleError = (msg) => (error) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar(`${msg}: ${error}`, {
      variant : "error", preventDuplicate : true,
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000,
    });
  };

  const handleMenuOpen = (e, index) => {
    setAnchorEl(e.currentTarget)
    let menu = [...showMenu];
    menu[index] = true;
    setShowMenu(menu);
  }

  const handleConfigSnackbars = ctxs => {
    updateProgress({ showProgress : false });

    for (let ctx of ctxs.inserted_contexts) {
      const msg = `Cluster ${ctx.name} at ${ctx.server} connected`
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

    for (let ctx of ctxs.updated_contexts) {
      const msg = `Cluster ${ctx.name} at ${ctx.server} already exists`
      enqueueSnackbar(msg, {
        variant : "info",
        action : (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration : 7000,
      });
    }

    for (let ctx of ctxs.errored_contexts) {
      const msg = `Failed to add cluster ${ctx.name} at ${ctx.server}`
      enqueueSnackbar(msg, {
        variant : "error",
        action : (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration : 7000,
      });
    }
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
      },
      handleError("Failed to get Kubernetes version")
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

  const updateCtxInfo = (ctxId, newInfo) => {
    if (newInfo.operator.error) {
      handleError("There is problem With operator")(newInfo.operator.error);
      return;
    }

    const state = _operatorStateRef.current;
    const op = state?.find(ctx => ctx.contextID === ctxId);
    if (!op) {
      return [...state, { contextID : ctxId, operatorStatus : newInfo.operator }];
    }

    let ctx = { ...op };
    const removeCtx = state?.filter(ctx => ctx.contextID !== ctxId);
    ctx.operatorStatus = newInfo.operator;
    return removeCtx ? [...removeCtx, ctx]: [ctx];
  }

  const handleOperatorSwitch = (index, checked) => {
    const contextId = contexts[index].id;
    const variables = {
      status : `${checked ? ENABLED : DISABLED}`,
      contextID : contextId
    };

    updateProgress({ showProgress : true });

    changeOperatorState((response, errors) => {
      updateProgress({ showProgress : false });

      if (errors !== undefined) {
        handleError(`Unable to ${!checked ? "Uni" : "I"}nstall operator`);
      }
      enqueueSnackbar("Operator " + response.operatorStatus.toLowerCase(), {
        variant : "success",
        autoHideDuration : 2000,
        action : (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        ),
      });

      const tempSubscription = fetchMesheryOperatorStatus({ k8scontextID : contextId }).subscribe({
        next : (res) => {
          _setOperatorState(updateCtxInfo(contextId, res))
          tempSubscription.unsubscribe();
        },
        error : (err) => console.log("error at operator scan: " + err),
      })

    }, variables);
  };

  const handleConfigDelete = (id, index) => {
    updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/kubernetes/contexts/" + id,
      {
        credentials : "same-origin",
        method : "DELETE"
      },
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
      textField.value = name;
      formData = formdata;

    }
  }

  function getOperatorStatus(ctxId) {
    const operator = _operatorState?.find(op => op.contextID === ctxId);
    if (!operator) {
      return {}
    }
    const operatorStatus = operator.operatorStatus;
    return {
      operatorState : operatorStatus.status === ENABLED,
      operatorVersion : operatorStatus.version,
    }
  }

  const getContextStatus = (ctxId) => {
    const operator = _operatorStateRef.current?.find(op => op.contextID === ctxId);
    if (!operator) {
      return {}
    }
    const operatorStatus = operator.operatorStatus;

    function getMeshSyncStats() {
      if (!operatorStatus) return {};
      const { controllers } = operatorStatus;
      // meshsync is at 1st idx
      if (controllers?.[1]) {
        const { status, version } = controllers[1];
        return {
          meshSyncState : status,
          meshSyncVersion : version
        }
      }
    }

    function getBrokerStats() {
      if (!operatorStatus) return {};
      const { controllers } = operatorStatus;
      // broker is at 0th idx
      if (controllers?.[0]) {
        const { status, version } = controllers[0];
        if (status != "") {
          return {
            natsState : status,
            natsVersion : version
          }
        }
      }
    }

    const defaultState = {
      operatorState : false,
      operatorVersion : null,
      meshSyncState : DISABLED,
      meshSyncVersion : "Not Available",
      natsState : "Not Active",
      natsVersion : "Not Available"
    }

    const actualOperatorState = {
      ...getOperatorStatus(ctxId),
      ...getMeshSyncStats(),
      ...getBrokerStats()
    }

    return _.merge(defaultState, actualOperatorState);
  }

  const uploadK8SConfig = async () => {
    return await promisifiedDataFetch(
      "/api/system/kubernetes",
      {
        method : "POST",
        body : formData,
      }
    )
  }

  const columns = [
    {
      name : "context",
      label : "Contexts",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={"center"} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (_, tableMeta,) => {
          return (
            <Tooltip title={`Server: ${tableMeta.rowData[2]}`}>
              <Chip
                label={data[tableMeta.rowIndex].context}
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
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={"center"} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      }
    },
    {
      name : "location",
      label : "Location",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={"center"} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      }
    },
    {
      name : "last_discovery",
      label : "Last Discovery",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={"center"} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value, tableMeta) => <p>{discover[tableMeta.rowIndex] || value}</p>
      }
    },
    {
      name : "Actions",
      options : {
        filter : true,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
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
                    checked={getOperatorStatus(contexts[tableMeta.rowIndex].id)?.operatorState}
                    onClick={(e) => handleOperatorSwitch(tableMeta.rowIndex, e.target.checked)}
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
    elevation : 0,
    label : "",
    responsive : "standard",
    fixedHeader : true,
    textLabels : {
      selectedRows : {
        text : "context(s) selected"
      }
    },
    onRowsDelete : (td) => {
      td.data.forEach((item) => {
        handleConfigDelete(data[item.index].id)
      })
    },


    renderExpandableRow : (rowData, rowMetaData) => {
      const contextId = contexts[rowMetaData.rowIndex].id;
      const { meshSyncState, meshSyncVersion, natsState, natsVersion, operatorState, operatorVersion } = getContextStatus(contextId);
      return (
        <NoSsr>
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
                              <Tooltip title={`Server: ${contexts[rowMetaData.rowIndex].server}`}
                              >
                                <Chip
                                  label={data[rowMetaData.rowIndex].context}
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
                              <ListItemText primary="Name" secondary={contexts[rowMetaData.rowIndex].name} />
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
                              } />
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
                              } />
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
                                title={operatorState
                                  ? `Version: ${operatorVersion}`
                                  : "Not Available"}
                                aria-label="meshSync"
                              >
                                <Chip
                                  // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
                                  label={"Operator"}
                                  style={!operatorState ? { opacity : 0.5 }: {}}
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

                        {(meshSyncState || natsState) &&
                          <>
                            <Grid item xs={12} md={4}>
                              <List>
                                <ListItem>
                                  <Tooltip
                                    title={meshSyncState !== DISABLED ? `Ping MeshSync` : "Not Available"}
                                    aria-label="meshSync"
                                  >
                                    <Chip
                                      label={"MeshSync"}
                                      style={ meshSyncState === DISABLED ? { opacity : 0.5 }: {}}
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
                                    title={natsState === "Not Active" ? "Not Available": `Reconnect NATS`}
                                    aria-label="nats"
                                  >
                                    <Chip
                                      label={"NATS"}
                                      onClick={() => handleNATSClick(rowMetaData.rowIndex)}
                                      style={natsState === "Not Active" ? { opacity : 0.5 }: {}}
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
                              <ListItemText primary="Operator State" secondary={operatorState ? "Active" : "Disabled"} />
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="Operator Version" secondary={operatorVersion} />
                            </ListItem>
                          </List>
                        </Grid>
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <ListItemText primary="MeshSync State" secondary={meshSyncState || "Disabled"} />
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="MeshSync Version" secondary={meshSyncVersion} />
                            </ListItem>
                          </List>
                        </Grid>
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <ListItemText primary="NATS State" secondary={natsState || "Not Connected"} />
                            </ListItem>
                            <ListItem>
                              <ListItemText primary="NATS Version" secondary={natsVersion} />
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
      title : "Add Kubernetes Cluster(s)",
      subtitle :
        <>
          <div className={classes.uploadCluster}>
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
                InputProps={{
                  readOnly : true,
                  endAdornment : (
                    <InputAdornment position="end">
                      <CloudUploadIcon />
                    </InputAdornment>
                  ),
                }}
              />
            </FormGroup>
          </div>
        </>,
      options : ["IMPORT","CANCEL"]
    })

    if (response === "IMPORT") {
      if (formData.get("k8sfile") === null) {
        handleError("No file selected")("Please select a valid kube config")
        return;
      }

      const inputFile = ( formData.get( "k8sfile" ).name );
      const invalidExtensions = /^.*\.(jpg|gif|jpeg|pdf|png|svg)$/i;

      if (invalidExtensions.test(inputFile)  ) {
        handleError("Invalid file selected")("Please select a valid kube config")
        return;
      }

      uploadK8SConfig().then((obj) => {
        handleConfigSnackbars(obj);
      }).
        catch(err => {
          handleError("failed to upload kubernetes config")(err)
        })
      formData.delete("k8sfile");
    }
  }

  const handleOperatorClick = (index) => {
    updateProgress({ showProgress : true });
    const ctxId = contexts[index].id
    const tempSubscription = fetchMesheryOperatorStatus({ k8scontextID : ctxId })
      .subscribe({
        next : (res) => {
          _setOperatorState(updateCtxInfo(ctxId, res))

          updateProgress({ showProgress : false });
          if (!res.operator.error && res.operator.status === ENABLED ) {
            enqueueSnackbar("Operator was successfully pinged!", {
              variant : "success",
              autoHideDuration : 2000,
              action : (key) => (
                <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                  <CloseIcon />
                </IconButton>
              ),
            });
          } else {
            handleError("Operator could not be reached")("Operator is disabled");
          }
          tempSubscription.unsubscribe();
        },
        error : handleError("Operator could not be pinged"),
      });
  };

  const handleNATSClick = (index) => {
    updateProgress({ showProgress : true });
    NatsStatusQuery({ k8scontextID : contexts[index].id }).subscribe({
      next : (res) => {
        updateProgress({ showProgress : false });
        if (res.controller.name === "broker" && res.controller.status.includes("CONNECTED")) {
          let runningEndpoint = res.controller.status.substring("CONNECTED".length).trim();
          enqueueSnackbar(`Broker was successfully pinged. ${runningEndpoint != "" ? `Running at ${runningEndpoint}` : ""}`, {
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

        stateUpdater(NATSState, setNATSState, res.controller.status.length !== 0 ? res.controller.status : "UNKNOWN", index)
        stateUpdater(NATSVersion, setNATSVersion, res.controller.version, index);
      },
      error : handleError("NATS status could not be retrieved"),
    });
  };

  const setMeshSyncStatusForGivenContext = (ctxId, meshsyncStatus) => {
    if (!ctxId) {
      return;
    }

    let operator = _operatorStateRef.current?.find(context => context.contextID === ctxId);
    if (operator) {
      let controllers = [operator.operatorStatus.controllers[0], meshsyncStatus]
      const newOperatorStatus = Object.assign({}, {
        contextID : ctxId,
        operatorStatus : {
          ...operator.operatorStatus,
          controllers,
        }
      })
      let remainingOperatorArray = _operatorStateRef.current?.filter(context => context.contextID !== ctxId)
      let finalOperatorArray = [...remainingOperatorArray, newOperatorStatus]
      _setOperatorState(finalOperatorArray)
    }
  }

  const handleMeshSyncClick = (index) => {
    updateProgress({ showProgress : true });
    const ctxId = contexts[index].id;
    MeshsyncStatusQuery(({ k8scontextID : ctxId })).subscribe({
      next : (res) => {
        updateProgress({ showProgress : false });
        if (res.controller.name === "meshsync") {
          setMeshSyncStatusForGivenContext(ctxId, res.controller)
        }
        if (res.controller.status === DISABLED) {
          handleError("MeshSync could not be reached")("MeshSync is unavailable");
        } else {
          let publishEndpoint = res.controller.status.substring("ENABLED".length).trim()
          enqueueSnackbar(`MeshSync was successfully pinged. ${publishEndpoint != "" ? `Publishing to ${publishEndpoint}` : ""}`, {
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
      error : handleError("MeshSync status could not be retrieved"),
    });

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
    <div style={{ display : 'table', tableLayout : 'fixed', width : '100%' }}>
      <DataTable
        title={
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={handleClick}
            className={classes.button}
            data-cy="btnResetDatabase"
          >
            <AddIcon fontSize="small" />
            <Typography className={classes.addIcon}> Add Cluster</Typography>
          </Button>
        }
        columns={columns}
        data={data}
        options={options}
        className={classes.table}
      />
      <PromptComponent ref={ref} />
      <PromptComponent ref={meshSyncResetRef} />
    </div>
  )
}
const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const selectedK8sContexts = state.get('selectedK8sContexts')
  const operatorState = state.get('operatorState');
  // const MeshSyncState = state.get('meshSyncState'); // disfunctional at this point of time
  return { k8sconfig, selectedK8sContexts, operatorState, /*MeshSyncState*/ };
}
const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
  // setMeshsyncSubscription : bindActionCreators(setMeshsyncSubscription, dispatch)
});

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesherySettingsNew)));

