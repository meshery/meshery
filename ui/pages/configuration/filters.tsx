import React from 'react';
import MesheryFilters from '../../components/MesheryFilters/Filters';
import { MesheryPage, PageContainer } from '../../components/MesheryPage';

function NewFilters() {
  return (
    <MesheryPage title="Filters">
      <PageContainer sx={{ maxWidth: '90%' }}>
        <MesheryFilters />
      </PageContainer>
    </MesheryPage>
  );
}

export default NewFilters;
