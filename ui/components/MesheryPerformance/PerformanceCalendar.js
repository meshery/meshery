//@ts-check
import React, { useEffect, useState } from "react";
import { Calendar, momentLocalizer } from "react-big-calendar";
import moment from "moment";
import { connect } from "react-redux";
import { updateProgress } from "../../lib/store";
import { bindActionCreators } from "redux";
import { withSnackbar } from "notistack";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import "react-big-calendar/lib/css/react-big-calendar.css";
import GenericModal from "../GenericModal";
import GrafanaCustomCharts from "../telemetry/grafana/GrafanaCustomCharts";
import MesheryChart from "../MesheryChart";
import { Paper, withStyles } from "@material-ui/core";
import { Typography } from "@material-ui/core";
import fetchAllResults from '../graphql/queries/FetchAllResultsQuery'
import { iconMedium } from "../../css/icons.styles";

const localizer = momentLocalizer(moment);
const styles = (theme) => ({
  paper : { padding : "3rem" },
  resultContainer : {
    display : "flex",
    flexDirection : "row",
    justifyContent : "space-between",
    ["@media (max-width: 830px)"] : {
      flexDirection : "column"
    },
  },
  vSep : {
    height : "10.4rem",
    width : "1px",
    background : "black",
    marginTop : "1.1rem",
    bottom : "0",
    left : "36%",
    backgroundColor : "#36454f",
    opacity : "0.7",
    ["@media (max-width: 830px)"] : {
      display : "none",
    }
  },
  hSep : {
    display : "none",
    ["@media (max-width: 830px)"] : {
      display : "block",
      width : "100%",
      height : "1px",
      background : "black",
      marginTop : "1.1rem",
      bottom : "0",
      left : "36%",
      backgroundColor : "#36454f",
      opacity : "0.7",
    }
  },
  calender : {
    "& .rbc-off-range-bg" : {
      backgroundColor : theme.palette.type === 'dark' ? "#a9a9a9" : "#e6e6e6",
      color : theme.palette.type === 'dark' ? "white" : "black",
    },
    "& .rbc-off-range " : {
      color : theme.palette.secondary.lightText,
    },
    "& .rbc-btn-group" : {
      color : theme.palette.secondary.lightText,
    },
    "& .rbc-toolbar button" : {

      color : theme.palette.secondary.lightText,


    },
    "& .rbc-today" : {
      backgroundColor : theme.palette.type === 'dark' ? "#505050" : "#eaf6ff",
    },
    "& .rbc-day-slot .rbc-time-slot" : {
      borderTop : `1px solid ${theme.palette.type === 'dark' ? "#555555" : "#eaf6ff"}`,
    },
    "& .rbc-toolbar button.rbc-active" : {
      backgroundColor : theme.palette.type === 'dark' ? "#505050" : "#eaf6ff",
      color : theme.palette.type === 'dark' ? "white" : "black",
    },
    "& .rbc-toolbar button:hover" : {
      backgroundColor : theme.palette.type === 'dark' ? "rgba(0, 0, 0, 0.54)" : "rgba(255, 255, 255, 0.54)"
    },
    "& .rbc-toolbar button.rbc-active:hover" : {
      backgroundColor : theme.palette.type === 'dark' ? "#909090" : "#eaf6ff",
    },
    "& .rbc-toolbar button:focus " : {
      backgroundColor : theme.palette.type === 'dark' ? "#505050" : "#eaf6ff",
      color : theme.palette.type === 'dark' ? "white" : "black",
    },
  },
  resultText : {
    color : theme.palette.secondary.lightText,
  },
  profileText : {
    color : theme.palette.secondary.lightText,
  },
});
// const PERFORMANCE_PROFILE_RESULTS_URL = "/api/user/performance/profiles/results";

/**
 * generateCalendarEventsFromResults takes in performance results data
 * and generate calendar events from them
 *
 * @param {{
 *  meshery_id: string,
 *  name: string,
 *  test_start_time: string,
 *  runner_results: {
 *    ActualDuration: number
 *  }
 * }[]} results performance results
 *
 * @returns {{
 *  title?: string,
 *  start?: Date,
 *  end?: Date,
 *  resource?: any
 * }[]}
 */
