import eventsReducer from './slices/events';
import globalEnvironmentContextReducer from './slices/globalEnvironmentContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';
import mesheryUiReducer from './slices/mesheryUi';
import prefTestReducer from './slices/prefTest';
import telemetryReducer from './slices/telemetry';
import adapterReducer from './slices/adapter';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalEnvironmentContext: globalEnvironmentContextReducer,
    ui: mesheryUiReducer,
    prefTest: prefTestReducer,
    telemetry: telemetryReducer,
    adapter: adapterReducer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
});
