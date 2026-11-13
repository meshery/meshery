'use strict';

const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { test } = require('node:test');

const productionScript = path.join(__dirname, 'audit-hex.js');

function writeFile(root, relativePath, contents) {
  const file = path.join(root, relativePath);
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, contents);
}

function createFixture(t) {
  const fixture = fs.mkdtempSync(path.join(os.tmpdir(), 'audit-hex-'));
  t.after(() => fs.rmSync(fixture, { recursive: true, force: true }));

  fs.mkdirSync(path.join(fixture, 'ui', 'scripts'), { recursive: true });
  fs.copyFileSync(productionScript, path.join(fixture, 'ui', 'scripts', 'audit-hex.js'));

  fs.mkdirSync(path.join(fixture, 'ui', 'components'), { recursive: true });

  return fixture;
}

function runAudit(fixture) {
  const script = path.join(fixture, 'ui', 'scripts', 'audit-hex.js');
  const result = spawnSync(process.execPath, [script, '--json'], {
    encoding: 'utf8',
  });

  assert.equal(result.status, 0, result.stderr);
  return JSON.parse(result.stdout);
}

test('excludes test files and nested __tests__ directories', (t) => {
  const fixture = createFixture(t);

  writeFile(fixture, 'ui/components/production.tsx', "const color = '#fff';\n");
  writeFile(fixture, 'ui/components/widget.test.tsx', "const color = '#000';\n");
  writeFile(fixture, 'ui/components/helper.test.js', "const color = '#123';\n");
  writeFile(fixture, 'ui/components/helper.test.ts', "const color = '#789';\n");
  writeFile(fixture, 'ui/components/another.test.jsx', "const color = '#abc';\n");
  writeFile(fixture, 'ui/components/__tests__/helper.ts', "const color = '#456';\n");

  const result = runAudit(fixture);

  assert.equal(result.files, 1);
  assert.equal(result.matches, 1);
  assert.deepEqual(
    result.offenders.map((offender) => offender.file),
    ['components/production.tsx'],
  );
});

test('does not exclude production files merely containing the word test', (t) => {
  const fixture = createFixture(t);

  writeFile(fixture, 'ui/components/test-helper.tsx', "const color = '#abc';\n");

  const result = runAudit(fixture);

  assert.equal(result.files, 1);
  assert.equal(result.matches, 1);
  assert.equal(result.offenders[0].file, 'components/test-helper.tsx');
});

test('detects case-insensitive RGB functions without matching identifier suffixes', (t) => {
  const fixture = createFixture(t);

  const cases = {
    'lowercase.tsx': ['rgb(255, 0, 0)', 1],
    'uppercase.tsx': ['RGB(255, 0, 0)', 1],
    'mixed-case.tsx': ['Rgb(255, 0, 0)', 1],
    'rgba.tsx': ['rgba(0, 0, 0, .5)', 1],
    'uppercase-rgba.tsx': ['RGBA(0, 0, 0, .5)', 1],
    'member-call.tsx': ['theme.rgb(255, 0, 0)', 1],
    'identifier.tsx': ['myrgb(255, 0, 0)', 0],
    'dollar-identifier.tsx': ['$rgb(255, 0, 0)', 0],
    'hyphen-identifier.tsx': ['-rgb(255, 0, 0)', 0],
  };

  for (const [filename, [content]] of Object.entries(cases)) {
    writeFile(fixture, `ui/components/${filename}`, content);
  }

  const result = runAudit(fixture);
  const actual = Object.fromEntries(
    Object.entries(result.offenders).map(([index, offender]) => [
      offender.file.split('/').pop(),
      offender.rgb,
    ]),
  );

  const expected = Object.fromEntries(
    Object.entries(cases).filter(([, [, count]]) => count > 0)
      .map(([filename, [, count]]) => [filename, count]),
  );

  assert.deepEqual(actual, expected);
  assert.equal(result.rgb, 6);
});
