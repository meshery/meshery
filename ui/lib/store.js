import { createStore, applyMiddleware } from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import thunkMiddleware from 'redux-thunk'
import { fromJS } from 'immutable'

const initialState = fromJS({
  page: {
    path: '',
    title: '',
  },
  user: {},
  k8sConfig: {
    inClusterConfig: false,
    k8sfile: '', 
    contextName: '', 
    clusterConfigured: false,
    configuredServer: '',
  },
  loadTest: {
    testName: '',
    meshName: '',
    url: '',
    qps: 0,
    c: 0,
    t: '30s',
    result: {},
  },
  meshAdapters: [],
  results: {
    startKey: '',
    results: [],
  },
  results_selection: {}, // { page: {index: content}}
  grafana: {
    grafanaURL: '',
    grafanaAPIKey: '',
    grafanaBoardSearch: '',
    grafanaBoards: [],
    selectedBoardsConfigs: [],
  },

  showProgress: false,
});

export const actionTypes = {
    UPDATE_PAGE: 'UPDATE_PAGE',
    UPDATE_USER: 'UPDATE_USER',
    UPDATE_CLUSTER_CONFIG: 'UPDATE_CLUSTER_CONFIG',
    UPDATE_LOAD_TEST_DATA: 'UPDATE_LOAD_TEST_DATA',
    UPDATE_ADAPTERS_INFO: 'UPDATE_ADAPTERS_INFO',
    // UPDATE_MESH_RESULTS: 'UPDATE_MESH_RESULTS',
    UPDATE_RESULTS_SELECTION: 'UPDATE_RESULTS_SELECTION',
    // DELETE_RESULTS_SELECTION: 'DELETE_RESULTS_SELECTION',
    CLEAR_RESULTS_SELECTION: 'CLEAR_RESULTS_SELECTION',
    UPDATE_GRAFANA_CONFIG: 'UPDATE_GRAFANA_CONFIG',
    UPDATE_PROGRESS: 'UPDATE_PROGRESS',
}

// REDUCERS
export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case actionTypes.UPDATE_PAGE:
      // console.log(`received an action to update page: ${action.path} title: ${action.title}`);
      return state.mergeDeep({
          page: {
            title: action.title,
            path: action.path,
          }
      });
    case actionTypes.UPDATE_USER:
      // console.log(`received an action to update user: ${JSON.stringify(action.user)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.user }))}`);
      return state.mergeDeep({ user: action.user });
    case actionTypes.UPDATE_CLUSTER_CONFIG:
      // console.log(`received an action to update k8sconfig: ${JSON.stringify(action.k8sConfig)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.k8sConfig }))}`);
      return state.mergeDeep({ k8sConfig: action.k8sConfig });
    case actionTypes.UPDATE_LOAD_TEST_DATA:
      // console.log(`received an action to update k8sconfig: ${JSON.stringify(action.loadTest)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.loadTest }))}`);
      return state.mergeDeep({ loadTest: action.loadTest });
    case actionTypes.UPDATE_ADAPTERS_INFO:
      // console.log(`received an action to update mesh info: ${JSON.stringify(action.mesh)} and New state: ${JSON.stringify(state.mergeDeep({ mesh: action.mesh }))}`);
      state = state.updateIn(['meshAdapters'], val => fromJS([]));
      return state.mergeDeep({ meshAdapters: action.meshAdapters });
    // case actionTypes.UPDATE_MESH_RESULTS:
    //   // console.log(`received an action to update mesh results: ${JSON.stringify(action.results)} and New state: ${JSON.stringify(state.mergeDeep({ results: action.results }))}`);
    //   // const results = state.get('results').get('results').toArray().concat(action.results);
    //   // do a more intelligent merge based on meshery_id
    //   const results = resultsMerge(state.get('results').get('results').toArray(), action.results);
    //   return state.mergeDeep({ results: { results }}); 
    case actionTypes.UPDATE_RESULTS_SELECTION:
      // console.log(`current results_selection: ${JSON.stringify(rs)}`);
      // if (typeof rs[action.page] === 'undefined'){
      //   rs[action.page] = {};
      // }
      if (Object.keys(action.results).length > 0){
        // const rs = state.get('results_selection').toObject();
        // rs[action.page] = action.results;
        return state.updateIn(['results_selection', action.page], val => action.results);
      } else {
        return state.deleteIn(['results_selection', action.page]);
      }
    case actionTypes.UPDATE_GRAFANA_CONFIG:
      return state.updateIn(['grafana'], val => fromJS(action.grafana));
      
    // case actionTypes.DELETE_RESULTS_SELECTION:
    //   const rs = state.get('results_selection').toObject();
    //   if (typeof rs[action.page] === 'undefined'){
    //     rs[action.page] = {};
    //   }
    //   delete(rs[action.page][action.index]);
    //   return state.mergeDeep({results_selection: rs});
    case actionTypes.CLEAR_RESULTS_SELECTION:
      state = state.deleteIn(['results_selection']);
      return state.mergeDeep({results_selection: fromJS({})});

    case actionTypes.UPDATE_PROGRESS:
      return state.mergeDeep({showProgress: action.showProgress});
      // return state.mergeDeep({results_selection: {}});
      // console.log(`keys: ${state.get('results_selection').keys()}`);
      // return state.deleteIn(state.get('results_selection').keys());
    
    // case actionTypes.INCREMENT:
    //   return Object.assign({}, state, {
    //     count: state.count + 1
    //   })
    // case actionTypes.DECREMENT:
    //   return Object.assign({}, state, {
    //     count: state.count - 1
    //   })
    // case actionTypes.RESET:
    //   return Object.assign({}, state, {
    //     count: exampleInitialState.count
    //   })
    default:
      return state
  }
}

