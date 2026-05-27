#!/usr/bin/env node
/**
 * Post-processes Relay-generated TypeScript artifacts to remove the
 * // @ts-nocheck directive that relay-compiler v20 unconditionally emits.
 *
 * The generated code is valid TypeScript without the directive, so this step
 * is safe and survives strict-mode tsc checks.  Run automatically as part of
 * `npm run relay` so that the suppression never re-appears on rebuild.
 */

const { globSync } = require('node:fs');
const fs = require('node:fs');
const path = require('node:path');

// Match all __generated__ directories under ui/ (excludes node_modules automatically)
const files = globSync('**/__generated__/*.graphql.ts', {
  cwd: path.join(__dirname, '..'),
  exclude: ['**/node_modules/**'],
});

let changed = 0;
for (const rel of files) {
  const abs = path.join(__dirname, '..', rel);
  const original = fs.readFileSync(abs, 'utf8');
  // Remove the @ts-nocheck line; handle both LF and CRLF line endings
  const patched = original.replace(/^\/\/ @ts-nocheck\r?\n/m, '');
  if (patched !== original) {
    fs.writeFileSync(abs, patched, 'utf8');
    changed++;
  }
}

console.log(
  `relay-strip-ts-nocheck: removed @ts-nocheck from ${changed}/${files.length} generated file(s).`,
);
