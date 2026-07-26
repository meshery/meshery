import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Tooltip, Grid2, FormControl, MenuItem, Table, FormattedTime } from '@sistent/sistent';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import {
  CustomColumnVisibilityControl,
  ResponsiveDataTable,
  SearchBar,
  UniversalFilter,
  TableCell,
  TableRow,
} from '@sistent/sistent';
import { MeshSyncDataFormatter } from '../metadata';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import { DefaultTableCell, SortableTableCell } from '../common';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import {
  JsonParse,
  camelcaseToSnakecase,
  getColumnValue,
  getVisibilityColums,
} from '../../../utils/utils';
import RegisterConnectionModal from '../RegisterConnectionModal';
import { CONNECTION_STATES, MESHSYNC_STATES } from '../../../utils/Enum';
import { getResponsiveColumnVisibility } from '../../../utils/responsive-column';
import { useWindowDimensions } from '../../../utils/dimension';
import { FormatId } from '../../data-formatter';
import {
  useGetMeshSyncResourceKindsQuery,
  useGetMeshSyncResourcesQuery,
} from '@/rtk-query/meshsync';
import { ConnectionStateChip } from '../ConnectionChip';
import { ContentContainer, ConnectionStyledSelect, InnerTableContainer } from '../styles';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import MeshSyncEmptyState from './MeshSyncEmptyState';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

