import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { PROMPT_VARIANTS, ResponsiveDataTable } from '@sistent/sistent';
import LoadingScreen from '../shared/LoadingState/LoadingComponent';
import { EVENT_TYPES } from '../../lib/event-types';
import _PromptComponent from '../PromptComponent';
import resetDatabase from '@/graphql/queries/ResetDatabaseQuery';

import { CONNECTION_KINDS, CONNECTION_STATES } from '../../utils/Enum';
import useKubernetesHook from '@/utils/hooks/useKubernetesHook';
import useGrafanaPingHook from '@/utils/hooks/useGrafanaPingHook';
import { getResponsiveColumnVisibility } from '../../utils/responsive-column';
import { useWindowDimensions } from '../../utils/dimension';
import { useGetEnvironmentsQuery } from '../../rtk-query/environments';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import { useTableUrlState } from '@/utils/hooks/useTableUrlState';
import { useColumnVisibilityPreference } from '@/utils/hooks/useColumnVisibilityPreference';

import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';

import type {
  ConnectionTableProps,
  ConnectionRow,
  EnvironmentOption,
  ExpansionFlags,
  RowData,
  SelectedFilters,
  SelectedRows,
} from './ConnectionTable.types';
import {
  ACTION_TYPES,
  CONNECTION_DOCS_URL,
  ENVIRONMENT_DOCS_URL,
  getErrorMessage,
  getStatusTransition,
} from './ConnectionTable.constants';
import type { ConnectionTransitionMap } from './ConnectionTable.constants';
import { useConnectionActions } from './ConnectionTable.hooks';
import { useConnectionColumns } from './ConnectionTable.columns';
import { useConnectionTableOptions } from './ConnectionTable.options';
import { ConnectionActionMenu, ConnectionDeploymentModeMenu } from './ConnectionActionMenu';
import { ConnectionTableToolbar } from './ConnectionTableToolbar';
import dynamic from 'next/dynamic';
import type { ConfigurableConnection } from './ConnectionConfigureModal';

// Lazy-loaded: it pulls in the RJSF/theme chain, which we keep out of the
// table's static import graph (smaller bundle + avoids eager theme init).
const ConnectionConfigureModal = dynamic(() => import('./ConnectionConfigureModal'), {
  ssr: false,
});

// Lazy-loaded for the same reason: only mounted for Kubernetes connections
// when the controllers configuration action is invoked.
const ConnectionControllersConfigModal = dynamic(
  () => import('./ConnectionControllersConfigModal'),
  { ssr: false },
);

