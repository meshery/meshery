import React, { useEffect, useMemo, useState } from 'react';
import { TableCell, TableSortLabel } from '@material-ui/core';
import dataFetch from '../../lib/data-fetch';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import ResponsiveDataTable from '../../utils/data-table';
// import { FormattedMetadata } from '../NotificationCenter/metadata';
import CustomColumnVisibilityControl from '../../utils/custom-column';
import useStyles from '../../assets/styles/general/tool.styles';
import SearchBar from '../../utils/custom-search';
import { getResourceStr, resourceParsers, timeAgo } from '../../utils/k8s-utils';
import { getClusterNameFromClusterId } from '../../utils/multi-ctx';
// import { TextWithLinks } from '../DataFormatter';
// import { Link } from '../DataFormatter';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

const Nodes = ({ classes, updateProgress, k8sConfig }) => {
  const ALL_NODES = 'all';
  const SINGLE_NODE = 'single';
  // const availableViews = [ALL_NODES, SINGLE_NODE];

  const [meshSyncResources, setMeshSyncResources] = useState([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [view, setView] = useState(ALL_NODES);

  const swtichView = (view) => {
    setView(view);
  };

  const StyleClass = useStyles();

  const { notify } = useNotification();

  const getMeshsyncResources = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';
    dataFetch(
      `/api/system/meshsync/resources?kind=Node&status=true&page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
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

  const columns = [
    {
      name: 'id',
      label: 'ID',
      options: {
        display: false,
      },
    },
    {
      name: 'metadata.name',
      label: 'Name',
      options: {
        sort: false,
        sortThirdClickReset: true,
        display: true,
        customBodyRender: (value) => {
          return (
            <>
              <div
                style={{
                  color: 'inherit',
                  textDecorationLine: 'underline',
                  cursor: 'pointer',
                  marginBottom: '0.5rem',
                }}
                onClick={() => swtichView(SINGLE_NODE)}
              >
                {value}
              </div>
            </>
          );
        },
      },
    },
    {
      name: 'apiVersion',
      label: 'API version',
      options: {
        sort: false,
        display: true,
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
      name: 'status.attribute',
      label: 'CPU',
      options: {
        sort: false,
        display: true,
        customBodyRender: function CustomBody(val) {
          let attribute = JSON.parse(val);
          let capacity = attribute.capacity;
          let cpu = getResourceStr(resourceParsers['cpu'](capacity.cpu), 'cpu');
          return <>{cpu}</>;
        },
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
      name: 'status.attribute',
      label: 'Memory',
      options: {
        sort: false,
        display: true,
        customBodyRender: function CustomBody(val) {
          let attribute = JSON.parse(val);
          let capacity = attribute.capacity;
          let memory = getResourceStr(resourceParsers['memory'](capacity.memory), 'memory');
          return <>{memory}</>;
        },
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
      name: 'cluster_id',
      label: 'Cluster',
      options: {
        sort: false,
        sortThirdClickReset: true,
        customBodyRender: function CustomBody(val) {
          let clusterName = getClusterNameFromClusterId(val, k8sConfig);
          return (
            <>
              <a
                href={'#'}
                target="_blank"
                rel="noopener noreferrer"
                style={{
                  color: 'inherit',
                  textDecorationLine: 'underline',
                  cursor: 'pointer',
                  marginBottom: '0.5rem',
                }}
              >
                {clusterName}
              </a>
            </>
          );
        },
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
      name: 'status.attribute',
      label: 'Internal IP',
      options: {
        sort: false,
        sortThirdClickReset: true,
        customBodyRender: function CustomBody(val) {
          let attribute = JSON.parse(val);
          let addresses = attribute.addresses;
          let internalIP =
            addresses.find((address) => address.type === 'InternalIP')?.address || '';
          return <>{internalIP}</>;
        },
      },
    },
    {
      name: 'status.attribute',
      label: 'External IP',
      options: {
        sort: false,
        sortThirdClickReset: true,
        customBodyRender: function CustomBody(val) {
          let attribute = JSON.parse(val);
          let addresses = attribute.addresses;
          let externalIP =
            addresses.find((address) => address.type === 'ExternalIP')?.address || '';
          return <>{externalIP}</>;
        },
      },
    },
    {
      name: 'metadata.creationTimestamp',
      label: 'Age',
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
        customBodyRender: function CustomBody(value) {
          let time = timeAgo(value);
          return <>{time}</>;
        },
      },
    },
  ];

  const options = useMemo(
    () => ({
      filter: false,
      viewColumns: false,
      search: false,
      responsive: 'standard',
      serverSide: true,
      selectableRows: false,
      count,
      rowsPerPage: pageSize,
      rowsPerPageOptions: [10, 25, 30],
      fixedHeader: true,
      page,
      print: false,
      download: false,
      textLabels: {
        selectedRows: {
          text: 'nodes(s) selected',
        },
      },
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
    }),
    [page, pageSize],
  );

  useEffect(() => {
    if (!loading) {
      getMeshsyncResources(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder]);

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = col.options?.display !== false;
    });
    return initialVisibility;
  });

  return (
    <>
      {view === ALL_NODES ? (
        <>
          <div
            className={StyleClass.toolWrapper}
            style={{ marginBottom: '5px', marginTop: '1rem' }}
          >
            <div className={classes.createButton}>{/* <MesherySettingsEnvButtons /> */}</div>
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
                placeholder="Search nodes..."
              />

              <CustomColumnVisibilityControl
                columns={columns}
                customToolsProps={{ columnVisibility, setColumnVisibility }}
              />
            </div>
          </div>
          <ResponsiveDataTable
            data={meshSyncResources}
            columns={columns}
            options={options}
            className={classes.muiRow}
            tableCols={tableCols}
            updateCols={updateCols}
            columnVisibility={columnVisibility}
          />
        </>
      ) : (
        <></>
      )}
    </>
  );
};

const SortableTableCell = ({ index, columnData, columnMeta, onSort }) => {
  return (
    <TableCell key={index} onClick={onSort}>
      <TableSortLabel
        active={columnMeta.name === columnData.name}
        // direction={columnMeta.direction || 'asc'}
      >
        <b>{columnData.label}</b>
      </TableSortLabel>
    </TableCell>
  );
};

export default Nodes;
