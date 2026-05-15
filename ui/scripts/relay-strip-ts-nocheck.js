#!/usr/bin/env node
/**
 * Post-processes Relay-generated TypeScript artifacts to remove the
 * // @ts-nocheck directive that relay-compiler v20 unconditionally emits.
 *
 * The generated code is valid TypeScript without the directive, so this step
 * is safe and survives strict-mode tsc checks.  Run automatically as part of
 * `npm run relay` so that the suppression never re-appears on rebuild.
 */

const fs = require('fs');
const path = require('path');
const { globSync } = require('glob');

const generatedGlob = 'graphql/{queries,mutations,subscriptions}/__generated__/*.graphql.ts';
const files = globSync(generatedGlob, { cwd: path.join(__dirname, '..') });

let changed = 0;
for (const rel of files) {
  const abs = path.join(__dirname, '..', rel);
  const original = fs.readFileSync(abs, 'utf8');
  // Remove the @ts-nocheck line (relay-compiler emits it on its own line)
  const patched = original.replace(/^\/\/ @ts-nocheck\n/m, '');
  if (patched !== original) {
    fs.writeFileSync(abs, patched, 'utf8');
    changed++;
  }
}

console.log(
  `relay-strip-ts-nocheck: removed @ts-nocheck from ${changed}/${files.length} generated file(s).`,
);
