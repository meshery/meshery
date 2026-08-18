#!/usr/bin/env node

// bats-to-allure.js — convert a BATS TAP report into Allure result JSON files.
//
// Beyond the raw TAP → Allure mapping this converter is the single injection
// point for the test-plan traceability contract shared with the UI lane and the
// meshery/qa Test-Group-keyed reports (e.g. "Connection Lifecycle"). The Test
// Plan "Latest" tab columns are: col A = Test #, col B = Test Group,
// col C = Client, col D = Component Under Test.
//
//   testId            = TC-<n>   (Test Plan "Test #", col A)   [title token: [TC-<n>]]
//   testGroup         = <col B value>  (report bucket)          [title token: [tg=<value>]]
//   client            = <col C value>  (CLI here; this converter only ever sees mesheryctl results)
//   componentUnderTest = <col D value>                          [title token: [cut=<value>]]
//   epic              = "Kubernetes Connections"  (legacy report bucket) [token [epic=<v>] or derived from cut]
//   tag               = TC-<n>   (chip / filter passthrough of the testId)
//
// A BATS test opts in by prefixing its `@test` title with leading bracket
// tokens, e.g.
//   @test "[TC-1042][cut=Kubernetes Connection][tg=Connection Lifecycle] connection create ... creates a connection" { ... }
// The tokens are stripped from the displayed Allure name.
//
// `testGroup` is the general report key: any Test Group (col B) can become its
// own filtered report in meshery/qa by keying on this label. "Connection
// Lifecycle" is the first consumer; the mechanism is not connection-specific.

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

// ---------------------------------------------------------------------------
// Test Plan deep-link (traceability contract, part A) -----------------------
//
// Each connection result carries an Allure `tms` link back to its exact row in
// the Meshery Test Plan Google Sheet "Latest" tab, so a reviewer can click
// straight from a report test to its source case. The row is derived from the
// Test # (col A, i.e. the TC-<n> id): for the connection block,
//   ROW = TestNum - CONN_ROW_OFFSET
//
// !!!  CURRENT-LAYOUT-DEPENDENT - REGENERATE IF THE SHEET IS RE-SORTED  !!!
// The offset below encodes the CURRENT "Latest" tab layout, in which the
// connection cases Test# 1012..1089 occupy the contiguous rows 234..311
// (1012 - 778 = 234 ... 1089 - 778 = 311, verified against the live sheet).
// If the Latest tab is re-sorted, or rows are inserted above the connection
// block, both the offset AND the [CONN_TEST_MIN, CONN_TEST_MAX] range MUST be
// recomputed - otherwise the deep-link silently points at the wrong row. The
// range guard is deliberate: results outside the connection block get NO link
// rather than a wrong one (the offset is only known-good for that block).
const CONN_TEST_MIN = 1012;
const CONN_TEST_MAX = 1089;
const CONN_ROW_OFFSET = 778; // ROW = TestNum - 778 (Test# 1012 -> row 234).
const TEST_PLAN_SHEET_ID = "13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs";
const TEST_PLAN_GID = "838298230";

// testPlanLink builds the Allure `tms` link back to the Test Plan row for a
// connection Test # (`TC-<n>`), or returns null when the id is malformed or
// falls outside the connection block the offset is valid for. See the caveat
// above: the offset is layout-specific, so a non-connection id must never be
// turned into a (wrong) deep-link.
function testPlanLink(testId) {
  const m = /^TC-(\d+)$/i.exec((testId || "").trim());
  if (!m) return null;
  const testNum = Number(m[1]);
  if (testNum < CONN_TEST_MIN || testNum > CONN_TEST_MAX) return null;
  const row = testNum - CONN_ROW_OFFSET;
  const url =
    `https://docs.google.com/spreadsheets/d/${TEST_PLAN_SHEET_ID}` +
    `/edit?gid=${TEST_PLAN_GID}#gid=${TEST_PLAN_GID}&range=A${row}`;
  return { name: `Test Plan ${testId}`, url, type: "tms" };
}

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
  const links = [];
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
      // [TC-<n>] — the Test Plan Test #. Emit both testId and a tag, plus the
      // deep-link back to the Test Plan row (connection block only; see
      // testPlanLink).
      const testId = token;
      labels.push({ name: "testId", value: testId });
      labels.push({ name: "tag", value: testId });
      const link = testPlanLink(testId);
      if (link) links.push(link);
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
      case "tg":
        // [tg=<Test Group>] — Test Plan Test Group (col B). The general report
        // key: any Test Group can drive its own filtered meshery/qa report.
        labels.push({ name: "testGroup", value });
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

  return { name: name.trim(), labels, links };
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

  const { name: cleanName, labels, links } = parseTitleTokens(name, extraLabels);
  return { status, name: cleanName, skipReason, labels, links };
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

function createAllureResult({ name, status, start, stop, details, labels, links }) {
  const result = {
    uuid: uuid(),
    name,
    status,
    stage: "finished",
    start,
    stop,
    statusDetails: details ? { message: details } : {},
    labels,
    // _diag accumulates the TAP `#` diagnostic lines that follow this test (its
    // captured command transcript on failure). It is a scratch field consumed
    // by finalizeDiagnostics() and deleted before the result JSON is written.
    _diag: []
  };
  if (links && links.length) result.links = links;
  return result;
}

// summarizeDiagnostics picks a one-line headline for statusDetails.message from
// the captured diagnostic lines: prefer the assertion "...' failed" line (the
// most informative), falling back to the first non-empty line.
function summarizeDiagnostics(diagLines) {
  const failLine = diagLines.find(l => /\bfailed\b/.test(l));
  return (failLine || diagLines.find(l => l.trim()) || "").trim();
}

// finalizeDiagnostics turns a failed result's captured BATS output (part B) into
// debuggable Allure detail: a concise headline (message), the full transcript
// (trace), AND a text attachment holding the same transcript (the command +
// captured $output/$stderr surfaced by `bats --print-output-on-failure`). It
// writes the attachment file into ALLURE_RESULTS_DIR and strips the scratch
// `_diag` field. Skipped results keep their skip-reason message untouched;
// passed results carry no diagnostics.
function finalizeDiagnostics(result) {
  const diag = result._diag || [];
  delete result._diag;

  if (result.status !== "failed" || diag.length === 0) return;

  const transcript = diag.join("\n");
  result.statusDetails.message =
    summarizeDiagnostics(diag) || result.statusDetails.message || "test failed";
  result.statusDetails.trace = transcript;

  const attachmentFile = `${result.uuid}-attachment.txt`;
  fs.writeFileSync(path.join(ALLURE_RESULTS_DIR, attachmentFile), transcript);
  result.attachments = [
    { name: "CLI output (bats)", source: attachmentFile, type: "text/plain" }
  ];
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
        labels: dedupeLabels([...baseLabels(extraLabels), ...parsed.labels]),
        links: parsed.links
      });

      results.push(result);
      current = result;
      continue;
    }

    // Diagnostic lines describe the test that PRECEDES them (a failing test's
    // stack + captured output in BATS TAP), so accumulate them on the current
    // result. finalizeDiagnostics() later turns them into message/trace/an
    // attachment for failed tests.
    const diagMatch = trimmed.match(TAP_DIAGNOSTIC_RE);
    if (diagMatch && current) {
      current._diag.push(diagMatch[1]);
    }
  }

  for (const result of results) {
    finalizeDiagnostics(result);
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
  testPlanLink,
  summarizeDiagnostics,
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
