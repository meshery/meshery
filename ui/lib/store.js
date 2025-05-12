//@ts-check
import { createStore, applyMiddleware } from 'redux';
import { composeWithDevTools } from '@redux-devtools/extension';
import thunkMiddleware from 'redux-thunk';
import { fromJS } from 'immutable';
import { createContext } from 'react';
import { createDispatchHook, createSelectorHook } from 'react-redux';

const initialState = fromJS({
  meshAdapters: [],
  meshAdaptersts: new Date(),
  grafana: {
    grafanaURL: '',
    grafanaAPIKey: '',
    grafanaBoardSearch: '',
    grafanaBoards: [],
    selectedBoardsConfigs: [],
    ts: new Date(-8640000000000000),
  },
  prometheus: {
    prometheusURL: '',
    selectedPrometheusBoardsConfigs: [],
    ts: new Date(-8640000000000000),
  },
  staticPrometheusBoardConfig: {},
  selectedAdapter: '',
});

export const actionTypes = {
  UPDATE_ADAPTERS_INFO: 'UPDATE_ADAPTERS_INFO',
  UPDATE_GRAFANA_CONFIG: 'UPDATE_GRAFANA_CONFIG',
  UPDATE_PROMETHEUS_CONFIG: 'UPDATE_PROMETHEUS_CONFIG',
  UPDATE_STATIC_BOARD_CONFIG: 'UPDATE_STATIC_BOARD_CONFIG',
  SET_ADAPTER: 'SET_ADAPTER',
  SET_CONNECTION_METADATA: 'SET_CONNECTION_METADATA',
};

// REDUCERS
export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case actionTypes.UPDATE_ADAPTERS_INFO:
      state = state.updateIn(['meshAdapters'], (val) => fromJS([]));
      state = state.updateIn(['meshAdaptersts'], (val) => fromJS(new Date()));
      return state.mergeDeep({ meshAdapters: action.meshAdapters });

    case actionTypes.UPDATE_GRAFANA_CONFIG:
      action.grafana.ts = new Date();
      return state.updateIn(['grafana'], (val) => fromJS(action.grafana));

    case actionTypes.UPDATE_PROMETHEUS_CONFIG:
      action.prometheus.ts = new Date();
      return state.updateIn(['prometheus'], (val) => fromJS(action.prometheus));

    case actionTypes.UPDATE_STATIC_BOARD_CONFIG:
      return state.updateIn(['staticPrometheusBoardConfig'], (val) =>
        fromJS(action.staticPrometheusBoardConfig),
      );

    case actionTypes.SET_ADAPTER:
      return state.mergeDeep({ selectedAdapter: action.selectedAdapter });

    default:
      return state;
  }
};

// ACTION CREATOR

export const updateAdaptersInfo =
  ({ meshAdapters }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_ADAPTERS_INFO, meshAdapters });
  };

export const updateGrafanaConfig =
  ({ grafana }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_GRAFANA_CONFIG, grafana });
  };

export const updatePrometheusConfig =
  ({ prometheus }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus });
  };

export const updateStaticPrometheusBoardConfig =
  ({ staticPrometheusBoardConfig }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_STATIC_BOARD_CONFIG, staticPrometheusBoardConfig });
  };

export const setAdapter =
  ({ selectedAdapter }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.SET_ADAPTER, selectedAdapter });
  };

export const makeStore = (initialState, options) => {
  return createStore(reducer, initialState, composeWithDevTools(applyMiddleware(thunkMiddleware)));
};

export const resultsMerge = (arr1, arr2) => {
  const keys = {};
  var arr = [];
  const compareAndAdd = (a) => {
    if (typeof keys[a.meshery_id] === 'undefined') {
      keys[a.meshery_id] = true;
      arr = arr.push(a);
    }
  };
  arr1.map(compareAndAdd);
  arr2.map(compareAndAdd);
  return arr;
};

export const LegacyStoreContext = createContext(null);

export const useLegacySelector = createSelectorHook(LegacyStoreContext);
export const useLegacyDispatch = createDispatchHook(LegacyStoreContext);
