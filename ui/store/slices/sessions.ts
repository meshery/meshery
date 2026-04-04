import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type SessionState = 'idle' | 'expiring' | 'expired';

interface SessionsState {
  status: SessionState;
}

const initialState: SessionsState = {
  status: 'idle',
};

const sessionsSlice = createSlice({
  name: 'sessions',
  initialState,
  reducers: {
    setSessionState: (state, action: PayloadAction<SessionState>) => {
      // expired is an absorbing state — cannot transition out of it
      if (state.status === 'expired' && action.payload !== 'expired') return;
      state.status = action.payload;
    },
  },
});

export const { setSessionState } = sessionsSlice.actions;
export default sessionsSlice.reducer;
