#!/usr/bin/env node
/**
 * meshery-ui color-literal audit
 *
 * Counts occurrences of hex colors (#RGB, #RGBA, #RRGGBB, #RRGGBBAA) and
 * rgb()/rgba() function literals in .ts/.tsx/.js/.jsx files. These
 * literals are being eliminated in favor of theme.palette.* (which comes
 * from @sistent/sistent) as part of the UI restructure.
 *
 * This script is broader than the companion ESLint rule
 * (no-restricted-syntax in ui/eslint.config.js), which only catches hex
 * strings that are exact Literal values. This audit also catches hex
 * colors embedded inside longer strings, like:
 *
 *   const x = { border: '1px solid #F91313' };
 *
 * Usage:
 *   node scripts/audit-hex.js           # summary + top 20 offenders
 *   node scripts/audit-hex.js --all     # print every offending file
 *   node scripts/audit-hex.js --json    # structured JSON to stdout
 *
 * Exit code is always 0.
 *
 * Final line of stdout is always:
 *   AUDIT hex files=<N> matches=<M>
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
  'store',
  'rtk-query',
  'machines',
  'hooks',
  'api',
  'css',
];

// Skipped entirely. These are the same exclusions applied in
// eslint.config.js to the hex/rgba no-restricted-syntax rule: the theme
// module, legacy theme/constants directories (scheduled for deletion),
// SVG icon components, third-party integration helpers, static assets,
// and tests.
const IGNORED_DIRS = new Set([
  'node_modules',
  '.next',
  'out',
  'playground',
  'playwright-report',
  'test-results',
  '__generated__',
  // Scoped opt-outs matching the ESLint rule
  'theme',
  'themes',
  'assets',
  'constants',
  'lib',
  'public',
  'tests',
]);

const EXTENSIONS = new Set(['.ts', '.tsx', '.js', '.jsx']);

// Hex colors: #RGB, #RGBA, #RRGGBB, #RRGGBBAA. Valid CSS hex lengths only.
// Uses a word boundary so #abcde (5 chars) does not match as #abc.
const HEX_PATTERN = /#(?:[0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})\b/g;

// rgb() / rgba() function calls inside strings.
const RGB_PATTERN = /rgba?\(/g;

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

function countMatches(src, pattern) {
  pattern.lastIndex = 0;
  let count = 0;
  while (pattern.exec(src) !== null) count++;
  return count;
}

function scanFile(file) {
  const src = fs.readFileSync(file, 'utf8');
  const hex = countMatches(src, HEX_PATTERN);
  const rgb = countMatches(src, RGB_PATTERN);
  return { hex, rgb };
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
  let totalHex = 0;
  let totalRgb = 0;
  for (const file of files) {
    const { hex, rgb } = scanFile(file);
    if (hex === 0 && rgb === 0) continue;
    const rel = path.relative(UI_ROOT, file);
    offenders.push({ file: rel, hex, rgb, total: hex + rgb });
    totalHex += hex;
    totalRgb += rgb;
  }
  offenders.sort((a, b) => b.total - a.total || a.file.localeCompare(b.file));

  const totalMatches = totalHex + totalRgb;

  if (asJson) {
    process.stdout.write(
      JSON.stringify(
        {
          audit: 'hex',
          files: offenders.length,
          matches: totalMatches,
          hex: totalHex,
          rgb: totalRgb,
          offenders,
        },
        null,
        2,
      ) + '\n',
    );
    return;
  }

  process.stdout.write('meshery-ui color-literal audit\n');
  process.stdout.write('==============================\n');
  process.stdout.write('Patterns: #RGB/#RGBA/#RRGGBB/#RRGGBBAA, rgb(), rgba()\n');
  process.stdout.write(`Scope:    ${SCAN_DIRS.join(', ')}\n`);
  process.stdout.write(
    'Ignored:  theme, themes, assets, constants, lib, public, tests, __generated__\n\n',
  );

  const limit = showAll ? offenders.length : Math.min(20, offenders.length);
  if (offenders.length === 0) {
    process.stdout.write('No color literals found outside the ignored paths. 🎯\n\n');
  } else {
    const header = showAll ? 'Offenders (hex + rgb count)' : `Top ${limit} files (hex + rgb count)`;
    process.stdout.write(`${header}:\n`);
    for (let i = 0; i < limit; i++) {
      const { file, hex, rgb } = offenders[i];
      const tag = `${hex}h+${rgb}r`.padStart(8);
      process.stdout.write(`  ${tag}  ${file}\n`);
    }
    if (!showAll && offenders.length > limit) {
      process.stdout.write(`  ... and ${offenders.length - limit} more (use --all to see them)\n`);
    }
    process.stdout.write('\n');
  }

  process.stdout.write('Summary:\n');
  process.stdout.write(`  Files:   ${offenders.length}\n`);
  process.stdout.write(`  Hex:     ${totalHex}\n`);
  process.stdout.write(`  rgb/a:   ${totalRgb}\n`);
  process.stdout.write(`  Total:   ${totalMatches}\n\n`);
  process.stdout.write(
    `AUDIT hex files=${offenders.length} matches=${totalMatches} hex=${totalHex} rgb=${totalRgb}\n`,
  );
}

main();
