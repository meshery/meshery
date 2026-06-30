// Unit tests for the /auth/redirect origin-allowlist guard.
// Run with `npm test` (provider-ui) or `node --test lib/safe-redirect.test.mjs`.

import { test } from "node:test";
import assert from "node:assert/strict";

import {
  buildTrustedHosts,
  isReturnToTrusted,
  resolveSafeRedirectURL,
} from "./safe-redirect.mjs";

const PROVIDER_URL = "https://cloud.meshery.io";
const CURRENT_ORIGIN = "https://kanvas.new";
const ALLOWED_HOSTS = ["cloud.meshery.io", "cloud.layer5.io", "perf.smp-spec.io"];

const resolve = (search, currentOrigin = CURRENT_ORIGIN) =>
  resolveSafeRedirectURL({
    search,
    providerUrl: PROVIDER_URL,
    currentOrigin,
    allowedHosts: ALLOWED_HOSTS,
  });

const isTrusted = (returnTo, currentOrigin = CURRENT_ORIGIN) =>
  isReturnToTrusted(returnTo, {
    trustedHosts: buildTrustedHosts({
      providerUrl: PROVIDER_URL,
      currentOrigin,
      allowedHosts: ALLOWED_HOSTS,
    }),
    currentOrigin,
  });

test("forwards to a trusted roster host, preserving query params", () => {
  const out = resolve(
    "?return_to=https%3A%2F%2Fcloud.layer5.io%2Fauth%2Fredirect%2Faccept&token=abc&ref=%2Fdashboard",
  );
  const url = new URL(out);
  assert.equal(url.hostname, "cloud.layer5.io");
  assert.equal(url.pathname, "/auth/redirect/accept");
  assert.equal(url.searchParams.get("token"), "abc");
  assert.equal(url.searchParams.get("ref"), "/dashboard");
});

test("merges params when return_to already carries a query string", () => {
  const out = resolve(
    "?return_to=https%3A%2F%2Fcloud.layer5.io%2Fauth%2Fredirect%2Faccept%3Fexisting%3D1&token=abc&ref=%2Fdashboard",
  );
  const url = new URL(out);
  assert.equal(url.hostname, "cloud.layer5.io");
  assert.equal(url.pathname, "/auth/redirect/accept");
  // Pre-existing query on return_to is preserved, not clobbered into a double `?`.
  assert.equal(url.searchParams.get("existing"), "1");
  assert.equal(url.searchParams.get("token"), "abc");
  assert.equal(url.searchParams.get("ref"), "/dashboard");
});

test("never forwards the return_to param to the destination (no bloat)", () => {
  const out = resolve(
    "?return_to=https%3A%2F%2Fcloud.layer5.io%2Fauth%2Fredirect%2Faccept&token=t",
  );
  assert.equal(new URL(out).searchParams.has("return_to"), false);
});

test("parameter pollution: a second, untrusted return_to is never forwarded", () => {
  // params.get("return_to") validates the first (trusted) value; the second
  // (evil) one must not leak onto the destination query string.
  const out = resolve(
    "?return_to=https%3A%2F%2Fcloud.layer5.io%2Fauth%2Fredirect%2Faccept&return_to=https%3A%2F%2Fevil.example.com&token=secret",
  );
  const url = new URL(out);
  assert.equal(url.hostname, "cloud.layer5.io");
  assert.equal(url.searchParams.has("return_to"), false);
  assert.equal(out.includes("evil.example.com"), false);
});

test("forwards to the configured provider host", () => {
  const out = resolve(
    "?return_to=https%3A%2F%2Fcloud.meshery.io%2Fauth%2Fredirect%2Faccept&token=t",
  );
  assert.equal(new URL(out).hostname, "cloud.meshery.io");
});

test("forwards to a same-origin http dev server reached by a LAN/Docker IP", () => {
  const origin = "http://192.168.1.40:3000";
  const out = resolve(
    "?return_to=http%3A%2F%2F192.168.1.40%3A3000%2Fauth%2Fredirect%2Faccept&token=t",
    origin,
  );
  assert.equal(new URL(out).host, "192.168.1.40:3000");
});

test("forwards to localhost http in dev", () => {
  const out = resolveSafeRedirectURL({
    search: "?return_to=http%3A%2F%2Flocalhost%3A3000%2Fauth%2Fredirect%2Faccept&token=t",
    providerUrl: "http://localhost:9876",
    currentOrigin: "http://localhost:3000",
    allowedHosts: [],
  });
  assert.equal(new URL(out).host, "localhost:3000");
});

test("falls back to PROVIDER_URL for an untrusted host (open redirect blocked)", () => {
  const out = resolve("?return_to=https%3A%2F%2Fevil.example.com%2Fsteal&token=secret");
  const url = new URL(out);
  assert.equal(url.hostname, "cloud.meshery.io");
  // The token rides to the trusted provider, never to evil.example.com.
  assert.equal(url.searchParams.get("token"), "secret");
});

test("falls back to PROVIDER_URL when return_to is absent", () => {
  assert.equal(new URL(resolve("?token=t")).hostname, "cloud.meshery.io");
});

test("rejects javascript: and data: schemes", () => {
  assert.equal(isTrusted("javascript:alert(document.cookie)"), false);
  assert.equal(isTrusted("data:text/html,<script>alert(1)</script>"), false);
});

test("rejects an http downgrade for an otherwise-trusted host", () => {
  // http to a roster host that is neither localhost nor the page's own origin.
  assert.equal(isTrusted("http://cloud.meshery.io/x"), false);
});

test("rejects the same hostname on a different port (matches host, not hostname)", () => {
  assert.equal(isTrusted("https://cloud.meshery.io:8443/x"), false);
});

test("rejects a malformed / relative return_to", () => {
  assert.equal(isTrusted("not-a-url"), false);
  assert.equal(isTrusted("/auth/redirect/accept"), false);
});

test("matches hosts case-insensitively", () => {
  assert.equal(isTrusted("https://CLOUD.MESHERY.IO/x"), true);
});

test("rejects a look-alike subdomain of a trusted host (exact match only)", () => {
  assert.equal(isTrusted("https://evil.cloud.meshery.io/x"), false);
  assert.equal(isTrusted("https://cloud.meshery.io.evil.com/x"), false);
});
