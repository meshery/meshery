import eventsReducer from './slices/events'
import { configureStore } from '@reduxjs/toolkit'
import { api } from "../rtk-query/index"

export const store = configureStore({
  reducer : {
    events : eventsReducer,
    [api.reducerPath] : api.reducer,
  },
  middleware : (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
})