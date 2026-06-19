// Origin-allowlist guard for the post-authentication `/auth/redirect` bounce.
//
// The remote-provider login flow lands the browser on the provider-ui
// `/auth/redirect` page with a `return_to` URL and the freshly minted provider
// `token` in the query string; the page then forwards the browser — carrying
// that token — to `return_to`. Because the page is a public, statically
// exported endpoint, anyone can craft `.../auth/redirect?return_to=<any-url>`,
// so an unguarded forward is an open redirect that hands the provider token to
// an arbitrary host (and lends the provider's own hostname to phishing, or —
// via a `javascript:`/`data:` `return_to` — enables script injection).
//
// This mirrors the trust model meshery-cloud already enforces for credentialed
// CORS (server/router/cors.go `corsOriginAllower`): when a credential rides on
// the request the decision MUST be an exact host match against a trusted
// allowlist, never a wildcard. The trusted set is the provider roster
// propagated from install/providers.env (ALLOWED_RETURN_TO_HOSTS) plus the
// page's own origin and the configured PROVIDER_URL. Anything else falls back
// to PROVIDER_URL, so the token only ever reaches a trusted host.
//
// Note: meshery-cloud additionally authorises tenant *custom* domains
// dynamically from its org-domain registry (the registry-driven half of
// cors.go). A statically exported page cannot perform that lookup, so a custom
// domain that is not in the propagated roster falls back to PROVIDER_URL.
// Closing that gap needs a provider-side host-validation endpoint the page can
// call; that is tracked as a follow-up.

const LOCALHOST_HOSTNAMES = new Set(["localhost", "127.0.0.1", "[::1]", "::1"]);

// Localhost (and *.localhost) is the only context where a non-HTTPS return_to
// is honoured, for the `npm run dev` node-server flow.
const isLocalhostHostname = (hostname) =>
  LOCALHOST_HOSTNAMES.has(hostname) || hostname.endsWith(".localhost");

// Lower-cased hostname of an absolute URL, or "" when it cannot be parsed.
const hostnameOf = (value) => {
  try {
    return new URL(value).hostname.toLowerCase();
  } catch {
    return "";
  }
};

// Build the set of lower-cased hostnames the bounce may forward to: the
// propagated provider roster, plus the configured provider and the page's own
// origin (redirecting to yourself never crosses an origin boundary).
export const buildTrustedHosts = ({ providerUrl, currentOrigin, allowedHosts = [] }) => {
  const hosts = new Set();
  for (const host of allowedHosts) {
    if (host) {
      hosts.add(String(host).toLowerCase());
    }
  }
  const providerHost = hostnameOf(providerUrl);
  if (providerHost) {
    hosts.add(providerHost);
  }
  const ownHost = hostnameOf(currentOrigin);
  if (ownHost) {
    hosts.add(ownHost);
  }
  return hosts;
};

// A return_to may receive the forwarded token only when it is an absolute HTTPS
// URL (HTTP allowed for localhost dev) whose host is exactly trusted. The
// scheme check also rejects non-navigational schemes such as
// `javascript:`/`data:`, and exact host matching rejects look-alike subdomains.
export const isReturnToTrusted = (returnTo, trustedHosts) => {
  let url;
  try {
    url = new URL(returnTo);
  } catch {
    return false;
  }
  const hostname = url.hostname.toLowerCase();
  const isHttps = url.protocol === "https:";
  const isLocalDevHttp = url.protocol === "http:" && isLocalhostHostname(hostname);
  if (!isHttps && !isLocalDevHttp) {
    return false;
  }
  return trustedHosts.has(hostname);
};

// Resolve the absolute URL the `/auth/redirect` page should navigate to.
// Forwards to `return_to` (carrying the existing query params) only when it is
// a trusted host; otherwise lands on the provider's canonical host so the token
// is never handed to an untrusted origin.
export const resolveSafeRedirectURL = ({ search, providerUrl, currentOrigin, allowedHosts }) => {
  const params = new URLSearchParams(search);
  const returnTo = params.get("return_to");
  const trustedHosts = buildTrustedHosts({ providerUrl, currentOrigin, allowedHosts });

  // Merge the incoming params into the destination via the URL API (rather than
  // string concatenation) so a return_to that already carries a query string
  // stays well-formed instead of producing a double `?`.
  let destination;
  if (returnTo && isReturnToTrusted(returnTo, trustedHosts)) {
    destination = new URL(returnTo);
  } else {
    // Untrusted or absent return_to: land on the provider's canonical root so
    // the token is never handed to an untrusted origin.
    destination = new URL(providerUrl);
    if (!destination.pathname.endsWith("/")) {
      destination.pathname += "/";
    }
  }
  for (const [key, value] of params) {
    destination.searchParams.set(key, value);
  }
  return destination.toString();
};
