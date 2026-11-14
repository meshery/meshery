
'use strict';

const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { test } = require('node:test');

const productionScript = path.join(__dirname, 'audit-hex.js');
const typescriptModuleRoot = path.resolve(path.dirname(require.resolve('typescript')), '..', '..');

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

  const scanDirs = [
    'components',
    'pages',
    'utils',
    'store',
    'rtk-query',
    'machines',
    'hooks',
    'api',
    'css',
  ];

  for (const dir of scanDirs) {
    fs.mkdirSync(path.join(fixture, 'ui', dir), { recursive: true });
  }

  return fixture;
}

function runAudit(fixture, args = ['--json'], options = {}) {
  const script = path.join(fixture, 'ui', 'scripts', 'audit-hex.js');

  const result = spawnSync(process.execPath, [...(options.nodeArgs || []), script, ...args], {
    encoding: 'utf8',
    env: {
      ...process.env,
      NODE_PATH: [typescriptModuleRoot, process.env.NODE_PATH]
        .filter(Boolean)
        .join(path.delimiter),
      ...options.env,
    },
  });

  assert.equal(result.status, 0, result.stderr);

  return args.includes('--json') ? { ...result, json: JSON.parse(result.stdout) } : result;
}

function addFilesystemFailureMock(fixture) {
  const preload = path.join(fixture, 'mock-fs-errors.js');

  fs.writeFileSync(
    preload,
    `
'use strict';

const fs = require('node:fs');
const path = require('node:path');

const originalReaddirSync = fs.readdirSync;
const originalReadFileSync = fs.readFileSync;

fs.readdirSync = function (target, ...args) {
  if (path.basename(String(target)) === 'blocked-directory') {
    const error = new Error('Permission denied');
    error.code = 'EACCES';
    throw error;
  }

  return originalReaddirSync.call(this, target, ...args);
};

fs.readFileSync = function (target, ...args) {
  if (path.basename(String(target)) === 'blocked-file.tsx') {
    const error = new Error('Permission denied');
    error.code = 'EACCES';
    throw error;
  }

  return originalReadFileSync.call(this, target, ...args);
};
`,
  );

  return preload;
}

test('ignores hex-like text in comments but preserves actual literals', (t) => {
  const fixture = createFixture(t);

  writeFile(
    fixture,
    'ui/components/comment-cases.tsx',
    [
      '// React error #185',
      '/* Documentation reference #123456 */',
      "const stringColor = '#F91313';",
      'const templateColor = `#AABBCC`;',
      'const jsx = <div title="#ABCDEF" />;',
      'const commentLikeString = "/* #654321 */";',
      'const templateInterpolation = `${ /* #112233 */ "#445566" }`;',
      'const slashLikeString = "https://example.com/path//segment#778899";',
      'const jsxComment = <div>{/* #987654 */}</div>;',
      'const jsxTextWithSlashes = <div>// label <span style={{ color: "#abcdef" }} /></div>;',
      'const regexWithSlashes = /[//]/; const regexColor = "#fedcba";',
      '',
    ].join('\n'),
  );

  const { json: result } = runAudit(fixture);

  assert.equal(result.files, 1);
  assert.equal(result.hex, 8);
  assert.equal(result.rgb, 0);
  assert.equal(result.matches, 8);
  assert.equal(result.offenders[0].hex, 8);
});

test('excludes test files and nested __tests__ directories', (t) => {
  const fixture = createFixture(t);

  writeFile(fixture, 'ui/components/production.tsx', "const color = '#fff';\n");
  writeFile(fixture, 'ui/components/widget.test.tsx', "const color = '#000';\n");
  writeFile(fixture, 'ui/components/helper.test.js', "const color = '#123';\n");
  writeFile(fixture, 'ui/components/helper.test.ts', "const color = '#789';\n");
  writeFile(fixture, 'ui/components/another.test.jsx', "const color = '#abc';\n");
  writeFile(fixture, 'ui/components/__tests__/helper.ts', "const color = '#456';\n");

  const { json: result } = runAudit(fixture);

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

  const { json: result } = runAudit(fixture);

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

  const { json: result } = runAudit(fixture);
  const actual = Object.fromEntries(
    Object.entries(result.offenders).map(([, offender]) => [
      offender.file.split('/').pop(),
      offender.rgb,
    ]),
  );

  const expected = Object.fromEntries(
    Object.entries(cases)
      .filter(([, [, count]]) => count > 0)
      .map(([filename, [, count]]) => [filename, count]),
  );

  assert.deepEqual(actual, expected);
  assert.equal(result.rgb, 6);
});

