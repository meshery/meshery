#!/usr/bin/env node
/**
 * meshery @mui import audit
 *
 * Counts files that still import from @mui/*, @material-ui/*, or @rjsf/mui.
 * These imports are being eliminated in favor of @sistent/sistent — the
 * design system — as part of the UI restructure.
 *
 * The companion ESLint rule (no-restricted-imports in ui/eslint.config.js)
 * warns developers when they add a new offending import; this script is
 * the trend tracker that CI uses to watch the count go down.
 *
 * Scans two roots:
 *   - ui/                          (main Next.js app)
 *   - install/docker-extension/ui/ (Docker Desktop extension sub-app)
 *
 * Each scan walks a list of `dirs` and also picks up explicit `files`
 * (top-level config files that no `dirs` entry would otherwise reach).
 * Each scan has its own `allowlist` of files whose @mui imports are
 * intentional; allowlisted files are excluded from the offender count
 * and from the `--fail-on-new` gate.
 *
 * Usage:
 *   node scripts/audit-mui.js                 # summary + top 20 offenders
 *   node scripts/audit-mui.js --all           # print every offending file
 *   node scripts/audit-mui.js --json          # emit structured JSON to stdout
 *   node scripts/audit-mui.js --fail-on-new   # exit 1 if any non-allowlisted hit
 *
 * Default exit code is 0 — this is an informational audit. The --fail-on-new
 * flag turns it into a gate; intended for CI once the count is stable.
 *
 * The final line of stdout is always:
 *   AUDIT mui files=<N> matches=<M>
 * …so CI can grep it for trend tracking without parsing the rest.
 */

'use strict';

const fs = require('fs');
const path = require('path');

const REPO_ROOT = path.resolve(__dirname, '..', '..');

// A scan is a {name, root, dirs, allowlist} entry. Each scan has its own
// root directory, its own list of directories to walk, and its own
// allowlist of files whose @mui imports are intentional (so they're
// excluded from the offender count and from the --fail-on-new gate).
const SCANS = [
  {
    name: 'ui',
    root: path.join(REPO_ROOT, 'ui'),
    dirs: [
      'components',
      'pages',
      'utils',
      'theme',
      'machines',
      'store',
      'rtk-query',
      'lib',
      'api',
      // Added 2026-05-17: ui/public/ holds drawer-icon hover React components
      // that import from @mui/icons-material. Without this entry the audit
      // missed an entire class of imports.
      'public',
    ],
    // Top-level config files that aren't under any scanned `dirs` entry but
    // can still hold @mui imports. The next.config.js webpack alias keys
    // (e.g. `'@mui/material': path.resolve(...)`) don't match the audit
    // regex (no `from`/`require(` prefix), so this list is currently just
    // remote-component.config.js, which has many real require() calls.
    // Listed by relative path; resolved from `root`.
    files: ['remote-component.config.js'],
    // Approved single-boundary wrappers around @mui/* and @rjsf/mui.
    //
    // These files are the *only* places in ui/ allowed to import from the
    // legacy MUI / RJSF packages, per the UI restructure plan §1.3 and the
    // Phase 3 issue (#18657). They exist because Sistent doesn't yet expose
    // the underlying primitives (`GlobalStyles`, `darken`, the date pickers,
    // the tree view, the RJSF theme adapter); every other consumer goes
    // through `@/theme` or the shared component wrappers, not `@mui/*`
    // directly.
    //
    // Adding a new entry should require an open issue and a code review —
    // the goal is to keep this list small and to drain it as Sistent grows.
    allowlist: new Set([
      'theme/index.ts',
      'components/shared/TreeView/TreeView.tsx',
      'components/shared/DatePicker/index.ts',
      'components/shared/DatePicker/MesheryDateTimePicker.tsx',
      'components/shared/FormFields/RJSFProvider.tsx',
      // remote-component.config.js resolves @mui/icons-material for the
      // @paciolan/remote-component plugin system; the require() entries
      // hand the host's icon implementations to remote components. Not
      // migratable without coordinating remote-component consumers.
      'remote-component.config.js',
    ]),
  },
  {
    name: 'docker-extension',
    root: path.join(REPO_ROOT, 'install', 'docker-extension', 'ui'),
    dirs: ['src'],
    // Documented exceptions for this sub-app. Sistent 0.16.10 (pinned here,
    // intentionally not bumped) ships its own nested `@mui/material@6.5.0`
    // under `node_modules/@sistent/sistent/node_modules/@mui/`. Because
    // `DockerMuiThemeProvider` resolves the root `@mui/material@^7`, the
    // two MUI instances do not share a React Context — Sistent primitives
    // rendered outside `SistentThemeProviderWithoutBaseLine` will not pick
    // up the Docker Desktop theme. Until this sub-app's Sistent is bumped
    // to a release that shares MUI with `@docker/docker-mui-theme` (see
    // design spec §7.1 follow-up), the docker-extension intentionally
    // imports MUI directly. These allowlisted files are exactly that.
    allowlist: new Set([
      'src/components/Catalog/CatalogCard.js',
      'src/components/Catalog/style.js',
      'src/components/ExtensionComponent/ExtensionComponent.js',
      'src/components/ExtensionComponent/styledComponents.js',
      'src/components/Walkthrough/Tour.js',
      'src/components/Walkthrough/tourStyledComponents.js',
    ]),
  },
];

