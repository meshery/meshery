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
    meshLocationURL: '', 
    reconfigureCluster: true,
  },
  loadTest: {
    url: '',
    qps: 0,
    c: 0,
    t: 1,
    result: {},
  },
  mesh: {
    Name: '',
    Ops: {},
  },
  results: {
    startKey: '',
    results: [],
  },
});

export const actionTypes = {
    UPDATE_PAGE: 'UPDATE_PAGE',
    UPDATE_USER: 'UPDATE_USER',
    UPDATE_CLUSTER_CONFIG: 'UPDATE_CLUSTER_CONFIG',
    UPDATE_LOAD_TEST_DATA: 'UPDATE_LOAD_TEST_DATA',
    UPDATE_MESH_INFO: 'UPDATE_MESH_INFO',
    UPDATE_MESH_RESULTS: 'UPDATE_MESH_RESULTS',
}

// REDUCERS
export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case actionTypes.UPDATE_PAGE:
      console.log(`received an action to update page: ${action.path} title: ${action.title}`);
      return state.mergeDeep({
          page: {
            title: action.title,
            path: action.path,
          }
      });
    case actionTypes.UPDATE_USER:
      console.log(`received an action to update user: ${JSON.stringify(action.user)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.user }))}`);
      return state.mergeDeep({ user: action.user });
    case actionTypes.UPDATE_CLUSTER_CONFIG:
      console.log(`received an action to update k8sconfig: ${JSON.stringify(action.k8sConfig)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.k8sConfig }))}`);
      return state.mergeDeep({ k8sConfig: action.k8sConfig });
    case actionTypes.UPDATE_LOAD_TEST_DATA:
      console.log(`received an action to update k8sconfig: ${JSON.stringify(action.loadTest)} and New state: ${JSON.stringify(state.mergeDeep({ user: action.loadTest }))}`);
      return state.mergeDeep({ loadTest: action.loadTest });
    case actionTypes.UPDATE_MESH_INFO:
      console.log(`received an action to update mesh info: ${JSON.stringify(action.mesh)} and New state: ${JSON.stringify(state.mergeDeep({ mesh: action.mesh }))}`);
      return state.mergeDeep({ mesh: action.mesh });
    case actionTypes.UPDATE_MESH_RESULTS:
      console.log(`received an action to update mesh results: ${JSON.stringify(action.results)} and New state: ${JSON.stringify(state.mergeDeep({ results: action.results }))}`);
      const results = state.get('results').get('results').toArray().concat(action.results);
      return state.mergeDeep({ results: { startKey: action.startKey, results }}); 
      
      
    
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

export const updateUser = ({user}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_USER, user });
}

export const updateK8SConfig = ({k8sConfig}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_CLUSTER_CONFIG, k8sConfig });
}

export const updateLoadTestData = ({loadTest}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_LOAD_TEST_DATA, loadTest });
}

export const updateMeshInfo = ({mesh}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_MESH_INFO, mesh });
}
export const updateMeshResults = ({startKey, results}) => dispatch => {
  return dispatch({ type: actionTypes.UPDATE_MESH_RESULTS, startKey, results });
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