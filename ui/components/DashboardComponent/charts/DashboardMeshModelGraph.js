import Grid from '@material-ui/core/Grid';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { IconButton, Typography } from '@material-ui/core';
import BBChart from '../../BBChart';
import { donut, pie } from 'billboard.js';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import theme from '../../../themes/app';
import { iconSmall } from '../../../css/icons.styles';
import {
  CustomTextTooltip,
  RenderTooltipContent,
} from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { InfoOutlined } from '@material-ui/icons';
import {
  useGetModelCategoriesQuery,
  useLazyGetComponentsQuery,
  useLazyGetMeshModelsQuery,
  useLazyGetModelFromCategoryQuery,
  useLazyGetRelationshipsQuery,
} from '@/rtk-query/meshModel';

function MeshModelContructs({ classes }) {
  const [getAllModels] = useLazyGetMeshModelsQuery();
  const [getAllComponents] = useLazyGetComponentsQuery();
  const [getAllRelationships] = useLazyGetRelationshipsQuery();

  // States to hold total counts
  const [totalModels, setTotalModels] = useState(0);
  const [totalComponents, setTotalComponents] = useState(0);
  const [totalRelationships, setTotalRelationships] = useState(0);

  // Fetch data and update state on component mount
  const fetchData = useCallback(async () => {
    try {
      const models = await getAllModels({
        page: 1,
        pagesize: 'all',
      });
      const components = await getAllComponents({
        page: 1,
        pagesize: 'all',
      });
      const relationships = await getAllRelationships({
        page: 1,
        pagesize: 'all',
      });

      setTotalModels(models.data.total_count);
      setTotalComponents(components.data.total_count);
      setTotalRelationships(relationships.data.total_count);
    } catch (error) {
      console.error('Error fetching Mesh Models data:', error);
    }
  }, [getAllModels, getAllComponents, getAllRelationships]);

  useEffect(() => {
    fetchData();
  }, []);

  // Data Cleanup
  const data = useMemo(() => {
    // TODO: Add Policies
    return [
      ['Models', totalModels],
      ['Components', totalComponents],
      ['Relationships', totalRelationships],
      // TODO: Add Policies
    ];
  }, [totalModels, totalRelationships, totalComponents]);

  const chartOptions = useMemo(
    () => ({
      data: {
        columns: data,
        type: donut(),
        colors: dataToColors(data),
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

  const url = `https://docs.meshery.io/concepts/logical/models`;

  return (
    <Link href="/settings#registry">
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Registry
          </Typography>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              backgroundColor="#3C494F"
              placement="left"
              interactive={true}
              title={RenderTooltipContent({
                showPriortext:
                  'Meshery uses a set of resource models to define concrete boundaries to ensure extensible and sustainable management.',
                showAftertext: 'to learn more about Models, Components, and Relationships',
                link: url,
              })}
            >
              <IconButton disableRipple={true} disableFocusRipple={true}>
                <InfoOutlined
                  color={theme.palette.secondary.iconMain}
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
      </div>
    </Link>
  );
}

function MeshModelCategories({ classes }) {
  const [categoryMap, setCategoryMap] = useState({});
  const { data: categories } = useGetModelCategoriesQuery();
  const [getModelFromCategory] = useLazyGetModelFromCategoryQuery();

  useEffect(() => {
    const fetchModelsForCategories = async () => {
      if (categories) {
        const updatedCategoryMap = { ...categoryMap };
        for (const category of categories.categories) {
          const categoryName = category.name;
          if (!updatedCategoryMap[categoryName]) {
            const { data: models } = await getModelFromCategory({
              page: 1,
              pagesize: 'all',
              category: categoryName,
            });
            updatedCategoryMap[categoryName] = models?.total_count || 0;
          }
        }
        setCategoryMap(updatedCategoryMap);
      }
    };

    fetchModelsForCategories();
  }, [categories]);

  const cleanedData = useMemo(
    () => Object.keys(categoryMap).map((key) => [key, categoryMap[key]]),
    [categoryMap],
  );

  const chartOptions = useMemo(
    () => ({
      data: {
        columns: cleanedData,
        colors: dataToColors(cleanedData),
        type: pie(),
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

  const url = `https://docs.meshery.io/concepts/logical/models`;

  return (
    <Link href="/settings#registry">
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Categories
          </Typography>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              backgroundColor="#3C494F"
              title={RenderTooltipContent({
                showPriortext:
                  'Each Model corresponds to a category, so the category shows the high-level use case of that model, e.g., prometheus is under “Observability and Analysis category.',
                showAftertext: 'to learn more about all Categories',
                link: url,
              })}
              placement="left"
              interactive={true}
            >
              <IconButton disableRipple={true} disableFocusRipple={true}>
                <InfoOutlined
                  color={theme.palette.secondary.iconMain}
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
      </div>
    </Link>
  );
}

const MeshModelGraph = ({ classes }) => {
  return (
    <Grid container spacing={2}>
      <Grid item xs={12} md={6}>
        <MeshModelCategories classes={classes} />
      </Grid>

      <Grid item xs={12} md={6}>
        <MeshModelContructs classes={classes} />
      </Grid>
    </Grid>
  );
};

MeshModelGraph.displayName = 'MeshModalGraph';
export default MeshModelGraph;
