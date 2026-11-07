'use strict';

const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { test } = require('node:test');

const productionScript = path.join(__dirname, 'audit-mui.js');

const uiDirectories = [
  'components',
  'pages',
  'utils',
  'theme',
  'machines',
  'store',
  'rtk-query',
  'lib',
  'api',
  'public',
];

function writeFile(root, relativePath, contents = '') {
  const file = path.join(root, relativePath);
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, contents);
}

function createFixture(t) {
  const fixture = fs.mkdtempSync(path.join(os.tmpdir(), 'audit-mui-'));
  t.after(() => fs.rmSync(fixture, { recursive: true, force: true }));

  fs.mkdirSync(path.join(fixture, 'ui', 'scripts'), { recursive: true });
  fs.copyFileSync(productionScript, path.join(fixture, 'ui', 'scripts', 'audit-mui.js'));

  for (const directory of uiDirectories) {
    fs.mkdirSync(path.join(fixture, 'ui', directory), { recursive: true });
  }
  fs.mkdirSync(
    path.join(fixture, 'install', 'docker-extension', 'ui', 'src'),
    { recursive: true },
  );
  writeFile(fixture, 'ui/remote-component.config.js');

  return fixture;
}

function runAudit(fixture, args = [], options = {}) {
  const script = path.join(fixture, 'ui', 'scripts', 'audit-mui.js');
  const result = spawnSync(
    process.execPath,
    [script, '--json', ...args],
    {
      encoding: 'utf8',
      env: { ...process.env, ...(options.env ?? {}) },
    },
  );

  let json;
  try {
    json = JSON.parse(result.stdout);
  } catch (error) {
    error.message = `Expected JSON stdout, got:\n${result.stdout}\n${error.message}`;
    throw error;
  }
  return { ...result, json };
}

function scanErrors(result) {
  return result.json.scanErrors ?? result.json.scans.flatMap((scan) => scan.scanErrors ?? []);
}

function scanErrorText(result) {
  return JSON.stringify(scanErrors(result));
}

test('a complete scan reports scanComplete true and no scan errors in JSON', (t) => {
  const fixture = createFixture(t);
  const result = runAudit(fixture);

  assert.equal(result.status, 0);
  assert.equal(result.json.scanComplete, true);
  assert.deepEqual(scanErrors(result), []);
});

test('a missing configured directory is recorded as a scan error', (t) => {
  const fixture = createFixture(t);
  fs.rmSync(path.join(fixture, 'ui', 'components'), { recursive: true });

  const result = runAudit(fixture);

  assert.equal(result.json.scanComplete, false);
  assert.match(scanErrorText(result), /missing scan dir .*ui[\\/]components/);
});

test('a missing explicit scan file is recorded as a scan error', (t) => {
  const fixture = createFixture(t);
  fs.rmSync(path.join(fixture, 'ui', 'remote-component.config.js'));

  const result = runAudit(fixture);

  assert.equal(result.json.scanComplete, false);
  assert.match(
    scanErrorText(result),
    /missing scan file .*remote-component\.config\.js/,
  );
});

test('a recursive directory read failure is recorded while sibling directories are scanned', (t) => {
  const fixture = createFixture(t);
  const inaccessible = path.join(fixture, 'ui', 'components', 'broken');

  fs.mkdirSync(inaccessible);
  writeFile(fixture, 'ui/components/broken/hidden.js', "require('@mui/material');\n");
  writeFile(fixture, 'ui/components/accessible.js', "require('@mui/material');\n");

  const preload = path.join(fixture, 'readdir-failure.cjs');
  writeFile(
    fixture,
    'readdir-failure.cjs',
    `
      const fs = require('node:fs');
      const original = fs.readdirSync;
      const blocked = ${JSON.stringify(inaccessible)};
      fs.readdirSync = (directory, options) => {
        if (directory === blocked) throw new Error('injected recursive read failure');
        return original(directory, options);
      };
    `,
  );

  const result = runAudit(fixture, [], {
    env: { NODE_OPTIONS: `--require=${preload}` },
  });

  assert.equal(result.json.scanComplete, false);
  assert.match(scanErrorText(result), /injected recursive read failure/);
  const offenders = result.json.scans.flatMap((scan) => scan.offenders);
  assert.ok(offenders.some((offender) => offender.file === 'components/accessible.js'));
  assert.ok(!offenders.some((offender) => offender.file === 'components/broken/hidden.js'));
});

test('a file read failure is recorded while remaining files are scanned', (t) => {
  const fixture = createFixture(t);
  const unreadable = path.join(fixture, 'ui', 'components', 'unreadable.js');

  writeFile(fixture, 'ui/components/unreadable.js', "require('@mui/material');\n");
  writeFile(fixture, 'ui/components/accessible.js', "require('@mui/material');\n");

  const preload = path.join(fixture, 'readfile-failure.cjs');
  writeFile(
    fixture,
    'readfile-failure.cjs',
    `
      const fs = require('node:fs');
      const path = require('node:path');
      const original = fs.readFileSync;
      const blocked = ${JSON.stringify(unreadable)};

      fs.readFileSync = (file, ...args) => {
        if (path.resolve(file) === blocked) {
          throw new Error('injected file read failure');
        }
        return original(file, ...args);
      };
    `,
  );

  const result = runAudit(fixture, ['--fail-on-new'], {
    env: { NODE_OPTIONS: '--require=' + preload },
  });

  assert.equal(result.status, 1);
  assert.equal(result.json.scanComplete, false);
  assert.match(scanErrorText(result), /injected file read failure/);

  const offenders = result.json.scans.flatMap((scan) => scan.offenders);

  assert.ok(offenders.some((offender) => offender.file === 'components/accessible.js'));
  assert.ok(!offenders.some((offender) => offender.file === 'components/unreadable.js'));
});

test('--fail-on-new exits 1 when the scan is incomplete even with no offenders', (t) => {
  const fixture = createFixture(t);
  fs.rmSync(path.join(fixture, 'ui', 'components'), { recursive: true });

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.scanComplete, false);
  assert.equal(result.json.files, 0);
});

test('--fail-on-new exits 1 when non-allowlisted offending imports exist', (t) => {
  const fixture = createFixture(t);
  writeFile(fixture, 'ui/components/offender.js', "import Button from '@mui/material/Button';\n");

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.scanComplete, true);
  assert.equal(result.json.files, 1);
});

test('allowlisted imports are excluded from offender counts', (t) => {
  const fixture = createFixture(t);
  writeFile(fixture, 'ui/theme/index.ts', "import { darken } from '@mui/material/styles';\n");

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 0);
  assert.equal(result.json.scanComplete, true);
  assert.equal(result.json.files, 0);
  assert.equal(result.json.scans[0].allowlistedFilesWithMuiImports, 1);
});

test('a complete scan with no offenders exits 0', (t) => {
  const fixture = createFixture(t);
  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 0);
  assert.equal(result.json.scanComplete, true);
  assert.equal(result.json.files, 0);
  assert.equal(result.json.matches, 0);
});

