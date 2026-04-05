import { createEntityAdapter, createSlice } from '@reduxjs/toolkit';
import { SEVERITY, STATUS } from '../../components/NotificationCenter/constants';
import { BellIcon } from '@sistent/sistent';

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
    history_mode: false, // used to determine if the notification center is in history mode . so we render in a different way
    title: 'Notifications', // title of the operation center
    empty_message: 'No notifications found', // message to show when there are no notifications
    icon: BellIcon,
  },
  isNotificationCenterOpen: false,
};

const defaultEventProperties = {
  severity: SEVERITY.INFO,
  status: STATUS.UNREAD,
};

const eventsEntityAdapter = createEntityAdapter({
  selectId: (event) => event.id,
  //sort based on update_at timestamp(utc)
  sortComparer: (a, b) => {
    if (b?.created_at?.localeCompare && a?.created_at?.localeCompare) {
      return b.created_at?.localeCompare(a.created_at);
    }
    return 0;
  },
});

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
      eventsEntityAdapter.addOne(state, event);
      // state.events = [event, ...state.events]
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
