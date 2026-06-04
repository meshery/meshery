#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

const ALLURE_RESULTS_DIR = process.env.ALLURE_OUTPUT_PATH || "allure-results";

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

function goStatusToAllure(action) {
  switch (action) {
    case "pass": return "passed";
    case "fail": return "failed";
    case "skip": return "skipped";
    default: return "broken";
  }
}

function createAllureResult({ name, fullName, status, start, stop, output, pkg, extraLabels }) {
  const suite = pkg.replace(/^github\.com\/meshery\/meshery\//, "");
  const projectLabel = extraLabels.find(l => l.name === "project");
  const baseLabels = [
    { name: "framework", value: "go-test" },
    { name: "language", value: "golang" },
    { name: "parentSuite", value: suite.split("/").slice(0, 2).join("/") },
    { name: "suite", value: suite },
    { name: "subSuite", value: name },
  ];

  if (!projectLabel) {
    baseLabels.push({ name: "project", value: "Meshery" });
  }

  return {
    uuid: uuid(),
    name,
    fullName: fullName || name,
    status,
    stage: "finished",
    start,
    stop,
    statusDetails: output ? { message: output } : {},
    labels: [...baseLabels, ...extraLabels]
  };
}

function convertGoTestToAllure(jsonFile) {
  ensureAllureDir();

  const content = fs.readFileSync(jsonFile, "utf-8");
  const lines = content.split("\n").filter(Boolean);

  const extraLabels = parseExtraLabels();
  const tests = new Map();

  for (const line of lines) {
    let event;
    try {
      event = JSON.parse(line);
    } catch {
      continue;
    }

    const { Action: action, Package: pkg, Test: testName, Time: time, Elapsed: elapsed, Output: output } = event;

    if (!testName) continue;

    const key = `${pkg}::${testName}`;

    if (action === "run") {
      tests.set(key, {
        name: testName,
        pkg,
        start: new Date(time).getTime(),
        output: []
      });
    } else if (action === "output" && tests.has(key)) {
      tests.get(key).output.push(output);
    } else if (action === "pass" || action === "fail" || action === "skip") {
      const test = tests.get(key);
      if (!test) continue;

      const stop = new Date(time).getTime();
      const outputText = test.output
        .filter(o => !o.startsWith("=== RUN") && !o.startsWith("--- "))
        .join("")
        .trim();

      const result = createAllureResult({
        name: test.name,
        fullName: `${test.pkg}::${test.name}`,
        status: goStatusToAllure(action),
        start: test.start,
        stop: stop || test.start + (elapsed || 0) * 1000,
        output: outputText || null,
        pkg: test.pkg,
        extraLabels
      });

      const outputPath = path.join(ALLURE_RESULTS_DIR, `${result.uuid}-result.json`);
      fs.writeFileSync(outputPath, JSON.stringify(result, null, 2));

      tests.delete(key);
    }
  }

  const resultFiles = fs.readdirSync(ALLURE_RESULTS_DIR).filter(f => f.endsWith("-result.json"));
  console.log(`Converted ${resultFiles.length} Go tests to Allure results in ${ALLURE_RESULTS_DIR}/`);
}

if (process.argv.length !== 3) {
  console.error("Usage: node gotest-to-allure.js <go-test-output.json>");
  console.error("  Generate input: go test -json ./... > go-test-output.json");
  process.exit(1);
}

convertGoTestToAllure(process.argv[2]);
