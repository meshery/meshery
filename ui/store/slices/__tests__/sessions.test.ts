import { describe, it, expect } from 'vitest';
import sessionsReducer, { setExpirationState } from '../sessions';

describe('sessions slice', () => {
  it('returns initial state with idle status', () => {
    const state = sessionsReducer(undefined, { type: 'unknown' });
    expect(state.status).toBe('idle');
  });

  it('transitions from idle to expiring', () => {
    const state = sessionsReducer(undefined, setExpirationState('expiring'));
    expect(state.status).toBe('expiring');
  });

  it('transitions from idle to expired', () => {
    const state = sessionsReducer(undefined, setExpirationState('expired'));
    expect(state.status).toBe('expired');
  });

  it('transitions from expiring to idle', () => {
    const initial = sessionsReducer(undefined, setExpirationState('expiring'));
    const state = sessionsReducer(initial, setExpirationState('idle'));
    expect(state.status).toBe('idle');
  });

  it('transitions from expiring to expired', () => {
    const initial = sessionsReducer(undefined, setExpirationState('expiring'));
    const state = sessionsReducer(initial, setExpirationState('expired'));
    expect(state.status).toBe('expired');
  });

  it('does not transition from expired to idle (absorbing state)', () => {
    const initial = sessionsReducer(undefined, setExpirationState('expired'));
    const state = sessionsReducer(initial, setExpirationState('idle'));
    expect(state.status).toBe('expired');
  });

  it('does not transition from expired to expiring (absorbing state)', () => {
    const initial = sessionsReducer(undefined, setExpirationState('expired'));
    const state = sessionsReducer(initial, setExpirationState('expiring'));
    expect(state.status).toBe('expired');
  });

  it('allows expired to expired (idempotent)', () => {
    const initial = sessionsReducer(undefined, setExpirationState('expired'));
    const state = sessionsReducer(initial, setExpirationState('expired'));
    expect(state.status).toBe('expired');
  });
});
