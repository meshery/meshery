import eventsReducer from './slices/events';
import globalEnvironmentContextReducer from './slices/globalEnvironmentContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';
import mesheryUiReducer from './slices/mesheryUi';
import prefTestReducer from './slices/prefTest';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalEnvironmentContext: globalEnvironmentContextReducer,
    ui: mesheryUiReducer,
    prefTest: prefTestReducer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
});
