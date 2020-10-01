import NoSsr from '@material-ui/core/NoSsr';
import React from 'react';
import { Controlled as CodeMirror } from 'react-codemirror2';
import {
  withStyles, Grid, TextField, IconButton, Dialog, DialogTitle, DialogContent, DialogActions, Divider, Card, CardHeader, CardActions, Menu, MenuItem, Chip, TableCell, TableRow, TableBody, TableHead, Table, Tooltip
} from '@material-ui/core';
import { blue } from '@material-ui/core/colors';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import CloseIcon from '@material-ui/icons/Close';
import { withSnackbar } from 'notistack';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import PlayIcon from '@material-ui/icons/PlayArrow';
// import { updateSMIResults } from '../lib/store';
import { updateProgress, } from '../lib/store';
import dataFetch from '../lib/data-fetch';
import MUIDataTable from "mui-datatables";
import Moment from 'react-moment';
import MesheryResultDialog from './MesheryResultDialog';


const styles = (theme) => ({
  root: {
    padding: theme.spacing(10),
    width: '100%',
  },
  chipGrid: {
    padding: theme.spacing(10),
    width: '100%',
    paddingBottom: '0',
  },
  buttons: {
    width: '100%',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  alreadyConfigured: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  chip: {
    height: '40px',
    marginRight: theme.spacing(5),
    marginTop: theme.spacing(4),
    marginBottom: theme.spacing(-5),
    fontSize: '15px',
  },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': {
        backgroundColor: blue[500],
      },
    },
  },
  colorBar: {},
  colorChecked: {},
  uploadButton: {
    margin: theme.spacing(1),
    marginTop: theme.spacing(3),
  },
  fileLabel: {
    width: '100%',
  },
  editorContainer: {
    width: '100%',
  },
  deleteLabel: {
    paddingRight: theme.spacing(2),
  },
  alignRight: {
    textAlign: 'right',
  },
  alignLeft: {
    textAlign: 'left',
    marginLeft: theme.spacing(1),
  },
  padLeft: {
    paddingLeft: theme.spacing(0.25),
  },
  padRight: {
    paddingRight: theme.spacing(0.25),
  },
  deleteRight: {
    float: 'right',
  },
  expTitleIcon: {
    width: theme.spacing(3),
    display: 'inline',
    verticalAlign: 'middle',
  },
  expIstioTitleIcon: {
    width: theme.spacing(2),
    display: 'inline',
    verticalAlign: 'middle',
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(0.5),
  },
  expTitle: {
    display: 'inline',
    verticalAlign: 'middle',
  },
  icon: {
    width: theme.spacing(2.5),
  },
  tableHeader: {
    fontWeight: 'bolder',
    fontSize: 18,
  },
  secondaryTable: {
    borderRadius: 10,
    backgroundColor: "#f7f7f7",
  },
});


class MesheryAdapterPlayComponent extends React.Component {

  
  constructor(props) {
    super(props);

    this.cmEditorAdd = null;
    this.cmEditorDel = null;

    const { adapter } = props;

    const menuState = {};

    this.addIconEles = {};
    this.delIconEles = {};
    // initializing menuState;
    if (adapter && adapter.ops) {
      // NOTE: this will have to updated to match the categories
      [0, 1, 2, 3, 4].forEach((i) => {
        menuState[i] = {
          add: false,
          delete: false,
        };
      });
    }

    this.state = {
      selectedOp: '',
      cmEditorValAdd: '',
      cmEditorValAddError: false,

      cmEditorValDel: '',
      cmEditorValDelError: false,

      selectionError: false,

      namespace: 'default',
      namespaceError: false,

      customDialogAdd: false,
      customDialogDel: false,
      customDialogSMI: false,

      open : false,

      menuState, // category: {add: 1, delete: 0}

      smi_result: [],
      selectedRowData: null,
      page: 0,
      search: '',
      sortOrder: '',
      pageSize: 10,
    };
  }

