import { createSlice } from '@reduxjs/toolkit';

const initialState: Record<string, any> = {};

const tableSlice = createSlice({
  name: 'tableState',
  initialState,
  reducers: {
    updateTableState: (state, action) => {
      const { tableId, update } = action.payload;
      if (!state[tableId]) {
        state[tableId] = {};
      }
      state[tableId] = {
        ...state[tableId],
        ...update,
      };
    },
  },
});

export const { updateTableState } = tableSlice.actions;
export const selectTableState = (state: any, tableId: string) => state.tableState[tableId] || {};
export default tableSlice.reducer;
