import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { getK8sClusterIdsFromCtxId } from '@/utils/multi-ctx';
import { mesheryEventBus } from '@/utils/eventBus';
import { store } from '..';
// TODO(schemas-canonical): import type { ControllerStatus } from '@meshery/schemas/constructs/ControllerStatus';
// Replace the local type below once meshery/schemas >= 1.3.0 is published and
// the ControllerStatus construct is generated there.
// Tracked in: meshery/meshery#19424
import type { ControllerStatus } from '@/rtk-query/kubernetes';

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
    // Merges incoming controller-status items into the store by the composite
    // key (connectionID, controller).  This avoids wholesale replacement so
    // that a partial SSE diff (e.g. only the OPERATOR status changed) does not
    // clobber the last-known MESHSYNC / BROKER states.
    setControllerState: (
      state,
      action: PayloadAction<{ controllerState: ControllerStatus[] | ControllerStatus }>,
    ) => {
      const incoming = Array.isArray(action.payload.controllerState)
        ? action.payload.controllerState
        : [action.payload.controllerState];

      const existing: ControllerStatus[] = Array.isArray(state.controllerState)
        ? (state.controllerState as ControllerStatus[])
        : [];

      const merged = [...existing];
      for (const item of incoming) {
        const idx = merged.findIndex(
          (s) => s.connectionID === item.connectionID && s.controller === item.controller,
        );
        if (idx === -1) {
          merged.push(item);
        } else {
          merged[idx] = item;
        }
      }
      state.controllerState = merged;
    },

    clearControllerState: (state) => {
      state.controllerState = null;
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
  clearControllerState,
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

  mesheryEventBus.publish({
    type: 'K8S_CONTEXTS_UPDATED',
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
  return getK8sClusterIdsFromCtxId(state.ui?.selectedK8sContexts, state.ui?.k8sConfig);
};

export const selectK8sConfig = (state) => {
  return state.ui.k8sConfig;
};

export const selectedOrg = (state) => {
  return state.ui.organization;
};
export default coreSlice.reducer;
