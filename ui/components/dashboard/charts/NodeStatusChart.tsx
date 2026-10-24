import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../general/BBChart';
import { CircularProgress, KEPPEL, Typography, Stack } from '@sistent/sistent';
import { useTheme } from '@/theme';
import { getLegendTemplate } from './utils';
import ConnectCluster from './ConnectCluster';
import { LoadingContainer, ChartSectionWithColumn, LegendSection } from '../style';

type NodeStatus = { status: string; count: number };

type NodeStatusChartProps = {
  nodeData?: NodeStatus[];
  isClusterLoading?: boolean;
};

export const NodeStatusChart = ({ nodeData, isClusterLoading }: NodeStatusChartProps) => {
  const theme = useTheme();
  const totalNodes = nodeData?.reduce((acc, node) => acc + node.count, 0);
  const columns = useMemo(() => nodeData?.map((node) => [node.status, node.count]), [nodeData]);
  const chartOptions = useMemo(
    () => ({
      data: {
        columns,
        type: donut(),
        colors: {
          Ready: KEPPEL,
          'Not Ready': theme.palette.error.main,
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
    }),
    [columns, totalNodes, theme.palette.error.main],
  );

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
            <LegendSection id="nodeLegend" />
          </>
        )}
      </>
    </ChartSectionWithColumn>
  );
};
