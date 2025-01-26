import React, { useEffect, useState } from 'react';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors } from '../../../utils/charts';
import Link from 'next/link';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { InfoOutlined } from '@material-ui/icons';
import { useGetPatternsQuery } from '@/rtk-query/design';
import { useGetFiltersQuery } from '@/rtk-query/filter';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useRouter } from 'next/router';
import { DashboardSection } from '../style';
import ConnectCluster from './ConnectCluster';
import { UsesSistent } from '@/components/SistentWrapper';
import { Box, IconButton, Typography, useTheme } from '@layer5/sistent';

export default function MesheryConfigurationChart() {
  const router = useRouter();
  const [chartData, setChartData] = useState([]);
  const theme = useTheme();

  const { data: patternsData, error: patternsError } = useGetPatternsQuery({
    page: 0,
    pagesize: 1,
  });

  const { data: filtersData, error: filtersError } = useGetFiltersQuery({
    page: 0,
    pagesize: 1,
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
      onclick: function (d) {
        const routeName = d.name.charAt(0).toLowerCase() + d.name.slice(1);
        router.push(`/configuration/${routeName}`);
      },
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

  return (
    <UsesSistent>
      <Link
        href="/configuration/designs"
        style={{
          pointerEvents: !CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject)
            ? 'none'
            : 'auto',
        }}
      >
        <DashboardSection>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="h6" gutterBottom>
              Configuration
            </Typography>
            <div onClick={(e) => e.stopPropagation()}>
              <CustomTextTooltip
                placement="left"
                interactive={true}
                variant="standard"
                leaveDelay={200}
                title={`Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured. [Learn more](https://docs.meshery.io/concepts/logical/designs)`}
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
