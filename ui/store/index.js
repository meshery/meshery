import eventsReducer from './slices/events'
import { configureStore } from '@reduxjs/toolkit'

export const store = configureStore({
  reducer: {
    events: eventsReducer,
  },
})