import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';
import OpenInNewIcon from '@mui/icons-material/OpenInNewOutlined';
import WarningIcon from '@mui/icons-material/Warning';
import CachedIcon from '@mui/icons-material/Cached';
import dataFetch from '../../../lib/data-fetch';
import { updateProgress } from '../../../lib/store';
import GrafanaCustomGaugeChart from './GrafanaCustomGaugeChart';
import bb, { area, line } from 'billboard.js';
import {
  IconButton,
  Card,
  CardContent,
  CardHeader,
  Tooltip,
  Box,
  styled,
  LinearProgress,
} from '@layer5/sistent';
import grafanaDateRangeToDate from './helper';

const StyledCard = styled(Card)(() => ({
  height: '100%',
  width: '100%',
}));

const StyledCardContent = styled(CardContent)(() => ({
  height: '100%',
  width: '100%',
}));

// const StyledSparklineContent = styled('div')({
//   display: 'grid',
//   gap: '0.5rem',
// });

const ChartContainer = styled('div')({
  width: '100%',
});

const HeaderIcon = styled(IconButton)(({ theme }) => ({
  '& svg': {
    fontSize: theme.spacing(2),
  },
}));

const ErrorText = styled('span')(({ theme }) => ({
  color: theme.palette.error.default,
}));

