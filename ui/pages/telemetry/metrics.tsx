import React from 'react';
import { Box } from '@sistent/sistent';
import TelemetryMetrics from '@/components/telemetry/metrics';
import { MesheryPage, PageContainer } from '../../components/general/MesheryPage';

const TelemetryMetricsPage = () => {
  return (
    <MesheryPage title="Metrics" headTitle="Telemetry Metrics">
      <PageContainer>
        <Box
          data-testid="telemetry-metrics"
          sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
        >
          <Box sx={{ flex: 1, minHeight: 0, overflow: 'auto' }}>
            <TelemetryMetrics />
          </Box>
        </Box>
      </PageContainer>
    </MesheryPage>
  );
};

export default TelemetryMetricsPage;
