import React from 'react';
import { Box, Typography, IconButton } from '@material-ui/core';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import ConnectClustersBtn from '../../General/ConnectClustersBtn';
import Link from 'next/link';
import theme from '../../../themes/app';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { useGetAllConnectionStatusQuery } from '@/rtk-query/connection';
import { InfoOutlined } from '@material-ui/icons';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';

export default function ConnectionStatsChart({ classes }) {
  const { data: statusData } = useGetAllConnectionStatusQuery();
  const router = useRouter();

  const chartData =
    statusData?.connections_status
      ?.filter((data) => isValidColumnName(data.status))
      .map((data) => [data.status, data.count]) || [];

  const chartOptions = {
    data: {
      columns: chartData,
      type: donut(),
      colors: dataToColors(chartData),
      onclick: function () {
        router.push('/management/connections');
      },
    },
    arc: {
      cornerRadius: {
        ratio: 0.05,
      },
    },
    donut: {
      title: 'Connections\n by Status',
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

  return (
    <Link
      href="/management/connections"
      style={{
        pointerEvents: !CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject)
          ? 'none'
          : 'auto',
      }}
    >
      <div className={classes.dashboardSection}>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" gutterBottom className={classes.link}>
            Connections
          </Typography>
          <div onClick={(e) => e.stopPropagation()}>
            <CustomTextTooltip
              interactive={true}
              variant="standard"
              title={`Meshery Connections are managed and unmanaged resources that either through discovery or manual entry can be assigned to one or more Environments. [Learn More](https://docs.meshery.io/concepts/logical/connections)`}
              placement="left"
            >
              <IconButton
                disableRipple={true}
                disableFocusRipple={true}
                disableTouchRipple={true}
                sx={{ padding: '0px' }}
              >
                <InfoOutlined
                  color={theme.palette.secondary.iconMain}
                  style={{ ...iconSmall, marginLeft: '0.5rem', cursor: 'pointer' }}
                  onClick={(e) => e.stopPropagation()}
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
