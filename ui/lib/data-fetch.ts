let sessionExpiredShown = false;

function showSessionExpiredAndRedirect(): void {
  // Only show once even if multiple requests fail simultaneously
  if (sessionExpiredShown) return;
  sessionExpiredShown = true;

  const redirectTo = window.location.host.endsWith('3000') ? '/user/login' : window.location.href;

  // If the document isn't interactive yet, redirect immediately
  if (!document.body) {
    window.location.href = redirectTo;
    return;
  }

  const overlay = document.createElement('div');
  overlay.setAttribute('role', 'dialog');
  overlay.setAttribute('aria-modal', 'true');
  overlay.style.cssText =
    'position:fixed;inset:0;z-index:99999;display:flex;align-items:center;' +
    'justify-content:center;background:rgba(0,0,0,0.5)';

  const box = document.createElement('div');
  box.style.cssText =
    'background:#fff;border-radius:8px;padding:32px;max-width:400px;text-align:center;' +
    'font-family:system-ui,sans-serif;color:#333;box-shadow:0 4px 24px rgba(0,0,0,0.2)';
  box.innerHTML =
    '<h2 style="margin:0 0 12px">Session Expired</h2>' +
    '<p style="margin:0 0 24px;color:#666">Your session has expired. You will be redirected to log in.</p>';

  const btn = document.createElement('button');
  btn.textContent = 'Log In';
  btn.style.cssText =
    'background:#477e96;color:#fff;border:none;border-radius:4px;padding:10px 32px;' +
    'font-size:14px;cursor:pointer';
  btn.onclick = () => {
    window.location.href = redirectTo;
  };
  box.appendChild(btn);
  overlay.appendChild(box);
  document.body.appendChild(overlay);

  // Auto-redirect after 5 seconds if user doesn't click
  setTimeout(() => {
    window.location.href = redirectTo;
  }, 5000);
}

import { recordActivity } from './sessionTimer';

const dataFetch = (
  url: string,
  options: RequestInit = {},
  successFn?: (data: unknown) => void,
  errorFn?: (err: unknown) => void
): void => {
  if (errorFn === undefined) {
    errorFn = (err: unknown) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    };
  }
  fetch(url, options)
    .then((res) => {
      if (res.status === 401 || res.redirected) {
        showSessionExpiredAndRedirect();
        return new Promise(() => {});
      }

      // Successful response — session is alive, reset the timeout warning
      recordActivity();

      let result;
      if (res.ok) {
        result = res.text().then((text) => {
          try {
            return JSON.parse(text);
          } catch {
            return text;
          }
        });

        return result;
      } else {
        throw res.text();
      }
    })
    .then(successFn)
    .catch((e) => {
      if (e && typeof e.then === 'function') {
        e.then((text: unknown) => {
          if (errorFn) errorFn(text);
        });
        return;
      }
      if (errorFn) errorFn(e);
    });
};

/**
 * promisifiedDataFetch adds a promise wrapper to the dataFetch function
 * and ideal for use inside async functions - which is most of the functions
 * @param {string} url url is the endpoint
 * @param {RequestInit} options HTTP request options
 * @returns
 */
export function promisifiedDataFetch<T = unknown>(url: string, options: RequestInit = {}): Promise<T> {
  return new Promise((resolve, reject) => {
    dataFetch(
      url,
      options,
      (result) => resolve(result as T),
      (err) => reject(err),
    );
  });
}

export default dataFetch;
