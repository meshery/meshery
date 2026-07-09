import {
  useGetEventsOfWorkspaceQuery,
  useAssignDesignToWorkspaceMutation,
  useAssignEnvironmentToWorkspaceMutation,
  useAssignTeamToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetTeamsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useUnassignTeamFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import CAN from '@/utils/can';
import {
  DesignIcon,
  EnvironmentIcon,
  TeamsIcon,
  useDesignAssignment,
  useEnvironmentAssignment,
  useTheme,
  useViewAssignment,
  ViewIcon,
  WorkspaceCard,
} from '@sistent/sistent';
import React from 'react';
import { useEffect, useState } from 'react';
import { WORKSPACE_ACTION_TYPES } from '.';
import { Keys } from '@meshery/schemas/permissions';
import { useTeamAssignment } from '@sistent/sistent';
import { AssignmentModal } from '@sistent/sistent';

const MesheryWorkspaceCard = ({
  workspaceDetails,
  handleWorkspaceModalOpen,
  handleDeleteWorkspaceConfirm,
  handleBulkSelect,
  selectedWorkspaces,
}) => {
  const [skip, setSkip] = useState(true);
  const [skipEvents, setSkipEvents] = useState(true);
  const isViewsVisible = CAN(Keys.KanvasViewViews.id, Keys.KanvasViewViews.function);
  const isDesignsVisible = CAN(
    Keys.CatalogManagementViewDesigns.id,
    Keys.CatalogManagementViewDesigns.function,
  );
  const isTeamsVisible = CAN(
    Keys.IdentityAccessManagementViewTeams.id,
    Keys.IdentityAccessManagementViewTeams.function,
  );
  const isEnvironmentsVisible = CAN(
    Keys.WorkspaceManagementViewEnvironment.id,
    Keys.WorkspaceManagementViewEnvironment.function,
  );
  const deleted = workspaceDetails.deletedAt.Valid;

  const { data: teamsOfWorkspace } = useGetTeamsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isTeamsVisible,
    },
  );

  const { data: environmentsOfWorkspace } = useGetEnvironmentsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isEnvironmentsVisible,
    },
  );

  const { data: designsOfWorkspace } = useGetDesignsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isDesignsVisible,
    },
  );
  const { data: viewsOfWorkspace } = useGetViewsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isViewsVisible,
    },
  );

  const { data: events, isLoading: isEventsLoading } = useGetEventsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 25,
    },
    {
      skip: skipEvents,
    },
  );

  useEffect(() => {
    if (!deleted) {
      setSkip(false);
    } else {
      setSkip(true);
    }
  }, [workspaceDetails, deleted]);

  const teamsOfWorkspaceCount = teamsOfWorkspace?.totalCount ? teamsOfWorkspace.totalCount : 0;

  const environmentsOfWorkspaceCount = environmentsOfWorkspace?.totalCount
    ? environmentsOfWorkspace.totalCount
    : 0;

  const designsOfWorkspaceCount = designsOfWorkspace?.totalCount
    ? designsOfWorkspace.totalCount
    : 0;

  const viewsOfWorkspaceCount = viewsOfWorkspace?.totalCount ? viewsOfWorkspace.totalCount : 0;

  const designsAndViewsCount = designsOfWorkspaceCount + viewsOfWorkspaceCount;
  const theme = useTheme();

  const teamAssignment = useTeamAssignment({
    workspaceId: workspaceDetails.id,
    isTeamsVisible: CAN(
      Keys.IdentityAccessManagementViewTeams.id,
      Keys.IdentityAccessManagementViewTeams.function,
    ),
    useAssignTeamToWorkspaceMutation: useAssignTeamToWorkspaceMutation,
    useGetTeamsOfWorkspaceQuery: useGetTeamsOfWorkspaceQuery,
    useUnassignTeamFromWorkspaceMutation: useUnassignTeamFromWorkspaceMutation,
  });

  const environmentAssignment = useEnvironmentAssignment({
    workspaceId: workspaceDetails.id,
    isEnvironmentsVisible: CAN(
      Keys.WorkspaceManagementViewEnvironment.id,
      Keys.WorkspaceManagementViewEnvironment.function,
    ),

    useAssignEnvironmentToWorkspaceMutation: useAssignEnvironmentToWorkspaceMutation,
    useGetEnvironmentsOfWorkspaceQuery: useGetEnvironmentsOfWorkspaceQuery,
    useUnassignEnvironmentFromWorkspaceMutation: useUnassignEnvironmentFromWorkspaceMutation,
  });

  const designAssignment = useDesignAssignment({
    workspaceId: workspaceDetails.id,
    isDesignsVisible: CAN(
      Keys.CatalogManagementViewDesigns.id,
      Keys.CatalogManagementViewDesigns.function,
    ),
    useAssignDesignToWorkspaceMutation: useAssignDesignToWorkspaceMutation,
    useGetDesignsOfWorkspaceQuery: useGetDesignsOfWorkspaceQuery,
    useUnassignDesignFromWorkspaceMutation: useUnassignDesignFromWorkspaceMutation,
  });

  const viewAssignment = useViewAssignment({
    workspaceId: workspaceDetails.id,
    isViewsVisible: CAN(Keys.KanvasViewViews.id, Keys.KanvasViewViews.function),
    useGetViewsOfWorkspaceQuery: useGetViewsOfWorkspaceQuery,
    useAssignViewToWorkspaceMutation: useAssignViewToWorkspaceMutation,
    useUnassignViewFromWorkspaceMutation: useUnassignViewFromWorkspaceMutation,
  });

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
  const handleAssignDesignModalOpen = (e) => {
    e.stopPropagation();
    designAssignment.handleAssignModal();
    viewAssignment.handleAssignModal();
  };
  return (
    <>
      <WorkspaceCard
        designAndViewOfWorkspaceCount={designsAndViewsCount}
        environmentsOfWorkspaceCount={environmentsOfWorkspaceCount}
        teamsOfWorkspaceCount={teamsOfWorkspaceCount}
        isDeleteWorkspaceAllowed={CAN(
          Keys.WorkspaceManagementDeleteWorkspace.id,
          Keys.WorkspaceManagementDeleteWorkspace.function,
        )}
        isTeamAllowed={
          CAN(
            Keys.WorkspaceManagementAssignTeamToWorkspace.id,
            Keys.WorkspaceManagementAssignTeamToWorkspace.function,
          ) ||
          CAN(
            Keys.WorkspaceManagementRemoveTeamFromWorkspace.id,
            Keys.WorkspaceManagementRemoveTeamFromWorkspace.function,
          )
        }
        isEditWorkspaceAllowed={CAN(
          Keys.WorkspaceManagementEditWorkspace.id,
          Keys.WorkspaceManagementEditWorkspace.function,
        )}
        isEnvironmentAllowed={
          CAN(
            Keys.WorkspaceManagementAssignEnvironmentToWorkspace.id,
            Keys.WorkspaceManagementAssignEnvironmentToWorkspace.function,
          ) ||
          CAN(
            Keys.WorkspaceManagementRemoveEnvironmentFromWorkspace.id,
            Keys.WorkspaceManagementRemoveEnvironmentFromWorkspace.function,
          )
        }
        onFlip={() => setSkipEvents(false)}
        onFlipBack={() => setSkipEvents(true)}
        workspaceDetails={workspaceDetails}
        onEdit={(e) => handleWorkspaceModalOpen(e, WORKSPACE_ACTION_TYPES.EDIT, workspaceDetails)}
        onDelete={(e) => handleDeleteWorkspaceConfirm(e, workspaceDetails)}
        onSelect={(e) => handleBulkSelect(e, workspaceDetails.id)}
        selectedWorkspaces={selectedWorkspaces}
        onAssignTeam={() => teamAssignment.handleAssignModal()}
        onAssignEnvironment={() => environmentAssignment.handleAssignModal()}
        onAssignDesign={(e) => handleAssignDesignModalOpen(e)}
        recentActivities={events?.data}
        loadingEvents={isEventsLoading}
        isDesignAllowed={
          CAN(
            Keys.WorkspaceManagementAssignDesignsToWorkspaces.id,
            Keys.WorkspaceManagementAssignDesignsToWorkspaces.function,
          ) ||
          CAN(
            Keys.WorkspaceManagementRemoveDesignsFromWorkspaces.id,
            Keys.WorkspaceManagementRemoveDesignsFromWorkspaces.function,
          )
        }
        isViewAllowed={
          CAN(Keys.KanvasAssignViewsToWorkspace.id, Keys.KanvasAssignViewsToWorkspace.function) ||
          CAN(
            Keys.KanvasUnassignViewsFromWorkspace.id,
            Keys.KanvasUnassignViewsFromWorkspace.function,
          )
        }
        isViewsVisible={false}
        isDesignsVisible={false}
        isTeamsVisible={isTeamsVisible}
        isEnvironmentsVisible={isEnvironmentsVisible}
      />
      <AssignmentModal
        key={`teams-assignment-${workspaceDetails.id}`}
        open={teamAssignment.assignModal}
        onClose={teamAssignment.handleAssignModalClose}
        title={`Assign Teams to ${workspaceDetails.name}`}
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
        helpText={`Assign Teams to ${workspaceDetails.name}`}
        isAssignAllowed={CAN(
          Keys.WorkspaceManagementAssignTeamToWorkspace.id,
          Keys.WorkspaceManagementAssignTeamToWorkspace.function,
        )}
        isRemoveAllowed={CAN(
          Keys.WorkspaceManagementRemoveTeamFromWorkspace.id,
          Keys.WorkspaceManagementRemoveTeamFromWorkspace.function,
        )}
      />
      <AssignmentModal
        key={`environments-assignment-${workspaceDetails.id}`}
        open={environmentAssignment.assignModal}
        onClose={environmentAssignment.handleAssignModalClose}
        title={`Assign Environments to ${workspaceDetails.name}`}
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
        helpText={`Assign Environments to ${workspaceDetails.name}`}
        isAssignAllowed={CAN(
          Keys.WorkspaceManagementAssignEnvironmentToWorkspace.id,
          Keys.WorkspaceManagementAssignEnvironmentToWorkspace.function,
        )}
        isRemoveAllowed={CAN(
          Keys.WorkspaceManagementRemoveEnvironmentFromWorkspace.id,
          Keys.WorkspaceManagementRemoveEnvironmentFromWorkspace.function,
        )}
      />

      <AssignmentModal
        key={`designs-assignment-${workspaceDetails.id}`}
        open={designAssignment.assignModal && viewAssignment.assignModal}
        onClose={designAssignment.handleAssignModalClose}
        title={`Assign Designs and Views to ${workspaceDetails.name}`}
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
        originalLeftCount={designAssignment.data?.totalCount}
        originalRightCount={designAssignment.workspaceData?.totalCount}
        onAssign={isDesignActivity || isViewActivity ? handleAssignments : null}
        disableTransfer={
          designAssignment.disableTransferButton && viewAssignment.disableTransferButton
        }
        helpText={`Assign Designs and Views to ${workspaceDetails.name}`}
        isAssignAllowed={CAN(
          Keys.WorkspaceManagementAssignDesignsToWorkspaces.id,
          Keys.WorkspaceManagementAssignDesignsToWorkspaces.function,
        )}
        isRemoveAllowed={CAN(
          Keys.WorkspaceManagementRemoveDesignsFromWorkspaces.id,
          Keys.WorkspaceManagementRemoveDesignsFromWorkspaces.function,
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
        originalLeftViewsCount={viewAssignment.data?.totalCount}
        originalRightViewsCount={viewAssignment.workspaceData?.totalCount}
        isAssignAllowedViews={CAN(
          Keys.KanvasAssignViewsToWorkspace.id,
          Keys.KanvasAssignViewsToWorkspace.function,
        )}
        isRemoveAllowedViews={CAN(
          Keys.KanvasUnassignViewsFromWorkspace.id,
          Keys.KanvasUnassignViewsFromWorkspace.function,
        )}
      />
    </>
  );
};

export default MesheryWorkspaceCard;
