import _ from 'lodash';
import { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { useDispatchRtk } from '../hooks';

const INVALID_REDUCER_PATH = Symbol('INVALID_REDUCER_PATH');
const REHYDRATE_STATE_ACTION = 'REHYDRATE_STATE_ACTION';

/**
 * Creates an action to rehydrate state.
 *
 * @param {string} reducerPath - The path of the reducer to rehydrate.
 * @param {*} inflatedState - The state to rehydrate with.
 * @returns {Object} An action object for rehydrating state.
 */
const rehydrateState = (reducerPath, inflatedState) => {
  return {
    type: REHYDRATE_STATE_ACTION,
    payload: {
      reducerPath,
      inflatedState,
    },
  };
};

/**
 * Reducer to handle state rehydration during Redux store updates.
 * This reducer intercepts a specific action type to rehydrate the state of specified reducers.
 *
 * @param {*} state - The current state.
 * @param {Object} action - The dispatched action.
 * @returns {*} The new state after rehydration.
 */
const rehydrateStateReducer = (state, action) => {
  if (action.type == REHYDRATE_STATE_ACTION) {
    const appState = _.cloneDeep(state);
    // Set the rehydrated state at the specified reducer path
    _.set(appState, action.payload.reducerPath.split('/'), action.payload.inflatedState);
    return appState;
  }
  return state;
};

/**
 * Initializes Redux persistence with given actions to persist on
 *
 * @param {Object} actionsToPersist - An object mapping action types to arrays of reducer paths.
 * Each action type is associated with an array of reducer paths whose state should be persisted.
 * @returns {Object} An object containing Redux persistence functions.
 */
export const initReduxPersist = (actionsToPersist) => {
  /**
   * Creates a new reducer with enhanced state rehydration logic for Redux persistence.
   * This function returns a new reducer that first rehydrates the state using the
   * rehydrateStateReducer, and then applies the original reducer to the rehydrated state.
   *
   * @param {function} reducer - The original reducer function to enhance.
   * @returns {function} A new enhanced reducer function with added state rehydration.
   */
  const createPersistEnhancedReducer = (reducer) => (state, action) => {
    // Rehydrate the state using the rehydrateStateReducer
    const newState = rehydrateStateReducer(state, action);
    // Apply the original reducer to the rehydrated state
    return reducer(newState, action);
  };

  /**
   * Redux middleware to persist state to local storage based on dispatched actions.
   * This middleware listens for specific actions and saves the relevant state to local storage.
   *
   * @param {Object} store - The Redux store.
   * @returns {Function} A middleware function.
   */
  const persistMiddleware = (store) => (next) => (action) => {
    const result = next(action);
    const reducersToPersist = actionsToPersist[action.type];
    if (reducersToPersist) {
      const appState = store.getState();
      reducersToPersist.forEach((reducerPath) => {
        const path = reducerPath.split('/');
        console.log('here', action, appState, path);
        const stateToPersist = _.get(appState, path, INVALID_REDUCER_PATH);

        if (stateToPersist == INVALID_REDUCER_PATH) {
          throw Error('Reducer Path to Persist Is Invalid', reducerPath);
        }

        localStorage.setItem(reducerPath, JSON.stringify(stateToPersist));
      });
    }
    return result;
  };

  /**
   * Action creator to load persisted state from local storage during Redux store initialization.
   * This function retrieves previously saved state from local storage and dispatches rehydration actions.
   *
   * @returns {Function} A thunk function.
   */
  const loadPersistedState = () => (dispatch) => {
    Object.values(actionsToPersist).forEach((reducerPaths) => {
      reducerPaths.forEach((path) => {
        let inflatedState = localStorage.getItem(path);
        try {
          inflatedState = JSON.parse(inflatedState);
        } catch (e) {
          console.error(`Error rehydrating state for reducer ${path}"`, inflatedState);
          return;
        }
        inflatedState && dispatch(rehydrateState(path, inflatedState));
      });
    });
  };

  return {
    persistMiddleware,
    createPersistEnhancedReducer,
    loadPersistedState,
  };
};

export const PersistedStateProvider = ({ children, loadPersistedState }) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const dispatch = useDispatchRtk();
  useEffect(() => {
    if (!loading) {
      return;
    }
    try {
      dispatch(loadPersistedState());
    } catch (e) {
      setError(e);
    }
    setLoading(false);

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [loading, dispatch]);

  error && console.error('Error Loading Persisted State', error);
  if (loading) {
    return null;
  }
  return <>{children}</>;
};
