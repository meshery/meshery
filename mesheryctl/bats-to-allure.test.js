#!/usr/bin/env node

// Unit tests for bats-to-allure.js. Run with: node --test bats-to-allure.test.js
// Dependency-free: uses Node's built-in test runner + assert (Node >= 18).

const test = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

const {
  parseTitleTokens,
  parseTestLine,
  dedupeLabels,
  convertTapToAllure,
  testPlanLink,
  summarizeDiagnostics,
  CONNECTION_EPIC
} = require("./bats-to-allure.js");

function labelValue(labels, name) {
  const found = labels.filter(l => l.name === name).map(l => l.value);
  return found.length <= 1 ? found[0] : found;
}

test("skipped tests are reported as skipped, not passed", () => {
  const r = parseTestLine("ok", "create minikube connection # skip minikube not installed");
  assert.equal(r.status, "skipped");
  assert.equal(r.skipReason, "minikube not installed");
  assert.equal(r.name, "create minikube connection");
});

test("skip directive matching is case-insensitive and reason is optional", () => {
  const r = parseTestLine("ok", "some test # SKIP");
  assert.equal(r.status, "skipped");
  assert.equal(r.skipReason, "");
  assert.equal(r.name, "some test");
});

test("passing and failing statuses are preserved", () => {
  assert.equal(parseTestLine("ok", "a passes").status, "passed");
  assert.equal(parseTestLine("not ok", "c fails").status, "failed");
});

test("TC token yields testId + tag and is stripped from the name", () => {
  const { name, labels } = parseTitleTokens("[TC-1042] connection create works");
  assert.equal(name, "connection create works");
  assert.equal(labelValue(labels, "testId"), "TC-1042");
  assert.equal(labelValue(labels, "tag"), "TC-1042");
});

test("cut token maps to componentUnderTest and derives the connection epic", () => {
  const { name, labels } = parseTitleTokens("[TC-1042][cut=Kubernetes Connection] view works");
  assert.equal(name, "view works");
  assert.equal(labelValue(labels, "componentUnderTest"), "Kubernetes Connection");
  assert.equal(labelValue(labels, "epic"), CONNECTION_EPIC);
  assert.equal(labelValue(labels, "testId"), "TC-1042");
});

test("tg token maps to a testGroup label and is stripped from the name", () => {
  const { name, labels } = parseTitleTokens(
    "[TC-1042][cut=Kubernetes Connection][tg=Connection Lifecycle] view works"
  );
  assert.equal(name, "view works");
  assert.equal(labelValue(labels, "testGroup"), "Connection Lifecycle");
  assert.equal(labelValue(labels, "componentUnderTest"), "Kubernetes Connection");
  assert.equal(labelValue(labels, "testId"), "TC-1042");
});

test("tg token is independent of cut/epic and carries an arbitrary Test Group", () => {
  const { labels } = parseTitleTokens("[TC-2001][cut=Model][tg=Model Import] import works");
  assert.equal(labelValue(labels, "testGroup"), "Model Import");
  // A non-connection cut still gets no connection epic.
  assert.equal(labelValue(labels, "epic"), undefined);
});

test("explicit epic token overrides the derived epic", () => {
  const { labels } = parseTitleTokens("[cut=Kubernetes Connection][epic=Something Else] t");
  assert.equal(labelValue(labels, "epic"), "Something Else");
});

test("suite-wide epic (ALLURE_LABELS) suppresses the derived connection epic", () => {
  const { labels } = parseTitleTokens(
    "[cut=Kubernetes Connection] t",
    [{ name: "epic", value: "MesheryCtl" }]
  );
  // The suite-wide epic already supplies an epic; deriving would file the test
  // under two epics.
  assert.equal(labels.filter(l => l.name === "epic").length, 0);
});

test("non-connection cut does not get the connection epic", () => {
  const { labels } = parseTitleTokens("[TC-2001][cut=Model] import works");
  assert.equal(labelValue(labels, "componentUnderTest"), "Model");
  assert.equal(labelValue(labels, "epic"), undefined);
});

test("client token is honored when present", () => {
  const { labels } = parseTitleTokens("[client=UI] t");
  assert.equal(labelValue(labels, "client"), "UI");
});

test("untagged titles are unchanged and produce no per-test labels", () => {
  const { name, labels } = parseTitleTokens("connection list works");
  assert.equal(name, "connection list works");
  assert.equal(labels.length, 0);
});

test("skip directive and title tokens compose", () => {
  const r = parseTestLine("ok", "[TC-1044][cut=Kubernetes Connection] positive delete # skip no id");
  assert.equal(r.status, "skipped");
  assert.equal(r.skipReason, "no id");
  assert.equal(r.name, "positive delete");
  assert.equal(labelValue(r.labels, "testId"), "TC-1044");
});

