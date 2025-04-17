import {
  useGetTeamsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useGetEventsOfWorkspaceQuery,
} from '@/rtk-query/workspace';
import CAN from '@/utils/can';
import { WorkspaceCard } from '@layer5/sistent';
import React from 'react';
import { useEffect, useState } from 'react';
import { WORKSPACE_ACTION_TYPES } from '.';
import { keys } from '@/utils/permission_constants';

const MesheryWorkspaceCard = ({
  workspaceDetails,
  handleWorkspaceModalOpen,
  handleDeleteWorkspaceConfirm,
  handleBulkSelect,
  selectedWorkspaces,
  handleAssignEnvironmentModalOpen,
  handleAssignDesignModalOpen,
  handleAssignTeamModalOpen,
}) => {
  const [skip, setSkip] = useState(true);
  const [skipEvents, setSkipEvents] = useState(true);
  const isViewsVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);
  const isTeamsVisible = CAN(keys.VIEW_TEAMS.action, keys.VIEW_TEAMS.subject);
  const isEnvironmentsVisible = CAN(keys.VIEW_ENVIRONMENTS.action, keys.VIEW_ENVIRONMENTS.subject);
  const deleted = workspaceDetails.deleted_at.Valid;

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

  const teamsOfWorkspaceCount = teamsOfWorkspace?.total_count ? teamsOfWorkspace.total_count : 0;

  const environmentsOfWorkspaceCount = environmentsOfWorkspace?.total_count
    ? environmentsOfWorkspace.total_count
    : 0;

  const designsOfWorkspaceCount = designsOfWorkspace?.total_count
    ? designsOfWorkspace.total_count
    : 0;

  const viewsOfWorkspaceCount = viewsOfWorkspace?.total_count ? viewsOfWorkspace.total_count : 0;

  const designsAndViewsCount = designsOfWorkspaceCount + viewsOfWorkspaceCount;

  return (
    <WorkspaceCard
      designAndViewOfWorkspaceCount={designsAndViewsCount}
      environmentsOfWorkspaceCount={environmentsOfWorkspaceCount}
      teamsOfWorkspaceCount={teamsOfWorkspaceCount}
      isDeleteWorkspaceAllowed={CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject)}
      isTeamAllowed={
        CAN(keys.ASSIGN_TEAM_TO_WORKSPACE.action, keys.ASSIGN_TEAM_TO_WORKSPACE.subject) ||
        CAN(keys.REMOVE_TEAM_FROM_WORKSPACE.action, keys.REMOVE_TEAM_FROM_WORKSPACE.subject)
      }
      isEditWorkspaceAllowed={CAN(keys.EDIT_WORKSPACE.action, keys.EDIT_WORKSPACE.subject)}
      isEnvironmentAllowed={
        CAN(
          keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.action,
          keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.subject,
        ) ||
        CAN(
          keys.REMOVE_ENVIRONMENT_FROM_WORKSPACE.action,
          keys.REMOVE_ENVIRONMENT_FROM_WORKSPACE.subject,
        )
      }
      onFlip={() => setSkipEvents(false)}
      onFlipBack={() => setSkipEvents(true)}
      workspaceDetails={workspaceDetails}
      onEdit={(e) => handleWorkspaceModalOpen(e, WORKSPACE_ACTION_TYPES.EDIT, workspaceDetails)}
      onDelete={(e) => handleDeleteWorkspaceConfirm(e, workspaceDetails)}
      onSelect={(e) => handleBulkSelect(e, workspaceDetails.id)}
      selectedWorkspaces={selectedWorkspaces}
      onAssignTeam={(e) => handleAssignTeamModalOpen(e, workspaceDetails)}
      onAssignEnvironment={(e) => handleAssignEnvironmentModalOpen(e, workspaceDetails)}
      onAssignDesign={(e) => handleAssignDesignModalOpen(e, workspaceDetails)}
      recentActivities={events?.data}
      loadingEvents={isEventsLoading}
      isDesignAllowed={
        CAN(keys.ASSIGN_DESIGNS_TO_WORKSPACE.action, keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject) ||
        CAN(keys.REMOVE_DESIGNS_FROM_WORKSPACE.action, keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject)
      }
      isViewAllowed={
        CAN(keys.ASSIGN_VIEWS_TO_WORKSPACE.action, keys.ASSIGN_VIEWS_TO_WORKSPACE.subject) ||
        CAN(keys.REMOVE_VIEWS_FROM_WORKSPACE.action, keys.REMOVE_VIEWS_FROM_WORKSPACE.subject)
      }
      isViewsVisible={isViewsVisible}
      isDesignsVisible={isDesignsVisible}
      isTeamsVisible={isTeamsVisible}
      isEnvironmentsVisible={isEnvironmentsVisible}
    />
  );
};

export default MesheryWorkspaceCard;
