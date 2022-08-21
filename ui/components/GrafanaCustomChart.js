import { Component, useEffect } from "react";
import PropTypes from "prop-types";
import { withStyles } from "@mui/styles";
import { NoSsr, IconButton, Card, CardContent, CardHeader, Tooltip, LinearProgress, Box } from "@mui/icons-material";
import moment from "moment";
import OpenInNewIcon from "@mui/icons-material/OpenInNewOutlined";
import WarningIcon from "@mui/icons-material/Warning";
import CachedIcon from "@mui/icons-material/Cached";
import GrafanaCustomGaugeChart from "./GrafanaCustomGaugeChart";

import bb, { area, line } from "billboard.js";
import { interval } from "rxjs";

const grafanaStyles = (theme) => ({
  chart: { width: "100%" },
  column: { flex: "1" },
  heading: { fontSize: theme.typography.pxToRem(15) },
  secondaryHeading: { fontSize: theme.typography.pxToRem(15), color: theme.palette.text.secondary },
  dateRangePicker: {
    display: "flex",
    justifyContent: "flex-end",
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(2),
  },
  cardHeader: { fontSize: theme.spacing(2), width: "100%" },
  cardHeaderIcon: { fontSize: theme.spacing(2) },
  card: { height: "100%", width: "100%" },
  sparklineCardContent: { display: "grid", gap: " 0.5rem" },
  cardContent: { height: "100%", width: "100%" },
  error: { color: "#D32F2F" },
});

