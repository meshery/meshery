// @ts-nocheck
import React, { useEffect, useState, useRef } from 'react';
import Moment from 'react-moment';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { TableCell, TableRow } from '@mui/material';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import EditIcon from '@mui/icons-material/Edit';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import {
  CustomColumnVisibilityControl,
  Modal,
  PROMPT_VARIANTS,
  ResponsiveDataTable,
  SearchBar,
  Button,
  Paper,
  Typography,
  IconButton,
  useTheme,
  CustomTooltip,
} from '@layer5/sistent';
import MesheryPerformanceComponent from './index';
import PerformanceProfileGrid from './PerformanceProfileGrid';
import PerformanceResults from './PerformanceResults';
import _PromptComponent from '../PromptComponent';
import ViewSwitch from '../ViewSwitch';

import { updateProgress } from '../../lib/store';
import { EVENT_TYPES } from '../../lib/event-types';
import fetchPerformanceProfiles from '../graphql/queries/PerformanceProfilesQuery';
import subscribePerformanceProfiles from '../graphql/subscriptions/PerformanceProfilesSubscription';
import { iconMedium } from '../../css/icons.styles';
import { useDeletePerformanceProfileMutation } from '@/rtk-query/performance-profile';
import { useNotification } from '@/utils/hooks/useNotification';
import { updateVisibleColumns } from '@/utils/responsive-column';
import { useWindowDimensions } from '@/utils/dimension';
import { ConditionalTooltip } from '@/utils/utils';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { ButtonTextWrapper, ProfileContainer, ViewSwitchBUtton } from './style';
import { DefaultTableCell, SortableTableCell } from '../connections/common';

/**
 * Type Definition for View Type
 * @typedef {"grid" | "table"} TypeView
 */

/**
 * ViewSwitch component renders a switch for toggling between
 * grid and table views
 * @param {{ view: TypeView, changeView: (view: TypeView) => void }} props
 */

