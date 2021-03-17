import React from 'react';
import { NoSsr } from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import {
  fortioResultToJsChartData, makeChart, makeOverlayChart, makeMultiChart,
} from '../lib/chartjs-formatter';

import bb, {areaStep, line} from 'billboard.js'

const styles = (theme) => ({
  title: {
    textAlign: 'center',
    fontSize: theme.spacing(1.75),
    marginBottom: theme.spacing(1)
  },
});

class MesheryChart extends React.Component {
  constructor(props) {
    super(props);
    this.chartRef = null;
    this.chart = null;
  }

  singleChart = (data) => {
    if (typeof data === 'undefined' || typeof data.StartTime === 'undefined') {
      return {};
    }
    return makeChart(fortioResultToJsChartData(data));
  }

  processChartData(chartData) {
    const self = this;
    if (self.chartRef && self.chartRef !== null) {
      if (chartData && chartData.data && chartData.options) {
        const xAxes = [];
        const yAxes = [];
        const colors = [];
        const types = {};
        const axes = {};
        const axis = {};
        const yAxisTracker = {};
        const xAxisTracker = {};

        if (chartData.data && chartData.data.datasets) {
          chartData.data.datasets.forEach((ds, ind) => {
            // xAxis.push('x');
            const yAxis = [ds.label];
            const xAxis = [`x${ind + 1}`];
            xAxisTracker[ds.label] = `x${ind + 1}`;
            yAxisTracker[ds.yAxisID] = `y${ind > 0 ? ind + 1 : ''}`;
            axes[ds.label] = `y${ind > 0 ? ind + 1 : ''}`;

            ds.data.forEach((d) => {
              // if(ind === 0){
              xAxis.push(d.x);
              // }
              yAxis.push(d.y);
            });
            yAxes.push(yAxis);
            xAxes.push(xAxis);
            if (ds.cubicInterpolationMode) {
              // types[ds.label] = "spline";
            } else {
              types[ds.label] = areaStep();
            }
            colors[ds.label] = ds.borderColor; // not sure which is better border or background
          });
        }

        if (chartData.options.scales.xAxes) {
          chartData.options.scales.xAxes.forEach((ya) => {
            axis.x = {
              show: true,
              label: {
                text: ya.scaleLabel.labelString,
                position: 'outer-middle',
              },
            };
          });
        }
        if (chartData.options.scales.yAxes) {
          chartData.options.scales.yAxes.forEach((ya) => {
            axis[yAxisTracker[ya.id]] = {
              show: true,
              label: {
                text: ya.scaleLabel.labelString,
                position: 'outer-middle',
              },
            };
          });
        }

        const grid = {};

        if (chartData.percentiles && chartData.percentiles.length > 0) {
          // position: "middle"
          // position: "start"
          let reTrack = 0;
          const percentiles = chartData.percentiles.map(({ Percentile, Value }) => {
            const re = {
              value: (Value * 1000).toFixed(2),
              text: `p${Percentile}`,
            };
            switch (reTrack % 3) {
              case 0:
              // re.position
                break;
              case 1:
                re.position = 'middle';
                break;
              case 2:
                re.position = 'start';
                break;
            }

            reTrack++;

            return re;
          });

          grid.x = {
            lines: percentiles,
          };
        }

        const chartConfig = {
          // oninit: function(args){
          //   console.log(JSON.stringify(args));
          // },
          // title: {
          //   text: chartData.options.title.text.join('\n'),
          // },
          bindto: self.chartRef,
          type: line(),
          data: {
            // x: 'x',
            xs: xAxisTracker,
            // xFormat: self.bbTimeFormat,
            columns: [...xAxes, ...yAxes],
            colors,
            axes,
            types,
            // groups,
            // type: 'area',
          },
          axis,

          grid,
          legend: {
            show: true,
          },
          point: {
            r: 0,
            focus: {
              expand: {
                r: 5,
              },
            },
          },
          tooltip: {
            show: true,
          },
        };
        if (!self.props.hideTitle) {
          self.titleRef.innerText = chartData.options.title.text.join('\n');
        }
        self.chart = bb.generate(chartConfig);
      } else {
        self.chart = bb.generate({
          type: line(),
          data: {
            columns: [
            ],
          },
          bindto: self.chartRef,
        });
      }
    }
  }

