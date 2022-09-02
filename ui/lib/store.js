import { createStore, applyMiddleware } from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import thunkMiddleware from 'redux-thunk'
import { fromJS } from 'immutable'

const initialState = fromJS({
  page : {
    path : '',
    title : '',
    isBeta : false,
  },
  user : {},
  k8sConfig : [], // k8sconfig stores kubernetes cluster configs
  selectedK8sContexts : ["all"], // The selected k8s context on which the operations should be performed
  loadTest : {
    testName : '',
    meshName : '',
    url : '',
    qps : 0,
    c : 0,
    t : '30s',
    result : {},
  },
  loadTestPref : {
    qps : 0,
    t : '30s',
    c : 0,
    gen : 'fortio',
    ts : new Date(),
  },
  meshAdapters : [],
  meshAdaptersts : new Date(),
  results : {
    startKey : '',
    results : [],
  },
  results_selection : {}, // format - { page: {index: content}}
  grafana : {
    grafanaURL : '',
    grafanaAPIKey : '',
    grafanaBoardSearch : '',
    grafanaBoards : [],
    selectedBoardsConfigs : [],
    ts : new Date(-8640000000000000),
  },
  prometheus : {
    prometheusURL : '',
    selectedPrometheusBoardsConfigs : [],
    ts : new Date(-8640000000000000),
  },
  staticPrometheusBoardConfig : {},
  anonymousUsageStats : true,
  anonymousPerfResults : true,
  showProgress : false,
  isDrawerCollapsed: false,
  selectedAdapter : '',
  events:[],
  catalogVisibility: true,

  // global gql-subscriptions
  operatorState: null,
  meshSyncState: null
});

export const actionTypes = {
  UPDATE_PAGE : 'UPDATE_PAGE',
  UPDATE_TITLE : 'UPDATE_TITLE',
  UPDATE_USER : 'UPDATE_USER',
  UPDATE_BETA_BADGE : 'UPDATE_BETA_BADGE',
  UPDATE_CLUSTER_CONFIG : 'UPDATE_CLUSTER_CONFIG',
  SET_K8S_CONTEXT : 'SET_K8S_CONTEXT',
  UPDATE_LOAD_TEST_DATA : 'UPDATE_LOAD_TEST_DATA',
  UPDATE_ADAPTERS_INFO : 'UPDATE_ADAPTERS_INFO',
  // UPDATE_MESH_RESULTS: 'UPDATE_MESH_RESULTS',
  UPDATE_RESULTS_SELECTION : 'UPDATE_RESULTS_SELECTION',
  // DELETE_RESULTS_SELECTION: 'DELETE_RESULTS_SELECTION',
  CLEAR_RESULTS_SELECTION : 'CLEAR_RESULTS_SELECTION',
  UPDATE_GRAFANA_CONFIG : 'UPDATE_GRAFANA_CONFIG',
  UPDATE_PROMETHEUS_CONFIG : 'UPDATE_PROMETHEUS_CONFIG',
  UPDATE_STATIC_BOARD_CONFIG : 'UPDATE_STATIC_BOARD_CONFIG',
  UPDATE_LOAD_GEN_CONFIG : 'UPDATE_LOAD_GEN_CONFIG',
  UPDATE_ANONYMOUS_USAGE_STATS : 'UPDATE_ANONYMOUS_USAGE_STATS',
  UPDATE_ANONYMOUS_PERFORMANCE_RESULTS : 'UPDATE_ANONYMOUS_PERFORMANCE_RESULTS',
  UPDATE_PROGRESS : 'UPDATE_PROGRESS',
  TOOGLE_DRAWER : 'TOOGLE_DRAWER',
  SET_ADAPTER : 'SET_ADAPTER',
  UPDATE_EVENTS : 'UPDATE_EVENTS',
  SET_CATALOG_CONTENT : 'SET_CATALOG_CONTENT',
  SET_OPERATOR_SUBSCRIPTION: 'SET_OPERATOR_SUBSCRIPTION',
  SET_MESHSYNC_SUBSCRIPTION: 'SET_MESHSYNC_SUBSCRIPTION'
  // UPDATE_SMI_RESULT: 'UPDATE_SMI_RESULT',
};

