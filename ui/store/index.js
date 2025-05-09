import eventsReducer from './slices/events';
import globalEnvironmentContextReducer from './slices/globalEnvironmentContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';
import coreRedcer from './slices/core';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalEnvironmentContext: globalEnvironmentContextReducer,
    core: coreRedcer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
});
