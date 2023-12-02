import React, { useEffect, useState } from 'react';
import { Box, Tooltip, Typography } from '@material-ui/core';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import { getConnectionStatusSummary } from '../../../api/connections';
import ConnectClustersBtn from '../../General/ConnectClustersBtn';
import Link from 'next/link';
import theme from '../../../themes/app';
import { iconSmall } from '../../../css/icons.styles';
import InfoIcon from '@material-ui/icons/Info';

export default function ConnectionStatsChart({ classes }) {
  const [chartData, setChartData] = useState([]);

  useEffect(() => {
    getConnectionStatusSummary().then((json) => {
      setChartData(
        json.connections_status
          .filter((data) => isValidColumnName(data.status))
          .map((data) => [data.status, data.count]),
      );
    });
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
      title: 'Connection Stats',
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

  const url = `https://docs.meshery.io/concepts/connections`;

  return (
    <Link href="/management/connections">
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Connections
          </Typography>
          <Tooltip title="Learn more about connections" placement="right">
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
                No connections found in your clusters
              </Typography>
              <ConnectClustersBtn />
            </div>
          )}
        </Box>
      </div>
    </Link>
  );
}
