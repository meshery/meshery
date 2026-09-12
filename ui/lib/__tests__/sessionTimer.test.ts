import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

describe('sessionTimer', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    // Reset module state between tests
    vi.resetModules();
    // Remove any leftover overlay from a prior test
    document.body.innerHTML = '';
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
    document.body.innerHTML = '';
  });

  it('recordActivity reschedules the warning timer without crashing', async () => {
    const { recordActivity } = await import('../sessionTimer');
    expect(() => recordActivity()).not.toThrow();

    // Fast-forward less than the warning window — no warning yet
    vi.advanceTimersByTime(60 * 1000);
    expect(document.getElementById('meshery-session-warning')).toBeNull();
  });

  it('shows a warning dialog after the inactivity threshold elapses', async () => {
    const { recordActivity } = await import('../sessionTimer');
    recordActivity();

    // Advance time beyond SESSION_TIMEOUT_MS - WARNING_BEFORE_MS (55 minutes)
    vi.advanceTimersByTime(55 * 60 * 1000);

    const overlay = document.getElementById('meshery-session-warning');
    expect(overlay).not.toBeNull();
    expect(overlay?.getAttribute('role')).toBe('dialog');
    expect(overlay?.getAttribute('aria-modal')).toBe('true');
  });

  it('contains "Stay Logged In" and "Log Out" buttons in the warning dialog', async () => {
    const { recordActivity } = await import('../sessionTimer');
    recordActivity();
    vi.advanceTimersByTime(55 * 60 * 1000);

    const overlay = document.getElementById('meshery-session-warning');
    expect(overlay).not.toBeNull();

    const buttons = overlay!.querySelectorAll('button');
    expect(buttons.length).toBe(2);
    const buttonTexts = Array.from(buttons).map((b) => b.textContent);
    expect(buttonTexts).toContain('Stay Logged In');
    expect(buttonTexts).toContain('Log Out');
  });

  it('clicking "Stay Logged In" sends a refresh fetch and dismisses the warning', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true });
    global.fetch = fetchMock as unknown as typeof fetch;

    const { recordActivity } = await import('../sessionTimer');
    recordActivity();
    vi.advanceTimersByTime(55 * 60 * 1000);

    const overlay = document.getElementById('meshery-session-warning')!;
    const stayBtn = Array.from(overlay.querySelectorAll('button')).find(
      (b) => b.textContent === 'Stay Logged In',
    )!;
    expect(stayBtn).toBeDefined();

    stayBtn.click();

    expect(fetchMock).toHaveBeenCalledWith('/api/user/prefs', { credentials: 'include' });
    expect(document.getElementById('meshery-session-warning')).toBeNull();
  });

  it('clicking "Stay Logged In" handles fetch rejection silently', async () => {
    const fetchMock = vi.fn().mockRejectedValue(new Error('network'));
    global.fetch = fetchMock as unknown as typeof fetch;

    const { recordActivity } = await import('../sessionTimer');
    recordActivity();
    vi.advanceTimersByTime(55 * 60 * 1000);

    const overlay = document.getElementById('meshery-session-warning')!;
    const stayBtn = Array.from(overlay.querySelectorAll('button')).find(
      (b) => b.textContent === 'Stay Logged In',
    )!;

    expect(() => stayBtn.click()).not.toThrow();
    expect(fetchMock).toHaveBeenCalled();
  });

  it('clicking "Log Out" navigates to /user/logout', async () => {
    const originalLocation = window.location;
    const setHrefSpy = vi.fn();

    Object.defineProperty(window, 'location', {
      configurable: true,
      writable: true,
      value: {
        ...originalLocation,
        set href(value: string) {
          setHrefSpy(value);
        },
        get href() {
          return originalLocation.href;
        },
      },
    });

    try {
      const { recordActivity } = await import('../sessionTimer');
      recordActivity();
      vi.advanceTimersByTime(55 * 60 * 1000);

      const overlay = document.getElementById('meshery-session-warning')!;
      const logoutBtn = Array.from(overlay.querySelectorAll('button')).find(
        (b) => b.textContent === 'Log Out',
      )!;
      logoutBtn.click();

      expect(setHrefSpy).toHaveBeenCalledWith('/user/logout');
    } finally {
      Object.defineProperty(window, 'location', {
        configurable: true,
        writable: true,
        value: originalLocation,
      });
    }
  });

  it('startSessionTimer registers click and keydown listeners that record activity', async () => {
    const addEventListenerSpy = vi.spyOn(window, 'addEventListener');
    const { startSessionTimer } = await import('../sessionTimer');
    startSessionTimer();

    const eventNames = addEventListenerSpy.mock.calls.map((c) => c[0]);
    expect(eventNames).toContain('click');
    expect(eventNames).toContain('keydown');
  });

  it('startSessionTimer schedules a warning even without explicit activity', async () => {
    const { startSessionTimer } = await import('../sessionTimer');
    startSessionTimer();
    vi.advanceTimersByTime(55 * 60 * 1000);
    expect(document.getElementById('meshery-session-warning')).not.toBeNull();
  });

  it('recordActivity dismisses an already-shown warning', async () => {
    const { recordActivity } = await import('../sessionTimer');
    recordActivity();
    vi.advanceTimersByTime(55 * 60 * 1000);
    expect(document.getElementById('meshery-session-warning')).not.toBeNull();

    // Subsequent activity should dismiss the warning
    recordActivity();
    expect(document.getElementById('meshery-session-warning')).toBeNull();
  });

  it('does not show the warning twice when timer fires repeatedly', async () => {
    const { recordActivity } = await import('../sessionTimer');
    recordActivity();
    vi.advanceTimersByTime(55 * 60 * 1000);
    const overlaysAfterFirst = document.querySelectorAll('#meshery-session-warning');
    expect(overlaysAfterFirst.length).toBe(1);

    // Advance more time. The warning timer is already done, but extra time
    // shouldn't produce additional dialogs unless rescheduling happens.
    vi.advanceTimersByTime(10 * 60 * 1000);
    const overlaysAfterMore = document.querySelectorAll('#meshery-session-warning');
    expect(overlaysAfterMore.length).toBe(1);
  });
});
