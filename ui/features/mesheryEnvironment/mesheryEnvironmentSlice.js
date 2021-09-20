import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { mesherySystemSync, promGrafMeshScan } from "./helpers";
import { fetchKuberernetesClusters, submitKubernetesClusterConfig } from "./kubernetesCluster";

/**
    this feature/module is responsible for all the data and functionalities regarding operations and
    configurations of Kubernetes clusters, connected Grafanas, Connected prometheus etc.
  */

/**
 * @typedef {{contextName: string, clusterName: string, currentContext: boolean | string}} Context
 */

/**
 * @typedef {{id:string, configuredServer: string, contextName: string, clusterConfigured: boolean, contexts: Array.<Context>, inClusterConfig: boolean}} KubernetesCluster
 */

/**
 * @typedef {Array.<KubernetesCluster>} KubernetesClusters
 */

/**
 * @typedef {{grafanaBoard: {}, grafanaPanels: {}, selectedTemplateVars: string[]}} SelectedGrafanaConfigType
 */

/**
 * @typedef {{grafanaURL :string, grafanaAPIKey :string, grafanaBoards: Array.<SelectedGrafanaConfigType>,}} GrafanaType
 */

/**
 * @typedef {Array.<GrafanaType>} connectedGrafanas
 */

/**
 * @typedef {{prometheusURL: string, selectedPromtheusBoardsConfig: Array.<SelectedGrafanaConfigType>}} PrometheusType
 */

/**
 * @typedef {Array.<PrometheusType>} connectedPrometheus
 */

const initialState = {
  /** @type {KubernetesClusters} */
  kubernetesClusters: [],

  /** @type {connectedGrafanas} */
  connectedGrafanas: [],

  /** @type {connectedPrometheus} */
  connectedPrometheus: [],
  // and other relavant data
};

export const mesherySystemSyncThunk = createAsyncThunk("mesheryEnvironment/mesherySystemSync", async () => {
  const response = await mesherySystemSync();
  console.log(response);
  return {
    // in the future, we might want to change it based on how meshery server sends
    // multiple clusters and metrics components to the frontend
    kubernetesClusters: [response?.k8sConfig],
    connectedGrafanas: [response?.grafana],
    connectedPrometheus: [response?.promethues],
  };
});

export const fetchKubernetesClustersThunk = createAsyncThunk("mesheryEnvironment/fetchKubernetesClusters", async () => {
  const response = await fetchKuberernetesClusters();
  return response;
});

export const submitKubernetesClusterConfigThunk = createAsyncThunk(
  "mesheryEnvironment/updateKubernetesCluster",
  async (payload) => {
    const response = await submitKubernetesClusterConfig(payload);
    // handle the response,
    // ask whether the operator has to be restarted or stopped or those kind of things
    // send a notification toast regarding the success or failure of the action performed
    return response;
  }
);

export const submitGrafanaConfigThunk = createAsyncThunk("mesheryEnvironment/submitGrafanaConfig", async () => {});

export const submitPrometheusConfigThunk = createAsyncThunk(
  "mesheryEnvironment/submitPrometheusConfig",
  async () => {}
);

export const promGrafMeshScanThunk = createAsyncThunk("mesheryEnvironment/promGrafMeshScan", async () => {
  const response = await promGrafMeshScan();
  return response;
});

const mesheryEnvironmentSlice = createSlice({
  name: "mesheryEnvironment",
  initialState,
  reducers: {
    // reducers to update the state in this slice
    updateKubernetesClusterData: (state, action) => {
      const cluster = state.kubernetesClusters.find((k8) => k8.id === action.payload.id);
      cluster.clusterConfigured = action.payload?.clusterConfigured;
      cluster.configuredServer = action.payload?.configuredServer;
      cluster.contextName = action.payload?.contextName;
      cluster.contexts = action.payload?.contexts;
    },
    addKubernetesCluster: (state, action) => {
      state.kubernetesClusters.push(action.payload.kubernetesCluster);
    },
    removeKubernetesCluster: (state, action) => {
      const index = state.kubernetesClusters.findIndex((cluster) => cluster.id === action.payload.id);
      state.kubernetesClusters.splice(index, 1);
      return state;
    },
    addGrafana: (state, action) => {
      state.connectedGrafanas.push(action.payload);
      return state;
    },
    addPrometheus: (state, action) => {
      state.connectedPrometheus.push(action.payload);
      return state;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchKubernetesClustersThunk.fulfilled, (state, action) => {
      const cluster = state.kubernetesClusters.find((k8) => k8.id === action.payload.id);
      cluster.clusterConfigured = action.payload.clusterConfigured;
      cluster.configuredServer = action.payload.configuredServer;
      cluster.contextName = action.payload.contextName;
      cluster.contexts = action.payload.contexts;
    });

    builder.addCase(mesherySystemSyncThunk.fulfilled, (state, action) => {
      state.loading = false;
      state.connectedGrafanas = action.payload.connectedGrafanas;
      state.connectedPrometheus = action.payload.connectedPrometheus;
      state.kubernetesClusters = action.payload.kubernetesClusters;
      return state;
    });
    builder.addCase(mesherySystemSyncThunk.pending, (state) => {
      state.loading = false;
      return state;
    });
    builder.addCase(mesherySystemSyncThunk.rejected, (state) => {
      state.loading = false;
      return state;
    });
  },
});

export default mesheryEnvironmentSlice.reducer;
export const { updateKubernetesClusterData, addKubernetesCluster, removeKubernetesCluster, addGrafana, addPrometheus } =
  mesheryEnvironmentSlice.actions;

// selectors should be written and exported
