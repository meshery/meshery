import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { initAddonStatusSubscription, initServiceMeshEvents } from "./helpers";

/**
   this feature/module handles all things related to service meshes like, getting the control plane data, proxy data of service mesh, 
   doing operations on service mesh like deploying meshes in cluster, deploying apps in cluster, all those stuff
   This might have a dependency on the MesheryEnvironment feature to get data like, connection status of cluster, adapters, metrics data etc.
  */

// eslint-disable-next-line

/**
 * @typedef {{name: string, members: Array.<{name: string, version: string, component: string, namespace: string}>}} ControlPlaneDataType
 */

/**
 * @typedef {{controlPlaneMemberName : string, containerName: string, image: string, status : {containerStatusName: string, image: string, state: string, lastState: string, restartCount: string | number, ready: boolean | string, started: Date | string | boolean, imageID: string, containerID: string}, ports: Array.<{name:string, containerPort: string, protocol: string}>, resources: Array.<{}>}} ProxyType
 */

/**
 * @typedef {{name: string, proxies: Array.<ProxyType>}} DataPlaneDataType
 */

/**
 * @typedef {{name: MeshNameType, controlPlaneData: ControlPlaneDataType , dataPlaneData: DataPlaneDataType , adapterId: string, addons: AddonType[] }} ServiceMeshType
 */

/**
 * @typedef {"all_mesh" | "invalid_mesh" | "app_mesh" | "citrix_service_mesh" | "consul" | "istio" | "kuma" | "linkerd" | "traefik_mesh" | "octarine" | "network_service_mesh" | "tanzu" | "open_service_mesh" | "nginx_service_mesh"} MeshNameType
 */

/**
 * @typedef {{name: string, endpoint:string, owner: MeshNameType}} AddonType
 */

const initialState = {
  /** @type {Array<ServiceMeshType>} */
  meshes: [],
  loading: false,
  /** @type {import("../provider/providerSlice").Error} */
  error: {},
};

export const initServiceMeshEventsThunk = createAsyncThunk(
  "serviceMeshes/initServiceMeshEvents",
  /**
   * @param {ServiceMeshType} filter
   * @param {(dataPlaneData: {dataPlanesState: import("./serviceMeshesSlice").DataPlaneDataType}, controlPlaneData: {controlPlanesState: import("./serviceMeshesSlice").ControlPlaneDataType}) => void} dataCB
   */
  async (dataCB) => {
    const response = await initServiceMeshEvents("istio", dataCB);

    return response;
  }
);

export const initAddonStatusSubscriptionThunk = createAsyncThunk(
  "serviceMeshes/initAddonStatusSubscriptions",
  async (dataCB) => {
    const ALL_MESH = {};
    const response = await initAddonStatusSubscription(ALL_MESH, dataCB);
    return response;
  }
);

const serviceMeshesSlice = createSlice({
  name: "serviceMeshes",
  initialState,
  reducers: {
    updateServiceMeshesData: (state, action) => {
      /** @type {Array.<ControlPlaneDataType>} */
      const controlPlanesData = action.payload.controlPlanesData?.controlPlanesState;

      /** @type {Array.<DataPlaneDataType>} */
      const dataPlanesData = action.payload.dataPlanesData?.dataPlanesState;

      controlPlanesData?.map((mesh) => {
        const currentMesh = state.meshes.find((_mesh) => _mesh.name === mesh.name);
        const dataPlaneDataForCurrentMesh = dataPlanesData.find((_mesh) => _mesh.name === mesh.name);
        // TODO: write the logic for finding adapter ID
        if (!currentMesh)
          state.meshes.push({
            name: mesh.name,
            controlPlaneData: mesh,
            dataPlaneData: dataPlaneDataForCurrentMesh,
            adapterId: "",
          });
        else {
          currentMesh.controlPlaneData = mesh;
          if (dataPlaneDataForCurrentMesh) currentMesh.dataPlaneData = dataPlaneDataForCurrentMesh;
        }
      });

      return state;
    },
    updateServiceMeshesAddonsStatus: (state, action) => {
      /** @type {AddonType[]} */
      const scannedAddons = action.payload;
      scannedAddons.forEach((addon) => {
        state.meshes.find((mesh) => mesh.name === addon.owner).addons.push(addon);
      });
    },
  },
  extraReducers: (builder) => {
    builder.addCase(initServiceMeshEventsThunk.fulfilled, (state) => {
      state.loading = false;
      return state;
    });

    //proper error handling has to be done for  `initServiceMeshEventsThunk` handler
    builder.addCase(initServiceMeshEventsThunk.rejected, (state, action) => {
      state.loading = false;
      state.error.description = action.error;
      return state;
    });
    builder.addCase(initAddonStatusSubscriptionThunk.fulfilled, (state) => {
      state.loading = false;
      return state;
    });

    //proper error handling has to be done for  `initServiceMeshEventsThunk` handler
    builder.addCase(initAddonStatusSubscriptionThunk.rejected, (state, action) => {
      state.loading = false;
      state.error.description = action.error;
      return state;
    });
  },
});

export default serviceMeshesSlice.reducer;
export const { updateServiceMeshesData } = serviceMeshesSlice.actions;

// selectors should be written and exported