export default function MeshSyncTable(props) {
  const { selectedResourceId, updateUrlWithResourceId, tabs } = props;
  const callbackRef = useRef();
  const [openRegistrationModal, setRegistrationModal] = useState(false);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('creation_timestamp desc');
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [modelFilter, setModeFilter] = useState();
  const [kindFilter, setKindFilter] = useState();
  const [namespaceFilter, setNamespaceFilter] = useState();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const { k8sConfig, selectedK8sContexts } = useSelector((state) => state.ui);
  const [selectedFilters, setSelectedFilters] = useState({
    kind: 'All',
    model: 'All',
    namespace: 'All',
  });
  const [registerConnection, setRegisterConnection] = useState({
    metadata: {},
    kind: '',
  });
  const { width } = useWindowDimensions();

  const { notify } = useNotification();
  const clusterIds = useMemo(
    () => getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig),
    [k8sConfig, selectedK8sContexts],
  );
  const serializedClusterIds = useMemo(() => JSON.stringify(clusterIds), [clusterIds]);

  const handleRegistrationModalClose = () => {
    setRegistrationModal(false);
  };

  const {
    data: meshSyncData,
    isError: isError,
    error: meshSyncError,
  } = useGetMeshSyncResourcesQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    kind: kindFilter,
    model: modelFilter,
    namespace: namespaceFilter,
    clusterIds: serializedClusterIds,
  });

  useEffect(() => {
    if (!isError) {
      return;
    }

    notify({
      message: 'Error fetching MeshSync Resources',
      event_type: EVENT_TYPES.ERROR,
      details: meshSyncError?.data,
    });
  }, [isError, meshSyncError?.data, notify]);

  const { data: clusterSummary } = useGetMeshSyncResourceKindsQuery({
    page: page,
    pagesize: 'all',
    search: search,
    order: sortOrder,
    clusterIds: clusterIds,
  });
  const availableKinds = (clusterSummary?.kinds || []).map((kind) => kind.Kind);
  const availableModels = [
    ...new Set((clusterSummary?.kinds || []).map((kind) => kind.Model).filter(Boolean)),
  ];
  const availableNamespaces = clusterSummary?.namespaces || [];
  const meshSyncResources = meshSyncData?.resources || [];

  const colViews = useMemo(
    () => [
      ['metadata.name', 'xs'],
      ['apiVersion', 'na'],
      ['kind', 'm'],
      ['model', 'm'],
      ['metadata.namespace', 'xs'],
      ['cluster_id', 'na'],
      ['pattern_resources', 'na'],
      ['metadata.creationTimestamp', 'l'],
      ['status', 'xs'],
      ['metadata', 'na'],
    ],
    [],
  );
  const columnNames = useMemo(() => colViews.map(([columnName]) => columnName), [colViews]);

  const columns = [
    {
      name: 'metadata.name',
      sortName: 'name',
      label: 'Name',
      options: {
        sort: true,
        sortThirdClickReset: true,
        sortName: 'name',
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={{ ...column, name: 'name' }}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: (value) => {
          const maxCharLength = 30;
          const shouldTruncate = value?.length > maxCharLength;

          return (
            <Tooltip title={value} placement="top">
              <div
                style={{
                  overflow: 'hidden',
                  whiteSpace: 'nowrap',
                  textOverflow: shouldTruncate ? 'ellipsis' : 'none',
                }}
              >
                {shouldTruncate ? `${value.slice(0, maxCharLength)}...` : value}
              </div>
            </Tooltip>
          );
        },
      },
    },
    {
      name: 'apiVersion',
      label: 'API version',
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
      name: 'model',
      label: 'Model',
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
      name: 'metadata.namespace',
      sortName: 'namespace',
      label: 'Namespace',
      options: {
        sort: true,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={{ ...column, name: 'namespace' }}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'cluster_id',
      label: 'Cluster ID',
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
        customBodyRender: (value) => <FormatId id={value} />,
      },
    },
    {
      name: 'pattern_resources',
      label: 'Pattern resources',
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
        customBodyRender: (value) => {
          const maxCharLength = 30;
          const shouldTruncate = value?.length > maxCharLength;

          return (
            <Tooltip title={value} placement="top">
              <div
                style={{
                  overflow: 'hidden',
                  whiteSpace: 'nowrap',
                  textOverflow: shouldTruncate ? 'ellipsis' : 'none',
                }}
              >
                {shouldTruncate ? `${value.slice(0, maxCharLength)}...` : value}
              </div>
            </Tooltip>
          );
        },
      },
    },
    {
      name: 'metadata.creationTimestamp',
      sortName: 'creation_timestamp',
      label: 'Discovered At',
      options: {
        sort: true,
        sortThirdClickReset: true,
        sortName: 'creation_timestamp',
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={{ ...column, name: 'creation_timestamp' }}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: function CustomBody(value) {
          return <FormattedTime date={value} />;
        },
      },
    },
    {
      name: 'status',
      label: 'Status',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          const componentMetadata = getColumnValue(
            tableMeta.rowData,
            'component_metadata',
            columns,
          );
          const DISCOVERED = {
            DISCOVERED: MESHSYNC_STATES.DISCOVERED,
          };
          const meshSyncStates =
            componentMetadata?.capabilities?.connection === true ? MESHSYNC_STATES : DISCOVERED;
          const disabled =
            componentMetadata?.capabilities?.connection === true &&
            value !== CONNECTION_STATES.REGISTERED
              ? false
              : true;
          return (
            <FormControl>
              <ConnectionStyledSelect
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                defaultValue={MESHSYNC_STATES.DISCOVERED}
                disabled={disabled}
                onClick={(e) => {
                  e.stopPropagation();
                  const clickedValue = e.target.value;
                  if (clickedValue !== MESHSYNC_STATES.DISCOVERED && clickedValue !== value) {
                    setRegistrationModal((open) => !open);
                  }
                }}
                onChange={() => {
                  callbackRef?.current?.(tableMeta);
                  setRegisterConnection({
                    capabilities: componentMetadata?.capabilities,
                    metadata: JsonParse(componentMetadata.metadata),
                    resourceID: tableMeta.rowData[tableMeta.rowData.length - 1],
                  });
                }}
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
                {Object.keys(meshSyncStates).map((s) => (
                  <MenuItem
                    disabled={
                      meshSyncStates[s] === value ||
                      meshSyncStates[s] === CONNECTION_STATES.REGISTERED
                        ? true
                        : false
                    }
                    value={meshSyncStates[s]}
                    key={meshSyncStates[s]}
                    style={{ padding: '0' }}
                  >
                    <ConnectionStateChip status={meshSyncStates[s]} />
                  </MenuItem>
                ))}
              </ConnectionStyledSelect>
            </FormControl>
          );
        },
      },
    },
    {
      name: 'component_metadata',
      label: 'Component Metadata',
      options: {
        display: false,
      },
    },
    {
      name: 'id',
      label: 'Resource ID',
      options: {
        display: false,
      },
    },
    {
      name: 'metadata',
      label: 'Metadata',
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
    serverSide: true,
    selectableRows: 'none',
    count: meshSyncData?.totalCount,
    rowsPerPage: pageSize,
    fixedHeader: true,
    page,
    print: false,
    download: false,
    sortOrder: {
      name: sortOrder.split(' ')[0].replace('metadata.', ''),
      direction: sortOrder.split(' ')[1],
    },
    textLabels: {
      selectedRows: {
        text: 'connection(s) selected',
      },
    },
    enableNestedDataAccess: '.',
    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      const activeColumn = columns[tableState.activeColumn];

      switch (action) {
        case 'changePage':
          setPage(tableState.page.toString());
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage.toString());
          break;
        case 'sort':
          if (sortInfo.length == 2) {
            const columnName = activeColumn?.sortName || camelcaseToSnakecase(activeColumn?.name);
            const direction = sortInfo[1] === 'ascending' ? 'asc' : 'desc';
            order = `${columnName} ${direction}`;
            if (order !== sortOrder) {
              setSortOrder(order);
            }
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
      if (expansionFlags.current.isUrlExpansion) return;

      expansionFlags.current.isHandlingExpansion = true;
      const expandedRows = allRowsExpanded.slice(-1);
      // Bail out on a value-equal write — see ConnectionTable.options.tsx
      // for the rationale.
      const nextExpandedRows = expandedRows.map((item) => item.index);
      const hasExpandedRowsChanged =
        nextExpandedRows.length !== rowsExpanded.length ||
        nextExpandedRows.some((rowIndex, index) => rowIndex !== rowsExpanded[index]);

      if (hasExpandedRowsChanged) {
        setRowsExpanded(nextExpandedRows);
      }

      if (expandedRows.length > 0 && meshSyncResources) {
        const index = expandedRows[0].index;
        const resource = meshSyncResources[index];

        if (
          resource &&
          updateUrlWithResourceId &&
          (!expansionFlags.current.isInitialLoad || resource.id !== selectedResourceId)
        ) {
          updateUrlWithResourceId(resource.id);
        }
      } else if (updateUrlWithResourceId && !expansionFlags.current.isInitialLoad) {
        updateUrlWithResourceId('');
      }

      expansionFlags.current.isHandlingExpansion = false;
    },
    renderExpandableRow: (rowData) => {
      const colSpan = rowData.length;
      const columnName = 'metadata'; // Name of the column containing the metadata
      const columnIndex = columns.findIndex((column) => column.name === columnName);

      // Access the metadata value using the column index
      const metadata = rowData[columnIndex];

      return (
        <TableCell colSpan={colSpan}>
          <InnerTableContainer>
            <Table>
              <TableRow>
                <TableCell>
                  <Grid2 container style={{ textTransform: 'lowercase' }} size="grow">
                    <ContentContainer
                      size={{
                        xs: 12,
                      }}
                      style={{
                        display: 'flex',
                        flexWrap: 'wrap',
                        padding: '0 20px',
                        gap: 30,
                      }}
                    >
                      <MeshSyncDataFormatter metadata={metadata} />
                    </ContentContainer>
                  </Grid2>
                </TableCell>
              </TableRow>
            </Table>
          </InnerTableContainer>
        </TableCell>
      );
    },
  };

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  useEffect(() => {
    if (meshSyncError) {
      handleError(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES)(meshSyncError);
    }
  }, [meshSyncError]);

  // Consolidate multiple useRef hooks into a single object
  const expansionFlags = useRef({
    isHandlingExpansion: false,
    isInitialLoad: true,
    isUrlExpansion: false,
    lastProcessedId: null,
  });

  // Mirror of the ConnectionTable fix: `meshSyncResources` is accessed via a
  // ref so RTK Query's identity churn on cache hits doesn't re-fire this
  // effect, while `meshSyncResourcesKey` (a primitive snapshot of the visible
  // id set) keys the effect to actual content changes — which is exactly the
  // condition under which a previously-missing deep link could now succeed
  // (data finished loading, user paginated, filter changed). Same-data
  // refetches produce the same key string and bail out via `Object.is`.
  const meshSyncResourcesRef = useRef(meshSyncResources);
  meshSyncResourcesRef.current = meshSyncResources;

  const meshSyncResourcesKey = useMemo(
    () => (meshSyncResources ?? []).map((resource) => resource.id).join('|'),
    [meshSyncResources],
  );

  useEffect(() => {
    if (!selectedResourceId || expansionFlags.current.isHandlingExpansion) return;
    if (expansionFlags.current.lastProcessedId === selectedResourceId) return;

    const resources = meshSyncResourcesRef.current;
    if (!resources || resources.length === 0) {
      // Data not loaded yet. Re-fires once `meshSyncResourcesKey` flips.
      return;
    }

    const index = resources.findIndex((resource) => resource.id === selectedResourceId);
    if (index === -1) {
      // Resource not on this page — do not lock out the effect by marking
      // `lastProcessedId` here. If the user paginates or filters into a page
      // that does include this id, `meshSyncResourcesKey` will change and the
      // effect re-runs. Intentionally do NOT clear the URL: pushing
      // `connectionId=""` started a URL-push → re-render → effect-re-fire
      // loop that surfaced as React error #185.
      return;
    }

    expansionFlags.current.isUrlExpansion = true;
    expansionFlags.current.lastProcessedId = selectedResourceId;
    setRowsExpanded([index]);
    expansionFlags.current.isUrlExpansion = false;
    expansionFlags.current.isInitialLoad = false;
  }, [selectedResourceId, meshSyncResourcesKey]);

  const filters = {
    kind: {
      name: 'Kind',
      options: [
        ...availableKinds.map((kind) => ({
          value: kind,
          label: kind,
        })),
      ],
    },
    model: {
      name: 'Model',
      options: [
        ...availableModels.map((model) => ({
          value: model,
          label: model,
        })),
      ],
    },
    namespace: {
      name: 'Namespace',
      options: [
        ...availableNamespaces.map((ns) => ({
          value: ns,
          label: ns,
        })),
      ],
    },
  };

  const handleApplyFilter = () => {
    const kindFilter = selectedFilters.kind === 'All' ? null : selectedFilters.kind;
    const modelFilter = selectedFilters.model === 'All' ? null : selectedFilters.model;
    const namespaceFilter = selectedFilters.namespace === 'All' ? null : selectedFilters.namespace;

    setKindFilter(kindFilter);
    setModeFilter(modelFilter);
    setNamespaceFilter(namespaceFilter);
  };

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() =>
    getResponsiveColumnVisibility(columnNames, colViews, width),
  );

  useEffect(() => {
    setColumnVisibility((previous) => {
      const next = getResponsiveColumnVisibility(columnNames, colViews, width);

      // Bail out on a value-equal recomputation. Without this guard
      // `getResponsiveColumnVisibility` (which always allocates a new object)
      // enqueued a commit on every parent re-render even when no responsive
      // breakpoint had crossed.
      const keys = new Set([...Object.keys(previous), ...Object.keys(next)]);
      for (const key of keys) {
        if (previous[key] !== next[key]) {
          return next;
        }
      }
      return previous;
    });
  }, [colViews, columnNames, width]);

  return (
    <>
      <ToolWrapper style={{ marginBottom: '5px', marginTop: '-30px' }}>
        <div
          style={{
            display: 'flex',
            borderRadius: '0.5rem 0.5rem 0 0',
            width: '100%',
            justifyContent: 'end',
          }}
        >
          <SearchBar
            onSearch={(value) => {
              setSearch(value);
            }}
            expanded={isSearchExpanded}
            setExpanded={setIsSearchExpanded}
            placeholder="Search Connections..."
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
      </ToolWrapper>

      {tabs}

      {!meshSyncResources || meshSyncResources.length === 0 ? (
        <MeshSyncEmptyState />
      ) : (
        <ResponsiveDataTable
          data={meshSyncResources}
          columns={columns}
          options={options}
          tableCols={tableCols}
          updateCols={updateCols}
          columnVisibility={columnVisibility}
        />
      )}
      <RegisterConnectionModal
        handleRegistrationModalClose={handleRegistrationModalClose}
        openRegistrationModal={openRegistrationModal}
        connectionData={registerConnection}
      />
    </>
  );
}
