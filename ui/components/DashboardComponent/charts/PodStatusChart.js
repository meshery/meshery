import React from 'react';
import { donut } from 'billboard.js';
import {
  Typography,
  SAFFRON,
  CircularProgress,
  KEPPEL,
  DARK_SLATE_GRAY,
  TEAL_BLUE,
  Stack,
} from '@layer5/sistent';
import BBChart from '@/components/BBChart';
import { getLegendTemplate } from './utils';
import ConnectCluster from './ConnectCluster';
import { LoadingContainer, ChartSectionWithColumn, LegendSection } from '../style';
import { ERROR_COLOR } from '@/constants/colors';

export const PodStatusChart = ({ podData, isClusterLoading }) => {
  const columns = podData?.map((pod) => {
    return [pod.status, pod.count];
  });

  const totalPods = podData?.reduce((acc, pod) => acc + pod.count, 0);
  const chartOptions = {
    data: {
      columns: columns,
      type: donut(),
      colors: {
        Running: KEPPEL,
        Pending: SAFFRON,
        Failed: ERROR_COLOR,
        Succeeded: TEAL_BLUE,
        Unknown: DARK_SLATE_GRAY,
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
    <ChartSectionWithColumn>
      <Stack direction="row" mb={2}>
        <Typography variant="h6">Pod Status Overview</Typography>
      </Stack>
      <>
        {isClusterLoading ? (
          <LoadingContainer>
            <CircularProgress />
          </LoadingContainer>
        ) : !podData ? (
          <ConnectCluster message="No workloads found in your cluster(s)." />
        ) : totalPods === 0 ? (
          <Typography variant="body1" align="center">
            No pods are currently in the cluster
          </Typography>
        ) : (
          <>
            <BBChart options={chartOptions} />
            <LegendSection id="podLegend"></LegendSection>
          </>
        )}
      </>
    </ChartSectionWithColumn>
  );
};