function PerformanceProfile({ updateProgress, user, handleDelete }) {
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ('grid'),
  );
  const modalRef = useRef(null);
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

  const [deletePerformanceProfile] = useDeletePerformanceProfileMutation();
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
          order: `${sortOrder}`,
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

    updateProgress({ showProgress: true });
    fetchPerformanceProfiles({
      selector: {
        pageSize: `${pageSize}`,
        page: `${page}`,
        search: `${encodeURIComponent(search)}`,
        order: `${sortOrder}`,
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
      primaryOption: 'DELETE',
    });
    return response;
  }

  function deleteProfile(id) {
    deletePerformanceProfile({ id: id })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({ message: 'Performance Profile Deleted!', event_type: EVENT_TYPES.SUCCESS });
        fetchTestProfiles(page, pageSize, search, sortOrder);
      })
      .catch(() => handleError('Failed To Delete Profile'));
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
  const theme = useTheme();

  let colViews = [
    ['name', 'xs'],
    ['endpoints', 'l'],
    ['last_run', 'l'],
    ['next_run', 'na'],
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
      name: 'endpoints',
      label: 'Endpoints',
      options: {
        filter: false,
        sort: true,
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
        customBodyRender: (value) => <ConditionalTooltip value={value} maxLength={20} />,
      },
    },
    {
      name: 'last_run',
      label: 'Last Run',
      options: {
        filter: false,
        sort: true,
        searchable: true,
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
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <div style={{ display: 'flex', gap: '1rem' }}>
              <CustomTooltip title="Edit">
                <div>
                  <IconButton
                    style={iconMedium}
                    onClick={(ev) => {
                      ev.stopPropagation();
                      setSelectedProfile(testProfiles[tableMeta.rowIndex]);
                    }}
                    aria-label="edit"
                    disabled={
                      !CAN(keys.EDIT_PERFORMANCE_TEST.action, keys.EDIT_PERFORMANCE_TEST.subject)
                    }
                  >
                    <EditIcon
                      style={{
                        fill: theme.palette.icon.secondary,
                        ...iconMedium,
                      }}
                    />
                  </IconButton>
                </div>
              </CustomTooltip>

              <CustomTooltip title="Run test">
                <div>
                  <IconButton
                    style={iconMedium}
                    onClick={(ev) => {
                      ev.stopPropagation();
                      setSelectedProfile({ ...testProfiles[tableMeta.rowIndex], runTest: true });
                    }}
                    aria-label="run"
                    disabled={!CAN(keys.RUN_TEST.action, keys.RUN_TEST.subject)}
                  >
                    <PlayArrowIcon
                      style={{
                        fill: theme.palette.icon.secondary,
                        ...iconMedium,
                      }}
                    />
                  </IconButton>
                </div>
              </CustomTooltip>
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
    selectableRows: 'multiple',
    count,
    rowsPerPage: pageSize,
    fixedHeader: true,
    page,
    print: false,
    download: false,
    sortOrder: {
      name: 'updated_at',
      direction: 'desc',
    },
    textLabels: {
      selectedRows: {
        text: 'profile(s) selected',
      },
    },

    onRowsDelete: async function handleDeleteRow(row) {
      let response = await showModal(Object.keys(row.lookup).length);
      if (response === 'DELETE') {
        const pids = Object.keys(row.lookup).map((idx) => testProfiles[idx]?.id);
        pids.forEach((pid) => handleDelete(pid));
      } else {
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
      <div style={{ padding: '0.5rem' }}>
        <ToolWrapper>
          {width < 550 && isSearchExpanded ? null : (
            <>
              {(testProfiles.length > 0 || viewType == 'table') && (
                <div style={{ width: 'fit-content', alignSelf: 'flex-start' }}>
                  <Button
                    aria-label="Add Performance Profile"
                    variant="contained"
                    color="primary"
                    size="large"
                    onClick={() => setProfileForModal({})}
                    disabled={
                      !CAN(
                        keys.ADD_PERFORMANCE_PROFILE.action,
                        keys.ADD_PERFORMANCE_PROFILE.subject,
                      )
                    }
                  >
                    <AddIcon style={{ paddingRight: '0.5', ...iconMedium }} />
                    <ButtonTextWrapper> Add Performance Profile </ButtonTextWrapper>
                  </Button>
                </div>
              )}
            </>
          )}
          <ViewSwitchBUtton>
            <SearchBar
              onSearch={(value) => {
                setSearch(value);
                fetchTestProfiles(page, pageSize, value, sortOrder);
              }}
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
              placeholder="Search Profiles..."
            />
            {viewType === 'table' && (
              <CustomColumnVisibilityControl
                id="ref"
                columns={columns}
                customToolsProps={{ columnVisibility, setColumnVisibility }}
              />
            )}
            <ViewSwitch view={viewType} changeView={setViewType} />
          </ViewSwitchBUtton>
        </ToolWrapper>

        {viewType === 'grid' ? (
          <PerformanceProfileGrid
            profiles={testProfiles}
            deleteHandler={deleteProfile}
            setProfileForModal={setProfileForModal}
            pages={Math.ceil(count / pageSize)}
            setPage={setPage}
            testHandler={setSelectedProfile}
          />
        ) : (
          <ResponsiveDataTable
            data={testProfiles}
            columns={columns}
            options={options}
            tableCols={tableCols}
            updateCols={updateCols}
            columnVisibility={columnVisibility}
          />
        )}
        {testProfiles.length === 0 && viewType === 'grid' && (
          <Paper sx={{ padding: '0.5rem' }}>
            <ProfileContainer>
              <Typography sx={{ fontSize: '1.5rem', marginBottom: '2rem' }} align="center">
                No Performance Profiles Found
              </Typography>
              <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => setProfileForModal({})}
                disabled={
                  !CAN(keys.ADD_PERFORMANCE_PROFILE.action, keys.ADD_PERFORMANCE_PROFILE.subject)
                }
              >
                <Typography className="addIcon">Add Performance Profile</Typography>
              </Button>
            </ProfileContainer>
          </Paper>
        )}
        <Modal
          open={!!profileForModal}
          title="Performance Profile Wizard"
          closeModal={() => setProfileForModal(undefined)}
          maxWidth="md"
        >
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
            closeModal={() => setProfileForModal(undefined)}
          />
        </Modal>
      </div>

      <_PromptComponent ref={modalRef} />
    </>
  );
}
const mapStateToProps = (state) => {
  return { user: state.get('user')?.toObject() };
};
const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(PerformanceProfile);
