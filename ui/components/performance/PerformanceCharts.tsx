import React from 'react';
import { Typography } from '@sistent/sistent';
import GrafanaCustomCharts from '../telemetry/grafana/GrafanaCustomCharts';

interface PerformanceChartsProps {
  localStaticPrometheusBoardConfig: any;
  prometheus: any;
  grafana: any;
  testUUID: string;
}

/**
 * PerformanceCharts renders the three optional telemetry chart sections that
 * appear under the form: pre-configured static Prometheus boards, additional
 * user-selected Prometheus boards, and Grafana boards.
 *
 * Extracted from `performance/index.tsx` in Phase 5.a so the entry point
 * stays under the 600-line size budget.
 */
const PerformanceCharts: React.FC<PerformanceChartsProps> = ({
  localStaticPrometheusBoardConfig,
  prometheus,
  grafana,
  testUUID,
}) => {
  let displayStaticCharts = null;
  let displayPromCharts = null;
  let displayGCharts = null;

  if (
    localStaticPrometheusBoardConfig &&
    localStaticPrometheusBoardConfig !== null &&
    Object.keys(localStaticPrometheusBoardConfig).length > 0 &&
    prometheus.prometheusURL !== ''
  ) {
    // only add testUUID to the board that should be persisted
    if (localStaticPrometheusBoardConfig.cluster) {
      localStaticPrometheusBoardConfig.cluster.testUUID = testUUID;
    }
    displayStaticCharts = (
      <React.Fragment>
        <Typography variant="h6" gutterBottom sx={{ textAlign: 'center' }}>
          Node Metrics
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={[
            localStaticPrometheusBoardConfig.cluster,
            localStaticPrometheusBoardConfig.node,
          ]}
          prometheusURL={prometheus.prometheusURL}
        />
      </React.Fragment>
    );
  }
  if (prometheus.selectedPrometheusBoardsConfigs.length > 0) {
    displayPromCharts = (
      <React.Fragment>
        <Typography variant="h6" gutterBottom sx={{ textAlign: 'center' }}>
          Prometheus charts
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={prometheus.selectedPrometheusBoardsConfigs}
          prometheusURL={prometheus.prometheusURL}
        />
      </React.Fragment>
    );
  }
  if (grafana.selectedBoardsConfigs.length > 0) {
    displayGCharts = (
      <React.Fragment>
        <Typography variant="h6" gutterBottom sx={{ textAlign: 'center' }}>
          Grafana charts
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={grafana.selectedBoardsConfigs}
          grafanaURL={grafana.grafanaURL}
          grafanaAPIKey={grafana.grafanaAPIKey}
        />
      </React.Fragment>
    );
  }

  return (
    <>
      {displayStaticCharts}
      {displayPromCharts}
      {displayGCharts}
    </>
  );
};

export default PerformanceCharts;
