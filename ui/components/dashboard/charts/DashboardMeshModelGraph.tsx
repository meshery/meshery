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

/**
 * Renders a Next.js Link only when `enabled` is true so denied users cannot focus or activate it via keyboard.
 */
function OptionalRegistryLink({
  enabled,
  href,
  children,
}: {
  enabled: boolean;
  href: string;
  children: React.ReactNode;
}) {
  return enabled ? (
    <Link href={href} style={{ textDecoration: 'none' }}>
      {children}
    </Link>
  ) : (
    <>{children}</>
  );
}

/**
 * Displays aggregate counts of registered models, components, relationships, and registrants.
 * Skips registry queries and disables click navigation unless `Keys.MesherySystemViewRegistry` is granted.
 */
function MeshModelContructs() {
  const canViewRegistry = useHasPermission(Keys.MesherySystemViewRegistry);
  const params = {
    page: 0,
    pagesize: '1',
  };
  const modelsQuery = useGetMeshModelsQuery({ params }, { skip: !canViewRegistry });
  const componentsQuery = useGetComponentsQuery({ params }, { skip: !canViewRegistry });
  const relationshipsQuery = useGetRelationshipsQuery({ params }, { skip: !canViewRegistry });
  const registrantsQuery = useGetRegistrantsQuery({ params }, { skip: !canViewRegistry });

  const modelCount = canViewRegistry ? modelsQuery.data?.totalCount || 0 : 0;
  const componentCount = canViewRegistry ? componentsQuery.data?.totalCount || 0 : 0;
  const relationshipCount = canViewRegistry ? relationshipsQuery.data?.totalCount || 0 : 0;
  const registrantsConut = canViewRegistry ? registrantsQuery.data?.totalCount || 0 : 0;
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
          if (!canViewRegistry) return;
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
    [data, theme, canViewRegistry, router],
  );

  return (
    <OptionalRegistryLink enabled={canViewRegistry} href="/settings?settingsCategory=Registry">
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
    </OptionalRegistryLink>
  );
}

/**
 * Displays counts of registered MeshModels grouped by category.
 * Skips category summary queries and disables navigation unless `Keys.MesherySystemViewRegistry` is granted.
 */
function MeshModelCategories() {
  const router = useRouter();
  const canViewRegistry = useHasPermission(Keys.MesherySystemViewRegistry);
  const categoryMap = useGetCategoriesSummary({ skip: !canViewRegistry });
  const theme = useTheme();

  const cleanedData = useMemo(
    () => (canViewRegistry ? Object.keys(categoryMap).map((key) => [key, categoryMap[key]]) : []),
    [categoryMap, canViewRegistry],
  );

  const chartOptions = useMemo(
    () => ({
      data: {
        columns: cleanedData,
        colors: dataToColors(cleanedData, theme),
        type: donut(),
        onclick: function () {
          if (!canViewRegistry) return;
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
    [cleanedData, theme, canViewRegistry, router],
  );

  return (
    <OptionalRegistryLink
      enabled={canViewRegistry}
      href="/settings?settingsCategory=Registry&tab=Models"
    >
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
    </OptionalRegistryLink>
  );
}

/**
 * Renders the Dashboard MeshModel category and construct summary charts side by side.
 */
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
