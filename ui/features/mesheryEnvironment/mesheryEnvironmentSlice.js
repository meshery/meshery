import { createSlice } from "@reduxjs/toolkit";

/**
    this feature/module is responsible for all the data and operations regarding operations and
    configurations regarding Kubernetes cluster, connected Grafana, Connected prometheus etc.
  */

const initialState = {
  kubernetesCluster: {},
  grafana: {},
  prometheus: {},
  // and other relavant data
};

const mesheryEnvironmentSlice = createSlice({
  name: "mesheryEnvironment",
  initialState,
  reducers: {
    // reducers to update the state in this slice
  },
  extraReducers: {
    // for dealing with thunk and observable
  },
});

export default mesheryEnvironmentSlice.reducer;
// export const {} = mesheryEnvironmentSlice.actions;

// selectors should be written and exported
