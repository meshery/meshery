import eventsReducer from './slices/events';

import globalContextReducer, {
  ORG_ACTIONS_TO_PERSIST,
  WORKSPACE_ACTIONS_TO_PERSIST,
} from './slices/globalContext';

import { combineReducers, configureStore } from '@reduxjs/toolkit';
import { api } from '../rtk-query/index';

import { initReduxPersist } from './redux-persist';

export const ACTIONS_TO_PERSIST = {
  ...ORG_ACTIONS_TO_PERSIST,
  ...WORKSPACE_ACTIONS_TO_PERSIST,
};

const reduxPersist = initReduxPersist(ACTIONS_TO_PERSIST);
const rootReducer = combineReducers({
  events: eventsReducer,
  globalContext: globalContextReducer,
  [api.reducerPath]: api.reducer,
});
export const store = configureStore({
  reducer: reduxPersist.createPersistEnhancedReducer(rootReducer),
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware).concat(reduxPersist.persistMiddleware),
});

export const { loadPersistedState } = reduxPersist;
