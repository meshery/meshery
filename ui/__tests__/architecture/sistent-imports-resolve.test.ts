import { describe, expect, it } from 'vitest';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import * as Sistent from '@sistent/sistent';

/**
 * Architecture guard: every value imported by name from `@sistent/sistent`
 * must exist as a real runtime export of the installed package.
 *
 * This catches a recurring, whole-page-killing bug class:
 *
 *   import { DeleteForever } from '@sistent/sistent';   // ❌ runtime undefined
 *
 * Sistent ships the `*Icon`-suffixed names (`DeleteForeverIcon`, `SearchIcon`,
 * `CompareArrowsIcon`, `GetAppIcon`, `IndeterminateCheckBoxIcon`, …) and the
 * MUI-v5 `Accordion*` family rather than the v4 `ExpansionPanel*` names. An
 * import that names something the package does not export resolves to
 * `undefined` at runtime; rendering it throws React error #130 ("Element type
 * is invalid … got: undefined") and the page-level ErrorBoundary swallows the
 * whole view.
 *
 * Neither `tsc` (Sistent's published types and runtime can diverge for a given
 * glyph) nor the component unit tests (which `vi.mock('@sistent/sistent')` and
 * therefore supply whatever name the test asks for) catch this. So this test
 * imports the REAL package and validates every named *value* import in the
 * source tree against `Object.keys(Sistent)`.
 *
 * Type-only imports are skipped: `import type { X }` and inline `import
 * { type X }` are erased under `verbatimModuleSyntax` and never reach the
 * runtime namespace, so they cannot cause this failure.
 */

const RUNTIME_EXPORTS = new Set(Object.keys(Sistent));

// Type-only Sistent exports imported somewhere WITHOUT the `type` keyword.
// Prefer to keep this empty: the correct fix for a new entry is a `type`
// import at the call site (required by `verbatimModuleSyntax` anyway), not an
// allowlist addition. Only add a name here if it is genuinely type-only and a
// `type` import is somehow not viable.
const ALLOWED_TYPE_ONLY: ReadonlySet<string> = new Set<string>();

const UI_ROOT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');

const PRUNE_DIRS = new Set([
  'node_modules',
  '.next',
  'out',
  'coverage',
  '.git',
  'public',
  'docs',
  'tests', // playwright e2e
  'cypress',
]);

const SOURCE_EXT = new Set(['.ts', '.tsx', '.js', '.jsx', '.mjs', '.cjs']);

function collectSourceFiles(dir: string, acc: string[]): string[] {
  let entries: fs.Dirent[];
  try {
    entries = fs.readdirSync(dir, { withFileTypes: true });
  } catch {
    return acc;
  }
  for (const entry of entries) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (!PRUNE_DIRS.has(entry.name)) collectSourceFiles(full, acc);
    } else if (entry.isFile()) {
      if (/\.(test|spec|stories)\.[cm]?[tj]sx?$/.test(entry.name)) continue;
      if (SOURCE_EXT.has(path.extname(entry.name))) acc.push(full);
    }
  }
  return acc;
}

// Matches `import [type] [Default,] { ...names... } from '@sistent/sistent'`.
// `[^}]*` spans newlines, so multi-line import blocks are covered.
const IMPORT_RE =
  /import\s+(type\s+)?(?:[A-Za-z_$][\w$]*\s*,\s*)?\{([^}]*)\}\s*from\s*['"]@sistent\/sistent['"]/g;

function unresolvedValueImports(source: string): string[] {
  // Strip block and line comments from the whole file before matching, so a
  // commented-out import cannot raise a false positive and a specifier after
  // an inline comment is still validated. String literals are preserved so a
  // `//` inside a URL (e.g. 'https://…') is never mistaken for a comment.
  const cleanSource = source.replace(
    /("([^"\\]|\\.)*"|'([^'\\]|\\.)*'|`([^`\\]|\\.)*`)|\/\*[\s\S]*?\*\/|\/\/.*/g,
    (full, stringLiteral) => (stringLiteral ? stringLiteral : ''),
  );

  const bad: string[] = [];
  for (const match of cleanSource.matchAll(IMPORT_RE)) {
    if (match[1]) continue; // whole statement is `import type { ... }`
    for (const rawSpecifier of match[2].split(',')) {
      const specifier = rawSpecifier.trim();
      if (!specifier) continue;
      if (/^type\s+/.test(specifier)) continue; // inline type-only specifier
      const sourceName = specifier.split(/\s+as\s+/)[0].trim();
      if (!/^[A-Za-z_$][\w$]*$/.test(sourceName)) continue;
      if (RUNTIME_EXPORTS.has(sourceName)) continue;
      if (ALLOWED_TYPE_ONLY.has(sourceName)) continue;
      bad.push(sourceName);
    }
  }
  return bad;
}

describe('@sistent/sistent named imports resolve at runtime', () => {
  it('reads a non-empty set of real runtime exports (guards a vacuous pass)', () => {
    expect(RUNTIME_EXPORTS.size).toBeGreaterThan(0);
  });

  it('imports no value that is missing from the installed package', () => {
    const files = collectSourceFiles(UI_ROOT, []);
    expect(files.length).toBeGreaterThan(0);

    const violations: string[] = [];
    for (const file of files) {
      for (const name of unresolvedValueImports(fs.readFileSync(file, 'utf8'))) {
        violations.push(`${name}  <-  ${path.relative(UI_ROOT, file)}`);
      }
    }

    // A non-empty list means a component imports a name that does not exist in
    // the installed @sistent/sistent — it will be `undefined` at runtime and
    // crash the page when rendered. Use the `*Icon` name (or the MUI-v5
    // Accordion* name), or a `type` import if the symbol is type-only.
    expect(violations).toEqual([]);
  });
});
