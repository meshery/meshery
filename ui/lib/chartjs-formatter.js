/* eslint block-scoped-var: 0 */
export const linearXAxe = {
  type: 'linear',
  scaleLabel: {
    display: true,
    labelString: 'Response time in ms',
    ticks: {
      min: 0,
      beginAtZero: true,
    },
  },
};
export const logXAxe = {
  type: 'logarithmic',
  scaleLabel: {
    display: true,
    labelString: 'Response time in ms (log scale)',
  },
  ticks: {
    // min: dataH[0].x, // newer chart.js are ok with 0 on x axis too
    callback: function (tick, index, ticks) {
      return tick.toLocaleString();
    },
  },
};
export const linearYAxe = {
  id: 'H',
  type: 'linear',
  ticks: {
    beginAtZero: true,
  },
  scaleLabel: {
    display: true,
    labelString: 'Count',
  },
};
export const logYAxe = {
  id: 'H',
  type: 'logarithmic',
  display: true,
  ticks: {
    // min: 1, // log mode works even with 0s
    // Needed to not get scientific notation display:
    callback: function (tick, index, ticks) {
      return tick.toString();
    },
  },
  scaleLabel: {
    display: true,
    labelString: 'Count (log scale)',
  },
};

/**
 * getMetadata takes in the test data and returns an object
 * with partially computed data and a "display" field
 */
export function getMetadata(rawdata, res) {
  return {
    title: {
      display: {
        key: 'Title',
        // value : res.Labels.split(' -_- ')?.[0] || "No Title"
        value: (rawdata ? rawdata[0].name : res.Labels.split(' -_- ')?.[0]) || 'No Title',
      },
    },
    url: {
      display: {
        key: 'URL',
        value:
          (rawdata ? rawdata[0].runner_results.URL : res.Labels.split(' -_- ')?.[1]) || 'No URL',
      },
    },
    startTime: {
      display: {
        key: 'Start Time',
        value: formatDate(res.StartTime),
      },
    },
    minimum: {
      display: {
        key: 'Minimum',
        value: `${myRound(1000.0 * res.DurationHistogram.Min, 3)} ms`,
      },
    },
    average: {
      display: {
        key: 'Average',
        value: `${myRound(1000.0 * res.DurationHistogram.Avg, 3)} ms`,
      },
    },
    maximum: {
      display: {
        key: 'Maximum',
        value: `${myRound(1000.0 * res.DurationHistogram.Max, 3)} ms`,
      },
    },
    qps: {
      display: {
        key: 'QPS',
        value: `Achieved ${myRound(res.ActualQPS, 1)} (Requested ${res?.RequestedQPS})`,
      },
    },
    numberOfConnections: {
      display: {
        key: 'Number Of Connections',
        value: res.NumThreads,
      },
    },
    duration: {
      display: {
        key: 'Duration',
        value: `Achieved ${myRound(res.ActualDuration / 1e9, 1)} (Requested ${
          res.RequestedDuration
        })`,
      },
    },
    errors: {
      display: {
        key: 'Errors',
        value: (() => {
          const status = res.RetCodes?.[200] || res.RetCodes?.SERVING || 0;
          const total = res.DurationHistogram.Count;

          if (status !== total) {
            if (status) return myRound((100.0 * (total - status)) / total, 2) + '% errors';

            return '100% errors!';
          }

          return 'No Errors';
        })(),
      },
    },
    percentiles: {
      display: {
        key: 'Percentiles',
        value: res.DurationHistogram?.Percentiles?.map((p) => {
          return {
            display: {
              key: `p${p.Percentile}`,
              value: `${myRound(1000 * p.Value, 2)} ms`,
            },
          };
        }),
      },
    },
    kubernetes: {
      display: {
        key: 'Kubernetes',
        value: [
          {
            display: {
              key: 'Server Version',
              value: res.kubernetes?.server_version,
            },
          },
          {
            display: {
              key: 'Nodes',
              value: res.kubernetes?.nodes?.map((node, i) => {
                return {
                  display: {
                    key: `Node ${i + 1}`,
                    value: [
                      {
                        display: {
                          key: 'Hostname',
                          value: node?.hostname,
                        },
                      },
                      {
                        display: {
                          key: 'CPU',
                          value: node?.allocatable_cpu,
                        },
                      },
                      {
                        display: {
                          key: 'Memory',
                          value: node?.allocatable_memory,
                        },
                      },
                      {
                        display: {
                          key: 'Arch',
                          value: node?.architecture,
                        },
                      },
                      {
                        display: {
                          key: 'OS',
                          value: node?.operating_system,
                        },
                      },
                      {
                        display: {
                          key: 'Kubelet Version',
                          value: node?.kubelet_version,
                        },
                      },
                      {
                        display: {
                          key: 'Container runtime',
                          value: node?.container_runtime_version,
                        },
                      },
                    ],
                  },
                };
              }),
            },
          },
        ],
      },
    },
  };
}

