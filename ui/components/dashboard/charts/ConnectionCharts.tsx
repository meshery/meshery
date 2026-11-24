import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../general/BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/meshery-mesh-interface/PatternService/CustomTextTooltip';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import { Keys } from '@meshery/schemas/permissions';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';
import { Box, InfoOutlinedIcon, Typography, useHasPermission, useTheme } from '@sistent/sistent';

/**
 * Displays aggregate connection counts by status.
 * Skips the connections query unless `Keys.WorkspaceManagementViewConnections` is granted.
 */
export default function ConnectionStatsChart() {
  const canViewConnections = useHasPermission(Keys.WorkspaceManagementViewConnections);
  const { data: connectionsData } = useGetConnectionsQuery(
    {
      page: 0,
      pagesize: 'all',
    },
    {
      skip: !canViewConnections,
    },
  );
  const router = useRouter();
  const theme = useTheme();

  // Compute status counts from connections data
  const chartData = useMemo(() => {
    if (!canViewConnections || !connectionsData?.connections) return [];
    const statusCounts: Record<string, number> = {};
    connectionsData.connections.forEach((conn: { status?: string }) => {
      const status = conn.status || 'unknown';
      statusCounts[status] = (statusCounts[status] || 0) + 1;
    });
    return Object.entries(statusCounts)
      .filter(([status]) => isValidColumnName(status))
      .map(([status, count]) => [status, count]);
  }, [canViewConnections, connectionsData]);

  const chartOptions = useMemo(
    () => ({
      data: {
        columns: chartData,
        type: donut(),
        colors: dataToColors(chartData, theme),
        onclick: function () {
          if (!canViewConnections) return;
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
    }),
    [chartData, canViewConnections, router, theme],
  );

  const section = (
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
              <InfoOutlinedIcon
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
  );

  return canViewConnections ? (
    <Link href="/management/connections" style={{ textDecoration: 'none' }}>
      {section}
    </Link>
  ) : (
    section
  );
}
