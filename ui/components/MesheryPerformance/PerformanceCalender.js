//@ts-check
import React, { useState } from "react";
import moment from "moment";
import { Calendar, momentLocalizer } from "react-big-calendar";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { Paper } from "@mui/material";
import { Typography } from "@mui/material";
import GenericModal from "../GenericModal";
import GrafanaCustomCharts from "../GrafanaCustomCharts";
import MesheryChart from "../MesheryChart";
const localizer = momentLocalizer(moment);

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
    const ntzStartTime = new Date(moment(test_start_time).utcOffset(test_start_time).format("MM/DD/YYYY HH:mm"));
    const ntzEndTime = ntzStartTime;
    ntzEndTime.setSeconds(ntzEndTime.getSeconds() + runner_results.ActualDuration / 1e9);

    return {
      title: name,
      start: ntzStartTime,
      end: ntzEndTime,
      resource: index,
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
    return { start: moment(from).format("YYYY-MM-DD"), end: moment(to).format("YYYY-MM-DD") };
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
 * }} props
 * @returns
 */
function PerformanceCalendar({ style }) {
  const [time, setTime] = useState(generateDateRange());
  const [results, setResults] = useState([]);
  const [selectedEvent, setSelectedEvent] = useState();

  function handleEventClick(result) {
    setSelectedEvent(results[result.resource]);
  }

  const ResultChart = ({ result }) => {
    if (!result) return <div />;

    const row = result.runner_results;
    const boardConfig = result.server_board_config;
    const serverMetrics = result.server_metrics;
    const startTime = new Date(row.StartTime);
    const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
    return (
      <Paper style={{ width: "100%", maxWidth: "90vw", padding: "0.5rem" }}>
        <div>
          <Typography variant="h6" gutterBottom align="center">
            Performance Graph
          </Typography>
          <MesheryChart
            rawdata={[result && result.runner_results ? result : {}]}
            data={[result && result.runner_results ? result.runner_results : {}]}
          />
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
        style={{ height: "100%" }}
        onRangeChange={(range) => setTime(generateDateRange(range.start, range.end))}
        onSelectEvent={(results) => handleEventClick(results)}
      />

      <GenericModal
        open={!!selectedEvent}
        Content={<ResultChart result={selectedEvent} />}
        handleClose={() => setSelectedEvent(undefined)}
      />
    </div>
  );
}

export default PerformanceCalendar;
