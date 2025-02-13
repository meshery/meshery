import React, { useEffect, useState } from 'react';
import { useLegacySelector } from 'lib/store';
import { WorkspaceActivityCard } from '@layer5/sistent';
import { useGetEventsOfWorkspaceQuery, useGetWorkspacesQuery } from '@/rtk-query/workspace';

const WorkspaceActivityWidget = () => {
  const currentOrg = useLegacySelector((state) => state.get('organization'));
  const { data: workspaces } = useGetWorkspacesQuery({
    orgId: currentOrg?.id,
  });

  const [selectedWorkspace, setSelectedWorkspace] = useState('');
  const { data: events, isLoading: isEventsLoading } = useGetEventsOfWorkspaceQuery(
    {
      workspaceId: selectedWorkspace,
      pagesize: 5,
    },
    {
      skip: !selectedWorkspace,
    },
  );

  useEffect(() => {
    setSelectedWorkspace(workspaces?.workspaces?.[0]?.id);
  }, [workspaces]);

  const handleWorkspaceChange = (event) => {
    const newWorkspaceId = event.target.value;
    setSelectedWorkspace(newWorkspaceId);
  };
  return (
    <WorkspaceActivityCard
      selectedWorkspace={selectedWorkspace}
      handleWorkspaceChange={handleWorkspaceChange}
      activities={events?.data}
      workspaces={workspaces?.workspaces}
      isEventsLoading={isEventsLoading}
      workspacePagePath= "/management/workspaces"
    />
  );
};

export default WorkspaceActivityWidget;
