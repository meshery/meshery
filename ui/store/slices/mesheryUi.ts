import { createSlice } from '@reduxjs/toolkit';
import { MESHERY_EXTENSION_EVENT } from '@sistent/sistent';
import { getK8sClusterIdsFromCtxId, persistSelectedK8sContexts } from '@/utils/multi-ctx';
import { mesheryEventBus } from '@/utils/eventBus';
import { store } from '..';

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
  providerCapabilities: null,
  controllerState: null,
  connectionMetadataState: null, // store connection definition metadata for state and connection kind management
  organization: null,
  workspace: null,
  keys: null,
};

const coreSlice = createSlice({
  name: 'core',
  initialState,
  reducers: {
    updatePage: (state, action) => {
      state.page = {
        ...state.page,
        ...action.payload,
      };
    },
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
      // Note: Side effects (session persistence, event bus publication) are
      // handled in the setK8sContexts thunk below; this reducer is side-effect free.
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
    setControllerState: (state, action) => {
      state.controllerState = action.payload.controllerState;
    },
    updateExtensionType: (state, action) => {
      state.extensionType = action.payload.extensionType;
    },
    updateProviderCapabilities: (state, action) => {
      state.providerCapabilities = action.payload.providerCapabilities;
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
  updatePage,
  updatePagePath,
  updateTitle,
  updateBetaBadge,
  updateUser,
  updateK8SConfig,
  setK8sContexts: setK8sContextsAction,
  updateProgress: updateProgressAction,
  toggleDrawer,
  toggleCatalogContent,
  setControllerState,
  setMeshsyncSubscription,
  updateExtensionType,
  updateProviderCapabilities,
  setConnectionMetadata,
  setOrganization,
  setWorkspace,
  setKeys,
} = coreSlice.actions;

// Add thunks for async operations or side effects
export const setK8sContexts = (payload) => (dispatch) => {
  dispatch(setK8sContextsAction(payload));

  // Session-persist the selection so it survives navigation and reloads.
  // Every selection change in the app flows through this thunk (header
  // checkboxes, deploy modal, context search), making it the single
  // persistence funnel while keeping the reducer pure.
  persistSelectedK8sContexts(payload.selectedK8sContexts);

  mesheryEventBus.publish({
    type: MESHERY_EXTENSION_EVENT.K8sContextsUpdated,
    data: {
      selectedK8sContexts: payload.selectedK8sContexts,
    },
  });
};

export const updateProgress = (progressData) => {
  store.dispatch(updateProgressAction(progressData));
};
// Core middleware configuration
export const coreMiddleware = (getDefaultMiddleware) =>
  getDefaultMiddleware({
    serializableCheck: {
      // Ignore these paths in the state
      ignoredPaths: ['core.loadTestPref.ts', 'core.meshAdaptersts'],
    },
  });

// Selectors
export const selectSelectedK8sClusters = (state) => {
  return getK8sClusterIdsFromCtxId(state.ui?.selectedK8sContexts, state.ui?.k8sConfig);
};

export const selectK8sConfig = (state) => {
  return state.ui.k8sConfig;
};

export const selectedOrg = (state) => {
  return state.ui.organization;
};
export default coreSlice.reducer;
