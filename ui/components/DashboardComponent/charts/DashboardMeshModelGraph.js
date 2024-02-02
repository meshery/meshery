import Grid from '@material-ui/core/Grid';
import React, { useEffect, useMemo, useState } from 'react';
import { IconButton, Typography } from '@material-ui/core';
import BBChart from '../../BBChart';
import { donut, pie } from 'billboard.js';
import {
  getAllComponents,
  getMeshModels,
  getRelationshipsDetail,
  fetchCategories,
  getModelFromCategoryApi,
} from '../../../api/meshmodel';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import theme from '../../../themes/app';
import { iconSmall } from '../../../css/icons.styles';
import {
  CustomTextTooltip,
  renderTooltipContent,
} from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { InfoOutlined } from '@material-ui/icons';

const useFetchTotal = (fetchr) => {
  const [total, setTotal] = useState(0);

  useEffect(() => {
    fetchr()
      .then((json) => {
        setTotal(json['total_count']);
      })
      .catch((e) => console.log('Api Error : ', e));
  }, []);
  return total;
};

function MeshModelContructs({ classes }) {
  // API Calls
  const totalModels = useFetchTotal(() => getMeshModels(1, 1));
  const totalComponents = useFetchTotal(() => getAllComponents(1, 1));
  const totalRelationships = useFetchTotal(() => getRelationshipsDetail(1, 1));

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
              title={renderTooltipContent({
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

  // API Calls
  useEffect(() => {
    fetchCategories().then((categoriesJson) => {
      categoriesJson['categories'].forEach((category) => {
        let categoryName = category.name;
        getModelFromCategoryApi(categoryName).then((modelsJson) => {
          setCategoryMap((prevState) => ({
            ...prevState,
            [categoryName]: modelsJson['total_count'],
          }));
        });
      });
    });
  }, []);

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
              title={renderTooltipContent({
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
