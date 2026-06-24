import React from 'react';
import MesheryCredentialComponent from '../../components/MesheryCredentialComponent';
import { MesheryPage, PageContainer } from '../../components/MesheryPage';

const Credentials = () => (
  <MesheryPage title="Credentials">
    <PageContainer>
      <MesheryCredentialComponent />
    </PageContainer>
  </MesheryPage>
);

export default Credentials;
