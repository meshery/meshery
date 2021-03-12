//@ts-check
import React, { useEffect, useState } from "react";
import { Calendar, momentLocalizer } from "react-big-calendar";
import moment from "moment";
import { connect } from "react-redux";
import { updateProgress } from "../../lib/store";
import { bindActionCreators } from "redux";
import { promisifiedDataFetch } from "../../lib/data-fetch";
import { withSnackbar } from "notistack";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import "react-big-calendar/lib/css/react-big-calendar.css";

const localizer = momentLocalizer(moment);

const PERFORMANCE_PROFILE_RESULTS_URL = "/api/user/performance/profiles/results";

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
 *  end?: Date
 * }[]}
 */
function generateCalendarEventsFromResults(results) {
  return results.map(({ test_start_time, name, runner_results }) => {
    // Remove incorrect timezone info
    const ntzStartTime = new Date(moment(test_start_time).utcOffset(test_start_time).format('MM/DD/YYYY HH:mm'));
    const ntzEndTime = ntzStartTime;
    ntzEndTime.setSeconds(ntzEndTime.getSeconds() + (runner_results.ActualDuration / 1e9));

    return {
      title: name,
      start: ntzStartTime,
      end: ntzEndTime,
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
      start: moment(from).format("YYYY-MM-DD"),
      end: moment(to).format("YYYY-MM-DD"),
    };
  }

  return {
    start: moment().startOf("M").format("YYYY-MM-DD"),
    end: moment().add(1, "M").startOf("M").format("YYYY-MM-DD"),
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
function PerformanceCalendar({ style, updateProgress, enqueueSnackbar, closeSnackbar }) {
  const [time, setTime] = useState(generateDateRange());
  const [results, setResults] = useState([]);

  useEffect(() => {
    fetchResults(time.start, time.end);
  }, [time]);

  async function fetchResults(start, end) {
    updateProgress({ showProgress: true });

    try {
      const res = await promisifiedDataFetch(
        `${PERFORMANCE_PROFILE_RESULTS_URL}?from=${start}&to=${end}`, 
        { credentials: "include" }
      );
      updateProgress({ showProgress: false });
      if (res) {
        setResults(res.results || []);
      }
    } catch (error) {
      console.error(error)
      handleError("Failed to Fetch Profiles")(error)
    }
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress: false });

      enqueueSnackbar(`${msg}: ${error}`, {
        variant: "error",
        action: function Action(key) {
          return (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          );
        },
        autoHideDuration: 8000,
      });
    };
  }

  return (
    <div style={style}>
      <Calendar
        events={generateCalendarEventsFromResults(results)}
        views={["month", "week", "day"]}
        step={60}
        showMultiDayTimes
        defaultDate={new Date()}
        localizer={localizer}
        style={{
          height: "100%",
        }}
        // @ts-ignore
        onRangeChange={(range) => setTime(generateDateRange(range.start, range.end))}
      />
    </div>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(null, mapDispatchToProps)(withSnackbar(PerformanceCalendar));
