import { NoSsr } from '@mui/material';
import { connect } from 'react-redux';
import { withRouter } from 'next/router';
import {
  AssignmentModal,
  Box,
  CustomColumnVisibilityControl,
  Pagination,
  PaginationItem,
  ResponsiveDataTable,
  TeamsIcon,
  useDesignAssignment,
  useEnvironmentAssignment,
  useTeamAssignment,
  useViewAssignment,
  ViewIcon,
  WorkspaceIcon,
  Modal as SisitentModal,
  createAndEditWorkspaceSchema,
  createAndEditWorkspaceUiSchema,
  Button,
  Grid,
  Typography,
  DeleteIcon,
  SearchBar,
  useTheme,
  PROMPT_VARIANTS,
  L5EditIcon,
  L5DeleteIcon,
  OutlinedPatternIcon,
} from '@layer5/sistent';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { EmptyState } from '../General';
import AddIconCircleBorder from '../../../assets/icons/AddIconCircleBorder';
import { useEffect, useRef, useState } from 'react';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignEnvironmentToWorkspaceMutation,
  useAssignTeamToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useCreateWorkspaceMutation,
  useDeleteWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetTeamsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useGetWorkspacesQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useUnassignTeamFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
  useUpdateWorkspaceMutation,
} from '../../../rtk-query/workspace';
import { updateProgress } from '../../../lib/store';
import { useNotification } from '../../../utils/hooks/useNotification';
import { RJSFModalWrapper } from '../../Modal';
import _PromptComponent from '../../PromptComponent';
import { debounce } from 'lodash';
import { EVENT_TYPES } from '../../../lib/event-types';
import EnvironmentIcon from '../../../assets/icons/Environment';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import DefaultError from '@/components/General/error-404/index';
import { UsesSistent } from '@/components/SistentWrapper';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import WorkSpaceDataTable from './workspace-table';
import { useWindowDimensions } from '@/utils/dimension';
import { updateVisibleColumns } from '@/utils/responsive-column';
import ViewSwitch from '@/components/ViewSwitch';
import { TableIconsContainer, IconWrapper, CreateButtonWrapper, BulkActionWrapper } from './styles';
import MesheryWorkspaceCard from './MesheryWorkspaceCard';

export const WORKSPACE_ACTION_TYPES = {
  CREATE: 'create',
  EDIT: 'edit',
};

