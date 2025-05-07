import React, { useState, useRef, useEffect } from 'react';
import {
  Box,
  Card,
  CardActions,
  CardHeader,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  FormControl,
  FormControlLabel,
  FormGroup,
  Grid,
  IconButton,
  Menu,
  MenuItem,
  Switch,
  Table,
  Tooltip,
  styled,
  FormLabel,
  TableBody,
  TableCell,
  TableRow,
  NoSsr,
  TableHead,
} from '@layer5/sistent';
import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import PlayIcon from '@mui/icons-material/PlayArrow';
import MUIDataTable from 'mui-datatables';
import { useRouter } from 'next/router';
import PropTypes from 'prop-types';
import { Controlled as CodeMirror } from 'react-codemirror2';
import Moment from 'react-moment';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch from '../lib/data-fetch';
import { setK8sContexts, updateProgress } from '../lib/store';
import { ctxUrl, getK8sClusterIdsFromCtxId } from '../utils/multi-ctx';
import fetchAvailableAddons from './graphql/queries/AddonsStatusQuery';
import fetchAvailableNamespaces from './graphql/queries/NamespaceQuery';
import MesheryMetrics from './MesheryMetrics';
import MesheryResultDialog from './MesheryResultDialog';
import ReactSelectWrapper from './ReactSelectWrapper';
import ConfirmationMsg from './ConfirmationModal';
import { iconMedium } from '../css/icons.styles';
import { ACTIONS } from '../utils/Enum';
import { getModelByName } from '../api/meshmodel';
import { EVENT_TYPES } from '../lib/event-types';
import { withNotify } from '../utils/hooks/useNotification';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';

export const AdapterChip = styled(Chip)(({ theme }) => ({
  height: '50px',
  fontSize: '15px',
  position: 'relative',
  top: theme.spacing(0.5),
  [theme.breakpoints.down('md')]: {
    fontSize: '12px',
  },
}));

const AdapterTableHeader = styled(TableCell)({
  fontWeight: 'bolder',
  fontSize: 18,
});

const AdapterSmWrapper = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
}));

const SecondaryTable = styled('div')({
  borderRadius: 10,
  backgroundColor: '#f7f7f7',
});

const PaneSection = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.tabs,
  padding: theme.spacing(3),
  borderRadius: 4,
}));

const ChipNamespaceContainer = styled(Grid)(() => ({
  gap: '2rem',
  margin: '0px',
}));

const InputWrapper = styled('div')(() => ({
  flex: '1',
  minWidth: '250px',
}));

const AdapterCard = styled(Card)(() => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
}));

