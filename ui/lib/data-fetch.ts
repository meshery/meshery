let sessionExpiredShown = false;

function showSessionExpiredAndRedirect() {
  // Only show once even if multiple requests fail simultaneously
  if (sessionExpiredShown) return;
  sessionExpiredShown = true;
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new Event('meshery:session-expired'));
  }
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
