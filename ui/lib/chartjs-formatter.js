export const linearXAxe = {
    type: 'linear',
    scaleLabel: {
      display: true,
      labelString: 'Response time in ms',
      ticks: {
        min: 0,
        beginAtZero: true
      }
    }
  }
export const logXAxe = {
    type: 'logarithmic',
    scaleLabel: {
      display: true,
      labelString: 'Response time in ms (log scale)'
    },
    ticks: {
          // min: dataH[0].x, // newer chart.js are ok with 0 on x axis too
      callback: function (tick, index, ticks) {
        return tick.toLocaleString()
      }
    }
  }
export const linearYAxe = {
    id: 'H',
    type: 'linear',
    ticks: {
      beginAtZero: true
    },
    scaleLabel: {
      display: true,
      labelString: 'Count'
    }
  }
export const logYAxe = {
    id: 'H',
    type: 'logarithmic',
    display: true,
    ticks: {
          // min: 1, // log mode works even with 0s
          // Needed to not get scientific notation display:
      callback: function (tick, index, ticks) {
        return tick.toString()
      }
    },
    scaleLabel: {
      display: true,
      labelString: 'Count (log scale)'
    }
  }
  
export function makeTitle (res) {
    var title = []
    if (res.Labels !== '') {
      if (res.URL) { // http results
          // title.push(res.Labels + ' - ' + res.URL + ' - ' + formatDate(res.StartTime))
          // title.push(res.URL + ' - ' + formatDate(res.StartTime))
          title.push(res.Labels + ' - ' + formatDate(res.StartTime))
      } else { // grpc results
          title.push(res.Destination + ' - ' + formatDate(res.StartTime))
      }
    }
    var percStr = 'min ' + myRound(1000.0 * res.DurationHistogram.Min, 3) + ' ms, average ' + myRound(1000.0 * res.DurationHistogram.Avg, 3) + ' ms'
    if (res.DurationHistogram.Percentiles) {
      for (var i = 0; i < res.DurationHistogram.Percentiles.length; i++) {
        var p = res.DurationHistogram.Percentiles[i]
        percStr += ', p' + p.Percentile + ' ' + myRound(1000 * p.Value, 2) + ' ms'
      }
    }
    percStr += ', max ' + myRound(1000.0 * res.DurationHistogram.Max, 3) + ' ms'
    var statusOk = typeof res.RetCodes !== 'undefined' && res.RetCodes !== null?res.RetCodes[200]:0;
    if (!statusOk) { // grpc results
      statusOk = typeof res.RetCodes !== 'undefined' && res.RetCodes !== null?res.RetCodes["SERVING"]:0;
    }
    var total = res.DurationHistogram.Count
    var errStr = 'no error'
    if (statusOk !== total) {
      if (statusOk) {
        errStr = myRound(100.0 * (total - statusOk) / total, 2) + '% errors'
      } else {
        errStr = '100% errors!'
      }
    }
    title.push('Response time histogram at ' + res.RequestedQPS + ' target qps (' +
          myRound(res.ActualQPS, 1) + ' actual) ' + res.NumThreads + ' connections for ' +
          res.RequestedDuration + ' (actual time ' + myRound(res.ActualDuration / 1e9, 1) + 's), ' +
          errStr)
    title.push(percStr)

    if(res.kubernetes){
      title.push(`\nKubernetes server version: ${res.kubernetes.server_version}`);
      title.push("\nNodes:");
      res.kubernetes.nodes.forEach((node, ind) => {
        title.push(`\nNode ${ind+1} - Hostname: ${node.hostname}, CPU: ${node.allocatable_cpu}, Memory: ${node.allocatable_memory}, ` +
                    `Arch: ${node.architecture} OS: ${node.os_image}, \n` +
                    `Kubelet version: ${node.kubelet_version}, Container runtime: ${node.container_runtime_version}`);
      });
    }

    return title
  }
  
