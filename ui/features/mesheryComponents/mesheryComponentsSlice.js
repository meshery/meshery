import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { fetchAvailableAdapters } from "./adapters";
import { fetchMesheryServerDetails } from "./helpers";
import {
  changeOperatorStatus,
  fetchMesheryComponentsStatus,
  initialiseOperatorStatusEventsSubscriptions,
} from "./operator";

/**
    this feature/module is responsible for all the data and operations concerning Connection status, 
    operations, configurations of all the meshery components,
    - Meshery operator
    - Meshsync
    - Meshery Broker (NATS)
    - Meshery Server
    - Meshery Adapters
  */

/**
 * @typedef {"UNKNOWN" | "DISABLED" | "ENABLED" | "CONNECTED"} ConnectionStatus
 */

/**
 * @typedef {Object} MesheryComponent
 * @property {ConnectionStatus} connectionStatus
 * @property {string | Object} version
 */

/**
 * @typedef {{build: string, latest: string, outdated: boolean, commitsha: string, release_channel: string}} MesheryServerVersionType
 */

/**
 * @typedef {Object} AdapterType
 * @property {string} adapter_location
 * @property {string} name
 * @property {string} version
 * @property {string} git_commit_sha
 * @property {Array | null} ops
 * @property {boolean} isActive
 */

/**
 * @typedef {Array.<AdapterType>} AdaptersListType
 */

const initialState = {
  /** @type {MesheryComponent}  */
  operator: {
    connectionStatus: "UNKNOWN",
    version: "UNKNOWN",
  },

  /** @type {MesheryComponent}  */
  meshsync: {
    connectionStatus: "UNKOWN",
    version: "UNKNOWN",
  },

  /** @type {MesheryComponent}  */
  broker: {
    connectionStatus: "UNKOWN",
    version: "UNKNOWN",
  },

  /** @type {MesheryComponent}  */
  server: {
    connectionStatus: "UNKOWN",
    /** @type {MesheryServerVersionType} */
    version: "UNKNOWN",
  },
  /**
   * @type {AdaptersListType}
   */
  adapters: [],
  loading: false,

  /** @type {{code: number | null, description: string | null}} */
  operatorError: {
    code: null,
    description: null,
  },

  subscription: {
    initialised: false,
    disposer: null,
  }, // subscription object
};

export const fetchComponentsStatusThunk = createAsyncThunk("mesheryComponents/fetchComponentsStatus", async () => {
  const response = await fetchMesheryComponentsStatus();
  return response;
});

export const initialiseOperatorStatusSubscriptionThunk = createAsyncThunk(
  "mesheryComponents/initialiseOperatorStatusEventsSubscription",
  async (cb) => {
    const response = await initialiseOperatorStatusEventsSubscriptions(cb); // cb should dispatch updateConnectionStatus action
    return response;
  }
);

export const fetchAvailableAdaptersThunk = createAsyncThunk("mesheryComponents/fetchAvailableAdapters", async () => {
  const response = await fetchAvailableAdapters();
  return response;
});

export const changeOperatorStateThunk = createAsyncThunk(
  "mesheryComponents/changeOperatorState",
  async (desiredState) => {
    // TODO: Error handling should be done as errors will be passed in resolved object
    const response = await changeOperatorStatus(desiredState);
    return response;
  }
);

export const fetchMesheryServerDetailsThunk = createAsyncThunk(
  "mesheryComponents/fetchMesheryServerDetails",
  async () => {
    const response = await fetchMesheryServerDetails();
    return response;
  }
);