export function makeTitle(rawdata, res) {
  var title = [];
  if (res.Labels !== '') {
    if (res.URL) {
      // http results
      // title.push(res.Labels + ' - ' + res.URL + ' - ' + formatDate(res.StartTime))
      // title.push(res.URL + ' - ' + formatDate(res.StartTime))
      console.log(res.Labels);
      var labels = res.Labels.split(' -_- ');
      // title.push(`Labels: ${labels.map(item => item + '\n')}`)
      title.push(`Title: ${rawdata ? rawdata[0].name : labels[0]}`);
      title.push(`URL: ${rawdata ? rawdata[0].runner_results.URL : labels[1]}`);
      title.push(`Start Time: ${formatDate(res.StartTime)}`);
    } else {
      // grpc results
      title.push(`Destination: ${res.Destination}`);
      title.push(`Start Time: ${formatDate(res.StartTime)}`);
    }
  }
  title.push(`Minimum: ${myRound(1000.0 * res.DurationHistogram.Min, 3)} ms`);
  title.push(`Average: ${myRound(1000.0 * res.DurationHistogram.Avg, 3)} ms`);
  title.push(`Maximum: ${myRound(1000.0 * res.DurationHistogram.Max, 3)} ms`);
  var percStr = `Minimum: ${myRound(1000.0 * res.DurationHistogram.Min, 3)} ms \nAverage: ${myRound(
    1000.0 * res.DurationHistogram.Avg,
    3,
  )} ms \nMaximum: ${myRound(1000.0 * res.DurationHistogram.Max, 3)} ms\n`;
  var percStr_2 = 'Percentiles: ';
  if (res.DurationHistogram.Percentiles) {
    for (var i = 0; i < res.DurationHistogram.Percentiles.length; i++) {
      var p = res.DurationHistogram.Percentiles[i];
      percStr_2 += `p${p.Percentile}: ${myRound(1000 * p.Value, 2)} ms; `;
      percStr += `p${p.Percentile}: ${myRound(1000 * p.Value, 2)} ms; `;
    }
    percStr = percStr.slice(0, -2);
  }
  var statusOk =
    typeof res.RetCodes !== 'undefined' && res.RetCodes !== null ? res.RetCodes[200] : 0;
  if (!statusOk) {
    // grpc results
    statusOk =
      typeof res.RetCodes !== 'undefined' && res.RetCodes !== null ? res.RetCodes.SERVING : 0;
  }
  var total = res.DurationHistogram.Count;
  var errStr = 'No Error';
  if (statusOk !== total) {
    if (statusOk) {
      errStr = myRound((100.0 * (total - statusOk)) / total, 2) + '% errors';
    } else {
      errStr = '100% errors!';
    }
  }
  title.push(`Target QPS: ${res.RequestedQPS} ( Actual QPS: ${myRound(res.ActualQPS, 1)} )`);
  title.push(`No of Connections: ${res.NumThreads}`);
  title.push(
    `Requested Duration: ${res.RequestedDuration} ( Actual Duration: ${myRound(
      res.ActualDuration / 1e9,
      1,
    )} )`,
  );
  title.push(`Errors: ${errStr}`);
  title.push(percStr_2);
  if (res.kubernetes) {
    title.push(`Kubernetes server version: ${res.kubernetes.server_version}`);
    title.push('Nodes:');
    res.kubernetes?.nodes?.forEach((node, ind) => {
      title.push(`Node ${ind + 1} - \nHostname: ${node.hostname} \nCPU: ${
        node.allocatable_cpu
      } \nMemory: ${node.allocatable_memory} \nArch: ${node.architecture} \nOS: ${node.os_image}
                    \nKubelet version: ${node.kubelet_version} \nContainer runtime: ${
                      node.container_runtime_version
                    }`);
    });
  }

  return title;
}

