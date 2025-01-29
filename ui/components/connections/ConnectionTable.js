import React, { useEffect, useRef, useState } from 'react';
import { TableCell, TableRow, Popover } from '@mui/material';
import {
  CustomTooltip,
  CustomColumnVisibilityControl,
  SearchBar,
  UniversalFilter,
  ResponsiveDataTable,
  PROMPT_VARIANTS,
  MenuItem,
  Box,
  IconButton,
  Typography,
  Table,
  Grid,
  Button,
  FormControl,
  useTheme,
} from '@layer5/sistent';
import {
  ActionButton,
  ContentContainer,
  CreateButton,
  InnerTableContainer,
  ActionListItem,
  ConnectionStyledSelect,
} from './styles';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import MesherySettingsEnvButtons from '../MesherySettingsEnvButtons';
import { getVisibilityColums } from '../../utils/utils';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { updateProgress, useLegacySelector } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { iconMedium } from '../../css/icons.styles';
import _PromptComponent from '../PromptComponent';
import resetDatabase from '../graphql/queries/ResetDatabaseQuery';
import SyncIcon from '@mui/icons-material/Sync';

import { CONNECTION_KINDS, CONNECTION_STATES } from '../../utils/Enum';
import FormatConnectionMetadata from './metadata';
import useKubernetesHook from '../hooks/useKubernetesHook';
import { ConnectionStateChip, TootltipWrappedConnectionChip } from './ConnectionChip';
import { DefaultTableCell, SortableTableCell } from './common';
import { getColumnValue } from '../../utils/utils';
import { updateVisibleColumns } from '../../utils/responsive-column';
import { useWindowDimensions } from '../../utils/dimension';
import MultiSelectWrapper from '../multi-select-wrapper';
import {
  useAddConnectionToEnvironmentMutation,
  useGetEnvironmentsQuery,
  useRemoveConnectionFromEnvironmentMutation,
  useSaveEnvironmentMutation,
} from '../../rtk-query/environments';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useGetConnectionsQuery, useUpdateConnectionMutation } from '@/rtk-query/connection';
import { CustomTextTooltip } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@/assets/icons/InfoOutlined';
import { DeleteIcon } from '@layer5/sistent';
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

