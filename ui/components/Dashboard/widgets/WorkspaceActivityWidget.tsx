import React, { useCallback, useMemo, useState } from 'react';
import { WorkspaceActivityCard } from '@sistent/sistent';
import { useGetEventsOfWorkspaceQuery, useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { useSelector } from 'react-redux';
import type { RootState } from '@/store/store';

type WorkspaceChangeHandler = NonNullable<
  React.ComponentProps<typeof WorkspaceActivityCard>['handleWorkspaceChange']
>;

const WorkspaceActivityWidget = () => {
  const { organization: currentOrg } = useSelector((state: RootState) => state.ui);
  const { data: workspaces } = useGetWorkspacesQuery(
    { orgId: currentOrg?.id },
    { skip: !currentOrg?.id },
  );

  const workspaceOptions = workspaces?.workspaces ?? [];
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
      skip: !activeWorkspaceId,
    },
  );

  const handleWorkspaceChange = useCallback<WorkspaceChangeHandler>((event) => {
    setSelectedWorkspace(String(event.target.value));
  }, []);

  return (
    <WorkspaceActivityCard
      selectedWorkspace={activeWorkspaceId}
      handleWorkspaceChange={handleWorkspaceChange}
      activities={isEventsError ? [] : (events?.data ?? [])}
      workspaces={workspaceOptions}
      isEventsLoading={isEventsLoading}
      workspacePagePath="/management/workspaces"
    />
  );
};

export default WorkspaceActivityWidget;
