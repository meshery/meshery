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
  convertTapToAllure,
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
