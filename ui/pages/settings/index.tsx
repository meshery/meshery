import React from 'react';
import MesherySettings from '@/components/settings/MesherySettings';
import { MesheryPage } from '../../components/MesheryPage';

function Settings() {
  return (
    <MesheryPage title="Settings">
      <MesherySettings />
    </MesheryPage>
  );
}

export default Settings;
