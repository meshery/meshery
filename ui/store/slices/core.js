import { createSlice } from '@reduxjs/toolkit';
import { createContext } from 'react';
import { createDispatchHook, createSelectorHook } from 'react-redux';
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
  anonymousUsageStats: true,
  anonymousPerfResults: true,
  showProgress: false,
  isDrawerCollapsed: false,
  selectedAdapter: '',
  events: [],
  notificationCenter: {
    openEventId: null,
    showFullNotificationCenter: false,
  },
  catalogVisibility: true,
  extensionType: '',
  capabilitiesRegistry: null,
  telemetryURLs: {
    grafana: [],
    prometheus: [],
  },
  // global gql-subscriptions
  operatorState: null,
  controllerState: null,
  meshSyncState: null,
  connectionMetadataState: null,
  organization: null,
  workspace: null,
  keys: null,
};

const coreSlice = createSlice({
  name: 'core',
  initialState,
  reducers: {
    updatePage: (state, action) => {
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
    updateLoadTestData: (state, action) => {
      state.loadTest = action.payload.loadTest;
    },
    updateLoadTestPref: (state, action) => {
      state.loadTestPref = { ...state.loadTestPref, ...action.payload.loadTestPref };
    },
    updateAnonymousUsageStats: (state, action) => {
      state.anonymousUsageStats = action.payload.anonymousUsageStats;
    },
    updateAnonymousPerformanceResults: (state, action) => {
      state.anonymousPerfResults = action.payload.anonymousPerfResults;
    },
    updateAdaptersInfo: (state, action) => {
      state.meshAdapters = action.payload.meshAdapters;
      state.meshAdaptersts = new Date();
    },
    updateResultsSelection: (state, action) => {
      const { page, results } = action.payload;
      if (Object.keys(results).length > 0) {
        state.results_selection = {
          ...state.results_selection,
          [page]: results,
        };
      } else {
        delete state.results_selection[page];
      }
    },
    clearResultsSelection: (state) => {
      state.results_selection = {};
    },
    updateGrafanaConfig: (state, action) => {
      state.grafana = {
        ...action.payload.grafana,
        ts: new Date(),
      };
    },
    updatePrometheusConfig: (state, action) => {
      state.prometheus = {
        ...action.payload.prometheus,
        ts: new Date(),
      };
    },
    updateStaticPrometheusBoardConfig: (state, action) => {
      state.staticPrometheusBoardConfig = action.payload.staticPrometheusBoardConfig;
    },
    updateProgress: (state, action) => {
      state.showProgress = action.payload.showProgress;
    },
    toggleDrawer: (state, action) => {
      state.isDrawerCollapsed = action.payload.isDrawerCollapsed;
    },
    setAdapter: (state, action) => {
      state.selectedAdapter = action.payload.selectedAdapter;
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
    updateTelemetryUrls: (state, action) => {
      state.telemetryURLs = action.payload.telemetryURLs;
    },
    openEventInNotificationCenter: (state, action) => {
      state.notificationCenter.openEventId = action.payload.eventId;
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
  updateTitle,
  updateBetaBadge,
  updateUser,
  updateK8SConfig,
  setK8sContexts: setK8sContextsAction,
  updateLoadTestData,
  updateLoadTestPref,
  updateAnonymousUsageStats,
  updateAnonymousPerformanceResults,
  updateAdaptersInfo,
  updateResultsSelection,
  clearResultsSelection,
  updateGrafanaConfig,
  updatePrometheusConfig,
  updateStaticPrometheusBoardConfig,
  updateProgress,
  toggleDrawer,
  setAdapter,
  toggleCatalogContent,
  setOperatorSubscription,
  setControllerState,
  setMeshsyncSubscription,
  updateExtensionType,
  updateCapabilities,
  updateTelemetryUrls,
  openEventInNotificationCenter,
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

export const LegacyStoreContext = createContext(null);
export const useLegacySelector = createSelectorHook(LegacyStoreContext);
export const useLegacyDispatch = createDispatchHook(LegacyStoreContext);

export default coreSlice.reducer;
