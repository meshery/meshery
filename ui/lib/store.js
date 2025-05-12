//@ts-check
import { createStore, applyMiddleware } from 'redux';
import { composeWithDevTools } from '@redux-devtools/extension';
import thunkMiddleware from 'redux-thunk';
import { fromJS } from 'immutable';
import { createContext } from 'react';
import { createDispatchHook, createSelectorHook } from 'react-redux';

const initialState = fromJS({
  loadTest: {
    testName: '',
    meshName: '',
    url: '',
    qps: 0,
    c: 0,
    t: '30s',
    result: {},
  },
  loadTestPref: {
    qps: 0,
    t: '30s',
    c: 0,
    gen: 'fortio',
    ts: new Date(),
  },
  meshAdapters: [],
  meshAdaptersts: new Date(),
  results: {
    startKey: '',
    results: [],
  },
  results_selection: {}, // format - { page: {index: content}}
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
  showProgress: false,
  selectedAdapter: '',
});

export const actionTypes = {
  UPDATE_LOAD_TEST_DATA: 'UPDATE_LOAD_TEST_DATA',
  UPDATE_ADAPTERS_INFO: 'UPDATE_ADAPTERS_INFO',
  UPDATE_RESULTS_SELECTION: 'UPDATE_RESULTS_SELECTION',
  CLEAR_RESULTS_SELECTION: 'CLEAR_RESULTS_SELECTION',
  UPDATE_GRAFANA_CONFIG: 'UPDATE_GRAFANA_CONFIG',
  UPDATE_PROMETHEUS_CONFIG: 'UPDATE_PROMETHEUS_CONFIG',
  UPDATE_STATIC_BOARD_CONFIG: 'UPDATE_STATIC_BOARD_CONFIG',
  UPDATE_LOAD_GEN_CONFIG: 'UPDATE_LOAD_GEN_CONFIG',
  UPDATE_PROGRESS: 'UPDATE_PROGRESS',
  SET_ADAPTER: 'SET_ADAPTER',
  UPDATE_EXTENSION_TYPE: 'UPDATE_EXTENSION_TYPE',
  SET_CONNECTION_METADATA: 'SET_CONNECTION_METADATA',
};

// REDUCERS
export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case actionTypes.UPDATE_USER:
      return state.mergeDeep({ user: action.user });
    case actionTypes.UPDATE_LOAD_TEST_DATA:
      return state.updateIn(['loadTest'], (val) => fromJS(action.loadTest));
    case actionTypes.UPDATE_LOAD_GEN_CONFIG:
      return state.mergeDeep({ loadTestPref: action.loadTestPref });
    case actionTypes.UPDATE_ADAPTERS_INFO:
      state = state.updateIn(['meshAdapters'], (val) => fromJS([]));
      state = state.updateIn(['meshAdaptersts'], (val) => fromJS(new Date()));
      return state.mergeDeep({ meshAdapters: action.meshAdapters });
    case actionTypes.UPDATE_RESULTS_SELECTION:
      if (Object.keys(action.results).length > 0) {
        return state.updateIn(['results_selection', action.page], (val) => action.results);
      } else {
        return state.deleteIn(['results_selection', action.page]);
      }
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
    case actionTypes.CLEAR_RESULTS_SELECTION:
      state = state.deleteIn(['results_selection']);
      return state.mergeDeep({ results_selection: fromJS({}) });

    case actionTypes.UPDATE_PROGRESS:
      return state.mergeDeep({ showProgress: action.showProgress });

    case actionTypes.SET_ADAPTER:
      return state.mergeDeep({ selectedAdapter: action.selectedAdapter });

    default:
      return state;
  }
};

// ACTION CREATOR

export const updateProgress =
  ({ showProgress }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_PROGRESS, showProgress });
  };

export const updateLoadTestData =
  ({ loadTest }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_LOAD_TEST_DATA, loadTest });
  };

export const updateLoadTestPref =
  ({ loadTestPref }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_LOAD_GEN_CONFIG, loadTestPref });
  };

export const updateAdaptersInfo =
  ({ meshAdapters }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_ADAPTERS_INFO, meshAdapters });
  };

export const updateResultsSelection =
  ({ page, results }) =>
  (dispatch) => {
    return dispatch({ type: actionTypes.UPDATE_RESULTS_SELECTION, page, results });
  };

export const clearResultsSelection = () => (dispatch) => {
  return dispatch({ type: actionTypes.CLEAR_RESULTS_SELECTION });
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
