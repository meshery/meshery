#!/usr/bin/env node
/**
 * meshery-ui component-size audit
 *
 * Counts .ts/.tsx/.js/.jsx source files that exceed the component size
 * budget set out in the UI restructure plan:
 *
 *   - Warn threshold:   600 lines
 *   - Hard ceiling:    1000 lines
 *
 * The companion ESLint rule (max-lines in ui/eslint.config.js) already
 * emits a warning when a file crosses 1000 lines. This audit is the
 * trend tracker that CI uses to watch both thresholds go down as the
 * eight giant files documented in the plan get split apart in phase 5.
 *
 * Line counting matches ESLint's max-lines behaviour with
 * { skipComments: true, skipBlankLines: true } — blank lines and
 * single-/multi-line comment lines are not counted toward the total.
 *
 * Usage:
 *   node scripts/audit-size.js          # summary + top 20 largest
 *   node scripts/audit-size.js --all    # print every file over the warn threshold
 *   node scripts/audit-size.js --json   # structured JSON to stdout
 *
 * Exit code is always 0.
 *
 * Final line of stdout is always:
 *   AUDIT size over_warn=<N> over_hard=<M>
 */

'use strict';

const fs = require('fs');
const path = require('path');

const UI_ROOT = path.resolve(__dirname, '..');

const WARN_THRESHOLD = 600;
const HARD_THRESHOLD = 1000;

// Directories under ui/ that are scanned for oversized source files.
const SCAN_DIRS = [
  'components',
  'pages',
  'utils',
  'store',
  'rtk-query',
  'machines',
  'hooks',
  'api',
  'themes',
  'theme',
  'lib',
];

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

// Count lines the way ESLint's max-lines does with
// { skipComments: true, skipBlankLines: true }. A line is skipped if
// it is blank after trimming, or if it lies inside a /* ... */ block,
// or if it starts with // after whitespace.
function countSignificantLines(src) {
  const lines = src.split(/\r?\n/);
  let count = 0;
  let inBlockComment = false;
  for (const rawLine of lines) {
    let line = rawLine.trim();
    if (line === '') continue;

    if (inBlockComment) {
      const end = line.indexOf('*/');
      if (end === -1) continue;
      line = line.slice(end + 2).trim();
      inBlockComment = false;
      if (line === '') continue;
    }

    // Strip a leading // comment.
    if (line.startsWith('//')) continue;

    // Strip a leading /* ... */ that begins and ends on this line.
    while (line.startsWith('/*')) {
      const end = line.indexOf('*/', 2);
      if (end === -1) {
        inBlockComment = true;
        line = '';
        break;
      }
      line = line.slice(end + 2).trim();
    }
    if (line === '') continue;

    count++;
  }
  return count;
}

function main() {
  const args = process.argv.slice(2);
  const showAll = args.includes('--all');
  const asJson = args.includes('--json');

  const files = [];
  for (const dir of SCAN_DIRS) {
    walk(path.join(UI_ROOT, dir), files);
  }

  const measured = [];
  for (const file of files) {
    const src = fs.readFileSync(file, 'utf8');
    const lines = countSignificantLines(src);
    if (lines < WARN_THRESHOLD) continue;
    measured.push({ file: path.relative(UI_ROOT, file), lines });
  }
  measured.sort((a, b) => b.lines - a.lines || a.file.localeCompare(b.file));

  const overHard = measured.filter((m) => m.lines >= HARD_THRESHOLD);
  const overWarn = measured; // already filtered to >= WARN_THRESHOLD

  if (asJson) {
    process.stdout.write(
      JSON.stringify(
        {
          audit: 'size',
          warn_threshold: WARN_THRESHOLD,
          hard_threshold: HARD_THRESHOLD,
          over_warn: overWarn.length,
          over_hard: overHard.length,
          files: measured,
        },
        null,
        2,
      ) + '\n',
    );
    return;
  }

  process.stdout.write('meshery-ui component-size audit\n');
  process.stdout.write('===============================\n');
  process.stdout.write(
    `Thresholds: warn=${WARN_THRESHOLD}, hard=${HARD_THRESHOLD} (blank lines and comments excluded)\n`,
  );
  process.stdout.write(`Scope:      ${SCAN_DIRS.join(', ')}\n\n`);

  if (overWarn.length === 0) {
    process.stdout.write(`No files over ${WARN_THRESHOLD} lines. The split is done.\n\n`);
  } else {
    const limit = showAll ? overWarn.length : Math.min(20, overWarn.length);
    const header = showAll ? `Files over ${WARN_THRESHOLD} lines` : `Top ${limit} largest files`;
    process.stdout.write(`${header}:\n`);
    for (let i = 0; i < limit; i++) {
      const { file, lines } = overWarn[i];
      const flag = lines >= HARD_THRESHOLD ? '!!' : '  ';
      process.stdout.write(`  ${flag} ${String(lines).padStart(5)}  ${file}\n`);
    }
    if (!showAll && overWarn.length > limit) {
      process.stdout.write(`  ... and ${overWarn.length - limit} more (use --all to see them)\n`);
    }
    process.stdout.write('\n');
  }

  process.stdout.write('Summary:\n');
  process.stdout.write(`  Over ${WARN_THRESHOLD}:  ${overWarn.length}\n`);
  process.stdout.write(`  Over ${HARD_THRESHOLD}: ${overHard.length}\n\n`);
  process.stdout.write(`AUDIT size over_warn=${overWarn.length} over_hard=${overHard.length}\n`);
}

main();
