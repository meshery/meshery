import { describe, it, expect } from 'vitest';
import sessionReducer, { clearSessionExpired, selectIsSessionExpired } from '../session';

describe('session slice', () => {
  it('returns initial state with isExpired false', () => {
    const state = sessionReducer(undefined, { type: 'unknown' });
    expect(state.isExpired).toBe(false);
  });

  it('sets isExpired true on SESSION_EXPIRED action', () => {
    const state = sessionReducer(undefined, { type: 'SESSION_EXPIRED' });
    expect(state.isExpired).toBe(true);
  });

  it('clearSessionExpired resets isExpired to false', () => {
    let state = sessionReducer(undefined, { type: 'SESSION_EXPIRED' });
    expect(state.isExpired).toBe(true);
    state = sessionReducer(state, clearSessionExpired());
    expect(state.isExpired).toBe(false);
  });

  it('selectIsSessionExpired selector works', () => {
    const rootState = {
      session: { isExpired: true },
    };
    expect(selectIsSessionExpired(rootState)).toBe(true);
  });
});
