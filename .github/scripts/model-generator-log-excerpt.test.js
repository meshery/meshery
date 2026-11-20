'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const {
  existsSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
} = require('node:fs');
const os = require('node:os');
const path = require('node:path');

const workflowPath = path.join(__dirname, '..', 'workflows', 'model-generator.yml');
const workflowLines = readFileSync(workflowPath, 'utf8').split(/\r?\n/);
const functionStart = workflowLines.findIndex((line) => line.trim() === 'append_log_excerpt() {');
const functionEnd = workflowLines.findIndex(
  (line, index) => index > functionStart && line === '          }',
);

assert.notEqual(functionStart, -1, 'append_log_excerpt must exist in the workflow');
assert.notEqual(functionEnd, -1, 'append_log_excerpt must have a closing brace');

const bashFunction = workflowLines
  .slice(functionStart, functionEnd + 1)
  .map((line) => line.slice(10))
  .join('\n');

const bashExecutable =
  process.env.BASH ||
  (process.platform === 'win32' &&
  existsSync(path.join(process.env.ProgramFiles || 'C:\\Program Files', 'Git', 'bin', 'bash.exe'))
    ? path.join(process.env.ProgramFiles || 'C:\\Program Files', 'Git', 'bin', 'bash.exe')
    : 'bash');

function toBashPath(filePath) {
  return process.platform === 'win32' ? filePath.replace(/\\/g, '/') : filePath;
}

function makeExcerpt(t, logLines) {
  const tempDir = mkdtempSync(path.join(os.tmpdir(), 'model-generator-excerpt-'));
  t.after(() => rmSync(tempDir, { recursive: true, force: true }));

  const emailPath = path.join(tempDir, 'email.txt');
  const logPath = path.join(tempDir, 'step.log');
  writeFileSync(emailPath, '');
  writeFileSync(logPath, `${logLines.join('\n')}\n`);

  const script = `${bashFunction}
EMAIL_FILE="$1"
append_log_excerpt "$2" "test step"
`;
  const result = spawnSync(
    bashExecutable,
    ['-c', script, 'model-generator-log-excerpt-test', toBashPath(emailPath), toBashPath(logPath)],
    { encoding: 'utf8' },
  );

  assert.equal(result.status, 0, result.stderr);
  return readFileSync(emailPath, 'utf8');
}

function lineCount(text) {
  return (text.match(/\n/g) || []).length;
}

test('long excerpts keep the beginning and terminating log lines within 100 lines', (t) => {
  const logLines = Array.from(
    { length: 250 },
    (_, index) => `line ${index + 1}: recoverable error ${index + 1}`,
  );
  logLines[244] = 'line 245: fatal terminating error';
  logLines[249] = 'line 250: final context after terminating error';

  const excerpt = makeExcerpt(t, logLines);

  assert.ok(lineCount(excerpt) <= 100);
  assert.ok(excerpt.includes(logLines[0]), 'the start of the error excerpt should be retained');
  assert.ok(excerpt.includes(logLines[244]), 'the terminating error should be retained');
  assert.ok(excerpt.includes(logLines[249]), 'the final log line should be retained');
});

test('short logs are not truncated', (t) => {
  const logLines = Array.from(
    { length: 25 },
    (_, index) => `short log line ${index + 1}: recoverable error`,
  );

  const excerpt = makeExcerpt(t, logLines);

  assert.ok(lineCount(excerpt) <= 100);
  for (const logLine of logLines) {
    assert.ok(excerpt.includes(logLine), `missing log line: ${logLine}`);
  }
});
