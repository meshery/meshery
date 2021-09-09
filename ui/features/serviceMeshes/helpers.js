import subscribeServiceMeshEvents from "./graphql/subscriptions/ServiceMeshSubscription";

/**
 *
 * @param {import("./serviceMeshesSlice").MeshNameType} meshFilter
 * @param {(dataPlaneData: {dataPlanesState: Array.<import("./serviceMeshesSlice").DataPlaneDataType>}, controlPlaneData: {controlPlanesState: Array.<import("./serviceMeshesSlice").ControlPlaneDataType>}) => void} dataCB
 * @returns
 */
// eslint-disable-next-line no-unused-vars
export const initServiceMeshEvents = (meshFilter = "ALL_MESH", dataCB) =>
  new Promise((res) => {
    const ALL_MESH = {};
    subscribeServiceMeshEvents(dataCB, ALL_MESH);
    // handle errors
    res("Subscription initialised succesfully");
  });
