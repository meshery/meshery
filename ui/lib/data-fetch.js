import fetch from 'isomorphic-unfetch';

const dataFetch = (url, options = {}, successFn, errorFn) => {
  // const controller = new AbortController();
  // const signal = controller.signal;
  // options.signal = signal;
  // setTimeout(() => controller.abort(), 10000); // nice to have but will mess with the load test
  if (errorFn === undefined) {
    errorFn = (err) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    };
  }
  fetch(url, options)
    .then((res) => {
      if (res.status === 401 || res.redirected) {
        if (window.location.host.endsWith('3000')) {
          window.location = '/user/login'; // for local dev thru node server
        } else {
          window.location.reload(); // for use with Go server
        }
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
