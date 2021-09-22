/* eslint-disable react/prop-types */
import React from "react";
import { render as rtlRender } from "@testing-library/react";
import { configureStore } from "@reduxjs/toolkit";
import { Provider } from "react-redux";
// Import your own reducer
import serviceMeshesReducer from "@/features/serviceMeshes/serviceMeshesSlice";
import mesheryEnvironmentReducer from "@/features/mesheryEnvironment/mesheryEnvironmentSlice";
import mesheryComponentsReducer from "@/features/mesheryComponents/mesheryComponentsSlice";
import providerReducer from "@/features/provider/providerSlice";
import performanceReducer from "@/features/performance/performanceSlice";

function render(
  ui,
  {
    preloadedState,
    store = configureStore({
      reducer: {
        serviceMeshes: serviceMeshesReducer,
        mesheryEnvironment: mesheryEnvironmentReducer,
        mesheryComponents: mesheryComponentsReducer,
        provider: providerReducer,
        performance: performanceReducer,
      },
      preloadedState,
    }),
    ...renderOptions
  } = {}
) {
  function Wrapper({ children }) {
    return <Provider store={store}>{children}</Provider>;
  }
  return rtlRender(ui, { wrapper: Wrapper, ...renderOptions });
}

// re-export everything
export * from "@testing-library/react";
// override render method
export { render };
