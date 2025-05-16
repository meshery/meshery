import React, { useEffect, useState } from 'react';
import { momentLocalizer } from 'react-big-calendar';
import moment from 'moment';
import { connect } from 'react-redux';
import { updateProgress } from '../../lib/store';
import { bindActionCreators } from 'redux';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import GenericModal from '../GenericModal';
import GrafanaCustomCharts from '../telemetry/grafana/GrafanaCustomCharts';
import MesheryChart from '../MesheryChart';
import { Typography, Paper } from '@layer5/sistent';
import fetchAllResults from '../graphql/queries/FetchAllResultsQuery';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { CalendarComponent } from './style';

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
    const ntzStartTime = new Date(
      moment(test_start_time).utcOffset(test_start_time).format('MM/DD/YYYY HH:mm'),
    );
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
    return {
      start: moment(from).format('YYYY-MM-DD'),
      end: moment(to).format('YYYY-MM-DD'),
    };
  }

  return {
    start: moment().startOf('M').format('YYYY-MM-DD'),
    end: moment().add(1, 'M').startOf('M').format('YYYY-MM-DD'),
  };
}

/**
 * PerformanceCalendar renders a calendar with details like tests run
 * @param {{
 *  style?: React.CSSProperties,
 *  updateProgress: any,
 * }} props
 * @returns
 */
function PerformanceCalendar({ style, updateProgress }) {
  const [time, setTime] = useState(generateDateRange());
  const [results, setResults] = useState([]);
  const [selectedEvent, setSelectedEvent] = useState();

  //hooks
  const { notify } = useNotification();

  useEffect(() => {
    fetchResults(time.start, time.end);
  }, [time]);

  async function fetchResults(start, end) {
    updateProgress({ showProgress: true });

    fetchAllResults({
      selector: {
        // default
        pageSize: `10`,
        page: `0`,
        search: ``,
        order: ``,
        from: start,
        to: end,
      },
    }).subscribe({
      next: (res) => {
        // @ts-ignore
        let result = res?.fetchAllResults;
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          if (result) {
            // @ts-ignore
            setResults(result.results || []);
          }
        }
      },
      error: handleError('Failed to Fetch Profiles'),
    });
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress: false });
      notify({
        message: `${msg}: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
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
          width: '100%',
          maxWidth: '90vw',
          padding: '0.5rem',
        }}
      >
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
      <CalendarComponent
        events={generateCalendarEventsFromResults(results)}
        views={['month', 'week', 'day']}
        defaultView="day"
        step={60}
        showMultiDayTimes
        defaultDate={new Date()}
        localizer={localizer}
        style={{ height: '100%' }}
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

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(null, mapDispatchToProps)(PerformanceCalendar);
