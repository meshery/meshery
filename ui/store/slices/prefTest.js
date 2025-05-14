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
    ts: new Date(),
  },
  results_selection: {}, // format - { page: {index: content}}
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
      state.loadTestPref = action.payload.loadTestPref;
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