function generateCalendarEventsFromResults(results) {
  return results.map(({ test_start_time, name, runner_results }, index) => {
    // Remove incorrect timezone info
    const ntzStartTime = new Date(moment(test_start_time).utcOffset(test_start_time).format('MM/DD/YYYY HH:mm'));
    const ntzEndTime = ntzStartTime;
    ntzEndTime.setSeconds(ntzEndTime.getSeconds() + (runner_results.ActualDuration / 1e9));

    return {
      title : name,
      start : ntzStartTime,
      end : ntzEndTime,
      resource : index
    };
  });
}

/**
 * generateDateRange generates date range from the given dates
 * in yyyy-mm-dd
 * @param {Date} [from] start date
 * @param {Date} [to] end date
 * @returns {{
 *  start: string,
 *  end: string
 * }}
 */
function generateDateRange(from, to) {
  if (from && to) {
    return {
      start : moment(from).format("YYYY-MM-DD"),
      end : moment(to).format("YYYY-MM-DD"),
    };
  }

  return {
    start : moment().startOf("M").format("YYYY-MM-DD"),
    end : moment().add(1, "M").startOf("M").format("YYYY-MM-DD"),
  };
}

/**
 * PerformanceCalendar renders a calendar with details like tests run
 * @param {{
 *  style?: React.CSSProperties,
 *  updateProgress: any,
 *  enqueueSnackbar: any,
 *  closeSnackbar: any
 * }} props
 * @returns
 */
function PerformanceCalendar({
  style, updateProgress, enqueueSnackbar, closeSnackbar, classes
}) {
  const [time, setTime] = useState(generateDateRange());
  const [results, setResults] = useState([]);
  const [selectedEvent, setSelectedEvent] = useState();

  useEffect(() => {
    fetchResults(time.start, time.end);
  }, [time]);

  async function fetchResults(start, end) {
    updateProgress({ showProgress : true });

    fetchAllResults({
      selector : {
        // default
        pageSize : `10`,
        page : `0`,
        search : ``,
        order : ``,
        from : start,
        to : end,
      },
    }).subscribe({
      next : (res) => {
        // @ts-ignore
        let result = res?.fetchAllResults;
        updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          if (result) {
            // @ts-ignore
            setResults(result.results || []);
          }
        }
      },
      error : handleError("Failed to Fetch Profiles"),
    });
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress : false });

      enqueueSnackbar(`${msg}: ${error}`, {
        variant : "error",
        action : function Action(key) {
          return (
            <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon style={iconMedium}/>
            </IconButton>
          );
        },
        autoHideDuration : 8000,
      });
    };
  }

  function handleEventClick(result) {
    setSelectedEvent(results[result.resource]);
  }

  function ResultChart({ result }) {
    if (!result) return <div />;

    const row = result.runner_results;
    const boardConfig = result.server_board_config;
    const serverMetrics = result.server_metrics;
    const startTime = new Date(row.StartTime);
    const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
    return (
      <Paper
        style={{
          width : "100%",
          maxWidth : "90vw",
          padding : "0.5rem"
        }}
      >
        <div>
          <Typography variant="h6" gutterBottom align="center">Performance Graph</Typography>
          <MesheryChart rawdata={[result && result.runner_results ? result : {}]} data={[result && result.runner_results
            ? result.runner_results
            : {}]} />
        </div>
        {boardConfig && boardConfig !== null && Object.keys(boardConfig).length > 0 && (
          <div>
            <GrafanaCustomCharts
              boardPanelConfigs={[boardConfig]}
              // @ts-ignore
              boardPanelData={[serverMetrics]}
              startDate={startTime}
              from={startTime.getTime().toString()}
              endDate={endTime}
              to={endTime.getTime().toString()}
              liveTail={false}
            />
          </div>
        )}
      </Paper>
    );
  }
  return (
    <div style={style}>
      <Calendar
        events={generateCalendarEventsFromResults(results)}
        views={["month", "week", "day"]}
        defaultView="day"
        step={60}
        showMultiDayTimes
        defaultDate={new Date()}
        localizer={localizer}
        className={classes.calender}
        style={{ height : "100%", }}
        // @ts-ignore
        onRangeChange={(range) => setTime(generateDateRange(range.start, range.end))}
        onSelectEvent={(results) => handleEventClick(results)}
      />

      <GenericModal
        open={!!selectedEvent}
        // @ts-ignore
        Content={<ResultChart result={selectedEvent} />}
        handleClose={() => setSelectedEvent(undefined)}
      />
    </div>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default withStyles(styles)(connect(null, mapDispatchToProps)(withSnackbar(PerformanceCalendar)));
