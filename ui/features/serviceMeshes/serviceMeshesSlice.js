import { createSlice } from "@reduxjs/toolkit";

/**
   this feature/module handles all things related to service meshes like, getting the control plane data, proxy data of service mesh, 
   doing operations on service mesh like deploying meshes in cluster, deploying apps in cluster, all those stuff
   This might have a dependency on the MesheryEnvironment feature to get data like, connection status of cluster, adapters, metrics data etc.
  */

// eslint-disable-next-line
const mesh = "istio" | "kuma" | "linkerd"; //  get this data from meshery adapter, basically which adapters are active

const initialState = {
  istio: {
    controlPlaneData: {},
    dataPlaneData: {},
    addons: {},

    // and other relavant data
  },
};

const serviceMeshesSlice = createSlice({
  name: "serviceMeshes",
  initialState,
  reducers: {
    // reducers to update the state in this slice
  },
  extraReducers: {
    // for dealing with thunk and observable
  },
});

export default serviceMeshesSlice.reducer;
// export const {} = serviceMeshesSlice.actions;

// selectors should be written and exported
