import React from 'react';
import { AddCircleIcon as AddIcon, Button, styled, Typography } from '@sistent/sistent';
import GrafanaCustomCharts from '../telemetry/grafana/GrafanaCustomCharts';
import { iconMedium } from '../../css/icons.styles';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const MetricsConfigButton = styled(Button)({
  '& .MuiSvgIcon-root': {
    paddingRight: '.35rem',
  },
});

const MetricsContainer = styled('div')({
  padding: '2rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
});

const MetricsTitle = styled(Typography)({
  margin: '0 0 2.5rem 0',
});

const NoMetricsText = styled(Typography)({
  fontSize: '1.5rem',
  marginBottom: '2rem',
});

function MesheryMetrics({
  boardConfigs = [],
  grafanaURL = '',
  grafanaAPIKey = '',
  handleGrafanaChartAddition,
}) {
  if (boardConfigs?.length)
    return (
      <>
        <MetricsTitle align="center" variant="h6">
          Performance Metrics
        </MetricsTitle>
        <GrafanaCustomCharts
          enableGrafanaChip
          boardPanelConfigs={boardConfigs || []}
          grafanaURL={grafanaURL || ''}
          grafanaAPIKey={grafanaAPIKey || ''}
        />
      </>
    );

  return (
    <MetricsContainer>
      <NoMetricsText align="center">No Metrics Configurations Found</NoMetricsText>
      <MetricsConfigButton
        aria-label="Add Grafana Charts"
        data-testid="configure-metrics-button"
        variant="contained"
        color="primary"
        size="large"
        onClick={() => handleGrafanaChartAddition()}
        disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
      >
        <AddIcon style={iconMedium} />
        Configure Metrics
      </MetricsConfigButton>
    </MetricsContainer>
  );
}

export default MesheryMetrics;