// ACTION CREATOR
export const updatepagepathandtitle = ({path, title}) => dispatch => {
    // console.log("invoking the updatepagepathandtitle action creator. . .");
  return dispatch({ type: actionTypes.UPDATE_PAGE, path, title });
}

export const updateProgress = ({showProgress}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_PROGRESS, showProgress });
}

export const updateUser = ({user}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_USER, user });
}

export const updateK8SConfig = ({k8sConfig}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_CLUSTER_CONFIG, k8sConfig });
}

export const updateLoadTestData = ({loadTest}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_LOAD_TEST_DATA, loadTest });
}

export const updateAdaptersInfo = ({meshAdapters}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_ADAPTERS_INFO, meshAdapters });
}
// export const updateMeshResults = ({startKey, results}) => dispatch => {
//   return dispatch({ type: actionTypes.UPDATE_MESH_RESULTS, startKey, results });
// }
export const updateResultsSelection = ({page, results}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_RESULTS_SELECTION, page, results });
}
// export const deleteFromResultsSelection = ({page, index}) => dispatch => {
//   return dispatch({ type: actionTypes.DELETE_RESULTS_SELECTION, page, index });
// }
export const clearResultsSelection = () => dispatch => {
  return dispatch({ type: actionTypes.CLEAR_RESULTS_SELECTION});
}
export const updateGrafanaConfig = ({grafana}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_GRAFANA_CONFIG, grafana });
}


// export const startClock = dispatch => {
//   return setInterval(() => {
//     // Dispatch `TICK` every 1 second
//     dispatch({ type: actionTypes.TICK, light: true, ts: Date.now() })
//   }, 1000)
// }

// export const incrementCount = () => dispatch => {
//   return dispatch({ type: actionTypes.INCREMENT })
// }

// export const decrementCount = () => dispatch => {
//   return dispatch({ type: actionTypes.DECREMENT })
// }

// export const resetCount = () => dispatch => {
//   return dispatch({ type: actionTypes.RESET })
// }

export const makeStore = (initialState, options) => {
  return createStore(
    reducer,
    initialState,
    composeWithDevTools(applyMiddleware(thunkMiddleware))
  )
}

export const resultsMerge = (arr1, arr2) => {
  const keys = {}
  const arr = [];
  const compareAndAdd = (a) => {
    if (typeof keys[a.meshery_id] === 'undefined'){
      keys[a.meshery_id] = true;
      arr = arr.push(a);
    }
  }
  arr1.map(compareAndAdd);
  arr2.map(compareAndAdd);
  return arr;
}