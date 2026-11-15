export const PROVIDER_URL = "https://cloud.meshery.io"; // AUTO-GENERATED from install/providers.env - run `make providers-propagate`

// Ports of the local Node dev servers that proxy /user/login to the Meshery
// server: the main UI's (ui/ui_dev_server.js, 3000) and provider-ui's own
// (ui_dev_server.js, 3001). Behind the Go server a reload re-runs its auth
// redirect instead.
const NODE_DEV_SERVER_PORTS = ["3000", "3001"];

// This can be migrated as a custom hook in React
const dataFetch = (url, options = {}, successFn, errorFn) => {
  if (errorFn === undefined) {
    errorFn = (err) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    };
  }
  fetch(url, options)
    .then(async (res) => {
      if (res.status === 401 || res.redirected) {
        if (NODE_DEV_SERVER_PORTS.includes(window.location.port)) {
          window.location = "/user/login";
        } else {
          window.location.reload();
        }
        return;
      }
      // Read the body once as text: res.json() rejects asynchronously, so a
      // try/catch around it never reached the text fallback.
      const body = await res.text();
      if (!res.ok) {
        errorFn(body);
        return;
      }
      let result;
      try {
        result = JSON.parse(body);
      } catch {
        result = body;
      }
      successFn(result);
    })
    .catch(errorFn);
};

export default dataFetch;
