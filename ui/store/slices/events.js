import { createEntityAdapter, createSlice } from '@reduxjs/toolkit';
import { SEVERITY, STATUS, validateEvents } from '../../components/NotificationCenter/constants';

const initialState = {
  current_view: {
    page: 1,
    page_size: 10,
    filters: {
      initial: true,
    },
    has_more: true,
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
      console.log('updates', updates);
      eventsEntityAdapter.updateMany(state, updates);
    },

    clearCurrentView: (state) => {
      state.current_view = initialState.current_view;
      state.events = [];
    },

    setCurrentView: (state, action) => {
      state.current_view = action.payload;
    },

    toggleNotificationCenter: (state) => {
      state.isNotificationCenterOpen = !state.isNotificationCenterOpen;
    },

    closeNotificationCenter: (state) => {
      state.isNotificationCenterOpen = false;
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
    if (page <= 1) {
      dispatch(setEvents(validateEvents(data?.events || [])));
      return;
    }
    dispatch(pushEvents(validateEvents(data?.events || [])));
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
  (dispatch, getState) => {
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

export const selectIsEventChecked = (state) => {
  return Boolean(selectEventById(state, id).checked);
};

export const selectAreAllEventsChecked = (state) => {
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
