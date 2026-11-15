import React, { useCallback, useMemo, useState } from 'react';
import { Keys } from '@meshery/schemas/permissions';
import { useHasPermission, WorkspaceActivityCard } from '@sistent/sistent';
import { useGetEventsOfWorkspaceQuery, useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { useSelector } from 'react-redux';
import type { RootState } from '../../../store';
import WidgetErrorFallback from './WidgetErrorFallback';

type WorkspaceChangeHandler = NonNullable<
  React.ComponentProps<typeof WorkspaceActivityCard>['handleWorkspaceChange']
>;

const WorkspaceActivityWidget = () => {
  const { organization: currentOrg } = useSelector((state: RootState) => state.ui);
  const canViewWorkspaces = useHasPermission(Keys.WorkspaceManagementViewWorkspace);
  const {
    data: workspaces,
    isError: isWorkspacesError,
    isLoading: isWorkspacesLoading,
  } = useGetWorkspacesQuery(
    { orgId: currentOrg?.id },
    { skip: !currentOrg?.id || !canViewWorkspaces },
  );

  const workspaceOptions = useMemo(
    () => (canViewWorkspaces ? (workspaces?.workspaces ?? []) : []),
    [canViewWorkspaces, workspaces?.workspaces],
  );
  const [selectedWorkspace, setSelectedWorkspace] = useState('');
  const activeWorkspaceId = useMemo(() => {
    if (
      selectedWorkspace &&
      workspaceOptions.some((workspace) => workspace.id === selectedWorkspace)
    ) {
      return selectedWorkspace;
    }

    return workspaceOptions[0]?.id ?? '';
  }, [selectedWorkspace, workspaceOptions]);

  const {
    data: events,
    isLoading: isEventsLoading,
    isError: isEventsError,
  } = useGetEventsOfWorkspaceQuery(
    {
      workspaceId: activeWorkspaceId,
      pagesize: 5,
    },
    {
      skip: !canViewWorkspaces || !activeWorkspaceId,
    },
  );

  const handleWorkspaceChange = useCallback<WorkspaceChangeHandler>((event) => {
    setSelectedWorkspace(String(event.target.value));
  }, []);

  if (canViewWorkspaces && isWorkspacesError) {
    return (
      <WidgetErrorFallback
        widgetTitle="Workspace Activity"
        message="Unable to load your workspaces. Please try again later."
      />
    );
  }

  return (
    <WorkspaceActivityCard
      selectedWorkspace={activeWorkspaceId}
      handleWorkspaceChange={handleWorkspaceChange}
      activities={canViewWorkspaces && !isEventsError ? (events?.data ?? []) : []}
      workspaces={workspaceOptions}
      isEventsLoading={canViewWorkspaces && (isEventsLoading || isWorkspacesLoading)}
      workspacePagePath="/management/workspaces"
    />
  );
};

export default WorkspaceActivityWidget;
