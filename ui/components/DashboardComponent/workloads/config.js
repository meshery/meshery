import React, { useEffect, useMemo, useState } from 'react';
import { TableCell, TableSortLabel } from '@material-ui/core';
import dataFetch from '../../../lib/data-fetch';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import ResponsiveDataTable from '../../../utils/data-table';
import CustomColumnVisibilityControl from '../../../utils/custom-column';
import useStyles from '../../../assets/styles/general/tool.styles';
import SearchBar from '../../../utils/custom-search';
import { timeAgo } from '../../../utils/k8s-utils';
import { getClusterNameFromClusterId } from '../../../utils/multi-ctx';
import View from '../view';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

export const WorkloadsTypes = {
  PODS: {
    name: 'Pod',
    columns: [],
  },
  DEPLOYMENT: {
    name: 'Deployment',
    columns: [],
  },
  DAEMONSETS: {
    name: 'DaemonSet',
    columns: [],
  },
  STATEFULSETS: {
    name: 'StatefulSet',
    columns: [],
  },
  REPLICASETS: {
    name: 'ReplicaSet',
    columns: [],
  },
  REPLICATIONCONTROLLERS: {
    name: 'ReplicationController',
    columns: [],
  },
  JOBS: {
    name: 'Job',
    columns: [],
  },
  CRONJOBS: {
    name: 'CronJob',
    columns: [],
  },
};

const StandardWorkloadTable = (props) => {
  console.log(props);
  const { classes, updateProgress, k8sConfig, workloadType } = props;
  const ALL_WORKLOAD = 'all';
  const SINGLE_WORKLOAD = 'single';
  const [meshSyncResources, setMeshSyncResources] = useState([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [selectedResource, setSelectedResource] = useState({});
  const [view, setView] = useState(ALL_WORKLOAD);

  const swtichView = (view, resource) => {
    setSelectedResource(resource);
    setView(view);
  };

  const workloadsTypes = {
    PODS: {
      name: 'Pod',
      columns: [
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
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <>
                  <div
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                    onClick={() =>
                      swtichView(SINGLE_WORKLOAD, meshSyncResources[tableMeta.rowIndex])
                    }
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
          label: 'Phase',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let phase = attribute.phase;
              return <>{phase}</>;
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
          label: 'Host IP',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let hostIP = attribute.hostIP;
              return <>{hostIP}</>;
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
          label: 'Pod IP',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let podIP = attribute.podIP;
              return <>{podIP}</>;
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
          name: 'metadata.namespace',
          label: 'Namespace',
          options: {
            sort: false,
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <>
                  <div
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                    onClick={() =>
                      swtichView(SINGLE_WORKLOAD, meshSyncResources[tableMeta.rowIndex])
                    }
                  >
                    {value}
                  </div>
                </>
              );
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Node',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let nodeName = attribute.nodeName;
              return <>{nodeName}</>;
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
      ],
    },
    DEPLOYMENT: {
      name: 'Deployment',
      columns: [],
    },
    DAEMONSETS: {
      name: 'DaemonSet',
      columns: [],
    },
    STATEFULSETS: {
      name: 'StatefulSet',
      columns: [],
    },
    REPLICASETS: {
      name: 'ReplicaSet',
      columns: [],
    },
    REPLICATIONCONTROLLERS: {
      name: 'ReplicationController',
      columns: [],
    },
    JOBS: {
      name: 'Job',
      columns: [],
    },
    CRONJOBS: {
      name: 'CronJob',
      columns: [],
    },
  };

  const StyleClass = useStyles();

  const { notify } = useNotification();

  const getMeshsyncResources = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';
    dataFetch(
      `/api/system/meshsync/resources?kind=${
        WorkloadsTypes[workloadType].name
      }&status=true&spec=true&annotations=true&labels=true&page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
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

  const [tableCols, updateCols] = useState();

  useEffect(() => {
    updateCols(workloadsTypes[workloadType].columns);
    if (!loading) {
      getMeshsyncResources(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder]);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    workloadsTypes[workloadType].columns.forEach((col) => {
      initialVisibility[col.name] = col.options?.display !== false;
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
          text: `${WorkloadsTypes[workloadType].name}(s) selected`,
        },
      },
      enableNestedDataAccess: '.',
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        let order = '';
        if (tableState.activeColumn) {
          order = `${workloadsTypes[workloadType].columns[tableState.activeColumn].name} desc`;
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
                order = `${workloadsTypes[workloadType].columns[tableState.activeColumn].name} asc`;
              } else {
                order = `${
                  workloadsTypes[workloadType].columns[tableState.activeColumn].name
                } desc`;
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
      {view === ALL_WORKLOAD ? (
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
                placeholder={`Search ${WorkloadsTypes[workloadType].name}...`}
              />

              <CustomColumnVisibilityControl
                columns={workloadsTypes[workloadType].columns}
                customToolsProps={{ columnVisibility, setColumnVisibility }}
              />
            </div>
          </div>
          <ResponsiveDataTable
            data={meshSyncResources}
            columns={workloadsTypes[workloadType].columns}
            options={options}
            className={classes.muiRow}
            tableCols={tableCols}
            updateCols={updateCols}
            columnVisibility={columnVisibility}
          />
        </>
      ) : (
        <>
          <View
            type={`${WorkloadsTypes[workloadType].name}`}
            setView={setView}
            resource={selectedResource}
            classes={classes}
          />
        </>
      )}
    </>
  );
};

const SortableTableCell = ({ index, columnData, columnMeta, onSort }) => {
  return (
    <TableCell key={index} onClick={onSort}>
      <TableSortLabel active={columnMeta.name === columnData.name}>
        <b>{columnData.label}</b>
      </TableSortLabel>
    </TableCell>
  );
};

export default StandardWorkloadTable;
