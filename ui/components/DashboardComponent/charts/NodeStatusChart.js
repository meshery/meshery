import React from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { Stack } from '@mui/material';
import { CircularProgress, Typography } from '@layer5/sistent';
import { getLegendTemplate } from './utils';
import ConnectCluster from './ConnectCluster';
import { LoadingContainer } from './style';

export const NodeStatusChart = ({ classes, nodeData, isClusterLoading }) => {
  const data = nodeData?.map((node) => {
    return [node.status, node.count];
  });
  const totalNodes = nodeData?.reduce((acc, node) => acc + node.count, 0);
  const chartOptions = {
    data: {
      columns: data,
      type: donut(),
      colors: {
        Ready: '#00A18F',
        'Not Ready': '#D32F2F',
        Unknown: '#757575',
      },
    },
    arc: {
      cornerRadius: {
        ratio: 0.05,
      },
    },
    donut: {
      title: `${totalNodes}\nNodes`,
      padAngle: 0.03,
      label: {
        format: (value) => `${((value / totalNodes) * 100).toFixed(1)}%`,
      },
    },
    tooltip: {
      format: {
        value: function (v) {
          return v;
        },
      },
    },
    legend: {
      show: true,
      contents: {
        bindto: '#nodeLegend',
        template: getLegendTemplate,
      },
    },
  };

  return (
    <div
      className={classes.dashboardSection}
      style={{
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <Stack direction="row" mb={2}>
        <Typography variant="h6">Node Status Overview</Typography>
      </Stack>
      <div className={classes.chartSection}>
        {isClusterLoading ? (
          <LoadingContainer>
            <CircularProgress />
          </LoadingContainer>
        ) : !nodeData ? (
          <ConnectCluster />
        ) : totalNodes === 0 ? (
          <Typography variant="body1" align="center">
            No nodes are currently in the cluster
          </Typography>
        ) : (
          <>
            <BBChart options={chartOptions} />
            <div id="nodeLegend" className={classes.legendSection}></div>
          </>
        )}
      </div>
    </div>
  );
};
