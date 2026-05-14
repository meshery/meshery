import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type ExpirationState = 'idle' | 'expiring' | 'expired';

interface SessionsState {
  status: ExpirationState;
}

const initialState: SessionsState = {
  status: 'idle',
};

const sessionsSlice = createSlice({
  name: 'sessions',
  initialState,
  reducers: {
    setExpirationState: (state, action: PayloadAction<ExpirationState>) => {
      // expired is an absorbing state — cannot transition out of it
      if (state.status === 'expired' && action.payload !== 'expired') return;
      state.status = action.payload;
    },
  },
});

export const { setExpirationState } = sessionsSlice.actions;

export const selectExpirationStatus = (state: { sessions: SessionsState }) => state.sessions.status;

export default sessionsSlice.reducer;
