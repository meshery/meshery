import { createEntityAdapter, createSlice } from '@reduxjs/toolkit';
import { SEVERITY, STATUS } from '../../components/layout/NotificationCenter/constants';

const initialState = {
  current_view: {
    page: 0,
    pagesize: 10,
    filters: {
      initial: true,
    },
    has_more: true,
    // Set when a page request fails so the UI can show an error + retry without
    // treating the failure as the end of the list.
    load_error: false,
    failed_request: null,
  },

  // this is used to fetch the view when notification center is opened
  view_to_fetch_on_open: {
    page: 0,
    filters: {
      status: STATUS.UNREAD,
    },
  },

  ui: {
    history_mode: false,
    title: 'Notifications',
    empty_message: 'No notifications found',
  },
  isNotificationCenterOpen: false,
};

export const isSerializable = (value, seen = new Set()) => {
  if (value === null || value === undefined) {
    return true;
  }
  const type = typeof value;
  if (type === 'string' || type === 'number' || type === 'boolean') {
    return true;
  }
  if (type !== 'object' || seen.has(value)) {
    return false;
  }
  const isArray = Array.isArray(value);
  const prototype = Object.getPrototypeOf(value);
  if (!isArray && prototype !== Object.prototype && prototype !== null) {
    return false;
  }
  seen.add(value);
  // Arrays can also carry named properties (e.g. list.onSelect = fn), which
  // iterating the array itself would miss, so read them through its keys.
  const entries = isArray ? Object.keys(value).map((key) => value[key]) : Object.values(value);
  const result = entries.every((entry) => isSerializable(entry, seen));
  seen.delete(value);
  return result;
};

const serializableUi = (ui) =>
  Object.fromEntries(Object.entries(ui || {}).filter(([, value]) => isSerializable(value)));

const defaultEventProperties = {
  severity: SEVERITY.INFO,
  status: STATUS.UNREAD,
};

const eventsEntityAdapter = createEntityAdapter({
  selectId: (event) => event.id,
  //sort based on createdAt timestamp(utc)
  sortComparer: (a, b) => {
    if (b?.createdAt?.localeCompare && a?.createdAt?.localeCompare) {
      return b.createdAt?.localeCompare(a.createdAt);
    }
    return 0;
  },
});

// An RFC 3339 time as Meshery Server's Go time.Time marshals one: whole
// seconds, an optional fraction of up to nine digits, and a zone.
const RFC3339 = /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.(\d{1,9}))?(Z|[+-]\d{2}:\d{2})$/;

const NANOSECONDS_PER_MILLISECOND = 1_000_000;

// instantOf returns value as [whole seconds since the epoch, in milliseconds;
// nanoseconds within that second], or undefined when it is not a time. Go
// sends nanosecond precision, which Date.parse would cut to milliseconds, so
// two revisions written within one millisecond would compare equal; the
// fraction is therefore read as integer nanoseconds. A time in another format
// falls back to Date.parse at millisecond precision.
const instantOf = (value) => {
  if (typeof value !== 'string') {
    return undefined;
  }
  const match = RFC3339.exec(value);
  if (match) {
    const [, wholeSeconds, fraction = '', zone] = match;
    const seconds = Date.parse(`${wholeSeconds}${zone}`);
    return Number.isNaN(seconds) ? undefined : [seconds, Number(fraction.padEnd(9, '0'))];
  }
  const at = Date.parse(value);
  if (Number.isNaN(at)) {
    return undefined;
  }
  const seconds = Math.floor(at / 1000) * 1000;
  return [seconds, (at - seconds) * NANOSECONDS_PER_MILLISECOND];
};

// compareTimes orders two times at full precision: negative when a is
// earlier, positive when later, and 0 when they are equal or either is not a
// time.
const compareTimes = (a, b) => {
  const x = instantOf(a);
  const y = instantOf(b);
  if (!x || !y) {
    return 0;
  }
  return x[0] - y[0] || x[1] - y[1];
};

const compareCounts = (a, b) => (typeof a === 'number' && typeof b === 'number' ? a - b : 0);

// compareRevisions orders two revisions of one stored event by when they were
// updated. When both update times are known and equal it breaks the tie on
// lastSeenAt and then occurrenceCount, which only grow from one revision to
// the next, so the order does not depend on which broadcast arrived first. It
// is negative when a is the earlier revision, and 0 when either update time is
// unknown or neither revision is known to be earlier.
const compareRevisions = (a, b) => {
  if (!instantOf(a.updatedAt) || !instantOf(b.updatedAt)) {
    return 0;
  }
  return (
    compareTimes(a.updatedAt, b.updatedAt) ||
    compareTimes(a.lastSeenAt, b.lastSeenAt) ||
    compareCounts(a.occurrenceCount, b.occurrenceCount)
  );
};

