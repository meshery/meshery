import { configureStore, type Reducer } from '@reduxjs/toolkit';

export const sampleCapabilities = {
  providerType: 'remote',
  providerName: 'Meshery Cloud',
  capabilities: {
    canPersistData: true,
  },
  extensions: {
    navigator: [],
  },
};

export const sampleAdapters = [
  {
    adapter_location: 'localhost:10000',
    name: 'Istio',
    version: '1.20.0',
  },
];

export const createTestStoreWithReducers = (
  reducers: { ui: Reducer; adapter: Reducer },
  initialState?: Record<string, unknown>,
) => {
  return configureStore({
    reducer: reducers,
    preloadedState: initialState,
  });
};
