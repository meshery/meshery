import React from 'react';
import DesignConfigurator from '@/components/designs/configurator/MeshModel';
import { MesheryPage } from '@/components/MesheryPage';

function DesignConfiguratorPage() {
  return (
    <MesheryPage title="Configure Design" headTitle="Designs Configurator">
      <DesignConfigurator />
    </MesheryPage>
  );
}

export default DesignConfiguratorPage;