export function fortioResultToJsChartData (res) {
    var dataP = [{
      x: 0.0,
      y: 0.0
    }]
    var len = res.DurationHistogram.Data.length
    var prevX = 0.0
    var prevY = 0.0
    for (var i = 0; i < len; i++) {
      var it = res.DurationHistogram.Data[i]
      var x = myRound(1000.0 * it.Start)
      if (i === 0) {
        // Extra point, 1/N at min itself
        dataP.push({
          x: x,
          // y: myRound(100.0 / res.DurationHistogram.Count, 3)
          y: myRound(100.0 / res.DurationHistogram.Count, 2)
        })
      } else {
        if (prevX !== x) {
          dataP.push({
            x: x,
            y: prevY
          })
        }
      }
      x = myRound(1000.0 * it.End)
      // var y = myRound(it.Percent, 3)
      var y = myRound(it.Percent, 2)
      dataP.push({
        x: x,
        y: y
      })
      prevX = x
      prevY = y
    }
    var dataH = []
    var prev = 1000.0 * res.DurationHistogram.Data[0].Start
    for (i = 0; i < len; i++) {
      it = res.DurationHistogram.Data[i]
      var startX = 1000.0 * it.Start
      var endX = 1000.0 * it.End
      if (startX !== prev) {
        dataH.push({
          x: myRound(prev),
          y: 0
        }, {
          x: myRound(startX),
          y: 0
        })
      }
      dataH.push({
        x: myRound(startX),
        y: it.Count
      }, {
        x: myRound(endX),
        y: it.Count
      })
      prev = endX
    }
    return {
      title: makeTitle(res),
      dataP: dataP,
      dataH: dataH,
      percentiles: res.DurationHistogram.Percentiles,
    }
  }
  
// export function myRound (v, digits = 6) {
export function myRound (v, digits = 2) {
    var p = Math.pow(10, digits)
    return Math.round(v * p) / p
  }
  
  export function pad (n) {
    return (n < 10) ? ('0' + n) : n
  }
  
  export function formatDate (dStr) {
    var d = new Date(dStr)
    return d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate()) + ' ' +
          pad(d.getHours()) + ':' + pad(d.getMinutes()) + ':' + pad(d.getSeconds())
  }
  
