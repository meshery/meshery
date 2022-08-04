import React, { useEffect } from 'react';
import { NoSsr } from '@material-ui/core';
import makeStyles from '@material-ui/styles/makeStyles';

import bb, { gauge } from 'billboard.js'

const useStyles = makeStyles({
  '@global' : { '.bb-chart-arcs-background' : { fill : '#e0e0e0',
    stroke : 'none', }, },
  grRoot : { width : '100%',
    height : '75%',
    minHeight : '18rem', },
  error : {
    color : '#D32F2F',
    width : '100%',
    textAlign : 'center',
    fontSize : '12px',
    // fontFamily: 'Helvetica Nueue',
    fontWeight : 'bold',
  },
  title : {
    fontSize : '12px',
    color : '#666666',
    // fontFamily: 'Helvetica Nueue',
    fontWeight : 'bold',
    textAlign : 'center',
    width : '100%',
  },
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
    let min = 0; let max = 100;
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

    let gdata = 0; let glabel = '';
    if (data && data.length > 0) {
      const dlind = data[0].length - 1;
      gdata = data[0][dlind]
        ? data[0][dlind]
        : 0;
      glabel = data[0][0];
    }

    if (chartRef && chartRef !== null) {
      self.chart = bb.generate({
        // oninit: function(args){
        //   console.log(JSON.stringify(args));
        // },
        bindto : chartRef,
        data : { columns : [
          [
            glabel,
            gdata,
          ],
        ],
        type : gauge(), },
        gauge : { min,
          max,
          // units,
          label : {
            // show: glabel && glabel !== '',
            format(value) {
              return value + units;
            },
            extents() {
              // return (isMax ? "Max:" : "Min:") + value;
              return '';
            }, },
          //    width: 39 // for adjusting arc thickness
        },
        color : { pattern : colors, // the three color levels for the percentage values.
          threshold : {
            //            unit: 'value', // percentage is default
            //            max: 200, // 100 is default
            values : thresholds, }, },
        legend : { show : false, },
        tooltip : { show : false, },
        // size: {
        //   height: '100%',
        // }
      });
    }
  };

  useEffect(() => {
    configChartData();
  });
  const { error } = props;
  const classes = useStyles();

  // const {chartData, options} = this.state;
  return (
    <NoSsr>
      {/* <div className={classes.title}>{panel.title}</div> */}
      <div className={classes.error}>{error && 'There was an error communicating with the server'}</div>
      <div ref={(ch) => chartRef = ch} className={classes.grRoot} />
    </NoSsr>
  );
}
