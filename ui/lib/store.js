import { createStore, applyMiddleware } from 'redux'
import { composeWithDevTools } from '@redux-devtools/extension'
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

  notificationCenter : {
    openEventId : null  ,
    showFullNotificationCenter: false,
  },
  catalogVisibility: true,
  extensionType: '',
  capabilitiesRegistry: null,
  telemetryURLs: {
    grafana: [],
    prometheus: []
  },
  // global gql-subscriptions
  operatorState: null,
  meshSyncState: null,
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
  UPDATE_RESULTS_SELECTION : 'UPDATE_RESULTS_SELECTION',
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
  PUSH_EVENT : 'PUSH_EVENT',
  SET_CATALOG_CONTENT : 'SET_CATALOG_CONTENT',
  SET_OPERATOR_SUBSCRIPTION: 'SET_OPERATOR_SUBSCRIPTION',
  SET_MESHSYNC_SUBSCRIPTION: 'SET_MESHSYNC_SUBSCRIPTION',
  // UPDATE_SMI_RESULT: 'UPDATE_SMI_RESULT',
  UPDATE_EXTENSION_TYPE: 'UPDATE_EXTENSION_TYPE',
  UPDATE_CAPABILITY_REGISTRY: 'UPDATE_CAPABILITY_REGISTRY',
  UPDATE_TELEMETRY_URLS : 'UPDATE_TELEMETRY_URLS',

  OPEN_EVENT_IN_NOTIFICATION_CENTER : "OPEN_EVENT_IN_NOTIFICATION_CENTER",
  TOGGLE_NOTIFICATION_CENTER : "TOGGLE_NOTIFICATION_CENTER"
};

// REDUCERS
export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case actionTypes.UPDATE_PAGE:
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
      return state.mergeDeep({ user : action.user });
    case actionTypes.UPDATE_CLUSTER_CONFIG:
      return state.merge({ k8sConfig : action.k8sConfig });
    case actionTypes.SET_K8S_CONTEXT:
      return state.merge({ selectedK8sContexts : action.selectedK8sContexts });
    case actionTypes.UPDATE_LOAD_TEST_DATA:
      return state.updateIn(['loadTest'], val => fromJS(action.loadTest));
    case actionTypes.UPDATE_LOAD_GEN_CONFIG:
      return state.mergeDeep({ loadTestPref : action.loadTestPref });
    case actionTypes.UPDATE_ANONYMOUS_USAGE_STATS:
      return state.mergeDeep({ anonymousUsageStats : action.anonymousUsageStats });
    case actionTypes.UPDATE_ANONYMOUS_PERFORMANCE_RESULTS:
      return state.mergeDeep({ anonymousPerfResults : action.anonymousPerfResults });
    case actionTypes.UPDATE_ADAPTERS_INFO:
      state = state.updateIn(['meshAdapters'], val => fromJS([]));
      state = state.updateIn(['meshAdaptersts'], val => fromJS(new Date()));
      return state.mergeDeep({ meshAdapters : action.meshAdapters });
       case actionTypes.UPDATE_RESULTS_SELECTION:
           if (Object.keys(action.results).length > 0){
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

    case actionTypes.UPDATE_EVENTS:
      return state.merge({ events : action.events.sort((a,b)=> b.get('timestamp') -  a.get('timestamp') ) })

    case actionTypes.PUSH_EVENT :
      return state.merge({
        events : state.get("events").push(action.event).sort((a,b)=>b.get('timestamp') - a.get('timestamp'))
      })
    case actionTypes.SET_CATALOG_CONTENT:
      return state.mergeDeep({ catalogVisibility : action.catalogVisibility })

    case actionTypes.SET_OPERATOR_SUBSCRIPTION:
      return state.merge({operatorState: action.operatorState});

    case actionTypes.SET_MESHSYNC_SUBSCRIPTION:
      return state.merge({meshSyncState: action.meshSyncState});

    case actionTypes.UPDATE_EXTENSION_TYPE:
        return state.merge({ extensionType: action.extensionType });

    case actionTypes.UPDATE_CAPABILITY_REGISTRY:
      return state.merge({capabilitiesRegistry: action.capabilitiesRegistry})

    case actionTypes.UPDATE_TELEMETRY_URLS:
      return state.updateIn(['telemetryURLs'], val => fromJS(action.telemetryURLs));

    case actionTypes.OPEN_EVENT_IN_NOTIFICATION_CENTER:

      return state
        .setIn(['notificationCenter','showFullNotificationCenter'],true)
        .setIn(['notificationCenter','openEventId'],action.eventId )

    case actionTypes.TOGGLE_NOTIFICATION_CENTER :
       return state.
        setIn(['notificationCenter','showFullNotificationCenter'],!state.get('notificationCenter').get('showFullNotificationCenter'))

    default:
      return state;
  }
};

