import { Component, useEffect } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, IconButton } from '@material-ui/core';
import { updateProgress } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import dataFetch from '../lib/data-fetch';
import { withSnackbar } from 'notistack';
import makeStyles from '@material-ui/styles/makeStyles';

let c3;
if (typeof window !== 'undefined') { 
  c3 = require('c3');
}

const useStyles = makeStyles({
  '@global': {
    '.c3-chart-arcs-background': {
      fill:'#e0e0e0',
      stroke: 'none',
    },
  },
  root: {
    width: '100%',
    height: '75%',
  },
  error: {
    color: '#D32F2F',
    width: '100%',
    textAlign: 'center',
    fontSize: '12px',
    // fontFamily: 'Helvetica Nueue',
    fontWeight: 'bold',
  },
  title: {
    fontSize: '12px',
    color: '#666666',
    // fontFamily: 'Helvetica Nueue',
    fontWeight: 'bold',
    textAlign: 'center',
    width: '100%',
  }
});

export default function GrafanaCustomGaugeChart(props) {
    let chartRef = null;
    const configChartData = () => {
      const { panel, data } = props;
      
      let units = '';
      if(panel.format){
        if(panel.format.startsWith('percent')){
          units = ' %';
        } else {
          units = ` ${panel.format}`;
        }
      }
      let min = 0, max = 100;
      if(panel.gauge){
        if(panel.gauge.minValue) min = panel.gauge.minValue;
        if(panel.gauge.maxValue) max = panel.gauge.maxValue;
      }
      let colors = [];
      if(panel.colors){
        colors = panel.colors;
      }
      let thresholds = [];
      if(panel.thresholds){
        thresholds = panel.thresholds.split(',').map(t => parseFloat(t.trim()));
      }

      let gdata = 0, glabel = '';
      if(data.datasets && data.datasets.length > 0 && 
          data.datasets[0].data && data.datasets[0].data.length > 0){
            const dlind = data.datasets[0].data.length - 1;
            gdata = data.datasets[0].data[dlind] && data.datasets[0].data[dlind].y?data.datasets[0].data[dlind].y:0;
            glabel = data.datasets[0].label;
          }

      const c3Chart = c3.generate({
        oninit: function(args){
          console.log(JSON.stringify(args));
        },
        bindto: chartRef,
        data: {
          columns: [
            [glabel, gdata],
          ],
          type: 'gauge',
        },
        gauge: {
             min,
             max,
             units,
             label: {
              show: glabel && glabel !== '',
             },
          //    width: 39 // for adjusting arc thickness
        },
        color: {
            pattern: colors, // the three color levels for the percentage values.
            threshold: {
    //            unit: 'value', // percentage is default
    //            max: 200, // 100 is default
                values: thresholds,
            }
        },
        legend: {
          show: false,
        },
        tooltip: {
          show: false
        },
        // size: {
        //   height: '100%',
        // }
      });
    }

    useEffect(() => {
      configChartData();
    });
    const { panel, error } = props;
    const classes = useStyles();
    
    // const {chartData, options} = this.state;
    return (
      <NoSsr>
        <div className={classes.title}>{panel.title}</div>
        <div className={classes.error}>{error && 'There was an error communicating with the server'}</div>
        <div ref={ch => chartRef = ch} className={classes.root}>
        </div>
      </NoSsr>
    );
    
}