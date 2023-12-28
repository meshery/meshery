import React, { useEffect, useMemo, useState } from 'react';
import { TableCell, Tooltip, TableContainer, Table, Grid, TableRow } from '@material-ui/core';
import Moment from 'react-moment';
import dataFetch from '../../../lib/data-fetch';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { ResponsiveDataTable } from '@layer5/sistent-components';
import CustomColumnVisibilityControl from '../../../utils/custom-column';
import useStyles from '../../../assets/styles/general/tool.styles';
import SearchBar from '../../../utils/custom-search';
import { MeshSyncDataFormatter } from '../metadata';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import { DefaultTableCell, SortableTableCell } from '../common';
import { camelcaseToSnakecase } from '../../../utils/utils';
import UniversalFilter from '../../../utils/custom-filter';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

export default function MeshSyncTable(props) {
  const { classes, updateProgress, selectedK8sContexts, k8sconfig } = props;
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState('');
  const [meshSyncResources, setMeshSyncResources] = useState([]);
  const [sortOrder, setSortOrder] = useState('');
  const [showMore, setShowMore] = useState(false);
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [loading, setLoading] = useState(false);
  const StyleClass = useStyles();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  const clusterIds = encodeURIComponent(
    JSON.stringify(getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sconfig)),
  );

  const { notify } = useNotification();

  const columns = [
    {
      name: 'metadata.name',
      label: 'Name',
      options: {
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
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
      label: 'Discovered At',
      options: {
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value) {
          return (
            <Tooltip
              title={
                <Moment startOf="day" format="LLL">
                  {value}
                </Moment>
              }
              placement="top"
              arrow
              interactive
            >
              <Moment format="LL">{value}</Moment>
            </Tooltip>
          );
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
      // resizableColumns: true,
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
          text: 'connection(s) selected',
        },
      },
      // customToolbarSelect: (selected) => (
      //   <Button
      //     variant="contained"
      //     color="primary"
      //     size="large"
      //     // @ts-ignore
      //     // onClick={() => handleDeleteConnections(selected)}
      //     style={{ background: '#8F1F00', marginRight: '10px' }}
      //   >
      //     <DeleteForeverIcon style={iconMedium} />
      //     Delete
      //   </Button>
      // ),
      enableNestedDataAccess: '.',
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        let order = '';
        const columnName = camelcaseToSnakecase(columns[tableState.activeColumn]?.name);
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
      expandableRows: true,
      expandableRowsHeader: false,
      expandableRowsOnClick: true,
      rowsExpanded: rowsExpanded,
      isRowExpandable: () => {
        return true;
      },
      onRowExpansionChange: (_, allRowsExpanded) => {
        setRowsExpanded(allRowsExpanded.slice(-1).map((item) => item.index));
        setShowMore(false);
      },
      renderExpandableRow: (rowData, tableMeta) => {
        const colSpan = rowData.length;
        const meshSyncResourcesMetaData =
          meshSyncResources && meshSyncResources[tableMeta.rowIndex];

        return (
          <TableCell colSpan={colSpan} className={classes.innerTableWrapper}>
            <TableContainer className={classes.innerTableContainer}>
              <Table>
                <TableRow className={classes.noGutter}>
                  <TableCell style={{ padding: '20px 0' }}>
                    <Grid container spacing={1} style={{ textTransform: 'lowercase' }}>
                      <Grid item xs={12} md={12} className={classes.contentContainer}>
                        <Grid container spacing={1}>
                          <Grid
                            item
                            xs={12}
                            md={12}
                            style={{
                              display: 'flex',
                              flexWrap: 'wrap',
                              padding: '0 20px',
                              gap: 30,
                            }}
                            className={classes.contentContainer}
                          >
                            <MeshSyncDataFormatter metadata={meshSyncResourcesMetaData.metadata} />
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
    }),
    [rowsExpanded, showMore, page, pageSize],
  );

  /**
   * fetch connections when the page loads
   */
  useEffect(() => {
    if (!loading) {
      getMeshsyncResources(page, pageSize, search, sortOrder);
    }
  }, [page, pageSize, search, sortOrder]);

  const getMeshsyncResources = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';
    dataFetch(
      `/api/system/meshsync/resources?clusterIds=${clusterIds}&page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
        search,
      )}&order=${encodeURIComponent(sortOrder)}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        const filteredData = res?.resources?.filter((item) => {
          if (selectedFilters.kind === 'All') {
            return true;
          }
          return item.kind === selectedFilters.kind;
        });
        setMeshSyncResources(filteredData);
        setCount(res?.total_count || 0);
        setPageSize(res?.page_size || 0);
        setLoading(false);
        setFilter(filter);
      },
      handleError(ACTION_TYPES.FETCH_MESHSYNC_RESOURCES),
    );
  };

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const filters = {
    kind: {
      name: 'Kind',
      options: [
        { label: 'Deployment', value: 'Deployment' },
        { label: 'Service', value: 'Service' },
        { label: 'Pod', value: 'Pod' },
        { label: 'Namespace', value: 'Namespace' },
        { label: 'StatefulSet', value: 'StatefulSet' },
        { label: 'DaemonSet', value: 'DaemonSet' },
        { label: 'Job', value: 'Job' },
        { label: 'CronJob', value: 'CronJob' },
        { label: 'ReplicaSet', value: 'ReplicaSet' },
        { label: 'ReplicationController', value: 'ReplicationController' },
        { label: 'HorizontalPodAutoscaler', value: 'HorizontalPodAutoscaler' },
        { label: 'Ingress', value: 'Ingress' },
        { label: 'ConfigMap', value: 'ConfigMap' },
        { label: 'Secret', value: 'Secret' },
        { label: 'ServiceAccount', value: 'ServiceAccount' },
        { label: 'PersistentVolume', value: 'PersistentVolume' },
        { label: 'PersistentVolumeClaim', value: 'PersistentVolumeClaim' },
        { label: 'StorageClass', value: 'StorageClass' },
        { label: 'VolumeAttachment', value: 'VolumeAttachment' },
        { label: 'CustomResourceDefinition', value: 'CustomResourceDefinition' },
        { label: 'ClusterRole', value: 'ClusterRole' },
        { label: 'ClusterRoleBinding', value: 'ClusterRoleBinding' },
        { label: 'Role', value: 'Role' },
        { label: 'RoleBinding', value: 'RoleBinding' },
        { label: 'NetworkPolicy', value: 'NetworkPolicy' },
        { label: 'PodSecurityPolicy', value: 'PodSecurityPolicy' },
        { label: 'Node', value: 'Node' },
        { label: 'CustomResource', value: 'CustomResource' },
        { label: 'CustomResourceDefinition', value: 'CustomResourceDefinition' },
        { label: 'Mesh', value: 'Mesh' },
        { label: 'MeshSync', value: 'MeshSync' },
        { label: 'MeshSyncResource', value: 'MeshSyncResource' },
        { label: 'MeshSyncResourceType', value: 'MeshSyncResourceType' },
      ],
    },
  };

  const [selectedFilters, setSelectedFilters] = useState({ kind: 'All' });

  const handleApplyFilter = () => {
    const columnName = Object.keys(selectedFilters)[0];
    const columnValue = selectedFilters[columnName];
    const filter = {
      [columnName]: columnValue === 'All' ? null : [columnValue],
    };
    setFilter(filter);
    getMeshsyncResources(page, pageSize, search, sortOrder, 'filter');
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
      <div className={StyleClass.toolWrapper} style={{ marginBottom: '5px', marginTop: '-30px' }}>
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
            expanded={isSearchExpanded}
            setExpanded={setIsSearchExpanded}
            placeholder="Search connections..."
          />

          <UniversalFilter
            id="ref"
            filters={filters}
            selectedFilters={selectedFilters}
            setSelectedFilters={setSelectedFilters}
            handleApplyFilter={handleApplyFilter}
            conditionForMaxHeight={true}
          />

          <CustomColumnVisibilityControl
            id="ref"
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
  );
}
