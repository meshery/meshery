import React from 'react';
import { NoSsr } from '@material-ui/core';
import { Line } from 'react-chartjs-2';
import { fortioResultToJsChartData, makeChart, makeOverlayChart, makeMultiChart } from '../lib/chartjs-formatter';

class MesheryChart extends React.Component {
  singleChart = (data) => {
      if (typeof data === 'undefined' || typeof data.URL === 'undefined'){
          data = {
            DurationHistogram: {
                Count: 0,
                Data: [
                    {
                        Start: new Date(),
                    }
                ],
                Max: 0,
                Min: 0,
                Avg: 0,
                Percentiles: []
            },
            StartTime: new Date(),
            URL: '',
            Labels: '',
            RetCodes: {}
        };
      }
      return makeChart(fortioResultToJsChartData(data));
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
    if (typeof chartData === 'undefined') {
      chartData = singleChart(this.props.data);
    }
    return (
      <NoSsr>
        <Line data={chartData.data} options={chartData.options} />
      </NoSsr>
    );
  }
}

export default MesheryChart;
