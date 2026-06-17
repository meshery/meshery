import React from 'react';
import TelemetryDashboards from '../components/telemetry/dashboards';
import { MesheryPage, PageContainer } from '../components/MesheryPage';

const Telemetry = () => (
  <MesheryPage title="Telemetry">
    <PageContainer>
      <TelemetryDashboards />
    </PageContainer>
  </MesheryPage>
);

export default Telemetry;
