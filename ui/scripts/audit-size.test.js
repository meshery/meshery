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
  assert.deepEqual(result.files, [
    { file: 'components/jsx.tsx', lines: 607 },
    { file: 'components/template.js', lines: 607 },
  ]);
});

test('counts code sharing a line with comments and regex literals containing slashes', (t) => {
  const fixture = createFixture(t);
  const mixedLines = Array.from(
    { length: 303 },
    (_, i) =>
      `const regex${i} = /[//]/; // trailing comment\n/* leading comment */ const value${i} = ${i};`,
  );

  writeFile(fixture, 'ui/components/mixed.js', `${mixedLines.join('\n')}\n`);

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 1);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files, [{ file: 'components/mixed.js', lines: 606 }]);
});

test('handles multiline block comments and code after their closing delimiter', (t) => {
  const fixture = createFixture(t);
  const blockCodeEntries = Array.from(
    { length: 605 },
    (_, i) => `/* documentation\n * continues\n */ const value${i} = ${i};`,
  );

  writeFile(
    fixture,
    'ui/components/block-comments.js',
    `/*\n${Array.from({ length: 605 }, (_, i) => ` * documentation ${i + 1}`).join('\n')}\n */\nconst value = 1;\n`,
  );

  writeFile(fixture, 'ui/components/block-comment-code.js', `${blockCodeEntries.join('\n')}\n`);

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 1);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files, [{ file: 'components/block-comment-code.js', lines: 605 }]);
});

test('counts JSX text in JavaScript files and regex literals containing comment-like text', (t) => {
  const fixture = createFixture(t);
  const jsxLines = Array.from({ length: 605 }, (_, i) => `// JSX text ${i + 1}`);
  const regexLines = Array.from({ length: 605 }, (_, i) => `const pattern${i} = /[//]/;`);

  writeFile(
    fixture,
    'ui/components/jsx.js',
    `const view = <div>\n${jsxLines.join('\n')}\n{title}\n</div>;\n`,
  );

  writeFile(fixture, 'ui/components/regex.js', `${regexLines.join('\n')}\n`);

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 2);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files, [
    { file: 'components/jsx.js', lines: 608 },
    { file: 'components/regex.js', lines: 605 },
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

test('skips comments before closing delimiters, keywords, JSX tags, and EOF', (t) => {
  const fixture = createFixture(t);
  const snippetLines = [
    "import {\n  foo,\n  // after last import\n} from 'pkg';",
    'interface Props {\n  a: string;\n  // after last prop\n}',
    'interface EmptyInterface {\n  // empty interface\n}',
    'type Box = {\n  x: number;\n  // after last type member\n};',
    'enum Kind {\n  A,\n  // after last enum member\n}',
    'class Widget {\n  run() {\n    return 1;\n    // after last statement\n  }\n  // after last class member\n}',
    'function emptyFn(\n  // empty params\n) {\n  // empty block\n}',
    'const obj = {\n  k: 1,\n  // after last property\n};',
    'const emptyObj = {\n  // empty object\n};',
    'const arr = [\n  1,\n  // after last element\n];',
    'const emptyArr = [\n  // empty array\n];',
    'call(\n  1,\n  // after last arg\n);',
    'call(\n  // empty args\n);',
    'if (cond) {\n  a();\n}\n// before else\nelse {\n  b();\n}',
    'try {\n  a();\n}\n// before catch\ncatch (err) {\n  b(err);\n}\n// before finally\nfinally {\n  c();\n}',
    'switch (v) {\n  case 1:\n    a();\n    // after case stmt\n  default:\n    // inside default\n}',
    'const view = (\n  <div\n    id="x"\n    // before >\n  >\n    {/*\n      multiline jsx comment\n    */}\n    {\n      // line jsx comment\n    }\n    <span\n      // before />\n    />\n  </div>\n);',
  ];
  const padding = Array.from({ length: 560 }, (_, i) => `const pad${i} = ${i};`);
  const source = `${snippetLines.join('\n')}\n${padding.join('\n')}\n// trailing eof comment\n`;

  writeFile(fixture, 'ui/components/punctuation-comments.tsx', source);

  const result = runAudit(fixture);
  const expectedSnippetCodeLines = 69;
  const expectedTotal = expectedSnippetCodeLines + padding.length;

  assert.equal(result.over_warn, 1);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files, [
    { file: 'components/punctuation-comments.tsx', lines: expectedTotal },
  ]);
});

test('does not treat unclosed or inline comment tokens in JSX text as comments', (t) => {
  const fixture = createFixture(t);
  const afterLines = Array.from({ length: 605 }, (_, i) => `export const after${i} = ${i};`);

  writeFile(
    fixture,
    'ui/components/jsx-slash-text.jsx',
    `const first = <div>{base} /* glob\n</div>;\nconst second = <div>{base} // text</div>;\n${afterLines.join('\n')}\n`,
  );

  const result = runAudit(fixture);

  assert.equal(result.over_warn, 1);
  assert.equal(result.over_hard, 0);
  assert.deepEqual(result.files, [{ file: 'components/jsx-slash-text.jsx', lines: 608 }]);
});
