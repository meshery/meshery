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
import { Line } from 'react-chartjs-2';
import 'chartjs-plugin-colorschemes';
// import 'chartjs-plugin-deferred';
import moment from 'moment';
if (typeof window !== 'undefined') { 
  require('chartjs-plugin-zoom');
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
      // return new Date(parseInt(dt.getTime().toString()));
      return new Date(parseFloat(dt));
  }
  return dto;
}

class GrafanaCustomChart extends Component {
    
    constructor(props){
      super(props);
      this.timeFormat = 'MM/DD/YYYY HH:mm:ss';
      this.datasetIndex = {};
      this.state = {
        data: [], // data for each target
        chartData: {
          datasets: [],
          labels: [],
        },
        options: {},
      };
    }

    componentDidMount() {
      console.log(`chart did mount`);
      this.configChartData();
    }

    configChartData = () => {
      const { panel } = this.props;
      // const {chartData} = this.state;
      const self = this;
      // let trackInitialChange = false;
      panel.targets.forEach((target, ind) => {
      //   if (typeof chartData.labels[ind] === 'undefined' || typeof chartData.datasets[ind] === 'undefined'){
      //     chartData.labels[ind] = target.legendFormat;
      //     chartData.datasets[ind] = {
      //       label: target.legendFormat,
      //       data: [],
      //       pointRadius: 0,
      //       fill: false,
      //     };
      //     trackInitialChange = true;
      //   }
        self.datasetIndex[`${ind}_0`] = ind;
      });
      // if(trackInitialChange){
      //   self.setState({chartData, options: self.createOptions(panel)});
      // } else {
      self.setState({options: self.createOptions()});
      // }
      self.collectChartData();
    }

    getOrCreateIndex(datasetInd) {
      if(typeof this.datasetIndex[datasetInd] !== 'undefined'){
        return this.datasetIndex[datasetInd];
      }
      let max = 0;
      Object.keys(this.datasetIndex).forEach(i => {
        if(this.datasetIndex[i] > max){
          max = this.datasetIndex[i];
        }
      });
      this.datasetIndex[datasetInd] = max+1;
      return max+1;
    }

    collectChartData = (chartInst) => {
      const { panel } = this.props;
      const self = this;
      panel.targets.forEach((target, ind) => {
        self.getData(ind, target, chartInst);
      });
    }