test('successful scans report no filesystem errors', (t) => {
  const fixture = createFixture(t);

  writeFile(fixture, 'ui/components/valid.tsx', "const color = '#fff';\n");

  const { json } = runAudit(fixture);

  assert.equal(json.files, 1);
  assert.equal(json.matches, 1);
  assert.deepEqual(json.scanErrors, []);
  assert.equal(json.scanComplete, true);
});

test('directory read failures are reported while accessible files are scanned', (t) => {
  const fixture = createFixture(t);
  const preload = addFilesystemFailureMock(fixture);

  writeFile(fixture, 'ui/components/valid.tsx', "const color = '#fff';\n");
  writeFile(fixture, 'ui/components/blocked-directory/hidden.tsx', "const color = '#000';\n");

  const { json } = runAudit(fixture, ['--json'], {
    nodeArgs: ['--require', preload],
  });

  assert.equal(json.files, 1);
  assert.equal(json.matches, 1);
  assert.equal(json.offenders[0].file, 'components/valid.tsx');
  assert.equal(json.scanComplete, false);
  assert.equal(json.scanErrors.length, 1);
  assert.match(json.scanErrors[0], /Unable to read directory/);
  assert.match(json.scanErrors[0], /blocked-directory/);
});

test('file read failures are reported while accessible files are scanned', (t) => {
  const fixture = createFixture(t);
  const preload = addFilesystemFailureMock(fixture);

  writeFile(fixture, 'ui/components/valid.tsx', "const color = '#abc';\n");
  writeFile(fixture, 'ui/components/blocked-file.tsx', "const color = '#123';\n");

  const { json } = runAudit(fixture, ['--json'], {
    nodeArgs: ['--require', preload],
  });

  assert.equal(json.files, 1);
  assert.equal(json.matches, 1);
  assert.equal(json.offenders[0].file, 'components/valid.tsx');
  assert.equal(json.scanComplete, false);
  assert.equal(json.scanErrors.length, 1);
  assert.match(json.scanErrors[0], /Unable to read file/);
  assert.match(json.scanErrors[0], /blocked-file\.tsx/);
});

test('text output warns about incomplete scans on stderr', (t) => {
  const fixture = createFixture(t);
  const preload = addFilesystemFailureMock(fixture);

  writeFile(fixture, 'ui/components/valid.tsx', "const color = '#fff';\n");
  writeFile(fixture, 'ui/components/blocked-file.tsx', "const color = '#123';\n");

  const result = runAudit(fixture, [], {
    nodeArgs: ['--require', preload],
  });

  assert.match(result.stderr, /WARNING: Scan incomplete/);
  assert.match(result.stderr, /Unable to read file/);
  assert.match(result.stdout, /AUDIT hex files=1 matches=1/);
});

test('ignores RGB functions in comments but counts actual functions', (t) => {
  const fixture = createFixture(t);

  writeFile(
    fixture,
    'ui/components/rgb-comments.tsx',
    [
      '// rgb(255, 0, 0)',
      '/* rgba(0, 0, 0, .5) */',
      'const color = rgb(12, 34, 56);',
      'const alpha = rgba(12, 34, 56, .5);',
      '',
    ].join('\n'),
  );

  const { json: result } = runAudit(fixture);

  assert.equal(result.files, 1);
  assert.equal(result.rgb, 2);
  assert.equal(result.matches, 2);
  assert.equal(result.offenders[0].rgb, 2);
});

test('scans nested production directories with root-scoped names', (t) => {
  const fixture = createFixture(t);

  const directoryNames = ['assets', 'lib', 'public', 'theme', 'tests'];

  for (const directory of directoryNames) {
    writeFile(
      fixture,
      `ui/components/feature/${directory}/source.tsx`,
      'const color = "#252E31";\n',
    );

    // These UI-root directories are intentionally outside SCAN_DIRS.
    writeFile(fixture, `ui/${directory}/root-source.tsx`, 'const color = "#000000";\n');
  }

  const { json: result } = runAudit(fixture);

  assert.equal(result.files, directoryNames.length);
  assert.equal(result.matches, directoryNames.length);
  assert.deepEqual(
    result.offenders.map((offender) => offender.file).sort(),
    directoryNames.map((directory) => `components/feature/${directory}/source.tsx`).sort(),
  );
});

test('excludes nested node_modules directories at every depth', (t) => {
  const fixture = createFixture(t);

  writeFile(fixture, 'ui/components/feature/source.tsx', 'const color = "#252E31";\n');
  writeFile(
    fixture,
    'ui/components/feature/node_modules/package/source.tsx',
    'const color = "#000000";\n',
  );

  const { json: result } = runAudit(fixture);

  assert.equal(result.files, 1);
  assert.equal(result.matches, 1);
  assert.deepEqual(
    result.offenders.map((offender) => offender.file),
    ['components/feature/source.tsx'],
  );
});
