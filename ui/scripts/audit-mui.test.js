'use strict';

const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { test } = require('node:test');

const productionScript = path.join(__dirname, 'audit-mui.js');
const typescriptEntry = require.resolve('typescript', {
  paths: [path.resolve(__dirname, '..')],
});
const typescriptNodeModules = path.resolve(path.dirname(typescriptEntry), '..', '..');
const childEnv = {
  ...process.env,
  NODE_PATH: [typescriptNodeModules, process.env.NODE_PATH].filter(Boolean).join(path.delimiter),
};

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
  fs.mkdirSync(path.join(fixture, 'install', 'docker-extension', 'ui', 'src'), { recursive: true });
  writeFile(fixture, 'ui/remote-component.config.js');

  return fixture;
}

function runAudit(fixture, args = [], options = {}) {
  const script = path.join(fixture, 'ui', 'scripts', 'audit-mui.js');
  const result = spawnSync(process.execPath, [script, '--json', ...args], {
    encoding: 'utf8',
    env: { ...childEnv, ...(options.env ?? {}) },
  });

  let json;
  try {
    json = JSON.parse(result.stdout);
  } catch (error) {
    error.message = `Expected JSON stdout, got:\n${result.stdout}\n${error.message}`;
    throw error;
  }
  return { ...result, json };
}

function runAuditHuman(fixture, args = []) {
  const script = path.join(fixture, 'ui', 'scripts', 'audit-mui.js');

  return spawnSync(process.execPath, [script, ...args], {
    encoding: 'utf8',
    env: childEnv,
  });
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
  assert.match(scanErrorText(result), /missing scan file .*remote-component\.config\.js/);
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

test('side-effect MUI imports are reported as offenders', (t) => {
  const fixture = createFixture(t);
  writeFile(fixture, 'ui/components/side-effect.js', "import '@mui/material/Button';\n");

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.scanComplete, true);
  assert.equal(result.json.files, 1);
  assert.equal(result.json.matches, 1);
  assert.ok(
    result.json.scans[0].offenders.some(
      (offender) => offender.file === 'components/side-effect.js',
    ),
  );
});

test('comments between import and its module specifier do not hide MUI imports', (t) => {
  const fixture = createFixture(t);
  writeFile(
    fixture,
    'ui/components/commented-import.js',
    "import /* polyfill */ '@mui/material/Button';\n",
  );

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.matches, 1);
  assert.ok(
    result.json.scans[0].offenders.some(
      (offender) => offender.file === 'components/commented-import.js',
    ),
  );
});

test('MUI import-like text in comments and strings is ignored', (t) => {
  const fixture = createFixture(t);
  writeFile(
    fixture,
    'ui/components/non-import-text.js',
    [
      "// import '@mui/material/Button';",
      'const example = "import \'@mui/material/Button\';";',
      "/* require('@mui/material'); */",
      '',
    ].join('\n'),
  );

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 0);
  assert.equal(result.json.matches, 0);
});

test('MUI re-exports are reported', (t) => {
  const fixture = createFixture(t);
  writeFile(
    fixture,
    'ui/components/re-export.js',
    "export { Button } from '@mui/material/Button';\n",
  );

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.matches, 1);
});

test('MUI require calls are reported', (t) => {
  const fixture = createFixture(t);
  writeFile(
    fixture,
    'ui/components/require.js',
    "const Button = require('@mui/material/Button');\n",
  );

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.matches, 1);
});

