import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  NoSsr,
  TableCell,
  Button,
  Tooltip,
  FormControl,
  Select,
  TableContainer,
  Table,
  Grid,
  TableRow,
  IconButton,
  Typography,
  Switch,
  Popover,
  AppBar,
  Tabs,
  Tab,
  MenuItem,
  Box,
  Chip,
} from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import Moment from 'react-moment';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import { updateProgress } from '../../lib/store';
import dataFetch from '../../lib/data-fetch';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import CustomColumnVisibilityControl from '../../utils/custom-column';
import SearchBar from '../../utils/custom-search';
import { ResponsiveDataTable } from '@layer5/sistent-components';
import useStyles from '../../assets/styles/general/tool.styles';
import Modal from '../Modal';
import { iconMedium } from '../../css/icons.styles';
import PromptComponent, { PROMPT_VARIANTS } from '../PromptComponent';
import resetDatabase from '../graphql/queries/ResetDatabaseQuery';
import changeOperatorState from '../graphql/mutations/OperatorStatusMutation';
import fetchMesheryOperatorStatus from '../graphql/queries/OperatorStatusQuery';
import MesherySettingsEnvButtons from '../MesherySettingsEnvButtons';
import styles from './styles';
import MeshSyncTable from './meshSync';
import ConnectionIcon from '../../assets/icons/Connection';
import MeshsyncIcon from '../../assets/icons/Meshsync';
import classNames from 'classnames';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import SyncIcon from '@mui/icons-material/Sync';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import {
  CONNECTION_KINDS,
  CONNECTION_STATES,
  CONTROLLERS,
  CONTROLLER_STATES,
} from '../../utils/Enum';
import FormatConnectionMetadata from './metadata';
import useKubernetesHook from '../hooks/useKubernetesHook';
import theme from '../../themes/app';
import { TootltipWrappedConnectionChip } from './ConnectionChip';
import InfoIcon from '@material-ui/icons/Info';
import { SortableTableCell } from './common';
import { getColumnValue, getVisibilityColums } from '../../utils/utils';
import HandymanIcon from '@mui/icons-material/Handyman';
import NotInterestedRoundedIcon from '@mui/icons-material/NotInterestedRounded';
import DisconnectIcon from '../../assets/icons/disconnect';
import { updateVisibleColumns } from '../../utils/responsive-column';
import { useWindowDimensions } from '../../utils/dimension';
import MultiSelectWrapper from '../multi-select-wrapper';
import {
  useGetEnvironmentsQuery,
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
  useSaveEnvironmentMutation,
} from '../../rtk-query/environments';
import ErrorBoundary from '../ErrorBoundary';
import { store } from '../../store';
import { Provider } from 'react-redux';

const ACTION_TYPES = {
  FETCH_CONNECTIONS: {
    name: 'FETCH_CONNECTIONS',
    error_msg: 'Failed to fetch connections',
  },
  UPDATE_CONNECTION: {
    name: 'UPDATE_CONNECTION',
    error_msg: 'Failed to update connection',
  },
  DELETE_CONNECTION: {
    name: 'DELETE_CONNECTION',
    error_msg: 'Failed to delete connection',
  },
  FETCH_CONNECTION_STATUS_TRANSITIONS: {
    name: 'FETCH_CONNECTION_STATUS_TRANSITIONS',
    error_msg: 'Failed to fetch connection transitions',
  },
  FETCH_ENVIRONMENT: {
    name: 'FETCH_ENVIRONMENT',
    error_msg: 'Failed to fetch environment',
  },
  CREATE_ENVIRONMENT: {
    name: 'CREATE_ENVIRONMENT',
    error_msg: 'Failed to create environment',
  },
};

/**
 * Parent Component for Connection Component
 *
 * @important
 * - Keep the component's responsibilities focused on connection management. Avoid adding unrelated functionality and state.
 */

