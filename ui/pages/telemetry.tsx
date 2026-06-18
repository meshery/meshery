import React, { useState } from 'react';
import { Box, Tab, Tabs } from '@sistent/sistent';
import TelemetryDashboards from '../components/telemetry/dashboards';
import TelemetryMetrics from '../components/telemetry/metrics';
import { MesheryPage, PageContainer } from '../components/MesheryPage';

const Telemetry = () => {
  const [tab, setTab] = useState(0);

  return (
    <MesheryPage title="Telemetry">
      <PageContainer>
        <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
          <Tabs
            value={tab}
            onChange={(_e, v) => setTab(v)}
            sx={{ borderBottom: 1, borderColor: 'divider' }}
          >
            <Tab label="Dashboards (Grafana)" />
            <Tab label="Metrics (Prometheus)" />
          </Tabs>
          <Box sx={{ flex: 1, minHeight: 0, overflow: 'auto' }}>
            {tab === 0 ? <TelemetryDashboards /> : <TelemetryMetrics />}
          </Box>
        </Box>
      </PageContainer>
    </MesheryPage>
  );
};

export default Telemetry;
