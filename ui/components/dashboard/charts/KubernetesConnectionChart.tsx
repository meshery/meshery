import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../general/BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/meshery-mesh-interface/PatternService/CustomTextTooltip';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { useRouter } from 'next/router';
import { DashboardSection, LoadingContainer } from '../style';
import ConnectCluster from './ConnectCluster';
import {
  Box,
  CircularProgress,
  InfoOutlinedIcon,
  KubernetesIcon,
  Typography,
  useTheme,
} from '@sistent/sistent';
import WidgetErrorFallback from '../widgets/WidgetErrorFallback';

export default function KubernetesConnectionStatsChart() {
  const {
    data: connectionData,
    isLoading,
    isError,
  } = useGetConnectionsQuery({
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

  const canViewConnections = CAN(
    Keys.WorkspaceManagementViewConnections.id,
    Keys.WorkspaceManagementViewConnections.function,
  );
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
        <CustomTextTooltip title="This chart shows the status of connections to your Kubernetes clusters.">
          <div>
            <InfoOutlinedIcon
              color={theme.palette.icon.default}
              style={{ ...iconSmall, marginLeft: '0.5rem', cursor: 'pointer' }}
            />
          </div>
        </CustomTextTooltip>
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <DashboardSection>
        {header}
        <LoadingContainer>
          <CircularProgress />
        </LoadingContainer>
      </DashboardSection>
    );
  }

  if (isError) {
    return (
      <DashboardSection>
        {header}
        <WidgetErrorFallback
          widgetTitle="Kubernetes Cluster Status"
          message="Unable to load your cluster connections. Please try again later."
        />
      </DashboardSection>
    );
  }

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
