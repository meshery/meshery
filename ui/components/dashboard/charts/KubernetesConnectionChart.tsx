import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';
import {
  Box,
  InfoOutlinedIcon,
  KubernetesIcon,
  Typography,
  useTheme,
  CustomTooltip,
  IconButton,
} from '@sistent/sistent';

export default function KubernetesConnectionStatsChart() {
  const { data: connectionData } = useGetConnectionsQuery({
    page: 0,
    pagesize: 'all',
    kind: JSON.stringify(['kubernetes']),
  });
  const router = useRouter();
  const theme = useTheme();

  const chartData = useMemo(
    () =>
      connectionData?.connections
        ? Object.entries(
            connectionData.connections.reduce(
              (acc: Record<string, number>, connection: { status: string }) => {
                if (isValidColumnName(connection.status)) {
                  acc[connection.status] = (acc[connection.status] || 0) + 1;
                }
                return acc;
              },
              {} as Record<string, number>,
            ),
          )
        : [],
    [connectionData?.connections],
  );

  const chartOptions = useMemo(
    () => ({
      size: {
        height: 250,
      },
      data: {
        columns: chartData,
        type: donut(),
        colors: dataToColors(chartData, theme),
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
    }),
    [chartData, router, theme],
  );

  const canViewConnections = CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject);
  const header = (
    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
        <KubernetesIcon
          fill={
            theme.palette.mode === 'light'
              ? theme.palette.icon.default
              : theme.palette.icon.disabled
          }
        />

        <Typography variant="h6" fontWeight="700">
          KUBERNETES CLUSTER STATUS
        </Typography>
      </div>
      <div onClick={(e) => e.stopPropagation()}>
        <CustomTooltip
          title="This chart shows the status of connections to your Kubernetes clusters."
          placement="top"
          interactive
        >
          <IconButton
            component="span"
            disableRipple={true}
            disableFocusRipple={true}
            onMouseDown={(e) => e.stopPropagation()}
            onTouchStart={(e) => e.stopPropagation()}
            style={{ marginLeft: '5px' }}
          >
            <InfoOutlinedIcon
              color={theme.palette.icon.default}
              style={{ ...iconSmall, cursor: 'pointer' }}
            />
          </IconButton>
        </CustomTooltip>
      </div>
    </div>
  );

  if (chartData.length === 0) {
    return (
      <DashboardSection>
        {header}
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            alignContent: 'center',
            height: '100%',
          }}
        >
          <ConnectCluster message={'No connections found in your clusters'} />
        </Box>
      </DashboardSection>
    );
  }

  return (
    <Link
      href="/management/connections"
      style={{
        textDecoration: 'none',
        pointerEvents: !canViewConnections ? 'none' : 'auto',
      }}
    >
      <DashboardSection>
        {header}
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            alignContent: 'center',
            height: '100%',
          }}
        >
          <BBChart options={chartOptions} />
        </Box>
      </DashboardSection>
    </Link>
  );
}
