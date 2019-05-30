import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, IconButton } from '@material-ui/core';
import { Line } from 'react-chartjs-2';
import { updateProgress } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import CloseIcon from '@material-ui/icons/Close';
import dataFetch from '../lib/data-fetch';
import { withSnackbar } from 'notistack';
import 'chartjs-plugin-colorschemes';

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

    componentDidMount() {
      this.collectChartData();
    }

    collectChartData = () => {
      const {board, panel, from, startDate, to, endDate, liveTail, refresh} = this.props;
      const self = this;
      let {grafanaURL} = this.props;
      if (grafanaURL.endsWith('/')){
        grafanaURL = grafanaURL.substring(0, grafanaURL.length - 1);
      }
      panel.targets.forEach((target, ind) => {
        const queryURL = `${grafanaURL}/api/datasources/proxy/${panel.datasource}/api/v1/query_range?` // TODO: need to check if it is ok to use datasource name instead of ID
                  +`query=${decodeURIComponent(target.expr)}&start=${Math.round(startDate.getTime()/1000)}&end=${Math.round(endDate.getTime()/1000)}&step=5`;
        self.getData(ind, queryURL, target);
      });
      self.setState({options: self.createOptions(panel)});
    }

    getData = async (ind, queryURL, target) => {
      const {refresh} = this.props;
      const {intervals, data, chartData} = this.state;
      const self = this;
      if(intervals[ind] && intervals[ind] !== null) {
        clearInterval(intervals[ind]);
      }
      const fetcher = () => {
        dataFetch(`${queryURL}`, { 
          method: 'GET',
        }, result => {
          self.props.updateProgress({showProgress: false});
          if (typeof result !== 'undefined'){
            data[ind] = result;
            chartData.datasets[ind] = self.transformDataForChart(result, target);
            chartData.labels[ind] = target.legendFormat;
            self.setState({data, chartData});
          }
        }, self.handleError);
      }
      fetcher();
      intervals[ind] = setInterval(fetcher, self.computeRefreshInterval(refresh) * 1000);
      self.setState({intervals});
    }

    transformDataForChart(data, target) {
      if (data && data.status === 'success' && data.data && data.data.resultType && data.data.resultType === 'matrix' 
          && data.data.result && data.data.result.length > 0){
            const localData = {
              label: target.legendFormat,
              data: [],
              pointStyle: 'line',
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
              return {
                x: new Date(arr[0] * 1000),
                y: arr[1],
              };
            })

            return localData;
      }
      return null;
    }

    createOptions(panel) {
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
            }
          },
          responsive: true,
          maintainAspectRatio: false,
          title: {
            // display: true,
            // fontStyle: 'normal',
            text: panel.title
          },
          scales: {
            xAxes: [
              {
                type: 'linear',
                // type: 'time',
						    distribution: 'series',
                scaleLabel: {
                  // display: true,
                  // labelString: 'Response time in ms',
                  // ticks: {
                  //   min: 0,
                  //   beginAtZero: true
                  // }
                }
              }
            ],
            // yAxes,
            yAxes: [{
            //       // id: target.refId,
                  // type: 'linear',
                  ticks: {
                    beginAtZero: true,
                  },
            //       // scaleLabel: {
            //       //   display: true,
            //       // },
              }],
          }
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
        autoHideDuration: 4000,
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