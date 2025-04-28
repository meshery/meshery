import {
  Breadcrumbs,
  NoSsr,
  WorkspaceRecentActivityModal,
  WorkspaceTeamsTable,
} from '@layer5/sistent';
import {
  Box,
  CustomColumnVisibilityControl,
  TeamsIcon,
  WorkspaceIcon,
  Modal,
  createAndEditWorkspaceSchema,
  createAndEditWorkspaceUiSchema,
  Button,
  Typography,
  SearchBar,
  useTheme,
  PROMPT_VARIANTS,
  ModalFooter,
} from '@layer5/sistent';
import { EmptyState } from '../General';
import AddIconCircleBorder from '../../../assets/icons/AddIconCircleBorder';
import { useContext, useRef, useState } from 'react';
import {
  useAssignTeamToWorkspaceMutation,
  useCreateWorkspaceMutation,
  useDeleteWorkspaceMutation,
  useGetEventsOfWorkspaceQuery,
  useGetTeamsOfWorkspaceQuery,
  useGetWorkspacesQuery,
  useUnassignTeamFromWorkspaceMutation,
  useUpdateWorkspaceMutation,
} from '../../../rtk-query/workspace';
import { updateProgress, useLegacySelector } from '../../../lib/store';
import { useNotification, useNotificationHandlers } from '../../../utils/hooks/useNotification';
import { RJSFModalWrapper } from '../../Modal';
import _PromptComponent from '../../PromptComponent';
import { EVENT_TYPES } from '../../../lib/event-types';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import DefaultError from '@/components/General/error-404/index';

import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import ViewSwitch from '@/components/ViewSwitch';
import { CreateButtonWrapper } from './styles';
import WorkspaceGridView from './WorkspaceGridView';
import RightArrowIcon from '@/assets/icons/RightArrowIcon';
import { useGetUsersForOrgQuery, useRemoveUserFromTeamMutation } from '@/rtk-query/user';
import WorkspaceDataTable from './WorkspaceDataTable';
import { iconMedium } from 'css/icons.styles';
import { WorkspaceSwitcherContext } from '@/components/SpacesSwitcher/WorkspaceSwitcher';

export const WORKSPACE_ACTION_TYPES = {
  CREATE: 'create',
  EDIT: 'edit',
};

const columnList = [
  {
    name: 'id',
    label: 'ID',
  },
  {
    name: 'owner_email',
    label: 'Owner Email',
  },
  {
    name: 'owner_id',
    label: 'Owner Id',
  },
  {
    name: 'org_name',
    label: 'Org Name',
  },
  {
    name: 'name',
    label: 'Name',
  },
  {
    name: 'description',
    label: 'Description',
  },
  {
    name: 'owner',
    label: 'Owner',
  },
  {
    name: 'environments',
    label: 'Environments',
  },
  {
    name: 'created_at',
    label: 'Created At',
  },
  {
    name: 'updated_at',
    label: 'Updated At',
  },
  {
    name: 'actions',
    label: 'Actions',
  },
];