export function makeChart (data) {
    return {
        percentiles: data.percentiles,
        data: {
          datasets: [{
            label: 'Cumulative %',
            data: data.dataP,
            fill: false,
            yAxisID: 'P',
            stepped: true,
            backgroundColor: 'rgba(134, 87, 167, 1)',
            borderColor: 'rgba(134, 87, 167, 1)',
            cubicInterpolationMode: 'monotone'
          },
          {
            label: 'Histogram: Count',
            data: data.dataH,
            yAxisID: 'H',
            pointStyle: 'rect',
            radius: 1,
            borderColor: 'rgba(87, 167, 134, .9)',
            backgroundColor: 'rgba(87, 167, 134, .75)',
            lineTension: 0
          }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          title: {
            display: true,
            fontStyle: 'normal',
            text: data.title
          },
          scales: {
            xAxes: [
              linearXAxe
            ],
            yAxes: [{
              id: 'P',
              position: 'right',
              ticks: {
                beginAtZero: true,
                max: 100
              },
              scaleLabel: {
                display: true,
                labelString: '%'
              }
            },
              linearYAxe
            ]
          }
        }
    }
  
    //     // TODO may need updateChart() if we persist settings even the first time
    // } else {
    //   chart.data.datasets[0].data = data.dataP
    //   chart.data.datasets[1].data = data.dataH
    //   chart.options.title.text = data.title
    //   updateChart(chart)
    // }
  }

  export function makeOverlayChartTitle (titleA, titleB) {
    // Each string in the array is a separate line
    return [
      'A: ' + titleA[0], titleA[1], // Skip 3rd line.
      '',
      'B: ' + titleB[0], titleB[1], // Skip 3rd line.
    ]
  }
  
  export function makeOverlayChart (dataA, dataB) {
    // var chartEl = document.getElementById('chart1')
    // chartEl.style.visibility = 'visible'
    // if (Object.keys(overlayChart).length !== 0) {
    //   return
    // }
    // deleteSingleChart()
    // deleteMultiChart()
    // var ctx = chartEl.getContext('2d')
    var title = makeOverlayChartTitle(dataA.title, dataB.title)
    return {
      data: {
        // "Cumulative %" datasets are listed first so they are drawn on top of the histograms.
        datasets: [{
          label: 'A: Cumulative %',
          data: dataA.dataP,
          fill: false,
          yAxisID: 'P',
          stepped: true,
          backgroundColor: 'rgba(134, 87, 167, 1)',
          borderColor: 'rgba(134, 87, 167, 1)',
          cubicInterpolationMode: 'monotone'
        }, {
          label: 'B: Cumulative %',
          data: dataB.dataP,
          fill: false,
          yAxisID: 'P',
          stepped: true,
          backgroundColor: 'rgba(204, 102, 0)',
          borderColor: 'rgba(204, 102, 0)',
          cubicInterpolationMode: 'monotone'
        }, {
          label: 'A: Histogram: Count',
          data: dataA.dataH,
          yAxisID: 'H',
          pointStyle: 'rect',
          radius: 1,
          borderColor: 'rgba(87, 167, 134, .9)',
          backgroundColor: 'rgba(87, 167, 134, .75)',
          lineTension: 0
        }, {
          label: 'B: Histogram: Count',
          data: dataB.dataH,
          yAxisID: 'H',
          pointStyle: 'rect',
          radius: 1,
          borderColor: 'rgba(36, 64, 238, .9)',
          backgroundColor: 'rgba(36, 64, 238, .75)',
          lineTension: 0
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        title: {
          display: true,
          fontStyle: 'normal',
          text: title
        },
        scales: {
          xAxes: [
            linearXAxe
          ],
          yAxes: [{
            id: 'P',
            position: 'right',
            ticks: {
              beginAtZero: true,
              max: 100
            },
            scaleLabel: {
              display: true,
              labelString: '%'
            }
          },
            linearYAxe
          ]
        }
      }
    }
    // updateChart(overlayChart)
  }
  
  export function makeMultiChart (results) {
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
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        title: {
          display: true,
          fontStyle: 'normal',
          text: ['Latency in milliseconds']
        },
        elements: {
          line: {
            tension: 0 // disables bezier curves
          }
        },
        scales: {
          yAxes: [{
            id: 'ms',
            ticks: {
              beginAtZero: true
            },
            scaleLabel: {
              display: true,
              labelString: 'ms'
            }
          }, {
            id: 'qps',
            position: 'right',
            ticks: {
              beginAtZero: true
            },
            scaleLabel: {
              display: true,
              labelString: 'QPS'
            }
          }]
        }
      }
    }

const multiLabel = (res) => {
  var l = formatDate(res.StartTime)
  if (res.Labels !== '') {
    if (res.Labels.indexOf(' -_- ') > -1) {
      const ls = res.Labels.split(' -_- '); // trying to match this with server side in fortio.go
      if (ls.length > 0){
        l += ' - ' + ls[0];
      } else {
        l += ' - ' + res.Labels;
      }
    } else {
      l += ' - ' + res.Labels
    }
  }
  return label_trunc(l);
}

const label_trunc = function(str) {
  if (str.length > length) {
    return str.match(/.{1,20}/g);;
  } else {
    return str;
  }
};
  
const findData = (slot, idx, res, p) => {
    // Not very efficient but there are only a handful of percentiles
    var pA = res.DurationHistogram.Percentiles
    if (!pA) {
  //    console.log('No percentiles in res', res)
      return
    }
    var pN = Number(p)
    for (var i = 0; i < pA.length; i++) {
      if (pA[i].Percentile === pN) {
        data.data.datasets[slot].data[idx] = 1000.0 * pA[i].Value
        return
      }
    }
    console.log('Not Found', p, pN, pA)
    // not found, not set
  }

  const fortioAddToMultiResult = (i, res) => {
    data.data.labels[i] = multiLabel(res)
    data.data.datasets[0].data[i] = 1000.0 * res.DurationHistogram.Min
    findData(1, i, res, '50')
    data.data.datasets[2].data[i] = 1000.0 * res.DurationHistogram.Avg
    findData(3, i, res, '75')
    findData(4, i, res, '90')
    findData(5, i, res, '99')
    findData(6, i, res, '99.9')
    data.data.datasets[7].data[i] = 1000.0 * res.DurationHistogram.Max
    data.data.datasets[8].data[i] = res.ActualQPS
  }

  const endMultiChart = (len) => {
    data.data.labels = data.data.labels.slice(0, len)
    for (var i = 0; i < data.data.datasets.length; i++) {
      data.data.datasets[i].data = data.data.datasets[i].data.slice(0, len)
    }
    // mchart.update()
  }

  for(var i=0;i<results.length; i++){
    fortioAddToMultiResult(i, results[i]);
  }
  endMultiChart(results.length);

  return data;
}