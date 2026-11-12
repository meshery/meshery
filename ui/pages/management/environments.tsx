import React from 'react';
import { EnvironmentComponent } from '@/components/lifecycle';
import { MesheryPage, PageContainer } from '../../components/general/MesheryPage';

const Environments = () => (
  <MesheryPage title="Environments">
    <PageContainer>
      <EnvironmentComponent />
    </PageContainer>
  </MesheryPage>
);

export default Environments;
