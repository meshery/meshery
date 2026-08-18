import React from 'react';
import { WorkspacesComponent } from '@/components/lifecycle';
import { MesheryPage, PageContainer } from '../../components/general/MesheryPage';

const Workspaces = () => (
  <MesheryPage title="Workspaces">
    <PageContainer>
      <WorkspacesComponent />
    </PageContainer>
  </MesheryPage>
);

export default Workspaces;