const grafanaDateRangeToDate = (dt, startDate) => {
  const dto = new Date();
  switch (dt) {
    case "now-2d":
      dto.setDate(dto.getDate() - 2);
      break;
    case "now-7d":
      dto.setDate(dto.getDate() - 7);
      break;
    case "now-30d":
      dto.setDate(dto.getDate() - 30);
      break;
    case "now-90d":
      dto.setDate(dto.getDate() - 90);
      break;
    case "now-6M":
      dto.setMonth(dto.getMonth() - 6);
      break;
    case "now-1y":
      dto.setFullYear(dto.getFullYear() - 1);
      break;
    case "now-2y":
      dto.setFullYear(dto.getFullYear() - 2);
      break;
    case "now-5y":
      dto.setFullYear(dto.getFullYear() - 5);
      break;
    case "now-1d/d":
      dto.setDate(dto.getDate() - 1);
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now-2d/d":
      dto.setDate(dto.getDate() - 2);
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now-7d/d":
      dto.setDate(dto.getDate() - 7);
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now-1w/w":
      dto.setDate(dto.getDate() - 6 - ((dto.getDay() + 8) % 7));
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setDate(dto.getDate() + 6);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now-1M/M":
      dto.setMonth(dto.getMonth() - 1);
      if (startDate) {
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(dto.getMonth());
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now-1y/y":
      dto.setFullYear(dto.getFullYear() - 1);
      if (startDate) {
        dto.setMonth(0);
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(12);
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now/d":
      dto.setDate(dto.getDate() - 6 - ((dto.getDay() + 8) % 7));
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now":
      break;
    case "now/w":
      dto.setDate(dto.getDate() - ((dto.getDay() + 7) % 7));
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setDate(dto.getDate() + 6);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now/M":
      if (startDate) {
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(dto.getMonth() + 1);
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now/y":
      if (startDate) {
        dto.setMonth(0);
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(12);
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case "now-5m":
      dto.setMinutes(dto.getMinutes() - 5);
      break;
    case "now-15m":
      dto.setMinutes(dto.getMinutes() - 15);
      break;
    case "now-30m":
      dto.setMinutes(dto.getMinutes() - 30);
      break;
    case "now-1h":
      dto.setHours(dto.getHours() - 1);
      break;
    case "now-3h":
      dto.setHours(dto.getHours() - 3);
      break;
    case "now-6h":
      dto.setHours(dto.getHours() - 6);
      break;
    case "now-12h":
      dto.setHours(dto.getHours() - 12);
      break;
    case "now-24h":
      dto.setHours(dto.getHours() - 24);
      break;
    default:
      return new Date(parseFloat(dt));
  }
  return dto;
};

const GrafanaCustomChart = (props) => {
  chartRef = null;
  chart = null;
  timeFormat = "MM/DD/YYYY HH:mm:ss";
  bbTimeFormat = "%Y-%m-%d %h:%M:%S %p";
  panelType = "";
  switch (props.panel.type) {
    case "graph":
      panelType = props.panel.type;
      break;
    case "singlestat":
      panelType =
        props.panel.type === "singlestat" && props.panel.sparkline && props.panel.sparkline.show === true
          ? "sparkline"
          : "gauge";
      // panelType = props.panel.type ==='singlestat' && props.panel.sparkline ? 'sparkline':'gauge';
      break;
  }
  const { sparkline } = props;
  datasetIndex = {};
  state = {
    xAxis: [],
    sparkline: !!sparkline,
    chartData: [],
    error: "",
    errorCount: 0,
  };

  useEffect(() => {
    configChartData();
  }, []);

  useEffect(() => {
    if (typeof interval !== "undefined") {
      clearInterval(interval);
    }
  }, [interval]);

  const configChartData = () => {
    const { panel, refresh, liveTail } = props;

    if (panel.targets) {
      panel.targets.forEach((target, ind) => {
        datasetIndex[`${ind}_0`] = ind;
      });
    }
    if (typeof interval !== "undefined") {
      clearInterval(interval);
    }
    if (liveTail) {
      interval = setInterval(() => {
        collectChartData();
      }, computeRefreshInterval(refresh) * 1000);
    }
    collectChartData();
  };

  const getOrCreateIndex = (datasetInd) => {
    if (typeof datasetIndex[datasetInd] !== "undefined") {
      return datasetIndex[datasetInd];
    }
    let max = 0;
    Object.keys(datasetIndex).forEach((i) => {
      if (datasetIndex[i] > max) {
        max = datasetIndex[i];
      }
    });
    datasetIndex[datasetInd] = max + 1;
    return max + 1;
  };

  const collectChartData = (chartInst) => {
    const { panel } = props;
    if (panel.targets) {
      panel.targets.forEach((target, ind) => {
        getData(ind, target, chartInst);
      });
    }
  };

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

  const getData = async (ind, target) => {
    const { prometheusURL, grafanaURL, grafanaAPIKey, panel, from, to, templateVars, testUUID, panelData } = props;
    const { chartData } = state;
    let { xAxis } = state;

    let queryRangeURL = "";
    let endpointURL = "";
    let endpointAPIKey = "";
    if (prometheusURL && prometheusURL !== "") {
      endpointURL = prometheusURL;
      queryRangeURL = "/api/prometheus/query_range";
    } else if (grafanaURL && grafanaURL !== "") {
      endpointURL = grafanaURL;
      endpointAPIKey = grafanaAPIKey;
      queryRangeURL = "/api/grafana/query_range";
    }

    let { expr } = target;
    if (templateVars && templateVars !== null && templateVars.length > 0) {
      templateVars.forEach((tv) => {
        const tvrs = tv.split("=");
        if (tvrs.length == 2) {
          expr = expr.replace(new RegExp(`$${tvrs[0]}`.replace(/[-\/\\^$*+?.()|[\]{}]/g, "\\$&"), "g"), tvrs[1]); //eslint-disable-line
        }
      });
    }
    const start = Math.round(grafanaDateRangeToDate(from).getTime() / 1000);
    const end = Math.round(grafanaDateRangeToDate(to).getTime() / 1000);
    let queryParams = `ds=${panel.datasource}&query=${encodeURIComponent(
      expr
    )}&start=${start}&end=${end}&step=${computeStep(start, end)}`;
    if (testUUID && testUUID.trim() !== "") {
      queryParams += `&uuid=${encodeURIComponent(testUUID)}`; // static_chart=true ?
    }

    const processReceivedData = (result) => {
      props.updateProgress({ showProgress: false });

      if (typeof result == "undefined" || result?.status != "success") {
        return;
      }

      if (typeof result !== "undefined") {
        const fullData = transformDataForChart(result);
        if (fullData.length === 0) {
          return;
        }

        fullData.forEach(({ metric, data }, di) => {
          const datasetInd = getOrCreateIndex(`${ind}_${di}`);
          const newData = [];

          // if (typeof cd.labels[datasetInd] === 'undefined' || typeof cd.datasets[datasetInd] === 'undefined'){
          let legend = typeof target.legendFormat !== undefined ? target.legendFormat : "";
          if (legend === "") {
            legend = Object.keys(metric).length > 0 ? JSON.stringify(metric) : "";
          } else {
            Object.keys(metric).forEach((metricKey) => {
              legend = legend
                .replace(`{{${metricKey}}}`, metric[metricKey])
                .replace(`{{ ${metricKey} }}`, metric[metricKey]);
            });
            legend = legend.replace("{{ ", "").replace("{{", "").replace(" }}", "").replace("}}", "");
          }

          // bb does NOT like labels which start with a number
          if (!isNaN(legend.charAt(0))) {
            legend = ` ${legend}`;
          }
          // if(legend.trim() === ''){
          //   legend = 'NO VALUE';
          // }
          newData.push(legend);
          xAxis = ["x"];
          // }
          data.forEach(({ x, y }) => {
            newData.push(y);
            xAxis.push(new Date(x));
          });
          chartData[datasetInd] = newData;
        });
        let groups = [];
        if (typeof panel.stack !== "undefined" && panel.stack) {
          const panelGroups = [];
          chartData.forEach((y) => {
            if (y.length > 0) {
              panelGroups.push(y[0]); // just the label
            }
          });
          groups = [panelGroups];
        }
        let chartDataFiltered = chartData.filter((element) => element !== undefined);
        if (chart && chart !== null) {
          chart.load({ columns: [xAxis, ...chartDataFiltered] });
        } else {
          createOptions(xAxis, chartDataFiltered, groups);
        }
        state.error &&
          setState({
            xAxis,
            chartData,
            error: "",
            errorCount: 0,
          });
      }
    };

    if (panelData && panelData[expr]) {
      processReceivedData(panelData[expr]);
    } else {
      queryParams += `&url=${encodeURIComponent(endpointURL)}&api-key=${encodeURIComponent(endpointAPIKey)}`;
      dataFetch(
        `${queryRangeURL}?${queryParams}`,
        {
          method: "GET",
          credentials: "include",
          // headers: headers,
        },
        processReceivedData,
        handleError
      );
    }
  };

  const transformDataForChart = (data) => {
    if (
      data &&
      data.status === "success" &&
      data.data &&
      data.data.resultType &&
      data.data.resultType === "matrix" &&
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

  const updateDateRange = () => {
    return function (domain) {
      if (domain.length === 2) {
        const min = domain[0];
        const max = domain[1];
        props.updateDateRange(
          `${min.getTime().toString()}`,
          min,
          `${max.getTime().toString()}`,
          max,
          false,
          props.refresh
        );
      }
    };
  };

  // createOptions() {
  //   const {panel, from, to, panelData} = props;
  //   const fromDate = grafanaDateRangeToDate(from);
  //   const toDate = grafanaDateRangeToDate(to);
  const createOptions = (xAxis, chartData, groups) => {
    const { panel, board, inDialog } = props;

    // const showAxis = panel.type ==='singlestat' && panel.sparkline && panel.sparkline.show === true?false:true;
    const showAxis = !(panel.type === "singlestat" && panel.sparkline);

    const xAxes = {
      type: "timeseries",
      // type : 'category',
      show: showAxis && !state.sparkline,
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

    const yAxes = { show: showAxis && !state.sparkline };
    if (panel.yaxes) {
      panel.yaxes.forEach((ya) => {
        if (typeof ya.label !== "undefined" && ya.label !== null) {
          yAxes.label = { text: ya.label, position: "outer-middle" };
        }
        if (ya.format.toLowerCase().startsWith("percent")) {
          const mulFactor = ya.format.toLowerCase() === "percentunit" ? 100 : 1;
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

    const linked = state.sparkline ? false : !inDialog ? { name: board && board.title ? board.title : "" } : false;

    let shouldDisplayLegend = Object.keys(datasetIndex).length <= 10;
    if (panel.type !== "graph") {
      shouldDisplayLegend = false;
    }
    if (chartRef && chartRef !== null) {
      const chartConfig = {
        // oninit: function(args){
        //   console.log(JSON.stringify(args));
        // },
        bindto: chartRef,
        size: state.sparkline
          ? {
              // width: 150,
              height: 50,
            }
          : null,
        data: {
          x: "x",
          xFormat: bbTimeFormat,
          columns: [xAxis, ...chartData],
          groups,
          type: state.sparkline ? line() : area(),
        },
        axis: { x: xAxes, y: yAxes },
        zoom: { enabled: true, type: "drag", onzoomend: updateDateRange() },
        grid,
        legend: { show: shouldDisplayLegend && !state.sparkline },

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

      if (panelType === "sparkline" && panel.sparkline && panel.sparkline.lineColor && panel.sparkline.fillColor) {
        // cd.datasets[datasetInd].borderColor = panel.sparkline.lineColor;
        // cd.datasets[datasetInd].backgroundColor = panel.sparkline.fillColor;

        const dataLength =
          chartConfig.data.columns && Array.isArray(chartConfig.data.columns) && chartConfig.data.columns.length > 1
            ? chartConfig.data.columns[1].length
            : 0; // 0 is for x axis
        if (dataLength > 0) {
          if (typeof chartConfig.data.colors === "undefined") {
            chartConfig.data.colors = {};
          }
          chartConfig.data.colors[chartConfig.data.columns[1][0]] = panel.sparkline.lineColor;

          if (dataLength > 1) {
            let content = "";
            if (panel.format.toLowerCase().startsWith("percent")) {
              const mulFactor = panel.format.toLowerCase() === "percentunit" ? 100 : 1;
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

      chart = bb.generate(chartConfig);
    }
  };

  const computeRefreshInterval = (refresh) => {
    refresh = refresh.toLowerCase();
    const l = refresh.length;
    const dur = refresh.substring(l - 1, l);
    refresh = refresh.substring(0, l - 1);
    let val = parseInt(refresh);
    if (dur === "d") {
      val *= 24;
    }
    if (dur === "h") {
      val *= 60;
    }
    if (dur === "m") {
      val *= 60;
    }
    if (dur === "s") {
      return val;
    }
    return 30; //fallback
  };

  const handleError = (error) => {
    props.updateProgress({ showProgress: false });
    if (error) {
      setState({
        error: error.message && error.message !== "" ? error.message : error !== "" ? error : "",
        errorCount: state.errorCount + 1,
      });
    }
  };

  const { classes, board, panel, inDialog, handleChartDialogOpen, panelData } = props;
  const { error, errorCount, chartData } = state;

  let loadingBar;
  let reloadButton;

  if (error) {
    createOptions([], [], []); // add empty data to charts
    loadingBar = (
      <Box sx={{ width: "100%" }}>
        <LinearProgress />
      </Box>
    );
  }

  if (errorCount > 3 * panel.targets.length && typeof interval !== "undefined") {
    clearInterval(interval); // clearing the interval to prevent further calls to get chart data
    loadingBar = null;
    reloadButton = (
      <IconButton key="Relaod" aria-label="reloadButton the Chart" color="inherit" onClick={() => configChartData()}>
        <CachedIcon className={classes.cardHeaderIcon} />
      </IconButton>
    );
  }

  const iconComponent = (
    <div>
      {reloadButton}
      <IconButton
        key="chartDialog"
        aria-label="Open chart in a dialog"
        color="inherit"
        onClick={() => handleChartDialogOpen(board, panel, panelData)}
      >
        <OpenInNewIcon className={classes.cardHeaderIcon} />
      </IconButton>
    </div>
  );

  let mainChart;
  if (panelType === "gauge") {
    mainChart = <GrafanaCustomGaugeChart data={chartData} panel={panel} error={error} />;
  } else {
    mainChart = (
      <div>
        <div ref={(ch) => (chartRef = ch)} className={classes.chart} />
      </div>
    );
  }
  // if (state.sparkline){
  //   return (
  //     <NoSsr>
  //       {loadingBar}
  //       <div className={classes.sparklineCardContent}>
  //         <div>{panel.title}</div>
  //         <div>{mainChart}</div>
  //         <div>{iconComponent}</div>
  //       </div>
  //     </NoSsr>
  //   )
  // }
  return (
    <NoSsr>
      {loadingBar}
      <Card className={classes.card}>
        {!inDialog && (
          <CardHeader
            disableTypography
            avatar={
              error && (
                <Tooltip title="There was an error communicating with the server" placement="top">
                  <WarningIcon className={classes.error} />
                </Tooltip>
              )
            }
            title={panel.title}
            action={iconComponent}
            className={classes.cardHeader}
          />
        )}
        <CardContent className={classes.cardContent}>{mainChart}</CardContent>
      </Card>
    </NoSsr>
  );
};

GrafanaCustomChart.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  // grafanaAPIKey: PropTypes.string.isRequired,
  board: PropTypes.object.isRequired,
  panel: PropTypes.object.isRequired,
  templateVars: PropTypes.array.isRequired,
  updateDateRange: PropTypes.func.isRequired,
  handleChartDialogOpen: PropTypes.func.isRequired,
  inDialog: PropTypes.bool.isRequired,
};


export default withStyles(grafanaStyles)(GrafanaCustomChart);
