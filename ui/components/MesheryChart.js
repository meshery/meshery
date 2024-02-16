import React from 'react';
import { Grid, NoSsr, Typography } from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import {
  fortioResultToJsChartData,
  makeChart,
  makeOverlayChart,
  makeMultiChart,
} from '../lib/chartjs-formatter';
import bb, { areaStep, line } from 'billboard.js';
import {
  TwitterShareButton,
  LinkedinShareButton,
  FacebookShareButton,
  TwitterIcon,
  LinkedinIcon,
  FacebookIcon,
} from 'react-share';
import ReplyIcon from '@material-ui/icons/Reply';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import { ClickAwayListener, Fade, Popper } from '@material-ui/core';
import classNames from 'classnames';

const styles = (theme) => ({
  title: {
    textAlign: 'center',
    fontSize: theme.spacing(1.75),
    marginBottom: theme.spacing(1),
  },
  percentiles: {
    height: '100%',
    justifyContent: 'center',
    display: 'flex',
    position: 'relative',
    alignItems: 'center',
    transform: 'translateY(-30%)',
  },
  chart: {
    width: 'calc( 100% - 150px )',
  },
  chartWrapper: {
    display: 'flex',
    flexWrap: 'no-wrap',
    justifyContent: 'center',
    alignItems: 'center',
  },
  share: {
    transform: 'scaleX(-1)',
  },
  popper: {
    width: 500,
  },
  paper: {
    padding: theme.spacing(1),
  },
  socialIcon: {
    margin: theme.spacing(0.4),
  },
  shareIcon: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
});

function NonRecursiveConstructDisplayCells(data) {
  return Object.keys(data).map((el) => {
    if (typeof data[el].display?.value === 'string' && !data[el].display?.hide) {
      return (
        <>
          <b>{data[el].display?.key}</b>: {data[el].display?.value}
        </>
      );
    }
  });
}

class MesheryChart extends React.Component {
  constructor(props) {
    super(props);
    this.chartRef = null;
    this.chart = null;
    this.percentileRef = null;
    this.state = {
      socialExpand: false,
      anchorEl: null,
      socialMessage: '',
    };
  }

  getSocialMessageForPerformanceTest(rps, percentile) {
    return `I achieved ${rps.trim()} RPS running my service at a P99.9 of ${percentile} ms using @mesheryio with @smp_spec! Find out how fast your service is with`;
  }

  handleSocialExpandClick(e, chartData) {
    this.setState({ anchorEl: e.currentTarget });
    this.setState({
      socialMessage: this.getSocialMessageForPerformanceTest(
        chartData.options.metadata.qps.display.value.split(' ')[1],
        chartData.percentiles[4].Value,
      ),
    });
    e.stopPropagation();
    this.setState((state) => ({ socialExpand: !state.socialExpand }));
  }

  singleChart = (rawdata, data) => {
    if (typeof data === 'undefined' || typeof data.StartTime === 'undefined') {
      return {};
    }
    return makeChart(fortioResultToJsChartData(rawdata, data));
  };

