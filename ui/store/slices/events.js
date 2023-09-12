import { createSlice } from '@reduxjs/toolkit'
import { SEVERITY, STATUS } from '../../components/NotificationCenter/constants'
import _ from 'lodash'

const initialState = {
  events: [],
  // summary: {
  //   count_by_severity_level: [],
  //   total_count: 0
  // },

  current_view: {
    page: 1,
    page_size: 10,
    filters: {
      initial: true,
    },
    has_more: true,
  }
}

const defaultEventProperties = {
  severity: SEVERITY.INFO,
  status: STATUS.UNREAD,
}

export const eventsSlice = createSlice({
  name: 'events',
  initialState,
  reducers: {

    clearEvents: (state) => {
      state.events = []
    },

    setEvents: (state, action) => {
      state.events = action.payload || []
      if (state.events.length == 0) {
        state.current_view.has_more = false
      }
    },


    pushEvents: (state, action) => {
      state.events = [...state.events, ...action.payload]
      if (action.payload.length == 0) {
        state.current_view.has_more = false
      }
    },

    pushEvent: (state, action) => {
      const event = {
        ...action.payload,
        severity: action.payload?.severity?.trim() || defaultEventProperties.severity,
        status: action.payload?.status?.trim() || defaultEventProperties.status,
      }
      state.events = [event, ...state.events]
    },

    clearCurrentView: (state) => {
      state.current_view = initialState.current_view
      state.events = []
    },

    setCurrentView: (state, action) => {
      state.current_view = action.payload
    }

  },
})

// Action creators are generated for each case reducer function
export const { pushEvent, clearEvents, setEvents,
  clearCurrentView,
  pushEvents, setCurrentView
} = eventsSlice.actions

export default eventsSlice.reducer


export const loadEvents = (fetch, page, filters) => async (dispatch, getState) => {
  console.log("loadEvents", page, filters,fetch)
  const currentView = getState().events.current_view
  if (currentView.page === page && _.isEqual(currentView.filters, filters)) {
    console.log("same page and filters")
    return
  }

  try {
    const { data } = await fetch({ page, filters })
    console.log("loadEventsRes", data)
    dispatch(setCurrentView({
      ...currentView,
      page,
      filters
    }))
    if (page <= 1) {
      dispatch(setEvents(data?.events))
      return
    }
    dispatch(pushEvents(data?.events || []))
  } catch (e) {
    console.error("Error while setting events in store --loadEvents", e)
    return
  }
}

export const loadNextPage = (fetch) => async (dispatch, getState) => {
  const currentView = getState().events.current_view
  dispatch(loadEvents(fetch,currentView.page + 1, currentView.filters))
}


