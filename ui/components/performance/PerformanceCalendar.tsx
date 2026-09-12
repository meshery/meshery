import React, { useEffect, useState } from 'react';
import { momentLocalizer } from 'react-big-calendar';
import moment from 'moment';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import GenericModal from '../shared/Modal/GenericModal';
import MesheryChart from '../MesheryChart';
import { Typography, Paper } from '@sistent/sistent';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { CalendarComponent } from './style';
import { updateProgressAction } from '@/store/slices/mesheryUi';
import { useGetPerformanceResultsQuery } from '@meshery/schemas/mesheryApi';
import { useDispatch } from 'react-redux';

const localizer = momentLocalizer(moment);

/**
 * generateCalendarEventsFromResults takes in performance results data
 * and generate calendar events from them
 *
 * @param {{
 *  mesheryId: string,
 *  name: string,
 *  testStartTime: string,
 *  runnerResults: {
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
  // Map-then-filter (rather than filter-then-map) so `resource: index` keeps
  // pointing into the original results array for the click-to-detail lookup.
  return results
    .map(({ testStartTime, name, runnerResults }, index) => {
      // parseZone preserves the offset embedded in the timestamp instead of
      // converting to local time. Guard against missing/unparseable values:
      // parseZone(undefined) silently returns "now" and parseZone(null)
      // returns an invalid moment, both of which would surface bogus events.
      const startMoment = moment.parseZone(testStartTime);
      if (!testStartTime || !startMoment.isValid() || !runnerResults?.ActualDuration) return null;

      const ntzStartTime = new Date(startMoment.format('MM/DD/YYYY HH:mm'));
      const ntzEndTime = new Date(ntzStartTime);
      ntzEndTime.setSeconds(ntzEndTime.getSeconds() + runnerResults.ActualDuration / 1e9);

      return {
        title: name,
        start: ntzStartTime,
        end: ntzEndTime,
        resource: index,
      };
    })
    .filter(Boolean);
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
function PerformanceCalendar({ style }) {
  const [time, setTime] = useState(generateDateRange());
  const [results, setResults] = useState([]);
  const [selectedEvent, setSelectedEvent] = useState();
  const dispatch = useDispatch();
  const {
    data: performanceResultsData,
    isFetching,
    isError,
    error,
  } = useGetPerformanceResultsQuery({
    page: '0',
    pagesize: '10',
    search: '',
    order: '',
    from: time.start,
    to: time.end,
  });

  //hooks
  const { notify } = useNotification();

  useEffect(() => {
    dispatch(updateProgressAction({ showProgress: isFetching }));
  }, [dispatch, isFetching]);

  useEffect(() => {
    setResults(performanceResultsData?.results || []);
  }, [performanceResultsData]);

  useEffect(() => {
    if (isError) handleError('Failed to Fetch Results')(error);
  }, [isError, error]);

  function handleError(msg) {
    return function (error) {
      dispatch(updateProgressAction({ showProgress: false }));
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

    const row = result.runnerResults;
    if (!row) return <div />;

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
            rawdata={[result && result.runnerResults ? result : {}]}
            data={[result && result.runnerResults ? result.runnerResults : {}]}
          />
        </div>
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

export default PerformanceCalendar;
