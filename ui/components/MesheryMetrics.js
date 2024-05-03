//@ts-check
import React from 'react';
import { Typography, Button } from '@material-ui/core';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import GrafanaCustomCharts from './telemetry/grafana/GrafanaCustomCharts';
import { makeStyles } from '@material-ui/core/styles';
import { iconMedium } from '../css/icons.styles';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const useStyles = makeStyles(() => ({
  addIcon: {
    paddingRight: '.35rem',
  },
}));

function MesheryMetrics({
  boardConfigs = [],
  grafanaURL = '',
  grafanaAPIKey = '',
  handleGrafanaChartAddition,
}) {
  const classes = useStyles();
  if (boardConfigs?.length)
    return (
      <>
        <Typography align="center" variant="h6" style={{ margin: '0 0 2.5rem 0' }}>
          Service Mesh Metrics
        </Typography>
        <GrafanaCustomCharts
          // @ts-ignore
          enableGrafanaChip
          boardPanelConfigs={boardConfigs || []}
          grafanaURL={grafanaURL || ''}
          grafanaAPIKey={grafanaAPIKey || ''}
        />
      </>
    );

  return (
    <div
      style={{
        padding: '2rem',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        flexDirection: 'column',
      }}
    >
      <Typography
        style={{ fontSize: '1.5rem', marginBottom: '2rem' }}
        align="center"
        color="textSecondary"
      >
        No Metrics Configurations Found
      </Typography>
      <Button
        aria-label="Add Grafana Charts"
        variant="contained"
        color="primary"
        size="large"
        onClick={() => handleGrafanaChartAddition()}
        disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
      >
        <AddIcon className={classes.addIcon} style={iconMedium} />
        Configure Metrics
      </Button>
    </div>
  );
}

export default MesheryMetrics;
