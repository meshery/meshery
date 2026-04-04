import { describe, it, expect } from 'vitest';
import sessionsReducer, { setSessionState } from '../sessions';

describe('sessions slice', () => {
  it('returns initial state with idle status', () => {
    const state = sessionsReducer(undefined, { type: 'unknown' });
    expect(state.status).toBe('idle');
  });

  it('transitions from idle to expiring', () => {
    const state = sessionsReducer(undefined, setSessionState('expiring'));
    expect(state.status).toBe('expiring');
  });

  it('transitions from idle to expired', () => {
    const state = sessionsReducer(undefined, setSessionState('expired'));
    expect(state.status).toBe('expired');
  });

  it('transitions from expiring to idle', () => {
    const initial = sessionsReducer(undefined, setSessionState('expiring'));
    const state = sessionsReducer(initial, setSessionState('idle'));
    expect(state.status).toBe('idle');
  });

  it('transitions from expiring to expired', () => {
    const initial = sessionsReducer(undefined, setSessionState('expiring'));
    const state = sessionsReducer(initial, setSessionState('expired'));
    expect(state.status).toBe('expired');
  });

  it('does not transition from expired to idle (absorbing state)', () => {
    const initial = sessionsReducer(undefined, setSessionState('expired'));
    const state = sessionsReducer(initial, setSessionState('idle'));
    expect(state.status).toBe('expired');
  });

  it('does not transition from expired to expiring (absorbing state)', () => {
    const initial = sessionsReducer(undefined, setSessionState('expired'));
    const state = sessionsReducer(initial, setSessionState('expiring'));
    expect(state.status).toBe('expired');
  });

  it('allows expired to expired (idempotent)', () => {
    const initial = sessionsReducer(undefined, setSessionState('expired'));
    const state = sessionsReducer(initial, setSessionState('expired'));
    expect(state.status).toBe('expired');
  });
});
