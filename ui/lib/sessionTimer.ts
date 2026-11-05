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
  if (warningShown || typeof document === 'undefined') return;
  warningShown = true;

  const overlay = document.createElement('div');
  overlay.id = 'meshery-session-warning';
  overlay.setAttribute('role', 'dialog');
  overlay.setAttribute('aria-modal', 'true');
  overlay.style.cssText =
    'position:fixed;inset:0;z-index:99998;display:flex;align-items:center;' +
    'justify-content:center;background:rgba(0,0,0,0.4)';

  const box = document.createElement('div');
  box.style.cssText =
    'background:#fff;border-radius:8px;padding:32px;max-width:420px;text-align:center;' +
    'font-family:system-ui,sans-serif;color:#333;box-shadow:0 4px 24px rgba(0,0,0,0.2)';
  box.innerHTML =
    '<h2 style="margin:0 0 12px">Session Expiring Soon</h2>' +
    '<p style="margin:0 0 24px;color:#666">' +
    'Your session will expire in a few minutes due to inactivity. ' +
    'Would you like to stay logged in?</p>';

  const btnRow = document.createElement('div');
  btnRow.style.cssText = 'display:flex;gap:12px;justify-content:center';

  const stayBtn = document.createElement('button');
  stayBtn.textContent = 'Stay Logged In';
  stayBtn.style.cssText =
    'background:#477e96;color:#fff;border:none;border-radius:4px;padding:10px 24px;' +
    'font-size:14px;cursor:pointer';
  stayBtn.onclick = () => {
    // A lightweight request to refresh the session
    fetch('/api/user/prefs', { credentials: 'include' }).catch(() => {});
    dismissWarning();
  };

  const logoutBtn = document.createElement('button');
  logoutBtn.textContent = 'Log Out';
  logoutBtn.style.cssText =
    'background:transparent;color:#666;border:1px solid #ccc;border-radius:4px;' +
    'padding:10px 24px;font-size:14px;cursor:pointer';
  logoutBtn.onclick = () => {
    window.location.href = '/user/logout';
  };

  btnRow.appendChild(stayBtn);
  btnRow.appendChild(logoutBtn);
  box.appendChild(btnRow);
  overlay.appendChild(box);
  document.body.appendChild(overlay);
}

function dismissWarning() {
  warningShown = false;
  const el = document.getElementById('meshery-session-warning');
  if (el) el.remove();
}

export function startSessionTimer() {
  if (typeof window === 'undefined') return;

  // Reset activity on any user interaction
  const resetOnInteraction = () => recordActivity();
  window.addEventListener('click', resetOnInteraction, { passive: true });
  window.addEventListener('keydown', resetOnInteraction, { passive: true });

  scheduleWarning();
}
