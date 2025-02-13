import React, { useEffect, useMemo, useState } from 'react';
import dataFetch from '../../../lib/data-fetch';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import {
  CustomColumnVisibilityControl,
  ResponsiveDataTable,
  SearchBar,
  Slide,
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
  const [selectedResource, setSelectedResource] = useState({});
  const [view, setView] = useState(ALL_VIEW);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const { width } = useWindowDimensions();
  const connectionMetadataState = useSelector((state) => state.get('connectionMetadataState'));

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

  const clusterIds = encodeURIComponent(
    JSON.stringify(getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig)),
  );

  const { notify } = useNotification();

  const getMeshsyncResources = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    const { query } = router;
    const resourceName =
      query.resourceName ||
      (['Node', 'Namespace'].includes(query.resource) ? query.resource : search);
    const resourceCategory = query.resource || tableConfig.name;
    const decodedClusterIds = JSON.parse(decodeURIComponent(clusterIds));
    if (decodedClusterIds.length === 0) {
      setLoading(false);
      return;
    }
    if (!resourceName) search = '';
    if (!sortOrder) sortOrder = '';
    dataFetch(
      `/api/system/meshsync/resources?kind=${resourceCategory}&status=true&spec=true&annotations=true&labels=true&clusterIds=${clusterIds}&page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
        resourceName,
      )}&order=${encodeURIComponent(sortOrder)}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        setMeshSyncResources(res?.resources || []);
        setPage(res?.page || 0);
        setCount(res?.total_count || 0);
        setPageSize(res?.page_size || 0);
        setLoading(false);
        if (query.resourceCategory && query.resourceName && res?.resources.length === 1) {
          switchView(SINGLE_VIEW, res?.resources[0]);
        }
      },
      handleError(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES),
    );
  };

  const [tableCols, updateCols] = useState(tableConfig.columns);

  useEffect(() => {
    if (!loading) {
      getMeshsyncResources(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder, clusterIds]);

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
                <div style={{ display: 'flex', justifyContent: 'end', width: '100%' }}>
                  <SearchBar
                    onSearch={(value) => {
                      setSearch(value);
                    }}
                    expanded={isSearchExpanded}
                    setExpanded={setIsSearchExpanded}
                    placeholder={`Search ${tableConfig.name}...`}
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