const MesheryAdapterPlayComponent = (props) => {
  const { adapter, updateProgress, notify, selectedK8sContexts, k8sconfig, grafana } = props;

  const router = useRouter();
  const cmEditorAddRef = useRef(null);
  const cmEditorDelRef = useRef(null);
  const addIconEles = useRef({});
  const delIconEles = useRef({});

  const [cmEditorValAdd, setCmEditorValAdd] = useState('');
  const [, setCmEditorValAddError] = useState(false);
  const [cmEditorValDel, setCmEditorValDel] = useState('');
  const [, setCmEditorValDelError] = useState(false);
  const [, setSelectionError] = useState(false);
  const [namespace, setNamespace] = useState({ value: 'default', label: 'default' });
  const [namespaceError, setNamespaceError] = useState(false);
  const [customDialogAdd, setCustomDialogAdd] = useState(false);
  const [customDialogDel, setCustomDialogDel] = useState(false);
  const [customDialogSMI, setCustomDialogSMI] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [addonSwitchGroup, setAddonSwitchGroup] = useState({});
  const [smi_result, setSmiResult] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [pageSize, setPageSize] = useState(10);
  const [namespaceList, setNamespaceList] = useState([]);
  const [namespaceSubscription, setNamespaceSubscription] = useState(null);
  const [, setActiveContexts] = useState([]);
  // const [deployModalOpen, setDeployModalOpen] = useState(false);
  const [category, setCategory] = useState(0);
  const [selectedOp, setSelectedOp] = useState('');
  const [isDeleteOp, setIsDeleteOp] = useState(false);
  const [operationName, setOperationName] = useState('');
  const [versionList, setVersionList] = useState([]);
  const [version, setVersion] = useState({ label: '', value: '' });
  const [versionError, setVersionError] = useState(false);

  // Initialize menu state based on adapter operations
  const initMenuState = () => {
    const menuState = {};
    if (adapter && adapter.ops) {
      [0, 1, 2, 3, 4].forEach((i) => {
        menuState[i] = {
          add: false,
          delete: false,
        };
      });
    }
    return menuState;
  };

  const [menuState, setMenuState] = useState(initMenuState());

  // Equivalent to componentDidMount
  useEffect(() => {
    const meshname = mapAdapterNameToMeshName(adapter.name);
    const variables = { type: meshname, k8sClusterIDs: getK8sClusterIds() };

    initSubscription();
    getMeshVersions();

    if (selectedK8sContexts) {
      if (selectedK8sContexts.includes('all')) {
        let active = [];
        k8sconfig.forEach((ctx) => {
          active.push(ctx.contextID);
        });
        setActiveContexts(active);
      } else {
        setActiveContexts(selectedK8sContexts);
      }
    }

    fetchAvailableAddons(variables).subscribe({
      next: (res) => {
        setAddonsState(res);
      },
      error: (err) => console.log('error at addon fetch: ' + err),
    });

    // Cleanup function (componentWillUnmount)
    return () => {
      disposeSubscriptions();
    };
  }, [adapter.name]);

  // Equivalent to componentDidUpdate for selectedK8sContexts
  useEffect(() => {
    if (props.selectedK8sContexts) {
      disposeSubscriptions();
      initSubscription();
    }
  }, [props.selectedK8sContexts]);

  // Equivalent to componentDidUpdate for adapter.name
  useEffect(() => {
    getMeshVersions();
  }, [props.adapter.name]);

  const initSubscription = () => {
    disposeSubscriptions();

    const subscription = fetchAvailableNamespaces({
      k8sClusterIDs: getK8sClusterIds(),
    }).subscribe({
      next: (res) => {
        let namespaces = [];
        res?.namespaces?.map((ns) => {
          namespaces.push({
            value: ns?.namespace,
            label: ns?.namespace,
          });
        });
        if (namespaces.length === 0) {
          namespaces.push({
            value: 'default',
            label: 'default',
          });
        }
        namespaces.sort((a, b) => (a.value > b.value ? 1 : -1));
        setNamespaceList(namespaces);
      },
      error: (err) => console.log('error at namespace fetch: ' + err),
    });

    setNamespaceSubscription(subscription);
  };

  const disposeSubscriptions = () => {
    if (namespaceSubscription) {
      namespaceSubscription.unsubscribe();
    }
  };

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sconfig);
  };

  const mapAdapterNameToMeshName = (name) => {
    if (name?.toLowerCase() === 'istio') return 'ISTIO';
    return 'ALL';
  };

  const getMeshVersions = () => {
    const activeMesh = adapter?.name;
    getModelByName(activeMesh.toLowerCase()).then((res) => {
      let uniqueVersions = [...new Set(res?.models?.map((model) => model?.version))].reverse();
      if (uniqueVersions.length === 0) {
        uniqueVersions = [''];
      }
      let versionList = uniqueVersions.map((version) => ({ value: version, label: version }));
      setVersionList(versionList);
      setVersion(versionList[0]);
    });
  };

  const setAddonsState = (data) => {
    const meshname = adapter.name;
    const localState = {};
    data?.addonsState?.forEach((addon) => {
      if (addon.owner === meshname) {
        const name = addon.name !== 'jaeger-collector' ? addon.name : 'jaeger';
        localState[`${name}-addon`] = true;
      }
    });
    setAddonSwitchGroup(localState);
  };

  // const handleChange =
  //   (name, isDelete = false) =>
  //   (event) => {
  //     if (name === 'selectedOp' && event.target.value !== '') {
  //       if (event.target.value === 'custom') {
  //         if (isDelete) {
  //           if (cmEditorValDel !== '' && cmEditorDelRef.current?.state.lint.marked.length === 0) {
  //             setSelectionError(false);
  //             setCmEditorValDelError(false);
  //           }
  //         } else if (
  //           cmEditorValAdd !== '' &&
  //           cmEditorAddRef.current?.state.lint.marked.length === 0
  //         ) {
  //           setSelectionError(false);
  //           setCmEditorValAddError(false);
  //         }
  //       } else {
  //         setSelectionError(false);
  //       }
  //     }

  //     if (name === 'selectedOp') {
  //       setSelectedOp(event.target.value);
  //     }
  //   };

  const handleNamespaceChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
      setNamespace(newValue);
      setNamespaceError(false);
    } else {
      setNamespaceError(true);
    }
  };

  const handleVersionChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
      setVersion(newValue);
      setVersionError(false);
    } else {
      setVersionError(true);
    }
  };

  const handleModalClose = (isDelete) => () => {
    if (isDelete) {
      setCustomDialogDel(false);
    } else {
      setCustomDialogAdd(false);
    }
  };

  const handleSMIClose = () => {
    setCustomDialogSMI(false);
  };

  const resetSelectedRowData = () => {
    setSelectedRowData(null);
  };

  // const handleModalOpen = (isDelete) => () => {
  //   if (isDelete) {
  //     setCustomDialogDel(true);
  //   } else {
  //     setCustomDialogAdd(true);
  //   }
  // };

  const handleOpen = () => {
    setModalOpen(true);
  };

  const handleClose = () => {
    setModalOpen(false);
  };

  // const handleDeployModalOpen = () => {
  //   setDeployModalOpen(true);
  // };

  // const handleDeployModalClose = () => {
  //   setDeployModalOpen(false);
  // };

  const handleSubmit =
    (cat, op, deleteOp = false) =>
    () => {
      handleOpen();

      const filteredOp = adapter.ops.filter(({ key }) => key === op);
      if (op === '' || typeof filteredOp === 'undefined' || filteredOp.length === 0) {
        setSelectionError(true);
        return;
      }

      if (deleteOp) {
        if (
          op === 'custom' &&
          (cmEditorValDel === '' || cmEditorDelRef.current?.state.lint.marked.length > 0)
        ) {
          setCmEditorValDelError(true);
          setSelectionError(true);
          return;
        }
      } else if (
        op === 'custom' &&
        (cmEditorValAdd === '' || cmEditorAddRef.current?.state.lint.marked.length > 0)
      ) {
        setCmEditorValAddError(true);
        setSelectionError(true);
        return;
      }

      if (namespace.value === '') {
        setNamespaceError(true);
        return;
      }

      if (version?.value === '') {
        setVersionError(true);
        return;
      }

      const opName = op
        .replaceAll('_', ' ')
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');

      setCategory(cat);
      setSelectedOp(op);
      setIsDeleteOp(deleteOp);
      setOperationName(opName);
    };

  const submitOp = (cat, op, deleteOp = false) => {
    const data = {
      adapter: adapter.adapter_location,
      query: op,
      namespace: namespace.value,
      customBody: deleteOp ? cmEditorValDel : cmEditorValAdd,
      deleteOp: deleteOp ? 'on' : '',
      version: version.value,
    };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join('&');

    updateProgress({ showProgress: true });
    handleClose();

    dataFetch(
      ctxUrl('/api/system/adapter/operation', selectedK8sContexts),
      {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: params,
      },
      (result) => {
        updateProgress({ showProgress: false });

        const newMenuState = { ...menuState };
        newMenuState[cat][deleteOp ? 'delete' : 'add'] = false;
        setMenuState(newMenuState);

        if (deleteOp) {
          setCustomDialogDel(false);
        } else {
          setCustomDialogAdd(false);
        }

        if (typeof result !== 'undefined') {
          notify({ message: 'Operation executing...', event_type: EVENT_TYPES.INFO });
        }
      },
      handleError(cat, deleteOp, op),
    );
  };

  const handleAdapterClick = (adapterLoc) => () => {
    updateProgress({ showProgress: true });

    dataFetch(
      `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        credentials: 'include',
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({ message: 'Adapter pinged!', event_type: EVENT_TYPES.SUCCESS });
        }
      },
      handleError('Could not ping adapter.'),
    );
  };

  const fetchSMIResults = (adapterName, page, pageSize, search, sortOrder) => {
    let query = '';
    if (typeof search === 'undefined' || search === null) {
      search = '';
    }
    if (typeof sortOrder === 'undefined' || sortOrder === null) {
      sortOrder = '';
    }
    query = `?page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
      search,
    )}&order=${encodeURIComponent(sortOrder)}`;

    dataFetch(
      `/api/smi/results${query}`,
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (typeof result !== 'undefined' && result.results) {
          const results = result.results.filter(
            (val) => val.mesh_name.toLowerCase() == adapterName.toLowerCase(),
          );
          setSmiResult({ ...result, results: results, total_count: results.length });
        }
      },
      (error) => console.log('Could not fetch SMI results.', error),
    );
  };

  // const handleSMIClick = (adapterName) => () => {
  //   updateProgress({ showProgress: true });
  //   fetchSMIResults(adapterName, page, pageSize, search, sortOrder);
  //   updateProgress({ showProgress: false });
  //   setCustomDialogSMI(true);
  // };

  const handleError = (cat, deleteOp, selectedOp) => {
    return (error) => {
      if (cat && deleteOp !== undefined) {
        const newMenuState = { ...menuState };
        newMenuState[cat][deleteOp ? 'delete' : 'add'] = false;
        setMenuState(newMenuState);

        if (deleteOp) {
          setCustomDialogDel(false);
        } else {
          setCustomDialogAdd(false);
        }
      }

      if (selectedOp) {
        setAddonSwitchGroup({ ...addonSwitchGroup, [selectedOp]: deleteOp });
      }

      updateProgress({ showProgress: false });
      notify({
        message: `Operation submission failed: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  };

  // const activeContextChangeCallback = (activeK8sContexts) => {
  //   if (activeK8sContexts.includes('all')) {
  //     activeK8sContexts = ['all'];
  //   }
  //   setK8sContexts({
  //     type: actionTypes.SET_K8S_CONTEXT,
  //     selectedK8sContexts: activeK8sContexts,
  //   });
  // };

  const addDelHandleClick = (cat, isDelete) => () => {
    const newMenuState = { ...menuState };
    newMenuState[cat][isDelete ? 'delete' : 'add'] =
      !newMenuState[cat][isDelete ? 'delete' : 'add'];
    setMenuState(newMenuState);

    if (cat === 4) {
      if (isDelete) {
        setCustomDialogDel(!customDialogDel);
      } else {
        setCustomDialogAdd(!customDialogAdd);
      }
    }
  };

  const generateMenu = (cat, isDelete, selectedAdapterOps) => {
    const ele = !isDelete ? addIconEles.current[cat] : delIconEles.current[cat];
    return (
      <Menu
        id="long-menu"
        anchorEl={ele}
        keepMounted
        open={menuState[cat][isDelete ? 'delete' : 'add']}
        onClose={addDelHandleClick(cat, isDelete)}
      >
        {selectedAdapterOps
          .sort((adap1, adap2) => adap1.value.localeCompare(adap2.value))
          .map(({ key, value }) => (
            <MenuItem
              key={`${key}_${new Date().getTime()}`}
              onClick={handleSubmit(cat, key, isDelete)}
            >
              {value}
            </MenuItem>
          ))}
      </Menu>
    );
  };

  const generateYAMLEditor = (cat, isDelete) => {
    return (
      <Dialog
        onClose={handleModalClose(isDelete)}
        aria-labelledby="adapter-dialog-title"
        open={isDelete ? customDialogDel : customDialogAdd}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle id="adapter-dialog-title" onClose={handleModalClose(isDelete)}>
          {adapter.name} Adapter - Custom YAML
          {isDelete ? '(delete)' : ''}
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent>
          <Grid container spacing={5}>
            <Grid item xs={6}>
              <ReactSelectWrapper
                label="Namespace"
                value={namespace}
                error={namespaceError}
                options={namespaceList}
                onChange={handleNamespaceChange}
              />
            </Grid>
            <Grid item xs={6}>
              <ReactSelectWrapper
                label="Version"
                value={version}
                error={versionError}
                options={versionList}
                onChange={handleVersionChange}
              />
            </Grid>
            <Grid item xs={12}>
              <CodeMirror
                editorDidMount={(editor) => {
                  if (isDelete) {
                    cmEditorDelRef.current = editor;
                  } else {
                    cmEditorAddRef.current = editor;
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
                    setCmEditorValDel(value);
                    if (value !== '' && editor.state.lint.marked.length === 0) {
                      setSelectionError(false);
                      setCmEditorValDelError(false);
                    }
                  } else {
                    setCmEditorValAdd(value);
                    if (value !== '' && editor.state.lint.marked.length === 0) {
                      setSelectionError(false);
                      setCmEditorValAddError(false);
                    }
                  }
                }}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <IconButton aria-label="Apply" onClick={handleSubmit(cat, 'custom', isDelete)}>
            {!isDelete && <PlayIcon style={iconMedium} />}
            {isDelete && <DeleteIcon style={iconMedium} />}
          </IconButton>
        </DialogActions>
      </Dialog>
    );
  };

  const generateSMIResult = () => {
    const smi_columns = [
      {
        name: 'ID',
        label: 'ID',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            );
          },
          customBodyRender: (value) => (
            <Tooltip title={value} placement="top">
              <div>{value.slice(0, 5) + '...'}</div>
            </Tooltip>
          ),
        },
      },
      {
        name: 'Date',
        label: 'Date',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            );
          },
          customBodyRender: (value) => <Moment format="LLLL">{value}</Moment>,
        },
      },
      {
        name: 'Service Mesh',
        label: 'Service Mesh',
        options: {
          filter: true,
          sort: true,
          searchable: true,
          customHeadRender: ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            );
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
          customHeadRender: ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            );
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
          customHeadRender: ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            );
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
          customHeadRender: ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>
            );
          },
        },
      },
    ];

    const smi_options = {
      sort: !(props.user && props.user.user_id === 'meshery'),
      search: !(props.user && props.user.user_id === 'meshery'),
      filterType: 'textField',
      expandableRows: true,
      selectableRows: 'none',
      rowsPerPage: pageSize,
      rowsPerPageOptions: [10, 20, 25],
      fixedHeader: true,
      print: false,
      download: false,
      renderExpandableRow: (rowData, rowMeta) => {
        const column = [
          'Specification',
          'Assertions',
          'Time',
          'Version',
          'Capability',
          'Result',
          'Reason',
        ];
        const data = smi_result.results[rowMeta.dataIndex].more_details.map((val) => {
          return [
            val.smi_specification,
            val.assertions,
            val.time,
            val.smi_version,
            val.capability,
            val.status,
            val.reason,
          ];
        });
        const colSpan = rowData.length + 1;
        return (
          <TableRow>
            <TableCell colSpan={colSpan}>
              <SecondaryTable>
                <Table aria-label="a dense table">
                  <TableHead>
                    <TableRow>
                      {column.map((val) => (
                        <TableCell colSpan={colSpan} key={val}>
                          {val}
                        </TableCell>
                      ))}
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {data.map((row, idx) => (
                      <TableRow key={idx}>
                        {row.map((val, vIdx) => {
                          if (val && val.match(/[0-9]+m[0-9]+.+[0-9]+s/i) != null) {
                            const time = val.split(/m|s/);
                            return (
                              <TableCell colSpan={colSpan} key={`${idx}-${vIdx}`}>
                                {time[0] + 'm ' + parseFloat(time[1]).toFixed(1) + 's'}
                              </TableCell>
                            );
                          } else {
                            return (
                              <TableCell colSpan={colSpan} key={`${idx}-${vIdx}`}>
                                {val}
                              </TableCell>
                            );
                          }
                        })}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </SecondaryTable>
            </TableCell>
          </TableRow>
        );
      },
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        let order = '';
        if (tableState.activeColumn) {
          order = `${smi_columns[tableState.activeColumn].name} desc`;
        }

        switch (action) {
          case 'changePage':
            fetchSMIResults(adapter.name, tableState.page, pageSize, search, sortOrder);
            setPage(tableState.page);
            break;
          case 'changeRowsPerPage':
            fetchSMIResults(adapter.name, page, tableState.rowsPerPage, search, sortOrder);
            setPageSize(tableState.rowsPerPage);
            break;
          case 'search':
            if (search !== tableState.searchText) {
              fetchSMIResults(
                adapter.name,
                page,
                pageSize,
                tableState.searchText !== null ? tableState.searchText : '',
                sortOrder,
              );
              setSearch(tableState.searchText);
            }
            break;
          case 'sort':
            if (sortInfo.length === 2) {
              if (sortInfo[1] === 'ascending') {
                order = `${smi_columns[tableState.activeColumn].name} asc`;
              } else {
                order = `${smi_columns[tableState.activeColumn].name} desc`;
              }
            }
            if (order !== sortOrder) {
              fetchSMIResults(adapter.name, page, pageSize, search, order);
              setSortOrder(order);
            }
            break;
        }
      },
    };

    var data = [];
    if (smi_result && smi_result.results) {
      data = smi_result.results.map((val) => {
        return [
          val.id,
          val.date,
          val.mesh_name,
          val.mesh_version,
          val.passing_percentage,
          val.status,
        ];
      });
    }

    return (
      <Dialog
        onClose={handleSMIClose}
        aria-labelledby="adapter-dialog-title"
        open={customDialogSMI}
        fullWidth
        maxWidth="md"
      >
        <MUIDataTable
          title={
            <AdapterTableHeader>Service Mesh Interface Conformance Results</AdapterTableHeader>
          }
          data={data}
          columns={smi_columns}
          options={smi_options}
        />
      </Dialog>
    );
  };

  const extractAddonOperations = (addonOpsCat) => {
    return adapter.ops.filter(
      ({ category, value }) => category === addonOpsCat && value?.startsWith('Add-on:'),
    );
  };

  const generateAddonSwitches = (selectedAdapterOps) => {
    if (!selectedAdapterOps.length) return null;

    return (
      <FormControl component="fieldset" style={{ padding: '1rem' }}>
        <FormLabel component="legend">Customize Addons</FormLabel>
        <FormGroup>
          {selectedAdapterOps
            .map((ops) => ({ ...ops, value: ops.value.replace('Add-on:', '') }))
            .sort((ops1, ops2) => ops1.value.localeCompare(ops2.value))
            .map((ops) => (
              <FormControlLabel
                control={
                  <Switch
                    color="primary"
                    checked={!!addonSwitchGroup[ops.key]}
                    onChange={(ev) => {
                      setAddonSwitchGroup({
                        ...addonSwitchGroup,
                        [ev.target.name]: ev.target.checked,
                      });
                      submitOp(ops.category, ops.key, !addonSwitchGroup[ops.key]);
                    }}
                    name={ops.key}
                  />
                }
                label={ops.value}
                key={ops.key}
              />
            ))}
        </FormGroup>
      </FormControl>
    );
  };

  const generateCardForCategory = (cat) => {
    if (typeof cat === 'undefined') {
      cat = 0;
    }

    let selectedAdapterOps =
      adapter && adapter.ops
        ? adapter.ops.filter(
            ({ category }) => (typeof category === 'undefined' && cat === 0) || category === cat,
          )
        : [];
    let content;
    let description;
    let permission;
    switch (cat) {
      case 0:
        content = 'Manage Cloud Native Infrastructure Lifecycle';
        description = 'Deploy cloud native infrastructure or SMI adapter on your cluster.';
        permission = {
          action: keys.MANAGE_CLOUD_NATIVE_INFRASTRUCTURE_LIFE_CYCLE.action,
          subject: keys.MANAGE_CLOUD_NATIVE_INFRASTRUCTURE_LIFE_CYCLE.subject,
        };
        break;

      case 1:
        content = 'Manage Sample Application Lifecycle';
        description = 'Deploy sample applications on/off the service mesh.';
        permission = {
          action: keys.MANAGE_CLOUD_NATIVE_INFRASTRUCTURE_LIFE_CYCLE.action,
          subject: keys.MANAGE_CLOUD_NATIVE_INFRASTRUCTURE_LIFE_CYCLE.subject,
        };
        break;

      case 2:
        content = 'Apply Cloud Native Infrastructure Configuration';
        description = 'Configure your cloud native infrastructure using some pre-defined options.';
        selectedAdapterOps = selectedAdapterOps.filter((ops) => !ops.value.startsWith('Add-on:'));
        permission = {
          action: keys.APPLY_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION.action,
          subject: keys.APPLY_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION.subject,
        };
        break;

      case 3:
        content = 'Validate Cloud Native Infrastructure Configuration';
        description =
          'Validate your cloud native infrastructure configuration against best practices.';
        permission = {
          action: keys.VALIDATE_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION.action,
          subject: keys.VALIDATE_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION.subject,
        };
        break;

      case 4:
        content = 'Apply Custom Configuration';
        description = 'Customize the configuration of your cloud native infrastructure.';
        permission = {
          action: keys.APPLY_CUSTOM_CLOUD_NATIVE_CONFIGURATION.action,
          subject: keys.APPLY_CUSTOM_CLOUD_NATIVE_CONFIGURATION.subject,
        };
        break;
    }

    return (
      <AdapterCard>
        <CardHeader title={content} subheader={description} style={{ flexGrow: 1 }} />
        <CardActions disableSpacing>
          <IconButton
            aria-label="install"
            ref={(ch) => (addIconEles.current[cat] = ch)}
            onClick={addDelHandleClick(cat, false)}
            disabled={!CAN(permission.action, permission.subject)}
          >
            {cat !== 4 ? <AddIcon style={iconMedium} /> : <PlayIcon style={iconMedium} />}
          </IconButton>
          {cat !== 4 && generateMenu(cat, false, selectedAdapterOps)}
          {cat === 4 && generateYAMLEditor(cat, false)}
          {cat !== 3 && (
            <Box width={'100%'}>
              <IconButton
                aria-label="delete"
                ref={(ch) => (delIconEles.current[cat] = ch)}
                style={{ float: 'right' }}
                onClick={addDelHandleClick(cat, true)}
                disabled={!CAN(permission.action, permission.subject)}
              >
                <DeleteIcon style={iconMedium} />
              </IconButton>
              {cat !== 4 && generateMenu(cat, true, selectedAdapterOps)}
              {cat === 4 && generateYAMLEditor(cat, true)}
            </Box>
          )}
        </CardActions>
      </AdapterCard>
    );
  };

  const renderGrafanaCustomCharts = (boardConfigs, grafanaURL, grafanaAPIKey) => {
    return (
      <MesheryMetrics
        boardConfigs={boardConfigs}
        grafanaAPIKey={grafanaAPIKey}
        grafanaURL={grafanaURL}
        handleGrafanaChartAddition={() => router.push('/settings?settingsCategory=Metrics')}
      />
    );
  };

  // Render component
  let adapterName = adapter.name.split(' ').join('').toLowerCase();
  let imageSrc = '/static/img/' + adapterName + '.svg';
  let adapterChip = (
    <AdapterChip
      label={adapter.adapter_location}
      data-cy="adapter-chip-ping"
      onClick={handleAdapterClick(adapter.adapter_location)}
      icon={<img src={imageSrc} width={'1.25rem'} />}
      variant="outlined"
    />
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
      {selectedRowData && selectedRowData !== null && Object.keys(selectedRowData).length > 0 && (
        <MesheryResultDialog rowData={selectedRowData} close={resetSelectedRowData} />
      )}
      <React.Fragment>
        <AdapterSmWrapper>
          <Grid container spacing={2} direction="row" alignItems="flex-start">
            {/* SECTION 1 */}
            <Grid item xs={12}>
              <PaneSection>
                <Grid container spacing={4}>
                  <ChipNamespaceContainer
                    container
                    item
                    xs={12}
                    alignItems="flex-start"
                    justify="space-between"
                  >
                    <div>{adapterChip}</div>
                    <InputWrapper>
                      <ReactSelectWrapper
                        label="Namespace"
                        value={namespace}
                        error={namespaceError}
                        options={namespaceList}
                        onChange={handleNamespaceChange}
                      />
                    </InputWrapper>
                    <InputWrapper>
                      <ReactSelectWrapper
                        label="Version"
                        value={version}
                        error={versionError}
                        options={versionList}
                        onChange={handleVersionChange}
                      />
                    </InputWrapper>
                  </ChipNamespaceContainer>
                  <Grid container spacing={1} style={{ margin: '1rem' }}>
                    <Grid
                      container
                      item
                      lg={!extractAddonOperations(2).length ? 12 : 10}
                      xs={12}
                      spacing={2}
                    >
                      {filteredOps.map((val, i) => (
                        <Grid item lg={3} md={4} xs={12} key={`adapter-card-${i}`}>
                          {generateCardForCategory(val)}
                        </Grid>
                      ))}
                    </Grid>
                    <Grid container item lg={2} xs={12}>
                      <Grid item xs={12} md={4}>
                        {generateAddonSwitches(extractAddonOperations(2))}
                      </Grid>
                    </Grid>
                  </Grid>
                </Grid>
              </PaneSection>
            </Grid>
            {/* SECTION 2 */}
            <Grid item xs={12}>
              <PaneSection>
                {renderGrafanaCustomCharts(
                  grafana.selectedBoardsConfigs,
                  grafana.grafanaURL,
                  grafana.grafanaAPIKey,
                )}
              </PaneSection>
            </Grid>
          </Grid>
        </AdapterSmWrapper>
        <ConfirmationMsg
          open={modalOpen}
          handleClose={handleClose}
          submit={{
            deploy: () => submitOp(category, selectedOp, false),
            unDeploy: () => submitOp(category, selectedOp, true),
          }}
          isDelete={isDeleteOp}
          title={operationName}
          tab={isDeleteOp ? ACTIONS.UNDEPLOY : ACTIONS.DEPLOY}
        />
        {customDialogSMI && generateSMIResult()}
      </React.Fragment>
    </NoSsr>
  );
};

MesheryAdapterPlayComponent.propTypes = {
  adapter: PropTypes.object.isRequired,
};

const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  const k8sconfig = st.get('k8sConfig');
  const selectedK8sContexts = st.get('selectedK8sContexts');

  return { grafana: { ...grafana, ts: new Date(grafana.ts) }, selectedK8sContexts, k8sconfig };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  setK8sContexts: bindActionCreators(setK8sContexts, dispatch),
});

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withNotify(MesheryAdapterPlayComponent));
