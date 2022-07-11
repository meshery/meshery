import React, {useState , useEffect, useRef} from 'react'
// import { DataGrid } from '@mui/x-data-grid'; 
import MUIDataTable from "mui-datatables";
import {
    Box, Grid, IconButton, List, ListItem, ListItemText, Menu, MenuItem, Modal, Switch, FormGroup,
    Tooltip, Paper,  TableCell, TableContainer, Table, Button, Typography, Toolbar, InputAdornment,
     TableBody, TextField, TableRow , TableHead, TableSortLabel, Checkbox, TablePagination, FormControlLabel
  } from '@mui/material';
  import AddIcon from '@mui/icons-material/Add';
  import CloudUploadIcon from '@mui/icons-material/CloudUpload';
  
function MeshConfigComponent({operatorState}) {
  
  const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 400,
    bgcolor: '#fff',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
  };

  const [open, setOpen] = useState(false);
  const [data, setData] = useState([])
  const [showMenu, setShowMenu] = useState([false])
  const [anchorEl, setAnchorEl] = useState(null);
  const [contexts, setContexts] = useState([]);
  const [k8sVersion, setK8sVersion] = useState(["N/A"]);
  const [discover, setLastDiscover] = useState(['']);
  const [_operatorState, _setOperatorState] = useState(operatorState || []);
  const ref = useRef(null);
  const meshSyncResetRef = useRef(null);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;
  
  const handleOpenModal = () => setOpen(true);
  const handleCloseModal = () => setOpen(false);

  let k8sfileElementVal = "";
  let formData = new FormData();

  const columns = [
    {
      name : "contexts",
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
                label={data[tableMeta.rowIndex].name}
                icon={<img src="/static/img/kubernetes.svg" />}
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
                // onClick={(e) => handleMenuOpen(e, tableMeta.rowIndex)}
              >
                <MoreVertIcon />
              </IconButton>
              <Menu
                id="long-menu"
                MenuListProps={{
                  'aria-labelledby' : 'long-button',
                }}
                anchorEl={anchorEl}
                // open={showMenu[tableMeta.rowIndex]}
                // onClose={() => handleMenuClose(tableMeta.rowIndex)}
              >
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  data-cy="btnResetDatabase"
                >
                  <Typography> Flush MeshSync </Typography>
                </Button>
                <MenuItem>
                  <Switch
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
        // handleConfigDelete(data[item.index].id)
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
                <TableCell>
                  <Paper >
                    <div>
                      <Grid container spacing={1} >
                        <Grid item xs={12} md={5}>
                          <List>
                            <ListItem>
                              <Tooltip title={`Server: ${contexts[rowMetaData.rowIndex].server}`}
                              >
                                <Chip
                                  label={data[rowMetaData.rowIndex].name}
                                  icon={<img src="/static/img/kubernetes.svg" />}
                                  variant="outlined"
                                  data-cy="chipContextName"
                                />
                              </Tooltip>
                            </ListItem>
                          </List>
                        </Grid>
                      </Grid>
                      <Grid container spacing={1} >
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
                              <ListItemText primary="Server" secondary={
                                contexts[rowMetaData.rowIndex].server
                              } />
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

                        <Grid item xs={12} md={4}>
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
                                  icon={<img src="/static/img/meshery-operator.svg" />}
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
                                    title={meshSyncState ? `Ping MeshSync` : "Not Available"}
                                    aria-label="meshSync"
                                  >
                                    <Chip
                                      label={"MeshSync"}
                                      icon={<img src="/static/img/meshsync.svg" />}
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
                                    title={natsState ? `Reconnect NATS` : "Not Available"}
                                    aria-label="nats"
                                  >
                                    <Chip
                                      label={"NATS"}
                                      icon={<img src="/static/img/nats-icon-color.svg"  />}
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

                      <Grid container spacing={1} >
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
              />
              <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                label="Upload kubeconfig"
                variant="outlined"
                fullWidth
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
        return;
      }
    }
  }

  return (
    <>
    <div style={{ display : 'table', tableLayout : 'fixed', width : '100%' }}>
    <MUIDataTable
      title={
        <Button
          type="submit"
          variant="contained"
          color="primary"
          size="large"
          onClick={handleOpenModal}
          data-cy="btnResetDatabase"
        >
          <AddIcon fontSize="small" />
          <Typography > Add Cluster</Typography>
        </Button>
      }
      columns={columns}
      data={data}
      options={options}
    />
  </div>
  <Modal 
        open={open}
        onClose={handleCloseModal}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={style}>
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
              />
              <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                label="Upload kubeconfig"
                variant="outlined"
                fullWidth
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
          </Box>
      </Modal>
  </>
  )
}

export default MeshConfigComponent