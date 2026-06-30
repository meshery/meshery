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
// the request the decision MUST be an exact match against a trusted allowlist,
// never a wildcard. The trusted set is the provider roster propagated from
// install/providers.env (ALLOWED_RETURN_TO_HOSTS) plus the configured
// PROVIDER_URL and the page's own origin. Matching is on `host` (hostname +
// any explicit port), not bare hostname, so a security origin's port is
// significant and different ports on the same domain are not conflated.
// Anything else falls back to PROVIDER_URL, so the token only ever reaches a
// trusted host.
//
// Note: meshery-cloud additionally authorises tenant *custom* domains
// dynamically from its org-domain registry (the registry-driven half of
// cors.go). A statically exported page cannot perform that lookup, so a custom
// domain that is not in the propagated roster falls back to PROVIDER_URL.
// Closing that gap needs a provider-side host-validation endpoint the page can
// call; that is tracked as a follow-up.

const LOCALHOST_HOSTNAMES = new Set(["localhost", "127.0.0.1", "[::1]", "::1"]);

// Localhost (and *.localhost) is one of the two contexts where a non-HTTPS
// return_to is honoured, for the `npm run dev` node-server flow.
const isLocalhostHostname = (hostname) =>
  LOCALHOST_HOSTNAMES.has(hostname) || hostname.endsWith(".localhost");

// Lower-cased host (hostname + any explicit, non-default port) of an absolute
// URL, or "" when it cannot be parsed. Port is part of a browser security
// origin, so the allowlist is matched on host, not bare hostname.
const hostOf = (value) => {
  try {
    return new URL(value).host.toLowerCase();
  } catch {
    return "";
  }
};

// True when `url` is the exact same origin (scheme + host + port) as the page.
const isSameOrigin = (url, currentOrigin) => {
  try {
    return new URL(currentOrigin).origin === url.origin;
  } catch {
    return false;
  }
};

// Build the set of lower-cased hosts the bounce may forward to over HTTPS: the
// propagated provider roster, plus the configured provider and the page's own
// host.
export const buildTrustedHosts = ({ providerUrl, currentOrigin, allowedHosts = [] }) => {
  const hosts = new Set();
  for (const host of allowedHosts) {
    if (host) {
      hosts.add(String(host).toLowerCase());
    }
  }
  const providerHost = hostOf(providerUrl);
  if (providerHost) {
    hosts.add(providerHost);
  }
  const ownHost = hostOf(currentOrigin);
  if (ownHost) {
    hosts.add(ownHost);
  }
  return hosts;
};

// Decide whether `returnTo` may receive the forwarded token.
//   - HTTPS: the host (hostname + port) must be exactly trusted. Exact matching
//     rejects look-alike subdomains and same-host-different-port.
//   - HTTP: honoured only for local dev / a self-hosted server reached over a
//     LAN/Docker IP — a localhost host, or an exact same-origin redirect back
//     to the page itself (same-origin never crosses a trust boundary, and the
//     full-origin check avoids an https->http downgrade to a roster host).
//   - Anything else (javascript:, data:, malformed, relative) is rejected.
export const isReturnToTrusted = (returnTo, { trustedHosts, currentOrigin }) => {
  let url;
  try {
    url = new URL(returnTo);
  } catch {
    return false;
  }
  if (url.protocol === "https:") {
    return trustedHosts.has(url.host.toLowerCase());
  }
  if (url.protocol === "http:") {
    return isLocalhostHostname(url.hostname.toLowerCase()) || isSameOrigin(url, currentOrigin);
  }
  return false;
};

// Resolve the absolute URL the `/auth/redirect` page should navigate to.
// Forwards to `return_to` only when it is trusted; otherwise lands on the
// provider's canonical host so the token is never handed to an untrusted origin.
export const resolveSafeRedirectURL = ({ search, providerUrl, currentOrigin, allowedHosts }) => {
  const params = new URLSearchParams(search);
  const returnTo = params.get("return_to");
  const trustedHosts = buildTrustedHosts({ providerUrl, currentOrigin, allowedHosts });

  let destination;
  if (returnTo && isReturnToTrusted(returnTo, { trustedHosts, currentOrigin })) {
    destination = new URL(returnTo);
  } else {
    // Untrusted or absent return_to: land on the provider's canonical root so
    // the token is never handed to an untrusted origin.
    destination = new URL(providerUrl);
    if (!destination.pathname.endsWith("/")) {
      destination.pathname += "/";
    }
  }

  // Merge the incoming params via the URL API (so a return_to that already
  // carries a query string stays well-formed), but never forward `return_to`
  // itself: the destination already IS that target, so copying it only bloats
  // the URL and — because validation reads the first return_to while a
  // parameter-polluted request can carry a second, untrusted one — could hand
  // an unvalidated URL to the next hop. Dropping it closes that bypass.
  for (const [key, value] of params) {
    if (key !== "return_to") {
      destination.searchParams.set(key, value);
    }
  }
  return destination.toString();
};
