/**
 * Proactive session timeout detection.
 *
 * Tracks the last successful API activity and warns the user before
 * their session expires. The backend remains authoritative — this is
 * purely a UX affordance so the user isn't surprised by a redirect.
 *
 * Call startSessionTimer() once during app bootstrap.
 */

import { store } from '../store';
import { setSessionState } from '../store/slices/sessions';

const SESSION_TIMEOUT_MS = 60 * 60 * 1000; // 1 hour (matches server default)
const WARNING_BEFORE_MS = 5 * 60 * 1000; // Warn 5 minutes before expiry

let lastActivity = Date.now();
let warningTimer: ReturnType<typeof setTimeout> | null = null;

export function recordActivity() {
  const sessionState = store.getState().sessions.status;
  if (sessionState === 'expired') return;

  lastActivity = Date.now();

  if (sessionState === 'expiring') {
    dismissWarning();
  }

  scheduleWarning();
}

function scheduleWarning() {
  if (warningTimer) clearTimeout(warningTimer);

  const timeUntilWarning = SESSION_TIMEOUT_MS - WARNING_BEFORE_MS;
  const elapsed = Date.now() - lastActivity;
  const delay = Math.max(0, timeUntilWarning - elapsed);

  warningTimer = setTimeout(showTimeoutWarning, delay);
}

function showTimeoutWarning() {
  const state = store.getState().sessions.status;

  if (state !== 'idle' || typeof document === 'undefined') return;

  store.dispatch(setSessionState('expiring'));
}

function dismissWarning() {
  store.dispatch(setSessionState('idle'));
}

export function triggerSessionExpired() {
  if (store.getState().sessions.status === 'expired') return;
  if (warningTimer) clearTimeout(warningTimer);
  store.dispatch(setSessionState('expired'));
}

export function startSessionTimer() {
  if (typeof window === 'undefined') return;

  // Reset activity on any user interaction
  const resetOnInteraction = () => recordActivity();
  window.addEventListener('click', resetOnInteraction, { passive: true });
  window.addEventListener('keydown', resetOnInteraction, { passive: true });

  scheduleWarning();
}
