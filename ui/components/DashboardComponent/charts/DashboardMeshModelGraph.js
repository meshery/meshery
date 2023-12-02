import Grid from '@material-ui/core/Grid';
import React, { useEffect, useMemo, useState } from 'react';
import { Typography, Tooltip } from '@material-ui/core';
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
import InfoIcon from '@material-ui/icons/Info';

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
        title: 'Capabilities\nRegistry',
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

  const url = `https://docs.meshery.io/concepts/models#models`;

  return (
    <Link href="/settings#registry">
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Models
          </Typography>
          <Tooltip title="Learn more about models" placement="right">
            <InfoIcon
              color={theme.palette.secondary.iconMain}
              style={{ ...iconSmall, marginLeft: '0.5rem' }}
              onClick={(e) => {
                e.stopPropagation();
                window.open(url, '_blank');
              }}
            />
          </Tooltip>
        </div>
        <div>
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

  const url = `https://docs.meshery.io/concepts/models`;

  return (
    <Link href="/settings#registry">
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Model Categories
          </Typography>
          <Tooltip title="Learn more about model categories" placement="right">
            <InfoIcon
              color={theme.palette.secondary.iconMain}
              style={{ ...iconSmall, marginLeft: '0.5rem' }}
              onClick={(e) => {
                e.stopPropagation();
                window.open(url, '_blank');
              }}
            />
          </Tooltip>
        </div>
        <div>
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