// REDUCERS
export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case actionTypes.UPDATE_PAGE:
      // console.log(`received an action to update page: ${action.path} title: ${action.title}`);
      return state.mergeDeep({
        page : {
          path : action.path,
        }
      });
    case actionTypes.UPDATE_TITLE:
      return state.mergeDeep({
        page : {
          title : action.title,
        }
      });
    case actionTypes.UPDATE_BETA_BADGE:
      return state.mergeDeep({
        page : {
          isBeta : action.isBeta,
        }
      });
    case actionTypes.UPDATE_USER:
      // console.log(`received an action to update user: ${JSON.stringify(action.user)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.user }))}`);
      return state.mergeDeep({ user : action.user });
    case actionTypes.UPDATE_CLUSTER_CONFIG:
      // console.log(`received an action to update k8sconfig: ${JSON.stringify(action.k8sConfig)} and New state: ${JSON.stringify(state.mergeDeep({ k8sConfig: action.k8sConfig }))}`);
      return state.merge({ k8sConfig : action.k8sConfig });
    case actionTypes.SET_K8S_CONTEXT:
      return state.merge({ selectedK8sContexts : action.selectedK8sContexts });
    case actionTypes.UPDATE_LOAD_TEST_DATA:
      // console.log(`received an action to update k8sconfig: ${JSON.stringify(action.loadTest)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.loadTest }))}`);
      return state.updateIn(['loadTest'], val => fromJS(action.loadTest));
    case actionTypes.UPDATE_LOAD_GEN_CONFIG:
      return state.mergeDeep({ loadTestPref : action.loadTestPref });
    case actionTypes.UPDATE_ANONYMOUS_USAGE_STATS:
      return state.mergeDeep({ anonymousUsageStats : action.anonymousUsageStats });
    case actionTypes.UPDATE_ANONYMOUS_PERFORMANCE_RESULTS:
      return state.mergeDeep({ anonymousPerfResults : action.anonymousPerfResults });
    case actionTypes.UPDATE_ADAPTERS_INFO:
      // console.log(`received an action to update mesh info: ${JSON.stringify(action.mesh)} and New state: ${JSON.stringify(state.mergeDeep({ mesh: action.mesh }))}`);
      state = state.updateIn(['meshAdapters'], val => fromJS([]));
      state = state.updateIn(['meshAdaptersts'], val => fromJS(new Date()));
      return state.mergeDeep({ meshAdapters : action.meshAdapters });
    // case actionTypes.UPDATE_MESH_RESULTS:
    //   // console.log(`received an action to update mesh results: ${JSON.stringify(action.results)} and New state: ${JSON.stringify(state.mergeDeep({ results: action.results }))}`);
    //   // const results = state.get('results').get('results').toArray().concat(action.results);
    //   // do a more intelligent merge based on meshery_id
    //   const results = resultsMerge(state.get('results').get('results').toArray(), action.results);
    //   return state.mergeDeep({ results: { results }});
    case actionTypes.UPDATE_RESULTS_SELECTION:
      // let lg = `current page: ${action.page}, results_selection: ${JSON.stringify(Object.keys(action.results))}`;
      // Object.keys(action.results).forEach(pg =>{
      //   lg += `- indices: ${JSON.stringify(Object.keys(action.results[pg]))}`;
      // });
      // alert(lg);

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
      action.grafana.ts = new Date();
      return state.updateIn(['grafana'], val => fromJS(action.grafana));

    case actionTypes.UPDATE_PROMETHEUS_CONFIG:
      action.prometheus.ts = new Date();
      return state.updateIn(['prometheus'], val => fromJS(action.prometheus));

    case actionTypes.UPDATE_STATIC_BOARD_CONFIG:
      return state.updateIn(['staticPrometheusBoardConfig'], val => fromJS(action.staticPrometheusBoardConfig));
    case actionTypes.CLEAR_RESULTS_SELECTION:
      state = state.deleteIn(['results_selection']);
      return state.mergeDeep({ results_selection : fromJS({}) });

    case actionTypes.UPDATE_PROGRESS:
      return state.mergeDeep({ showProgress : action.showProgress });

    case actionTypes.TOOGLE_DRAWER:
      return state.mergeDeep({ isDrawerCollapsed : action.isDrawerCollapsed });

    case actionTypes.SET_ADAPTER:
      return state.mergeDeep({ selectedAdapter : action.selectedAdapter });   
      // case actionTypes.UPDATE_SMI_RESULT:
      //   console.log(`received an action to update smi result`,action.smi_result);
      //   if(action.smi_result!==undefined)
      //     return state.updateIn(['smi_result'], val => fromJS(action.smi_result));
      //   else
      //     return state

    case actionTypes.UPDATE_EVENTS:
      return state.merge({ events : action.events })

    case actionTypes.SET_CATALOG_CONTENT:
      return state.mergeDeep({ catalogVisibility : action.catalogVisibility })

    case actionTypes.SET_OPERATOR_SUBSCRIPTION: 
      return state.merge({operatorState: action.operatorState});

    case actionTypes.SET_MESHSYNC_SUBSCRIPTION: 
      return state.merge({meshSyncState: action.meshSyncState});

    default:
      return state;
  }
};