const Workspaces = ({ organization }) => {
  const theme = useTheme();
  const [workspaceModal, setWorkspaceModal] = useState({
    open: false,
    schema: {},
  });
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [sortOrder, setSortOrder] = useState('');
  const [search, setSearch] = useState('');

  const [orgId, setOrgId] = useState('');
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [actionType, setActionType] = useState('');
  const [initialData, setInitialData] = useState({});
  const [editWorkspaceId, setEditWorkspaceId] = useState('');
  const [viewType, setViewType] = useState('grid');
  const [environmentAssignWorkspace, setEnvironmentAssignWorkspace] = useState({});
  const [teamAssignWorkspace, setTeamAssignWorkspace] = useState({});
  const [designAssignWorkspace, setDesignAssignWorkspace] = useState({});
  const [viewsAssignWorkspace, setViewsAssignWorkspace] = useState({});
  const [selectedWorkspaces, setSelectedWorkspaces] = useState([]);

  const ref = useRef(null);
  const bulkDeleteRef = useRef(null);
  const { notify } = useNotification();

  const {
    data: workspacesData,
    isError: isWorkspacesError,
    error: workspacesError,
  } = useGetWorkspacesQuery(
    {
      page: page,
      pagesize: pageSize,
      search: search,
      order: sortOrder,
      orgId: orgId,
    },
    {
      skip: !orgId ? true : false,
    },
  );

  const [createWorkspace] = useCreateWorkspaceMutation();

  const [updateWorkspace] = useUpdateWorkspaceMutation();

  const [deleteWorkspace] = useDeleteWorkspaceMutation();

  const workspaces = workspacesData?.workspaces ? workspacesData.workspaces : [];
  const handleCreateWorkspace = ({ organization, name, description }) => {
    createWorkspace({
      workspacePayload: {
        name: name,
        description: description,
        organization_id: organization,
      },
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace "${name}" created `))
      .catch((error) => handleError(`Workspace Create Error: ${error?.data}`));
    handleWorkspaceModalClose();
  };

  const handleEditWorkspace = ({ organization, name, description }) => {
    updateWorkspace({
      workspaceId: editWorkspaceId,
      workspacePayload: {
        name: name,
        description: description,
        organization_id: organization,
      },
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace "${name}" updated`))
      .catch((error) => handleError(`Workspace Update Error: ${error?.data}`));
    handleWorkspaceModalClose();
  };

  const handleDeleteWorkspace = (id, name) => {
    deleteWorkspace({
      workspaceId: id,
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace "${name}" deleted`))
      .catch((error) => handleError(`Workspace Delete Error: ${error?.data}`));
  };

  useEffect(() => {
    if (isWorkspacesError) {
      handleError(`Workspaces Fetching Error: ${workspacesError?.data}`);
    }
  }, [isWorkspacesError, workspacesError, handleError]);

  useEffect(() => {
    setOrgId(organization?.id);
  }, [organization]);

  const fetchSchema = () => {
    const updatedSchema = {
      schema: createAndEditWorkspaceSchema,
      uiSchema: createAndEditWorkspaceUiSchema,
    };
    updatedSchema.schema?.properties?.organization &&
      ((updatedSchema.schema = {
        ...updatedSchema.schema,
        properties: {
          ...updatedSchema.schema.properties,
          organization: {
            ...updatedSchema.schema.properties.organization,
            enum: [organization?.id],
            enumNames: [organization?.name],
          },
        },
      }),
      (updatedSchema.uiSchema = {
        ...updatedSchema.uiSchema,
        organization: {
          ...updatedSchema.uiSchema.organization,
          ['ui:widget']: 'hidden',
        },
      }));
    setWorkspaceModal({
      open: true,
      schema: updatedSchema,
    });
  };

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleSuccess = (msg) => {
    updateProgress({ showProgress: false });
    notify({
      message: msg,
      event_type: EVENT_TYPES.SUCCESS,
    });
  };

  const handleWorkspaceModalOpen = (e, actionType, workspaceObject) => {
    e.stopPropagation();
    if (actionType === WORKSPACE_ACTION_TYPES.EDIT) {
      setActionType(WORKSPACE_ACTION_TYPES.EDIT);
      setInitialData({
        name: workspaceObject.name,
        description: workspaceObject.description,
        organization: workspaceObject.organization_id,
      });
      setEditWorkspaceId(workspaceObject.id);
    } else {
      setActionType(WORKSPACE_ACTION_TYPES.CREATE);
      setInitialData({
        name: undefined,
        description: '',
        organization: orgId,
      });
      setEditWorkspaceId('');
    }
    fetchSchema();
  };

  const handleWorkspaceModalClose = () => {
    setWorkspaceModal({
      open: false,
      schema: {},
    });
    setActionType('');
  };

  const handleBulkDeleteWorkspace = () => {
    selectedWorkspaces.map((workspaceId) => {
      handleDeleteWorkspace(
        workspaceId,
        workspaces.find((workspace) => workspace.name === name),
      );
    });
    setSelectedWorkspaces([]);
  };

  const handleBulkSelect = (e, id) => {
    const isChecked = e.target.checked;
    if (isChecked) {
      setSelectedWorkspaces([...selectedWorkspaces, id]);
    } else {
      const newSelectedEnv = selectedWorkspaces.filter((env) => env !== id);
      setSelectedWorkspaces(newSelectedEnv);
    }
  };

  const handleDeleteWorkspaceConfirm = async (e, workspace) => {
    e.stopPropagation();
    let response = await ref.current.show({
      title: `Delete workspace ?`,
      subtitle: deleteWorkspaceModalContent(workspace.name),
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
    });
    if (response === 'DELETE') {
      handleDeleteWorkspace(workspace.id, workspace.name);
    }
  };

  const deleteWorkspaceModalContent = (workspace) => (
    <>
      <p>Are you sure you want to delete this workspace? (This action is irreversible)</p>
      <p>
        Workspace Name:
        <i>
          <b>{workspace}</b>
        </i>
      </p>
    </>
  );

  const handleBulkDeleteWorkspaceConfirm = async (e) => {
    e.stopPropagation();
    let response = await bulkDeleteRef.current.show({
      title: `Delete ${selectedWorkspaces.length} workspaces ?`,
      subtitle: deleteBlukWorkspaceModalContent(),
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
    });
    if (response === 'DELETE') {
      handleBulkDeleteWorkspace();
    }
  };

  const deleteBlukWorkspaceModalContent = () => (
    <p>Are you sure you want to delete these workspaces? (This action is irreversible)</p>
  );
  const environmentAssignment = useEnvironmentAssignment({
    workspaceId: environmentAssignWorkspace?.id,
    useAssignEnvironmentToWorkspaceMutation: useAssignEnvironmentToWorkspaceMutation,
    useGetEnvironmentsOfWorkspaceQuery: useGetEnvironmentsOfWorkspaceQuery,
    useUnassignEnvironmentFromWorkspaceMutation: useUnassignEnvironmentFromWorkspaceMutation,
    isEnvironmentsVisible: CAN(keys.VIEW_ENVIRONMENTS.action, keys.VIEW_ENVIRONMENTS.subject),
  });

  const teamAssignment = useTeamAssignment({
    workspaceId: teamAssignWorkspace?.id,
    useAssignTeamToWorkspaceMutation: useAssignTeamToWorkspaceMutation,
    useGetTeamsOfWorkspaceQuery: useGetTeamsOfWorkspaceQuery,
    useUnassignTeamFromWorkspaceMutation: useUnassignTeamFromWorkspaceMutation,
    isTeamsVisible: CAN(keys.VIEW_TEAMS.action, keys.VIEW_TEAMS.subject),
  });

  const designAssignment = useDesignAssignment({
    workspaceId: designAssignWorkspace?.id,
    useAssignDesignToWorkspaceMutation: useAssignDesignToWorkspaceMutation,
    useGetDesignsOfWorkspaceQuery: useGetDesignsOfWorkspaceQuery,
    useUnassignDesignFromWorkspaceMutation: useUnassignDesignFromWorkspaceMutation,
    isDesignsVisible: CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject),
  });

  const viewAssignment = useViewAssignment({
    workspaceId: viewsAssignWorkspace?.id,
    useGetViewsOfWorkspaceQuery: useGetViewsOfWorkspaceQuery,
    useAssignViewToWorkspaceMutation: useAssignViewToWorkspaceMutation,
    useUnassignViewFromWorkspaceMutation: useUnassignViewFromWorkspaceMutation,
    isViewsVisible: CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject),
  });

  const handleAssignEnvironmentModalOpen = (e, workspace) => {
    e.stopPropagation();
    setEnvironmentAssignWorkspace(workspace);
    environmentAssignment.handleAssignModal();
  };

  const handleAssignTeamModalOpen = (e, workspace) => {
    e.stopPropagation();
    teamAssignment.handleAssignModal();
    setTeamAssignWorkspace(workspace);
  };

  const handleAssignDesignModalOpen = (e, workspace) => {
    e.stopPropagation();
    setDesignAssignWorkspace(workspace);
    setViewsAssignWorkspace(workspace);
    designAssignment.handleAssignModal();
    viewAssignment.handleAssignModal();
  };

  const isDesignActivity = designAssignment?.isActivityOccurred(designAssignment?.assignedItems);
  const isViewActivity = viewAssignment?.isActivityOccurred(viewAssignment?.assignedItems);
  const handleAssignments = () => {
    if (isDesignActivity) {
      designAssignment.handleAssign();
    }
    if (isViewActivity) {
      viewAssignment.handleAssign();
    }
  };

  let colViews = [
    ['id', 'na'],
    ['name', 'xs'],
    ['description', 'm'],
    ['owner', 'xs'],
    ['actions', 'xs'],
    ['updated_at', 'l'],
    ['created_at', 'na'],
  ];

  const columns = [
    {
      name: 'id',
      label: 'ID',
      filter: false,
    },
    {
      name: 'name',
      label: 'Name',
      options: {
        filter: false,
        sort: true,
        searchable: true,
      },
    },
    {
      name: 'description',
      label: 'Description',
    },
    {
      name: 'owner',
      label: 'Owner',
      options: {
        filter: false,
        sort: true,
        searchable: true,
      },
    },
    {
      name: 'created_at',
      label: 'Created At',
      options: {
        filter: false,
        sort: true,
        searchable: true,
      },
    },
    {
      name: 'updated_at',
      label: 'Updated At',
      options: {
        filter: false,
        sort: true,
        searchable: true,
      },
    },
    {
      name: 'actions',
      label: 'Actions',
      options: {
        sort: false,
        customBodyRender: (value, tableMeta) => {
          return (
            <UsesSistent>
              <TableIconsContainer>
                <IconWrapper>
                  {
                    <>
                      <L5EditIcon
                        title="Edit Workspace"
                        key={`edit_role-${tableMeta.rowIndex}`}
                        disabled={!CAN(keys.EDIT_WORKSPACE.action, keys.EDIT_WORKSPACE.subject)}
                        onClick={(e) =>
                          handleWorkspaceModalOpen(
                            e,
                            WORKSPACE_ACTION_TYPES.EDIT,
                            workspacesData.workspaces[tableMeta.rowIndex],
                          )
                        }
                      />

                      <L5DeleteIcon
                        title="Delete Workspace"
                        key={`delete_role-${tableMeta.rowIndex}`}
                        disabled={!CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject)}
                        onClick={(e) =>
                          handleDeleteWorkspaceConfirm(
                            e,
                            workspacesData.workspaces[tableMeta.rowIndex],
                          )
                        }
                      />
                    </>
                  }
                </IconWrapper>
              </TableIconsContainer>
            </UsesSistent>
          );
        },
      },
    },
  ];
  const [tableCols, updateCols] = useState(columns);
  const { width } = useWindowDimensions();
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
    viewColumns: false,
    filterType: 'dropdown',
    selectableRows: 'none',
    responsive: 'standard',
    expandableRows: true,
    expandableRowsHeader: false,
    expandableRowsOnClick: true,
    count: workspaces?.total_count,
    rowsPerPage: pageSize,
    page,
    print: false,
    download: false,
    elevation: 0,
    serverSide: true,
    search: false,
    textLabels: {
      selectedRows: {
        text: 'user(s) selected',
      },
    },
    sortOrder: {
      name: sortOrder.split(' ')[0],
      direction: sortOrder.split(' ')[1],
    },
    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          if (tableState.selectedRows.data.length) {
            tableState.selectedRows = {
              data: [],
              lookup: {},
            };
          }

          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          setSearch(tableState.searchText !== null ? tableState.searchText : '');
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
        default:
          break;
      }
    },
    renderExpandableRow: (rowData) => <WorkSpaceDataTable rowData={rowData} />,
  };

  return (
    <NoSsr>
      <UsesSistent>
        {CAN(keys.VIEW_WORKSPACE.action, keys.VIEW_WORKSPACE.subject) ? (
          <>
            <ToolWrapper>
              <CreateButtonWrapper>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={(e) => handleWorkspaceModalOpen(e, WORKSPACE_ACTION_TYPES.CREATE)}
                  sx={{
                    backgroundColor: '#607d8b',
                    padding: '8px',
                    borderRadius: '5px',
                    marginRight: '2rem',
                  }}
                  disabled={!CAN(keys.CREATE_WORKSPACE.action, keys.CREATE_WORKSPACE.subject)}
                  data-cy="btnResetDatabase"
                >
                  <AddIconCircleBorder sx={{ width: '20px', height: '20px' }} />
                  <Typography
                    sx={{
                      paddingLeft: '4px',
                      marginRight: '4px',
                      textTransform: 'none',
                    }}
                  >
                    Create
                  </Typography>
                </Button>
              </CreateButtonWrapper>
              <Box display={'flex'}>
                <SearchBar
                  onSearch={(value) => {
                    setSearch(value);
                  }}
                  placeholder="Search Workspaces..."
                  expanded={isSearchExpanded}
                  setExpanded={setIsSearchExpanded}
                />
                {viewType !== 'grid' && (
                  <CustomColumnVisibilityControl
                    columns={columns}
                    customToolsProps={{ columnVisibility, setColumnVisibility }}
                  />
                )}
                <ViewSwitch view={viewType} changeView={setViewType} />
              </Box>
            </ToolWrapper>
            {selectedWorkspaces.length > 0 && (
              <BulkActionWrapper>
                <Typography>
                  {selectedWorkspaces.length > 1
                    ? `${selectedWorkspaces.length} workspaces selected`
                    : `${selectedWorkspaces.length} workspace selected`}
                </Typography>
                <Button>
                  <DeleteIcon
                    fill={theme.palette.text.default}
                    onClick={handleBulkDeleteWorkspaceConfirm}
                    disabled={
                      CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject) &&
                      selectedWorkspaces.length > 0
                        ? false
                        : true
                    }
                  />
                </Button>
              </BulkActionWrapper>
            )}
            <>
              {workspaces.length === 0 ? (
                <EmptyState
                  icon={<WorkspaceIcon height="6rem" width="6rem" fill="#808080" />}
                  message="No workspace available"
                  pointerLabel="Click “Create” to establish your first workspace."
                />
              ) : viewType === 'grid' ? (
                <>
                  <Grid container spacing={2}>
                    {workspaces.map((workspace) => (
                      <MesheryWorkspaceCard
                        key={workspace.id}
                        workspaceDetails={workspace}
                        handleAssignDesignModalOpen={handleAssignDesignModalOpen}
                        handleAssignTeamModalOpen={handleAssignTeamModalOpen}
                        handleAssignEnvironmentModalOpen={handleAssignEnvironmentModalOpen}
                        handleWorkspaceModalOpen={handleWorkspaceModalOpen}
                        handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
                        handleBulkSelect={handleBulkSelect}
                        selectedWorkspaces={selectedWorkspaces}
                      />
                    ))}
                  </Grid>
                  <Grid
                    container
                    sx={{ padding: '2rem 0', marginTop: '20px' }}
                    flex
                    justifyContent="center"
                    spacing={2}
                  >
                    <Pagination
                      count={Math.ceil(workspacesData?.total_count / pageSize)}
                      page={page + 1}
                      onChange={debounce((_, page) => setPage(page - 1), 150)}
                      boundaryCount={3}
                      renderItem={(item) => (
                        <PaginationItem
                          slots={{ previous: ChevronLeftIcon, next: ChevronRightIcon }}
                          {...item}
                        />
                      )}
                    />
                  </Grid>
                </>
              ) : (
                <UsesSistent>
                  <ResponsiveDataTable
                    columns={columns}
                    data={workspacesData?.workspaces}
                    options={options}
                    columnVisibility={columnVisibility}
                    tableCols={tableCols}
                    updateCols={updateCols}
                    backgroundColor={theme.palette.background.card}
                  />
                </UsesSistent>
              )}
            </>
            {(actionType === WORKSPACE_ACTION_TYPES.CREATE
              ? CAN(keys.CREATE_WORKSPACE.action, keys.CREATE_WORKSPACE.subject)
              : CAN(keys.EDIT_WORKSPACE.action, keys.EDIT_WORKSPACE.subject)) &&
              workspaceModal.open && (
                <SisitentModal
                  open={workspaceModal.open}
                  closeModal={handleWorkspaceModalClose}
                  title={
                    actionType === WORKSPACE_ACTION_TYPES.CREATE
                      ? 'Create Workspace'
                      : 'Edit Workspace'
                  }
                >
                  <RJSFModalWrapper
                    schema={workspaceModal.schema.schema}
                    uiSchema={workspaceModal.schema.uiSchema}
                    handleSubmit={
                      actionType === WORKSPACE_ACTION_TYPES.CREATE
                        ? handleCreateWorkspace
                        : handleEditWorkspace
                    }
                    submitBtnText={actionType === WORKSPACE_ACTION_TYPES.CREATE ? 'Save' : 'Update'}
                    initialData={initialData}
                    handleClose={handleWorkspaceModalClose}
                  />
                </SisitentModal>
              )}
            <AssignmentModal
              open={environmentAssignment.assignModal}
              onClose={environmentAssignment.handleAssignModalClose}
              title={`Assign Environments to ${environmentAssignWorkspace.name}`}
              headerIcon={
                <EnvironmentIcon height="40" width="40" fill={theme.palette.common.white} />
              }
              name="Environments"
              assignableData={environmentAssignment.data}
              handleAssignedData={environmentAssignment.handleAssignData}
              originalAssignedData={environmentAssignment.workspaceData}
              emptyStateIcon={
                <EnvironmentIcon
                  height="5rem"
                  width="5rem"
                  fill={theme.palette.icon.disabled}
                  secondaryFill={theme.palette.icon.disabled}
                />
              }
              handleAssignablePage={environmentAssignment.handleAssignablePage}
              handleAssignedPage={environmentAssignment.handleAssignedPage}
              originalLeftCount={environmentAssignment.data?.length}
              originalRightCount={environmentAssignment.assignedItems?.length}
              onAssign={environmentAssignment.handleAssign}
              disableTransfer={environmentAssignment.disableTransferButton}
              helpText={`Assign Environments to ${environmentAssignWorkspace.name}`}
              isAssignAllowed={CAN(
                keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.action,
                keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.subject,
              )}
              isRemoveAllowed={CAN(
                keys.REMOVE_ENVIRONMENT_FROM_WORKSPACE.action,
                keys.REMOVE_ENVIRONMENT_FROM_WORKSPACE.subject,
              )}
            />

            <AssignmentModal
              open={teamAssignment.assignModal}
              onClose={teamAssignment.handleAssignModalClose}
              title={`Assign Teams to ${teamAssignWorkspace.name}`}
              headerIcon={
                <TeamsIcon height="40" width="40" primaryFill={theme.palette.common.white} />
              }
              name="Teams"
              assignableData={teamAssignment.data}
              handleAssignedData={teamAssignment.handleAssignData}
              originalAssignedData={teamAssignment.workspaceData}
              emptyStateIcon={<TeamsIcon height="5rem" width="5rem" />}
              handleAssignablePage={teamAssignment.handleAssignablePage}
              handleAssignedPage={teamAssignment.handleAssignedPage}
              originalLeftCount={teamAssignment.data?.length}
              originalRightCount={teamAssignment.assignedItems?.length}
              onAssign={teamAssignment.handleAssign}
              disableTransfer={teamAssignment.disableTransferButton}
              helpText={`Assign Teams to ${teamAssignWorkspace.name}`}
              isAssignAllowed={CAN(
                keys.ASSIGN_TEAM_TO_WORKSPACE.action,
                keys.ASSIGN_TEAM_TO_WORKSPACE.subject,
              )}
              isRemoveAllowed={CAN(
                keys.REMOVE_TEAM_FROM_WORKSPACE.action,
                keys.REMOVE_TEAM_FROM_WORKSPACE.subject,
              )}
            />

            <AssignmentModal
              open={designAssignment.assignModal && viewAssignment.assignModal}
              onClose={designAssignment.handleAssignModalClose}
              title={`Assign Designs and Views to ${designAssignWorkspace.name}`}
              headerIcon={
                <OutlinedPatternIcon height="40" width="40" fill={theme.palette.icon.default} />
              }
              name="Designs"
              assignableData={designAssignment.data}
              handleAssignedData={designAssignment.handleAssignData}
              originalAssignedData={designAssignment.workspaceData}
              emptyStateIcon={
                <OutlinedPatternIcon
                  height="5rem"
                  width="5rem"
                  fill={theme.palette.icon.disabled}
                />
              }
              handleAssignablePage={designAssignment.handleAssignablePage}
              handleAssignedPage={designAssignment.handleAssignedPage}
              originalLeftCount={designAssignment.data?.total_count}
              originalRightCount={designAssignment.workspaceData?.total_count}
              onAssign={isDesignActivity || isViewActivity ? handleAssignments : null}
              disableTransfer={
                designAssignment.disableTransferButton && viewAssignment.disableTransferButton
              }
              helpText={`Assign Designs and Views to ${designAssignWorkspace.name}`}
              isAssignAllowed={CAN(
                keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
                keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
              )}
              isRemoveAllowed={CAN(
                keys.REMOVE_DESIGNS_FROM_WORKSPACE.action,
                keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject,
              )}
              showViews={CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject)}
              emptyStateViewsIcon={
                <ViewIcon height="5rem" width="5rem" fill={theme.palette.icon.disabled} />
              }
              nameViews="Views"
              assignableViewsData={viewAssignment.data}
              handleAssignedViewsData={viewAssignment.handleAssignData}
              originalAssignedViewsData={viewAssignment.workspaceData}
              handleAssignableViewsPage={viewAssignment.handleAssignablePage}
              handleAssignedViewsPage={viewAssignment.handleAssignedPage}
              originalLeftViewsCount={viewAssignment.data?.total_count}
              originalRightViewsCount={viewAssignment.workspaceData?.total_count}
              isAssignAllowedViews={CAN(
                keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
                keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
              )}
              isRemoveAllowedViews={CAN(
                keys.REMOVE_VIEWS_FROM_WORKSPACE.action,
                keys.REMOVE_VIEWS_FROM_WORKSPACE.subject,
              )}
            />
            <_PromptComponent ref={ref} />
            <_PromptComponent ref={bulkDeleteRef} />
          </>
        ) : (
          <DefaultError />
        )}
      </UsesSistent>
    </NoSsr>
  );
};

const mapStateToProps = (state) => {
  const organization = state.get('organization');
  return {
    organization,
  };
};

const WorkspacesPageWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <UsesSistent>
        <Workspaces {...props} />
      </UsesSistent>
    </NoSsr>
  );
};

export default connect(mapStateToProps)(withRouter(WorkspacesPageWithErrorBoundary));
