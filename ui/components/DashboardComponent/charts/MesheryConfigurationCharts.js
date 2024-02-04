import React, { useEffect, useState } from 'react';
import { Box, IconButton, Typography } from '@material-ui/core';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import dataFetch from '../../../lib/data-fetch';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import CreateDesignBtn from '../../General/CreateDesignBtn';
import theme from '../../../themes/app';
import { iconSmall } from '../../../css/icons.styles';
import {
  CustomTextTooltip,
  renderTooltipContent,
} from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { InfoOutlined } from '@material-ui/icons';

const ACTION_TYPES = {
  FETCH_PATTERNS: {
    name: 'FETCH_PATTERNS',
    error_msg: 'Failed to fetch designs',
  },
  FETCH_FILTERS: {
    name: 'FETCH_FILTERS',
    error_msg: 'Failed to fetch WASM filters',
  },
};

export default function MesheryConfigurationChart({ classes }) {
  const [chartData, setChartData] = useState([]);

  const { notify } = useNotification();

  const handleError = (action) => (error) => {
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
    });
  };

  function fetchDesigns() {
    dataFetch(
      `/api/pattern`,
      { credentials: 'include' },
      (result) => {
        if (result) {
          setChartData((prevData) => [...prevData, ['Designs', result.total_count]]);
        }
      },
      handleError(ACTION_TYPES.FETCH_PATTERNS),
    );
  }

  function fetchFilters() {
    dataFetch(
      `/api/filter`,
      { credentials: 'include' },
      (result) => {
        if (result) {
          setChartData((prevData) => [...prevData, ['Filters', result.total_count]]);
        }
      },
      handleError(ACTION_TYPES.FETCH_FILTERS),
    );
  }

  useEffect(() => {
    fetchDesigns();
    fetchFilters();
  }, []);

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
              title={renderTooltipContent({
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