  processChartData(chartData) {
    const self = this;
    if (self.chartRef && self.chartRef !== null) {
      if (chartData && chartData.data && chartData.options) {
        const xAxes = [];
        const yAxes = [];
        const colors = {};
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
              label: { text: ya.scaleLabel.labelString, position: 'outer-middle' },
            };
          });
        }
        if (chartData.options.scales.yAxes) {
          chartData.options.scales.yAxes.forEach((ya) => {
            axis[yAxisTracker[ya.id]] = {
              show: true,
              label: { text: ya.scaleLabel.labelString, position: 'outer-middle' },
            };
          });
        }

        const grid = {};

        if (chartData.percentiles && chartData.percentiles.length > 0) {
          // position: "middle"
          // position: "start"
          let reTrack = 0;
          const percentiles = chartData.percentiles.map(({ Percentile, Value }) => {
            const re = { value: (Value * 1000).toFixed(2), text: `p${Percentile}` };
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

          grid.x = { lines: percentiles };
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
            colors: { ...colors, 'Cumulative %': 'rgb(71,126,150)' },
            axes,
            types,
            // groups,
            // type: 'area',
          },
          axis,

          grid,
          legend: { show: true },
          point: { r: 0, focus: { expand: { r: 5 } } },
          tooltip: { show: true },
        };
        if (!self.props.hideTitle) {
          if (this.props.data.length == 4) {
            self.titleRef.innerText =
              chartData.options.title.text.slice(0, 2).join('\n') +
              '\n' +
              chartData.options.title.text[2].split('\n')[0];
            if (chartData.options.title.text[2])
              self.percentileRef.innerText = chartData.options.title.text[2]
                .split('\n')[1]
                .split('|')
                .join('\n');
          } else {
            self.titleRef.innerText = chartData.options.title.text.join('\n');
          }
        }

        self.chart = bb.generate(chartConfig);
      } else {
        self.chart = bb.generate({
          type: line(),
          data: { columns: [] },
          bindto: self.chartRef,
        });
      }
    }
  }

  processMultiChartData(chartData) {
    // >= 3 datasets
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
          if (
            typeof ds.yAxisID !== 'undefined' &&
            typeof yAxisTracker[ds.yAxisID] === 'undefined'
          ) {
            yTrack++;
            yAxisTracker[ds.yAxisID] = `y${yTrack}`;
            axes[ds.label] = `y${yTrack}`;
          }
          // axes[ds.label] = `y${ind>0?ind+1:''}`;

          ds.data.forEach((d) => yAxis.push(d));
          yAxes.push(yAxis);
          colors[ds.label] = ds.borderColor; // not sure which is better border or background
        });
      }
      if (chartData.data && chartData.data.labels) {
        chartData.data.labels.forEach((l) => {
          categories.push(l.join(' '));
        });
      }

      axis.x = {
        show: true,
        label: {},
        type: 'category',
        categories,
      };

      if (chartData.options.scales.yAxes) {
        chartData.options.scales.yAxes.forEach((ya) => {
          let lab;
          if (typeof yAxisTracker[ya.id] !== 'undefined') lab = yAxisTracker[ya.id];
          else lab = 'y';

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

      if (self.chartRef && self.chartRef !== null) {
        const chartConfig = {
          bindto: self.chartRef,
          data: {
            columns: [...xAxes, ...yAxes],
            colors,
            axes,
          },
          axis,
          legend: { show: true, position: 'right' },
          point: { r: 0, focus: { expand: { r: 5 } } },
          tooltip: { show: true },
        };
        if (!self.props.hideTitle) {
          self.titleRef = chartData.options.title.text;
        }
        self.chart = bb.generate(chartConfig);
      }
    }
  }

  render() {
    let chartData;

    if (typeof this.props.data !== 'undefined') {
      const results = this.props.data;
      if (results.length === 2) {
        chartData = makeOverlayChart(
          fortioResultToJsChartData(this.props.rawdata, results[0]),
          fortioResultToJsChartData(this.props.rawdata, results[1]),
        );
      } else if (results.length > 2) {
        chartData = makeMultiChart(this.props.rawdata, results);
      }
    }
    const self = this;
    if (typeof chartData === 'undefined') {
      const tmpData =
        typeof this.props.data !== 'undefined'
          ? this.props.data.length === 1
            ? this.props.data[0]
            : {}
          : {};
      chartData = this.singleChart(this.props.rawdata, tmpData);
    }

    const { classes } = this.props;

    return (
      <NoSsr>
        <div className={classes.shareIcon}>
          <IconButton
            aria-label="Share"
            className={classes.expand}
            onClick={(e) => this.handleSocialExpandClick(e, chartData)}
          >
            <ReplyIcon className={classNames(classes.share, classes.iconColor)} />
          </IconButton>
        </div>
        <Popper
          open={this.state.socialExpand}
          anchorEl={this.state.anchorEl}
          transition
          style={{ zIndex: '1301' }}
        >
          {({ TransitionProps }) => (
            <ClickAwayListener onClickAway={() => this.setState({ socialExpand: false })}>
              <Fade {...TransitionProps} timeout={350}>
                <Paper className={classes.paper}>
                  <TwitterShareButton
                    className={classes.socialIcon}
                    url={'https://meshery.io'}
                    title={this.state.socialMessage}
                    hashtags={['opensource']}
                  >
                    <TwitterIcon size={32} />
                  </TwitterShareButton>
                  <LinkedinShareButton
                    className={classes.socialIcon}
                    url={'https://meshery.io'}
                    summary={this.state.socialMessage}
                  >
                    <LinkedinIcon size={32} />
                  </LinkedinShareButton>
                  <FacebookShareButton
                    className={classes.socialIcon}
                    url={'https://meshery.io'}
                    quote={this.state.socialMessage}
                    hashtag={'#opensource'}
                  >
                    <FacebookIcon size={32} />
                  </FacebookShareButton>
                </Paper>
              </Fade>
            </ClickAwayListener>
          )}
        </Popper>
        <div>
          <div
            ref={(ch) => (this.titleRef = ch)}
            className={classes.title}
            style={{ display: 'none' }}
          />
          <Grid container spacing={1} style={{ margin: '1rem' }} justifyContent="center">
            {NonRecursiveConstructDisplayCells(chartData?.options?.metadata || {})?.map((el, i) => {
              return (
                <Grid item xs={4} key={`nri-${i}`}>
                  {el}
                </Grid>
              );
            })}
          </Grid>
          <div className={classes.chartWrapper}>
            <div
              className={classes.chart}
              ref={(ch) => {
                this.chartRef = ch;
              }}
            ></div>
            <div
              className={classes.percentiles}
              ref={(ch) => {
                this.percentileRef = ch;
                if (this.props.data.length > 2) {
                  self.processMultiChartData(chartData);
                } else {
                  self.processChartData(chartData);
                }
              }}
            >
              {this.props.data.length === 1 ? (
                <div style={{ margin: '1rem' }}>
                  <Typography style={{ whiteSpace: 'nowrap' }} gutterBottom>
                    Percentile Summary
                  </Typography>
                  <div>
                    {NonRecursiveConstructDisplayCells(
                      chartData?.options?.metadata?.percentiles?.display?.value || {},
                    ).map((el, i) => {
                      return <div key={`percentile-${i}`}>{el}</div>;
                    })}
                  </div>
                </div>
              ) : null}
            </div>
          </div>
        </div>
      </NoSsr>
    );
  }
}

export default withStyles(styles)(MesheryChart);