  processMultiChartData(chartData) { // >= 3 datasets
    const self = this;
    if (chartData && chartData.data && chartData.options) {
      const xAxes = [];
      const categories = [];
      const yAxes = [];
      const colors = [];
      const axes = {};
      const axis = {};
      let yTrack = 1;
      const yAxisTracker = {};
      // const xAxisTracker = {};

      if (chartData.data && chartData.data.datasets) {
        chartData.data.datasets.forEach((ds) => {
          // xAxis.push('x');
          const yAxis = [ds.label];
          // xAxisTracker[ds.label] = `x${ind+1}`;
          if (typeof ds.yAxisID !== 'undefined' && typeof yAxisTracker[ds.yAxisID] === 'undefined') {
            yTrack++;
            yAxisTracker[ds.yAxisID] = `y${yTrack}`;
            axes[ds.label] = `y${yTrack}`;
          }
          // axes[ds.label] = `y${ind>0?ind+1:''}`;

          ds.data.forEach((d) => {
            yAxis.push(d);
          });
          yAxes.push(yAxis);
          colors[ds.label] = ds.borderColor; // not sure which is better border or background
        });
      }
      if (chartData.data && chartData.data.labels) {
        chartData.data.labels.forEach((l) => {
          categories.push(l.join(' '));
        });
      }

      // if(chartData.options.scales.xAxes){
      // chartData.options.scales.xAxes.forEach((ya) => {
      axis.x = {
        show: true,
        label: {
          // text: ya.scaleLabel.labelString,
          // position: 'outer-middle',
        },
        type: 'category',
        categories,
      };

      if (chartData.options.scales.yAxes) {
        chartData.options.scales.yAxes.forEach((ya) => {
          let lab;
          if (typeof yAxisTracker[ya.id] !== 'undefined') {
            lab = yAxisTracker[ya.id];
          } else {
            lab = 'y';
          }
          axis[lab] = {
            show: true,
            min: 0,
            label: {
              text: ya.scaleLabel.labelString,
              position: 'outer-middle',
            },
          };
        });
      }

      // chartConfig.data.colors[chartConfig.data.columns[1][0]] = panel.sparkline.lineColor;
      // chartConfig.color = {
      //     pattern: [
      //       panel.sparkline.fillColor,
      //     ]
      //   };

      if (self.chartRef && self.chartRef !== null) {
        const chartConfig = {
          // oninit: function(args){
          //   console.log(JSON.stringify(args));
          // },
          // title: {
          //   text: chartData.options.title.text.join('\n'),
          // },
          bindto: self.chartRef,
          data: {
            // x: 'x',
            // xs: xAxisTracker,
            // xFormat: self.bbTimeFormat,
            columns: [...xAxes, ...yAxes],
            colors,
            axes,
            // types,
            // groups,
            // type: 'area',
          },
          axis,
          // zoom: {
          //   enabled: {
          //     type: "drag"
          //   },
          //   onzoomend: self.updateDateRange(),
          // },
          // grid: grid,
          legend: {
            show: true,
            // position: 'inset',
            // inset: {
            //   // anchor: 'bottom-left',
            //   // x: 125,
            //   // y: -200,
            //   anchor: 'top-right',
            //   // x: 125,
            //   y: -5,
            //   step: 1,
            // },
            position: 'right',
          },
          // point: {
          //   show: false
          // },
          point: {
            r: 0,
            focus: {
              expand: {
                r: 5,
              },
            },
          },
          tooltip: {
            show: true,
            // linked: !inDialog,
            // format: {
            //   title: function(x) {
            //     // return d3.timeFormat(self.bbTimeFormat)(x);
            //     return moment(x).format(self.timeFormat);
            //   }
            // }
          },
          // area: {
          //   linearGradient: true
          // },
        };

        // if(self.panelType === 'sparkline' && panel.sparkline && panel.sparkline.lineColor && panel.sparkline.fillColor){
        //   // cd.datasets[datasetInd].borderColor = panel.sparkline.lineColor;
        //   // cd.datasets[datasetInd].backgroundColor = panel.sparkline.fillColor;

        //   const dataLength = chartConfig.data.columns && chartConfig.data.columns.length > 1? chartConfig.data.columns[1].length:0; // 0 is for x axis
        //   if(dataLength > 0){
        //     if(typeof chartConfig.data.colors === 'undefined'){
        //       chartConfig.data.colors = {};
        //     }
        //     chartConfig.data.colors[chartConfig.data.columns[1][0]] = panel.sparkline.lineColor;

        //     if(dataLength > 1){
        //       let content = '';
        //       if(panel.format.toLowerCase().startsWith('percent')){
        //         const mulFactor = panel.format.toLowerCase() === 'percentunit'?100:1;
        //         if(!isNaN(chartConfig.data.columns[1][dataLength-1])){
        //           const tk = (chartConfig.data.columns[1][dataLength-1] * mulFactor).toFixed(2);
        //           content = `${tk}%`;
        //         }
        //       } else {
        //         content = `${chartConfig.data.columns[1][dataLength-1]} ${panel.format}`;
        //       }
        //       chartConfig.title = {
        //         text: "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n" + content, // for sparkline, we want to print the value as title
        //       };
        //     }
        //   }

        //   chartConfig.color = {
        //     pattern: [
        //       panel.sparkline.fillColor,
        //     ]
        //   };

        // }
        if (!self.props.hideTitle) {
          self.titleRef.innerText = chartData.options.title.text.join('\n');
        }
        self.chart = bb.generate(chartConfig);
      }
    }
  }

  render() {
    let chartData;
    if (typeof this.props.data !== 'undefined') {
      const results = this.props.data;
      if (results.length == 2) {
        chartData = makeOverlayChart(fortioResultToJsChartData(results[0]), fortioResultToJsChartData(results[1]));
      } else if (results.length > 2) {
        chartData = makeMultiChart(results);
      }
    }
    const self = this;
    if (typeof chartData === 'undefined') {
      const tmpData = (typeof this.props.data !== 'undefined') ? (this.props.data.length == 1 ? this.props.data[0] : {}) : {};
      chartData = this.singleChart(tmpData);
    }
    const { classes } = this.props;
    return (
      <NoSsr>
        <div>
          <div ref={(ch) => this.titleRef = ch} className={classes.title} />
          <div ref={(ch) => {
            this.chartRef = ch;
            if (this.props.data.length > 2) {
              self.processMultiChartData(chartData);
            } else {
              self.processChartData(chartData);
            }
          }}
          />

          {/* <div style={{
            minHeight: '300px',
          }}>
          <Line data={chartData.data} options={chartData.options} />
          </div> */}
        </div>
      </NoSsr>
    );
  }
}

export default withStyles(styles)(MesheryChart);
