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
      const pref = action.payload?.loadTestPref;
      if (!pref || typeof pref !== 'object') return;
      const ts = pref.ts instanceof Date ? pref.ts.getTime() : pref.ts;
      state.loadTestPref = { ...pref, ts };
    },
    updateLoadTest: (state, action) => {
      const loadTest = action.payload?.loadTest;
      if (!loadTest || typeof loadTest !== 'object') return;
      state.loadTest = loadTest;
    },
    updateResultsSelection: (state, action) => {
      const results = action.payload?.results;
      const page = action.payload?.page;
      if (results == null || typeof results !== 'object') return;
      if (page == null) return;
      if (typeof page !== 'number' && typeof page !== 'string') return;
      if (typeof page === 'number' && !isFinite(page)) return;
      const pageKey = String(page);
      if (Object.keys(results).length > 0) {
        state.results_selection[pageKey] = results;
      } else {
        delete state.results_selection[pageKey];
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
