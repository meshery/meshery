'use strict';

const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { test } = require('node:test');

const productionScript = path.join(__dirname, 'audit-size.js');

function writeFile(root, relativePath, contents) {
  const file = path.join(root, relativePath);
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, contents);
}

function createFixture(t) {
  const fixture = fs.mkdtempSync(path.join(os.tmpdir(), 'audit-size-'));
  t.after(() => fs.rmSync(fixture, { recursive: true, force: true }));

  fs.mkdirSync(path.join(fixture, 'ui', 'scripts'), { recursive: true });
  fs.copyFileSync(productionScript, path.join(fixture, 'ui', 'scripts', 'audit-size.js'));

  return fixture;
}

function runAudit(fixture) {
  const script = path.join(fixture, 'ui', 'scripts', 'audit-size.js');
  const typescriptRoot = path.resolve(path.dirname(require.resolve('typescript')), '..', '..');

  const result = spawnSync(process.execPath, [script, '--json'], {
    encoding: 'utf8',
    env: {
      ...process.env,
      NODE_PATH: [typescriptRoot, process.env.NODE_PATH].filter(Boolean).join(path.delimiter),
    },
  });

  assert.equal(result.status, 0, result.stderr);
  return JSON.parse(result.stdout);
}

test('counts comment-like template and JSX text while skipping actual comments', (t) => {
  const fixture = createFixture(t);
  const templateLines = Array.from({ length: 605 }, (_, i) => `// template content ${i + 1}`);
  const jsxLines = Array.from({ length: 605 }, (_, i) => `// JSX content ${i + 1}`);

  writeFile(
    fixture,
    'ui/components/template.js',
    `const template = \`\n${templateLines.join('\n')}\n\`;\n`,
  );

  writeFile(
    fixture,
    'ui/components/jsx.tsx',
    `const view = <div>\n${jsxLines.join('\n')}\n</div>;\n`,
  );

  writeFile(
    fixture,
    'ui/components/comments.js',
    `${Array.from({ length: 605 }, (_, i) => `// actual comment ${i + 1}`).join('\n')}\nconst value = 1;\n`,
  );

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 2);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files.map(({ file }) => file).sort(), [
    'components/jsx.tsx',
    'components/template.js',
  ]);
});

test('counts code sharing a line with comments and regex literals containing slashes', (t) => {
  const fixture = createFixture(t);

  writeFile(
    fixture,
    'ui/components/mixed.js',
    'const regex = /[//]/; // trailing comment\n/* leading comment */ const value = 1;\n',
  );

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 0);
  assert.equal(result.over_hard, 0);
});

test('handles multiline block comments and code after their closing delimiter', (t) => {
  const fixture = createFixture(t);

  writeFile(
    fixture,
    'ui/components/block-comments.js',
    `/*\n${Array.from({ length: 605 }, (_, i) => ` * documentation ${i + 1}`).join('\n')}\n */\nconst value = 1;\n`,
  );

  writeFile(
    fixture,
    'ui/components/block-comment-code.js',
    `/* documentation\n * continues\n */ const value = 1;\n`,
  );

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 0);
  assert.equal(result.over_hard, 0);
});

test('counts JSX text in JavaScript files and regex literals containing comment-like text', (t) => {
  const fixture = createFixture(t);
  const jsxLines = Array.from({ length: 605 }, (_, i) => `// JSX text ${i + 1}`);
  const regexLines = Array.from({ length: 605 }, () => 'const pattern = /[//]/;');

  writeFile(
    fixture,
    'ui/components/jsx.js',
    `const view = <div>\n${jsxLines.join('\n')}\n</div>;\n`,
  );

  writeFile(fixture, 'ui/components/regex.js', `${regexLines.join('\n')}\n`);

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 2);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files.map(({ file }) => file).sort(), [
    'components/jsx.js',
    'components/regex.js',
  ]);
});

test('counts code on lines containing block comments', (t) => {
  const fixture = createFixture(t);
  const lines = Array.from(
    { length: 605 },
    (_, i) => `/* comment ${i + 1} */ const value${i} = ${i};`,
  );

  writeFile(fixture, 'ui/components/inline-block-comments.js', `${lines.join('\n')}\n`);

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 1);
  assert.equal(result.over_hard, 0);
  assert.equal(result.files[0].file, 'components/inline-block-comments.js');
  assert.equal(result.files[0].lines, 605);
});
