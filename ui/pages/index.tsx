import React from 'react';
import Dashboard from '@/components/dashboard';
import { MesheryPage } from '../components/MesheryPage';

function Index() {
  return (
    <MesheryPage title="Dashboard">
      <Dashboard />
    </MesheryPage>
  );
}

export default Index;
