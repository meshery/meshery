import { configureStore } from "@reduxjs/toolkit";
import serviceMeshesReducer from "@/features/serviceMeshes/serviceMeshesSlice";
import mesheryEnvironmentReducer from "@/features/mesheryEnvironment/mesheryEnvironmentSlice";
import mesheryComponentsReducer from "@/features/mesheryComponents/mesheryComponentsSlice";
import providerReducer from "@/features/provider/providerSlice";

export default configureStore({
  reducer: {
    serviceMeshes: serviceMeshesReducer,
    mesheryEnvironment: mesheryEnvironmentReducer,
    mesheryComponents: mesheryComponentsReducer,
    provider: providerReducer,
  },
});
