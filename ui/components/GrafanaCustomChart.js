import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, IconButton } from '@material-ui/core';

import { updateProgress } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import CloseIcon from '@material-ui/icons/Close';
import dataFetch from '../lib/data-fetch';
import { withSnackbar } from 'notistack';
import dynamic from 'next/dynamic'
import { Line } from 'react-chartjs-2';
import 'chartjs-plugin-colorschemes';
// import Moment from "react-moment";
// import moment from 'moment';
if (typeof window !== 'undefined') { 
 require('chartjs-plugin-zoom');

// const zoom = dynamic(() => import('chartjs-plugin-zoom').then(mod => mod.zoom), {
//   ssr: false
// })
// import 'chartjs-plugin-streaming';
// const streaming = dynamic(() => import('chartjs-plugin-streaming'), {
//   ssr: false
// })
  require('chartjs-plugin-streaming');
}
const grafanaStyles = theme => ({
    root: {
      width: '100%',
    },
    column: {
      flexBasis: '33.33%',
    },
    heading: {
      fontSize: theme.typography.pxToRem(15),
    },
    secondaryHeading: {
      fontSize: theme.typography.pxToRem(15),
      color: theme.palette.text.secondary,
    },
    dateRangePicker: {
      display: 'flex',
      justifyContent: 'flex-end',
      marginRight: theme.spacing(1),
      marginBottom: theme.spacing(2),
    },
    iframe: {
      minHeight: theme.spacing(55),
      minWidth: theme.spacing(55),
    }
  });

