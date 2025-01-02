import React from 'react';
import { donut } from 'billboard.js';
import { Stack } from '@mui/material';
import { Typography, SAFFRON, CircularProgress } from '@layer5/sistent';
import BBChart from '@/components/BBChart';
import { getLegendTemplate } from './utils';
import ConnectCluster from './ConnectCluster';
import { LoadingContainer } from './style';
export const PodStatusChart = ({ classes, podData, isClusterLoading }) => {
  const columns = podData?.map((pod) => {
    return [pod.status, pod.count];
  });

  const totalPods = podData?.reduce((acc, pod) => acc + pod.count, 0);
  const chartOptions = {
    data: {
      columns: columns,
      type: donut(),
      colors: {
        Running: '#00A18F',
        Pending: SAFFRON,
        Failed: '#D32F2F',
        Succeeded: '#1976D2',
        Unknown: '#757575',
      },
    },
    arc: {
      cornerRadius: {
        ratio: 0.05,
      },
    },
    donut: {
      title: `${totalPods}\nPods`,
      padAngle: 0.03,
      label: {
        format: (value) => `${((value / totalPods) * 100).toFixed(1)}%`,
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
        bindto: '#podLegend',
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
        <Typography variant="h6">Pod Status Overview</Typography>
      </Stack>
      <div className={classes.chartSection}>
        {isClusterLoading ? (
          <LoadingContainer>
            <CircularProgress />
          </LoadingContainer>
        ) : !podData ? (
          <ConnectCluster />
        ) : totalPods === 0 ? (
          <Typography variant="body1" align="center">
            No pods are currently in the cluster
          </Typography>
        ) : (
          <>
            <BBChart options={chartOptions} />
            <div id="podLegend" className={classes.legendSection}></div>
          </>
        )}
      </div>
    </div>
  );
};
