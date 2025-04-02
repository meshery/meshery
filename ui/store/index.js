import eventsReducer from './slices/events';
import globalEnvironmentContextReducer from './slices/globalEnvironmentContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalEnvironmentContext: globalEnvironmentContextReducer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
});
