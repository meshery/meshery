// @ts-nocheck
import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
// import { InfoOutlined } from '@mui/icons-material';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';
import { Box, Typography, useTheme, InfoIcon } from '@sistent/sistent';

export default function ConnectionStatsChart() {
  const { data: connectionsData } = useGetConnectionsQuery({
    page: 0,
    pagesize: 'all',
  });
  const router = useRouter();
  const theme = useTheme();

  // Compute status counts from connections data
  const chartData = useMemo(() => {
    if (!connectionsData?.connections) return [];
    const statusCounts = {};
    connectionsData.connections.forEach((conn) => {
      const status = conn.status || 'unknown';
      statusCounts[status] = (statusCounts[status] || 0) + 1;
    });
    return Object.entries(statusCounts)
      .filter(([status]) => isValidColumnName(status))
      .map(([status, count]) => [status, count]);
  }, [connectionsData]);

  const chartOptions = {
    data: {
      columns: chartData,
      type: donut(),
      colors: dataToColors(chartData),
      onclick: function () {
        router.push('/management/connections');
      },
    },
    arc: {
      cornerRadius: {
        ratio: 0.05,
      },
    },
    donut: {
      title: 'Connections\n by Status',
      padAngle: 0.03,
      label: {
        format: function (value) {
          return value;
        },
      },
    },
    tooltip: {
      format: {
        value: function (v) {
          return v;
        },
      },
    },
  };

  return (
    <Link
      href="/management/connections"
      style={{
        pointerEvents: !CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject)
          ? 'none'
          : 'auto',
      }}
    >
      <DashboardSection>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom>
            Connections
          </Typography>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              title={`Meshery Connections are managed and unmanaged resources that either through discovery or manual entry can be assigned to one or more Environments. [Learn More](https://docs.meshery.io/concepts/logical/connections)`}
              placement="left"
            >
              <div>
                <InfoIcon
                  color={theme.palette.icon.default}
                  style={{ ...iconSmall, marginLeft: '0.5rem', cursor: 'pointer' }}
                />
              </div>
            </CustomTextTooltip>
          </div>
        </div>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            alignContent: 'center',
            height: '100%',
          }}
        >
          {chartData.length > 0 ? (
            <BBChart options={chartOptions} />
          ) : (
            <ConnectCluster message={'No connections found in your clusters'} />
          )}
        </Box>
      </DashboardSection>
    </Link>
  );
}