  handleChange = (name, isDelete = false) => {
    const self = this;
    return (event) => {
      if (name === 'namespace' && event.target.value !== '') {
        this.setState({ namespaceError: false });
      }

      if (name === 'selectedOp' && event.target.value !== '') {
        if (event.target.value === 'custom') {
          if (isDelete) {
            if (self.state.cmEditorValDel !== '' && self.cmEditorDel.state.lint.marked.length === 0) {
              self.setState({ selectionError: false, cmEditorValDelError: false });
            }
          } else if (self.state.cmEditorValAdd !== '' && self.cmEditorAdd.state.lint.marked.length === 0) {
            self.setState({ selectionError: false, cmEditorValAddError: false });
          }
        } else {
          self.setState({ selectionError: false });
        }
      }

      self.setState({ [name]: event.target.value });
    };
  }

  handleModalClose(isDelete) {
    const self = this;
    return () => {
      const item = isDelete ? 'customDialogDel' : 'customDialogAdd';
      self.setState({ [item]: false });
    };
  }

  handleSMIClose() {
    const self = this;
    return () => {
      self.setState({['customDialogSMI']: false });
    }
  }

  resetSelectedRowData() {
    const self = this;
    return () => {
      self.setState({ selectedRowData: null });
    };
  }

  handleModalOpen(isDelete) {
    const self = this;
    return () => {
      const item = isDelete ? 'customDialogDel' : 'customDialogAdd';
      self.setState({ [item]: true });
    };
  }

  handleSubmit = (cat, selectedOp, deleteOp = false) => {
    const self = this;
    return () => {
      const {
        namespace, cmEditorValAdd, cmEditorValDel,
      } = self.state;
      const { adapter } = self.props;
      const filteredOp = adapter.ops.filter(({ key }) => key === selectedOp);
      if (selectedOp === '' || typeof filteredOp === 'undefined' || filteredOp.length === 0) {
        self.setState({ selectionError: true });
        return;
      }
      if (deleteOp) {
        if (selectedOp === 'custom' && (cmEditorValDel === '' || self.cmEditorDel.state.lint.marked.length > 0)) {
          self.setState({ cmEditorValDelError: true, selectionError: true });
          return;
        }
      } else if (selectedOp === 'custom' && (cmEditorValAdd === '' || self.cmEditorAdd.state.lint.marked.length > 0)) {
        self.setState({ cmEditorValAddError: true, selectionError: true });
        return;
      }
      if (namespace === '') {
        self.setState({ namespaceError: true });
        return;
      }
      self.submitOp(cat, selectedOp, deleteOp);
    };
  }

