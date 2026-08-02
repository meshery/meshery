#!/usr/bin/env node

// bats-to-allure.js — convert a BATS TAP report into Allure result JSON files.
//
// Beyond the raw TAP → Allure mapping this converter is the single injection
// point for the Kubernetes Connection test-plan traceability contract shared
// with the UI lane and the meshery/qa "Kubernetes Connections" report:
//
//   testId            = TC-<n>   (Test Plan "Test #", col A)   [title token: [TC-<n>]]
//   componentUnderTest = <col C value>                          [title token: [cut=<value>]]
//   epic              = "Kubernetes Connections"  (report bucket) [token [epic=<v>] or derived from cut]
//   client            = "CLI"    (this converter only ever sees mesheryctl results)
//   tag               = TC-<n>   (chip / filter passthrough of the testId)
//
// A BATS test opts in by prefixing its `@test` title with leading bracket
// tokens, e.g.
//   @test "[TC-1042][cut=Kubernetes Connection] connection create ... creates a connection" { ... }
// The tokens are stripped from the displayed Allure name.

const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

const ALLURE_RESULTS_DIR = "allure-results";

// TAP regex patterns.
const TAP_TEST_RE = /^(ok|not ok)\s+\d+\s+(.*)/;
const TAP_DIAGNOSTIC_RE = /^\s*#\s?(.*)/;
// TAP directive appended to a test line, e.g. "... # skip minikube not installed".
const TAP_SKIP_DIRECTIVE_RE = /\s*#\s*skip\b\s*(.*)$/i;

// componentUnderTest values (Test Plan col C) that belong to the
// "Kubernetes Connections" report epic. Kept in lockstep with the qa
// allurerc.mjs "Kubernetes Connections" plugin filter (meshery/qa). When a
// test carries one of these as its cut and sets no explicit [epic=...], the
// converter buckets it under the connection epic so results group correctly.
const CONNECTION_COMPONENTS = new Set(["kubernetes connection"]);
const CONNECTION_EPIC = "Kubernetes Connections";

// Parse suite-wide labels from env: ALLURE_LABELS="key=value,key=value".
function parseExtraLabels() {
  const raw = process.env.ALLURE_LABELS;
  if (!raw) return [];

  return raw.split(",")
    .map(pair => pair.trim())
    .filter(Boolean)
    .map(pair => {
      const [name, value] = pair.split("=");
      if (!name || !value) return null;
      return { name, value };
    })
    .filter(Boolean);
}

function ensureAllureDir() {
  fs.mkdirSync(ALLURE_RESULTS_DIR, { recursive: true });
}

function uuid() {
  return crypto.randomUUID();
}

function tapStatusToAllure(status) {
  return status === "ok" ? "passed" : "failed";
}

// parseTitleTokens splits leading "[...]" tokens off a BATS test title and maps
// them to Allure labels per the traceability contract above. Returns the
// display name (tokens stripped) and the derived per-test labels. extraLabels
// are the suite-wide labels (from ALLURE_LABELS) that will also be attached, so
// the cut->epic derivation can defer to a suite-wide epic instead of emitting a
// duplicate.
function parseTitleTokens(rawName, extraLabels = []) {
  const labels = [];
  let name = rawName;
  let cut = null;
  let epic = null;

  // Consume consecutive leading [ ... ] tokens.
  const tokenRe = /^\s*\[([^\]]*)\]/;
  let match;
  while ((match = tokenRe.exec(name)) !== null) {
    const token = match[1].trim();
    name = name.slice(match[0].length);

    if (/^TC-/i.test(token)) {
      // [TC-<n>] — the Test Plan Test #. Emit both testId and a tag.
      const testId = token;
      labels.push({ name: "testId", value: testId });
      labels.push({ name: "tag", value: testId });
      continue;
    }

    const eq = token.indexOf("=");
    if (eq === -1) {
      // Unknown bare token — surface it as a tag rather than dropping it.
      if (token) labels.push({ name: "tag", value: token });
      continue;
    }

    const key = token.slice(0, eq).trim();
    const value = token.slice(eq + 1).trim();
    if (!key || !value) continue;

    switch (key.toLowerCase()) {
      case "cut":
        cut = value;
        labels.push({ name: "componentUnderTest", value });
        break;
      case "epic":
        epic = value;
        labels.push({ name: "epic", value });
        break;
      case "client":
        labels.push({ name: "client", value });
        break;
      default:
        labels.push({ name: key, value });
    }
  }

  // Derive the connection epic from cut only when no epic is set by any source
  // — neither a title [epic=...] token nor a suite-wide ALLURE_LABELS epic —
  // mirroring the qa "Kubernetes Connections" plugin's componentUnderTest
  // fallback. Emitting a second epic here would file the test under two epics.
  const epicAlreadyPresent = epic !== null || extraLabels.some(l => l.name === "epic");
  if (!epicAlreadyPresent && cut && CONNECTION_COMPONENTS.has(cut.toLowerCase())) {
    labels.push({ name: "epic", value: CONNECTION_EPIC });
  }

  return { name: name.trim(), labels };
}

