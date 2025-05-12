import { createSlice } from "@reduxjs/toolkit";

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
  },
});

export const { updateLoadTestPref, updateLoadTest } = prefTestSlice.actions;

export default prefTestSlice.reducer;
