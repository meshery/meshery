#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

const ALLURE_RESULTS_DIR = "allure-results";

// TAP regex patterns
const TAP_TEST_RE = /^(ok|not ok)\s+\d+\s+(.*)/;
const TAP_DIAGNOSTIC_RE = /^\s*#\s?(.*)/;

// Parse labels from env: key=value,key=value
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

function createAllureResult({
  name,
  status,
  start,
  stop,
  details,
  extraLabels
}) {
  return {
    uuid: uuid(),
    name,
    status,
    stage: "finished",
    start,
    stop,
    statusDetails: details ? { message: details } : {},
    labels: [
      { name: "framework", value: "bats" },
      { name: "language", value: "bash" },
      { name: "project" , value: "mesheryctl" },
      ...extraLabels
    ]
  };
}

function convertTapToAllure(tapFile) {
  ensureAllureDir();

  const content = fs.readFileSync(tapFile, "utf-8");
  const lines = content.split("\n");

  const extraLabels = parseExtraLabels();

  let diagnostics = [];
  let testCount = 0;

  for (const line of lines) {
    const trimmed = line.trimEnd();

    // Diagnostics
    const diagMatch = trimmed.match(TAP_DIAGNOSTIC_RE);
    if (diagMatch) {
      diagnostics.push(diagMatch[1]);
      continue;
    }

    // Test line
    const testMatch = trimmed.match(TAP_TEST_RE);
    if (!testMatch) continue;

    testCount++;

    const [, rawStatus, name] = testMatch;
    const status = tapStatusToAllure(rawStatus);

    const start = Date.now();
    const stop = start + 1;

    const result = createAllureResult({
      name,
      status,
      start,
      stop,
      details: diagnostics.length ? diagnostics.join("\n") : null,
      extraLabels
    });

    diagnostics = [];

    const outputPath = path.join(
      ALLURE_RESULTS_DIR,
      `${result.uuid}-result.json`
    );

    fs.writeFileSync(outputPath, JSON.stringify(result, null, 2));
  }

  console.log(`Converted ${testCount} BATS tests to Allure results.`);
}

// CLI
if (process.argv.length !== 3) {
  console.error("Usage: node bats-to-allure.js <bats-report.tap>");
  process.exit(1);
}

convertTapToAllure(process.argv[2]);

