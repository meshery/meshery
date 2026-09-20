import React, { useMemo } from 'react';
import BBChart from '../../general/BBChart';
import { donut } from 'billboard.js';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/meshery-mesh-interface/PatternService/CustomTextTooltip';
import {
  useGetCategoriesSummary,
  useGetComponentsQuery,
  useGetMeshModelsQuery,
  useGetRelationshipsQuery,
  useGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import { DashboardSection } from '../style';
import { Keys } from '@meshery/schemas/permissions';
import { useRouter } from 'next/router';
import { Grid2, InfoOutlinedIcon, Typography, useHasPermission, useTheme } from '@sistent/sistent';

function MeshModelContructs() {
  const params = {
    page: 0,
    pagesize: '1',
  };
  const modelCount = useGetMeshModelsQuery({ params }).data?.totalCount || 0;
  const componentCount = useGetComponentsQuery({ params }).data?.totalCount || 0;
  const relationshipCount = useGetRelationshipsQuery({ params }).data?.totalCount || 0;
  const registrantsConut = useGetRegistrantsQuery({ params }).data?.totalCount || 0;
  const theme = useTheme();

  // Data Cleanup
  const data = [
    ['Models', modelCount],
    ['Components', componentCount],
    ['Relationships', relationshipCount],
    ['Registrants', registrantsConut],
  ];
  const router = useRouter();
  const chartOptions = useMemo(
    () => ({
      data: {
        columns: data,
        type: donut(),
        colors: dataToColors(data, theme),
        onclick: function (d: { name: string }) {
          router.push(`/settings?settingsCategory=Registry&tab=${d.name}`);
        },
      },
      arc: {
        cornerRadius: {
          ratio: 0.05,
        },
      },
      donut: {
        title: 'Registered\nCapabilities\nby Type',
        padAngle: 0.03,
      },
      tooltip: {
        format: {
          value: function (v) {
            return v;
          },
        },
      },
    }),
    [data, theme],
  );

  const canViewRegistry = useHasPermission(Keys.MesherySystemViewRegistry);

  return (
    <Link
      href="/settings?settingsCategory=Registry"
      style={{
        textDecoration: 'none',
        pointerEvents: !canViewRegistry ? 'none' : 'auto',
      }}
    >
      <DashboardSection>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom>
            Registry
          </Typography>

          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              placement="left"
              interactive={true}
              title={`The Meshery Registry is a critical component acting as the central repository for all capabilities known to Meshery. [Learn More](https://docs.meshery.io/concepts/logical/registry)`}
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
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <BBChart options={chartOptions} />
        </div>
      </DashboardSection>
    </Link>
  );
}

function MeshModelCategories() {
  const router = useRouter();
  const categoryMap = useGetCategoriesSummary();
  const theme = useTheme();

  const cleanedData = useMemo(
    () => Object.keys(categoryMap).map((key) => [key, categoryMap[key]]),
    [categoryMap],
  );

  const chartOptions = useMemo(
    () => ({
      data: {
        columns: cleanedData,
        colors: dataToColors(cleanedData, theme),
        type: donut(),
        onclick: function () {
          router.push('/settings?settingsCategory=Registry&tab=Models');
        },
      },
      arc: {
        cornerRadius: {
          ratio: 0.05,
        },
      },
      donut: {
        title: 'Models\nby Category',
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
            return `${v} Models`;
          },
        },
      },
      legend: {
        show: false,
      },
    }),
    [cleanedData, theme],
  );

  return (
    <Link href="/settings?settingsCategory=Registry&tab=Models" style={{ textDecoration: 'none' }}>
      <DashboardSection>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom>
            Models by Category
          </Typography>

          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              title={`Meshery Models represent the fundamental building blocks of your infrastructure. Models are categorized by their function. For example, a model for Prometheus belongs in the "Observability and Analysis" category. [Learn More](https://docs.meshery.io/concepts/logical/models)`}
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
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <BBChart options={chartOptions} />
        </div>
      </DashboardSection>
    </Link>
  );
}

const MeshModelGraph = () => {
  return (
    <Grid2 container spacing={2} size="grow">
      <Grid2 style={{ marginBottom: '0.5rem' }} size={{ xs: 12, md: 6 }}>
        <MeshModelCategories />
      </Grid2>

      <Grid2 style={{ marginBottom: '0.5rem' }} size={{ xs: 12, md: 6 }}>
        <MeshModelContructs />
      </Grid2>
    </Grid2>
  );
};

MeshModelGraph.displayName = 'MeshModalGraph';
export default MeshModelGraph;
