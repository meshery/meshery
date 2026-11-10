export const PROVIDER_URL = "https://cloud.meshery.io"; // AUTO-GENERATED from install/providers.env - run `make providers-propagate`

// Exact-match hosts the post-auth `/auth/redirect` bounce may forward the
// browser (and the provider token) to. Propagated from the active provider
// roster in install/providers.env; enforced by lib/safe-redirect.mjs.
export const ALLOWED_RETURN_TO_HOSTS = [
  // BEGIN AUTO-GENERATED from install/providers.env - run `make providers-propagate`
  "cloud.meshery.io",
  "perf.smp-spec.io",
  "cloud.layer5.io",
  "platform.tata-consulting.co.uk",
  "collab.eti.cisco.com",
  "kickstart.metabit.com",
  "provider.od10.in",
  // END AUTO-GENERATED
];

// This can be migrated as a custom hook in React
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
        if (window.location.host.endsWith("3000")) {
          window.location = "/user/login"; // for local dev thru node server
        } else {
          window.location.reload(); // for use with Go server
        }
      }
      let result;
      if (res.ok) {
        // console.log(`res type: ${res.type}`);
        try {
          result = res.json();
        } catch {
          result = res.text();
        }
        return result;
      } else {
        res.text().then(errorFn);
      }
    })
    .then(successFn)
    .catch(errorFn);
};

export default dataFetch;