// parseTestLine turns a single TAP "ok/not ok" line into a normalized record:
// status (passed|failed|skipped), the clean display name, any skip reason, and
// the per-test traceability labels parsed from the title tokens.
function parseTestLine(rawStatus, rawName, extraLabels = []) {
  let status = tapStatusToAllure(rawStatus);
  let name = rawName;
  let skipReason = null;

  const skipMatch = name.match(TAP_SKIP_DIRECTIVE_RE);
  if (skipMatch) {
    // A skipped BATS test is emitted by TAP as "ok N <name> # skip <reason>".
    // Without this it was reported as "passed" — masking never-run tests.
    status = "skipped";
    skipReason = skipMatch[1].trim();
    name = name.slice(0, skipMatch.index);
  }

  const { name: cleanName, labels } = parseTitleTokens(name, extraLabels);
  return { status, name: cleanName, skipReason, labels };
}

function baseLabels(extraLabels) {
  return [
    { name: "framework", value: "bats" },
    { name: "language", value: "bash" },
    { name: "project", value: "mesheryctl" },
    // This converter only ever processes mesheryctl (CLI) results.
    { name: "client", value: "CLI" },
    ...extraLabels
  ];
}

// Label names that may legitimately appear more than once on one result.
// Everything else is single-valued and must not be duplicated (a duplicate
// e.g. `client` or `epic` label breaks downstream report filtering).
const MULTI_VALUED_LABELS = new Set(["tag"]);

// dedupeLabels collapses single-valued labels so the LAST occurrence wins,
// while preserving multi-valued labels (e.g. `tag`). Precedence therefore runs
// base -> ALLURE_LABELS -> per-test title tokens (most specific wins), which is
// why the base `client=CLI` is overridable by a `[client=...]` token.
function dedupeLabels(labels) {
  const indexByName = new Map();
  const out = [];
  for (const label of labels) {
    if (MULTI_VALUED_LABELS.has(label.name)) {
      out.push(label);
      continue;
    }
    if (indexByName.has(label.name)) {
      out[indexByName.get(label.name)] = label; // last value wins, in place
    } else {
      indexByName.set(label.name, out.length);
      out.push(label);
    }
  }
  return out;
}

function createAllureResult({ name, status, start, stop, details, labels }) {
  return {
    uuid: uuid(),
    name,
    status,
    stage: "finished",
    start,
    stop,
    statusDetails: details ? { message: details } : {},
    labels
  };
}

function convertTapToAllure(tapFile) {
  ensureAllureDir();

  const content = fs.readFileSync(tapFile, "utf-8");
  const lines = content.split("\n");

  const extraLabels = parseExtraLabels();

  const results = [];
  let current = null; // the most recently emitted result — trailing diagnostics attach here.
  let testCount = 0;

  for (const line of lines) {
    const trimmed = line.trimEnd();

    const testMatch = trimmed.match(TAP_TEST_RE);
    if (testMatch) {
      testCount++;
      const [, rawStatus, rawName] = testMatch;
      const parsed = parseTestLine(rawStatus, rawName, extraLabels);

      const start = Date.now();
      const result = createAllureResult({
        name: parsed.name,
        status: parsed.status,
        start,
        stop: start + 1,
        details: parsed.skipReason || null,
        labels: dedupeLabels([...baseLabels(extraLabels), ...parsed.labels])
      });

      results.push(result);
      current = result;
      continue;
    }

    // Diagnostic lines describe the test that PRECEDES them (e.g. a failing
    // test's stack in BATS TAP), so attach them to the current result.
    const diagMatch = trimmed.match(TAP_DIAGNOSTIC_RE);
    if (diagMatch && current) {
      const existing = current.statusDetails.message;
      current.statusDetails.message = existing
        ? `${existing}\n${diagMatch[1]}`
        : diagMatch[1];
    }
  }

  for (const result of results) {
    const outputPath = path.join(ALLURE_RESULTS_DIR, `${result.uuid}-result.json`);
    fs.writeFileSync(outputPath, JSON.stringify(result, null, 2));
  }

  console.log(`Converted ${testCount} BATS tests to Allure results.`);
  return results;
}

module.exports = {
  parseTitleTokens,
  parseTestLine,
  dedupeLabels,
  convertTapToAllure,
  CONNECTION_EPIC
};

// CLI
if (require.main === module) {
  if (process.argv.length !== 3) {
    console.error("Usage: node bats-to-allure.js <bats-report.tap>");
    process.exit(1);
  }
  convertTapToAllure(process.argv[2]);
}
