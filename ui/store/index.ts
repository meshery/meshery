import eventsReducer from './slices/events';
import globalEnvironmentContextReducer from './slices/globalEnvironmentContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';
import mesheryUiReducer from './slices/mesheryUi';
import prefTestReducer from './slices/prefTest';
import telemetryReducer from './slices/telemetry';
import adapterReducer from './slices/adapter';
import sessionReducer from './slices/session';
import { authMiddleware } from './middleware/authMiddleware';
import { rtkErrorMiddleware } from './middleware/rtkErrorMiddleware';
import { mesheryEventBus } from '@/utils/eventBus';
import { installSessionInterceptor } from '../lib/sessionInterceptor';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalEnvironmentContext: globalEnvironmentContextReducer,
    ui: mesheryUiReducer,
    prefTest: prefTestReducer,
    telemetry: telemetryReducer,
    adapter: adapterReducer,
    session: sessionReducer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware).concat(authMiddleware).concat(rtkErrorMiddleware),
});

// Install global fetch interceptor for 401/redirect -> SESSION_EXPIRED.
// This catches ALL HTTP requests (dataFetch, RTK Query, Relay, raw fetch).
installSessionInterceptor(store);

mesheryEventBus.on('DISPATCH_TO_MESHERY_STORE').subscribe((event) => {
  console.log('Dispatching to Meshery Store:', event.data);
  store.dispatch(event.data);
});
