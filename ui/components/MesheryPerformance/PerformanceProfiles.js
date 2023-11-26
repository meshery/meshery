// @ts-nocheck
import React, { useEffect, useState, useRef } from 'react';
import PromptComponent, { PROMPT_VARIANTS } from '../PromptComponent';
import CloseIcon from '@material-ui/icons/Close';
import PerformanceProfileGrid from './PerformanceProfileGrid';
import dataFetch from '../../lib/data-fetch';
import IconButton from '@material-ui/core/IconButton';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { updateProgress } from '../../lib/store';
import GenericModal from '../GenericModal';
import MesheryPerformanceComponent from './index';
import { Paper, Typography, Button, DialogTitle, TableCell, TableRow } from '@material-ui/core';
import fetchPerformanceProfiles from '../graphql/queries/PerformanceProfilesQuery';
import { withStyles } from '@material-ui/core/styles';
import { iconMedium } from '../../css/icons.styles';
import subscribePerformanceProfiles from '../graphql/subscriptions/PerformanceProfilesSubscription';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import ResponsiveDataTable from '../../utils/data-table';
import Moment from 'react-moment';
import { withSnackbar } from 'notistack';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import PerformanceResults from './PerformanceResults';
import EditIcon from '@material-ui/icons/Edit';
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import ReusableTooltip from '../reusable-tooltip';
import CustomColumnVisibilityControl from '../../utils/custom-column';
import ViewSwitch from '../ViewSwitch';
import SearchBar from '../../utils/custom-search';
import useStyles from '../../assets/styles/general/tool.styles';
import { updateVisibleColumns } from '../../utils/responsive-column';
import { useWindowDimensions } from '../../utils/dimension';

const MESHERY_PERFORMANCE_URL = '/api/user/performance/profiles';
const styles = (theme) => ({
  title: {
    textAlign: 'center',
    minWidth: 400,
    padding: '10px',
    color: '#fff',
    flexGrow: 1,
  },
  btnText: {
    display: 'block',
    '@media (max-width: 1450px)': {
      display: 'none',
    },
  },
  dialogHeader: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.headerColor
        : theme.palette.secondary.mainBackground,
    display: 'flex',
    alignItems: 'center',
  },
  addButton: {
    width: 'fit-content',
    alignSelf: 'flex-start',
  },
  viewSwitchButton: {
    justifySelf: 'flex-end',
    marginLeft: 'auto',
    paddingLeft: '1rem',
    display: 'flex',
  },
  pageContainer: {
    padding: '0.5rem',
  },
  noProfileContainer: {
    padding: '2rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
  },
  noProfilePaper: {
    padding: '0.5rem',
  },
  noProfileText: {
    fontSize: '1.5rem',
    marginBottom: '2rem',
  },
  addProfileModal: {
    margin: 'auto',
    maxWidth: '90%',
    outline: 'none',
  },
  addIcon: {
    paddingRight: '0.5',
  },
  grid: { padding: theme.spacing(2) },
  tableHeader: {
    fontWeight: 'bolder',
    fontSize: 18,
  },
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
});
/**
 * Type Definition for View Type
 * @typedef {"grid" | "table"} TypeView
 */

/**
 * ViewSwitch component renders a switch for toggling between
 * grid and table views
 * @param {{ view: TypeView, changeView: (view: TypeView) => void }} props
 */

