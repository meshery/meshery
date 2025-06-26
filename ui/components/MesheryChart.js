import React, { useRef } from 'react';
import { Grid2, Typography, styled } from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import {
  fortioResultToJsChartData,
  makeChart,
  makeOverlayChart,
  makeMultiChart,
} from '../lib/chartjs-formatter';
import bb, { areaStep, line } from 'billboard.js';
const ChartTitle = styled(Typography)(({ theme }) => ({
  textAlign: 'center',
  fontSize: theme.spacing(1.75),
  marginBottom: theme.spacing(1),
}));

const ChartPercentiles = styled('div')({
  height: '100%',
  justifyContent: 'center',
  display: 'flex',
  position: 'relative',
  alignItems: 'center',
  transform: 'translateY(-30%)',
});

const ChartContainer = styled('div')({
  width: 'calc( 100% - 150px )',
});

const ChartWrapper = styled('div')({
  display: 'flex',
  flexWrap: 'no-wrap',
  justifyContent: 'center',
  alignItems: 'center',
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

function MesheryChart(props) {
  const chartRef = useRef(null);
  const chart = useRef(null);
  const percentileRef = useRef(null);
  const titleRef = useRef(null);

  const singleChart = (rawdata, data) => {
    if (typeof data === 'undefined' || typeof data.StartTime === 'undefined') {
      return {};
    }
    return makeChart(fortioResultToJsChartData(rawdata, data));
  };

  const processChartData = (chartData) => {
    if (chartRef.current && chartRef.current !== null) {
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
          bindto: chartRef.current,
          type: line(),
          data: {
            xs: xAxisTracker,
            columns: [...xAxes, ...yAxes],
            colors: { ...colors, 'Cumulative %': 'rgb(71,126,150)' },
            axes,
            types,
          },
          axis,

          grid,
          legend: { show: true },
          point: { r: 0, focus: { expand: { r: 5 } } },
          tooltip: { show: true },
        };
        if (!props.hideTitle) {
          if (props.data.length == 4) {
            titleRef.current.innerText =
              chartData.options.title.text.slice(0, 2).join('\n') +
              '\n' +
              chartData.options.title.text[2].split('\n')[0];
            if (chartData.options.title.text[2])
              percentileRef.current.innerText = chartData.options.title.text[2]
                .split('\n')[1]
                .split('|')
                .join('\n');
          } else {
            titleRef.current.innerText = chartData.options.title.text.join('\n');
          }
        }

        chart.current = bb.generate(chartConfig);
      } else {
        chart.current = bb.generate({
          type: line(),
          data: { columns: [] },
          bindto: chartRef.current,
        });
      }
    }
  };

  const processMultiChartData = (chartData) => {
    // >= 3 datasets

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

      if (chartRef.current && chartRef.current !== null) {
        const chartConfig = {
          bindto: chartRef.current,
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
        if (!props.hideTitle) {
          titleRef.current = chartData.options.title.text;
        }
        chart.current = bb.generate(chartConfig);
      }
    }
  };

  let chartData;

  if (typeof props.data !== 'undefined') {
    const results = props.data;
    if (results.length === 2) {
      chartData = makeOverlayChart(
        fortioResultToJsChartData(props.rawdata, results[0]),
        fortioResultToJsChartData(props.rawdata, results[1]),
      );
    } else if (results.length > 2) {
      chartData = makeMultiChart(props.rawdata, results);
    }
  }

  if (typeof chartData === 'undefined') {
    const tmpData =
      typeof props.data !== 'undefined' ? (props.data.length === 1 ? props.data[0] : {}) : {};
    chartData = singleChart(props.rawdata, tmpData);
  }

  return (
    <NoSsr>
      <div>
        <ChartTitle ref={titleRef} style={{ display: 'none' }} />
        <Grid2 container justifyContent="center" style={{ padding: '0.5rem' }} size="grow">
          {NonRecursiveConstructDisplayCells(chartData?.options?.metadata || {})?.map((el, i) => (
            <Grid2 key={`nri-${i}`} size={{ xs: 4 }}>
              {el}
            </Grid2>
          ))}
        </Grid2>
        <ChartWrapper>
          <ChartContainer ref={chartRef}></ChartContainer>
          <ChartPercentiles
            ref={(ch) => {
              percentileRef.current = ch;
              if (props.data.length > 2) {
                processMultiChartData(chartData);
              } else {
                processChartData(chartData);
              }
            }}
          >
            {props.data.length === 1 ? (
              <div>
                <Typography style={{ whiteSpace: 'nowrap' }} gutterBottom>
                  Percentile Summary
                </Typography>
                <div>
                  {NonRecursiveConstructDisplayCells(
                    chartData?.options?.metadata?.percentiles?.display?.value || {},
                  ).map((el, i) => (
                    <div key={`percentile-${i}`}>{el}</div>
                  ))}
                </div>
              </div>
            ) : null}
          </ChartPercentiles>
        </ChartWrapper>
      </div>
    </NoSsr>
  );
}
export default MesheryChart;