    getData = async (ind, target, chartInst) => {
      const {grafanaURL, panel, from, to, templateVars, liveTail} = this.props;
      const {data, chartData} = this.state;

      const cd = (typeof chartInst === 'undefined'?chartData:chartInst.data);
      
      if (grafanaURL.endsWith('/')){
        grafanaURL = grafanaURL.substring(0, grafanaURL.length - 1);
      }
      const self = this;
      let expr = target.expr;
      templateVars.forEach(tv => {
        const tvrs = tv.split('=');
        if (tvrs.length == 2){
          expr = expr.replace(`$${tvrs[0]}`,tvrs[1]);
        }
      });
      
      const start = Math.round(grafanaDateRangeToDate(from).getTime()/1000);
      const end = Math.round(grafanaDateRangeToDate(to).getTime()/1000);
      const queryURL = `${grafanaURL}/api/datasources/proxy/${panel.datasource}/api/v1/query_range?` // TODO: need to check if it is ok to use datasource name instead of ID
                +`query=${decodeURIComponent(expr)}&start=${start}&end=${end}&step=10`; // step 5 or 10
      dataFetch(`${queryURL}`, { 
        method: 'GET',
      }, result => {
        self.props.updateProgress({showProgress: false});
        if (typeof result !== 'undefined'){
          // data[ind] = result;
          
          // chartData.datasets[ind] = self.transformDataForChart(result, target);
          // // chartData.labels[ind] = target.legendFormat;
          // self.setState({data, chartData});
          const fullData = self.transformDataForChart(result, target);
          fullData.forEach(({metric, data}, di) => {
            const datasetInd = self.getOrCreateIndex(`${ind}_${di}`);
            const newData = [];

            if (typeof cd.labels[datasetInd] === 'undefined' || typeof cd.datasets[datasetInd] === 'undefined'){
              let legend = target.legendFormat;
              Object.keys(metric).forEach(metricKey => {
                legend = legend.replace(`{{${metricKey}}}`, metric[metricKey])
                            .replace(`{{ ${metricKey} }}`, metric[metricKey]);
              });
              legend = legend.replace(`{{ `, '')
                          .replace(` }}`, '');

              cd.labels[datasetInd] = legend;
              cd.datasets[datasetInd] = {
                label: legend,
                data: [],
                pointRadius: 0,
                fill: false,
              };
            }



            data.forEach(({x, y}) => {
              let toadd = true;
              cd.datasets[datasetInd].data.forEach(({x: x1, y: y1}) => {
                if(x === x1) {
                  toadd = false;
                }
              });
              if(toadd){
                newData.push({x, y});  
              }
            });
            Array.prototype.push.apply(cd.datasets[datasetInd].data, newData);
            cd.datasets[datasetInd].data.sort((a, b) => {
              return new Date(a.x).getTime() - new Date(b.x).getTime();
            })
          });
          console.log(`current from: ${from} to: ${to} chart labels: ${JSON.stringify(cd.labels)}`);
          if(typeof chartInst === 'undefined'){
            for(let cddi=0;cddi < cd.datasets.length; cddi++){
              if(typeof cd.datasets[cddi] === 'undefined'){
                cd.datasets[cddi] = {data:[], label: ''};
              }
            }
            self.setState({chartData});
          } else {
            chartInst.update({
              preservation: true,
            });
            if(!liveTail) {
              chartInst.options.scales.xAxes[0].realtime.pause = true
            }
            // const min = chartInst.scales["x-axis-0"].min;
            // const max = chartInst.scales["x-axis-0"].max;
            // if(isNaN(min) || isNaN(max)){
            //   // chartInst.scales["x-axis-0"].min = grafanaDateRangeToDate(from).getTime();
            //   // chartInst.scales["x-axis-0"].max = grafanaDateRangeToDate(to).getTime();
            //   chartInst.reset();
            // }
          }
        }
      }, self.handleError);
    }

    transformDataForChart(data, target) {
      if (data && data.status === 'success' && data.data && data.data.resultType && data.data.resultType === 'matrix' 
          && data.data.result && data.data.result.length > 0){
            // const localData = {
            //   label: target.legendFormat,
            //   data: [],
            //   // pointStyle: 'line',
            //   pointRadius: 0,
            //   // lineTension: 0,
            //   // borderWidth: 2,
            //   fill: false,
            //   // type: 'line',
            //   // stepped: true,
            //   // cubicInterpolationMode: 'monotone',
            //   // yAxisID: target.refId,
            //   // stepped: true,
            //   // backgroundColor: 'rgba(134, 87, 167, 1)',
            //   // borderColor: 'rgba(134, 87, 167, 1)',
            //   // cubicInterpolationMode: 'monotone'
            // };
            let fullData = [];
            data.data.result.forEach(r => {
              const localData = r.values.map(arr => {
                const x = moment(arr[0] * 1000).format(this.timeFormat);
                const y = parseInt(arr[1]);
                return {
                  x,
                  y,
                };
              })
              fullData.push({
                data: localData,
                metric: r.metric,
              })
            })
            // localData.data = data.data.result[0].values.map(arr => {
            

            return fullData;
      }
      return [];
    }