function PerformanceProfile({ updateProgress, classes, user, handleDelete }) {
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ('grid'),
  );
  const modalRef = useRef(null);
  const StyleClass = useStyles();

  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [testProfiles, setTestProfiles] = useState([]);
  const [profileForModal, setProfileForModal] = useState();
  const { notify } = useNotification();
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  // const [loading, setLoading] = useState(false);
  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles(page, pageSize, search, sortOrder);
    const subscription = subscribePerformanceProfiles(
      (res) => {
        let result = res?.subscribePerfProfiles;
        if (typeof result !== 'undefined') {
          if (result) {
            setCount(result.total_count || 0);
            setPageSize(result.page_size || 0);
            setTestProfiles(result.profiles || []);
            setPage(result.page || 0);
          }
        }
      },
      {
        selector: {
          pageSize: `${pageSize}`,
          page: `${page}`,
          search: `${encodeURIComponent(search)}`,
          order: `${encodeURIComponent(sortOrder)}`,
        },
      },
    );
    return () => {
      subscription.dispose();
    };
  }, [page, pageSize, search, sortOrder]);

  /**
   * fetchTestProfiles constructs the queries based on the parameters given
   * and fetches the performance profiles
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  function fetchTestProfiles(page, pageSize, search, sortOrder) {
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';
    // const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
    //   sortOrder
    // )}`;

    updateProgress({ showProgress: true });
    fetchPerformanceProfiles({
      selector: {
        pageSize: `${pageSize}`,
        page: `${page}`,
        search: `${encodeURIComponent(search)}`,
        order: `${encodeURIComponent(sortOrder)}`,
      },
    }).subscribe({
      next: (res) => {
        let result = res?.getPerformanceProfiles;
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          if (result) {
            setCount(result.total_count || 0);
            setPageSize(result.page_size || 0);
            setTestProfiles(result.profiles || []);
            // setPage(result.page || 0);
          }
        }
      },
      error: handleError('Failed to Fetch Profiles'),
    });
  }

  async function showModal(count) {
    let response = await modalRef.current.show({
      title: `Delete ${count ? count : ''} Performance Profile${count > 1 ? 's' : ''}?`,
      subtitle: `Are you sure you want to delete ${count > 1 ? 'these' : 'this'} ${
        count ? count : ''
      } performance profile${count > 1 ? 's' : ''}?`,
      variant: PROMPT_VARIANTS.DANGER,
      options: ['Yes', 'No'],
    });
    return response;
  }

  function deleteProfile(id) {
    dataFetch(
      `${MESHERY_PERFORMANCE_URL}/${id}`,
      {
        method: 'DELETE',
        credentials: 'include',
      },
      () => {
        updateProgress({ showProgress: false });
        notify({ message: 'Performance Profile Deleted!', event_type: EVENT_TYPES.SUCCESS });
        fetchTestProfiles(page, pageSize, search, sortOrder);
      },
      handleError('Failed To Delete Profile'),
    );
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress: false });
      notify({
        message: `${msg}: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  }

  const [selectedProfile, setSelectedProfile] = useState();
  useEffect(() => {
    setProfileForModal(selectedProfile);
  }, [selectedProfile]);

  const searchTimeout = useRef(null);

  let colViews = [
    ['name', 'xs'],
    ['endpoints', 'l'],
    ['last_run', 'l'],
    ['next_run', 'xl'],
    ['updated_at', 'l'],
    ['Actions', 'xs'],
  ];

  const columns = [
    {
      name: 'name',
      label: 'Profile',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'endpoints',
      label: 'Endpoints',
      options: {
        filter: false,
        sort: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'last_run',
      label: 'Last Run',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'next_run',
      label: 'Next Run',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'updated_at',
      label: 'Updated On',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'Actions',
      label: 'Actions',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <div>
              <ReusableTooltip title="Edit">
                <IconButton
                  style={iconMedium}
                  onClick={(ev) => {
                    ev.stopPropagation();
                    setSelectedProfile(testProfiles[tableMeta.rowIndex]);
                  }}
                  aria-label="edit"
                  // @ts-ignore
                  color="rgba(0, 0, 0, 0.54)"
                >
                  <EditIcon style={iconMedium} />
                </IconButton>
              </ReusableTooltip>

              <ReusableTooltip title="Run test">
                <IconButton
                  style={iconMedium}
                  onClick={(ev) => {
                    ev.stopPropagation();
                    setSelectedProfile({ ...testProfiles[tableMeta.rowIndex], runTest: true });
                  }}
                  aria-label="run"
                  // @ts-ignore
                  color="rgba(0, 0, 0, 0.54)"
                >
                  <PlayArrowIcon style={iconMedium} />
                </IconButton>
              </ReusableTooltip>
            </div>
          );
        },
      },
    },
  ];

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(' ')[0]) {
      columns[idx].options.sortDirection = sortOrder.split(' ')[1];
    }
  });

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

  const options = {
    filter: false,
    search: false,
    viewColumns: false,
    sort: !(user && user.user_id === 'meshery'),
    // search : !(user && user.user_id === "meshery"),
    filterType: 'textField',
    responsive: 'standard',
    resizableColumns: true,
    serverSide: true,
    selectableRows: true,
    count,
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    print: false,
    download: false,
    textLabels: {
      selectedRows: {
        text: 'profile(s) selected',
      },
    },

    onRowsDelete: async function handleDeleteRow(row) {
      let response = await showModal(Object.keys(row.lookup).length);
      console.log(response);
      if (response === 'Yes') {
        const pids = Object.keys(row.lookup).map((idx) => testProfiles[idx]?.id);
        pids.forEach((pid) => handleDelete(pid));
      }
      if (response === 'No') {
        fetchTestProfiles(page, pageSize, search, sortOrder);
      }
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          if (searchTimeout.current) clearTimeout(searchTimeout.current);
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) setSearch(tableState.searchText);
          }, 500);
          break;
        case 'sort':
          if (sortInfo.length === 2) {
            if (sortInfo[1] === 'ascending') {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) setSortOrder(order);
          break;
      }
    },
    expandableRows: true,
    renderExpandableRow: function ExpandableRow(rowData, rowMeta) {
      const colSpan = rowData.length;
      return (
        <TableRow>
          <TableCell />
          <TableCell colSpan={colSpan}>
            <PerformanceResults
              // @ts-ignore
              CustomHeader={<Typography variant="h6">Test Results</Typography>}
              // @ts-ignore
              endpoint={`/api/user/performance/profiles/${
                testProfiles[rowMeta.rowIndex].id
              }/results`}
              // @ts-ignore
              elevation={0}
            />
          </TableCell>
        </TableRow>
      );
    },
  };

  return (
    <>
      <div className={classes.pageContainer}>
        <div className={StyleClass.toolWrapper}>
          {width < 550 && isSearchExpanded ? null : (
            <>
              {(testProfiles.length > 0 || viewType == 'table') && (
                <div className={classes.addButton}>
                  <Button
                    aria-label="Add Performance Profile"
                    variant="contained"
                    color="primary"
                    size="large"
                    onClick={() => setProfileForModal({})}
                  >
                    <AddIcon style={iconMedium} className={classes.addIcon} />
                    <span className={classes.btnText}> Add Performance Profile </span>
                  </Button>
                </div>
              )}
            </>
          )}
          <div className={classes.viewSwitchButton}>
            <SearchBar
              onSearch={(value) => {
                setSearch(value);
                fetchTestProfiles(page, pageSize, value, sortOrder);
              }}
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
              placeholder="Search Profiles..."
            />
            <CustomColumnVisibilityControl
              classes={classes}
              columns={columns}
              customToolsProps={{ columnVisibility, setColumnVisibility }}
            />
            <ViewSwitch view={viewType} changeView={setViewType} />
          </div>
        </div>
        {viewType === 'grid' ? (
          <PerformanceProfileGrid
            profiles={testProfiles}
            deleteHandler={deleteProfile}
            setProfileForModal={setProfileForModal}
            pages={Math.ceil(count / pageSize)}
            setPage={setPage}
          />
        ) : (
          <ResponsiveDataTable
            data={testProfiles}
            columns={columns}
            // @ts-ignore
            options={options}
            tableCols={tableCols}
            updateCols={updateCols}
            columnVisibility={columnVisibility}
          />
        )}
        {testProfiles.length === 0 && viewType === 'grid' && (
          <Paper className={classes.noProfilePaper}>
            <div className={classes.noProfileContainer}>
              <Typography className={classes.noProfileText} align="center" color="textSecondary">
                No Performance Profiles Found
              </Typography>
              <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => setProfileForModal({})}
              >
                <Typography className="addIcon">Add Performance Profile</Typography>
              </Button>
            </div>
          </Paper>
        )}
        <GenericModal
          open={!!profileForModal}
          Content={
            <Paper className={classes.addProfileModal}>
              <div className={classes.dialogHeader}>
                <DialogTitle className={classes.title}>Performance Profile Wizard</DialogTitle>
                <IconButton
                  aria-label="close"
                  style={{ color: 'white' }}
                  onClick={() => setProfileForModal(undefined)}
                >
                  <CloseIcon />
                </IconButton>
              </div>

              <MesheryPerformanceComponent
                loadAsPerformanceProfile
                performanceProfileID={profileForModal?.id}
                profileName={profileForModal?.name}
                meshName={profileForModal?.service_mesh}
                url={profileForModal?.endpoints?.[0]}
                qps={profileForModal?.qps}
                loadGenerator={profileForModal?.load_generators?.[0]}
                t={profileForModal?.duration}
                c={profileForModal?.concurrent_request}
                reqBody={profileForModal?.request_body}
                headers={profileForModal?.request_headers}
                cookies={profileForModal?.request_cookies}
                contentType={profileForModal?.content_type}
                runTestOnMount={!!profileForModal?.runTest}
                metadata={profileForModal?.metadata}
              />
            </Paper>
          }
          handleClose={() => {
            setProfileForModal(undefined);
          }}
        />
      </div>

      <PromptComponent ref={modalRef} />
    </>
  );
}
const mapStateToProps = (state) => {
  return { user: state.get('user')?.toObject() };
};
const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withSnackbar(PerformanceProfile)),
);
