import React, { useEffect, useMemo, useState } from 'react';
import dataFetch from '../../../lib/data-fetch';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { CustomColumnVisibilityControl, ResponsiveDataTable, SearchBar } from '@layer5/sistent';
import useStyles from '../../../assets/styles/general/tool.styles';
import View from '../view';
import { ALL_VIEW, SINGLE_VIEW } from './config';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import { updateVisibleColumns } from '../../../utils/responsive-column';
import { useWindowDimensions } from '../../../utils/dimension';
import { camelcaseToSnakecase } from '../../../utils/utils';
import { useSelector } from 'react-redux';
import { UsesSistent } from '@/components/SistentWrapper';
import { Slide } from '@material-ui/core';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

const ResourcesTable = (props) => {
  const {
    classes,
    updateProgress,
    k8sConfig,
    resourceConfig,
    submenu,
    workloadType,
    selectedK8sContexts,
  } = props;
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
    ? resourceConfig(switchView, meshSyncResources, k8sConfig, connectionMetadataState)[
        workloadType
      ]
    : resourceConfig(switchView, meshSyncResources, k8sConfig, connectionMetadataState);

  const clusterIds = encodeURIComponent(
    JSON.stringify(getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig)),
  );

  const StyleClass = useStyles();

  const { notify } = useNotification();

  const getMeshsyncResources = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';
    dataFetch(
      `/api/system/meshsync/resources?kind=${
        tableConfig.name
      }&status=true&spec=true&annotations=true&labels=true&clusterIds=${clusterIds}&page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
        search,
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
      },
      handleError(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES),
    );
  };

  const [tableCols, updateCols] = useState(tableConfig.columns);

  useEffect(() => {
    if (!loading) {
      getMeshsyncResources(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder]);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(tableConfig.colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    tableConfig.columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  });

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
      onCellClick: (_, meta) =>
        meta.columnName !== 'cluster_id' &&
        switchView(SINGLE_VIEW, meshSyncResources[meta.rowIndex]),

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
    [page, pageSize],
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
      {view !== ALL_VIEW ? (
        <Slide
          in={view !== ALL_VIEW}
          timeout={400}
          direction={'left'}
          exit={true}
          enter={true}
          mountOnEnter
          unmountOnExit
        >
          <div>
            <View
              type={`${tableConfig.name}`}
              setView={setView}
              resource={selectedResource}
              classes={classes}
            />
          </div>
        </Slide>
      ) : (
        <div>
          <div
            className={StyleClass.toolWrapper}
            style={{ marginBottom: '5px', marginTop: '1rem' }}
          >
            <div className={classes.createButton}>{/* <MesherySettingsEnvButtons /> */}</div>
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
            </UsesSistent>
          </div>
          <UsesSistent>
            <ResponsiveDataTable
              data={meshSyncResources}
              columns={tableConfig.columns}
              options={options}
              className={classes.muiRow}
              tableCols={tableCols}
              updateCols={updateCols}
              columnVisibility={columnVisibility}
            />
          </UsesSistent>
        </div>
      )}
    </>
  );
};

export default ResourcesTable;