const ConnectionTable = ({ meshsyncControllerState, connectionMetadataState, selectedFilter }) => {
  const organization = useLegacySelector((state) => state.get('organization'));
  const ping = useKubernetesHook();
  const { width } = useWindowDimensions();
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState();
  const [sortOrder, setSortOrder] = useState('');
  const [rowData, setRowData] = useState(null);
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [selectedFilters, setSelectedFilters] = useState({
    status: 'All',
    kind: 'All',
  });
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState();
  const [kindFilter, setKindFilter] = useState();

  const [useUpdateConnectionMutator] = useUpdateConnectionMutation();
  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();
  const [saveEnvironmentMutator] = useSaveEnvironmentMutation();
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const modalRef = useRef(null);

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
  // lock for not allowing multiple updates at the same time
  // needs to be a ref because it needs to be shared between renders
  // and useState loses reactivity when down table custom cells
  const updatingConnection = useRef(false);
  const { notify } = useNotification();

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
    kind: selectedFilter
      ? JSON.stringify([selectedFilter])
      : kindFilter
        ? JSON.stringify([kindFilter])
        : '',
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
  const url = `https://docs.meshery.io/concepts/logical/connections#states-and-the-lifecycle-of-connections`;
  const envUrl = `https://docs.meshery.io/concepts/logical/environments`;

  let colViews = [
    ['name', 'xs'],
    ['environments', 'm'],
    ['kind', 'm'],
    ['type', 's'],
    ['sub_type', 'na'],
    ['created_at', 'na'],
    ['status', 'xs'],
    ['Actions', 'xs'],
    ['ConnectionID', 'xs'],
  ];

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

  const handleDeleteConnection = async (connectionId, connectionKind) => {
    if (connectionId) {
      let response = await modalRef.current.show({
        title: `Delete Connection`,
        subtitle: `Are you sure that you want to delete the connection?`,
        primaryOption: 'DELETE',
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
        variant: PROMPT_VARIANTS.DANGER,
      });
      if (response === 'DELETE') {
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
        primaryOption: 'DELETE',
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
        variant: PROMPT_VARIANTS.DANGER,
      });
      if (response === 'DELETE') {
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

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleActionMenuClose = () => {
    setAnchorEl(null);
    setRowData(null);
  };

  const handleFlushMeshSync = () => {
    return async () => {
      handleActionMenuClose();
      const connection = connections[rowData.rowIndex];

      let response = await modalRef.current.show({
        title: `Flush MeshSync data for ${connection.metadata?.name} ?`,
        subtitle: `Are you sure to Flush MeshSync data for “${connection.metadata?.name}”? Fresh MeshSync data will be repopulated for this context, if MeshSync is actively running on this cluster.`,
        primaryOption: 'PROCEED',
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
          k8scontextID: connection.metadata?.id,
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

  const UpdateConnectionStatus = (connectionKind, requestBody) => {
    useUpdateConnectionMutator({
      connectionKind: connectionKind,
      connectionPayload: requestBody,
    })
      .unwrap()
      .then(() => {
        notify({
          message: `Connection status updated`,
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
      primaryOption: 'Confirm',
      showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
      variant: PROMPT_VARIANTS.WARNING,
    });
    if (response === 'Confirm') {
      const requestBody = JSON.stringify({
        [connectionId]: e.target.value,
      });
      UpdateConnectionStatus(connectionKind, requestBody);
    }
  };

  const handleActionMenuOpen = (event, tableMeta) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
    setRowData(tableMeta);
  };
  const theme = useTheme();

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
                    <IconButton color="primary">
                      <InfoOutlinedIcon height={20} width={20} />
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
      name: 'ConnectionID',
      label: 'Connection ID',
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
          const connectionId = getColumnValue(tableMeta.rowData, 'id', columns);
          return <span>{connectionId}</span>;
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
            <UsesSistent>
              <FormControl>
                <ConnectionStyledSelect
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
                        style={{
                          padding: 0,
                          display: status === value ? 'none' : 'flex',
                          justifyContent: 'center',
                        }}
                        value={status}
                        key={status}
                      >
                        <ConnectionStateChip status={status} />
                      </MenuItem>
                    ))}
                </ConnectionStyledSelect>
              </FormControl>
            </UsesSistent>
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
            <Box display={'flex'} justifyContent={'center'}>
              {getColumnValue(tableMeta.rowData, 'kind', columns) ===
              CONNECTION_KINDS.KUBERNETES ? (
                <UsesSistent>
                  <IconButton
                    aria-label="more"
                    id="long-button"
                    aria-haspopup="true"
                    onClick={(e) => handleActionMenuOpen(e, tableMeta)}
                  >
                    <MoreVertIcon style={iconMedium} />
                  </IconButton>
                </UsesSistent>
              ) : (
                '-'
              )}
            </Box>
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
    sortOrder: {
      name: 'name',
      direction: 'asc',
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
        <DeleteIcon style={iconMedium} />
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
        <UsesSistent>
          <TableCell colSpan={colSpan}>
            <InnerTableContainer>
              <Table>
                <TableRow style={{ padding: 0 }}>
                  <TableCell style={{ padding: '20px 0', overflowX: 'hidden' }}>
                    <Grid container spacing={1} style={{ textTransform: 'lowercase' }}>
                      <ContentContainer item xs={12} md={12}>
                        <Grid container spacing={1}>
                          <ContentContainer item xs={12} md={12}>
                            <FormatConnectionMetadata
                              connection={connection}
                              meshsyncControllerState={meshsyncControllerState}
                            />
                          </ContentContainer>
                        </Grid>
                      </ContentContainer>
                    </Grid>
                  </TableCell>
                </TableRow>
              </Table>
            </InnerTableContainer>
          </TableCell>
        </UsesSistent>
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
  return (
    <UsesSistent>
      <ToolWrapper style={{ marginBottom: '5px', marginTop: '-30px' }}>
        <CreateButton>
          <MesherySettingsEnvButtons />
        </CreateButton>
        <div
          style={{
            display: 'flex',
            borderRadius: '0.5rem 0.5rem 0 0',
            width: '100%',
            justifyContent: 'flex-end',
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
            style={{ zIndex: 1300 }}
            id="ref"
            columns={getVisibilityColums(columns)}
            customToolsProps={{ columnVisibility, setColumnVisibility }}
          />
        </div>
      </ToolWrapper>

      <ResponsiveDataTable
        data={connections}
        columns={columns}
        options={options}
        tableCols={tableCols}
        updateCols={updateCols}
        columnVisibility={columnVisibility}
      />
      <_PromptComponent ref={modalRef} />

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
          <ActionListItem>
            <ActionButton
              type="submit"
              onClick={handleFlushMeshSync()}
              data-cy="btnResetDatabase"
              disabled={!CAN(keys.FLUSH_MESHSYNC_DATA.action, keys.FLUSH_MESHSYNC_DATA.subject)}
            >
              <SyncIcon {...iconMedium} />
              <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
                Flush MeshSync
              </Typography>
            </ActionButton>
          </ActionListItem>
        </Grid>
      </Popover>
    </UsesSistent>
  );
};

export default ConnectionTable;
