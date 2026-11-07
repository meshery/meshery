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

const timeOf = (value) => {
  const at = typeof value === 'string' ? Date.parse(value) : NaN;
  return Number.isNaN(at) ? undefined : at;
};

// isOlderRevision reports whether incoming is an earlier revision of the
// stored event than existing: it was updated before existing was, or it shows
// open a condition existing already shows resolved (a stored condition is
// never reopened).
const isOlderRevision = (incoming, existing) => {
  if (existing.resolvedAt && !incoming.resolvedAt) {
    return true;
  }
  const incomingAt = timeOf(incoming.updatedAt);
  const existingAt = timeOf(existing.updatedAt);
  return incomingAt !== undefined && existingAt !== undefined && incomingAt < existingAt;
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
    },

    pushEvents: (state, action) => {
      // state.events = [...state.events, ...action.payload]
      eventsEntityAdapter.addMany(state, action.payload);
      state.current_view.has_more = action.payload.length == 0 ? false : true;
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
      state.view_to_fetch_on_open =
        action.payload.view_to_open || initialState.view_to_fetch_on_open;
      if (action.payload.view_to_open) {
        state.current_view = action.payload.view_to_open;
      }

      if (action.payload.ui) {
        state.ui = {
          ...state.ui,
          ...action.payload.ui,
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
  setCurrentView,
  updateEvent,
  toggleNotificationCenter,
  closeNotificationCenter,
  updateEvents,
} = eventsSlice.actions;

export default eventsSlice.reducer;

export const loadEvents = (fetch, page, filters) => async (dispatch, getState) => {
  const currentView = getState().events.current_view;
  try {
    const { data } = await fetch({ page, filters });
    dispatch(
      setCurrentView({
        ...currentView,
        page,
        filters,
      }),
    );
    if (page <= 0) {
      dispatch(setEvents(data?.events || []));
      return;
    }
    dispatch(pushEvents(data?.events || []));
  } catch (e) {
    console.error('Error while setting events in store --loadEvents', e);
    return;
  }
};

export const loadNextPage = (fetch) => async (dispatch, getState) => {
  const currentView = getState().events.current_view;
  dispatch(loadEvents(fetch, currentView.page + 1, currentView.filters));
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

// does a soft deletion on ui
export const deleteEvent =
  ({ id }) =>
  (dispatch) => {
    dispatch(updateEvent({ id, changes: { is_deleted: true } }));
    //mutator({ id });
  };

export const deleteEvents =
  ({ ids }) =>
  (dispatch) => {
    dispatch(
      updateEvents(
        ids.map((id) => ({
          id,
          changes: {
            is_deleted: true,
          },
        })),
      ),
    );
  };

//selectors

//select all events
export const selectEvents = (state) => {
  return eventsEntityAdapter.getSelectors().selectAll(state.events);
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