function GrafanaCustomChart(props) {
  const {
    panel,
    board,
    inDialog,
    handleChartDialogOpen,
    panelData,
    prometheusURL,
    grafanaURL,
    grafanaAPIKey,
    from,
    to,
    templateVars,
    testUUID,
    connectionID,
    refresh,
    liveTail,
    updateProgress,
    updateDateRange,
    sparkline,
  } = props;

  const chartRef = useRef(null);
  const [chart, setChart] = useState(null);
  const timeFormat = 'MM/DD/YYYY HH:mm:ss';
  const bbTimeFormat = '%Y-%m-%d %h:%M:%S %p';

  // Determine panelType based on panel properties
  let panelType = '';
  switch (panel.type) {
    case 'graph':
      panelType = panel.type;
      break;
    case 'singlestat':
      panelType =
        panel.type === 'singlestat' && panel.sparkline && panel.sparkline.show === true
          ? 'sparkline'
          : 'gauge';
      // panelType = panel.type ==='singlestat' && panel.sparkline ? 'sparkline':'gauge';
      break;
  }

  // State declarations
  const [datasetIndex, setDatasetIndex] = useState({});
  const [xAxis, setXAxis] = useState([]);
  const [chartData, setChartData] = useState([]);
  const [error, setError] = useState('');
  const [errorCount, setErrorCount] = useState(0);
  const [isSparkline] = useState(!!sparkline);
  const [interval, setIntervalState] = useState(undefined);

  // Helper function to get or create index
  const getOrCreateIndex = (datasetInd) => {
    if (typeof datasetIndex[datasetInd] !== 'undefined') {
      return datasetIndex[datasetInd];
    }
    let max = 0;
    Object.keys(datasetIndex).forEach((i) => {
      if (datasetIndex[i] > max) {
        max = datasetIndex[i];
      }
    });
    const newDatasetIndex = { ...datasetIndex };
    newDatasetIndex[datasetInd] = max + 1;
    setDatasetIndex(newDatasetIndex);
    return max + 1;
  };

  // Compute refresh interval
  const computeRefreshInterval = (refreshStr) => {
    refreshStr = refreshStr.toLowerCase();
    const l = refreshStr.length;
    const dur = refreshStr.substring(l - 1, l);
    refreshStr = refreshStr.substring(0, l - 1);
    let val = parseInt(refreshStr);
    if (dur === 'd') {
      val *= 24;
    }
    if (dur === 'h') {
      val *= 60;
    }
    if (dur === 'm') {
      val *= 60;
    }
    if (dur === 's') {
      return val;
    }
    return 30; //fallback
  };

  // Compute step for queries
  const computeStep = (start, end) => {
    let step = 10;
    const diff = end - start;
    const min = 60;
    const hrs = 60 * min;
    const days = 24 * hrs;
    const month = 30 * days; // approx.
    const year = 12 * month; // approx.

    if (diff <= 10 * min) {
      // 10 mins
      step = 5;
    } else if (diff <= 30 * min) {
      // 30 mins
      step = 10;
    } else if (diff > 30 * min && diff <= 1 * hrs) {
      // 60 mins/1hr
      step = 20;
    } else if (diff > 1 * hrs && diff <= 3 * hrs) {
      // 3 hrs
      step = 1 * min;
    } else if (diff > 3 * hrs && diff <= 6 * hrs) {
      // 6 hrs
      step = 2 * min;
    } else if (diff > 6 * hrs && diff <= 1 * days) {
      // 24 hrs/1 day
      step = 8 * min;
    } else if (diff > 1 * days && diff <= 2 * days) {
      // 2 days
      step = 16 * min;
    } else if (diff > 2 * days && diff <= 4 * days) {
      // 4 days
      step = 32 * min;
    } else if (diff > 4 * days && diff <= 7 * days) {
      // 7 days
      step = 56 * min;
    } else if (diff > 7 * days && diff <= 15 * days) {
      // 15 days
      step = 2 * hrs;
    } else if (diff > 15 * days && diff <= 1 * month) {
      // 30 days/1 month
      step = 4 * hrs;
    } else if (diff > 1 * month && diff <= 3 * month) {
      // 3 months
      step = 12 * hrs;
    } else if (diff > 3 * month && diff <= 6 * month) {
      // 6 months
      step = 1 * days;
    } else if (diff > 6 * month && diff <= 1 * year) {
      // 1 year/12 months
      step = 2 * days;
    } else if (diff > 1 * year && diff <= 2 * year) {
      // 2 years
      step = 4 * days;
    } else if (diff > 2 * year && diff <= 5 * year) {
      // 5 years
      step = 10 * days;
    } else {
      step = 30 * days;
    }
    return step;
  };

  // Transform data for chart
  const transformDataForChart = (data) => {
    if (
      data &&
      data.status === 'success' &&
      data.data &&
      data.data.resultType &&
      data.data.resultType === 'matrix' &&
      data.data.result &&
      data.data.result.length > 0
    ) {
      const fullData = [];
      data.data.result.forEach((r) => {
        const localData = r.values.map((arr) => {
          const x = moment(arr[0] * 1000).format(timeFormat);
          const y = parseFloat(parseFloat(arr[1]).toFixed(2));
          return { x, y };
        });
        fullData.push({ data: localData, metric: r.metric });
      });
      return fullData;
    }
    return [];
  };

  // Handle error
  const handleError = (error) => {
    updateProgress({ showProgress: false });
    if (error) {
      setError(error.message && error.message !== '' ? error.message : error !== '' ? error : '');
      setErrorCount((prevCount) => prevCount + 1);
    }
  };

  // Update date range function
  const updateDateRangeFunc = () => {
    return function (domain) {
      if (domain.length === 2) {
        const min = domain[0];
        const max = domain[1];
        updateDateRange(
          `${min.getTime().toString()}`,
          min,
          `${max.getTime().toString()}`,
          max,
          false,
          refresh,
        );
      }
    };
  };

  // Create chart options
  const createOptions = (xAxis, chartData, groups) => {
    // const fromDate = grafanaDateRangeToDate(from);
    // const toDate = grafanaDateRangeToDate(to);

    // const showAxis = panel.type ==='singlestat' && panel.sparkline && panel.sparkline.show === true?false:true;
    const showAxis = !(panel.type === 'singlestat' && panel.sparkline);

    const xAxes = {
      type: 'timeseries',
      // type : 'category',
      show: showAxis && !isSparkline,
      tick: {
        // format: c3TimeFormat,
        // fit: true,
        fit: false,
        count: 5,
        // centered: true
      },
      // label: {
      //   text: "X Label",
      //   position: "outer-center"
      // }
    };

    const yAxes = { show: showAxis && !isSparkline };
    if (panel.yaxes) {
      panel.yaxes.forEach((ya) => {
        if (typeof ya.label !== 'undefined' && ya.label !== null) {
          yAxes.label = { text: ya.label, position: 'outer-middle' };
        }
        if (ya.format.toLowerCase().startsWith('percent')) {
          const mulFactor = ya.format.toLowerCase() === 'percentunit' ? 100 : 1;
          yAxes.tick = {
            format(d) {
              const tk = (d * mulFactor).toFixed(2);
              return `${tk}%`;
            },
          };
        }
      });
    }

    const grid = {
      // x: {
      //     show: showAxis,
      // },
      // y: {
      //     show: showAxis,
      // }
    };

    const linked = isSparkline
      ? false
      : !inDialog
        ? { name: board && board.title ? board.title : '' }
        : false;

    let shouldDisplayLegend = Object.keys(datasetIndex).length <= 10;
    if (panel.type !== 'graph') {
      shouldDisplayLegend = false;
    }

    if (chartRef && chartRef.current) {
      const chartConfig = {
        // oninit: function(args){
        //   console.log(JSON.stringify(args));
        // },
        bindto: chartRef.current,
        size: isSparkline
          ? {
              // width: 150,
              height: 50,
            }
          : null,
        data: {
          x: 'x',
          xFormat: bbTimeFormat,
          columns: [xAxis, ...chartData],
          groups,
          type: isSparkline ? line() : area(),
        },
        axis: { x: xAxes, y: yAxes },
        zoom: { enabled: true, type: 'drag', onzoomend: updateDateRangeFunc() },
        grid,
        legend: { show: shouldDisplayLegend && !isSparkline },

        point: { r: 0, focus: { expand: { r: 5 } } },
        tooltip: {
          show: showAxis,
          linked,
          format: {
            title(x) {
              // return d3.timeFormat(bbTimeFormat)(x);
              return moment(x).format(timeFormat);
            },
          },
        },
        area: { linearGradient: true },
      };

      if (
        panelType === 'sparkline' &&
        panel.sparkline &&
        panel.sparkline.lineColor &&
        panel.sparkline.fillColor
      ) {
        // cd.datasets[datasetInd].borderColor = panel.sparkline.lineColor;
        // cd.datasets[datasetInd].backgroundColor = panel.sparkline.fillColor;

        const dataLength =
          chartConfig.data.columns &&
          Array.isArray(chartConfig.data.columns) &&
          chartConfig.data.columns.length > 1
            ? chartConfig.data.columns[1].length
            : 0; // 0 is for x axis
        if (dataLength > 0) {
          if (typeof chartConfig.data.colors === 'undefined') {
            chartConfig.data.colors = {};
          }
          chartConfig.data.colors[chartConfig.data.columns[1][0]] = panel.sparkline.lineColor;

          if (dataLength > 1) {
            let content = '';
            if (panel.format.toLowerCase().startsWith('percent')) {
              const mulFactor = panel.format.toLowerCase() === 'percentunit' ? 100 : 1;
              if (!isNaN(chartConfig.data.columns[1][dataLength - 1])) {
                const tk = (chartConfig.data.columns[1][dataLength - 1] * mulFactor).toFixed(2);
                content = `${tk}%`;
              }
            } else {
              content = `${chartConfig.data.columns[1][dataLength - 1]} ${panel.format}`;
            }
            chartConfig.title = {
              text: `\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n${content}`, // for sparkline, we want to print the value as title
            };
          }
        }

        chartConfig.color = { pattern: [panel.sparkline.fillColor] };
      }

      const newChart = bb.generate(chartConfig);
      setChart(newChart);
    }
  };

  // Get data for chart
  const getData = async (ind, target, chartInst, datasource) => {
    let { xAxis: currentXAxis, chartData: currentChartData } = { xAxis, chartData };

    let queryRangeURL = '';
    let endpointURL = '';
    let endpointAPIKey = '';
    if (prometheusURL && prometheusURL !== '') {
      endpointURL = prometheusURL;
      queryRangeURL = `/api/prometheus/query_range/${connectionID}`;
    } else if (grafanaURL && grafanaURL !== '') {
      endpointURL = grafanaURL;
      endpointAPIKey = grafanaAPIKey;
      queryRangeURL = `/api/grafana/query_range/${connectionID}`;
    }

    let { expr } = target;
    if (templateVars && templateVars !== null && templateVars.length > 0) {
      templateVars.forEach((tv) => {
        const tvrs = tv.split('=');
        if (tvrs.length === 2) {
          expr = expr.replace(
            new RegExp(`$${tvrs[0]}`.replace(/[-/^$*+?.()|[\]{}]/g, '\\$&'), 'g'),
            tvrs[1],
          ); //eslint-disable-line
        }
      });
    }
    const start = Math.round(grafanaDateRangeToDate(from).getTime() / 1000);
    const end = Math.round(grafanaDateRangeToDate(to).getTime() / 1000);

    let ds = datasource?.charAt(0).toUpperCase() + datasource?.substring(1);

    let queryParams = `ds=${ds}&query=${encodeURIComponent(
      expr,
    )}&start=${start}&end=${end}&step=${computeStep(start, end)}`;
    if (testUUID && testUUID.trim() !== '') {
      queryParams += `&uuid=${encodeURIComponent(testUUID)}`; // static_chart=true ?
    }

    const processReceivedData = (result) => {
      updateProgress({ showProgress: false });

      if (typeof result == 'undefined' || result?.status != 'success') {
        return;
      }

      if (typeof result !== 'undefined') {
        const fullData = transformDataForChart(result);
        if (fullData.length === 0) {
          return;
        }

        const newChartData = [...currentChartData];
        let newXAxis = [...currentXAxis];

        fullData.forEach(({ metric, data }, di) => {
          const datasetInd = getOrCreateIndex(`${ind}_${di}`);
          const newData = [];

          // if (typeof cd.labels[datasetInd] === 'undefined' || typeof cd.datasets[datasetInd] === 'undefined'){
          let legend = typeof target.legendFormat !== 'undefined' ? target.legendFormat : '';
          if (legend === '') {
            legend = Object.keys(metric).length > 0 ? JSON.stringify(metric) : '';
          } else {
            Object.keys(metric).forEach((metricKey) => {
              legend = legend
                .replace(`{{${metricKey}}}`, metric[metricKey])
                .replace(`{{ ${metricKey} }}`, metric[metricKey]);
            });
            legend = legend
              .replace('{{ ', '')
              .replace('{{', '')
              .replace(' }}', '')
              .replace('}}', '');
          }

          // bb does NOT like labels which start with a number
          if (!isNaN(legend.charAt(0))) {
            legend = ` ${legend}`;
          }
          // if(legend.trim() === ''){
          //   legend = 'NO VALUE';
          // }
          newData.push(legend);
          newXAxis = ['x'];
          // }
          data.forEach(({ x, y }) => {
            newData.push(y);
            newXAxis.push(new Date(x));
          });
          newChartData[datasetInd] = newData;
        });

        let groups = [];
        if (typeof panel.stack !== 'undefined' && panel.stack) {
          const panelGroups = [];
          newChartData.forEach((y) => {
            if (y && y.length > 0) {
              panelGroups.push(y[0]); // just the label
            }
          });
          groups = [panelGroups];
        }

        let chartDataFiltered = newChartData.filter((element) => element !== undefined);

        if (chart && chart !== null) {
          chart.load({ columns: [newXAxis, ...chartDataFiltered] });
        } else {
          createOptions(newXAxis, chartDataFiltered, groups);
        }

        if (error) {
          setXAxis(newXAxis);
          setChartData(newChartData);
          setError('');
          setErrorCount(0);
        } else {
          setXAxis(newXAxis);
          setChartData(newChartData);
        }
      }
    };

    if (panelData && panelData[expr]) {
      processReceivedData(panelData[expr]);
    } else {
      queryParams += `&url=${encodeURIComponent(endpointURL)}&api-key=${encodeURIComponent(
        endpointAPIKey,
      )}`;
      dataFetch(
        `${queryRangeURL}?${queryParams}`,
        {
          method: 'GET',
          credentials: 'include',
          // headers: headers,
        },
        processReceivedData,
        handleError,
      );
    }
  };

  // Collect chart data
  const collectChartData = (chartInst) => {
    if (panel.targets) {
      panel.targets.forEach((target, ind) => {
        getData(ind, target, chartInst, target.datasource?.type);
      });
    }
  };

  // Config chart data
  const configChartData = () => {
    // Initialize datasetIndex
    const initialDatasetIndex = {};
    if (panel.targets) {
      panel.targets.forEach((target, ind) => {
        initialDatasetIndex[`${ind}_0`] = ind;
      });
    }
    setDatasetIndex(initialDatasetIndex);

    // Clear existing interval if it exists
    if (interval) {
      clearInterval(interval);
    }

    // Set new interval if liveTail is enabled
    if (liveTail) {
      const newInterval = setInterval(
        () => {
          collectChartData();
        },
        computeRefreshInterval(refresh) * 1000,
      );
      setIntervalState(newInterval);
    }

    // Initial data collection
    collectChartData();
  };

  // ComponentDidMount equivalent
  useEffect(() => {
    configChartData();

    // ComponentWillUnmount equivalent
    return () => {
      if (interval) {
        clearInterval(interval);
      }
    };
  }, []); // Empty dependency array means this runs once on mount

  // Render
  let loadingBar;
  let reloadButton;

  if (error) {
    createOptions([], [], []); // add empty data to charts
    loadingBar = (
      <Box sx={{ width: '100%' }}>
        <LinearProgress />
      </Box>
    );
  }

  if (errorCount > 3 * panel?.targets?.length && interval) {
    clearInterval(interval); // clearing the interval to prevent further calls to get chart data
    loadingBar = null;
    reloadButton = (
      <HeaderIcon
        key="Reload"
        aria-label="reload the Chart"
        color="inherit"
        onClick={() => configChartData()}
      >
        <CachedIcon />
      </HeaderIcon>
    );
  }

  const iconComponent = (
    <div>
      {reloadButton}
      <HeaderIcon
        key="chartDialog"
        aria-label="Open chart in a dialog"
        color="inherit"
        onClick={() => handleChartDialogOpen(board, panel, panelData)}
      >
        <OpenInNewIcon />
      </HeaderIcon>
    </div>
  );

  let mainChart;
  if (panelType === 'gauge') {
    mainChart = <GrafanaCustomGaugeChart data={chartData} panel={panel} error={error} />;
  } else {
    mainChart = (
      <ChartContainer>
        <div ref={chartRef} />
      </ChartContainer>
    );
  }

  //  if (isSparkline){
  //    return (
  //       <NoSsr>
  //       {loadingBar}
  //       <StyledSparklineContent>
  //         <div>{panel.title}</div>
  //         <div>{mainChart}</div>
  //         <div>{iconComponent}</div>
  //       </StyledSparklineContent>
  //     </NoSsr>
  //   )
  // }

  return (
    <NoSsr>
      {loadingBar}
      <StyledCard>
        {!inDialog && (
          <CardHeader
            disableTypography
            avatar={
              error && (
                <Tooltip title="There was an error communicating with the server" placement="top">
                  <WarningIcon component={ErrorText} />
                </Tooltip>
              )
            }
            title={panel.title}
            action={iconComponent}
            sx={{ fontSize: (theme) => theme.spacing(2), width: '100%' }}
          />
        )}
        <StyledCardContent>{mainChart}</StyledCardContent>
      </StyledCard>
    </NoSsr>
  );
}

GrafanaCustomChart.propTypes = {
  grafanaURL: PropTypes.string.isRequired,
  connectionID: PropTypes.string.isRequired,
  // grafanaAPIKey: PropTypes.string.isRequired,
  board: PropTypes.object.isRequired,
  panel: PropTypes.object.isRequired,
  templateVars: PropTypes.array.isRequired,
  updateDateRange: PropTypes.func.isRequired,
  handleChartDialogOpen: PropTypes.func.isRequired,
  inDialog: PropTypes.bool.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(null, mapDispatchToProps)(GrafanaCustomChart);
