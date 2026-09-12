import React from 'react';
import MesheryConnections from '../../components/connections';
import { MesheryPage, PageContainer } from '../../components/general/MesheryPage';

const Connections = () => (
  <MesheryPage title="Connections">
    <PageContainer>
      <MesheryConnections />
    </PageContainer>
  </MesheryPage>
);

export default Connections;
