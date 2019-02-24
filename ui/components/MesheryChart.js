import React from 'react';
import { NoSsr } from '@material-ui/core';
import { Line } from 'react-chartjs-2';
import { fortioResultToJsChartData, makeChart } from '../lib/chartjs-formatter';

class MesheryChart extends React.Component {

  render() {
      let { data } = this.props;
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
      const chartData = makeChart(fortioResultToJsChartData(data));
    return (
      <NoSsr>
        <Line data={chartData.data} options={chartData.options} />
      </NoSsr>
    );
  }
}

export default MesheryChart;