test('TypeScript import-equals MUI imports are reported', (t) => {
  const fixture = createFixture(t);

  writeFile(
    fixture,
    'ui/components/import-equals.ts',
    "import Mui = require('@mui/material/Button');\n",
  );

  const result = runAudit(fixture, ['--fail-on-new']);

  assert.equal(result.status, 1);
  assert.equal(result.json.matches, 1);

  assert.ok(
    result.json.scans[0].offenders.some(
      (offender) => offender.file === 'components/import-equals.ts',
    ),
  );
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

test('allowlist status reports files with MUI imports as active', (t) => {
  const fixture = createFixture(t);
  const relativePath = 'theme/index.ts';

  writeFile(fixture, `ui/${relativePath}`, "import { darken } from '@mui/material/styles';\n");

  const result = runAudit(fixture);
  const uiScan = result.json.scans.find((scan) => scan.name === 'ui');

  assert.ok(uiScan.allowlistStatus.active.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.unused.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.missing.includes(relativePath));
});

test('allowlist status reports files without MUI imports as unused', (t) => {
  const fixture = createFixture(t);
  const relativePath = 'theme/index.ts';

  writeFile(fixture, `ui/${relativePath}`, 'export const theme = {};\n');

  const result = runAudit(fixture, ['--fail-on-new']);
  const uiScan = result.json.scans.find((scan) => scan.name === 'ui');

  assert.equal(result.status, 0);
  assert.equal(result.json.scanComplete, true);
  assert.ok(uiScan.allowlistStatus.unused.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.active.includes(relativePath));
  const humanResult = runAuditHuman(fixture);

  assert.match(humanResult.stdout, /\[unused\].*theme\/index\.ts/);
});

test('allowlist status reports nonexistent files as missing', (t) => {
  const fixture = createFixture(t);
  const relativePath = 'theme/index.ts';

  const result = runAudit(fixture, ['--fail-on-new']);
  const uiScan = result.json.scans.find((scan) => scan.name === 'ui');

  assert.ok(uiScan.allowlistStatus.missing.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.active.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.unused.includes(relativePath));
  assert.equal(result.status, 0);
  assert.equal(result.json.scanComplete, true);

  const humanResult = runAuditHuman(fixture);

  assert.equal(humanResult.status, 0);
  assert.match(humanResult.stdout, /\[missing\].*theme\/index\.ts/);
});

test('a failed read does not classify an allowlisted file as unused or missing', (t) => {
  const fixture = createFixture(t);
  const relativePath = 'theme/index.ts';
  const unreadable = path.join(fixture, 'ui', relativePath);

  writeFile(fixture, `ui/${relativePath}`, "import { darken } from '@mui/material/styles';\n");

  const preload = path.join(fixture, 'allowlist-read-failure.cjs');

  writeFile(
    fixture,
    'allowlist-read-failure.cjs',
    `
      const fs = require('node:fs');
      const path = require('node:path');
      const original = fs.readFileSync;
      const blocked = ${JSON.stringify(unreadable)};

      fs.readFileSync = (file, ...args) => {
        if (path.resolve(file) === blocked) {
          throw new Error('injected allowlist read failure');
        }
        return original(file, ...args);
      };
    `,
  );

  const result = runAudit(fixture, [], {
    env: { NODE_OPTIONS: '--require=' + preload },
  });

  const uiScan = result.json.scans.find((scan) => scan.name === 'ui');

  assert.equal(result.json.scanComplete, false);
  assert.match(scanErrorText(result), /injected allowlist read failure/);
  assert.ok(!uiScan.allowlistStatus.active.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.unused.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.missing.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.uncovered.includes(relativePath));
});

test('existing allowlist files outside the scanned file set are uncovered', (t) => {
  const fixture = createFixture(t);
  const relativePath = 'components/not-scanned.json';
  const script = path.join(fixture, 'ui/scripts/audit-mui.js');

  writeFile(fixture, `ui/${relativePath}`, '{}\n');

  const original = fs.readFileSync(script, 'utf8');
  const anchor = "      'theme/index.ts',";
  assert.equal(original.split(anchor).length - 1, 1);
  fs.writeFileSync(script, original.replace(anchor, `${anchor}\n      '${relativePath}',`));

  const result = runAudit(fixture, ['--fail-on-new']);
  const uiScan = result.json.scans.find((scan) => scan.name === 'ui');

  assert.equal(result.status, 0);
  assert.equal(result.json.scanComplete, true);
  assert.ok(uiScan.allowlistStatus.uncovered.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.active.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.unused.includes(relativePath));
  assert.ok(!uiScan.allowlistStatus.missing.includes(relativePath));

  const humanResult = runAuditHuman(fixture);
  assert.equal(humanResult.status, 0);
  assert.match(humanResult.stdout, /\[uncovered\].*components\/not-scanned\.json/);
});

test('scannedFiles counts only files that were successfully read', (t) => {
  const fixture = createFixture(t);
  const unreadable = path.join(fixture, 'ui', 'components', 'unreadable.js');

  writeFile(fixture, 'ui/components/unreadable.js', "require('@mui/material');\n");
  writeFile(fixture, 'ui/components/readable-a.js', 'export const a = 1;\n');
  writeFile(fixture, 'ui/components/readable-b.js', 'export const b = 1;\n');

  const preload = path.join(fixture, 'scanned-count-read-failure.cjs');
  writeFile(
    fixture,
    'scanned-count-read-failure.cjs',
    `
      const fs = require('node:fs');
      const path = require('node:path');
      const original = fs.readFileSync;
      const blocked = ${JSON.stringify(unreadable)};

      fs.readFileSync = (file, ...args) => {
        if (path.resolve(file) === blocked) {
          throw new Error('injected read failure for scannedFiles test');
        }
        return original(file, ...args);
      };
    `,
  );

  const result = runAudit(fixture, [], {
    env: { NODE_OPTIONS: '--require=' + preload },
  });

  const uiScan = result.json.scans.find((scan) => scan.name === 'ui');

  // 3 readable files: readable-a, readable-b, and remote-component.config.js
  // created by createFixture. The unreadable file must not be counted.
  assert.equal(uiScan.scannedFiles, 3);

  // overall total equals the sum of the per-scan counts
  const perScanSum = result.json.scans.reduce((acc, scan) => acc + scan.scannedFiles, 0);
  assert.equal(result.json.scannedFiles, perScanSum);
});
