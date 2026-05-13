import React from 'react';
import MesherySettings from '../../components/Settings/MesherySettings';
import { MesheryPage } from '../../components/MesheryPage';

function Settings() {
  return (
    <MesheryPage title="Settings">
      <MesherySettings />
    </MesheryPage>
  );
}

export default Settings;
