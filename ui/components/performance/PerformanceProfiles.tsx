import React, { useEffect, useState, useRef } from 'react';
import Moment from 'react-moment';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import {
  AddCircleIcon as AddIcon,
  Button,
  CustomColumnVisibilityControl,
  CustomTooltip,
  EditIcon,
  IconButton,
  Modal,
  Paper,
  PlayArrowIcon,
  PROMPT_VARIANTS,
  ResponsiveDataTable,
  SearchBar,
  TableCell,
  TableRow,
  Typography,
  useTheme,
} from '@sistent/sistent';
import MesheryPerformanceComponent from './index';
import PerformanceProfileGrid from './PerformanceProfileGrid';
import PerformanceResults from './PerformanceResults';
import _PromptComponent from '../PromptComponent';
import ViewSwitch from '../ViewSwitch';
import { EVENT_TYPES } from '../../lib/event-types';
import { iconMedium } from '../../css/icons.styles';
import {
  useDeletePerformanceProfileMutation,
  useGetPerformanceProfilesQuery,
} from '@/rtk-query/performance-profile';
import { useNotification } from '@/utils/hooks/useNotification';
import { updateVisibleColumns } from '@/utils/responsive-column';
import { useWindowDimensions } from '@/utils/dimension';
import { ConditionalTooltip } from '@/utils/utils';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { isLocalProvider } from '@/utils/provider';
import { ButtonTextWrapper, ProfileContainer, ViewSwitchBUtton } from './style';
import { DefaultTableCell, SortableTableCell } from '../connections/common';
import { useDispatch, useSelector } from 'react-redux';
import { updateProgressAction } from '@/store/slices/mesheryUi';
import type { GetPerformanceProfilesApiResponse } from '@meshery/schemas/mesheryApi';

/**
 * Type Definition for View Type
 * @typedef {"grid" | "table"} TypeView
 */

/**
 * ViewSwitch component renders a switch for toggling between
 * grid and table views
 * @param {{ view: TypeView, changeView: (view: TypeView) => void }} props
 */

type PerformanceProfileItem = GetPerformanceProfilesApiResponse['profiles'][number];
type SelectablePerformanceProfile = PerformanceProfileItem & { runTest?: boolean };

function PerformanceProfile({ handleDelete }) {
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    'grid',
  );
  const modalRef = useRef(null);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('updatedAt desc');
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [testProfiles, setTestProfiles] = useState<PerformanceProfileItem[]>([]);
  const [profileForModal, setProfileForModal] = useState<Partial<SelectablePerformanceProfile>>();
  const { notify } = useNotification();
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const { providerCapabilities } = useSelector((state) => state.ui);
  const dispatch = useDispatch();

  const [deletePerformanceProfile] = useDeletePerformanceProfileMutation();
  const {
    data: performanceProfilesData,
    isFetching: isFetchingProfiles,
    isError: isProfileFetchError,
    error: profileFetchError,
    refetch: refetchProfiles,
  } = useGetPerformanceProfilesQuery({
    page,
    pagesize: pageSize,
    search,
    order: sortOrder,
  });

  useEffect(() => {
    dispatch(updateProgressAction({ showProgress: isFetchingProfiles }));
  }, [dispatch, isFetchingProfiles]);

  useEffect(() => {
    if (!performanceProfilesData) return;

    setCount(performanceProfilesData.totalCount || 0);
    setPageSize(performanceProfilesData.pageSize || pageSize);
    setTestProfiles(performanceProfilesData.profiles || []);
  }, [performanceProfilesData]);

  useEffect(() => {
    if (isProfileFetchError) {
      handleError('Failed to Fetch Profiles')(profileFetchError);
    }
  }, [isProfileFetchError, profileFetchError]);

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
        dispatch(updateProgressAction({ showProgress: false }));
        notify({ message: 'Performance Profile Deleted!', event_type: EVENT_TYPES.SUCCESS });
        refetchProfiles();
      })
      .catch(() => handleError('Failed To Delete Profile'));
  }

  function handleError(msg) {
    return function (error) {
      dispatch(updateProgressAction({ showProgress: false }));
      notify({
        message: `${msg}: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  }

  const [selectedProfile, setSelectedProfile] = useState<SelectablePerformanceProfile>();
  useEffect(() => {
    setProfileForModal(selectedProfile);
  }, [selectedProfile]);

  const searchTimeout = useRef(null);
  const theme = useTheme();

  let colViews = [
    ['name', 'xs'],
    ['endpoints', 'l'],
    ['lastRun', 'l'],
    ['nextRun', 'na'],
    ['updatedAt', 'l'],
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
      name: 'lastRun',
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
      name: 'nextRun',
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
      name: 'updatedAt',
      label: 'Updated At',
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
    sort: !isLocalProvider(providerCapabilities),
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
      name: 'updatedAt',
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
        refetchProfiles();
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
              endpoint={`/api/user/performance/profiles/${testProfiles[rowMeta.rowIndex].id}/results`}
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
            meshName={profileForModal?.serviceMesh}
            url={profileForModal?.endpoints?.[0]}
            qps={profileForModal?.qps}
            loadGenerator={profileForModal?.loadGenerators?.[0]}
            t={profileForModal?.duration}
            c={profileForModal?.concurrentRequest}
            reqBody={profileForModal?.requestBody}
            headers={profileForModal?.requestHeaders}
            cookies={profileForModal?.requestCookies}
            contentType={profileForModal?.contentType}
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

export default PerformanceProfile;
