#!/usr/bin/env node
/**
 * meshery-ui @mui import audit
 *
 * Counts files that still import from @mui/*, @material-ui/*, or @rjsf/mui.
 * These imports are being eliminated in favor of @sistent/sistent — the
 * design system — as part of the UI restructure.
 *
 * The companion ESLint rule (no-restricted-imports in ui/eslint.config.js)
 * warns developers when they add a new offending import; this script is
 * the trend tracker that CI uses to watch the count go down.
 *
 * Usage:
 *   node scripts/audit-mui.js            # summary + top 20 offenders
 *   node scripts/audit-mui.js --all      # print every offending file
 *   node scripts/audit-mui.js --json     # emit structured JSON to stdout
 *
 * Exit code is always 0 — this is an informational audit, not a gate.
 *
 * The final line of stdout is always:
 *   AUDIT mui files=<N> matches=<M>
 * …so CI can grep it for trend tracking without parsing the rest.
 */

'use strict';

const fs = require('fs');
const path = require('path');

const UI_ROOT = path.resolve(__dirname, '..');

// Directories under ui/ that are scanned.
const SCAN_DIRS = [
  'components',
  'pages',
  'utils',
  'themes',
  'theme',
  'machines',
  'store',
  'rtk-query',
  'lib',
  'hooks',
  'api',
];

// Directories (absolute or ui-relative) that are never scanned.
const IGNORED_DIRS = new Set([
  'node_modules',
  '.next',
  'out',
  'playground',
  'playwright-report',
  'test-results',
  '__generated__',
]);

const EXTENSIONS = new Set(['.ts', '.tsx', '.js', '.jsx']);

// Match any import/require from @mui/*, @material-ui/*, or @rjsf/mui.
// Covers:
//   import X from '@mui/material'
//   import { X } from '@mui/material'
//   import '@mui/material/Button'
//   const X = require('@mui/material')
const IMPORT_PATTERN =
  /(?:from\s+|require\(\s*)['"](@mui\/[^'"]+|@material-ui\/[^'"]+|@rjsf\/mui)['"]/g;

function walk(dir, out) {
  let entries;
  try {
    entries = fs.readdirSync(dir, { withFileTypes: true });
  } catch {
    return;
  }
  for (const entry of entries) {
    if (IGNORED_DIRS.has(entry.name)) continue;
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walk(full, out);
    } else if (entry.isFile() && EXTENSIONS.has(path.extname(entry.name))) {
      out.push(full);
    }
  }
}

function scanFile(file) {
  const src = fs.readFileSync(file, 'utf8');
  const matches = [];
  let m;
  IMPORT_PATTERN.lastIndex = 0;
  while ((m = IMPORT_PATTERN.exec(src)) !== null) {
    matches.push(m[1]);
  }
  return matches;
}

function main() {
  const args = process.argv.slice(2);
  const showAll = args.includes('--all');
  const asJson = args.includes('--json');

  const files = [];
  for (const dir of SCAN_DIRS) {
    walk(path.join(UI_ROOT, dir), files);
  }

  const offenders = [];
  let totalMatches = 0;
  for (const file of files) {
    const matches = scanFile(file);
    if (matches.length === 0) continue;
    const rel = path.relative(UI_ROOT, file);
    offenders.push({ file: rel, count: matches.length, imports: matches });
    totalMatches += matches.length;
  }
  offenders.sort((a, b) => b.count - a.count || a.file.localeCompare(b.file));

  if (asJson) {
    process.stdout.write(
      JSON.stringify(
        {
          audit: 'mui',
          files: offenders.length,
          matches: totalMatches,
          offenders,
        },
        null,
        2,
      ) + '\n',
    );
    return;
  }

  process.stdout.write('meshery-ui @mui/* import audit\n');
  process.stdout.write('==============================\n');
  process.stdout.write('Patterns: @mui/*, @material-ui/*, @rjsf/mui\n');
  process.stdout.write(`Scope:    ${SCAN_DIRS.join(', ')}\n\n`);

  const limit = showAll ? offenders.length : Math.min(20, offenders.length);
  if (offenders.length === 0) {
    process.stdout.write('No @mui/* imports found. The restructure is done.\n\n');
  } else {
    const header = showAll ? 'Offenders' : `Top ${limit} files (by match count)`;
    process.stdout.write(`${header}:\n`);
    for (let i = 0; i < limit; i++) {
      const { file, count } = offenders[i];
      process.stdout.write(`  ${String(count).padStart(3)}  ${file}\n`);
    }
    if (!showAll && offenders.length > limit) {
      process.stdout.write(`  ... and ${offenders.length - limit} more (use --all to see them)\n`);
    }
    process.stdout.write('\n');
  }

  process.stdout.write('Summary:\n');
  process.stdout.write(`  Files:   ${offenders.length}\n`);
  process.stdout.write(`  Matches: ${totalMatches}\n\n`);
  process.stdout.write(`AUDIT mui files=${offenders.length} matches=${totalMatches}\n`);
}

main();
