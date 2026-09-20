import React from 'react';
import MesheryFilters from '@/components/filters/Filters';
import { MesheryPage, PageContainer } from '../../components/general/MesheryPage';

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
