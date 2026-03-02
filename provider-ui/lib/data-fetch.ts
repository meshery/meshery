import fetch from 'isomorphic-unfetch';

export const PROVIDER_URL = 'https://cloud.layer5.io';

type SuccessCallback<T> = (_data: T) => void;
type ErrorCallback = (_errorValue: string | Error) => void;

// This can be migrated as a custom hook in React
const dataFetch = <T = unknown>(
  url: string,
  options: globalThis.RequestInit = {},
  successFn?: SuccessCallback<T>,
  errorFn?: ErrorCallback
): void => {
  const handleError: ErrorCallback =
    errorFn ||
    ((err) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    });

  fetch(url, options)
    .then((res) => {
      if (res.status === 401 || res.redirected) {
        if (window.location.host.endsWith('3000')) {
          window.location.href = '/user/login'; // for local dev thru node server
        } else {
          window.location.reload(); // for use with Go server
        }
      }
      if (res.ok) {
        try {
          return res.json() as Promise<T>;
        } catch {
          return res.text() as unknown as Promise<T>;
        }
      } else {
        res.text().then(handleError);
        return undefined;
      }
    })
    .then((result) => {
      if (result !== undefined && successFn) {
        successFn(result);
      }
    })
    .catch(handleError);
};

export default dataFetch;
