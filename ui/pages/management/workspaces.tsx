import React from 'react';
import { WorkspacesComponent } from '../../components/Lifecycle';
import { MesheryPage, PageContainer } from '../../components/MesheryPage';

const Workspaces = () => (
  <MesheryPage title="Workspaces">
    <PageContainer>
      <WorkspacesComponent />
    </PageContainer>
  </MesheryPage>
);

export default Workspaces;