const Workspaces = () => {
  const theme = useTheme();
  const [workspaceModal, setWorkspaceModal] = useState({
    open: false,
    schema: {},
  });
  const organization = useLegacySelector((state) => {
    return typeof state?.get === 'function' ? state.get('organization') : state?.organization || {};
  });
  const [page, setPage] = useState(0);
  const pageSize = 10;
  const sortOrder = 'updated_at desc';
  const [search, setSearch] = useState('');
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [actionType, setActionType] = useState('');
  const [initialData, setInitialData] = useState({});
  const [editWorkspaceId, setEditWorkspaceId] = useState('');
  let [selectedWorkspace, setSelectedWorkspace] = useState({
    id: '',
    name: '',
  });
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);
  if (workspaceSwitcherContext.selectedWorkspace.id) {
    selectedWorkspace = workspaceSwitcherContext.selectedWorkspace;
    setSelectedWorkspace = workspaceSwitcherContext.setSelectedWorkspace;
  }
  const [viewType, setViewType] = useState(selectedWorkspace.id ? 'table' : 'grid');

  const [teamsModal, setTeamsModal] = useState({
    open: false,
    workspaceId: '',
    workspaceName: '',
  });
  const [activityModal, setActivityModal] = useState({
    open: false,
    workspaceId: '',
    workspaceName: '',
  });

  const handleRowClick = (rowData) => {
    const workspaceId = rowData[0];
    const workspaceName = rowData[4].props.children.at(-1); // Get the last child of the name cell
    setSelectedWorkspace({
      id: workspaceId,
      name: workspaceName,
    });
  };
  const ref = useRef(null);
  const bulkDeleteRef = useRef(null);
  const { notify } = useNotification();

  const { data: workspacesData } = useGetWorkspacesQuery(
    {
      page: page,
      pagesize: pageSize,
      search: search,
      order: sortOrder,
      orgId: organization?.id,
    },
    {
      skip: !organization?.id ? true : false,
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
        organization: organization?.id,
      });
      setEditWorkspaceId('');
    }
    fetchSchema();
  };

  const handleTeamsModalOpen = (e, workspaceId, workspaceName) => {
    e.stopPropagation();
    setTeamsModal({
      open: true,
      workspaceId: workspaceId,
      workspaceName: workspaceName,
    });
  };
  const handleTeamsModalClose = () => {
    setTeamsModal({
      open: false,
      workspaceId: '',
      workspaceName: '',
    });
  };

  const handleActivityModalOpen = (e, workspaceId, workspaceName) => {
    e.stopPropagation();
    setActivityModal({
      open: true,
      workspaceId: workspaceId,
      workspaceName: workspaceName,
    });
  };

  const handleActivityModalClose = () => {
    setActivityModal({
      open: false,
      workspaceId: '',
      workspaceName: '',
    });
  };

  const handleWorkspaceModalClose = () => {
    setWorkspaceModal({
      open: false,
      schema: {},
    });
    setActionType('');
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

  const handleViewChange = (val) => {
    if (val === viewType) return;

    setPage(0);
    setSelectedWorkspace({ id: '', name: '' });
    setViewType(val);
  };

  const [columnVisibility, setColumnVisibility] = useState({});

  return (
    <NoSsr>
      {CAN(keys.VIEW_WORKSPACE.action, keys.VIEW_WORKSPACE.subject) ? (
        <>
          <div style={{ marginBottom: '1.5rem' }}>
            <Breadcrumbs
              separator={
                <RightArrowIcon height={20} width={20} primaryFill={theme.palette.icon.default} />
              }
              aria-label="breadcrumb"
            >
              <div
                style={{
                  cursor: selectedWorkspace.id ? 'pointer' : 'default',
                  color: selectedWorkspace.id
                    ? theme.palette.background.brand.default
                    : theme.palette.text.default,
                  textDecoration: 'none',
                }}
                onClick={() => {
                  if (selectedWorkspace.id) {
                    setSelectedWorkspace({ id: '', name: '' });
                  }
                }}
              >
                All Workspaces
              </div>
              {selectedWorkspace.id && <Typography>{selectedWorkspace.name}</Typography>}
            </Breadcrumbs>
          </div>
          {!selectedWorkspace.id && (
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
              <Box display={'flex'} alignItems={'center'}>
                {!selectedWorkspace?.id && (
                  <>
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
                        columns={columnList}
                        customToolsProps={{ columnVisibility, setColumnVisibility }}
                      />
                    )}
                  </>
                )}
                <ViewSwitch
                  view={viewType}
                  changeView={handleViewChange}
                  key={`view-switch-${viewType}`} // Add key to force re-render when viewType changes
                />
              </Box>
            </ToolWrapper>
          )}
          <>
            {workspaces.length === 0 ? (
              <EmptyState
                icon={<WorkspaceIcon height="6rem" width="6rem" fill="#808080" />}
                message="No workspace available"
                pointerLabel="Click “Create” to establish your first workspace."
              />
            ) : viewType === 'grid' ? (
              <WorkspaceGridView
                handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
                handleWorkspaceModalOpen={handleWorkspaceModalOpen}
                page={page}
                setPage={setPage}
                totalPages={Math.ceil(workspacesData?.total_count / pageSize)}
                workspacesData={workspaces}
                key={`grid-view-${viewType}`}
              />
            ) : (
              <WorkspaceDataTable
                handleWorkspaceModalOpen={handleWorkspaceModalOpen}
                handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
                handleTeamsModalOpen={handleTeamsModalOpen}
                handleActivityModalOpen={handleActivityModalOpen}
                columnVisibility={columnVisibility}
                handleRowClick={handleRowClick}
                selectedWorkspace={selectedWorkspace}
                setColumnVisibility={setColumnVisibility}
                search={search}
                viewType={viewType}
              />
            )}
          </>
          {(actionType === WORKSPACE_ACTION_TYPES.CREATE
            ? CAN(keys.CREATE_WORKSPACE.action, keys.CREATE_WORKSPACE.subject)
            : CAN(keys.EDIT_WORKSPACE.action, keys.EDIT_WORKSPACE.subject)) &&
            workspaceModal.open && (
              <Modal
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
              </Modal>
            )}
          <Modal
            maxWidth="lg"
            open={teamsModal.open}
            closeModal={handleTeamsModalClose}
            title={`Manage "${teamsModal.workspaceName}" Teams`}
            headerIcon={<TeamsIcon {...iconMedium} primaryFill={theme.palette.common.white} />}
          >
            <WorkspaceTeamsTable
              workspaceId={teamsModal.workspaceId}
              isAssignTeamAllowed={CAN(
                keys.ASSIGN_TEAM_TO_WORKSPACE.action,
                keys.ASSIGN_TEAM_TO_WORKSPACE.subject,
              )}
              isDeleteTeamAllowed={CAN(keys.DELETE_TEAM.action, keys.DELETE_TEAM.subject)}
              isEditTeamAllowed={CAN(keys.EDIT_TEAM.action, keys.EDIT_TEAM.subject)}
              isLeaveTeamAllowed={CAN(keys.LEAVE_TEAM.action, keys.LEAVE_TEAM.subject)}
              useAssignTeamToWorkspaceMutation={useAssignTeamToWorkspaceMutation}
              useGetTeamsOfWorkspaceQuery={useGetTeamsOfWorkspaceQuery}
              useUnassignTeamFromWorkspaceMutation={useUnassignTeamFromWorkspaceMutation}
              workspaceName={teamsModal.workspaceName}
              fetchTeamUsers={() => {}}
              org_id={organization?.id}
              useGetUsersForOrgQuery={useGetUsersForOrgQuery}
              useNotificationHandlers={useNotificationHandlers}
              useRemoveUserFromTeamMutation={useRemoveUserFromTeamMutation}
              isRemoveTeamFromWorkspaceAllowed={CAN(
                keys.REMOVE_TEAM_FROM_WORKSPACE.action,
                keys.REMOVE_TEAM_FROM_WORKSPACE.subject,
              )}
            />
            <ModalFooter variant="filled"></ModalFooter>
          </Modal>
          <WorkspaceRecentActivityModal
            workspaceId={activityModal.workspaceId}
            workspaceName={activityModal.workspaceName}
            open={activityModal.open}
            handleClose={handleActivityModalClose}
            useGetEventsOfWorkspaceQuery={useGetEventsOfWorkspaceQuery}
          />

          <_PromptComponent ref={ref} />
          <_PromptComponent ref={bulkDeleteRef} />
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
  );
};

const WorkspacesPageWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <Workspaces {...props} />
    </NoSsr>
  );
};

export default WorkspacesPageWithErrorBoundary;
