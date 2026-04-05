import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

let mockStatus = 'idle';

vi.mock('../../store', () => ({
  store: {
    getState: () => ({ sessions: { status: mockStatus } }),
    dispatch: vi.fn(),
  },
}));

import { store } from '../../store';
import { setSessionState } from '@/store/slices/sessions';
import { recordActivity, triggerSessionExpired, startSessionTimer } from '../sessionTimer';

describe('sessionTimer', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    mockStatus = 'idle';
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('startSessionTimer registers click and keydown listeners', () => {
    const addEventSpy = vi.spyOn(window, 'addEventListener');
    startSessionTimer();

    const registeredEvents = addEventSpy.mock.calls.map((call) => call[0]);
    expect(registeredEvents).toContain('click');
    expect(registeredEvents).toContain('keydown');

    addEventSpy.mockRestore();
  });

  it('startSessionTimer does not register duplicate listeners', () => {
    const addEventSpy = vi.spyOn(window, 'addEventListener');
    startSessionTimer();

    expect(addEventSpy).not.toHaveBeenCalled();

    addEventSpy.mockRestore();
  });

  it('recordActivity dispatches idle when session is expiring', () => {
    mockStatus = 'expiring';
    recordActivity();
    expect(store.dispatch).toHaveBeenCalledWith(setSessionState('idle'));
  });

  it('recordActivity does nothing when session is expired', () => {
    mockStatus = 'expired';
    recordActivity();
    expect(store.dispatch).not.toHaveBeenCalled();
  });

  it('dispatches expiring after timeout minus warning period', () => {
    recordActivity();

    // Advance to 55 minutes (SESSION_TIMEOUT_MS - WARNING_BEFORE_MS)
    vi.advanceTimersByTime(55 * 60 * 1000);

    expect(store.dispatch).toHaveBeenCalledWith(setSessionState('expiring'));
  });

  it('does not dispatch expiring before the warning period', () => {
    recordActivity();
    vi.clearAllMocks();

    // Advance to 54 minutes — should not fire yet
    vi.advanceTimersByTime(54 * 60 * 1000);

    expect(store.dispatch).not.toHaveBeenCalledWith(setSessionState('expiring'));
  });

  it('triggerSessionExpired dispatches expired', () => {
    triggerSessionExpired();
    expect(store.dispatch).toHaveBeenCalledWith(setSessionState('expired'));
  });

  it('triggerSessionExpired is a no-op when already expired', () => {
    mockStatus = 'expired';
    triggerSessionExpired();
    expect(store.dispatch).not.toHaveBeenCalled();
  });
});
