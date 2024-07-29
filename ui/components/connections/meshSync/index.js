import React, { useEffect, useRef, useState } from 'react';
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
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import {
  CustomColumnVisibilityControl,
  ResponsiveDataTable,
  SearchBar,
  UniversalFilter,
} from '@layer5/sistent';
import useStyles from '../../../assets/styles/general/tool.styles';
import { MeshSyncDataFormatter } from '../metadata';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import { DefaultTableCell, SortableTableCell } from '../common';
import {
  JsonParse,
  camelcaseToSnakecase,
  getColumnValue,
  getVisibilityColums,
} from '../../../utils/utils';
import RegisterConnectionModal from './RegisterConnectionModal';
import classNames from 'classnames';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import { CONNECTION_STATES, MESHSYNC_STATES } from '../../../utils/Enum';
import { updateVisibleColumns } from '../../../utils/responsive-column';
import { useWindowDimensions } from '../../../utils/dimension';
import { FormatId } from '../../DataFormatter';
import {
  useGetMeshSyncResourceKindsQuery,
  useGetMeshSyncResourcesQuery,
} from '@/rtk-query/meshsync';
import { UsesSistent } from '@/components/SistentWrapper';

const ACTION_TYPES = {
  FETCH_MESHSYNC_RESOURCES: {
    name: 'FETCH_MESHSYNC_RESOURCES',
    error_msg: 'Failed to fetch meshsync resources',
  },
};

export default function MeshSyncTable(props) {
  const { classes, updateProgress, selectedK8sContexts, k8sconfig } = props;
  const callbackRef = useRef();
  const [openRegistrationModal, setRegistrationModal] = useState(false);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [rowsExpanded, setRowsExpanded] = useState([]);
  const [selectedKind, setSelectedKind] = useState('');

  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [selectedFilters, setSelectedFilters] = useState({ kind: 'All' });
  const [registerConnection, setRegisterConnection] = useState({
    metadata: {},
    kind: '',
  });
  const StyleClass = useStyles();
  const { width } = useWindowDimensions();

  const icons = {
    [MESHSYNC_STATES.REGISTER]: () => <AssignmentTurnedInIcon />,
    [MESHSYNC_STATES.DISCOVERED]: () => <ExploreIcon />,
  };

  const { notify } = useNotification();

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
    kind: selectedKind,
    clusterIds: JSON.stringify(getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sconfig)),
  });
  if (isError) {
    if (isError) {
      notify({
        message: 'Error fetching MeshSync Resources',
        event_type: EVENT_TYPES.ERROR,
        details: meshSyncError?.data,
      });
    }
  }
  const { data: allKinds } = useGetMeshSyncResourceKindsQuery({
    page: page,
    pagesize: 'all',
    search: search,
    order: sortOrder,
    clusterIds: JSON.stringify(getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sconfig)),
  });
  const availableKinds = allKinds?.kinds || [];

  const meshSyncResources = meshSyncData?.resources || [];

  let colViews = [
    ['metadata.name', 'xs'],
    ['apiVersion', 'xs'],
    ['kind', 'm'],
    ['cluster_id', 'na'],
    ['pattern_resources', 'na'],
    ['metadata.creationTimestamp', 'l'],
    ['status', 'xs'],
    ['metadata', 'na'],
  ];

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
            <>
              <FormControl className={classes.chipFormControl}>
                <Select
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
    // resizableColumns: true,
    serverSide: true,
    selectableRows: 'none',
    count: meshSyncData?.total_count,
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
      // setShowMore(false);
    },
    renderExpandableRow: (rowData) => {
      const colSpan = rowData.length;
      const columnName = 'metadata'; // Name of the column containing the metadata
      const columnIndex = columns.findIndex((column) => column.name === columnName);

      // Access the metadata value using the column index
      const metadata = rowData[columnIndex];

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
                          <MeshSyncDataFormatter metadata={metadata} />
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
  };

  const handleApplyFilter = () => {
    const columnName = Object.keys(selectedFilters)[0];
    const columnValue = selectedFilters[columnName];

    // Check if the selected value is "All"
    const newSelectedKind = columnValue === 'All' ? '' : columnValue;
    setSelectedKind(newSelectedKind);
  };

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  });

  useEffect(() => {
    updateCols(columns);
  }, []);

  return (
    <>
      <div className={StyleClass.toolWrapper} style={{ marginBottom: '5px', marginTop: '-30px' }}>
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
        </UsesSistent>
      </div>
      <UsesSistent>
        <ResponsiveDataTable
          data={meshSyncResources}
          columns={columns}
          options={options}
          className={classes.muiRow}
          tableCols={tableCols}
          updateCols={updateCols}
          columnVisibility={columnVisibility}
        />
      </UsesSistent>
      <RegisterConnectionModal
        handleRegistrationModalClose={handleRegistrationModalClose}
        openRegistrationModal={openRegistrationModal}
        connectionData={registerConnection}
      />
    </>
  );
}
