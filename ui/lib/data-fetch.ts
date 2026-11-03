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

/**
 * The error passed to `errorFn` (and rejected by `promisifiedDataFetch`) when a
 * request fails. `message` is always human-readable; `status` is the HTTP status
 * (undefined when no response arrived) and `body` is the raw response text.
 */
export class DataFetchError extends Error {
  url;
  status;
  body;

  constructor(message, { url, status = undefined, body = undefined, cause = undefined }) {
    super(message, { cause });
    this.name = 'DataFetchError';
    this.url = url;
    this.status = status;
    this.body = body;
  }
}

// JSON error bodies carry their text in a field: `{"error": "..."}`, or
// `message` in a MeshKit envelope. Fall back to the raw text otherwise.
function messageFromBody(body) {
  try {
    const parsed = JSON.parse(body);
    // `message` first, but skip a blank or non-string one so `error` still counts.
    for (const field of [parsed?.message, parsed?.error]) {
      if (typeof field === 'string' && field.trim() !== '') return field;
    }
  } catch {
    // Not JSON; the text itself is the message.
  }
  return body;
}

function errorFromResponse(url, res, body) {
  // An empty body or an HTML page (e.g. a proxy's 502/504) is not a useful
  // message; describe the status instead and keep the raw text on `body`.
  const readable = body.trim() !== '' && !/^\s*</.test(body);
  const status = res.statusText ? `${res.status} ${res.statusText}` : `${res.status}`;
  const message = readable
    ? messageFromBody(body)
    : `Request to ${url} failed with status ${status}`;
  return new DataFetchError(message, { url, status: res.status, body });
}

/**
 * @param {string} url
 * @param {RequestInit & { timeout?: number }} options `timeout` (ms) is opt-in;
 *   without it a request waits as long as the browser allows.
 * @param {Function} successFn called with the parsed JSON, or the text if it is not JSON
 * @param {Function} errorFn called with a DataFetchError, or the caller's own AbortError
 */
const dataFetch = (url, options = {}, successFn, errorFn) => {
  if (errorFn === undefined) {
    errorFn = (err) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    };
  }

  const { timeout, ...fetchOptions } = options as RequestInit & { timeout?: number };
  let timer;
  let timedOut = false;
  let detachCallerSignal = () => {};
  if (timeout > 0) {
    const controller = new AbortController();
    const callerSignal = fetchOptions.signal;
    if (callerSignal?.aborted) {
      controller.abort(callerSignal.reason);
    } else {
      // Detached once the request settles, so a long-lived signal reused across
      // requests doesn't keep every finished request's controller alive.
      const onCallerAbort = () => controller.abort(callerSignal.reason);
      callerSignal?.addEventListener('abort', onCallerAbort, { once: true });
      detachCallerSignal = () => callerSignal?.removeEventListener('abort', onCallerAbort);
    }
    timer = setTimeout(() => {
      timedOut = true;
      controller.abort();
    }, timeout);
    fetchOptions.signal = controller.signal;
  }

  const toDataFetchError = (err) => {
    if (err instanceof DataFetchError) return err;
    if (err?.name === 'AbortError') {
      // An abort the caller asked for is passed through untouched.
      return timedOut
        ? new DataFetchError(`Request to ${url} timed out after ${timeout}ms`, { url, cause: err })
        : err;
    }
    if (err instanceof TypeError) {
      return new DataFetchError(`Unable to reach the Meshery server (${url}): ${err.message}`, {
        url,
        cause: err,
      });
    }
    return err;
  };

  fetch(url, fetchOptions)
    .then((res) => {
      if (res.status === 401 || res.redirected) {
        clearTimeout(timer);
        detachCallerSignal();
        showSessionExpiredAndRedirect();
        return new Promise(() => {});
      }

      // Successful response — session is alive, reset the timeout warning
      recordActivity();

      if (!res.ok) {
        return res.text().then(
          (text) => {
            throw errorFromResponse(url, res, text ?? '');
          },
          (err) => {
            // An abort mid-body (timeout or the caller's signal) must stay an
            // abort, so toDataFetchError can report it as such.
            if (err?.name === 'AbortError') throw err;
            throw errorFromResponse(url, res, '');
          },
        );
      }

      return res.text().then((text) => {
        try {
          return JSON.parse(text);
        } catch (e) {
          return text;
        }
      });
    })
    .catch((err) => {
      throw toDataFetchError(err);
    })
    .then(successFn)
    .catch((err) => errorFn(err))
    .finally(() => {
      clearTimeout(timer);
      detachCallerSignal();
    });
};

/**
 * promisifiedDataFetch adds a promise wrapper to the dataFetch function
 * and ideal for use inside async functions - which is most of the functions
 * @param {string} url url is the endpoint
 * @param {Record<string, any>} options HTTP request options, plus an optional `timeout` in ms
 * @returns {Promise<any>} rejects with a DataFetchError (or the caller's own AbortError)
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
