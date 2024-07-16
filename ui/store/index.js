import eventsReducer from './slices/events';
import globalContextReducer from './slices/globalContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalContext: globalContextReducer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
});
