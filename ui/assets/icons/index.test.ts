import { describe, expect, it } from 'vitest';

import * as iconBarrel from './index';

/**
 * Integrity guard for the canonical icon barrel.
 *
 * The barrel re-exports design-system glyphs from `@sistent/sistent` under
 * legacy names (e.g. `DeleteForeverIcon as DeleteForever`). When a re-export
 * names a member that the installed Sistent build does not actually export,
 * the alias silently resolves to `undefined` at runtime. Rendering such an
 * alias throws React error #130 ("Element type is invalid ... got: undefined")
 * and takes down the whole page — this is exactly what happened to the
 * Connections page, where `DeleteForever` had been pointed at a non-existent
 * `@sistent/sistent` export instead of `DeleteForeverIcon`.
 *
 * Neither `tsc` (Sistent's types and runtime diverge for this case) nor the
 * mocked component unit tests can catch the mismatch, so this test imports the
 * REAL barrel (no `vi.mock`) and asserts that every value export resolves to a
 * defined value.
 */
describe('icon barrel integrity', () => {
  it('exposes no undefined value exports (every re-export resolves)', () => {
    const undefinedExports = Object.entries(iconBarrel)
      .filter(([, value]) => value === undefined)
      .map(([name]) => name);

    expect(undefinedExports).toEqual([]);
  });

  it('re-exports DeleteForever from the real Sistent glyph', () => {
    // Direct regression assertion for the Connections-page crash.
    expect(iconBarrel.DeleteForever).toBeDefined();
  });
});
