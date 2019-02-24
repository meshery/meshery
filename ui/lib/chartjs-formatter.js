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
          title.push(res.Labels + ' - ' + res.URL + ' - ' + formatDate(res.StartTime))
      } else { // grpc results
          title.push(res.Labels + ' - ' + res.Destination + ' - ' + formatDate(res.StartTime))
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
    var statusOk = res.RetCodes[200]
    if (!statusOk) { // grpc results
      statusOk = res.RetCodes["SERVING"]
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
          y: myRound(100.0 / res.DurationHistogram.Count, 3)
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
      var y = myRound(it.Percent, 3)
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
      dataH: dataH
    }
  }
  
export function myRound (v, digits = 6) {
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