function ConnectionManagementPage(props) {
  const [createConnectionModal, setCreateConnectionModal] = useState({
    open: false,
  });
  const [createConnection, setCreateConnection] = useState({});

  const handleCreateConnectionModalOpen = () => {
    setCreateConnectionModal({ open: true });
  };

  const handleCreateConnectionModalClose = () => {
    setCreateConnectionModal({ open: false });
  };

  const handleCreateConnectionSubmit = () => {};

  useEffect(() => {
    dataFetch(
      '/api/schema/resource/helmRepo',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        setCreateConnection(result);
      },
    );
  }, []);

  return (
    <>
      <Connections
        createConnectionModal={createConnectionModal}
        onOpenCreateConnectionModal={handleCreateConnectionModalOpen}
        onCloseCreateConnectionModal={handleCreateConnectionModalClose}
        {...props}
      />
      {createConnectionModal.open && (
        <Modal
          open={true}
          schema={createConnection.rjsfSchema}
          uiSchema={createConnection.uiSchema}
          handleClose={handleCreateConnectionModalClose}
          handleSubmit={handleCreateConnectionSubmit}
          title="Connect Helm Repository"
          submitBtnText="Connect"
        />
      )}
    </>
  );
}
function Connections(props) {
  const {
    classes,
    updateProgress,
    /*onOpenCreateConnectionModal,*/ operatorState,
    selectedK8sContexts,
    k8sconfig,
    connectionMetadataState,
    meshsyncControllerState,
    organization,
  } = props;
  console.log('props: ', props);
  const modalRef = useRef(null);
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [connections, setConnections] = useState([]);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [showMore, setShowMore] = useState(false);
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [loading, setLoading] = useState(false);
  const [rowData, setSelectedRowData] = useState({});
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [_operatorState, _setOperatorState] = useState(operatorState || []);
  const [tab, setTab] = useState(0);
  const ping = useKubernetesHook();
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();
  const [saveEnvironmentMutator] = useSaveEnvironmentMutation();

  const addConnectionToEnvironment = async (
    environmentId,
    environmentName,
    connectionId,
    connectionName,
  ) => {
    addConnectionToEnvironmentMutator({ environmentId, connectionId })
      .unwrap()
      .then(() => {
        getConnections(page, pageSize, search, sortOrder);
        notify({
          message: `Connection: ${connectionName} assigned to environment: ${environmentName}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch((err) => {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      });
  };

  const removeConnectionFromEnvironment = (
    environmentId,
    environmentName,
    connectionId,
    connectionName,
  ) => {
    removeConnectionFromEnvMutator({ environmentId, connectionId })
      .unwrap()
      .then(() => {
        getConnections(page, pageSize, search, sortOrder);
        notify({
          message: `Connection: ${connectionName} removed from environment: ${environmentName}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch((err) => {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      });
  };

  const saveEnvironment = (connectionId, connectionName, environmentName) => {
    saveEnvironmentMutator({
      body: {
        name: environmentName,
        organization_id: organization?.id,
      },
    })
      .unwrap()
      .then((resp) => {
        notify({
          message: `Environment: ${resp.Name} saved`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        addConnectionToEnvironment(resp.id, resp.name, connectionId, connectionName);
        getConnections(page, pageSize, search, sortOrder);
      })
      .catch((err) => {
        notify({
          message: `${ACTION_TYPES.CREATE_ENVIRONMENT.error_msg}: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      });
  };

  const {
    data: environmentsResponse,
    isSuccess: isEnvironmentsSuccess,
    isError: isEnvironmentsError,
    error: environmentsError,
  } = useGetEnvironmentsQuery(
    { orgId: organization?.id },
    {
      skip: !organization?.id,
    },
  );
  let environments = environmentsResponse?.environments || [];

  const open = Boolean(anchorEl);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;
  const meshSyncResetRef = useRef(null);
  const { notify } = useNotification();
  const StyleClass = useStyles();
  const url = `https://docs.meshery.io/concepts/connections`;

  const icons = {
    [CONNECTION_STATES.IGNORED]: () => <RemoveCircleIcon />,
    [CONNECTION_STATES.CONNECTED]: () => <CheckCircleIcon />,
    [CONNECTION_STATES.REGISTERED]: () => <AssignmentTurnedInIcon />,
    [CONNECTION_STATES.DISCOVERED]: () => <ExploreIcon />,
    [CONNECTION_STATES.DELETED]: () => <DeleteForeverIcon />,
    [CONNECTION_STATES.MAINTENANCE]: () => <HandymanIcon />,
    [CONNECTION_STATES.DISCONNECTED]: () => (
      <DisconnectIcon fill="#E75225" width={24} height={24} />
    ),
    [CONNECTION_STATES.NOTFOUND]: () => <NotInterestedRoundedIcon />,
  };

  const handleEnvironmentSelect = (
    connectionId,
    connName,
    assignedEnvironments,
    selectedEnvironments,
    unSelectedEnvironments,
  ) => {
    let newlySelectedEnvironments = selectedEnvironments.filter((env) => {
      return !assignedEnvironments.some((assignedEnv) => assignedEnv.value === env.value);
    });

    newlySelectedEnvironments.forEach((environment) => {
      let envName = environment.label;
      let environmentId = environment.value || '';
      let isNew = environment.__isNew__ || false;

      if (isNew) {
        saveEnvironment(connectionId, connName, envName);
        return;
      }

      addConnectionToEnvironment(environmentId, envName, connectionId, connName);
    });
    unSelectedEnvironments.forEach((environment) => {
      let envName = environment.label;
      let environmentId = environment.value || '';

      removeConnectionFromEnvironment(environmentId, envName, connectionId, connName);
    });
  };

  let colViews = [
    ['name', 'xs'],
    ['environments', 'm'],
    ['kind', 'm'],
    ['type', 's'],
    ['sub_type', 'm'],
    ['updated_at', 'l'],
    ['created_at', 'na'],
    ['status', 'xs'],
    ['Actions', 'xs'],
  ];

  const columns = [
    {
      name: 'id',
      label: 'ID',
      options: {
        display: false,
      },
    },
    {
      name: 'metadata.server_location',
      label: 'Server Location',
      options: {
        display: false,
      },
    },
    {
      name: 'metadata.server',
      label: 'Server',
      options: {
        display: false,
      },
    },
    {
      name: 'name',
      label: 'Name',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: (value, tableMeta) => {
          const server =
            getColumnValue(tableMeta.rowData, 'metadata.server', columns) ||
            getColumnValue(tableMeta.rowData, 'metadata.server_location', columns);
          return (
            <TootltipWrappedConnectionChip
              tooltip={'Server: ' + server}
              title={value}
              status={getColumnValue(tableMeta.rowData, 'status', columns)}
              onDelete={() =>
                handleDeleteConnection(
                  getColumnValue(tableMeta.rowData, 'id', columns),
                  getColumnValue(tableMeta.rowData, 'kind', columns),
                )
              }
              handlePing={(e) => {
                e.stopPropagation();
                if (tableMeta.rowData[4] === CONNECTION_KINDS.KUBERNETES) {
                  ping(tableMeta.rowData[3], tableMeta.rowData[2], tableMeta.rowData[0]);
                }
              }}
              iconSrc={`/${getColumnValue(tableMeta.rowData, 'kindLogo', columns)}`}
              style={{ maxWidth: '120px' }}
            />
          );
        },
      },
    },
    {
      name: 'environments',
      label: 'Environments',
      options: {
        sort: false,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: (value, tableMeta) => {
          const getOptions = () => {
            return environments.map((env) => ({ label: env.name, value: env.id })) || [];
          };
          let cleanedEnvs = value?.map((env) => ({ label: env.name, value: env.id })) || [];
          return (
            isEnvironmentsSuccess && (
              <Grid item xs={12} style={{ height: '5rem', width: '15rem' }}>
                <Grid item xs={12} style={{ marginTop: '2rem', cursor: 'pointer' }}>
                  <MultiSelectWrapper
                    onChange={(selected, unselected) =>
                      handleEnvironmentSelect(
                        connections[tableMeta.rowIndex].id,
                        connections[tableMeta.rowIndex].name,
                        cleanedEnvs,
                        selected,
                        unselected,
                      )
                    }
                    options={getOptions()}
                    value={cleanedEnvs}
                    placeholder={`Assigned Environments`}
                    isSelectAll={true}
                    menuPlacement={'bottom'}
                  />
                </Grid>
              </Grid>
            )
          );
        },
      },
    },
    {
      name: 'kind',
      label: 'Kind',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'type',
      label: 'Category',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'sub_type',
      label: 'Sub Category',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'updated_at',
      label: 'Updated At',
      options: {
        sort: true,
        sortThirdClickReset: true,
        display: false,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: function CustomBody(value) {
          return (
            <Tooltip
              title={
                <Moment startOf="day" format="LLL">
                  {value}
                </Moment>
              }
              placement="top"
              arrow
              interactive
            >
              <Moment format="LL">{value}</Moment>
            </Tooltip>
          );
        },
      },
    },
    {
      name: 'created_at',
      label: 'Discovered At',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: function CustomBody(value) {
          return (
            <Tooltip
              title={
                <Moment startOf="day" format="LLL">
                  {value}
                </Moment>
              }
              placement="top"
              arrow
              interactive
            >
              <Moment format="LL">{value}</Moment>
            </Tooltip>
          );
        },
      },
    },
    {
      name: 'status',
      label: 'Status',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
              icon={
                <InfoIcon
                  color={theme.palette.secondary.iconMain}
                  style={{
                    cursor: 'pointer',
                    height: 20,
                    width: 20,
                  }}
                  onClick={(e) => {
                    e.stopPropagation();
                    window.open(url, '_blank');
                  }}
                />
              }
              tooltip="Click to learn about connection and status"
            />
          );
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          const nextStatusCol = getColumnValue(tableMeta.rowData, 'nextStatus', columns);
          const originalNextStatus = nextStatusCol && nextStatusCol[value];
          let nextStatus = [];
          if (originalNextStatus !== undefined) {
            nextStatus = Object.values(originalNextStatus);
            nextStatus.push(value);
          } else {
            nextStatus.push(value);
          }
          const disabled = value === 'deleted' ? true : false;
          return (
            <>
              <FormControl className={classes.chipFormControl}>
                <Select
                  labelId="connection-status-select-label"
                  id="connection-status-select"
                  disabled={disabled}
                  value={value}
                  defaultValue={value}
                  onClick={(e) => e.stopPropagation()}
                  onChange={(e) =>
                    handleStatusChange(
                      e,
                      getColumnValue(tableMeta.rowData, 'id', columns),
                      getColumnValue(tableMeta.rowData, 'kind', columns),
                    )
                  }
                  className={classes.statusSelect}
                  disableUnderline
                  MenuProps={{
                    anchorOrigin: {
                      vertical: 'bottom',
                      horizontal: 'left',
                    },
                    transformOrigin: {
                      vertical: 'top',
                      horizontal: 'left',
                    },
                    getContentAnchorEl: null,
                    MenuListProps: { disablePadding: true },
                    PaperProps: { square: true },
                  }}
                >
                  {nextStatus &&
                    nextStatus.map((status) => (
                      <MenuItem
                        disabled={status === value ? true : false}
                        style={{ padding: '0', display: status === value ? 'none' : 'flex' }}
                        value={status}
                        key={status}
                      >
                        <Chip
                          className={classNames(classes.statusChip, classes[status])}
                          avatar={icons[status] ? icons[status]() : ''}
                          label={status}
                        />
                      </MenuItem>
                    ))}
                </Select>
              </FormControl>
            </>
          );
        },
      },
    },
    {
      name: 'Actions',
      label: 'Actions',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return (
            <TableCell>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <div className={classes.centerContent}>
              {getColumnValue(tableMeta.rowData, 'kind', columns) ===
              CONNECTION_KINDS.KUBERNETES ? (
                <IconButton
                  aria-label="more"
                  id="long-button"
                  aria-haspopup="true"
                  onClick={(e) => handleActionMenuOpen(e, tableMeta)}
                >
                  <MoreVertIcon style={iconMedium} />
                </IconButton>
              ) : (
                '-'
              )}
            </div>
          );
        },
      },
    },
    {
      name: 'nextStatus',
      label: 'nextStatus',
      options: {
        display: false,
      },
    },
    {
      name: 'kindLogo',
      label: 'kindLogo',
      options: {
        display: false,
      },
    },
  ];

  const options = useMemo(
    () => ({
      filter: false,
      viewColumns: false,
      search: false,
      responsive: 'standard',
      resizableColumns: true,
      serverSide: true,
      count,
      rowsPerPage: pageSize,
      rowsPerPageOptions: [10, 20, 30],
      fixedHeader: true,
      page,
      print: false,
      download: false,
      textLabels: {
        selectedRows: {
          text: 'connection(s) selected',
        },
      },
      customToolbarSelect: (selected) => (
        <Button
          variant="contained"
          color="primary"
          size="large"
          // @ts-ignore
          onClick={() => handleDeleteConnections(selected)}
          style={{ background: '#8F1F00', marginRight: '10px' }}
        >
          <DeleteForeverIcon style={iconMedium} />
          Delete
        </Button>
      ),
      enableNestedDataAccess: '.',
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        let order = '';
        if (tableState.activeColumn) {
          order = `${columns[tableState.activeColumn].name} desc`;
        }
        switch (action) {
          case 'changePage':
            setPage(tableState.page.toString());
            break;
          case 'changeRowsPerPage':
            setPageSize(tableState.rowsPerPage.toString());
            break;
          case 'sort':
            if (sortInfo.length == 2) {
              if (sortInfo[1] === 'ascending') {
                order = `${columns[tableState.activeColumn].name} asc`;
              } else {
                order = `${columns[tableState.activeColumn].name} desc`;
              }
            }
            if (order !== sortOrder) {
              setSortOrder(order);
            }
            break;
        }
      },
      expandableRows: true,
      expandableRowsHeader: false,
      expandableRowsOnClick: true,
      rowsExpanded: rowsExpanded,
      isRowExpandable: () => {
        return true;
      },
      onRowExpansionChange: (_, allRowsExpanded) => {
        setRowsExpanded(allRowsExpanded.slice(-1).map((item) => item.index));
        setShowMore(false);
      },
      renderExpandableRow: (rowData, tableMeta) => {
        const colSpan = rowData.length;
        const connection = connections && connections[tableMeta.rowIndex];
        return (
          <TableCell colSpan={colSpan} className={classes.innerTableWrapper}>
            <TableContainer className={classes.innerTableContainer}>
              <Table>
                <TableRow className={classes.noGutter}>
                  <TableCell style={{ padding: '20px 0', overflowX: 'hidden' }}>
                    <Grid container spacing={1} style={{ textTransform: 'lowercase' }}>
                      <Grid item xs={12} md={12} className={classes.contentContainer}>
                        <Grid container spacing={1}>
                          <Grid item xs={12} md={12} className={classes.contentContainer}>
                            <FormatConnectionMetadata
                              connection={connection}
                              meshsyncControllerState={meshsyncControllerState}
                            />
                          </Grid>
                        </Grid>
                      </Grid>
                    </Grid>
                  </TableCell>
                </TableRow>
              </Table>
            </TableContainer>
          </TableCell>
        );
      },
    }),
    [rowsExpanded, showMore, page, pageSize],
  );

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  });

  /**
   * fetch connections when the page loads
   */
  useEffect(() => {
    updateCols(columns);
    if (!loading && connectionMetadataState) {
      getConnections(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder, connectionMetadataState, isEnvironmentsSuccess]);

  useEffect(() => {
    if (isEnvironmentsError) {
      notify({
        message: `${ACTION_TYPES.FETCH_ENVIRONMENT.error_msg}: ${environmentsError}`,
        event_type: EVENT_TYPES.ERROR,
        details: environmentsError.toString(),
      });
    }
  }, [environmentsError]);

  const getConnections = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';
    dataFetch(
      `/api/integrations/connections?page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
        search,
      )}&order=${encodeURIComponent(sortOrder)}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        res?.connections.forEach((connection) => {
          (connection.nextStatus =
            connection.nextStatus === undefined &&
            connectionMetadataState[connection.kind]?.transitions),
            (connection.kindLogo =
              connection.kindLogo === undefined && connectionMetadataState[connection.kind]?.icon);
        });
        setConnections(res?.connections || []);
        setPage(res?.page || 0);
        setCount(res?.total_count || 0);
        setPageSize(res?.page_size || 0);
        setLoading(false);
      },
      handleError(ACTION_TYPES.FETCH_CONNECTIONS),
    );
  };

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const updateConnectionStatus = (connectionKind, requestBody) => {
    dataFetch(
      `/api/integrations/connections/${connectionKind}/status`,
      {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: requestBody,
      },
      () => {
        getConnections(page, pageSize, search, sortOrder);
      },
      handleError(ACTION_TYPES.UPDATE_CONNECTION),
    );
  };

  const handleStatusChange = async (e, connectionId, connectionKind) => {
    e.stopPropagation();
    let response = await modalRef.current.show({
      title: `Connection status transition`,
      subtitle: `Are you sure that you want to transition the connection status to ${e.target.value.toUpperCase()}?`,
      options: ['Confirm', 'No'],
      variant: PROMPT_VARIANTS.CONFIRMATION,
    });
    if (response === 'Confirm') {
      const requestBody = JSON.stringify({
        [connectionId]: e.target.value,
      });
      updateConnectionStatus(connectionKind, requestBody);
    }
  };

  const handleDeleteConnection = async (connectionId, connectionKind) => {
    if (connectionId) {
      let response = await modalRef.current.show({
        title: `Delete Connection`,
        subtitle: `Are you sure that you want to delete the connection?`,
        options: ['Delete', 'No'],
        variant: PROMPT_VARIANTS.DANGER,
      });
      if (response === 'Delete') {
        const requestBody = JSON.stringify({
          [connectionId]: CONNECTION_STATES.DELETED,
        });
        updateConnectionStatus(connectionKind, requestBody);
      }
    }
  };

  const handleDeleteConnections = async (selected) => {
    if (selected) {
      let response = await modalRef.current.show({
        title: `Delete Connections`,
        subtitle: `Are you sure that you want to delete the connections?`,
        options: ['Delete', 'No'],
        variant: PROMPT_VARIANTS.DANGER,
      });
      if (response === 'Delete') {
        // let bulkConnections = {}
        // selected.data.map(({ index }) => {
        //   bulkConnections = {
        //     ...bulkConnections,
        //     [connections[index].id]: CONNECTION_STATES.DELETED
        //   };
        // })
        // const requestBody = JSON.stringify(bulkConnections);
        // updateConnectionStatus(requestBody);
        selected.data.map(({ index }) => {
          const requestBody = JSON.stringify({
            [connections[index].id]: CONNECTION_STATES.DELETED,
          });
          updateConnectionStatus(connections[index].kind, requestBody);
        });
      }
    }
  };

  const handleActionMenuOpen = (event, tableMeta) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
    setSelectedRowData(tableMeta);
  };

  const handleActionMenuClose = () => {
    setAnchorEl(null);
  };

  const handleFlushMeshSync = (index) => {
    return async () => {
      handleActionMenuClose();
      let response = await meshSyncResetRef.current.show({
        title: `Flush MeshSync data for ${connections[index].metadata?.name} ?`,
        subtitle: `Are you sure to Flush MeshSync data for “${connections[index].metadata?.name}”? Fresh MeshSync data will be repopulated for this context, if MeshSync is actively running on this cluster.`,
        options: ['PROCEED', 'CANCEL'],
        variant: PROMPT_VARIANTS.WARNING,
      });
      if (response === 'PROCEED') {
        updateProgress({ showProgress: true });
        resetDatabase({
          selector: {
            clearDB: 'true',
            ReSync: 'true',
            hardReset: 'false',
          },
          k8scontextID: connections[index].metadata?.id,
        }).subscribe({
          next: (res) => {
            updateProgress({ showProgress: false });
            if (res.resetStatus === 'PROCESSING') {
              notify({ message: `Database reset successful.`, event_type: EVENT_TYPES.SUCCESS });
            }
          },
          error: handleError('Database is not reachable, try restarting server.'),
        });
      }
    };
  };

  function getOperatorStatus(index) {
    const ctxId = connections[index]?.metadata?.id;
    const operator = meshsyncControllerState?.find(
      (op) => op.contextId === ctxId && op.controller === CONTROLLERS.OPERATOR,
    );
    if (!operator) {
      return {};
    }
    return {
      operatorState: operator.status === CONTROLLER_STATES.DEPLOYED,
      operatorVersion: operator?.version,
    };
  }

  const handleOperatorSwitch = (index, checked) => {
    const contextId = connections[index].metadata?.id;
    const variables = {
      status: `${checked ? CONTROLLER_STATES.DEPLOYED : CONTROLLER_STATES.DISABLED}`,
      contextID: contextId,
    };

    updateProgress({ showProgress: true });

    changeOperatorState((response, errors) => {
      updateProgress({ showProgress: false });

      if (errors !== undefined) {
        handleError(`Unable to ${!checked ? 'Uni' : 'I'}nstall operator`);
      }
      notify({
        message: `Operator ${response.operatorStatus.toLowerCase()}`,
        event_type: EVENT_TYPES.SUCCESS,
      });

      const tempSubscription = fetchMesheryOperatorStatus({ k8scontextID: contextId }).subscribe({
        next: (res) => {
          _setOperatorState(updateCtxInfo(contextId, res));
          tempSubscription.unsubscribe();
        },
        error: (err) => console.log('error at operator scan: ' + err),
      });
    }, variables);
  };

  const updateCtxInfo = (ctxId, newInfo) => {
    if (newInfo.operator.error) {
      handleError('There is problem With operator')(newInfo.operator.error.description);
      return;
    }

    const state = _operatorStateRef.current;
    const op = state?.find((ctx) => ctx.contextID === ctxId);
    if (!op) {
      return [...state, { contextID: ctxId, operatorStatus: newInfo.operator }];
    }

    let ctx = { ...op };
    const removeCtx = state?.filter((ctx) => ctx.contextID !== ctxId);
    ctx.operatorStatus = newInfo.operator;
    return removeCtx ? [...removeCtx, ctx] : [ctx];
  };

  return (
    <>
      <NoSsr>
        <AppBar position="static" color="default" className={classes.appBar}>
          <Tabs
            value={tab}
            className={classes.tabs}
            onChange={(e, newTab) => {
              e.stopPropagation();
              setTab(newTab);
            }}
            indicatorColor="primary"
            textColor="primary"
            variant="fullWidth"
            sx={{
              height: '10%',
            }}
          >
            <Tab
              className={classes.tab}
              label={
                <div className={classes.iconText}>
                  <span style={{ marginRight: '0.3rem' }}>Connections</span>
                  <ConnectionIcon width="20" height="20" />
                  {/* <img src="/static/img/connection-light.svg" className={classes.icon} /> */}
                </div>
              }
            />
            <Tab
              className={classes.tab}
              label={
                <div className={classes.iconText}>
                  <span style={{ marginRight: '0.3rem' }}>MeshSync</span>
                  <MeshsyncIcon width="20" height="20" />
                </div>
              }
            />
          </Tabs>
        </AppBar>
        {tab === 0 && (
          <div
            className={StyleClass.toolWrapper}
            style={{ marginBottom: '5px', marginTop: '-30px' }}
          >
            <div className={classes.createButton}>
              {/* <div>
              <Button
                aria-label="Rediscover"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={onOpenCreateConnectionModal}
                style={{ marginRight: '1rem', borderRadius: '5px' }}
              >
                Connect Helm Repository
              </Button>
            </div> */}
              <MesherySettingsEnvButtons />
            </div>
            <div
              className={classes.searchAndView}
              style={{
                display: 'flex',
                borderRadius: '0.5rem 0.5rem 0 0',
              }}
            >
              <SearchBar
                onSearch={(value) => {
                  setSearch(value);
                }}
                placeholder="Search connections..."
                expanded={isSearchExpanded}
                setExpanded={setIsSearchExpanded}
              />

              <CustomColumnVisibilityControl
                id="ref"
                columns={getVisibilityColums(columns)}
                customToolsProps={{ columnVisibility, setColumnVisibility }}
              />
            </div>
          </div>
        )}
        {tab === 0 && (
          <ResponsiveDataTable
            data={connections}
            columns={columns}
            options={options}
            className={classes.muiRow}
            tableCols={tableCols}
            updateCols={updateCols}
            columnVisibility={columnVisibility}
          />
        )}
        {tab === 1 && (
          <MeshSyncTable
            classes={classes}
            updateProgress={updateProgress}
            search={search}
            selectedK8sContexts={selectedK8sContexts}
            k8sconfig={k8sconfig}
          />
        )}
        <PromptComponent ref={modalRef} />
        <Popover
          open={open}
          anchorEl={anchorEl}
          onClose={handleActionMenuClose}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'left',
          }}
        >
          <Grid style={{ margin: '10px' }}>
            <div className={classNames(classes.list, classes.listButton)}>
              <Box className={classes.listItem} sx={{ width: '100%' }}>
                <Button
                  type="submit"
                  onClick={handleFlushMeshSync(rowData.rowIndex)}
                  data-cy="btnResetDatabase"
                  className={classes.button}
                >
                  <SyncIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                  <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
                    Flush MeshSync
                  </Typography>
                </Button>
              </Box>
            </div>
            <div className={classes.list}>
              <Box className={classes.listItem} sx={{ width: '100%' }}>
                <div className={classes.listContainer}>
                  <Switch
                    defaultChecked={getOperatorStatus(rowData.rowIndex)?.operatorState}
                    onClick={(e) => handleOperatorSwitch(rowData.rowIndex, e.target.checked)}
                    name="OperatorSwitch"
                    color="primary"
                    className={classes.OperatorSwitch}
                  />
                  <Typography variant="body1">Operator</Typography>
                </div>
              </Box>
            </div>
          </Grid>
        </Popover>
        <PromptComponent ref={meshSyncResetRef} />
      </NoSsr>
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const selectedK8sContexts = state.get('selectedK8sContexts');
  const operatorState = state.get('operatorState');
  const connectionMetadataState = state.get('connectionMetadataState');
  const meshsyncControllerState = state.get('controllerState');
  const organization = state.get('organization');

  return {
    k8sconfig,
    selectedK8sContexts,
    operatorState,
    connectionMetadataState,
    meshsyncControllerState,
    organization,
  };
};

const ConnectionManagementPageWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in Connection Management', e)}
      >
        <Provider store={store}>
          <ConnectionManagementPage {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

// @ts-ignore
export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(ConnectionManagementPageWithErrorBoundary),
);
