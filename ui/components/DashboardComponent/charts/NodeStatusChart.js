import React from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { Stack } from '@mui/material';
import { CircularProgress, KEPPEL, Typography } from '@layer5/sistent';
import { getLegendTemplate } from './utils';
import ConnectCluster from './ConnectCluster';
import { LoadingContainer, ChartSectionWithColumn, LegendSection } from '../style';
import { ERROR_COLOR } from '@/constants/colors';

export const NodeStatusChart = ({ nodeData, isClusterLoading }) => {
  const data = nodeData?.map((node) => {
    return [node.status, node.count];
  });
  const totalNodes = nodeData?.reduce((acc, node) => acc + node.count, 0);
  const chartOptions = {
    data: {
      columns: data,
      type: donut(),
      colors: {
        Ready: KEPPEL,
        'Not Ready': ERROR_COLOR,
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
    <ChartSectionWithColumn>
      <Stack direction="row" mb={2}>
        <Typography variant="h6">Node Status Overview</Typography>
      </Stack>
      <>
        {isClusterLoading ? (
          <LoadingContainer>
            <CircularProgress />
          </LoadingContainer>
        ) : !nodeData ? (
          <ConnectCluster message="No workloads found in your cluster(s)." />
        ) : totalNodes === 0 ? (
          <Typography variant="body1" align="center">
            No nodes are currently in the cluster
          </Typography>
        ) : (
          <>
            <BBChart options={chartOptions} />
            <LegendSection id="nodeLegend"></LegendSection>
          </>
        )}
      </>
    </ChartSectionWithColumn>
  );
};
