import React from 'react';
import { Box, Typography, MenuItem, Select } from '@material-ui/core';
import { donut } from 'billboard.js';
import BBChart from '../../BBChart';
import { dataToColors, isValidColumnName } from '../../../utils/charts';
import ConnectClustersBtn from '../../General/ConnectClustersBtn';
import Link from 'next/link';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

export default function WorkloadChart({
  classes,
  resourses = [],
  namespaces = [],
  selectedNamespace = '',
  handleSetNamespace,
}) {
  resourses = resourses || [];
  namespaces = namespaces || [];
  const chartData = resourses
    .filter((resource) => isValidColumnName(resource?.kind))
    .map((resource) => [resource?.kind, resource?.count]);

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
      title: 'Workloads',
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
    legend: {
      show: false,
    },
  };

  return (
    <div
      className={classes.dashboardSection}
      style={{
        padding: '0.5rem',
        paddingTop: '2rem',
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Link
          href="/management/connections"
          style={{
            pointerEvents: !CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject)
              ? 'none'
              : 'auto',
          }}
        >
          <Typography variant="h6" gutterBottom className={classes.link}>
            Workloads
          </Typography>
        </Link>
        {namespaces?.length > 0 && (
          <Select value={selectedNamespace} onChange={(e) => handleSetNamespace(e.target.value)}>
            {namespaces.map((ns) => (
              <MenuItem key={ns.uniqueID} value={ns}>
                {ns}
              </MenuItem>
            ))}
          </Select>
        )}
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
              No workloads found in your cluster(s).
            </Typography>
            <ConnectClustersBtn />
          </div>
        )}
      </Box>
    </div>
  );
}
