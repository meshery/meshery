let sessionExpiredShown = false;

function showSessionExpiredAndRedirect() {
  // Only show once even if multiple requests fail simultaneously
  if (sessionExpiredShown) return;
  sessionExpiredShown = true;

  const redirectTo = window.location.host.endsWith('3000') ? '/user/login' : window.location.href;

  // If the document isn't interactive yet, redirect immediately
  if (!document.body) {
    window.location = redirectTo;
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
    window.location = redirectTo;
  };
  box.appendChild(btn);
  overlay.appendChild(box);
  document.body.appendChild(overlay);

  // Auto-redirect after 5 seconds if user doesn't click
  setTimeout(() => {
    window.location = redirectTo;
  }, 5000);
}

import { recordActivity } from './sessionTimer';

const dataFetch = (url, options = {}, successFn, errorFn) => {
  if (errorFn === undefined) {
    errorFn = (err) => {
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
          } catch (e) {
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
      if (e.then) {
        e.then((text) => errorFn(text));
        return;
      }
      errorFn(e);
    });
};

/**
 * promisifiedDataFetch adds a promise wrapper to the dataFetch function
 * and ideal for use inside async functions - which is most of the functions
 * @param {string} url url is the endpoint
 * @param {Record<string, any>} options HTTP request options
 * @returns
 */
export function promisifiedDataFetch(url, options = {}) {
  return new Promise((resolve, reject) => {
    dataFetch(
      url,
      options,
      (result) => resolve(result),
      (err) => reject(err),
    );
  });
}

export default dataFetch;
