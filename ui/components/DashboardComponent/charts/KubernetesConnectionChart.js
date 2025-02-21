import React from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import Link from 'next/link';
import { iconMedium, iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import { InfoOutlined } from '@mui/icons-material';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';
import { Box, Typography, useTheme } from '@layer5/sistent';

export default function KubernetesConnectionStatsChart() {
  const { data: connectionData } = useGetConnectionsQuery({
    page: 0,
    pagesize: 'all',
    kind: JSON.stringify(['kubernetes']),
  });
  const router = useRouter();
  const theme = useTheme();

  const chartData = connectionData?.connections
    ? Object.entries(
        connectionData.connections.reduce((acc, connection) => {
          if (isValidColumnName(connection.status)) {
            acc[connection.status] = (acc[connection.status] || 0) + 1;
          }
          return acc;
        }, {}),
      )
    : [];

  const chartOptions = {
    size: {
      height: 250,
    },
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
      title: 'Clusters\n  Status',
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
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <img style={iconMedium} src="/static/img/kubernetes.svg" />

            <Typography variant="h6" fontWeight="700">
              KUBERNETES CLUSTER STATUS
            </Typography>
          </div>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip title="This chart shows the status of the connections status to your Kubernetes clusters.">
              <div>
                <InfoOutlined
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
