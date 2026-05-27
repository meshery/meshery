import React from 'react';
import PerformanceProfiles from '@/components/performance/PerformanceProfiles';
import { MesheryPage } from '../../components/MesheryPage';

function Results() {
  return (
    <MesheryPage title="Profiles" headTitle="Performance Profiles">
      <PerformanceProfiles />
    </MesheryPage>
  );
}

export default Results;
