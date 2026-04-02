import { createSlice } from '@reduxjs/toolkit';

interface SessionState {
  isExpired: boolean;
}

const initialState: SessionState = {
  isExpired: false,
};

const sessionSlice = createSlice({
  name: 'session',
  initialState,
  reducers: {
    clearSessionExpired: (state) => {
      state.isExpired = false;
    },
  },
  extraReducers: (builder) => {
    builder.addMatcher(
      (action) => action.type === 'SESSION_EXPIRED',
      (state) => {
        state.isExpired = true;
      },
    );
  },
});

export const { clearSessionExpired } = sessionSlice.actions;

export const selectIsSessionExpired = (state: { session: SessionState }) => state.session.isExpired;

export default sessionSlice.reducer;
