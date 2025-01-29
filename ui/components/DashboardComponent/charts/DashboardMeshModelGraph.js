import Grid from '@layer5/sistent';
import React, { useMemo } from 'react';
import BBChart from '../../BBChart';
import { donut } from 'billboard.js';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { InfoOutlined } from '@mui/icons-material';
import {
  useGetCategoriesSummary,
  useGetComponentsQuery,
  useGetMeshModelsQuery,
  useGetRelationshipsQuery,
  useGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import { DashboardSection } from '../style';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { IconButton, Typography, useTheme } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

function MeshModelContructs() {
  const params = {
    page: 0,
    pagesize: '1',
  };
  const modelCount = useGetMeshModelsQuery({ params }).data?.total_count || 0;
  const componentCount = useGetComponentsQuery({ params }).data?.total_count || 0;
  const relationshipCount = useGetRelationshipsQuery({ params }).data?.total_count || 0;
  const registrantsConut = useGetRegistrantsQuery({ params }).data?.total_count || 0;
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
        colors: dataToColors(data),
        onclick: function (d) {
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
    [data],
  );

  return (
    <UsesSistent>
      <Link
        href="/settings?settingsCategory=Registry"
        style={{
          pointerEvents: !CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject)
            ? 'none'
            : 'auto',
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
                variant="standard"
                leaveDelay={200}
                title={`The Meshery Registry is a critical component acting as the central repository for all capabilities known to Meshery. [Learn More](https://docs.meshery.io/concepts/logical/registry)`}
              >
                <IconButton disableRipple={true} disableFocusRipple={true}>
                  <InfoOutlined
                    color={theme.palette.icon.default}
                    style={{ ...iconSmall, marginLeft: '0.5rem', cursor: 'pointer' }}
                    onClick={(e) => {
                      e.stopPropagation();
                    }}
                  />
                </IconButton>
              </CustomTextTooltip>
            </div>
          </div>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <BBChart options={chartOptions} />
          </div>
        </DashboardSection>
      </Link>
    </UsesSistent>
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
        colors: dataToColors(cleanedData),
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
    [cleanedData],
  );

  return (
    <UsesSistent>
      <Link href="/settings?settingsCategory=Registry&tab=Models">
        <DashboardSection>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="h6" gutterBottom>
              Models by Category
            </Typography>

            <div onClick={(e) => e.stopPropagation()}>
              <CustomTextTooltip
                title={`Meshery Models represent the fundamental building blocks of your infrastructure. Models are categorized by their function. For example, a model for Prometheus belongs in the "Observability and Analysis" category. [Learn More](https://docs.meshery.io/concepts/logical/models)`}
                placement="left"
                variant="standard"
                interactive={true}
                leaveDelay={200}
              >
                <IconButton disableRipple={true} disableFocusRipple={true}>
                  <InfoOutlined
                    color={theme.palette.icon.default}
                    style={{ ...iconSmall, marginLeft: '0.5rem', cursor: 'pointer' }}
                    onClick={(e) => {
                      e.stopPropagation();
                    }}
                  />
                </IconButton>
              </CustomTextTooltip>
            </div>
          </div>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <BBChart options={chartOptions} />
          </div>
        </DashboardSection>
      </Link>
    </UsesSistent>
  );
}

const MeshModelGraph = () => {
  return (
    <Grid container spacing={2}>
      <Grid item xs={12} md={6}>
        <MeshModelCategories />
      </Grid>

      <Grid item xs={12} md={6}>
        <MeshModelContructs />
      </Grid>
    </Grid>
  );
};

MeshModelGraph.displayName = 'MeshModalGraph';
export default MeshModelGraph;
