import React from 'react';
import { EnvironmentComponent } from '../../components/Lifecycle';
import { MesheryPage, PageContainer } from '../../components/MesheryPage';

const Environments = () => (
  <MesheryPage title="Environments">
    <PageContainer>
      <EnvironmentComponent />
    </PageContainer>
  </MesheryPage>
);

export default Environments;
