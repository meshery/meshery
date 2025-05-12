import { createSlice } from '@reduxjs/toolkit';
import { getK8sClusterIdsFromCtxId } from '@/utils/multi-ctx';
import { mesheryEventBus } from '@/utils/eventBus';

const initialState = {
  page: {
    path: '',
    title: '',
    isBeta: false,
  },
  user: {},
  k8sConfig: [], // k8sconfig stores kubernetes cluster configs
  selectedK8sContexts: ['all'], // The selected k8s context on which the operations should be performed
  showProgress: false,
  isDrawerCollapsed: false,
  catalogVisibility: true,
  extensionType: '',
  capabilitiesRegistry: null,
  operatorState: null,
  controllerState: null,
  meshSyncState: null,
  connectionMetadataState: null, // store connection definition metadata for state and connection kind management
  organization: null,
  workspace: null,
  keys: null,
};

const coreSlice = createSlice({
  name: 'core',
  initialState,
  reducers: {
    updatePagePath: (state, action) => {
      state.page.path = action.payload.path;
    },
    updateTitle: (state, action) => {
      state.page.title = action.payload.title;
    },
    updateBetaBadge: (state, action) => {
      state.page.isBeta = action.payload.isBeta;
    },
    updateUser: (state, action) => {
      state.user = action.payload.user;
    },
    updateK8SConfig: (state, action) => {
      state.k8sConfig = action.payload.k8sConfig;
    },
    setK8sContexts: (state, action) => {
      state.selectedK8sContexts = action.payload.selectedK8sContexts;
      // Note: Event bus publication would be handled in the thunk action
    },
    updateProgress: (state, action) => {
      state.showProgress = action.payload.showProgress;
    },
    toggleDrawer: (state, action) => {
      state.isDrawerCollapsed = action.payload.isDrawerCollapsed;
    },
    toggleCatalogContent: (state, action) => {
      state.catalogVisibility = action.payload.catalogVisibility;
    },
    setOperatorSubscription: (state, action) => {
      state.operatorState = action.payload.operatorState;
    },
    setControllerState: (state, action) => {
      state.controllerState = action.payload.controllerState;
    },
    setMeshsyncSubscription: (state, action) => {
      state.meshSyncState = action.payload.meshSyncState;
    },
    updateExtensionType: (state, action) => {
      state.extensionType = action.payload.extensionType;
    },
    updateCapabilities: (state, action) => {
      state.capabilitiesRegistry = action.payload.capabilitiesRegistry;
    },
    setConnectionMetadata: (state, action) => {
      state.connectionMetadataState = action.payload.connectionMetadataState;
    },
    setOrganization: (state, action) => {
      state.organization = action.payload.organization;
      sessionStorage.setItem('currentOrg', JSON.stringify(action.payload.organization));
    },
    setWorkspace: (state, action) => {
      state.workspace = action.payload.workspace;
      sessionStorage.setItem('currentWorkspace', JSON.stringify(action.payload.workspace));
    },
    setKeys: (state, action) => {
      state.keys = action.payload.keys;
      sessionStorage.setItem('keys', JSON.stringify(action.payload.keys));
    },
  },
});

// Extract the action creators and the reducer
export const {
  updatePagePath,
  updateTitle,
  updateBetaBadge,
  updateUser,
  updateK8SConfig,
  setK8sContexts: setK8sContextsAction,
  updateProgress,
  toggleDrawer,
  toggleCatalogContent,
  setOperatorSubscription,
  setControllerState,
  setMeshsyncSubscription,
  updateExtensionType,
  updateCapabilities,
  setConnectionMetadata,
  setOrganization,
  setWorkspace,
  setKeys,
} = coreSlice.actions;

// Add thunks for async operations or side effects
export const setK8sContexts = (payload) => (dispatch) => {
  dispatch(setK8sContextsAction(payload));

  mesheryEventBus.publish({
    type: 'K8S_CONTEXTS_UPDATED',
    data: {
      selectedK8sContexts: payload.selectedK8sContexts,
    },
  });
};

// Core middleware configuration
export const coreMiddleware = (getDefaultMiddleware) =>
  getDefaultMiddleware({
    serializableCheck: {
      // Ignore these action types
      ignoredActions: ['core/updateGrafanaConfig', 'core/updatePrometheusConfig'],
      // Ignore these field paths in all actions
      ignoredActionPaths: ['payload.grafana.ts', 'payload.prometheus.ts'],
      // Ignore these paths in the state
      ignoredPaths: [
        'core.grafana.ts',
        'core.prometheus.ts',
        'core.loadTestPref.ts',
        'core.meshAdaptersts',
      ],
    },
  });

// Selectors
export const selectSelectedK8sClusters = (state) => {
  return getK8sClusterIdsFromCtxId(state.selectedK8sContexts, state.k8sConfig);
};

export const selectK8sConfig = (state) => {
  return state.k8sConfig;
};

export default coreSlice.reducer;