const grafanaDateRangeToDate = (dt, startDate) => {
  let dto = new Date();
  switch (dt) {
    case 'now-2d':
        dto.setDate(dto.getDate() - 2);
        break;
    case 'now-7d':
        dto.setDate(dto.getDate() - 7);
        break;
    case 'now-30d':
        dto.setDate(dto.getDate() - 30);
        break;
    case 'now-90d':
        dto.setDate(dto.getDate() - 90);
        break;
    case 'now-6M':
        dto.setMonth(dto.getMonth() - 6);
        break;
    case 'now-1y':
        dto.setFullYear(dto.getFullYear() - 1);
        break;
    case 'now-2y':
        dto.setFullYear(dto.getFullYear() - 2);
        break;
    case 'now-5y':
        dto.setFullYear(dto.getFullYear() - 5);
        break;
    case 'now-1d/d':
        dto.setDate(dto.getDate() - 1);
        if(startDate){
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now-2d/d':
        dto.setDate(dto.getDate() - 2);
        if(startDate){
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now-7d/d':
        dto.setDate(dto.getDate() - 7);
        if(startDate){
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now-1w/w':
        dto.setDate(dto.getDate() - 6 - (dto.getDay() + 8) % 7);
        if(startDate){
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setDate(dto.getDate() + 6);
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now-1M/M':
        dto.setMonth(dto.getMonth() - 1);
        if(startDate){
          dto.setDate(1);
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setMonth(dto.getMonth());
          dto.setDate(0);
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now-1y/y':
        dto.setFullYear(dto.getFullYear() - 1);
        if(startDate){
          dto.setMonth(0);
          dto.setDate(1);
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setMonth(12);
          dto.setDate(0);
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now/d':
        dto.setDate(dto.getDate() - 6 - (dto.getDay() + 8) % 7);
        if(startDate){
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now':
        break;
    case 'now/w':
        dto.setDate(dto.getDate() - (dto.getDay() + 7) % 7);
        if(startDate){
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setDate(dto.getDate() + 6);
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now/M':
        if(startDate){
          dto.setDate(1);
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setMonth(dto.getMonth()+1);
          dto.setDate(0);
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now/y':
        if(startDate){
          dto.setMonth(0);
          dto.setDate(1);
          dto.setHours(0);
          dto.setMinutes(0);
          dto.setSeconds(0);
          dto.setMilliseconds(0);
        } else {
          dto.setMonth(12);
          dto.setDate(0);
          dto.setHours(23);
          dto.setMinutes(59);
          dto.setSeconds(59);
          dto.setMilliseconds(999);
        }
        break;
    case 'now-5m':
        dto.setMinutes(dto.getMinutes() - 5);
        break;
    case 'now-15m':
        dto.setMinutes(dto.getMinutes() - 15);
        break;
    case 'now-30m':
        dto.setMinutes(dto.getMinutes() - 30);
        break;
    case 'now-1h':
        dto.setHours(dto.getHours() - 1);
        break;
    case 'now-3h':
        dto.setHours(dto.getHours() - 3);
        break;
    case 'now-6h':
        dto.setHours(dto.getHours() - 6);
        break;
    case 'now-12h':
        dto.setHours(dto.getHours() - 12);
        break;
    case 'now-24h':
        dto.setHours(dto.getHours() - 24);
        break;
    default:
      return new Date(parseInt(dt.getTime().toString()));
  }
  return dto;
}

class GrafanaCustomChart extends Component {
    state = {
      intervals: [], // to store the different intervals/timeouts instances
      data: [], // data for each target
      chartData: {
        datasets: [],
        labels: [],
      },
      options: {},
    };
    constructor(props){
      super(props);
      this.timeFormat = 'MM/DD/YYYY HH:mm:ss';
    }

    componentDidMount() {
      this.collectChartData();
    }

    collectChartData = (chartInst) => {
      const { panel } = this.props;
      const self = this;
      
      panel.targets.forEach((target, ind) => {
        self.getData(ind, target, chartInst);
      });
      self.setState({options: self.createOptions(panel)});
    }

    getData = async (ind, target, chartInst) => {
      const {refresh, grafanaURL, panel, from, startDate, to, endDate, liveTail} = this.props;
      const {intervals, data, chartData} = this.state;
      
      if (grafanaURL.endsWith('/')){
        grafanaURL = grafanaURL.substring(0, grafanaURL.length - 1);
      }
      const self = this;
      // if(intervals[ind] && intervals[ind] !== null) {
      //   clearInterval(intervals[ind]);
      // }
      
      const start = Math.round(grafanaDateRangeToDate(from).getTime()/1000); //startDate.getTime()/1000);
      const end = Math.round(grafanaDateRangeToDate(to).getTime()/1000); //endDate.getTime()/1000);
      const queryURL = `${grafanaURL}/api/datasources/proxy/${panel.datasource}/api/v1/query_range?` // TODO: need to check if it is ok to use datasource name instead of ID
                +`query=${decodeURIComponent(target.expr)}&start=${start}&end=${end}&step=10`; // step 5 or 10
      dataFetch(`${queryURL}`, { 
        method: 'GET',
      }, result => {
        self.props.updateProgress({showProgress: false});
        if (typeof result !== 'undefined'){
          data[ind] = result;
          chartData.datasets[ind] = self.transformDataForChart(result, target);
          chartData.labels[ind] = target.legendFormat;
          self.setState({data, chartData});


          chartInst.update({
            preservation: true,
          });
        }
      }, self.handleError);
      
      // intervals[ind] = setInterval(fetcher, self.computeRefreshInterval(refresh) * 1000);
      // self.setState({intervals});
    }

    transformDataForChart(data, target) {
      if (data && data.status === 'success' && data.data && data.data.resultType && data.data.resultType === 'matrix' 
          && data.data.result && data.data.result.length > 0){
            const localData = {
              label: target.legendFormat,
              data: [],
              // pointStyle: 'line',
              pointRadius: 0,
              // lineTension: 0,
              // borderWidth: 2,
              fill: false,
              // type: 'line',
              // stepped: true,
              // cubicInterpolationMode: 'monotone',
              // yAxisID: target.refId,
              // stepped: true,
              // backgroundColor: 'rgba(134, 87, 167, 1)',
              // borderColor: 'rgba(134, 87, 167, 1)',
              // cubicInterpolationMode: 'monotone'
            };

            localData.data = data.data.result[0].values.map(arr => {
              const x = new Date(arr[0] * 1000);
              // const x = new Date(arr[0] * 1000).toString();
              // const x = <Moment format="LLLL">{new Date(arr[0] * 1000)}</Moment>;
              // const x = moment(new Date(arr[0] * 1000)).format(this.timeFormat);
              const y = parseInt(arr[1]);
              return {
                // x: new Date(arr[0] * 1000).toLocaleDateString() +' '+ new Date(arr[0] * 1000).toLocaleTimeString(),
                x,
                y,
              };
            })

            return localData;
      }
      return null;
    }

    createOptions(panel) {
      const {refresh, from, to} = this.props;
      const fromDate = grafanaDateRangeToDate(from);
      const toDate = grafanaDateRangeToDate(to);
      const diff = Math.floor((toDate - fromDate) / 1000);
      const refreshPeriod = this.computeRefreshInterval(refresh);
      const self = this;
      // const yAxes = [];
      // panel.targets.forEach(target => {
      //   yAxes.push({
      //     id: target.refId,
      //     type: 'linear',
      //     // ticks: {
      //     //   beginAtZero: true,
      //     // },
      //     // scaleLabel: {
      //     //   display: true,
      //     // },
      //   });
      // });
      return {
          plugins: {
            colorschemes: {
              // scheme: 'office.Office2007-2010-6'
              scheme: 'brewer.RdYlGn4'
            },
            streaming: {            // per-chart option
                frameRate: 5       // chart is drawn 30 times every second
            }
          },
          responsive: true,
          maintainAspectRatio: false,
          title: {
            display: true,
            // fontStyle: 'normal',
            text: panel.title
          },
          tooltips: {
            mode: 'nearest',
            intersect: false
          },
          hover: {
            mode: 'nearest',
            intersect: false
          },
          pan: {
            enabled: true,    // Enable panning
            mode: 'x',        // Allow panning in the x direction
            rangeMin: {
                x: null       // Min value of the delay option
            },
            rangeMax: {
                x: null       // Max value of the delay option
            }
          },
          zoom: {
              enabled: true,    // Enable zooming
              mode: 'x',        // Allow zooming in the x direction
              rangeMin: {
                  x: null       // Min value of the duration option
              },
              rangeMax: {
                  x: null       // Max value of the duration option
              }
          },
          scales: {
            xAxes: [
              {
                type: 'realtime',
                realtime: {
                  // delay: (refreshPeriod + 10) * 1000,
                  delay: 5000,
                  duration: diff * 1000,
                  refresh: refreshPeriod * 1000,
                  onRefresh: self.collectChartData,
                  // pause: true,
                }
                // // type: 'linear',
                // type: 'time',
                // time: {
                // //   parser: this.timeFormat,
                // // //   // unit: 'minute',
                // // //   // round: 'day'
                // // //   // tooltipFormat: 'll HH:mm'
                // unit: 'seconds',
                // stepSize: 1,
                // // displayFormats: {
                // //   'minute': 'HH:mm',
                // //   'hour': 'HH:mm'
                // // }
                // },
                // // distribution: 'linear',
                // distribution: 'series',
                // // bounds: 'data',
                // ticks: {
                // //   source: 'data',
                //   source: 'data',
                //   autoSkip: true
                // },
                // // scaleLabel: {
                //   // display: true,
                //   // labelString: 'Response time in ms',
                //   // ticks: {
                //   //   min: 0,
                //   //   beginAtZero: true
                //   // }
                // // }
              }
            ],
            // yAxes,
            yAxes: [{
            //       // id: target.refId,
                  type: 'linear',
                  // ticks: {
                  //   beginAtZero: true,
                  // },
            //       // scaleLabel: {
            //       //   display: true,
            //       // },
              }],
          },
        }
    }

    componentWillUnmount(){
      const {intervals} = this.state;
      intervals.forEach(interval => {
        if (interval && interval !== null){
          clearInterval(interval);
        }
      })
    }

    computeRefreshInterval = (refresh) => {
      refresh = refresh.toLowerCase();
      const l = refresh.length;
      const dur = refresh.substring(l - 1, l);
      refresh = refresh.substring(0, l - 1);
      let val = parseInt(refresh);
      switch (dur){
        case 'd':
          val *= 24;
        case 'h':
          val *= 60;
        case 'm':
          val *= 60;
        case 's':
          return val;
      }
      return 30; // fallback
    }
  
    handleError = error => {
      const self = this;
      this.props.updateProgress({showProgress: false});
      this.props.enqueueSnackbar(`There was an error communicating with Grafana`, {
        variant: 'error',
        action: (key) => (
          <IconButton
                key="close"
                aria-label="Close"
                color="inherit"
                onClick={() => self.props.closeSnackbar(key) }
              >
                <CloseIcon />
          </IconButton>
        ),
        autoHideDuration: 1000,
      });
    }
    
    render() {
      const { classes } = this.props;
      const {chartData, options} = this.state;
      if (chartData.datasets.length > 0) { 
        const filteredData = chartData.datasets.filter(x => typeof x !== 'undefined') 
        if (filteredData.length === chartData.datasets.length) {
          return (
              <NoSsr>
              {/* <React.Fragment> */}
              {/* <div className={classes.root}> */}
                <Line data={chartData} options={options} />
              {/* </div> */}
              {/* </React.Fragment> */}
              </NoSsr>
            );
        }
      }
      return null;
    }
}

GrafanaCustomChart.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  board: PropTypes.object.isRequired,
  panel: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
      updateProgress: bindActionCreators(updateProgress, dispatch),
  }
}
// const mapStateToProps = st => {
//   return null;
// }

export default withStyles(grafanaStyles)(connect(
  null,
  mapDispatchToProps
)(withSnackbar(GrafanaCustomChart)));