// ACTION CREATOR
export const updatepagepath = ({ path }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_PAGE, path });
}

export const updatepagetitle = ({ path, title }) => dispatch => {
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

export const setK8sContexts = ({ selectedK8sContexts }) => dispatch => {
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

export const updateResultsSelection = ({ page, results }) => dispatch => {
  return dispatch({ type : actionTypes.UPDATE_RESULTS_SELECTION, page, results });
}

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

// the localStorage key for events
const getEventsKey = (userId) => {
  // whenever the event schema changes bump the version
  const version = 1
  return `events-version-${version}-${userId}`
}

const persistEvents = (events,userId) => {
  try {
    localStorage.setItem(getEventsKey(userId),JSON.stringify(events))
  } catch (e) {
    console.error("Error Persisting events ",e)
  }
}

export const updateEvents = ({ events }) => (dispatch,getState) => {
  if (typeof events === 'object') {
    events = fromJS(events)
  }
  const result = dispatch({ type: actionTypes.UPDATE_EVENTS, events });
  const user = getState().get('user').toJS()
  const newEvents = getState().get('events').toJS()
  if (user && user.user_id) {
    persistEvents(newEvents,user.user_id)
  }
  return result
}


export const pushEvent = ({event}) => (dispatch,getState) => {
  if( typeof event == 'object') {
    event = fromJS(event)
  }
  const res = dispatch({
    type : actionTypes.PUSH_EVENT ,
    event
  })

  const user = getState().get('user').toJS()
  const events = getState().get('events').toJS()
  if (user && user.user_id) {
    persistEvents(events,user.user_id)
  }
  return res
}

export const loadEventsFromPersistence = () => (dispatch,getState) => {
  const user = getState().get('user').toJS()
  if(getState().get("events").length > 0 ||  !user.user_id)  {
    return
  }
  const rawData = localStorage.getItem(getEventsKey(user.user_id))
  let events ;
  try {
    events = JSON.parse(rawData) || []
  }catch(e) {
    console.error("error parsing events json",e,rawData)
    events = []
  }
  return dispatch({
      type : actionTypes.UPDATE_EVENTS,
      events : fromJS(events)
  })
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

export const updateExtensionType = ({ extensionType }) => dispatch => {
  return dispatch({type: actionTypes.UPDATE_EXTENSION_TYPE, extensionType})
}

export const updateCapabilities = ({capabilitiesRegistry}) => dispatch => {
  return dispatch({type: actionTypes.UPDATE_CAPABILITY_REGISTRY, capabilitiesRegistry})
}

export const updateTelemetryUrls = ({ telemetryURLs }) => dispatch => {
  return dispatch({type: actionTypes.UPDATE_TELEMETRY_URLS, telemetryURLs })
}

export const openEventInNotificationCenter = ({eventId}) => dispatch => {
  return dispatch({
    type : actionTypes.OPEN_EVENT_IN_NOTIFICATION_CENTER , eventId
  })
}

export const toggleNotificationCenter = () => dispatch => {
  return dispatch({
    type : actionTypes.TOGGLE_NOTIFICATION_CENTER
  })
}
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