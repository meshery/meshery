import React, { useEffect, useRef, useState } from 'react';
import {
  NoSsr,
  TableCell,
  Button,
  FormControl,
  Select,
  TableContainer,
  Table,
  Grid,
  TableRow,
  IconButton,
  Typography,
  Popover,
  AppBar,
  Tabs,
  Tab,
  MenuItem,
  Box,
  Chip,
} from '@material-ui/core';
import {
  CustomColumnVisibilityControl,
  CustomTooltip,
  SearchBar,
  UniversalFilter,
} from '@layer5/sistent';
import { withStyles } from '@material-ui/core/styles';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import { updateProgress } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
// import CustomColumnVisibilityControl from '../../utils/custom-column';
// import SearchBar from '../../utils/custom-search';
import { ResponsiveDataTable } from '@layer5/sistent';
import useStyles from '../../assets/styles/general/tool.styles';
import Modal from '../Modal';
import { iconMedium } from '../../css/icons.styles';
import PromptComponent, { PROMPT_VARIANTS } from '../PromptComponent';
import resetDatabase from '../graphql/queries/ResetDatabaseQuery';
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
  CONNECTION_STATE_TO_TRANSITION_MAP,
} from '../../utils/Enum';
import FormatConnectionMetadata from './metadata';
import useKubernetesHook from '../hooks/useKubernetesHook';
import theme from '../../themes/app';
import { TootltipWrappedConnectionChip } from './ConnectionChip';
import { DefaultTableCell, SortableTableCell } from './common';
import { getColumnValue, getVisibilityColums } from '../../utils/utils';
import HandymanIcon from '@mui/icons-material/Handyman';
import NotInterestedRoundedIcon from '@mui/icons-material/NotInterestedRounded';
import DisconnectIcon from '../../assets/icons/disconnect';
import { updateVisibleColumns } from '../../utils/responsive-column';
import { useWindowDimensions } from '../../utils/dimension';
// import UniversalFilter from '../../utils/custom-filter';
import MultiSelectWrapper from '../multi-select-wrapper';
import {
  useGetEnvironmentsQuery,
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
  useSaveEnvironmentMutation,
} from '../../rtk-query/environments';
import ErrorBoundary from '../ErrorBoundary';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '../General/error-404/index';
import { useGetConnectionsQuery, useUpdateConnectionMutation } from '@/rtk-query/connection';
import { useGetSchemaQuery } from '@/rtk-query/schema';
import { CustomTextTooltip } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@/assets/icons/InfoOutlined';
import { DeleteIcon } from '@layer5/sistent';
import { withRouter } from 'next/router';
import { UsesSistent } from '../SistentWrapper';
import { formatDate } from '../DataFormatter';
import { getFallbackImageBasedOnKind } from '@/utils/fallback';

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

  const { data: schemaResponse } = useGetSchemaQuery({
    schemaName: 'helmRepo',
  });

  const createConnection = schemaResponse ?? {};

  const handleCreateConnectionModalOpen = () => {
    setCreateConnectionModal({ open: true });
  };

  const handleCreateConnectionModalClose = () => {
    setCreateConnectionModal({ open: false });
  };

  const handleCreateConnectionSubmit = () => {};

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
  const modalRef = useRef(null);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState();
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [rowData, setSelectedRowData] = useState({});
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [_operatorState] = useState(operatorState || []);
  const [tab, setTab] = useState(0);
  const ping = useKubernetesHook();
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [statusFilter, setStatusFilter] = useState();
  const [kindFilter, setKindFilter] = useState();
  const [selectedFilters, setSelectedFilters] = useState({ status: 'All', kind: 'All' });

  const [useUpdateConnectionMutator] = useUpdateConnectionMutation();
  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();
  const [saveEnvironmentMutator] = useSaveEnvironmentMutation();

  // lock for not allowing multiple updates at the same time
  // needs to be a ref because it needs to be shared between renders
  // and useState loses reactivity when down table custom cells
  const updatingConnection = useRef(false);

  const {
    data: connectionData,
    isError: isConnectionError,
    error: connectionError,
    refetch: getConnections,
  } = useGetConnectionsQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    status: statusFilter ? JSON.stringify([statusFilter]) : '',
    kind: kindFilter ? JSON.stringify([kindFilter]) : '',
  });

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

  const modifiedConnections = connectionData?.connections.map((connection) => {
    return {
      ...connection,
      nextStatus:
        connection.nextStatus === undefined &&
        connectionMetadataState[connection.kind]?.transitions,
      kindLogo: connection.kindLogo === undefined && connectionMetadataState[connection.kind]?.icon,
    };
  });

  const connections = modifiedConnections?.filter((connection) => {
    if (selectedFilters.status === 'All' && selectedFilters.kind === 'All') {
      return true;
    }
    return (
      (statusFilter === null || connection.status === statusFilter) &&
      (kindFilter === null || connection.kind === kindFilter)
    );
  });

  const addConnectionToEnvironment = async (
    environmentId,
    environmentName,
    connectionId,
    connectionName,
  ) => {
    return addConnectionToEnvironmentMutator({ environmentId, connectionId })
      .unwrap()
      .then(() => {
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

  const removeConnectionFromEnvironment = async (
    environmentId,
    environmentName,
    connectionId,
    connectionName,
  ) => {
    return removeConnectionFromEnvMutator({ environmentId, connectionId })
      .unwrap()
      .then(() => {
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

  const saveEnvironment = async (connectionId, connectionName, environmentName) => {
    return saveEnvironmentMutator({
      body: {
        name: environmentName,
        organization_id: organization?.id,
      },
    })
      .unwrap()
      .then((resp) => {
        notify({
          message: `Environment "${resp.name}" created`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        environments = [...environments, resp];
        addConnectionToEnvironment(resp.id, resp.name, connectionId, connectionName);
      })
      .catch((err) => {
        notify({
          message: `${ACTION_TYPES.CREATE_ENVIRONMENT.error_msg}: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      });
  };

  const open = Boolean(anchorEl);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;
  const meshSyncResetRef = useRef(null);
  const { notify } = useNotification();
  const StyleClass = useStyles();
  const url = `https://docs.meshery.io/concepts/logical/connections#states-and-the-lifecycle-of-connections`;
  const envUrl = `https://docs.meshery.io/concepts/logical/environments`;

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

  const handleEnvironmentSelect = async (
    connectionId,
    connName,
    assignedEnvironments,
    selectedEnvironments,
    unSelectedEnvironments,
  ) => {
    if (updatingConnection.current) {
      return;
    }

    updatingConnection.current = true;

    try {
      let newlySelectedEnvironments = selectedEnvironments.filter((env) => {
        return !assignedEnvironments.some((assignedEnv) => assignedEnv.value === env.value);
      });

      for (let environment of newlySelectedEnvironments) {
        let envName = environment.label;
        let environmentId = environment.value || '';
        let isNew = environment.__isNew__ || false;

        if (isNew) {
          await saveEnvironment(connectionId, connName, envName);
          return;
        }

        addConnectionToEnvironment(environmentId, envName, connectionId, connName);
      }
      for (let environment of unSelectedEnvironments) {
        let envName = environment.label;
        let environmentId = environment.value || '';

        await removeConnectionFromEnvironment(environmentId, envName, connectionId, connName);
      }
    } finally {
      getConnections();
      updatingConnection.current = false;
    }
  };
  let colViews = [
    ['name', 'xs'],
    ['environments', 'm'],
    ['kind', 'm'],
    ['type', 's'],
    ['sub_type', 'na'],
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
          const name = getColumnValue(tableMeta.rowData, 'metadata.name', columns);
          const kind = getColumnValue(tableMeta.rowData, 'kind', columns);
          return (
            <>
              <TootltipWrappedConnectionChip
                tooltip={'Server: ' + server}
                title={kind === CONNECTION_KINDS.KUBERNETES ? name : value}
                status={getColumnValue(tableMeta.rowData, 'status', columns)}
                onDelete={() =>
                  handleDeleteConnection(
                    getColumnValue(tableMeta.rowData, 'id', columns),
                    getColumnValue(tableMeta.rowData, 'kind', columns),
                  )
                }
                handlePing={() => {
                  // e.stopPropagation();
                  if (getColumnValue(tableMeta.rowData, 'kind', columns) === 'kubernetes') {
                    ping(
                      getColumnValue(tableMeta.rowData, 'metadata.name', columns),
                      getColumnValue(tableMeta.rowData, 'metadata.server', columns),
                      getColumnValue(tableMeta.rowData, 'id', columns),
                    );
                  }
                }}
                iconSrc={`/${
                  getColumnValue(tableMeta.rowData, 'kindLogo', columns) ||
                  getFallbackImageBasedOnKind(kind)
                }`}
                width="12rem"
              />
              {kind == 'kubernetes' && (
                <UsesSistent>
                  <CustomTextTooltip
                    placement="top"
                    interactive={true}
                    title="Learn more about connection status and how to [troubleshoot Kubernetes connections](https://docs.meshery.io/guides/troubleshooting/meshery-operator-meshsync)"
                  >
                    <IconButton className={classes.infoIconButton} color="primary">
                      <InfoOutlinedIcon height={20} width={20} className={classes.infoIcon} />
                    </IconButton>
                  </CustomTextTooltip>
                </UsesSistent>
              )}
            </>
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
        customHeadRender: function CustomHead({ ...column }) {
          return (
            <DefaultTableCell
              columnData={column}
              icon={
                <IconButton disableRipple={true} disableFocusRipple={true}>
                  <InfoOutlinedIcon
                    fill={theme.palette.secondary.iconMain}
                    style={{
                      cursor: 'pointer',
                      height: 20,
                      width: 20,
                    }}
                    onClick={(e) => {
                      e.stopPropagation();
                    }}
                  />
                </IconButton>
              }
              tooltip={`Meshery Environments allow you to logically group related Connections and their associated Credentials. [Learn more](${envUrl})`}
            />
          );
        },
        customBodyRender: (value, tableMeta) => {
          const getOptions = () => {
            return environments.map((env) => ({ label: env.name, value: env.id }));
          };
          let cleanedEnvs = value?.map((env) => ({ label: env.name, value: env.id })) || [];
          let updatingEnvs = updatingConnection.current;
          console.log('cleanedEnvs', updatingEnvs);
          return (
            isEnvironmentsSuccess && (
              <div onClick={(e) => e.stopPropagation()}>
                <Grid item xs={12} style={{ height: '5rem', width: '15rem' }}>
                  <Grid item xs={12} style={{ marginTop: '2rem', cursor: 'pointer' }}>
                    <MultiSelectWrapper
                      updating={updatingEnvs}
                      onChange={(selected, unselected) =>
                        handleEnvironmentSelect(
                          getColumnValue(tableMeta.rowData, 'id', columns),
                          getColumnValue(tableMeta.rowData, 'name', columns),
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
                      disabled={
                        !CAN(
                          keys.ASSIGN_CONNECTIONS_TO_ENVIRONMENT.action,
                          keys.ASSIGN_CONNECTIONS_TO_ENVIRONMENT.subject,
                        )
                      }
                    />
                  </Grid>
                </Grid>
              </div>
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
          const renderValue = formatDate(value);
          return (
            <UsesSistent>
              <CustomTooltip title={renderValue} placement="top" arrow interactive>
                {renderValue}
              </CustomTooltip>
            </UsesSistent>
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
                <IconButton disableRipple={true} disableFocusRipple={true}>
                  <InfoOutlinedIcon
                    fill={theme.palette.secondary.iconMain}
                    style={{
                      cursor: 'pointer',
                      height: 20,
                      width: 20,
                    }}
                    onClick={(e) => {
                      e.stopPropagation();
                    }}
                  />
                </IconButton>
              }
              tooltip={`Every connection can be in one of the states at any given point of time. Eg: Connected, Registered, Discovered, etc. It allow users more control over whether the discovered infrastructure is to be managed or not (registered for use or not). [Learn more](${url})`}
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
          const disabled =
            value === 'deleted'
              ? true
              : !CAN(keys.CHANGE_CONNECTION_STATE.action, keys.CHANGE_CONNECTION_STATE.subject);
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
                          label={
                            status == value
                              ? status
                              : CONNECTION_STATE_TO_TRANSITION_MAP?.[status] || status
                          }
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
    {
      name: 'metadata.name',
      label: 'Name',
      options: {
        display: false,
      },
    },
  ];

  const options = {
    filter: false,
    viewColumns: false,
    search: false,
    responsive: 'standard',
    resizableColumns: true,
    serverSide: true,
    count: connectionData?.total_count,
    rowsPerPage: pageSize,
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
        onClick={() => handleDeleteConnections(selected)}
        style={{ background: theme.palette.secondary.danger, marginRight: '10px' }}
        disabled={!CAN(keys.DELETE_A_CONNECTION.action, keys.DELETE_A_CONNECTION.subject)}
      >
        <DeleteIcon fill={theme.palette.secondary.whiteIcon} style={iconMedium} />
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
  };

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

  useEffect(() => {
    updateCols(columns);
    if (isEnvironmentsError) {
      notify({
        message: `${ACTION_TYPES.FETCH_ENVIRONMENT.error_msg}: ${environmentsError}`,
        event_type: EVENT_TYPES.ERROR,
        details: environmentsError.toString(),
      });
    }

    if (isConnectionError) {
      notify({
        message: `${ACTION_TYPES.FETCH_CONNECTIONS.error_msg}: ${connectionError}`,
        event_type: EVENT_TYPES.ERROR,
        details: connectionError.toString(),
      });
    }
  }, [environmentsError, connectionError, isEnvironmentsSuccess]);

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const UpdateConnectionStatus = (connectionKind, requestBody) => {
    useUpdateConnectionMutator({
      connectionKind: connectionKind,
      connectionPayload: requestBody,
    })
      .unwrap()
      .then(() => {
        notify({
          message: `Connection status updated successfully`,
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

  const handleStatusChange = async (e, connectionId, connectionKind) => {
    e.stopPropagation();
    let response = await modalRef.current.show({
      title: `Connection Status Transition`,
      subtitle: `Are you sure that you want to transition the connection status to ${e.target.value.toUpperCase()}?`,
      options: ['Confirm', 'Cancel'],
      showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
      // variant: PROMPT_VARIANTS.CONFIRMATION,
    });
    if (response === 'Confirm') {
      const requestBody = JSON.stringify({
        [connectionId]: e.target.value,
      });
      UpdateConnectionStatus(connectionKind, requestBody);
    }
  };

  const handleDeleteConnection = async (connectionId, connectionKind) => {
    if (connectionId) {
      let response = await modalRef.current.show({
        title: `Delete Connection`,
        subtitle: `Are you sure that you want to delete the connection?`,
        options: ['Delete', 'Cancel'],
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
        variant: PROMPT_VARIANTS.DANGER,
      });
      if (response === 'Delete') {
        const requestBody = JSON.stringify({
          [connectionId]: CONNECTION_STATES.DELETED,
        });
        UpdateConnectionStatus(connectionKind, requestBody);
      }
    }
  };

  const handleDeleteConnections = async (selected) => {
    if (selected) {
      let response = await modalRef.current.show({
        title: `Delete Connections`,
        subtitle: `Are you sure that you want to delete the connections?`,
        options: ['Delete', 'Cancel'],
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
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
          UpdateConnectionStatus(connections[index].kind, requestBody);
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

  const filters = {
    status: {
      name: 'Status',
      options: [
        { label: 'Connected', value: 'connected' },
        { label: 'Registered', value: 'registered' },
        { label: 'Discovered', value: 'discovered' },
        { label: 'Ignored', value: 'ignored' },
        { label: 'Deleted', value: 'deleted' },
        { label: 'Maintenance', value: 'maintenance' },
        { label: 'Disconnected', value: 'disconnected' },
        { label: 'Not Found', value: 'not found' },
      ],
    },
    kind: {
      name: 'Kind',
      options: Object.entries(CONNECTION_KINDS).map(([key, value]) => ({ label: key, value })),
    },
  };

  const handleApplyFilter = () => {
    const statusFilter = selectedFilters.status === 'All' ? null : selectedFilters.status;
    const kindFilter = selectedFilters.kind === 'All' ? null : selectedFilters.kind;

    setKindFilter(kindFilter);
    setStatusFilter(statusFilter);
  };
  return (
    <NoSsr>
      {CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject) ? (
        <>
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
              <UsesSistent>
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
                    placeholder="Search Connections..."
                    expanded={isSearchExpanded}
                    setExpanded={setIsSearchExpanded}
                  />

                  <UniversalFilter
                    id="ref"
                    filters={filters}
                    selectedFilters={selectedFilters}
                    setSelectedFilters={setSelectedFilters}
                    handleApplyFilter={handleApplyFilter}
                  />

                  <CustomColumnVisibilityControl
                    id="ref"
                    columns={getVisibilityColums(columns)}
                    customToolsProps={{ columnVisibility, setColumnVisibility }}
                  />
                </div>
              </UsesSistent>
            </div>
          )}
          {tab === 0 && CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject) && (
            <UsesSistent>
              <ResponsiveDataTable
                data={connections}
                columns={columns}
                options={options}
                className={classes.muiRow}
                tableCols={tableCols}
                updateCols={updateCols}
                columnVisibility={columnVisibility}
              />
            </UsesSistent>
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
                    disabled={
                      !CAN(keys.FLUSH_MESHSYNC_DATA.action, keys.FLUSH_MESHSYNC_DATA.subject)
                    }
                  >
                    <SyncIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                    <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
                      Flush MeshSync
                    </Typography>
                  </Button>
                </Box>
              </div>
            </Grid>
          </Popover>
          <PromptComponent ref={meshSyncResetRef} />
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
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
        <ConnectionManagementPage {...props} />
      </ErrorBoundary>
    </NoSsr>
  );
};

// @ts-ignore
export default withStyles(styles)(
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(withRouter(ConnectionManagementPageWithErrorBoundary)),
);
