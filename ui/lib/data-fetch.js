import fetch from "isomorphic-unfetch";

const dataFetch = (url, options = {}, successFn, errorFn) => {
  fetch(url, options)
    .then((res) => {
      if (res.status === 401 || res.redirected) {
        if (window.location.host.endsWith("3000")) {
          window.location = "/user/login"; // for local dev thru node server
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
        throw new Error(res.statusText);
      }
    })
    .then(successFn)
    .catch(errorFn);
};

export function promisifiedDataFetch(url, options = {}) {
  return new Promise((resolve, reject) => {
    dataFetch(
      url,
      options,
      (result) => resolve(result),
      (err) => reject(err)
    );
  });
}

export default dataFetch;
