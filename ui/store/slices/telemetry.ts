import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  grafana: {
    grafanaURL: '',
    grafanaAPIKey: '',
    grafanaBoardSearch: '',
    grafanaBoards: [],
    selectedBoardsConfigs: [],
    ts: -8640000000000000,
  },
  prometheus: {
    prometheusURL: '',
    selectedPrometheusBoardsConfigs: [],
    ts: -8640000000000000,
  },
  staticPrometheusBoardConfig: {},
};

const telemetrySlice = createSlice({
  name: 'telemetry',
  initialState,
  reducers: {
    updateGrafanaConfig: (state, action) => {
      state.grafana = { ...state.grafana, ...action.payload };
      state.grafana.ts = Date.now();
    },
    updatePrometheusConfig: (state, action) => {
      state.prometheus = { ...state.prometheus, ...action.payload };
      state.prometheus.ts = Date.now();
    },
    updateStaticPrometheusBoardConfig: (state, action) => {
      state.staticPrometheusBoardConfig = action.payload;
    },
  },
});
export const { updateGrafanaConfig, updatePrometheusConfig, updateStaticPrometheusBoardConfig } =
  telemetrySlice.actions;

export default telemetrySlice.reducer;