// isOlderRevision reports whether incoming is an earlier revision of the
// stored event than existing: it shows open a condition existing already
// shows resolved (a stored condition is never reopened), or compareRevisions
// orders it before existing.
const isOlderRevision = (incoming, existing) => {
  if (existing.resolvedAt && !incoming.resolvedAt) {
    return true;
  }
  return compareRevisions(incoming, existing) < 0;
};

export const eventsSlice = createSlice({
  name: 'events',
  initialState: eventsEntityAdapter.getInitialState(initialState),
  reducers: {
    clearEvents: (state) => {
      state.events = [];
    },

    setEvents: (state, action) => {
      // state.events = action.payload || []
      eventsEntityAdapter.removeAll(state);
      eventsEntityAdapter.addMany(state, action.payload);

      state.current_view.has_more = action.payload.length == 0 ? false : true;
      state.current_view.load_error = false;
      state.current_view.failed_request = null;
    },

    pushEvents: (state, action) => {
      // state.events = [...state.events, ...action.payload]
      eventsEntityAdapter.addMany(state, action.payload);
      state.current_view.has_more = action.payload.length == 0 ? false : true;
      state.current_view.load_error = false;
      state.current_view.failed_request = null;
    },

    // A page request failed: flag the error and record the { page, filters } that
    // failed so retry repeats that request, not the last page that loaded.
    setLoadError: (state, action) => {
      state.current_view.load_error = true;
      state.current_view.failed_request = action.payload || null;
    },

    pushEvent: (state, action) => {
      const event = {
        ...action.payload,
        severity: action.payload?.severity?.trim() || defaultEventProperties.severity,
        status: action.payload?.status?.trim() || defaultEventProperties.status,
      };
      // The server broadcasts an occurrence counted on, or the resolution of,
      // an open condition as the condition's stored event, so an event already
      // shown is updated in place rather than ignored.
      const existing = state.entities[event.id];
      if (!existing) {
        eventsEntityAdapter.addOne(state, event);
        return;
      }
      // Broadcasts of one condition can arrive out of order, so an older
      // revision of the stored event never replaces a newer one.
      if (isOlderRevision(event, existing)) {
        return;
      }
      // Whether the user has read the notification is the user's, not the
      // broadcast's: a re-push keeps it, and the selection, as they were.
      eventsEntityAdapter.upsertOne(state, {
        ...event,
        status: existing.status,
        checked: existing.checked,
      });
    },

    updateEvent: eventsEntityAdapter.updateOne,
    updateEvents: eventsEntityAdapter.updateMany,
    removeEvent: (state, { payload }) => {
      eventsEntityAdapter.removeOne(state, payload.id);
    },
    removeEvents: (state, { payload }) => {
      eventsEntityAdapter.removeMany(state, payload.ids);
    },
    updateIsEventChecked: (state, { payload }) => {
      const { id, value } = payload;
      eventsEntityAdapter.updateOne(state, {
        id,
        changes: {
          checked: value,
        },
      });
    },

    updateCheckAllEvents: (state, { payload }) => {
      const updates = Object.keys(state.entities).map((id) => ({
        id,
        changes: {
          checked: payload,
        },
      }));
      eventsEntityAdapter.updateMany(state, updates);
    },

    clearCurrentView: (state) => {
      state.current_view = initialState.current_view;
      state.events = [];
    },

    setCurrentView: (state, action) => {
      console.log('Setting current view in events slice:', action);
      state.current_view = action.payload;
    },

    toggleNotificationCenter: (state, action) => {
      console.log('Toggling notification center state', action);
      state.isNotificationCenterOpen = !state.isNotificationCenterOpen;
    },

    closeNotificationCenter: (state) => {
      state.isNotificationCenterOpen = false;
      state.view_to_fetch_on_open = initialState.view_to_fetch_on_open;
      state.current_view = initialState.current_view;
      state.ui = initialState.ui;
    },

    openNotificationCenter: (state, action) => {
      console.log('Opening notification center with view:', action.payload);
      state.isNotificationCenterOpen = true;
      const viewToOpen = isSerializable(action.payload.view_to_open)
        ? action.payload.view_to_open
        : undefined;
      state.view_to_fetch_on_open = viewToOpen || initialState.view_to_fetch_on_open;
      if (viewToOpen) {
        state.current_view = viewToOpen;
      }

      if (action.payload.ui) {
        state.ui = {
          ...state.ui,
          ...serializableUi(action.payload.ui),
        };
      }
    },
  },
});

