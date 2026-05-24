import React from 'react';
import MesheryPerformanceComponent from '@/components/performance/Dashboard';
import { MesheryPage } from '../../components/MesheryPage';

function Performance() {
  return (
    <MesheryPage title="Performance" headTitle="Performance Dashboard">
      <MesheryPerformanceComponent />
    </MesheryPage>
  );
}

export default Performance;