const ConnectionTable = ({
  selectedFilter,
  selectedConnectionId,
  updateUrlWithConnectionId,
}: ConnectionTableProps) => {
  const {
    organization,
    connectionMetadataState,
    controllerState: meshsyncControllerState,
  } = useSelector(
    (state: {
      ui: {
        organization?: { id?: string };
        // `null` matches the Redux initial state (see
        // `store/slices/mesheryUi.ts`). The slice is only populated after
        // `_app.tsx`'s async `loadMeshModelComponent` resolves, which can
        // race the first render of this page.
        connectionMetadataState: Record<
          string,
          { transitions?: string[]; icon?: string; transitionMap?: ConnectionTransitionMap }
        > | null;
        controllerState: unknown;
      };
    }) => state.ui,
  );
  const ping = useKubernetesHook();
  const pingGrafana = useGrafanaPingHook();
  const { width } = useWindowDimensions();

  const { tableState, updateTableState, copyRowDeepLink } = useTableUrlState({
    tableKey: 'con',
    // Row deeplinks reuse the existing `connectionId` param so the parent's
    // expansion logic keeps working without changes.
    rowParam: 'connectionId',
    defaults: {
      page: 0,
      pageSize: 10,
      sortOrder: 'created_at desc',
      search: '',
      filters: { status: '', kind: '' },
    },
  });

  const { page, pageSize, sortOrder, search } = tableState;
  const setPage = useCallback((p: number) => updateTableState({ page: p }), [updateTableState]);
  const setPageSize = useCallback(
    (ps: number) => updateTableState({ pageSize: ps }),
    [updateTableState],
  );
  const setSortOrder = useCallback(
    (so: string) => updateTableState({ sortOrder: so }),
    [updateTableState],
  );
  const setSearch = useCallback(
    (s: string) => updateTableState({ search: s, page: 0 }),
    [updateTableState],
  );

  // Applied filters come from URL state so they survive navigation.
  const statusFilter = tableState.filters.status || null;
  const kindFilter = tableState.filters.kind || null;

  const [rowData, setRowData] = useState<RowData | null>(null);
  const [rowsExpanded, setRowsExpanded] = useState<number[]>([]);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [selectedFilters, setSelectedFilters] = useState<SelectedFilters>(() => ({
    status: tableState.filters.status || 'All',
    kind: tableState.filters.kind || 'All',
  }));
  const {
    notify,
    updateConnectionByIdMutator,
    addConnectionToEnvironment,
    removeConnectionFromEnvironment,
    saveEnvironment,
    updateConnectionStatus,
  } = useConnectionActions({ organizationId: organization?.id });
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const [deploymentModeAnchorEl, setDeploymentModeAnchorEl] = useState<HTMLElement | null>(null);
  const [configureConnection, setConfigureConnection] = useState<ConfigurableConnection | null>(
    null,
  );
  const [controllersConfigConnection, setControllersConfigConnection] =
    useState<ConfigurableConnection | null>(null);
  const open = Boolean(anchorEl);
  const deploymentModeOpen = Boolean(deploymentModeAnchorEl);
  const modalRef = useRef<{ show: (options: unknown) => Promise<string | null> } | null>(null);
  const lastNotifiedErrorsRef = useRef<{ environments: string; connections: string }>({
    environments: '',
    connections: '',
  });

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
    updateTableState({
      filters: {
        status: selectedFilters.status === 'All' ? '' : selectedFilters.status,
        kind: selectedFilters.kind === 'All' ? '' : selectedFilters.kind,
      },
      page: 0,
    });
  };
  // lock for not allowing multiple updates at the same time
  // needs to be a ref because it needs to be shared between renders
  // and useState loses reactivity when down table custom cells
  const updatingConnection = useRef(false);

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
    // RTK Query's `error` objects can change identity across renders while
    // remaining semantically the same (e.g. a failed request that stays in
    // error state). Emitting a snackbar on every such render can create a
    // feedback loop (snackbar state update -> re-render -> effect re-fire).
    // De-duplicate notifications by their rendered message.
    //
    // Note: this is intentionally local state (a ref) so it doesn't add
    // another state update into the render cycle.
    const last = lastNotifiedErrorsRef.current;

    if (isEnvironmentsError) {
      const message = `${ACTION_TYPES.FETCH_ENVIRONMENT.error_msg}: ${getErrorMessage(environmentsError)}`;
      if (last.environments !== message) {
        notify({
          message,
          event_type: EVENT_TYPES.ERROR,
        });
        last.environments = message;
      }
    } else {
      last.environments = '';
    }

    if (isConnectionError) {
      const message = `${ACTION_TYPES.FETCH_CONNECTIONS.error_msg}: ${getErrorMessage(connectionError)}`;
      if (last.connections !== message) {
        notify({
          message,
          event_type: EVENT_TYPES.ERROR,
        });
        last.connections = message;
      }
    } else {
      last.connections = '';
    }
  }, [connectionError, environmentsError, isConnectionError, isEnvironmentsError, notify]);

  const enhancedConnections = useMemo(() => {
    if (!connectionData?.connections) return [];

    // `connectionMetadataState` is `null` in the Redux initial state and is
    // only populated after `_app.tsx`'s async `loadMeshModelComponent`
    // completes. The pages-router routes to /management/connections before
    // that promise resolves, so this memo must tolerate a null map.
    // A connection only needs a kind and a status to render; the display name
    // falls back to `metadata.name`/kind in the Name column. Requiring a
    // top-level `name` here wrongly hid connections (e.g. kubernetes, grafana)
    // whose name lives only in `metadata.name`.
    return connectionData.connections
      .filter((conn) => conn && conn.kind && conn.status)
      .map((connection) => ({
        ...connection,
        nextStatus:
          connection.nextStatus || connectionMetadataState?.[connection.kind]?.transitions,
        kindLogo: connection.kindLogo || connectionMetadataState?.[connection.kind]?.icon,
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
      ['transitionMap', 'xs'],
      ['ConnectionID', 'na'],
    ],
    [],
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

  const handleDeploymentModeAnchorOpen = useCallback((event: React.MouseEvent<HTMLElement>) => {
    setDeploymentModeAnchorEl(event.currentTarget);
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

  const handleConfigureConnection = useCallback(() => {
    const connection = getConnectionAtRowIndex(rowData?.rowIndex);
    handleActionMenuClose();
    if (connection) {
      setConfigureConnection(connection as ConfigurableConnection);
    }
  }, [getConnectionAtRowIndex, handleActionMenuClose, rowData?.rowIndex]);

  const handleConfigureControllers = useCallback(() => {
    const connection = getConnectionAtRowIndex(rowData?.rowIndex);
    handleActionMenuClose();
    // Operator / MeshSync / Broker configuration applies to Kubernetes
    // connections only.
    if (connection && (connection as ConfigurableConnection).kind === 'kubernetes') {
      setControllersConfigConnection(connection as ConfigurableConnection);
    }
  }, [getConnectionAtRowIndex, handleActionMenuClose, rowData?.rowIndex]);

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

  // The previous shape was `useCallback(() => async () => {...})` invoked as
  // `handleFlushMeshSync()` in JSX, which created a new async closure on
  // every render and minted a fresh `onFlushMeshSync` prop reference each
  // commit. Returning the async function directly keeps the prop stable.
  const handleFlushMeshSync = useCallback(async () => {
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
  }, [getConnectionAtRowIndex, handleActionMenuClose, handleError, notify, rowData?.rowIndex]);

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
      const subtitle = getStatusTransition(
        connectionMetadataState?.[connectionKind]?.transitionMap,
        connectionStatus,
        status.toLowerCase(),
      );
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
    [connectionMetadataState, updateConnectionStatus],
  );

  const handleActionMenuOpen = useCallback((event, tableMeta: RowData) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
    setRowData(tableMeta);
  }, []);

  // Consolidate multiple useRef hooks into a single object
  const expansionFlags = useRef<ExpansionFlags>({
    isHandlingExpansion: false,
    isInitialLoad: true,
    isUrlExpansion: false,
    lastProcessedId: null,
  });

  // `filteredConnections` is accessed via a ref so RTK Query's identity
  // churn on every cache hit doesn't re-fire this effect. The effect re-fires
  // through `filteredConnectionsKey` (a primitive snapshot of the visible id
  // set), which only changes when the *content* of the visible page changes
  // — which is exactly the condition under which a previously-missing deep
  // link could now succeed (data finished loading, user paginated, filter
  // changed). Same-data refetches produce the same key string, so they bail
  // out of the effect via `Object.is` equality on the dep.
  const filteredConnectionsRef = useRef(filteredConnections);
  filteredConnectionsRef.current = filteredConnections;

  const filteredConnectionsKey = useMemo(
    () => filteredConnections.map((conn) => conn.id).join('|'),
    [filteredConnections],
  );

  useEffect(() => {
    if (!selectedConnectionId || expansionFlags.current.isHandlingExpansion) return;
    if (expansionFlags.current.lastProcessedId === selectedConnectionId) return;

    const connections = filteredConnectionsRef.current;
    if (!connections || connections.length === 0) {
      // Data not loaded yet. The effect will re-fire as soon as
      // `filteredConnectionsKey` flips on first arrival.
      return;
    }

    const index = connections.findIndex((conn) => conn.id === selectedConnectionId);
    if (index === -1) {
      // The deep-linked connection isn't on the current page. Do not mark
      // `lastProcessedId` — that would lock the effect out for the rest of
      // the session. If the user paginates or filters into a page that does
      // include this id, `filteredConnectionsKey` will change and the effect
      // re-runs to expand the row. Intentionally do NOT clear the URL: the
      // pre-fix code pushed `connectionId=""` here, which kicked off a
      // URL-push → re-render → effect-re-fire loop that surfaced as React
      // error #185.
      return;
    }

    expansionFlags.current.isUrlExpansion = true;
    expansionFlags.current.lastProcessedId = selectedConnectionId;
    setRowsExpanded([index]);
    expansionFlags.current.isUrlExpansion = false;
    expansionFlags.current.isInitialLoad = false;
  }, [selectedConnectionId, filteredConnectionsKey]);

  // Project the per-kind connection definitions down to just their state
  // machines for the status-transition dropdown.
  const transitionMapByKind = useMemo(() => {
    if (!connectionMetadataState) return null;
    return Object.fromEntries(
      Object.entries(connectionMetadataState).map(([kind, meta]) => [kind, meta?.transitionMap]),
    );
  }, [connectionMetadataState]);

  const columns = useConnectionColumns({
    url: CONNECTION_DOCS_URL,
    envUrl: ENVIRONMENT_DOCS_URL,
    environmentOptions,
    isEnvironmentsSuccess,
    updatingConnection,
    handleDeleteConnection,
    handleEnvironmentSelect,
    handleStatusChange,
    handleActionMenuOpen,
    ping,
    pingGrafana,
    transitionMapByKind,
  });
  const columnNames = useMemo(
    () => columns.map((column) => column.name),
    [columns, isEnvironmentsSuccess],
  );

  const options = useConnectionTableOptions({
    totalCount: connectionData?.totalCount,
    page,
    pageSize,
    setPage,
    setPageSize,
    sortOrder,
    setSortOrder,
    rowsExpanded,
    setRowsExpanded,
    columns,
    filteredConnections,
    meshsyncControllerState,
    selectedConnectionId,
    updateUrlWithConnectionId,
    expansionFlags,
    handleDeleteConnections,
  });

  const [tableCols, setTableCols] = useState(columns);

  // Keep the latest `columns` in a ref so the sync effect below can read them
  // without depending on `columns` identity — `columns` is rebuilt on most
  // renders (not all of its inputs are referentially stable), so a `[columns]`
  // dependency would setState every render and loop infinitely.
  const columnsRef = useRef(columns);
  columnsRef.current = columns;

  // ResponsiveDataTable renders cells from this `tableCols` snapshot and only
  // re-syncs it on columnVisibility identity changes, so a cell whose
  // `customBodyRender` closes over async data (the environments select is gated
  // on `isEnvironmentsSuccess`) would stay frozen at its first-render output.
  // Re-push the freshly built columns once those inputs settle — keyed on the
  // settling signals (not `columns`) so it runs only when the rendered output
  // can actually change.
  useEffect(() => {
    setTableCols(columnsRef.current);
  }, [isEnvironmentsSuccess, environmentOptions]);

  const { columnVisibility, setColumnVisibilityByUser, setColumnVisibilityByResponsive } =
    useColumnVisibilityPreference(
      'connections',
      getResponsiveColumnVisibility(columnNames, colViews, width),
    );

  useEffect(() => {
    const next = getResponsiveColumnVisibility(columnNames, colViews, width);

    // Only apply responsive update when the computed layout actually changed so
    // we avoid flushing user-preference overrides on every unrelated re-render.
    setColumnVisibilityByResponsive(next);
  }, [colViews, columnNames, width, setColumnVisibilityByResponsive]);

  if (isConnectionLoading) {
    return <LoadingScreen animatedIcon="AnimatedMeshery" message="Loading Connections" />;
  }

  return (
    <>
      <ConnectionTableToolbar
        isSearchExpanded={isSearchExpanded}
        setIsSearchExpanded={setIsSearchExpanded}
        onSearch={setSearch}
        filters={filters}
        selectedFilters={selectedFilters}
        setSelectedFilters={setSelectedFilters}
        handleApplyFilter={handleApplyFilter}
        columns={columns}
        columnVisibility={columnVisibility}
        setColumnVisibility={setColumnVisibilityByUser}
      />

      <ResponsiveDataTable
        data={filteredConnections}
        columns={columns}
        options={options}
        tableCols={tableCols}
        updateCols={setTableCols}
        columnVisibility={columnVisibility}
      />

      <_PromptComponent ref={modalRef} />
      <ConnectionActionMenu
        anchorEl={anchorEl}
        open={open}
        onClose={handleActionMenuClose}
        onFlushMeshSync={handleFlushMeshSync}
        onDeploymentModeAnchor={handleDeploymentModeAnchorOpen}
        onConfigure={handleConfigureConnection}
        onConfigureControllers={
          rowData?.rowIndex != null &&
          (getConnectionAtRowIndex(rowData.rowIndex) as ConfigurableConnection | null)?.kind ===
            'kubernetes'
            ? handleConfigureControllers
            : undefined
        }
        onCopyLink={
          rowData?.rowIndex != null
            ? () => {
                const connection = filteredConnections[rowData.rowIndex];
                if (connection?.id) copyRowDeepLink(connection.id);
              }
            : undefined
        }
      />

      {/* Only mount (and thus load) the configure modal once a row is chosen. */}
      {configureConnection && (
        <ConnectionConfigureModal
          isOpen={Boolean(configureConnection)}
          connection={configureConnection}
          onClose={() => setConfigureConnection(null)}
        />
      )}

      {controllersConfigConnection?.id && (
        <ConnectionControllersConfigModal
          isOpen={Boolean(controllersConfigConnection)}
          connectionId={String(controllersConfigConnection.id)}
          connectionName={controllersConfigConnection.name}
          onClose={() => setControllersConfigConnection(null)}
        />
      )}

      <ConnectionDeploymentModeMenu
        anchorEl={deploymentModeAnchorEl}
        open={deploymentModeOpen}
        onClose={handleDeploymentModeMenuClose}
        onSelectMode={handleDeploymentModeChange}
      />
    </>
  );
};

export default ConnectionTable;