// Action creators are generated for each case reducer function
export const {
  pushEvent,
  clearEvents,
  setEvents,
  clearCurrentView,
  updateIsEventChecked,
  updateCheckAllEvents,
  pushEvents,
  setLoadError,
  setCurrentView,
  updateEvent,
  toggleNotificationCenter,
  closeNotificationCenter,
  openNotificationCenter,
  updateEvents,
  removeEvent,
  removeEvents,
} = eventsSlice.actions;

export default eventsSlice.reducer;

// Monotonic token identifying the latest in-flight load. A filter change or a
// newer page load claims a fresh token; any earlier request whose result
// arrives late is superseded and its success/error is ignored, so a stale
// failure cannot overwrite the current view (e.g. re-flag an obsolete filter).
let latestRequestId = 0;

export const loadEvents = (fetch, page, filters) => async (dispatch, getState) => {
  const currentView = getState().events.current_view;
  const requestId = ++latestRequestId;
  const isSuperseded = () => requestId !== latestRequestId;
  try {
    // RTK Query's lazy trigger resolves with { error } on failure instead of
    // throwing, so a failed request must be detected here explicitly. Treating
    // it as an empty page would set has_more=false and advance the page,
    // making a transient failure look like the end of the list.
    const { data, error } = await fetch({ page, filters });
    if (isSuperseded()) {
      return;
    }
    if (error) {
      dispatch(setLoadError({ page, filters }));
      console.error('Error while fetching events in store --loadEvents', error);
      return;
    }
    dispatch(
      setCurrentView({
        ...currentView,
        page,
        filters,
        load_error: false,
        failed_request: null,
      }),
    );
    if (page <= 0) {
      dispatch(setEvents(data?.events || []));
      return;
    }
    dispatch(pushEvents(data?.events || []));
  } catch (e) {
    if (isSuperseded()) {
      return;
    }
    dispatch(setLoadError({ page, filters }));
    console.error('Error while setting events in store --loadEvents', e);
    return;
  }
};

export const loadNextPage = (fetch) => async (dispatch, getState) => {
  const currentView = getState().events.current_view;
  // Don't chase a page that failed to load as if it were the next one; let the
  // user retry the failed page instead (see retryCurrentPage).
  if (currentView.load_error) {
    return;
  }
  dispatch(loadEvents(fetch, currentView.page + 1, currentView.filters));
};

// Re-request the exact page + filters that failed (not the last page that
// loaded), so a retry does not skip a page or duplicate notifications.
export const retryCurrentPage = (fetch) => async (dispatch, getState) => {
  const currentView = getState().events.current_view;
  const failed = currentView.failed_request;
  if (!failed) {
    return;
  }
  dispatch(loadEvents(fetch, failed.page, failed.filters));
};

export const updateEventStatus =
  ({ id, status }) =>
  (dispatch) => {
    //const currentView = getState().events.current_view;
    dispatch(
      updateEvent({
        id,
        changes: {
          status,
        },
      }),
    );
  };

export const deleteEvent =
  ({ id }) =>
  (dispatch) => {
    dispatch(removeEvent({ id }));
  };

export const deleteEvents =
  ({ ids }) =>
  (dispatch) => {
    dispatch(removeEvents({ ids }));
  };

//selectors

//select all events
export const selectEvents = (state) => {
  return eventsEntityAdapter.getSelectors().selectAll(state.events);
};

// true when the last page request failed and can be retried
export const selectLoadError = (state) => {
  return Boolean(state.events.current_view.load_error);
};

export const selectCheckedEvents = (state) => {
  return selectEvents(state).filter((e) => e.checked);
};

export const selectEventById = (state, id) => {
  return eventsEntityAdapter.getSelectors().selectById(state.events, id);
};

export const selectIsEventChecked = (state, id) => {
  return Boolean(selectEventById(state, id).checked);
};

export const selectAreAllEventsChecked = (state) => {
  if (selectEvents(state).length == 0) {
    return false;
  }
  return selectEvents(state).reduce((selected, event) => (event.checked ? selected : false), true);
};

export const selectIsEventVisible = (state, id) => {
  const event = selectEventById(state, id);
  const currentFilters = state.events.current_view?.filters || {};
  const shouldBeInCurrentFilteredView = currentFilters.status
    ? currentFilters.status == event.status
    : true;
  const isDeleted = event.is_deleted || false;
  return !isDeleted && shouldBeInCurrentFilteredView;
};
export const selectSeverity = (state) => {
  const currentSeverityList = state.events?.current_view?.filters?.severity;
  return currentSeverityList ? currentSeverityList[0] : undefined;
};
