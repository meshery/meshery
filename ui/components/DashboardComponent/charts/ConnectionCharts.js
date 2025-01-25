import React from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { useGetAllConnectionStatusQuery } from '@/rtk-query/connection';
import { InfoOutlined } from '@material-ui/icons';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';
import { Box, IconButton, Typography, useTheme } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

export default function ConnectionStatsChart() {
  const { data: statusData } = useGetAllConnectionStatusQuery();
  const router = useRouter();
  const theme = useTheme();

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
    <UsesSistent>
      <Link
        href="/management/connections"
        style={{
          pointerEvents: !CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject)
            ? 'none'
            : 'auto',
        }}
      >
        <DashboardSection>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="h6" gutterBottom>
              Connections
            </Typography>
            <div onClick={(e) => e.stopPropagation()}>
              <CustomTextTooltip
                interactive={true}
                variant="standard"
                title={`Meshery Connections are managed and unmanaged resources that either through discovery or manual entry can be assigned to one or more Environments. [Learn More](https://docs.meshery.io/concepts/logical/connections)`}
                placement="left"
                leaveDelay={200}
              >
                <IconButton
                  disableRipple={true}
                  disableFocusRipple={true}
                  disableTouchRipple={true}
                  sx={{ padding: '0px' }}
                >
                  <InfoOutlined
                    color={theme.palette.icon.default}
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
              <ConnectCluster message={'No connections found in your clusters'} />
            )}
          </Box>
        </DashboardSection>
      </Link>
    </UsesSistent>
  );
}
