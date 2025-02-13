import React, { useEffect } from 'react';
import { Box, styled, useTheme } from '@layer5/sistent';
import bb, { gauge } from 'billboard.js';

import { NoSsr } from '@mui/material';

const ChartRoot = styled(Box)(() => ({
  width: '100%',
  height: '75%',
  minHeight: '18rem',
  '& .bb-chart-arcs-background': {
    fill: '#e0e0e0',
    stroke: 'none',
  },
}));

const ErrorMessage = styled(Box)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.error.main,
    width: '100%',
    textAlign: 'center',
    fontSize: '12px',
    fontWeight: 'bold',
  };
});

export default function GrafanaCustomGaugeChart(props) {
  let chartRef = null;

  const configChartData = () => {
    const { panel, data } = props;

    let units = '';
    if (panel.format) {
      if (panel.format.startsWith('percent')) {
        units = '%';
      } else {
        units = ` ${panel.format}`;
      }
    }

    let min = 0;
    let max = 100;
    if (panel.gauge) {
      if (panel.gauge.minValue) min = panel.gauge.minValue;
      if (panel.gauge.maxValue) max = panel.gauge.maxValue;
    }

    let colors = [];
    if (panel.colors) {
      colors = panel.colors;
    }

    let thresholds = [];
    if (panel.thresholds) {
      thresholds = panel.thresholds.split(',').map((t) => parseFloat(t.trim()));
    }

    let gdata = 0;
    let glabel = '';
    if (data && data.length > 0) {
      const dlind = data[0].length - 1;
      gdata = data[0][dlind] ? data[0][dlind] : 0;
      glabel = data[0][0];
    }

    if (chartRef && chartRef !== null) {
      self.chart = bb.generate({
        bindto: chartRef,
        data: {
          columns: [[glabel, gdata]],
          type: gauge(),
        },
        gauge: {
          min,
          max,
          label: {
            format(value) {
              return value + units;
            },
            extents() {
              return '';
            },
          },
        },
        color: {
          pattern: colors,
          threshold: {
            values: thresholds,
          },
        },
        legend: { show: false },
        tooltip: { show: false },
      });
    }
  };

  useEffect(() => {
    configChartData();
  });

  const { error } = props;

  return (
    <NoSsr>
      <Box>
        {/* <ChartTitle>{props.panel?.title}</ChartTitle> */}
        <ErrorMessage>{error && 'There was an error communicating with the server'}</ErrorMessage>
        <ChartRoot ref={(ch) => (chartRef = ch)} />
      </Box>
    </NoSsr>
  );
}