test("ALLURE_LABELS epic is not duplicated by the derived connection epic", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "b2a-"));
  const cwd = process.cwd();
  const prevLabels = process.env.ALLURE_LABELS;
  try {
    process.chdir(tmp);
    process.env.ALLURE_LABELS = "epic=MesheryCtl,layer=cli";
    const tap = [
      "1..1",
      "ok 1 [TC-1041][cut=Kubernetes Connection] list works"
    ].join("\n");
    fs.writeFileSync(path.join(tmp, "sample.tap"), tap);

    const results = convertTapToAllure(path.join(tmp, "sample.tap"));
    assert.equal(results.length, 1);

    const epics = results[0].labels.filter(l => l.name === "epic");
    assert.equal(epics.length, 1);
    assert.equal(epics[0].value, "MesheryCtl");
  } finally {
    process.chdir(cwd);
    if (prevLabels === undefined) delete process.env.ALLURE_LABELS;
    else process.env.ALLURE_LABELS = prevLabels;
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test("dedupeLabels collapses single-valued labels (last wins) and keeps tags", () => {
  const out = dedupeLabels([
    { name: "client", value: "CLI" },
    { name: "epic", value: "MesheryCtl" },
    { name: "client", value: "UI" },   // token override wins
    { name: "tag", value: "TC-1" },
    { name: "tag", value: "TC-2" }      // multiple tags preserved
  ]);
  assert.deepEqual(out.filter(l => l.name === "client"), [{ name: "client", value: "UI" }]);
  assert.equal(out.filter(l => l.name === "epic").length, 1);
  assert.deepEqual(out.filter(l => l.name === "tag").map(l => l.value), ["TC-1", "TC-2"]);
});

test("a [client=...] token does not produce a duplicate client label", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "b2a-"));
  const cwd = process.cwd();
  try {
    process.chdir(tmp);
    fs.writeFileSync(path.join(tmp, "s.tap"), "1..1\nok 1 [client=UI] t\n");
    const [r] = convertTapToAllure(path.join(tmp, "s.tap"));
    const clients = r.labels.filter(l => l.name === "client");
    assert.equal(clients.length, 1);
    assert.equal(clients[0].value, "UI");
  } finally {
    process.chdir(cwd);
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test("convertTapToAllure end-to-end writes results with base + per-test labels", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "b2a-"));
  const cwd = process.cwd();
  try {
    process.chdir(tmp);
    const tap = [
      "1..3",
      "ok 1 [TC-1041][cut=Kubernetes Connection] list works",
      "ok 2 [TC-1042][cut=Kubernetes Connection] create # skip minikube not installed",
      "not ok 3 [TC-1043][cut=Kubernetes Connection] delete fails",
      "# (in test file x.bats, line 3)",
      "#   `false' failed"
    ].join("\n");
    fs.writeFileSync(path.join(tmp, "sample.tap"), tap);

    const results = convertTapToAllure(path.join(tmp, "sample.tap"));
    assert.equal(results.length, 3);

    const byId = Object.fromEntries(
      results.map(r => [r.labels.find(l => l.name === "testId").value, r])
    );

    assert.equal(byId["TC-1041"].status, "passed");
    assert.equal(byId["TC-1042"].status, "skipped");
    assert.equal(byId["TC-1042"].statusDetails.message, "minikube not installed");
    assert.equal(byId["TC-1043"].status, "failed");

    // Every result carries the base contract labels.
    for (const r of results) {
      assert.equal(r.labels.find(l => l.name === "client").value, "CLI");
      assert.equal(r.labels.find(l => l.name === "project").value, "mesheryctl");
      assert.equal(r.labels.find(l => l.name === "epic").value, CONNECTION_EPIC);
    }

    // Trailing diagnostics attach to the failing (preceding) test, not the next one.
    assert.match(byId["TC-1043"].statusDetails.message, /failed/);

    // Results are written to allure-results/.
    const files = fs.readdirSync(path.join(tmp, "allure-results"));
    assert.equal(files.filter(f => f.endsWith("-result.json")).length, 3);
  } finally {
    process.chdir(cwd);
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

// --- Part A: Test Plan deep-link -------------------------------------------

test("testPlanLink maps a connection Test # to its Latest-tab row (ROW = n - 778)", () => {
  const link = testPlanLink("TC-1012");
  assert.equal(link.type, "tms");
  assert.equal(link.name, "Test Plan TC-1012");
  // Test# 1012 -> row 234; the deep-link targets A234 on the Latest tab.
  assert.match(link.url, /range=A234$/);
  assert.match(link.url, /13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/);
  assert.match(link.url, /gid=838298230/);
});

test("testPlanLink covers both ends of the connection block", () => {
  assert.match(testPlanLink("TC-1012").url, /range=A234$/); // first
  assert.match(testPlanLink("TC-1089").url, /range=A311$/); // last (1089 - 778)
});

test("testPlanLink returns null outside the connection block or for malformed ids", () => {
  assert.equal(testPlanLink("TC-1011"), null); // below the block
  assert.equal(testPlanLink("TC-1090"), null); // above the block
  assert.equal(testPlanLink("TC-500"), null); // unrelated case
  assert.equal(testPlanLink("nonsense"), null);
  assert.equal(testPlanLink(""), null);
  assert.equal(testPlanLink(undefined), null);
});

test("parseTitleTokens emits a Test Plan tms link for a connection TC", () => {
  const { links } = parseTitleTokens("[TC-1042][cut=Kubernetes Connection] view works");
  assert.equal(links.length, 1);
  assert.equal(links[0].type, "tms");
  assert.equal(links[0].name, "Test Plan TC-1042");
  assert.match(links[0].url, /range=A264$/); // 1042 - 778
});

test("parseTitleTokens emits no link for a non-connection Test #", () => {
  const { links } = parseTitleTokens("[TC-2001][cut=Model] import works");
  assert.equal(links.length, 0);
});

test("convertTapToAllure attaches the Test Plan link to the result", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "b2a-"));
  const cwd = process.cwd();
  try {
    process.chdir(tmp);
    fs.writeFileSync(
      path.join(tmp, "s.tap"),
      "1..1\nok 1 [TC-1013][cut=Kubernetes Connection][tg=Connection Lifecycle] create\n"
    );
    const [r] = convertTapToAllure(path.join(tmp, "s.tap"));
    assert.equal(r.links.length, 1);
    assert.equal(r.links[0].name, "Test Plan TC-1013");
    assert.match(r.links[0].url, /range=A235$/); // 1013 - 778
  } finally {
    process.chdir(cwd);
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

// --- Part B: failure output capture ----------------------------------------

test("summarizeDiagnostics prefers the assertion '...failed' line", () => {
  const diag = [
    "(in test file 007-connection/02-connection-list.bats, line 25)",
    "  `assert_success' failed",
    "-- command failed --",
    "status : 1"
  ];
  assert.equal(summarizeDiagnostics(diag), "`assert_success' failed");
});

test("summarizeDiagnostics falls back to the first non-empty line", () => {
  assert.equal(summarizeDiagnostics(["", "   ", "boom happened"]), "boom happened");
  assert.equal(summarizeDiagnostics([]), "");
});

test("failed tests get a trace + text attachment holding the captured output", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "b2a-"));
  const cwd = process.cwd();
  try {
    process.chdir(tmp);
    const tap = [
      "1..2",
      "ok 1 [TC-1041][cut=Kubernetes Connection] list works",
      "not ok 2 [TC-1043][cut=Kubernetes Connection] create fails",
      "# (in test file 007-connection/01-connection-create.bats, line 55)",
      "#   `assert_success' failed",
      "# -- command failed --",
      "# status : 1",
      "# output : Error: unable to get: .connections"
    ].join("\n");
    fs.writeFileSync(path.join(tmp, "sample.tap"), tap);

    const results = convertTapToAllure(path.join(tmp, "sample.tap"));
    const passed = results.find(r => r.status === "passed");
    const failed = results.find(r => r.status === "failed");

    // Headline message is the concise assertion line.
    assert.equal(failed.statusDetails.message, "`assert_success' failed");
    // Full transcript is in the trace...
    assert.match(failed.statusDetails.trace, /unable to get: \.connections/);
    assert.match(failed.statusDetails.trace, /command failed/);
    // ...and duplicated as a text attachment file on disk.
    assert.equal(failed.attachments.length, 1);
    assert.equal(failed.attachments[0].type, "text/plain");
    assert.equal(failed.attachments[0].source, `${failed.uuid}-attachment.txt`);
    const attachmentBody = fs.readFileSync(
      path.join(tmp, "allure-results", failed.attachments[0].source),
      "utf-8"
    );
    assert.match(attachmentBody, /unable to get: \.connections/);

    // A passing test carries no trace/attachment and no leaked scratch field.
    assert.equal(passed.statusDetails.trace, undefined);
    assert.equal(passed.attachments, undefined);
    assert.equal(passed._diag, undefined);
    assert.equal(failed._diag, undefined);
  } finally {
    process.chdir(cwd);
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test("skipped tests keep their skip-reason message and get no attachment", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "b2a-"));
  const cwd = process.cwd();
  try {
    process.chdir(tmp);
    fs.writeFileSync(
      path.join(tmp, "s.tap"),
      "1..1\nok 1 [TC-1013][cut=Kubernetes Connection] create # skip minikube not installed\n"
    );
    const [r] = convertTapToAllure(path.join(tmp, "s.tap"));
    assert.equal(r.status, "skipped");
    assert.equal(r.statusDetails.message, "minikube not installed");
    assert.equal(r.statusDetails.trace, undefined);
    assert.equal(r.attachments, undefined);
  } finally {
    process.chdir(cwd);
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});
