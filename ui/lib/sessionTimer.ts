/**
 * Proactive session timeout detection.
 *
 * Tracks the last successful API activity and warns the user before
 * their session expires. The backend remains authoritative — this is
 * purely a UX affordance so the user isn't surprised by a redirect.
 *
 * Call startSessionTimer() once during app bootstrap.
 */

const SESSION_TIMEOUT_MS = 60 * 60 * 1000; // 1 hour (matches server default)
const WARNING_BEFORE_MS = 5 * 60 * 1000; // Warn 5 minutes before expiry

let lastActivity = Date.now();
let warningTimer: ReturnType<typeof setTimeout> | null = null;
let warningShown = false;

export function recordActivity() {
  lastActivity = Date.now();
  if (warningShown) {
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
  if (warningShown || typeof window === 'undefined') return;
  warningShown = true;
  window.dispatchEvent(new Event('meshery:session-warning'));
}

function dismissWarning() {
  warningShown = false;
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new Event('meshery:session-warning-dismiss'));
  }
}

export function startSessionTimer() {
  if (typeof window === 'undefined') return;

  // Reset activity on any user interaction
  const resetOnInteraction = () => recordActivity();
  window.addEventListener('click', resetOnInteraction, { passive: true });
  window.addEventListener('keydown', resetOnInteraction, { passive: true });

  scheduleWarning();
}
