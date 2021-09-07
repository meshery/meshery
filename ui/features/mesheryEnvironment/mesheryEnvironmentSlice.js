import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { fetchKuberernetesClusters, submitKubernetesClusterConfig } from "./kubernetesCluster";

/**
    this feature/module is responsible for all the data and functionalities regarding operations and
    configurations of Kubernetes clusters, connected Grafanas, Connected prometheus etc.
  */

/**
 * @typedef {{contextName: string, clusterName: string, currentContext: boolean | string}} Context
 */

/**
 * @typedef {{id:string, configuredServer: string, contextName: string, clusterConfigured: boolean, contexts: Array.<Context>}} KubernetesCluster
 */

/**
 * @typedef {Array.<KubernetesCluster>} KubernetesClusters
 */

/**
 * @typedef {{url :string, APIKey :string, boardSearch: string, boards: [], selectedBoardsConfigs: []}} grafana
 */

/**
 * @typedef {Array.<grafana>} connectedGrafanas
 */

/**
 * @typedef {{config: any}} prometheus
 */

/**
 * @typedef {Array.<prometheus>} connectedPrometheus
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

export const fetchKubernetesClustersThunk = createAsyncThunk("mesheryEnvironment/fetchKubernetesClusters", async () => {
  const response = await fetchKuberernetesClusters();
  return response;
});

export const submitKubernetesClusterConfigThunk = createAsyncThunk(
  "mesheryEnvironment/updateKubernetesCluster",
  async (payload) => {
    const response = await submitKubernetesClusterConfig(payload);
    return response;
  }
);

export const submitGrafanaConfigThunk = createAsyncThunk("mesheryEnvironment/submitGrafanaConfig", async () => {});

export const submitPrometheusConfigThunk = createAsyncThunk(
  "mesheryEnvironment/submitPrometheusConfig",
  async () => {}
);

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
  },
});

export default mesheryEnvironmentSlice.reducer;
export const { updateKubernetesClusterData, addKubernetesCluster, removeKubernetesCluster, addGrafana, addPrometheus } =
  mesheryEnvironmentSlice.actions;

// selectors should be written and exported
