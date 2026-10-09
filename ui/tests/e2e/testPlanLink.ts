// Test Plan deep-link helper (part A), shared by the connection test map and its
// unit test. Kept as a standalone, dependency-free module (no allure / Playwright
// imports) so the pure row-derivation logic can be unit-tested under vitest -
// which excludes `tests/e2e` from test discovery but can still import a helper
// from it.
//
// Each tracked connection case links back to its exact row in the Meshery Test
// Plan "Latest" tab so a reviewer can click from a report test straight to its
// source case. The row is derived from the Test # (col A, the `TC-<n>` id):
//   ROW = TestNum - CONN_ROW_OFFSET
//
// !!!  CURRENT-LAYOUT-DEPENDENT - REGENERATE IF THE SHEET IS RE-SORTED  !!!
// The offset and range below encode the CURRENT "Latest" tab layout, in which
// the connection cases Test# 1012..1089 occupy the contiguous rows 234..311
// (1012 - 778 = 234 ... 1089 - 778 = 311). This is the SAME contract as the CLI
// lane in `mesheryctl/bats-to-allure.js`; the two MUST be kept in lockstep. If
// the Latest tab is re-sorted, or rows are inserted above the connection block,
// recompute the offset AND the [CONN_TEST_MIN, CONN_TEST_MAX] range in BOTH
// files. The range guard is deliberate: an id outside the connection block (or a
// malformed one) yields NO link rather than a wrong one (e.g. `range=ANaN`).
export const CONN_TEST_MIN = 1012;
export const CONN_TEST_MAX = 1089;
export const CONN_ROW_OFFSET = 778;
const TEST_PLAN_SHEET_ID = '13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs';
const TEST_PLAN_GID = '838298230';

/**
 * Build the Test Plan "Latest" tab deep-link for a connection Test # (`TC-<n>`),
 * or return `null` when the id is malformed or falls outside the connection
 * block the offset is valid for. The link name ("Test Plan TC-<n>") matches the
 * CLI lane. See the offset caveat above.
 */
export function testPlanLink(testId: string): { url: string; name: string } | null {
  const match = /^TC-(\d+)$/i.exec((testId ?? '').trim());
  if (!match) return null;
  const testNum = Number(match[1]);
  if (testNum < CONN_TEST_MIN || testNum > CONN_TEST_MAX) return null;
  const row = testNum - CONN_ROW_OFFSET;
  const url =
    `https://docs.google.com/spreadsheets/d/${TEST_PLAN_SHEET_ID}` +
    `/edit?gid=${TEST_PLAN_GID}#gid=${TEST_PLAN_GID}&range=A${row}`;
  return { url, name: `Test Plan ${testId}` };
}