    updateDateRange(){
      const self = this;
      let tm;
      return  function({chart}){
        // from={from} startDate={startDate} to={to} endDate={endDate} liveTail={liveTail}  refresh={refresh}
        // this.props.updateDateRange(this.props.from, this.props.startDate, 'now', this.props.endDate, event.target.checked, this.props.refresh);  
        // also we should pause the chart. . . probably NOT
        if (typeof tm !== 'undefined'){
          clearTimeout(tm);
        }
        if(typeof chart !== 'undefined'){
          const min = chart.scales["x-axis-0"].min;
          const max = chart.scales["x-axis-0"].max;
          tm = setTimeout(function(){
            if(!isNaN(min) && !isNaN(max)){
              console.log(`zoom/pan - processing min: ${min} and max: ${max}`);
              self.props.updateDateRange(`${min}`, new Date(min), `${max}`, new Date(max), false, self.props.refresh);
            } else {
              // self.props.updateDateRange(self.props., new Date(min), `${max}`, new Date(max), false, self.props.refresh);
              self.props.updateDateRange(self.props.from, self.props.startDate, self.props.to, self.props.endDate, self.props.liveTail, self.props.refresh);  
            }
          }, 1000);
        }
        return false;
      }
    }

    createOptions() {
      const {panel, refresh, from, to, liveTail} = this.props;
      const fromDate = grafanaDateRangeToDate(from);
      const toDate = grafanaDateRangeToDate(to);
      const diff = Math.floor((toDate - fromDate) / 1000);
      const refreshPeriod = this.computeRefreshInterval(refresh);
      const self = this;

      let xAx;
      // if (liveTail){
        xAx = {
          type: 'realtime',
          realtime: {
            // delay: (refreshPeriod + 10) * 1000,
            delay: 5000,
            duration: diff * 1000,
            refresh: refreshPeriod * 1000,
            onRefresh: self.collectChartData,
            pause: false,
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
        };
      // } 
      let streaming = {
          frameRate: 5,
      }
      if (!liveTail){
        streaming = false;
      }

      

      return {
          plugins: {
            // deferred: {
            //   xOffset: 150,
            //   yOffset: '50%',
            //   delay: 500
            // },
            colorschemes: {
              // scheme: 'office.Office2007-2010-6'
              scheme: 'brewer.RdYlGn4'
            },
            streaming: streaming,
          },
          responsive: true,
          maintainAspectRatio: false,
          title: {
            display: true,
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
            enabled: true,
            mode: 'x',
            rangeMin: {
                x: null
            },
            rangeMax: {
                x: null
            },
            onPan: self.updateDateRange(),
          },
          zoom: {
              enabled: true,
              // drag: true, // if set to true will turn off pinch
              mode: 'x',
              speed: 0.05,
              rangeMin: {
                  x: null
              },
              rangeMax: {
                  x: null
              },
              onZoom: self.updateDateRange(),
          },
          scales: {
            xAxes: [
              xAx,
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
      console.log(`chart will unmount`);
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
      // this.props.enqueueSnackbar(`There was an error communicating with Grafana`, {
      //   variant: 'error',
      //   action: (key) => (
      //     <IconButton
      //           key="close"
      //           aria-label="Close"
      //           color="inherit"
      //           onClick={() => self.props.closeSnackbar(key) }
      //         >
      //           <CloseIcon />
      //     </IconButton>
      //   ),
      //   autoHideDuration: 1000,
      // });
    }
    
    render() {
      const { classes } = this.props;
      const {chartData, options} = this.state;
      let finalChartData = {
        datasets: [],
        labels: [],
      }
      const filteredData = chartData.datasets.filter(x => typeof x !== 'undefined')
      if(chartData.datasets.length === filteredData.length){
        finalChartData = chartData;
      }
      return (
        <NoSsr>
          <Line data={finalChartData} options={options} />
        </NoSsr>
      );
      // return null;
    }
}

GrafanaCustomChart.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  board: PropTypes.object.isRequired,
  panel: PropTypes.object.isRequired,
  templateVars: PropTypes.array.isRequired,
  updateDateRange: PropTypes.func.isRequired,
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