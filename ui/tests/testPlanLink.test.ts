import { describe, expect, it } from 'vitest';

// Unit tests for the Test Plan deep-link helper (part A). This lives outside
// `tests/e2e` on purpose: vitest excludes `tests/e2e` from test discovery, and
// Playwright (testDir `tests/e2e`) would otherwise try to run a `*.test.ts`
// placed there. It imports the pure helper (no allure/Playwright coupling),
// mirroring the CLI-lane contract tests in `mesheryctl/bats-to-allure.test.js`.
import { testPlanLink } from './e2e/testPlanLink';

describe('testPlanLink', () => {
  it('maps a connection Test # to its Latest-tab row (ROW = n - 778)', () => {
    const link = testPlanLink('TC-1012');
    expect(link).not.toBeNull();
    expect(link?.name).toBe('Test Plan TC-1012');
    // Test# 1012 -> row 234.
    expect(link?.url).toMatch(/range=A234$/);
    expect(link?.url).toContain('13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs');
    expect(link?.url).toContain('gid=838298230');
  });

  it('covers both ends of the connection block', () => {
    expect(testPlanLink('TC-1012')?.url).toMatch(/range=A234$/); // first
    expect(testPlanLink('TC-1089')?.url).toMatch(/range=A311$/); // last (1089 - 778)
  });

  it('returns null outside the connection block', () => {
    expect(testPlanLink('TC-1011')).toBeNull(); // below the block
    expect(testPlanLink('TC-1090')).toBeNull(); // above the block
    expect(testPlanLink('TC-500')).toBeNull(); // unrelated case
  });

  it('returns null for malformed or empty ids (no ANaN link)', () => {
    expect(testPlanLink('nonsense')).toBeNull();
    expect(testPlanLink('1012')).toBeNull(); // missing TC- prefix
    expect(testPlanLink('')).toBeNull();
    // @ts-expect-error - exercise the runtime guard against nullish input
    expect(testPlanLink(undefined)).toBeNull();
  });
});
