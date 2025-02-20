import React, { useEffect, useState } from 'react';
import { Grid, WorkspaceCard } from '@layer5/sistent';
import {
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetEventsOfWorkspaceQuery,
  useGetTeamsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
} from '../../../rtk-query/workspace';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import { WORKSPACE_ACTION_TYPES } from '.';

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

  const { data: environmentsOfWorkspace } = useGetEnvironmentsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isEnvironmentsVisible,
    },
  );

  const { data: teamsofWorkspace } = useGetTeamsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isTeamsVisible,
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
  const { data: designsOfWorkspace } = useGetDesignsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
      pagesize: 1,
    },
    {
      skip: skip || !isDesignsVisible,
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

  const environmentsOfWorkspaceCount = environmentsOfWorkspace?.total_count
    ? environmentsOfWorkspace.total_count
    : 0;

  const designsOfWorkspaceCount = designsOfWorkspace?.total_count
    ? designsOfWorkspace.total_count
    : 0;

  const viewsOfWorkspaceCount = viewsOfWorkspace?.total_count ? viewsOfWorkspace.total_count : 0;
  const designsAndViewCount = designsOfWorkspaceCount + viewsOfWorkspaceCount;
  const teamsOfWorkspaceCount = teamsofWorkspace?.total_count ? teamsofWorkspace.total_count : 0;

  return (
    <Grid item xs={12} md={6} key={workspaceDetails.id}>
      <WorkspaceCard
        workspaceDetails={workspaceDetails}
        onEdit={(e) => handleWorkspaceModalOpen(e, WORKSPACE_ACTION_TYPES.EDIT, workspaceDetails)}
        onDelete={(e) => handleDeleteWorkspaceConfirm(e, workspaceDetails)}
        onSelect={(e) => handleBulkSelect(e, workspaceDetails.id)}
        selectedWorkspaces={selectedWorkspaces}
        onAssignEnvironment={(e) => handleAssignEnvironmentModalOpen(e, workspaceDetails)}
        onAssignDesign={(e) => handleAssignDesignModalOpen(e, workspaceDetails)}
        onAssignTeam={(e) => handleAssignTeamModalOpen(e, workspaceDetails)}
        teamsOfWorkspaceCount={teamsOfWorkspaceCount}
        environmentsOfWorkspaceCount={environmentsOfWorkspaceCount}
        designAndViewOfWorkspaceCount={designsAndViewCount}
        isDeleteWorkspaceAllowed={CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject)}
        isEditWorkspaceAllowed={CAN(keys.EDIT_WORKSPACE.action, keys.EDIT_WORKSPACE.subject)}
        isDesignAllowed={
          CAN(keys.ASSIGN_DESIGNS_TO_WORKSPACE.action, keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject) ||
          CAN(keys.REMOVE_DESIGNS_FROM_WORKSPACE.action, keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject)
        }
        isViewAllowed={
          CAN(keys.ASSIGN_VIEWS_TO_WORKSPACE.action, keys.ASSIGN_VIEWS_TO_WORKSPACE.subject) ||
          CAN(keys.REMOVE_VIEWS_FROM_WORKSPACE.action, keys.REMOVE_VIEWS_FROM_WORKSPACE.subject)
        }
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
        isTeamAllowed={
          CAN(keys.ASSIGN_TEAM_TO_WORKSPACE.action, keys.ASSIGN_TEAM_TO_WORKSPACE.subject) ||
          CAN(keys.REMOVE_TEAM_FROM_WORKSPACE.action, keys.REMOVE_TEAM_FROM_WORKSPACE.subject)
        }
        recentActivities={events?.data}
        loadingEvents={isEventsLoading}
        onFlip={() => setSkipEvents(false)}
        onFlipBack={() => setSkipEvents(true)}
        isViewsVisible={isViewsVisible}
        isDesignsVisible={isDesignsVisible}
        isTeamsVisible={isTeamsVisible}
        isEnvironmentsVisible={isEnvironmentsVisible}
      />
    </Grid>
  );
};

export default MesheryWorkspaceCard;
