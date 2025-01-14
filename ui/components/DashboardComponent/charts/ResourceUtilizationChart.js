import React from 'react';
import { gauge } from 'billboard.js';
import { Box, Typography, Stack, KEPPEL, SAFFRON, CircularProgress } from '@layer5/sistent';
import BBChart from '@/components/BBChart';
import ConnectCluster from './ConnectCluster';
import { LoadingContainer, ChartSectionWithColumn } from '../style';

export const ResourceUtilizationChart = ({ usageData, isClusterLoading }) => {
  const commonConfig = {
    gauge: {
      max: 100,
      label: {
        format: function (value) {
          return value + '%';
        },
      },
    },
    legend: {
      show: false,
    },
    color: {
      pattern: [KEPPEL, SAFFRON, '#D32F2F'],
      threshold: {
        values: [40, 70, 95],
      },
    },
    size: {
      height: 150,
    },
  };
  const cpuChartOptions = {
    ...commonConfig,
    data: {
      columns: [[usageData?.[0]?.['resource'], usageData?.[0]?.['percentage']]],
      type: gauge(),
    },
  };
  const memoryChartOptions = {
    ...commonConfig,
    data: {
      columns: [[usageData?.[1]?.['resource'], usageData?.[1]?.['percentage']]],
      type: gauge(),
    },
  };

  const diskChartOptions = {
    ...commonConfig,
    data: {
      columns: [[usageData?.[2]?.['resource'], usageData?.[2]?.['percentage']]],
      type: gauge(),
    },
  };

  return (
    <ChartSectionWithColumn>
      <Stack direction="row" mb={4}>
        <Typography variant="h6">Resource Utilization</Typography>
      </Stack>
      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          alignItems: 'center',
          gap: 3,
          '& > *': {
            flex: '1 1 calc(33.333% - 16px)',
            minWidth: '250px',
          },
        }}
      >
        {isClusterLoading ? (
          <LoadingContainer>
            <CircularProgress />
          </LoadingContainer>
        ) : !usageData ? (
          <ConnectCluster message="No workloads found in your cluster(s)." />
        ) : (
          <>
            <Box>
              <Typography variant="subtitle1" align="center" gutterBottom>
                CPU Utilization
              </Typography>
              <BBChart options={cpuChartOptions} />
            </Box>
            <Box>
              <Typography variant="subtitle1" align="center" gutterBottom>
                Memory Utilization
              </Typography>
              <BBChart options={memoryChartOptions} />
            </Box>
            <Box>
              <Typography variant="subtitle1" align="center" gutterBottom>
                Disk Utilization
              </Typography>
              <BBChart options={diskChartOptions} />
            </Box>
          </>
        )}
      </Box>
    </ChartSectionWithColumn>
  );
};
