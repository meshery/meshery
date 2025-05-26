import React, { useEffect, useState } from 'react';
import { WorkspaceActivityCard } from '@layer5/sistent';
import { useGetEventsOfWorkspaceQuery, useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { useSelector } from 'react-redux';

const WorkspaceActivityWidget = () => {
  const { organization: currentOrg } = useSelector((state) => state.ui);
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
      workspacePagePath="/management/workspaces"
    />
  );
};

export default WorkspaceActivityWidget;
