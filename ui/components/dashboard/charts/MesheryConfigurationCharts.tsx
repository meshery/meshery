import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/meshery-mesh-interface/PatternService/CustomTextTooltip';
import { useGetPatternsQuery } from '@/rtk-query/design';
import { useGetFiltersQuery } from '@/rtk-query/filter';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';

import { Box, InfoOutlined, Typography, useTheme } from '@sistent/sistent';

type ChartColumn = [string, number];

export default function MesheryConfigurationChart() {
  const router = useRouter();
  const theme = useTheme();

  const { data: patternsData, error: patternsError } = useGetPatternsQuery({
    page: 0,
    pagesize: 1,
  });

  const { data: filtersData, error: filtersError } = useGetFiltersQuery({
    page: 0,
    pagesize: 1,
  });

  const chartData = useMemo<ChartColumn[]>(() => {
    const nextChartData: ChartColumn[] = [];

    if (!patternsError && patternsData?.patterns) {
      nextChartData.push(['Designs', patternsData.totalCount ?? 0]);
    }

    if (!filtersError && filtersData?.filters) {
      nextChartData.push(['Filters', filtersData.totalCount ?? 0]);
    }

    return nextChartData;
  }, [
    filtersData?.filters,
    filtersData?.totalCount,
    filtersError,
    patternsData?.patterns,
    patternsData?.totalCount,
    patternsError,
  ]);

  const chartOptions = useMemo(
    () => ({
      data: {
        columns: chartData,
        type: donut(),
        colors: dataToColors(chartData),
        onclick: (dataPoint: { name: string }) => {
          const routeName = dataPoint.name.charAt(0).toLowerCase() + dataPoint.name.slice(1);
          router.push(`/configuration/${routeName}`);
        },
      },
      arc: {
        cornerRadius: {
          ratio: 0.05,
        },
      },
      donut: {
        title: 'Content\nby Type',
        padAngle: 0.03,
        label: {
          format: (value: number) => value,
        },
      },
      tooltip: {
        format: {
          value: (value: number) => value,
        },
      },
    }),
    [chartData, router],
  );

  return (
    <Link
      href="/configuration/designs"
      style={{
        textDecoration: 'none',
        pointerEvents: !CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject) ? 'none' : 'auto',
      }}
    >
      <DashboardSection>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom>
            Configuration
          </Typography>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              placement="left"
              title={`Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured. [Learn more](https://docs.meshery.io/concepts/logical/designs)`}
            >
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
