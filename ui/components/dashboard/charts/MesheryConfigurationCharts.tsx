import React, { useMemo } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../general/BBChart';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/meshery-mesh-interface/PatternService/CustomTextTooltip';
import { useGetPatternsQuery } from '@/rtk-query/design';
import { useGetFiltersQuery } from '@/rtk-query/filter';
import { Keys } from '@meshery/schemas/permissions';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';

import { Box, InfoOutlinedIcon, Typography, useHasPermission, useTheme } from '@sistent/sistent';

type ChartColumn = [string, number];

/**
 * Displays summary counts of Meshery Designs and WASM Filters.
 * Skips each query unless its respective view permission (`CatalogManagementViewDesigns`
 * / `CatalogManagementViewFilters`) is granted.
 */
export default function MesheryConfigurationChart() {
  const router = useRouter();
  const theme = useTheme();
  const canViewDesigns = useHasPermission(Keys.CatalogManagementViewDesigns);
  const canViewFilters = useHasPermission(Keys.CatalogManagementViewFilters);

  const { data: patternsData, error: patternsError } = useGetPatternsQuery(
    {
      page: 0,
      pagesize: 1,
    },
    {
      skip: !canViewDesigns,
    },
  );

  const { data: filtersData, error: filtersError } = useGetFiltersQuery(
    {
      page: 0,
      pagesize: 1,
    },
    {
      skip: !canViewFilters,
    },
  );

  const chartData = useMemo<ChartColumn[]>(() => {
    const nextChartData: ChartColumn[] = [];

    if (canViewDesigns && !patternsError && patternsData?.patterns) {
      nextChartData.push(['Designs', patternsData.totalCount ?? 0]);
    }

    if (canViewFilters && !filtersError && filtersData?.filters) {
      nextChartData.push(['Filters', filtersData.totalCount ?? 0]);
    }

    return nextChartData;
  }, [
    canViewDesigns,
    canViewFilters,
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
        colors: dataToColors(chartData, theme),
        onclick: (dataPoint: { name: string }) => {
          const isAllowed =
            dataPoint.name === 'Designs'
              ? canViewDesigns
              : dataPoint.name === 'Filters'
                ? canViewFilters
                : false;
          if (!isAllowed) return;
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
    [chartData, canViewDesigns, canViewFilters, router, theme],
  );

  const section = (
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
        {!canViewDesigns && !canViewFilters ? (
          <Typography variant="body1" align="center">
            You do not have permission to view Designs or Filters.
          </Typography>
        ) : chartData.length > 0 ? (
          <BBChart options={chartOptions} />
        ) : (
          <ConnectCluster message={'No connections found in your clusters'} />
        )}
      </Box>
    </DashboardSection>
  );

  let configurationHref: string | null = null;
  if (canViewDesigns) {
    configurationHref = '/configuration/designs';
  } else if (canViewFilters) {
    configurationHref = '/configuration/filters';
  }

  return configurationHref ? (
    <Link href={configurationHref} style={{ textDecoration: 'none' }}>
      {section}
    </Link>
  ) : (
    section
  );
}
