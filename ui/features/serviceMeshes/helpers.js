import subscribeAddonStatusEvents from "./graphql/subscriptions/AddonStatusSubscription";
import subscribeServiceMeshEvents from "./graphql/subscriptions/ServiceMeshSubscription";

/**
 *
 * @param {import("./serviceMeshesSlice").MeshNameType} meshFilter
 * @param {(dataPlaneData: {dataPlanesState: Array.<import("./serviceMeshesSlice").DataPlaneDataType>}, controlPlaneData: {controlPlanesState: Array.<import("./serviceMeshesSlice").ControlPlaneDataType>}) => void} dataCB
 * @returns
 */
// eslint-disable-next-line no-unused-vars
export const initServiceMeshEvents = (meshFilter = {}, dataCB) =>
  new Promise((res) => {
    const ALL_MESH = {};
    subscribeServiceMeshEvents(dataCB, ALL_MESH);
    // handle errors
    res("Subscription initialised succesfully");
  });

export const initAddonStatusSubscription = (meshFilter = {}, dataCB) =>
  new Promise((res) => {
    subscribeAddonStatusEvents(dataCB, meshFilter);
    // handle errors
    res("Subscription initialised succesfully");
  });
