import { createSlice } from '@reduxjs/toolkit'
import { SEVERITY } from '../../components/NotificationCenter/constants'

const initialState = {
  events: [],
  summary: {
    count_by_severity_level: [],
    total_count: 0
  }
}

const defaultEventProperties = {
  severity: SEVERITY.INFO,
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
    },

    setEventsSummary: (state, action) => {
      state.summary = action.payload || []
    },

    pushEvent: (state, action) => {
      const event = {
        ...action.payload,
        severity: action.payload?.severity?.trim() || defaultEventProperties.severity,
      }
      state.events = [event, ...state.events]
      //update severity count
      const severity = event.severity
      const severityIndex = state.summary.count_by_severity_level.findIndex((item) => item.severity === severity)
      if (severityIndex === -1) {
        state.summary.count_by_severity_level.push({
          severity: severity,
          count: 1
        })
      } else {
        state.summary.count_by_severity_level[severityIndex].count += 1
      }
      state.summary.total_count += 1
    },

  },
})

// Action creators are generated for each case reducer function
export const { pushEvent, clearEvents, setEvents, setEventsSummary } = eventsSlice.actions

export default eventsSlice.reducer