  submitOp = (cat, selectedOp, deleteOp = false) => {
    const {
      namespace, cmEditorValAdd, cmEditorValDel, menuState,
    } = this.state;
    const { adapter } = this.props;
    // const fileInput = document.querySelector('#k8sfile') ;

    const data = {
      adapter: adapter.adapter_location,
      query: selectedOp,
      namespace,
      customBody: deleteOp ? cmEditorValDel : cmEditorValAdd,
      deleteOp: deleteOp ? 'on' : '',
    };

    const params = Object.keys(data).map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`).join('&');
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch('/api/mesh/ops', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
      },
      body: params,
    }, (result) => {
      self.props.updateProgress({ showProgress: false });
      menuState[cat][deleteOp ? 'delete' : 'add'] = false;
      const dlg = deleteOp ? 'customDialogDel' : 'customDialogAdd';
      self.setState({ menuState, [dlg]: false });

      if (typeof result !== 'undefined') {
        self.props.enqueueSnackbar('Operation executing...', {
          variant: 'info',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key)}
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError(cat, deleteOp));
  }

  handleAdapterClick = (adapterLoc) => () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(`/api/mesh/adapter/ping?adapter=${encodeURIComponent(adapterLoc)}`, {
      credentials: 'same-origin',
      credentials: 'include',
    }, (result) => {
      this.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        this.props.enqueueSnackbar('Adapter successfully pinged!', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key)}
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError('Could not ping adapter.'));
  }

  fetchSMIResults= (adapterName, page, pageSize, search, sortOrder) => {
    const self = this;
    let query = '';
    if (typeof search === 'undefined' || search === null) {
      search = '';
    }
    if (typeof sortOrder === 'undefined' || sortOrder === null) {
      sortOrder = '';
    }
    query = `?page=${page}&pageSize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(sortOrder)}`;

    dataFetch(`/api/smi/results${query}`, {
      credentials: 'same-origin',
      method: 'GET',
      credentials: 'include',
    }, (result) => {
      if (typeof result !== 'undefined' && result.results) {
        const results  = result.results.filter(val => val.mesh_name.toLowerCase()==adapterName.toLowerCase())
        self.setState({smi_result: {...result, results:results, total_count:results.length}});
      }
    }, self.handleError('Could not fetch SMI results.'));
  }

  handleSMIClick = (adapterName) => () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    const {page, pageSize, search, sortOrder} = self.state;
    self.fetchSMIResults(adapterName, page, pageSize, search, sortOrder);
    this.props.updateProgress({ showProgress: false });
    self.setState({ ['customDialogSMI']: true })
    

  }

  handleError = (cat, deleteOp) => {
    const self = this;
    return (error) => {
      if(cat && deleteOp) {
        const { menuState } = self.state;
        menuState[cat][deleteOp ? 'delete' : 'add'] = false;
        const dlg = deleteOp ? 'customDialogDel' : 'customDialogAdd';
        self.setState({ menuState, [dlg]: false });
      }
      self.props.updateProgress({ showProgress: false });
      self.props.enqueueSnackbar(`Operation submission failed: ${error}`, {
        variant: 'error',
        action: (key) => (
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={() => self.props.closeSnackbar(key)}
          >
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration: 8000,
      });
    };
  }

  handleExpandClick() {
    // setExpanded(!expanded);
  }

  generateMenu(cat, isDelete, selectedAdapterOps) {
    const { menuState } = this.state;
    const ele = !isDelete ? this.addIconEles[cat] : this.delIconEles[cat];
    return (
      <Menu
        id="long-menu"
        anchorEl={ele}
        keepMounted
        open={menuState[cat][isDelete ? 'delete' : 'add']}
        onClose={this.addDelHandleClick(cat, isDelete)}
      >
        {selectedAdapterOps.map(({ key, value }) => (
          <MenuItem key={`${key}_${new Date().getTime()}`} onClick={this.handleSubmit(cat, key, isDelete)}>
            {value}
          </MenuItem>
        ))}
      </Menu>
    );
  }

  handleOpen = () => {
    setOpen(true);
  }

  handleClose = () => {
    setOpen(false);
  }


  generateSMIResult() {
    const self = this;
    
    const {
      customDialogSMI, smi_result, pageSize
    } = self.state;

    const {
      user, classes,
    } = self.props;
    
    
    const smi_columns = [
      {
        name: 'ID',
        label: 'ID',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({index, ...column}) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
              
            )
          },
          customBodyRender: (value) => (
            <Tooltip title={value} placement="top">
              <div>{value.slice(0,5)+ "..."}</div>
            </Tooltip>
          )
        },
      },
      {
        name: 'Date',
        label: 'Date',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({index, ...column}) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
                
            )
          },
          customBodyRender: (value) => (
            <Moment format="LLLL">{value}</Moment>
          ),
        },
      },
      {
        name: 'Service Mesh',
        label: 'Service Mesh',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({index, ...column}) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
                
            )
          },
        },
      },
      {
        name: 'Service Mesh Version',
        label: 'Service Mesh Version',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({index, ...column}) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
                
            )
          },
        },
      },
      {
        name: '% Passed',
        label: '% Passed',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({index, ...column}) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
                
            )
          },
        },
      },
      {
        name: 'status',
        label: 'Status',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({index, ...column}) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            )
          },
        },
      }
    ]

    const smi_options = {
      sort: !(user && user.user_id === 'meshery'),
      search: !(user && user.user_id === 'meshery'),
      filterType: 'textField',
      expandableRows: true,
      selectableRows: false,
      rowsPerPage: pageSize,
      rowsPerPageOptions: [10, 20, 25],
      fixedHeader: true,
      print: false,
      download: false,
      renderExpandableRow: (rowData, rowMeta) => {
        const column = ["Specification","Assertions", "Time","Version", "Capability", "Result", "Reason"]
        const data = smi_result.results[rowMeta.dataIndex].more_details.map((val) => {
          return [val.smi_specification,val.assertions,val.time,"alpha1/v1",val.capability,val.status,val.reason] 
        })
        const colSpan = rowData.length + 1
        return (
          <TableRow>
            <TableCell colSpan={colSpan}>
              <div className={classes.secondaryTable}>
                <Table  aria-label="a dense table">
                  <TableHead>
                    <TableRow>
                      {column.map((val) => (<TableCell colSpan={colSpan}>{val}</TableCell>))}
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {data.map((row) => (
                      <TableRow >
                        {row.map(val => {
                          if(val && val.match(/[0-9]+m[0-9]+.+[0-9]+s/i)!=null) {
                            const time = val.split(/m|s/)
                            return <TableCell colSpan={colSpan}>{time[0]+"m " + parseFloat(time[1]).toFixed(1) + "s"}</TableCell>
                          } else {
                            return <TableCell colSpan={colSpan}>{val}</TableCell>
                          }
                        }
                        )}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </TableCell>
          </TableRow>
        );
      },
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText? tableState.announceText.split(' : '):[];
        let order='';
        if(tableState.activeColumn){
          order = `${columns[tableState.activeColumn].name} desc`;
        }

        switch (action) {
          case 'changePage':
            self.fetchSMIResults(self.props.adapter.name,tableState.page, self.state.pageSize, self.state.search, self.state.sortOrder);
            break;
          case 'changeRowsPerPage':
            self.fetchSMIResults(self.props.adapter.name,self.state.page, tableState.rowsPerPage, self.state.search, self.state.sortOrder);
            break;
          case 'search':
            if (self.searchTimeout) {
              clearTimeout(self.searchTimeout);
            }
            self.searchTimeout = setTimeout(() => {
              if (self.state.search !== tableState.searchText) {
                self.fetchSMIResults(self.props.adapter.name,self.state.page, self.state.pageSize, tableState.searchText !== null ? tableState.searchText : '', self.state.sortOrder);
              }
            }, 500);
            break;
          case 'sort':
            if (sortInfo.length == 2) {
              if (sortInfo[1] === 'ascending') {
                order = `${columns[tableState.activeColumn].name} asc`;
              } else {
                order = `${columns[tableState.activeColumn].name} desc`;
              }
            }
            if (order !== this.state.sortOrder) {
              self.fetchSMIResults(self.props.adapter.name,self.state.page, self.state.pageSize, self.state.search, order);
            }
            break;
        }
      },
    }

    var data = [];
    if(smi_result&&smi_result.results) {
      data = smi_result.results.map((val) => {
        return [val.id,val.date,val.mesh_name,val.mesh_version,val.passing_percentage,val.status]
      }) 
    }


    return (
      <Dialog
        onClose={this.handleSMIClose()}
        aria-labelledby="adapter-dialog-title"
        open={customDialogSMI}
        fullWidth
        maxWidth="md"
      >
        <MUIDataTable
          title={<div className={classes.tableHeader}>Service Mesh Interface Conformance Results</div>}
          data={data}
          columns={smi_columns}
          options={smi_options}
        />
      </Dialog>
    );
  }

  generateYAMLEditor(cat, isDelete) {
    const { adapter } = this.props;
    const {
      customDialogAdd, customDialogDel, namespace, namespaceError, cmEditorValAdd, cmEditorValDel,
    } = this.state;
    const self = this;
    return (
      <Dialog
        onClose={this.handleModalClose(isDelete)}
        aria-labelledby="adapter-dialog-title"
        open={isDelete ? customDialogDel : customDialogAdd}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle id="adapter-dialog-title" onClose={this.handleModalClose(isDelete)}>
          {adapter.name}
          {' '}
          Adapter - Custom YAML
          {isDelete ? '(delete)' : ''}
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent>
          <Grid container spacing={5}>
            <Grid item xs={12}>
              <TextField
                required
                id="namespace"
                name="namespace"
                label="Namespace"
                fullWidth
                value={namespace}
                error={namespaceError}
                margin="normal"
                variant="outlined"
                onChange={this.handleChange('namespace')}
              />
            </Grid>
            <Grid item xs={12}>
              <CodeMirror
                editorDidMount={(editor) => {
                  if (isDelete) {
                    self.cmEditorDel = editor;
                  } else {
                    self.cmEditorAdd = editor;
                  }
                }}
                value={isDelete ? cmEditorValDel : cmEditorValAdd}
                options={{
                  theme: 'material',
                  lineNumbers: true,
                  lineWrapping: true,
                  gutters: ['CodeMirror-lint-markers'],
                  lint: true,
                  mode: 'text/x-yaml',
                }}
                onBeforeChange={(editor, data, value) => {
                  if (isDelete) {
                    self.setState({ cmEditorValDel: value });
                  } else {
                    self.setState({ cmEditorValAdd: value });
                  }
                  if (isDelete) {
                    if (value !== '' && self.cmEditorDel.state.lint.marked.length === 0) {
                      self.setState({ selectionError: false, cmEditorValDelError: false });
                    }
                  } else if (value !== '' && self.cmEditorAdd.state.lint.marked.length === 0) {
                    self.setState({ selectionError: false, cmEditorValAddError: false });
                  }
                }}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <IconButton aria-label="Apply" color="primary" onClick={this.handleSubmit(cat, 'custom', isDelete)}>
            {/* <FontAwesomeIcon icon={faArrowRight} transform="shrink-4" fixedWidth /> */}
            {!isDelete && <PlayIcon />}
            {isDelete && <DeleteIcon />}
          </IconButton>

        </DialogActions>
      </Dialog>
    );
  }


  addDelHandleClick = (cat, isDelete) => {
    const self = this;
    return () => {
      const { menuState, customDialogAdd, customDialogDel } = self.state;
      menuState[cat][isDelete ? 'delete' : 'add'] = !menuState[cat][isDelete ? 'delete' : 'add'];

      const dlg = isDelete ? 'customDialogDel' : 'customDialogAdd';
      let dlgv = isDelete ? customDialogDel : customDialogAdd;
      if (cat === 4) {
        dlgv = !dlgv;
      }
      self.setState({ menuState, [dlg]: dlgv });
    };
  }

  generateCardForCategory(cat) {
    if (typeof cat === 'undefined') {
      cat = 0;
    }
    const { classes, adapter } = this.props;
    // const expanded = false;

    const selectedAdapterOps = adapter && adapter.ops ? adapter.ops.filter(({ category }) => typeof category === 'undefined' && cat === 0 || category === cat) : [];
    let content;
    let description;
    switch (cat) {
      case 0:
        content = 'Manage Service Mesh Lifecycle';
        description = 'Deploy a service mesh or SMI adapter on your cluster.';
        break;

      case 1:
        content = 'Manage Sample Application Lifecycle';
        description = 'Deploy sample applications on/off the service mesh.';
        break;

      case 2:
        content = 'Apply Service Mesh Configuration';
        description = 'Configure your service mesh using some pre-defined options.';
        break;

      case 3:
        content = 'Validate Service Mesh Configuration';
        description = 'Validate your service mesh configuration against best practices.';
        break;

      case 4:
        content = 'Apply Custom Configuration';
        description = 'Customize the configuration of your service mesh.';
        break;
    }
    return (
      <Card className={classes.card}>
        <CardHeader
          title={content}
          subheader={description}
        />
        <CardActions disableSpacing>
          <IconButton aria-label="install" ref={(ch) => this.addIconEles[cat] = ch} onClick={this.addDelHandleClick(cat, false)}>
            {cat !== 4 ? <AddIcon /> : <PlayIcon />}
          </IconButton>
          {cat !== 4 && this.generateMenu(cat, false, selectedAdapterOps)}
          {cat === 4 && this.generateYAMLEditor(cat, false)}
          {cat !== 3 && (
            <div className={classes.fileLabel}>
              <IconButton aria-label="delete" ref={(ch) => this.delIconEles[cat] = ch} className={classes.deleteRight} onClick={this.addDelHandleClick(cat, true)}>
                <DeleteIcon />
              </IconButton>
              {cat !== 4 && this.generateMenu(cat, true, selectedAdapterOps)}
              {cat === 4 && this.generateYAMLEditor(cat, true)}
            </div>
          )}
        </CardActions>
      </Card>
    );
  }

  render() {
    const {
      classes, adapter,
    } = this.props;
    const {
      namespace,
      namespaceError,
      selectedRowData
    } = this.state;

    let adapterName = (adapter.name).split(" ").join("").toLowerCase();
    let imageSrc = "/static/img/" + adapterName + ".svg";
    let adapterChip = (
      <Chip
        label={adapter.adapter_location}
        onClick={this.handleAdapterClick(adapter.adapter_location)}
        icon={<img src={imageSrc} className={classes.icon} />}
        className={classes.chip}
        variant="outlined"
      />
    );

    let imageSMISrc = "/static/img/smi.png";
    let smiChip = (
      <React.Fragment>
        <Chip
          label="View SMI Conformance results"
          onClick={this.handleSMIClick(adapter.name)}
          icon={<img src={imageSMISrc} className={classes.icon} />}
          className={classes.chip}
          variant="outlined"
        />
        {this.generateSMIResult()}
      </React.Fragment>
    );

    const filteredOps = [];
    if (adapter && adapter.ops && adapter.ops.length > 0) {
      adapter.ops.forEach(({ category }) => {
        if (typeof category === 'undefined') {
          category = 0;
        }
        if (filteredOps.indexOf(category) === -1) {
          filteredOps.push(category);
        }
      });
      filteredOps.sort();
    }

    return (
      <NoSsr>
        {selectedRowData && selectedRowData !== null && Object.keys(selectedRowData).length > 0
        && (
          <MesheryResultDialog
            rowData={selectedRowData}
            close={self.resetSelectedRowData()}
          />
        )}
        <React.Fragment>
          <div className={classes.chipGrid}>
            <Grid container spacing={3}>
              <Grid item xs={3}>
                {adapterChip}
              </Grid>
              <Grid item xs={3}>
              </Grid>
              <Grid item xs={3}>
              </Grid>
              <Grid item xs={3}>
                {smiChip}
              </Grid>
            </Grid>
          </div>
          <div className={classes.root}>
            <Grid container spacing={5}>
              <Grid item xs={12}>
                <TextField
                  required
                  id="namespace"
                  name="namespace"
                  label="Namespace"
                  fullWidth
                  value={namespace}
                  error={namespaceError}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange('namespace')}
                />
              </Grid>
              {filteredOps.map((val) => (
                <Grid item xs={12} md={4}>
                  {this.generateCardForCategory(val)}
                </Grid>
              ))}
            </Grid>
          </div>
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryAdapterPlayComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  adapter: PropTypes.object.isRequired,
};

// const mapStateToProps = (state) => {
//   const smi_result = state.get('smi_result').toJS();
//   return { smi_result, };
// };

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  // updateSMIResults: bindActionCreators(updateSMIResults, dispatch),
});

export default withStyles(styles)(connect(
  // mapStateToProps,
  null,
  mapDispatchToProps,
)(withRouter(withSnackbar(MesheryAdapterPlayComponent))));