// ACTION CREATOR
export const updatepagepath = ({ path }) => dispatch => {
  // console.log("invoking the updatepagepathandtitle action creator. . .");
  return dispatch({ type : actionTypes.UPDATE_PAGE, path });
}

export const updatepagetitle = ({ path, title }) => dispatch => {
  // console.log("invoking the updatepagepathandtitle action creator. . .");
  return dispatch({ type : actionTypes.UPDATE_TITLE, title });
}

export const updateProgress = ({ showProgress }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_PROGRESS, showProgress });
}

export const updatebetabadge = ({ isBeta }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_BETA_BADGE, isBeta });
}

export const updateUser = ({ user }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_USER, user });
}

export const updateK8SConfig = ({ k8sConfig }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_CLUSTER_CONFIG, k8sConfig });
}

export const setK8sContexts = ({  selectedK8sContexts}) => dispatch => {
  return dispatch({ type : actionTypes.SET_K8S_CONTEXT, selectedK8sContexts });
}

export const updateLoadTestData = ({ loadTest }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_LOAD_TEST_DATA, loadTest });
}

export const updateLoadTestPref = ({ loadTestPref }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_LOAD_GEN_CONFIG, loadTestPref });
}
export const updateAnonymousUsageStats = ({ anonymousUsageStats }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_ANONYMOUS_USAGE_STATS, anonymousUsageStats });
}
export const updateAnonymousPerformanceResults = ({ anonymousPerfResults }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_ANONYMOUS_PERFORMANCE_RESULTS, anonymousPerfResults });
}

export const updateAdaptersInfo = ({ meshAdapters }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_ADAPTERS_INFO, meshAdapters });
}
// export const updateMeshResults = ({startKey, results}) => dispatch => {
//   return dispatch({ type: actionTypes.UPDATE_MESH_RESULTS, startKey, results });
// }
export const updateResultsSelection = ({ page, results }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_RESULTS_SELECTION, page, results });
}
// export const deleteFromResultsSelection = ({page, index}) => dispatch => {
//   return dispatch({ type: actionTypes.DELETE_RESULTS_SELECTION, page, index });
// }
export const clearResultsSelection = () => dispatch => {
  return dispatch({ type : actionTypes.CLEAR_RESULTS_SELECTION });
}
export const updateGrafanaConfig = ({ grafana }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_GRAFANA_CONFIG, grafana });
}

export const updatePrometheusConfig = ({ prometheus }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus });
}

export const updateStaticPrometheusBoardConfig = ({ staticPrometheusBoardConfig }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_STATIC_BOARD_CONFIG, staticPrometheusBoardConfig });
}

export const toggleDrawer = ({isDrawerCollapsed}) => dispatch => {
  return dispatch({ type : actionTypes.TOOGLE_DRAWER, isDrawerCollapsed });
}

export const setAdapter = ({selectedAdapter}) => dispatch => {
  return dispatch({ type : actionTypes.SET_ADAPTER, selectedAdapter });
}

export const updateEvents = ({ events }) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_EVENTS, events });
}

export const toggleCatalogContent = ({ catalogVisibility }) => dispatch => {
  return dispatch({ type: actionTypes.SET_CATALOG_CONTENT, catalogVisibility });
}

export const setOperatorSubscription = ({operatorState}) => dispatch => {
  return dispatch({type: actionTypes.SET_OPERATOR_SUBSCRIPTION, operatorState})
}

export const setMeshsyncSubscription = ({meshSyncState}) => dispatch => {
  return dispatch({type: actionTypes.SET_MESHSYNC_SUBSCRIPTION, meshSyncState})
}

// export const updateSMIResults = ({smi_result}) => dispatch => {
//   console.log("invoking the updateSMIResults action creator. . .",smi_result);
//   return dispatch({ type: actionTypes.UPDATE_SMI_RESULT, smi_result });
// }

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
  );
}

export const resultsMerge = (arr1, arr2) => {
  const keys = {};
  var arr = [];
  const compareAndAdd = (a) => {
    if (typeof keys[a.meshery_id] === 'undefined'){
      keys[a.meshery_id] = true;
      arr = arr.push(a);
    }
  };
  arr1.map(compareAndAdd);
  arr2.map(compareAndAdd);
  return arr;
};
