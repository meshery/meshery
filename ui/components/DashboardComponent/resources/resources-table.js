import React, { useEffect, useMemo, useState } from 'react';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import {
  CustomColumnVisibilityControl,
  ResponsiveDataTable,
  SearchBar,
  Slide,
  UniversalFilter,
} from '@layer5/sistent';
import View from '../view';
import { ALL_VIEW, SINGLE_VIEW } from './config';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import { updateVisibleColumns } from '../../../utils/responsive-column';
import { useWindowDimensions } from '../../../utils/dimension';
import { camelcaseToSnakecase } from '../../../utils/utils';
import { useSelector } from 'react-redux';

import { useRouter } from 'next/router';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import { useGetMeshSyncResourceKindsQuery } from '@/rtk-query/meshsync';
import { api } from '@/rtk-query/index';
import { initiateQuery } from '@/rtk-query/utils';

export const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

const ResourcesTable = (props) => {
  const { updateProgress, k8sConfig, resourceConfig, submenu, workloadType, selectedK8sContexts } =
    props;
  const [meshSyncResources, setMeshSyncResources] = useState([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState();
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [namespaceFilter, setNamespaceFilter] = useState('');
  const [selectedFilters, setSelectedFilters] = useState({
    namespace: 'All',
  });
  const [selectedResource, setSelectedResource] = useState({});
  const [view, setView] = useState(ALL_VIEW);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const { width } = useWindowDimensions();
  const connectionMetadataState = useSelector((state) => state.get('connectionMetadataState'));
  const handleApplyFilter = () => {
    const namespaceFilter = selectedFilters.namespace === 'All' ? null : selectedFilters.namespace;
    setNamespaceFilter(namespaceFilter);
  };
  const clusterIds = getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig);
  const { data: clusterSummary } = useGetMeshSyncResourceKindsQuery(
    {
      page: page,
      pagesize: 'all',
      clusterIds: clusterIds,
    },
    { skip: clusterIds.length === 0 },
  );

  const filters = {
    namespace: {
      name: 'Namespace',
      options: [
        ...(clusterSummary?.namespaces || []).map((ns) => ({
          value: ns,
          label: ns,
        })),
      ],
    },
  };

  const switchView = (view, resource) => {
    setSelectedResource(resource);
    setView(view);
  };

  const tableConfig = submenu
    ? resourceConfig(
        switchView,
        meshSyncResources,
        k8sConfig,
        connectionMetadataState,
        workloadType,
        selectedK8sContexts,
      )[workloadType]
    : resourceConfig(
        switchView,
        meshSyncResources,
        k8sConfig,
        connectionMetadataState,
        workloadType,
        selectedK8sContexts,
      );

  const encodedClusterIds = encodeURIComponent(JSON.stringify(clusterIds));

  const { notify } = useNotification();

  const getMeshsyncResources = async (page, pageSize, search, sortOrder) => {
    setLoading(true);
    const { query } = router;
    const resourceName =
      query.resourceName ||
      (['Node', 'Namespace'].includes(query.resource) ? query.resource : search);
    const resourceCategory = query.resource || tableConfig.name;
    const decodedClusterIds = JSON.parse(decodeURIComponent(encodedClusterIds));
    if (decodedClusterIds.length === 0) {
      setLoading(false);
      return;
    }
    if (!resourceName) search = '';
    if (!sortOrder) sortOrder = '';

    const queryParams = {};
    queryParams.kind = resourceCategory;
    queryParams.status = true;
    queryParams.spec = true;
    queryParams.annotations = true;
    queryParams.labels = true;
    queryParams.clusterIds = encodedClusterIds;
    queryParams.page = page;
    queryParams.pagesize = pageSize;
    queryParams.search = encodeURIComponent(resourceName);
    queryParams.order = encodeURIComponent(sortOrder);

    if (namespaceFilter) {
      queryParams.namespace = encodeURIComponent(namespaceFilter);
    }

    try {
      const result = await initiateQuery(api.endpoints.getMeshSyncResources, queryParams);

      setMeshSyncResources(result?.resources || []);
      setPage(result?.page || 0);
      setCount(result?.total_count || 0);
      setPageSize(result?.page_size || 0);
      setLoading(false);

      if (query.resourceCategory && query.resourceName && result?.resources.length === 1) {
        switchView(SINGLE_VIEW, result?.resources[0]);
      }
    } catch (error) {
      handleError(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES)(error);
    }
  };

  const [tableCols, updateCols] = useState(tableConfig.columns);

  useEffect(() => {
    if (!loading) {
      getMeshsyncResources(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder, encodedClusterIds, namespaceFilter]);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(tableConfig.colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    tableConfig.columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  });
  const appendNameToQuery = (name) => {
    const currentQuery = { ...router.query, resourceName: name };
    router.push(
      {
        pathname: router.pathname,
        query: currentQuery,
      },
      undefined,
      { shallow: true },
    );
  };
  const router = useRouter();
  const options = useMemo(
    () => ({
      filter: false,
      viewColumns: false,
      search: false,
      responsive: 'standard',
      serverSide: true,
      selectableRows: 'none',
      count,
      rowsPerPage: pageSize,
      fixedHeader: true,
      page,
      print: false,
      download: false,
      textLabels: {
        selectedRows: {
          text: `${tableConfig.name}(s) selected`,
        },
      },
      enableNestedDataAccess: '.',
      onCellClick: (_, meta) => {
        if (meta.columnName !== 'cluster_id') {
          const currentResource = meshSyncResources[meta.rowIndex];
          if (currentResource) {
            switchView(SINGLE_VIEW, currentResource);
            appendNameToQuery(currentResource.metadata.name);
          }
        }
      },
      expandableRowsOnClick: true,
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        const columnName = camelcaseToSnakecase(tableConfig.columns[tableState.activeColumn]?.name);

        let order = '';
        if (tableState.activeColumn) {
          order = `${columnName} desc`;
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
                order = `${columnName} asc`;
              } else {
                order = `${columnName} desc`;
              }
            }
            if (order !== sortOrder) {
              setSortOrder(order);
            }
            break;
        }
      },
    }),
    [page, pageSize, meshSyncResources],
  );

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  return (
    <>
      <Slide in={view === SINGLE_VIEW} timeout={400} direction="left" mountOnEnter unmountOnExit>
        <div>
          {view === SINGLE_VIEW && (
            <View
              type={`${tableConfig.name}`}
              setView={setView}
              resource={selectedResource}
              k8sConfig={k8sConfig}
            />
          )}
        </div>
      </Slide>

      <Slide in={view === ALL_VIEW} timeout={400} direction="right" mountOnEnter unmountOnExit>
        <div>
          {view === ALL_VIEW && (
            <>
              <ToolWrapper style={{ marginBottom: '5px', marginTop: '1rem' }}>
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'end',
                    width: '100%',
                    alignItems: 'center',
                    gap: '0.5rem',
                  }}
                >
                  <SearchBar
                    onSearch={(value) => {
                      setSearch(value);
                    }}
                    expanded={isSearchExpanded}
                    setExpanded={setIsSearchExpanded}
                    placeholder={`Search ${tableConfig.name}...`}
                  />
                  <UniversalFilter
                    filters={filters}
                    selectedFilters={selectedFilters}
                    setSelectedFilters={setSelectedFilters}
                    handleApplyFilter={handleApplyFilter}
                  />
                  <CustomColumnVisibilityControl
                    id="ref"
                    columns={tableConfig.columns}
                    customToolsProps={{ columnVisibility, setColumnVisibility }}
                  />
                </div>
              </ToolWrapper>
              <ResponsiveDataTable
                data={meshSyncResources}
                columns={tableConfig.columns}
                options={options}
                tableCols={tableCols}
                updateCols={updateCols}
                columnVisibility={columnVisibility}
              />
            </>
          )}
        </div>
      </Slide>
    </>
  );
};

export default ResourcesTable;