export function fortioResultToJsChartData(rawdata, res) {
  var dataP = [
    {
      x: 0.0,
      y: 0.0,
    },
  ];
  var len = res.DurationHistogram.Data.length;
  var prevX = 0.0;
  var prevY = 0.0;
  for (var i = 0; i < len; i++) {
    var it = res.DurationHistogram.Data[i];
    var x = myRound(1000.0 * it.Start);
    if (i === 0) {
      // Extra point, 1/N at min itself
      dataP.push({
        x: x,
        // y: myRound(100.0 / res.DurationHistogram.Count, 3)
        y: myRound(100.0 / res.DurationHistogram.Count, 2),
      });
    } else {
      if (prevX !== x) {
        dataP.push({
          x: x,
          y: prevY,
        });
      }
    }
    x = myRound(1000.0 * it.End);
    // var y = myRound(it.Percent, 3)
    var y = myRound(it.Percent, 2);
    dataP.push({
      x: x,
      y: y,
    });
    prevX = x;
    prevY = y;
  }
  var dataH = [];
  var prev = 1000.0 * res.DurationHistogram.Data[0].Start;
  for (i = 0; i < len; i++) {
    it = res.DurationHistogram.Data[i];
    var startX = 1000.0 * it.Start;
    var endX = 1000.0 * it.End;
    if (startX !== prev) {
      dataH.push(
        {
          x: myRound(prev),
          y: 0,
        },
        {
          x: myRound(startX),
          y: 0,
        },
      );
    }
    dataH.push(
      {
        x: myRound(startX),
        y: it.Count,
      },
      {
        x: myRound(endX),
        y: it.Count,
      },
    );
    prev = endX;
  }
  return {
    title: makeTitle(rawdata, res),
    metadata: getMetadata(rawdata, res),
    dataP: dataP,
    dataH: dataH,
    percentiles: res.DurationHistogram.Percentiles,
  };
}

// export function myRound (v, digits = 6) {
export function myRound(v, digits = 2) {
  var p = Math.pow(10, digits);
  return Math.round(v * p) / p;
}

export function pad(n) {
  return n < 10 ? '0' + n : n;
}

export function formatDate(dStr) {
  var d = new Date(dStr);
  return (
    d.getFullYear() +
    '-' +
    pad(d.getMonth() + 1) +
    '-' +
    pad(d.getDate()) +
    ' ' +
    pad(d.getHours()) +
    ':' +
    pad(d.getMinutes()) +
    ':' +
    pad(d.getSeconds())
  );
}

