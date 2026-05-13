import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/router';
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
  Grid2,
  Button,
  FormControl,
  useTheme,
  TableCell,
  TableRow,
  Popover,
} from '@sistent/sistent';
import {
  ContentContainer,
  CreateButton,
  InnerTableContainer,
  ActionListItem,
  ConnectionStyledSelect,
} from './styles';
import { FormatId } from '../data-formatter';
import LoadingScreen from '../shared/LoadingState/LoadingComponent';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import MesherySettingsEnvButtons from '../MesherySettingsEnvButtons';
import { getVisibilityColums } from '../../utils/utils';
import { MoreVert as MoreVertIcon, Sync as SyncIcon } from '@/assets/icons';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { iconMedium } from '../../css/icons.styles';
import _PromptComponent from '../PromptComponent';
import resetDatabase from '@/graphql/queries/ResetDatabaseQuery';

import { CONNECTION_KINDS, CONNECTION_STATES, MESHSYNC_DEPLOYMENT_TYPE } from '../../utils/Enum';
import FormatConnectionMetadata from './metadata';
import useKubernetesHook from '@/utils/hooks/useKubernetesHook';
import { ConnectionStateChip, TooltipWrappedConnectionChip } from './ConnectionChip';
import { DefaultTableCell, SortableTableCell } from './common';
import { getColumnValue } from '../../utils/utils';
import { getResponsiveColumnVisibility } from '../../utils/responsive-column';
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
import { useGetConnectionsQuery, useUpdateConnectionByIdMutation } from '@/rtk-query/connection';
import { CustomTextTooltip } from '../meshery-mesh-interface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@/assets/icons/InfoOutlined';
import { DeleteIcon } from '@sistent/sistent';

import { formatDate } from '../data-formatter';
import { getFallbackImageBasedOnKind, normalizeStaticImagePath } from '@/utils/fallback';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';

type ConnectionTableProps = {
  selectedFilter?: string;
  selectedConnectionId?: string;
  updateUrlWithConnectionId?: (connectionId: string) => void;
};

type EnvironmentOption = {
  label: string;
  value: string;
  __isNew__?: boolean;
};

type ConnectionRow = {
  id: string;
  kind: string;
  status: string;
  name?: string;
  kindLogo?: string;
  nextStatus?: string[];
  metadata?: Record<string, any>;
  environments?: Array<{ id: string; name: string }>;
  [key: string]: any;
};

type SelectedFilters = {
  status: string;
  kind: string;
};

type SelectedRows = {
  data?: Array<{ index: number }>;
};

type RowData = {
  rowIndex: number;
};

type ExpansionFlags = {
  isHandlingExpansion: boolean;
  isInitialLoad: boolean;
  isUrlExpansion: boolean;
  lastProcessedId: string | null;
};

const getErrorMessage = (error: any, fallback = 'Unknown error') => {
  if (!error) {
    return fallback;
  }

  if (typeof error === 'string') {
    return error;
  }

  if (typeof error === 'object') {
    if ('error' in error && typeof error.error === 'string') {
      return error.error;
    }

    if ('message' in error && typeof error.message === 'string') {
      return error.message;
    }

    if ('data' in error && typeof error.data === 'string') {
      return error.data;
    }
  }

  return fallback;
};

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

const kubernetesConnectionTransitions = {
  connected: {
    disconnected:
      'Are you sure you want to transition from CONNECTED to DISCONNECTED? This will perform planned maintenance by removing the operator but keeping the cluster registered.',
    ignored:
      'Are you sure you want to transition from CONNECTED to IGNORED? This will mark the connection as ignored due to unplanned maintenance, without deleting the registration.',
    deleted:
      'Are you sure you want to transition from CONNECTED to DELETED? This will undeploy the operator and unregister the cluster completely.',
    'not found':
      'Are you sure you want to transition from CONNECTED to NOT FOUND? Meshery could not connect to the cluster or it is currently unavailable. You can either delete the connection or try re-registering.',
  },
  disconnected: {
    connected:
      'Are you sure you want to transition from DISCONNECTED to CONNECTED? This will reconnect the cluster and redeploy the operator after maintenance.',
    deleted:
      'Are you sure you want to transition from DISCONNECTED to DELETED? This will remove the cluster completely by undeploying the operator and unregistering.',
  },
  ignored: {
    deleted:
      'Are you sure you want to transition from IGNORED to DELETED? This will completely remove the ignored cluster by undeploying the operator and unregistering.',
    registered:
      'Are you sure you want to transition from IGNORED to REGISTER? This will reinitiate the registration process for the ignored connection and attempt to connect it again.',
  },
  'not found': {
    discovered:
      'Are you sure you want to transition from NOT FOUND to DISCOVERED? You are trying to re-register the cluster. Meshery will attempt to reconnect to the cluster.',
    deleted:
      'Are you sure you want to transition from NOT FOUND to DELETED? This will remove the unreachable connection completely by unregistering it.',
  },
};

