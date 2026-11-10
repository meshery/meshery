import eventsReducer from './slices/events';
import globalEnvironmentContextReducer from './slices/globalEnvironmentContext';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';
import mesheryUiReducer from './slices/mesheryUi';
import prefTestReducer from './slices/prefTest';
import adapterReducer from './slices/adapter';
import { rtkErrorMiddleware } from './middleware/rtkErrorMiddleware';
import { MESHERY_EXTENSION_EVENT } from '@sistent/sistent';
import { mesheryEventBus } from '@/utils/eventBus';

export const store = configureStore({
  reducer: {
    events: eventsReducer,
    globalEnvironmentContext: globalEnvironmentContextReducer,
    ui: mesheryUiReducer,
    prefTest: prefTestReducer,
    adapter: adapterReducer,
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware).concat(rtkErrorMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;

mesheryEventBus.on(MESHERY_EXTENSION_EVENT.DispatchToMesheryStore).subscribe((event) => {
  // `EventBus.on` filters at runtime but is typed as the full event union, so the
  // discriminant has to be re-checked here for `event.data` to narrow to a
  // dispatchable action rather than the union of every event payload.
  if (event.type !== MESHERY_EXTENSION_EVENT.DispatchToMesheryStore) return;
  console.log('Dispatching to Meshery Store:', event.data);
  store.dispatch(event.data);
});
