import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  TableCell,
  Tooltip,
  TableContainer,
  Table,
  Grid,
  TableRow,
  FormControl,
  Select,
  MenuItem,
  Chip,
} from '@material-ui/core';
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
import { camelcaseToSnakecase, getColumnValue } from '../../../utils/utils';
import RegisterConnectionModal from './RegisterConnectionModal';
import classNames from 'classnames';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import { MESHSYNC_STATES } from '../../../utils/Enum';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

export default function MeshSyncTable(props) {
  const { classes, updateProgress, selectedK8sContexts, k8sconfig } = props;
  const callbackRef = useRef();
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [search, setSearch] = useState('');
  const [meshSyncResources, setMeshSyncResources] = useState([]);
  const [sortOrder, setSortOrder] = useState('');
  const [showMore, setShowMore] = useState(false);
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [registerConnection, setRegisterConnection] = useState({
    metadata: {},
    kind: '',
  });
  const StyleClass = useStyles();

  const icons = {
    [MESHSYNC_STATES.REGISTER]: () => <AssignmentTurnedInIcon />,
    [MESHSYNC_STATES.DISCOVERED]: () => <ExploreIcon />,
  };

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
    {
      name: 'status',
      label: 'Status',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          const metadata = getColumnValue(
            tableMeta.rowData,
            'component_metadata',
            columns,
          )?.metadata;
          const DISCOVERED = {
            DISCOVERED: MESHSYNC_STATES.DISCOVERED,
          };
          const meshSyncStates =
            metadata && JSON.parse(metadata)?.capabilities?.connection === true
              ? MESHSYNC_STATES
              : DISCOVERED;
          const disabled =
            metadata && JSON.parse(metadata)?.capabilities?.connection === true ? false : true;
          return (
            <>
              <FormControl className={classes.chipFormControl}>
                <Select
                  labelId="demo-simple-select-label"
                  id="demo-simple-select"
                  defaultValue={MESHSYNC_STATES.DISCOVERED}
                  disabled={disabled}
                  onClick={(e) => e.stopPropagation()}
                  onChange={() => {
                    callbackRef?.current?.(tableMeta);
                    const componentMetadata = getColumnValue(
                      tableMeta.rowData,
                      'component_metadata',
                      columns,
                    );
                    const connectionKind = getColumnValue(tableMeta.rowData, 'kind', columns);
                    setRegisterConnection({
                      metadata: JSON.parse(componentMetadata.metadata),
                      kind: connectionKind,
                    });
                  }}
                  className={classes.statusSelect}
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
                      disabled={meshSyncStates[s] === value ? true : false}
                      value={meshSyncStates[s]}
                      key={meshSyncStates[s]}
                      style={{ padding: '0' }}
                    >
                      <Chip
                        className={classNames(classes.statusChip, classes[meshSyncStates[s]])}
                        avatar={icons[meshSyncStates[s]]()}
                        label={meshSyncStates[s]}
                      />
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </>
          );
        },
      },
    },
    {
      name: 'component_metadata',
      label: 'Component Metadata',
      options: {
        sort: false,
        sortThirdClickReset: false,
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
        setMeshSyncResources(res?.resources || []);
        setCount(res?.total_count || 0);
        setPageSize(res?.page_size || 0);
        setLoading(false);
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

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = col.options?.display !== false;
    });
    return initialVisibility;
  });

  const handleOpenRegisterModal = (callback) => {
    callbackRef.current = callback;
  };

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
      <RegisterConnectionModal
        handleOpen={handleOpenRegisterModal}
        connectionData={registerConnection}
      />
    </>
  );
}
