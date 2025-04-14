//@ts-check
import React, { useState } from 'react';
import {
  Grid,
  L5DeleteIcon,
  Modal,
  Pagination,
  PaginationItem,
  useDesignAssignment,
  useEnvironmentAssignment,
  useTeamAssignment,
  useViewAssignment,
  ModalBody,
  ModalButtonPrimary,
  ModalFooter,
  ModalButtonSecondary,
  useTheme,
  AssignmentModal,
  TeamsIcon,
  EnvironmentIcon,
  DesignIcon,
  ViewIcon,
} from '@layer5/sistent';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignEnvironmentToWorkspaceMutation,
  useAssignTeamToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useDeleteWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetTeamsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useUnassignTeamFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '../../../rtk-query/workspace';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import { useNotificationHandlers } from '@/utils/hooks/useNotification';
import LeftArrowIcon from '@/assets/icons/LeftArrowIcon';
import RightArrowIcon from '@/assets/icons/RightArrowIcon';
import { UserCommonBox } from './styles';
import MesheryWorkspaceCard from './MesheryWorkspaceCard';

const WorkspaceGridView = ({
  workspacesData,
  handleWorkspaceModalOpen,
  handleDeleteWorkspaceConfirm,
  totalPages,
  page,
  setPage,
}) => {
  const [deleteWorkspacesModal, setDeleteWorkspacesModal] = useState(false);
  const [teamAssignWorkspace, setTeamAssignWorkspace] = useState({});
  const [environmentAssignWorkspace, setEnvironmentAssignWorkspace] = useState({});
  const [designAssignWorkspace, setDesignAssignWorkspace] = useState({});
  const [viewsAssignWorkspace, setViewsAssignWorkspace] = useState({});
  const [selectedWorkspaces, setSelectedWorkspaces] = useState([]);
  const [deleteWorkspace] = useDeleteWorkspaceMutation();

  const handleDeleteWorkspacesModalClose = () => {
    setDeleteWorkspacesModal(false);
  };

  const handleDeleteWorkspacesModalOpen = () => {
    setDeleteWorkspacesModal(true);
  };

  const { handleSuccess, handleError } = useNotificationHandlers();
  const handleDeleteWorkspace = (id) => {
    deleteWorkspace({
      workspaceId: id,
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace deleted`))
      .catch((error) => handleError(`Workspace Delete Error: ${error?.data}`));
  };

  const handleBulkDeleteEnv = () => {
    selectedWorkspaces.map((workspaceId) => {
      handleDeleteWorkspace(workspaceId);
    });
    setSelectedWorkspaces([]);
    handleDeleteWorkspacesModalClose();
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

  const teamAssignment = useTeamAssignment({
    workspaceId: teamAssignWorkspace?.id,
    isTeamsVisible: CAN(keys.VIEW_TEAMS.action, keys.VIEW_TEAMS.subject),
    useAssignTeamToWorkspaceMutation: useAssignTeamToWorkspaceMutation,
    useGetTeamsOfWorkspaceQuery: useGetTeamsOfWorkspaceQuery,
    useUnassignTeamFromWorkspaceMutation: useUnassignTeamFromWorkspaceMutation,
  });

  const environmentAssignment = useEnvironmentAssignment({
    workspaceId: environmentAssignWorkspace?.id,
    isEnvironmentsVisible: CAN(keys.VIEW_ENVIRONMENTS.action, keys.VIEW_ENVIRONMENTS.subject),

    useAssignEnvironmentToWorkspaceMutation: useAssignEnvironmentToWorkspaceMutation,
    useGetEnvironmentsOfWorkspaceQuery: useGetEnvironmentsOfWorkspaceQuery,
    useUnassignEnvironmentFromWorkspaceMutation: useUnassignEnvironmentFromWorkspaceMutation,
  });

  const designAssignment = useDesignAssignment({
    workspaceId: designAssignWorkspace?.id,
    isDesignsVisible: CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject),
    useAssignDesignToWorkspaceMutation: useAssignDesignToWorkspaceMutation,
    useGetDesignsOfWorkspaceQuery: useGetDesignsOfWorkspaceQuery,
    useUnassignDesignFromWorkspaceMutation: useUnassignDesignFromWorkspaceMutation,
  });

  const viewAssignment = useViewAssignment({
    workspaceId: viewsAssignWorkspace?.id,
    isViewsVisible: CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject),
    useGetViewsOfWorkspaceQuery: useGetViewsOfWorkspaceQuery,
    useAssignViewToWorkspaceMutation: useAssignViewToWorkspaceMutation,
    useUnassignViewFromWorkspaceMutation: useUnassignViewFromWorkspaceMutation,
  });

  const handleAssignTeamModalOpen = (e, workspace) => {
    e.stopPropagation();
    setTeamAssignWorkspace(workspace);
    teamAssignment.handleAssignModal();
  };

  const handleAssignEnvironmentModalOpen = (e, workspace) => {
    e.stopPropagation();
    setEnvironmentAssignWorkspace(workspace);
    environmentAssignment.handleAssignModal();
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

  const theme = useTheme();
  return (
    <>
      <Grid container spacing={2} sx={{ marginTop: '0px' }}>
        {selectedWorkspaces?.length > 0 && (
          <UserCommonBox
            sx={{
              width: '100%',
              p: '0.8rem',
              justifyContent: 'space-between',
              marginLeft: '1rem',
              marginTop: '0.18rem',
            }}
          >
            <span>
              {selectedWorkspaces?.length > 1
                ? `${selectedWorkspaces?.length} workspaces selected`
                : `${selectedWorkspaces?.length} workspace selected`}
            </span>

            <L5DeleteIcon
              onClick={handleDeleteWorkspacesModalOpen}
              disabled={
                selectedWorkspaces.length > 0
                  ? !CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject)
                  : true
              }
            />
          </UserCommonBox>
        )}
        {workspacesData?.map((workspace) => (
          <Grid item xs={12} md={6} key={workspace.id}>
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
          </Grid>
        ))}
      </Grid>
      <Grid
        container
        sx={{ padding: '2rem 0' }}
        flex
        justifyContent="center"
        spacing={2}
        position="relative"
      >
        <Pagination
          count={totalPages}
          page={page + 1}
          sx={{
            backgroundColor: 'white',
            borderRadius: '1rem',
            padding: '0.5rem',
          }}
          onChange={() => setPage(page - 1)}
          boundaryCount={3}
          renderItem={(item) => (
            <PaginationItem slots={{ previous: LeftArrowIcon, next: RightArrowIcon }} {...item} />
          )}
        />
      </Grid>
      {CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject) && (
        <Modal
          open={deleteWorkspacesModal}
          closeModal={handleDeleteWorkspacesModalClose}
          title={'Delete Workspace'}
        >
          <ModalBody
            style={{
              padding: '2rem 4rem',
            }}
          >
            <div>{`Do you want to delete ${selectedWorkspaces.length} workspace(s) ?`}</div>
          </ModalBody>
          <ModalFooter variant="filled">
            <ModalButtonSecondary
              style={{ marginRight: '16PX' }}
              onClick={handleDeleteWorkspacesModalClose}
            >
              Cancel
            </ModalButtonSecondary>
            <ModalButtonPrimary
              style={{
                background: theme.palette.background.error.default,
              }}
              onClick={handleBulkDeleteEnv}
            >
              Delete
            </ModalButtonPrimary>
          </ModalFooter>
        </Modal>
      )}
      <AssignmentModal
        open={teamAssignment.assignModal}
        onClose={teamAssignment.handleAssignModalClose}
        title={`Assign Teams to ${teamAssignWorkspace.name}`}
        headerIcon={
          <TeamsIcon height="40" width="40" primaryFill={theme.palette.background.constant.white} />
        }
        name="Teams"
        assignableData={teamAssignment.data}
        handleAssignedData={teamAssignment.handleAssignData}
        originalAssignedData={teamAssignment.workspaceData}
        emptyStateIcon={
          <TeamsIcon
            height="5rem"
            width="5rem"
            primaryFill={theme.palette.background.supplementary}
          />
        }
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
        open={environmentAssignment.assignModal}
        onClose={environmentAssignment.handleAssignModalClose}
        title={`Assign Environments to ${environmentAssignWorkspace.name}`}
        headerIcon={<EnvironmentIcon height="40" width="40" fill="white" />}
        name="Environments"
        assignableData={environmentAssignment.data}
        handleAssignedData={environmentAssignment.handleAssignData}
        originalAssignedData={environmentAssignment.workspaceData}
        emptyStateIcon={
          <EnvironmentIcon
            height="5rem"
            width="5rem"
            fill={theme.palette.background.supplementary}
            secondaryFill={theme.palette.text.secondary}
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
        open={designAssignment.assignModal && viewAssignment.assignModal}
        onClose={designAssignment.handleAssignModalClose}
        title={`Assign Designs and Views to ${designAssignWorkspace.name}`}
        headerIcon={<DesignIcon height="40" width="40" secondaryFill="white" />}
        name="Designs"
        assignableData={designAssignment.data}
        handleAssignedData={designAssignment.handleAssignData}
        originalAssignedData={designAssignment.workspaceData}
        emptyStateIcon={
          <DesignIcon
            height="5rem"
            width="5rem"
            secondaryFill={theme.palette.background.supplementary}
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
        showViews={true}
        emptyStateViewsIcon={
          <ViewIcon height="5rem" width="5rem" fill={theme.palette.background.supplementary} />
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
    </>
  );
};

export default WorkspaceGridView;