// <connection-kind>: TransitionMap (status:{status:description})
const CONNECTION_STATE_TRANSITIONS = {
  kubernetes: kubernetesConnectionTransitions,
};

const getStatusTransition = (
  connectionKind: string,
  connectionState: string,
  transitionState: string,
) => {
  // This is for one connection kind that is kubernetes, and adding other connection kinds
  // here will make it more complex.
  // This issue can be resolved if we add the transition messages in the connection schemas
  // and use the same schema to get the transition messages.
  // Github issue: https://github.com/meshery/schemas/issues/303

  switch (connectionKind) {
    case 'kubernetes':
      return kubernetesConnectionTransitions[connectionState][transitionState];
    default:
      return `Are you sure you want to transition from ${connectionState} to ${transitionState}?`;
  }
};

const ConnectionTable = ({
  selectedFilter,
  selectedConnectionId,
  updateUrlWithConnectionId,
}: ConnectionTableProps) => {
  const router = useRouter();
  const {
    organization,
    connectionMetadataState,
    controllerState: meshsyncControllerState,
  } = useSelector(
    (state: {
      ui: {
        organization?: { id?: string };
        connectionMetadataState: Record<string, { transitions?: string[]; icon?: string }>;
        controllerState: unknown;
      };
    }) => state.ui,
  );
  const ping = useKubernetesHook();
  const { width } = useWindowDimensions();
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [sortOrder, setSortOrder] = useState('created_at desc');
  const [rowData, setRowData] = useState<RowData | null>(null);
  const [rowsExpanded, setRowsExpanded] = useState<number[]>([]);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [selectedFilters, setSelectedFilters] = useState<SelectedFilters>({
    status: 'All',
    kind: 'All',
  });
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<string | null>(null);
  const [kindFilter, setKindFilter] = useState<string | null>(null);
  const [updateConnectionByIdMutator] = useUpdateConnectionByIdMutation();
  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();
  const [saveEnvironmentMutator] = useSaveEnvironmentMutation();
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const [deploymentModeAnchorEl, setDeploymentModeAnchorEl] = useState<HTMLElement | null>(null);
  const open = Boolean(anchorEl);
  const deploymentModeOpen = Boolean(deploymentModeAnchorEl);
  const modalRef = useRef<{ show: (options: unknown) => Promise<string | null> } | null>(null);

  useEffect(() => {
    if (typeof router.query.searchText === 'string') {
      setSearch(router.query.searchText);
    }
  }, [router.query.searchText]);

  const filters = useMemo(
    () => ({
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
    }),
    [],
  );

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
    isLoading: isConnectionLoading,
  } = useGetConnectionsQuery(
    {
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
    },
    undefined,
  );
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

  const environmentOptions = useMemo(
    () =>
      (environmentsResponse?.environments || []).map((env) => ({
        label: env.name,
        value: env.id,
      })),
    [environmentsResponse?.environments],
  );

  useEffect(() => {
    if (isEnvironmentsError) {
      notify({
        message: `${ACTION_TYPES.FETCH_ENVIRONMENT.error_msg}: ${getErrorMessage(environmentsError)}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }

    if (isConnectionError) {
      notify({
        message: `${ACTION_TYPES.FETCH_CONNECTIONS.error_msg}: ${getErrorMessage(connectionError)}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [connectionError, environmentsError, isConnectionError, isEnvironmentsError, notify]);

  const enhancedConnections = useMemo(() => {
    if (!connectionData?.connections) return [];

    return connectionData.connections
      .filter((conn) => conn.name && conn.kind && conn.status)
      .map((connection) => ({
        ...connection,
        nextStatus: connection.nextStatus || connectionMetadataState[connection.kind]?.transitions,
        kindLogo: connection.kindLogo || connectionMetadataState[connection.kind]?.icon,
      }));
  }, [connectionData?.connections, connectionMetadataState]) as ConnectionRow[];

  const filteredConnections = useMemo(
    () =>
      enhancedConnections.filter(({ status, kind }) => {
        const statusMatch = selectedFilters.status === 'All' || status === selectedFilters.status;
        const kindMatch = selectedFilters.kind === 'All' || kind === selectedFilters.kind;
        return statusMatch && kindMatch;
      }),
    [enhancedConnections, selectedFilters],
  );

  const url = `https://docs.meshery.io/concepts/logical/connections#states-and-the-lifecycle-of-connections`;
  const envUrl = `https://docs.meshery.io/concepts/logical/environments`;

  const colViews = useMemo(
    () => [
      ['name', 'xs'],
      ['environments', 'm'],
      ['kind', 'm'],
      ['type', 's'],
      ['sub_type', 'na'],
      ['created_at', 'na'],
      ['status', 'xs'],
      ['Actions', 'xs'],
      ['ConnectionID', 'na'],
    ],
    [],
  );
  const addConnectionToEnvironment = useCallback(
    async (
      environmentId: string,
      environmentName: string,
      connectionId: string,
      connectionName: string,
    ) => {
      try {
        await addConnectionToEnvironmentMutator({ environmentId, connectionId }).unwrap();
        notify({
          message: `Connection: ${connectionName} assigned to environment: ${environmentName}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [addConnectionToEnvironmentMutator, notify],
  );

  const removeConnectionFromEnvironment = useCallback(
    async (
      environmentId: string,
      environmentName: string,
      connectionId: string,
      connectionName: string,
    ) => {
      try {
        await removeConnectionFromEnvMutator({ environmentId, connectionId }).unwrap();
        notify({
          message: `Connection: ${connectionName} removed from environment: ${environmentName}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [notify, removeConnectionFromEnvMutator],
  );

  const saveEnvironment = useCallback(
    async (connectionId: string, connectionName: string, environmentName: string) => {
      try {
        const response = await saveEnvironmentMutator({
          body: {
            name: environmentName,
            organization_id: organization?.id,
          },
        }).unwrap();

        notify({
          message: `Environment "${response.name}" created`,
          event_type: EVENT_TYPES.SUCCESS,
        });

        await addConnectionToEnvironment(response.id, response.name, connectionId, connectionName);
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.CREATE_ENVIRONMENT.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [addConnectionToEnvironment, notify, organization?.id, saveEnvironmentMutator],
  );

  const updateConnectionStatus = useCallback(
    async (connectionId: string, newStatus: string) => {
      try {
        await updateConnectionByIdMutator({
          connectionId,
          body: { status: newStatus },
        }).unwrap();

        notify({
          message: `Connection status updated`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [notify, updateConnectionByIdMutator],
  );

  const handleDeleteConnection = useCallback(
    async (connectionId: string) => {
      if (!connectionId || !modalRef.current) {
        return;
      }

      const response = await modalRef.current.show({
        title: `Delete Connection`,
        subtitle: `Are you sure that you want to delete the connection?`,
        primaryOption: 'DELETE',
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
        variant: PROMPT_VARIANTS.DANGER,
      });

      if (response === 'DELETE') {
        await updateConnectionStatus(connectionId, CONNECTION_STATES.DELETED);
      }
    },
    [updateConnectionStatus],
  );

  const handleDeleteConnections = useCallback(
    async (selected: SelectedRows) => {
      if (!selected?.data?.length || !modalRef.current) {
        return;
      }

      // Capture the connection IDs up front. The user has to acknowledge the
      // confirmation modal before delete fires, and `filteredConnections` can
      // be invalidated/reordered by an in-flight refetch in that window — using
      // the index after-the-fact dereferenced stale rows and silently no-op'd
      // (no PUT, no notification), which surfaced as a hung e2e snackbar wait.
      const ids = selected.data
        .map(({ index }) => filteredConnections[index]?.id)
        .filter(Boolean) as string[];

      if (ids.length === 0) {
        return;
      }

      const response = await modalRef.current.show({
        title: `Delete Connections`,
        subtitle: `Are you sure that you want to delete the connections?`,
        primaryOption: 'DELETE',
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
        variant: PROMPT_VARIANTS.DANGER,
      });

      if (response === 'DELETE') {
        await Promise.all(
          ids.map((connectionId) =>
            updateConnectionStatus(connectionId, CONNECTION_STATES.DELETED),
          ),
        );
      }
    },
    [filteredConnections, updateConnectionStatus],
  );

  const handleError = useCallback(
    (action: { error_msg?: string } | string) => (error: unknown) => {
      updateProgress({ showProgress: false });

      const message =
        typeof action === 'string'
          ? action
          : `${action.error_msg}: ${getErrorMessage(error, 'Request failed')}`;

      notify({
        message,
        event_type: EVENT_TYPES.ERROR,
        details: String(error),
      });
    },
    [notify],
  );

  const handleActionMenuClose = useCallback(() => {
    setAnchorEl(null);
    setRowData(null);
  }, []);

  const handleDeploymentModeMenuClose = useCallback(() => {
    setDeploymentModeAnchorEl(null);
  }, []);

  const getConnectionAtRowIndex = useCallback(
    (rowIndex?: number | null) => {
      if (rowIndex == null) {
        return null;
      }

      return filteredConnections[rowIndex] ?? null;
    },
    [filteredConnections],
  );

  const handleDeploymentModeChange = useCallback(
    async (newMode: string) => {
      const connection = getConnectionAtRowIndex(rowData?.rowIndex);

      if (!connection) {
        handleDeploymentModeMenuClose();
        handleActionMenuClose();
        return;
      }

      try {
        await updateConnectionByIdMutator({
          connectionId: connection.id,
          body: {
            ...connection,
            metadata: {
              ...connection.metadata,
              meshsync_deployment_mode: newMode,
            },
          },
        }).unwrap();

        notify({
          message: `Deployment mode changed to ${newMode}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `Failed to change deployment mode: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
        });
      }

      handleDeploymentModeMenuClose();
      handleActionMenuClose();
    },
    [
      getConnectionAtRowIndex,
      handleActionMenuClose,
      handleDeploymentModeMenuClose,
      notify,
      rowData?.rowIndex,
      updateConnectionByIdMutator,
    ],
  );

  const handleFlushMeshSync = useCallback(
    () => async () => {
      handleActionMenuClose();

      const connection = getConnectionAtRowIndex(rowData?.rowIndex);
      const connectionName = connection?.metadata?.name;

      if (!connection || !modalRef.current) {
        return;
      }

      const response = await modalRef.current.show({
        title: `Flush MeshSync data for ${connectionName} ?`,
        subtitle: `Are you sure to Flush MeshSync data for “${connectionName}”? Fresh MeshSync data will be repopulated for this context, if MeshSync is actively running on this cluster.`,
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
          k8scontextID: connection.metadata?.id || '',
        }).subscribe({
          next: (result) => {
            updateProgress({ showProgress: false });
            if (result.resetStatus === 'PROCESSING') {
              notify({ message: `Database reset successful.`, event_type: EVENT_TYPES.SUCCESS });
            }
          },
          error: handleError('Database is not reachable, try restarting server.'),
        });
      }
    },
    [getConnectionAtRowIndex, handleActionMenuClose, handleError, notify, rowData?.rowIndex],
  );

  const handleEnvironmentSelect = useCallback(
    async (
      connectionId: string,
      connName: string,
      assignedEnvironments: EnvironmentOption[],
      selectedEnvironments: EnvironmentOption[],
      unSelectedEnvironments: EnvironmentOption[],
    ) => {
      if (updatingConnection.current) {
        return;
      }

      updatingConnection.current = true;

      try {
        const newlySelectedEnvironments = selectedEnvironments.filter((environment) => {
          return !assignedEnvironments.some(
            (assignedEnvironment) => assignedEnvironment.value === environment.value,
          );
        });

        const selectedExistingEnvironments = newlySelectedEnvironments.filter(
          (environment) => !environment.__isNew__,
        );
        const selectedNewEnvironments = newlySelectedEnvironments.filter(
          (environment) => environment.__isNew__,
        );

        await Promise.all([
          ...selectedExistingEnvironments.map((environment) =>
            addConnectionToEnvironment(
              environment.value || '',
              environment.label,
              connectionId,
              connName,
            ),
          ),
          ...selectedNewEnvironments.map((environment) =>
            saveEnvironment(connectionId, connName, environment.label),
          ),
          ...unSelectedEnvironments.map((environment) =>
            removeConnectionFromEnvironment(
              environment.value || '',
              environment.label,
              connectionId,
              connName,
            ),
          ),
        ]);
      } finally {
        getConnections();
        updatingConnection.current = false;
      }
    },
    [addConnectionToEnvironment, getConnections, removeConnectionFromEnvironment, saveEnvironment],
  );

  const handleStatusChange = useCallback(
    async (event, connectionId: string, connectionKind: string, connectionStatus: string) => {
      event.stopPropagation();

      if (!modalRef.current) {
        return;
      }

      const status = event.target.value;
      const subtitle = getStatusTransition(connectionKind, connectionStatus, status.toLowerCase());
      const response = await modalRef.current.show({
        title: `Transition connection to ${status.toUpperCase()}?`,
        subtitle,
        primaryOption: 'Confirm',
        showInfoIcon: `Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.`,
        variant: PROMPT_VARIANTS.WARNING,
      });

      if (response === 'Confirm') {
        await updateConnectionStatus(connectionId, status);
      }
    },
    [updateConnectionStatus],
  );

  const handleActionMenuOpen = useCallback((event, tableMeta: RowData) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
    setRowData(tableMeta);
  }, []);
  const theme = useTheme();

  // Consolidate multiple useRef hooks into a single object
  const expansionFlags = useRef<ExpansionFlags>({
    isHandlingExpansion: false,
    isInitialLoad: true,
    isUrlExpansion: false,
    lastProcessedId: null,
  });

  // Update rowsExpanded when a specific connection ID is selected
  useEffect(() => {
    if (!selectedConnectionId || expansionFlags.current.isHandlingExpansion) return;
    if (expansionFlags.current.lastProcessedId === selectedConnectionId) return;

    if (filteredConnections && filteredConnections.length > 0) {
      expansionFlags.current.isUrlExpansion = true;
      expansionFlags.current.lastProcessedId = selectedConnectionId;

      const index = filteredConnections?.findIndex((conn) => conn.id === selectedConnectionId);
      if (index !== -1) {
        setRowsExpanded([index]);
      } else {
        updateUrlWithConnectionId?.('');
      }

      expansionFlags.current.isUrlExpansion = false;
      expansionFlags.current.isInitialLoad = false;
    }
  }, [filteredConnections, selectedConnectionId, updateUrlWithConnectionId]);

  const columns = useMemo(() => {
    const nextColumns = [
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
                icon={null}
                tooltip=""
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const server =
              getColumnValue(tableMeta.rowData, 'metadata.server', nextColumns) ||
              getColumnValue(tableMeta.rowData, 'metadata.server_location', nextColumns);
            const name = getColumnValue(tableMeta.rowData, 'metadata.name', nextColumns);
            const kind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);
            const iconSrc = normalizeStaticImagePath(
              getColumnValue(tableMeta.rowData, 'kindLogo', nextColumns) ||
                getFallbackImageBasedOnKind(kind),
            );

            return (
              <>
                <TooltipWrappedConnectionChip
                  tooltip={server ? `Server: ${server}` : ''}
                  title={kind === CONNECTION_KINDS.KUBERNETES ? name : value}
                  status={getColumnValue(tableMeta.rowData, 'status', nextColumns)}
                  onDelete={() =>
                    handleDeleteConnection(getColumnValue(tableMeta.rowData, 'id', nextColumns))
                  }
                  handlePing={() => {
                    if (getColumnValue(tableMeta.rowData, 'kind', nextColumns) === 'kubernetes') {
                      ping(
                        getColumnValue(tableMeta.rowData, 'metadata.name', nextColumns),
                        getColumnValue(tableMeta.rowData, 'metadata.server', nextColumns),
                        getColumnValue(tableMeta.rowData, 'id', nextColumns),
                      );
                    }
                  }}
                  iconSrc={iconSrc}
                  width="12rem"
                />
                {kind === 'kubernetes' && (
                  <CustomTextTooltip
                    placement="top"
                    title="Learn more about connection status and how to [troubleshoot Kubernetes connections](https://docs.meshery.io/guides/troubleshooting/meshery-operator-meshsync)"
                  >
                    <div style={{ display: 'inline-block' }}>
                      <IconButton
                        color="default"
                        onClick={(event) => {
                          event.stopPropagation();
                          event.preventDefault();
                        }}
                      >
                        <InfoOutlinedIcon height={20} width={20} />
                      </IconButton>
                    </div>
                  </CustomTextTooltip>
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
                  <IconButton
                    disableRipple={true}
                    disableFocusRipple={true}
                    onClick={(event) => {
                      event.stopPropagation();
                    }}
                  >
                    <InfoOutlinedIcon
                      style={{
                        cursor: 'pointer',
                        height: 20,
                        width: 20,
                      }}
                    />
                  </IconButton>
                }
                tooltip={`Meshery Environments allow you to logically group related Connections and their associated Credentials. [Learn more](${envUrl})`}
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const cleanedEnvs =
              value?.map((environment) => ({
                label: environment.name,
                value: environment.id,
              })) || [];

            return (
              isEnvironmentsSuccess && (
                <div onClick={(event) => event.stopPropagation()}>
                  <Grid2 size={{ xs: 12 }} style={{ height: '5rem', width: '15rem' }}>
                    <Grid2 size={{ xs: 12 }} style={{ marginTop: '2rem', cursor: 'pointer' }}>
                      <MultiSelectWrapper
                        updating={updatingConnection.current}
                        onChange={(selected, unselected) =>
                          handleEnvironmentSelect(
                            getColumnValue(tableMeta.rowData, 'id', nextColumns),
                            getColumnValue(tableMeta.rowData, 'name', nextColumns),
                            cleanedEnvs,
                            selected,
                            unselected,
                          )
                        }
                        options={environmentOptions}
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
                    </Grid2>
                  </Grid2>
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
                icon={null}
                tooltip=""
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
                icon={null}
                tooltip=""
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
                icon={null}
                tooltip=""
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
                icon={null}
                tooltip=""
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
                icon={null}
                tooltip=""
              />
            );
          },
          customBodyRender: function CustomBody(value) {
            const renderValue = formatDate(value);
            return (
              <CustomTooltip title={renderValue} placement="top" arrow interactive>
                <span>{renderValue}</span>
              </CustomTooltip>
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
                icon={null}
                tooltip=""
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const connectionId = getColumnValue(tableMeta.rowData, 'id', nextColumns);
            return <FormatId id={connectionId} />;
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
                  <IconButton
                    disableRipple={true}
                    disableFocusRipple={true}
                    onClick={(event) => {
                      event.stopPropagation();
                    }}
                  >
                    <InfoOutlinedIcon
                      style={{
                        cursor: 'pointer',
                        height: 20,
                        width: 20,
                      }}
                    />
                  </IconButton>
                }
                tooltip={`Every connection can be in one of the states at any given point of time. Eg: Connected, Registered, Discovered, etc. It allow users more control over whether the discovered infrastructure is to be managed or not (registered for use or not). [Learn more](${url})`}
              />
            );
          },
          customBodyRender: function CustomBody(value, tableMeta) {
            const currentStatus = value;
            const kind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);

            const nextStatus = Object.keys(
              CONNECTION_STATE_TRANSITIONS?.[kind]?.[currentStatus] ?? {},
            );
            nextStatus.push(currentStatus);

            const disabled =
              value === 'deleted'
                ? true
                : !CAN(keys.CHANGE_CONNECTION_STATE.action, keys.CHANGE_CONNECTION_STATE.subject);

            return (
              <FormControl>
                <ConnectionStyledSelect
                  labelId="connection-status-select-label"
                  id="connection-status-select"
                  disabled={disabled}
                  value={value}
                  defaultValue={value}
                  onClick={(event) => event.stopPropagation()}
                  onChange={(event) =>
                    handleStatusChange(
                      event,
                      getColumnValue(tableMeta.rowData, 'id', nextColumns),
                      getColumnValue(tableMeta.rowData, 'kind', nextColumns),
                      getColumnValue(tableMeta.rowData, 'status', nextColumns),
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
                  {nextStatus.length === 1 && (
                    <MenuItem disabled>No transitions Available</MenuItem>
                  )}
                  {nextStatus.map((status) => (
                    <MenuItem
                      disabled={status === value}
                      style={{
                        padding: 0,
                        display: status === value ? 'none' : 'flex',
                        justifyContent: 'center',
                      }}
                      value={status}
                      key={status}
                    >
                      <ConnectionStateChip status={status} actionable={status !== value} />
                    </MenuItem>
                  ))}
                </ConnectionStyledSelect>
              </FormControl>
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
                {getColumnValue(tableMeta.rowData, 'kind', nextColumns) ===
                CONNECTION_KINDS.KUBERNETES ? (
                  <IconButton
                    aria-label="more"
                    id="long-button"
                    aria-haspopup="true"
                    onClick={(event) => handleActionMenuOpen(event, tableMeta)}
                  >
                    <MoreVertIcon style={iconMedium} />
                  </IconButton>
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

    return nextColumns;
  }, [
    envUrl,
    environmentOptions,
    handleActionMenuOpen,
    handleDeleteConnection,
    handleEnvironmentSelect,
    handleStatusChange,
    isEnvironmentsSuccess,
    ping,
    url,
  ]);
  const columnNames = useMemo(() => columns.map((column) => column.name), [columns]);

  const options = useMemo(
    () => ({
      filter: false,
      viewColumns: false,
      search: false,
      responsive: 'standard',
      resizableColumns: true,
      serverSide: true,
      count: connectionData?.totalCount,
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
        name: sortOrder.split(' ')[0],
        direction: sortOrder.split(' ')[1],
      },
      customToolbarSelect: (selected) => (
        <Button
          color="error"
          variant="contained"
          size="large"
          onClick={() => handleDeleteConnections(selected)}
          sx={{ backgroundColor: `${theme.palette.error.dark} !important`, marginRight: '10px' }}
          disabled={!CAN(keys.DELETE_A_CONNECTION.action, keys.DELETE_A_CONNECTION.subject)}
          data-testid="Button-delete-connections"
        >
          <DeleteIcon style={iconMedium} fill={theme.palette.common.white} />
          Delete
        </Button>
      ),
      enableNestedDataAccess: '.',
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        let order = '';

        if (typeof tableState.activeColumn === 'number') {
          order = `${columns[tableState.activeColumn].name} desc`;
        }

        switch (action) {
          case 'changePage':
            setPage(tableState.page);
            break;
          case 'changeRowsPerPage':
            setPageSize(tableState.rowsPerPage);
            break;
          case 'sort':
            if (sortInfo.length === 2 && typeof tableState.activeColumn === 'number') {
              order =
                sortInfo[1] === 'ascending'
                  ? `${columns[tableState.activeColumn].name} asc`
                  : `${columns[tableState.activeColumn].name} desc`;
            }

            if (order && order !== sortOrder) {
              setSortOrder(order);
            }
            break;
        }
      },
      expandableRows: true,
      expandableRowsHeader: false,
      expandableRowsOnClick: true,
      rowsExpanded,
      isRowExpandable: () => true,
      onRowExpansionChange: (_, allRowsExpanded) => {
        if (expansionFlags.current.isUrlExpansion) {
          return;
        }

        expansionFlags.current.isHandlingExpansion = true;
        const expandedRows = allRowsExpanded.slice(-1);
        setRowsExpanded(expandedRows.map((item) => item.index));

        if (expandedRows.length > 0) {
          const index = expandedRows[0].index;
          const connection = filteredConnections[index];

          if (
            connection &&
            updateUrlWithConnectionId &&
            (!expansionFlags.current.isInitialLoad || connection.id !== selectedConnectionId)
          ) {
            updateUrlWithConnectionId(connection.id);
          }
        } else if (updateUrlWithConnectionId && !expansionFlags.current.isInitialLoad) {
          updateUrlWithConnectionId('');
        }

        expansionFlags.current.isHandlingExpansion = false;
      },
      renderExpandableRow: (rowData, tableMeta) => {
        const connection = filteredConnections[tableMeta.rowIndex];
        return (
          <TableCell colSpan={rowData.length}>
            <InnerTableContainer>
              <Table>
                <TableRow style={{ padding: 0 }}>
                  <TableCell style={{ overflowX: 'hidden', padding: 0 }}>
                    <Grid2 container style={{ textTransform: 'lowercase' }} size="grow">
                      <ContentContainer size={{ xs: 12 }}>
                        <FormatConnectionMetadata
                          connection={connection}
                          meshsyncControllerState={meshsyncControllerState}
                        />
                      </ContentContainer>
                    </Grid2>
                  </TableCell>
                </TableRow>
              </Table>
            </InnerTableContainer>
          </TableCell>
        );
      },
    }),
    [
      columns,
      connectionData?.totalCount,
      filteredConnections,
      handleDeleteConnections,
      meshsyncControllerState,
      page,
      pageSize,
      rowsExpanded,
      selectedConnectionId,
      sortOrder,
      theme.palette.common.white,
      theme.palette.error.dark,
      updateUrlWithConnectionId,
    ],
  );

  const [tableCols, updateCols] = useState(columns);

  useEffect(() => {
    updateCols(columns);
  }, [columns]);

  const [columnVisibility, setColumnVisibility] = useState<Record<string, boolean | undefined>>(
    () => getResponsiveColumnVisibility(columnNames, colViews, width),
  );

  useEffect(() => {
    setColumnVisibility(getResponsiveColumnVisibility(columnNames, colViews, width));
  }, [colViews, columnNames, width]);

  if (isConnectionLoading) {
    return <LoadingScreen animatedIcon="AnimatedMeshery" message="Loading Connections" />;
  }

  return (
    <>
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
          <div data-testid="ConnectionTable-search">
            <SearchBar
              onSearch={(value) => {
                setSearch(value);
              }}
              placeholder="Search Connections..."
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
            />
          </div>

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
        data={filteredConnections}
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
        <ActionListItem>
          <Button
            type="submit"
            onClick={handleFlushMeshSync()}
            data-cy="btnResetDatabase"
            disabled={!CAN(keys.FLUSH_MESHSYNC_DATA.action, keys.FLUSH_MESHSYNC_DATA.subject)}
          >
            <SyncIcon {...iconMedium} />
            <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
              Flush MeshSync
            </Typography>
          </Button>
        </ActionListItem>
        <ActionListItem>
          <Button
            type="submit"
            onClick={(e) => {
              e.stopPropagation();
              setDeploymentModeAnchorEl(e.currentTarget);
            }}
            data-cy="btnChangeDeploymentMode"
          >
            <Typography variant="body1">Modify Deployment Mode</Typography>
          </Button>
        </ActionListItem>
      </Popover>

      <Popover
        open={deploymentModeOpen}
        anchorEl={deploymentModeAnchorEl}
        onClose={handleDeploymentModeMenuClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
      >
        <ActionListItem>
          <Button onClick={() => handleDeploymentModeChange(MESHSYNC_DEPLOYMENT_TYPE.OPERATOR)}>
            <Typography variant="body1">Operator</Typography>
          </Button>
        </ActionListItem>
        <ActionListItem>
          <Button onClick={() => handleDeploymentModeChange(MESHSYNC_DEPLOYMENT_TYPE.EMBEDDED)}>
            <Typography variant="body1">Embedded</Typography>
          </Button>
        </ActionListItem>
      </Popover>
    </>
  );
};

export default ConnectionTable;
