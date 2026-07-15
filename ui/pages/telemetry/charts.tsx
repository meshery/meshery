import React from 'react';
import { Box } from '@sistent/sistent';
import TelemetryDashboards from '@/components/telemetry/dashboards';
import { MesheryPage, PageContainer } from '../../components/MesheryPage';

const TelemetryCharts = () => {
  return (
    <MesheryPage title="Charts" headTitle="Telemetry Charts">
      <PageContainer>
        <Box
          data-testid="telemetry-charts"
          sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
        >
          <Box sx={{ flex: 1, minHeight: 0, overflow: 'auto' }}>
            <TelemetryDashboards />
          </Box>
        </Box>
      </PageContainer>
    </MesheryPage>
  );
};

export default TelemetryCharts;
