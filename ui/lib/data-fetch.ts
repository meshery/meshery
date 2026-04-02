// Store reference for dispatching SESSION_EXPIRED on 401.
// Lazy-initialized to avoid circular dependency (data-fetch -> store -> rtk-query -> data-fetch).
let _store: { dispatch: (action: unknown) => void } | null = null;

export function setDataFetchStore(store: { dispatch: (action: unknown) => void }) {
  _store = store;
}

const dataFetch = (url, options = {}, successFn, errorFn) => {
  if (errorFn === undefined) {
    errorFn = (err) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    };
  }
  fetch(url, options)
    .then((res) => {
      if (res.status === 401 || res.redirected) {
        if (_store) {
          _store.dispatch({ type: 'SESSION_EXPIRED' });
        }
        return new Promise(() => {});
      }

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
