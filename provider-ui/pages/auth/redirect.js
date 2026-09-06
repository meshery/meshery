import { useEffect } from "react";
import { ALLOWED_RETURN_TO_HOSTS, PROVIDER_URL } from "../../lib/data-fetch";
import { resolveSafeRedirectURL } from "../../lib/safe-redirect.mjs";

const Redirect = () => {
  useEffect(() => {
    // Forward the browser (carrying the provider token in the query string)
    // only to a trusted host. An unguarded forward here is an open redirect
    // that leaks the token to an arbitrary origin; resolveSafeRedirectURL
    // falls back to PROVIDER_URL for any untrusted return_to. See
    // ../../lib/safe-redirect.mjs.
    window.location.href = resolveSafeRedirectURL({
      search: window.location.search,
      providerUrl: PROVIDER_URL,
      currentOrigin: window.location.origin,
      allowedHosts: ALLOWED_RETURN_TO_HOSTS,
    });
  }, []);
  return <></>;
};

export default Redirect;
