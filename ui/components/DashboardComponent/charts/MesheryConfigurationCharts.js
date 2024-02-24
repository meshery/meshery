import React, { useEffect, useState } from 'react';
import { Box, IconButton, Typography } from '@material-ui/core';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
// import { useNotification } from '../../../utils/hooks/useNotification';
import CreateDesignBtn from '../../General/CreateDesignBtn';
import theme from '../../../themes/app';
import { iconSmall } from '../../../css/icons.styles';
import {
  CustomTextTooltip,
  RenderTooltipContent,
} from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { InfoOutlined } from '@material-ui/icons';
import { useGetPatternsQuery } from '@/rtk-query/design';
import { useGetFiltersQuery } from '@/rtk-query/filter';

export default function MesheryConfigurationChart({ classes }) {
  // const { notify } = useNotification();

  const [chartData, setChartData] = useState([]);

  const { data: patternsData, error: patternsError } = useGetPatternsQuery({
    page: 0,
    pagesize: 25,
  });

  const { data: filtersData, error: filtersError } = useGetFiltersQuery({
    page: 0,
    pagesize: 25,
  });

  useEffect(() => {
    if (!patternsError && patternsData?.patterns) {
      setChartData((prevData) => [...prevData, ['Designs', patternsData.total_count]]);
    }
  }, [patternsData, patternsError]);

  useEffect(() => {
    if (!filtersError && filtersData?.filters) {
      setChartData((prevData) => [...prevData, ['Filters', filtersData.total_count]]);
    }
  }, [filtersData, filtersError]);

  const chartOptions = {
    data: {
      columns: chartData,
      type: donut(),
      colors: dataToColors(chartData),
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
  };

  const url = `https://docs.meshery.io/guides/configuration-management`;

  return (
    <Link href="/configuration/designs">
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Configuration
          </Typography>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              backgroundColor="#3C494F"
              placement="left"
              interactive={true}
              title={RenderTooltipContent({
                showPriortext: 'Meshery’s ability to configure infrastructure and applications.',
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
            <div
              style={{
                padding: '2rem',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                flexDirection: 'column',
              }}
            >
              <Typography style={{ fontSize: '1.5rem', marginBottom: '1rem' }} align="center">
                No Meshery configuration found
              </Typography>
              <CreateDesignBtn />
            </div>
          )}
        </Box>
      </div>
    </Link>
  );
}
