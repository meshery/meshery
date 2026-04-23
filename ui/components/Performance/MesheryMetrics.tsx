import React from 'react';
import { Typography, Button, styled } from '@sistent/sistent';
import { AddCircleOutline as AddIcon } from '../icons';
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

const MetricsTitle = styled(Typography)(({ theme }) => ({
  margin: '0 0 2.5rem 0',
  color: theme.palette.mode === 'dark' ? theme.palette.common.white : theme.palette.text.primary,
}));

const NoMetricsText = styled(Typography)(({ theme }) => ({
  fontSize: '1.5rem',
  marginBottom: '2rem',
  color: theme.palette.mode === 'dark' ? theme.palette.common.white : theme.palette.text.secondary,
}));

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
