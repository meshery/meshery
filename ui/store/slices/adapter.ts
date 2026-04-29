import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  meshAdapters: [],
  meshAdaptersts: Date.now(),
  selectedAdapter: '',
};

const adapterSlice = createSlice({
  name: 'adapter',
  initialState,
  reducers: {
    updateAdaptersInfo: (state, action) => {
      state.meshAdapters = action.payload.meshAdapters;
      state.meshAdaptersts = Date.now();
    },
    setAdapter: (state, action) => {
      state.selectedAdapter = action.payload.selectedAdapter;
    },
  },
});

export const { updateAdaptersInfo, setAdapter } = adapterSlice.actions;

export default adapterSlice.reducer;