export function makeChart(data) {
  return {
    percentiles: data.percentiles,
    data: {
      datasets: [
        {
          label: 'Cumulative %',
          data: data.dataP,
          fill: false,
          yAxisID: 'P',
          stepped: true,
          backgroundColor: 'rgba(134, 87, 167, 1)',
          borderColor: 'rgba(134, 87, 167, 1)',
          cubicInterpolationMode: 'monotone',
        },
        {
          label: 'Histogram: Count',
          data: data.dataH,
          yAxisID: 'H',
          pointStyle: 'rect',
          radius: 1,
          borderColor: 'rgba(87, 167, 134, .9)',
          backgroundColor: 'rgba(87, 167, 134, .75)',
          lineTension: 0,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      metadata: data?.metadata,
      title: {
        display: true,
        fontStyle: 'normal',
        text: data.title,
      },
      scales: {
        xAxes: [linearXAxe],
        yAxes: [
          {
            id: 'P',
            position: 'right',
            ticks: {
              beginAtZero: true,
              max: 100,
            },
            scaleLabel: {
              display: true,
              labelString: '%',
            },
          },
          linearYAxe,
        ],
      },
    },
  };

  //     // TODO may need updateChart() if we persist settings even the first time
  // } else {
  //   chart.data.datasets[0].data = data.dataP
  //   chart.data.datasets[1].data = data.dataH
  //   chart.options.title.text = data.title
  //   updateChart(chart)
  // }
}

export function makeOverlayChartTitle(titleA, titleB) {
  // Each string in the array is a separate line
  return [
    'A: ' + titleA[0],
    titleA[1], // Skip 3rd line.
    '',
    'B: ' + titleB[0],
    titleB[1], // Skip 3rd line.
  ];
}

export function makeOverlayChart(dataA, dataB) {
  // var chartEl = document.getElementById('chart1')
  // chartEl.style.visibility = 'visible'
  // if (Object.keys(overlayChart).length !== 0) {
  //   return
  // }
  // deleteSingleChart()
  // deleteMultiChart()
  // var ctx = chartEl.getContext('2d')
  var title = makeOverlayChartTitle(dataA.title, dataB.title);
  return {
    data: {
      // "Cumulative %" datasets are listed first so they are drawn on top of the histograms.
      datasets: [
        {
          label: 'A: Cumulative %',
          data: dataA.dataP,
          fill: false,
          yAxisID: 'P',
          stepped: true,
          backgroundColor: 'rgba(134, 87, 167, 1)',
          borderColor: 'rgba(134, 87, 167, 1)',
          cubicInterpolationMode: 'monotone',
        },
        {
          label: 'B: Cumulative %',
          data: dataB.dataP,
          fill: false,
          yAxisID: 'P',
          stepped: true,
          backgroundColor: 'rgba(204, 102, 0)',
          borderColor: 'rgba(204, 102, 0)',
          cubicInterpolationMode: 'monotone',
        },
        {
          label: 'A: Histogram: Count',
          data: dataA.dataH,
          yAxisID: 'H',
          pointStyle: 'rect',
          radius: 1,
          borderColor: 'rgba(87, 167, 134, .9)',
          backgroundColor: 'rgba(87, 167, 134, .75)',
          lineTension: 0,
        },
        {
          label: 'B: Histogram: Count',
          data: dataB.dataH,
          yAxisID: 'H',
          pointStyle: 'rect',
          radius: 1,
          borderColor: 'rgba(36, 64, 238, .9)',
          backgroundColor: 'rgba(36, 64, 238, .75)',
          lineTension: 0,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      metadata: [dataA?.metadata, dataB?.metadata],
      title: {
        display: true,
        fontStyle: 'normal',
        text: title,
      },
      scales: {
        xAxes: [linearXAxe],
        yAxes: [
          {
            id: 'P',
            position: 'right',
            ticks: {
              beginAtZero: true,
              max: 100,
            },
            scaleLabel: {
              display: true,
              labelString: '%',
            },
          },
          linearYAxe,
        ],
      },
    },
  };
  // updateChart(overlayChart)
}

export function makeMultiChart(rawdata, results) {
  // document.getElementById('running').style.display = 'none'
  // document.getElementById('update').style.visibility = 'hidden'
  // var chartEl = document.getElementById('chart1')
  // chartEl.style.visibility = 'visible'
  // if (Object.keys(mchart).length !== 0) {
  //   return
  // }
  // deleteSingleChart()
  // deleteOverlayChart()
  // var ctx = chartEl.getContext('2d')
  let data = {
    // type: 'line',
    data: {
      labels: [],
      datasets: [
        {
          label: 'Min',
          fill: false,
          stepped: true,
          borderColor: 'hsla(111, 100%, 40%, .8)',
          backgroundColor: 'hsla(111, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'Median',
          fill: false,
          stepped: true,
          borderDash: [5, 5],
          borderColor: 'hsla(220, 100%, 40%, .8)',
          backgroundColor: 'hsla(220, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'Avg',
          fill: false,
          stepped: true,
          backgroundColor: 'hsla(266, 100%, 40%, .8)',
          borderColor: 'hsla(266, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'p75',
          fill: false,
          stepped: true,
          backgroundColor: 'hsla(60, 100%, 40%, .8)',
          borderColor: 'hsla(60, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'p90',
          fill: false,
          stepped: true,
          backgroundColor: 'hsla(45, 100%, 40%, .8)',
          borderColor: 'hsla(45, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'p99',
          fill: false,
          stepped: true,
          backgroundColor: 'hsla(30, 100%, 40%, .8)',
          borderColor: 'hsla(30, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'p99.9',
          fill: false,
          stepped: true,
          backgroundColor: 'hsla(15, 100%, 40%, .8)',
          borderColor: 'hsla(15, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'Max',
          fill: false,
          stepped: true,
          borderColor: 'hsla(0, 100%, 40%, .8)',
          backgroundColor: 'hsla(0, 100%, 40%, .8)',
          data: [],
        },
        {
          label: 'QPS',
          yAxisID: 'qps',
          fill: false,
          stepped: true,
          borderColor: 'rgba(0, 0, 0, .8)',
          backgroundColor: 'rgba(0, 0, 0, .8)',
          data: [],
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      metadata: results?.metadata,
      title: {
        display: true,
        fontStyle: 'normal',
        text: ['Latency in milliseconds'],
      },
      elements: {
        line: {
          tension: 0, // disables bezier curves
        },
      },
      scales: {
        yAxes: [
          {
            id: 'ms',
            ticks: {
              beginAtZero: true,
            },
            scaleLabel: {
              display: true,
              labelString: 'ms',
            },
          },
          {
            id: 'qps',
            position: 'right',
            ticks: {
              beginAtZero: true,
            },
            scaleLabel: {
              display: true,
              labelString: 'QPS',
            },
          },
        ],
      },
    },
  };

  const multiLabel = (res) => {
    var l = formatDate(res.StartTime);
    if (res.Labels !== '') {
      if (res.Labels.indexOf(' -_- ') > -1) {
        const ls = res.Labels.split('-_-'); // trying to match this with server side in fortio.go
        if (ls.length > 0) {
          l += ' - ' + ls[0];
        } else {
          l += ' - ' + res.Labels;
        }
      } else {
        l += ' - ' + res.Labels;
      }
    }
    return label_trunc(l);
  };

  const label_trunc = function (str) {
    if (str.length > length) {
      return str.match(/.{1,20}/g);
    } else {
      return str;
    }
  };

  const findData = (slot, idx, res, p) => {
    // Not very efficient but there are only a handful of percentiles
    var pA = res.DurationHistogram.Percentiles;
    if (!pA) {
      //    console.log('No percentiles in res', res)
      return;
    }
    var pN = Number(p);
    for (var i = 0; i < pA.length; i++) {
      if (pA[i].Percentile === pN) {
        data.data.datasets[slot].data[idx] = 1000.0 * pA[i].Value;
        return;
      }
    }
    console.log('Not Found', p, pN, pA);
    // not found, not set
  };

  const fortioAddToMultiResult = (i, res) => {
    data.data.labels[i] = multiLabel(res);
    data.data.datasets[0].data[i] = 1000.0 * res.DurationHistogram.Min;
    findData(1, i, res, '50');
    data.data.datasets[2].data[i] = 1000.0 * res.DurationHistogram.Avg;
    findData(3, i, res, '75');
    findData(4, i, res, '90');
    findData(5, i, res, '99');
    findData(6, i, res, '99.9');
    data.data.datasets[7].data[i] = 1000.0 * res.DurationHistogram.Max;
    data.data.datasets[8].data[i] = res.ActualQPS;
  };

  const endMultiChart = (len) => {
    data.data.labels = data.data.labels.slice(0, len);
    for (var i = 0; i < data.data.datasets.length; i++) {
      data.data.datasets[i].data = data.data.datasets[i].data.slice(0, len);
    }
    // mchart.update()
  };

  for (var i = 0; i < results.length; i++) {
    fortioAddToMultiResult(i, results[i]);
  }
  endMultiChart(results.length);

  return data;
}
