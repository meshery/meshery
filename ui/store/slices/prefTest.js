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

const prefTestSlice = createSlice({
  name: 'prefTest',
  initialState,
  reducers: {
    updateLoadTestPref: (state, action) => {
      const payload = action.payload.loadTestPref;
      if (payload.ts && payload.ts instanceof Date) {
        payload.ts = payload.ts.getTime();
      }
      state.loadTestPref = payload;
    },
    updateLoadTest: (state, action) => {
      state.loadTest = action.payload.loadTest;
    },
    updateResultsSelection: (state, action) => {
      if (Object.keys(action.payload.results).length > 0) {
        state.results_selection[action.payload.page] = action.payload.results;
      } else {
        delete state.results_selection[action.payload.page];
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