// Directories (basename match) that are never scanned, in any root.
const IGNORED_DIRS = new Set([
  'node_modules',
  '.next',
  'out',
  'playground',
  'playwright-report',
  'test-results',
  '__generated__',
  'build',
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
  } catch (err) {
    // Always report — this function is only called for subdirectories
    // discovered during recursion, so any read failure is unexpected.
    // The "is this configured scan dir missing?" check happens in
    // runScan() *before* walk() runs.
    process.stderr.write(`audit-mui: ${dir}: ${err.message}\n`);
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

function runScan(scan) {
  const files = [];
  for (const dir of scan.dirs) {
    const full = path.join(scan.root, dir);
    // A configured scan dir that doesn't exist is a config error, not
    // expected absence. Warn loudly — otherwise a typo (or a renamed
    // directory) would silently drop a whole scan target and `--fail-on-new`
    // would still pass. This is exactly the audit blind-spot pattern this
    // script exists to prevent.
    if (!fs.existsSync(full)) {
      process.stderr.write(`audit-mui: missing scan dir ${full}\n`);
      continue;
    }
    walk(full, files);
  }
  // Explicit top-level files (e.g., config files that live at the scan
  // root and aren't reached by walking any `dirs` entry).
  for (const file of scan.files ?? []) {
    const full = path.join(scan.root, file);
    if (!fs.existsSync(full)) {
      process.stderr.write(`audit-mui: missing scan file ${full}\n`);
      continue;
    }
    files.push(full);
  }

  const offenders = [];
  let totalMatches = 0;
  let skippedAllowlisted = 0;
  for (const file of files) {
    const matches = scanFile(file);
    if (matches.length === 0) continue;
    const rel = path.relative(scan.root, file);
    if (scan.allowlist.has(rel)) {
      skippedAllowlisted += 1;
      continue;
    }
    offenders.push({ scan: scan.name, file: rel, count: matches.length, imports: matches });
    totalMatches += matches.length;
  }
  offenders.sort((a, b) => b.count - a.count || a.file.localeCompare(b.file));
  return {
    name: scan.name,
    files: offenders.length,
    matches: totalMatches,
    skippedAllowlisted,
    allowlistSize: scan.allowlist.size,
    offenders,
  };
}

function main() {
  const args = process.argv.slice(2);
  const showAll = args.includes('--all');
  const asJson = args.includes('--json');
  const failOnNew = args.includes('--fail-on-new');

  const results = SCANS.map(runScan);
  const totalFiles = results.reduce((acc, r) => acc + r.files, 0);
  const totalMatches = results.reduce((acc, r) => acc + r.matches, 0);
  const totalAllowlist = results.reduce((acc, r) => acc + r.allowlistSize, 0);
  // Sort across all scans so the "Top N by match count" header is honest
  // when there are offenders in more than one scan.
  const allOffenders = results
    .flatMap((r) => r.offenders)
    .sort((a, b) => b.count - a.count || a.file.localeCompare(b.file));

  if (asJson) {
    process.stdout.write(
      JSON.stringify(
        {
          audit: 'mui',
          files: totalFiles,
          matches: totalMatches,
          scans: results.map((r) => ({
            name: r.name,
            files: r.files,
            matches: r.matches,
            allowlistSize: r.allowlistSize,
            allowlistedFilesWithMuiImports: r.skippedAllowlisted,
            offenders: r.offenders,
          })),
        },
        null,
        2,
      ) + '\n',
    );
    process.exit(failOnNew && totalFiles > 0 ? 1 : 0);
  }

  process.stdout.write('meshery @mui/* import audit\n');
  process.stdout.write('===========================\n');
  process.stdout.write('Patterns: @mui/*, @material-ui/*, @rjsf/mui\n');
  for (const r of results) {
    const scan = SCANS.find((s) => s.name === r.name);
    process.stdout.write(`Scan [${r.name}]: ${scan.dirs.join(', ')}\n`);
  }
  process.stdout.write('\n');

  const limit = showAll ? allOffenders.length : Math.min(20, allOffenders.length);
  if (allOffenders.length === 0) {
    process.stdout.write(
      'No @mui/* imports found outside the allowlists. The restructure is done.\n\n',
    );
  } else {
    const header = showAll ? 'Offenders' : `Top ${limit} files (by match count)`;
    process.stdout.write(`${header}:\n`);
    for (let i = 0; i < limit; i++) {
      const { scan: scanName, file, count } = allOffenders[i];
      process.stdout.write(`  ${String(count).padStart(3)}  [${scanName}] ${file}\n`);
    }
    if (!showAll && allOffenders.length > limit) {
      process.stdout.write(
        `  ... and ${allOffenders.length - limit} more (use --all to see them)\n`,
      );
    }
    process.stdout.write('\n');
  }

  process.stdout.write('Per-scan summary:\n');
  for (const r of results) {
    process.stdout.write(
      `  [${r.name.padEnd(20)}] files=${r.files} matches=${r.matches} allowlist=${r.allowlistSize}\n`,
    );
  }
  process.stdout.write('\n');

  process.stdout.write('Summary:\n');
  process.stdout.write(`  Files:                ${totalFiles}\n`);
  process.stdout.write(`  Matches:              ${totalMatches}\n`);
  process.stdout.write(`  Allowlisted (total):  ${totalAllowlist}\n\n`);
  process.stdout.write(`AUDIT mui files=${totalFiles} matches=${totalMatches}\n`);

  if (failOnNew && totalFiles > 0) {
    process.stderr.write(
      `\nFAIL: ${totalFiles} non-allowlisted file(s) import @mui/*. ` +
        `Migrate to @sistent/sistent or update the allowlist with justification.\n`,
    );
    process.exit(1);
  }
}

main();
