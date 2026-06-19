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

const trusted = buildTrustedHosts({
  providerUrl: PROVIDER_URL,
  currentOrigin: CURRENT_ORIGIN,
  allowedHosts: ALLOWED_HOSTS,
});

const resolve = (search) =>
  resolveSafeRedirectURL({
    search,
    providerUrl: PROVIDER_URL,
    currentOrigin: CURRENT_ORIGIN,
    allowedHosts: ALLOWED_HOSTS,
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

test("forwards to the configured provider host", () => {
  const out = resolve(
    "?return_to=https%3A%2F%2Fcloud.meshery.io%2Fauth%2Fredirect%2Faccept&token=t",
  );
  assert.equal(new URL(out).hostname, "cloud.meshery.io");
});

test("forwards to the same (dev) origin over http", () => {
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
  assert.equal(isReturnToTrusted("javascript:alert(document.cookie)", trusted), false);
  assert.equal(isReturnToTrusted("data:text/html,<script>alert(1)</script>", trusted), false);
});

test("rejects an http downgrade for an otherwise-trusted host", () => {
  assert.equal(isReturnToTrusted("http://cloud.meshery.io/x", trusted), false);
});

test("rejects a malformed / relative return_to", () => {
  assert.equal(isReturnToTrusted("not-a-url", trusted), false);
  assert.equal(isReturnToTrusted("/auth/redirect/accept", trusted), false);
});

test("matches hosts case-insensitively", () => {
  assert.equal(isReturnToTrusted("https://CLOUD.MESHERY.IO/x", trusted), true);
});

test("rejects a look-alike subdomain of a trusted host (exact match only)", () => {
  assert.equal(isReturnToTrusted("https://evil.cloud.meshery.io/x", trusted), false);
  assert.equal(isReturnToTrusted("https://cloud.meshery.io.evil.com/x", trusted), false);
});
