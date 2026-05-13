import React from 'react';
import Dashboard from '../components/Dashboard';
import { MesheryPage } from '../components/MesheryPage';

function Index() {
  return (
    <MesheryPage title="Dashboard">
      <Dashboard />
    </MesheryPage>
  );
}

export default Index;
