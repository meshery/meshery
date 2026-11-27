import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  loadTest: {
    testName: '',
    meshName: '',
    url: '',
    qps: 0,
    c: 0,
    t: '30s',
    result: {},
  },
  loadTestPref: {
    qps: 0,
    t: '30s',
    c: 0,
    gen: 'fortio',
    ts: Date.now(),
  },
  results_selection: {},
  results: {
    startKey: '',
    results: [],
  },
};

const isPlainObject = (value: unknown): boolean =>
  typeof value === 'object' && value !== null && !Array.isArray(value);

const isValidSelectionPage = (page: unknown): boolean =>
  (typeof page === 'number' && Number.isFinite(page) && page >= 0) ||
  (typeof page === 'string' && page.trim().length > 0);

const prefTestSlice = createSlice({
  name: 'prefTest',
  initialState,
  reducers: {
    updateLoadTestPref: (state, action) => {
      const payload = action.payload?.loadTestPref;
      if (!isPlainObject(payload)) {
        return;
      }
      state.loadTestPref =
        payload.ts instanceof Date ? { ...payload, ts: payload.ts.getTime() } : payload;
    },
    updateLoadTest: (state, action) => {
      const payload = action.payload?.loadTest;
      if (!isPlainObject(payload)) {
        return;
      }
      state.loadTest = payload;
    },
    updateResultsSelection: (state, action) => {
      const page = action.payload?.page;
      if (!isValidSelectionPage(page)) {
        return;
      }
      const results = action.payload?.results;
      if (isPlainObject(results) && Object.keys(results).length > 0) {
        state.results_selection[page] = results;
      } else {
        delete state.results_selection[page];
      }
    },
    clearResultsSelection: (state) => {
      state.results_selection = {};
    },
  },
});

export const { updateLoadTestPref, updateLoadTest, updateResultsSelection, clearResultsSelection } =
  prefTestSlice.actions;

export default prefTestSlice.reducer;