const mesheryComponentsSlice = createSlice({
  name: "mesheryComponents",
  initialState,
  reducers: {
    // use `reduce-reducers` to combine multiple reducers
    // (Refer: https://github.com/reduxjs/redux-toolkit/issues/259#issuecomment-604496169)

    // redux-toolkit uses IMMER under the hood which allows us to mutate the state without actually mutating it
    updateConnectionStatus: (state, action) => {
      state.operator.connectionStatus = action.payload?.operator?.status || "UNKNOWN";

      state.broker.connectionStatus = action.payload?.broker?.status || "UNKNOWN";
      state.broker.version = action.payload?.broker?.version || "UNKNOWN";

      state.meshsync.connectionStatus = action.payload?.meshsync?.status || "UNKNOWN";
      state.meshsync.version = action.payload?.meshsync?.version || "UNKNOWN";

      if (
        action.payload.operator.error &&
        (action.payload.operator.error.code || action.payload.operator.error.description)
      )
        state.operatorError = action.payload.operator.error;
      return state;
    },
    updateMesheryServerDetails: (state, action) => {
      state.server.connectionStatus = action.payload?.status || "UNKNOWN";
      state.server.version = action.payload?.version || "UNKNOWN";
      return state;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchComponentsStatusThunk.pending, (state) => {
      state.loading = true;
      return state;
    });
    builder.addCase(fetchComponentsStatusThunk.fulfilled, (state, action) => {
      state.loading = false;
      return mesheryComponentsSlice.caseReducers.updateConnectionStatus(state, action);
    });
    builder.addCase(fetchComponentsStatusThunk.rejected, (state) => {
      state.loading = false;
      return state;
    });

    builder.addCase(fetchAvailableAdaptersThunk.pending, (state) => {
      state.loading = true;
      return state;
    });
    builder.addCase(fetchAvailableAdaptersThunk.fulfilled, (state, action) => {
      state.loading = false;
      state.adapters = [];
      action.payload.forEach((adapter) => {
        state.adapters.push({
          name: adapter.name,
          version: adapter.version,
          git_commit_sha: adapter.git_commit_sha,
          ops: adapter.ops,
          adapter_location: adapter.adapter_location,
          isActive: adapter.git_commit_sha !== "" && adapter.ops.length !== 0 ? true : false,
        });
      });
      return state;
    });
    builder.addCase(fetchAvailableAdaptersThunk.rejected, (state) => {
      state.loading = false;
      return state;
    });

    builder.addCase(initialiseOperatorStatusSubscriptionThunk.fulfilled, (state) => {
      state.subscription.initialised = true;
      // state.subscription.disposer = action.payload;
      return state;
    });
    builder.addCase(initialiseOperatorStatusSubscriptionThunk.rejected, (state) => {
      state.subscription.initialised = false;
      return state;
    });

    builder.addCase(changeOperatorStateThunk.pending, (state) => {
      state.loading = true;
      return state;
    }),
      builder.addCase(changeOperatorStateThunk.fulfilled, (state) => {
        state.loading = false;
        // handle other state updates
        // no need to handle the state change since it will be handled by subscription

        return state;
      });
    builder.addCase(changeOperatorStateThunk.rejected, (state) => {
      state.loading = false;
      return state;
    });

    builder.addCase(fetchMesheryServerDetailsThunk.pending, (state) => {
      state.loading = true;
      return state;
    });
    builder.addCase(fetchMesheryServerDetailsThunk.fulfilled, (state, action) => {
      state.loading = false;
      state.server.version = action.payload;
      return state;
    });
    builder.addCase(fetchMesheryServerDetailsThunk.rejected, (state) => {
      state.loading = false;
      return state;
    });
  },
});

export default mesheryComponentsSlice.reducer;
export const { updateConnectionStatus } = mesheryComponentsSlice.actions;

export const operatorSelector = (state) => state.mesheryComponents.operator;
export const adaptersSelector = (state) => state.mesheryComponents.adapters;
export const loadingSelector = (state) => state.mesheryComponents.loading;

export const mesheryComponentsSelector = (state) => {
  return {
    operator: state.mesheryComponents.operator,
    meshsync: state.mesheryComponents.meshsync,
    broker: state.mesheryComponents.broker,
    server: state.mesheryComponents.server,
  };
};
