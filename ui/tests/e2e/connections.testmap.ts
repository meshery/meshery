import type { TestInfo } from '@playwright/test';
import * as allure from 'allure-js-commons';
import { testPlanLink } from './testPlanLink';

/**
 * Traceability map for the Kubernetes Connection Playwright suite.
 *
 * Each tracked case mirrors a row in the Meshery Test Plan spreadsheet
 * ("Latest" tab). The tab's columns are: A = Test #, B = Test Group,
 * C = Client, D = Component Under Test. So `testId` is column A ("Test #") and
 * `componentUnderTest` is column D, verbatim. Every case's Test Group (column B)
 * is "Connection Lifecycle", emitted once as the `testGroup` Allure label in the
 * spec's describe-level beforeEach (see connections.spec.ts). Keeping the same
 * `testId` here that the CLI (BATS) and reporting lanes use lets the Connection
 * Lifecycle Allure report line up results across clients on a single Test #.
 *
 * IMPORTANT: `testId`s are reused from EXISTING sheet rows, never invented. The
 * Kubernetes Connection lifecycle cases live at Test # 1012-1089 (the "Latest"
 * tab is contiguous 1..1011 then 1012..1089), and each row's `[matrix Rn]`
 * remark ties it to the connection scenario matrix. Only the Client=UI / Both
 * rows are automated here; scenarios without a rewritten test yet are grouped by
 * feature only (see {@link connTagsUntracked}) until a test is added - at which
 * point they graduate into {@link CONN_CASES}. This file is the single
 * reconciliation point between the sheet and the specs.
 */

/** Allure epic shared by every connection case (groups the Connections report). */
export const CONN_EPIC = 'Kubernetes Connections';

/** Client label for the shared cross-lane contract (UI vs CLI). */
export const CONN_CLIENT = 'UI';

// Test Plan deep-link (part A): the row-derivation contract lives in the pure,
// unit-tested `./testPlanLink` helper (guards the `TC-<n>` shape and the
// 1012..1089 connection block, returning null otherwise so a malformed/
// out-of-block id never becomes a wrong link). Every CONN_CASES testId is inside
// that block by construction (see CONN_CASES below), so the guard normally
// passes - it exists so a future out-of-block id fails safe instead of emitting
// `range=ANaN`. Re-exported here for back-compat with existing importers.
export { testPlanLink } from './testPlanLink';

export interface ConnCase {
  /** Meshery Test Plan "Latest" tab, column A ("Test #"), e.g. `TC-98`. */
  testId: string;
  /** Numeric sheet row, for diffing this map against the sheet. */
  sheetRow: number;
  /** Meshery Test Plan "Latest" tab, column D ("Component Under Test"), verbatim. */
  componentUnderTest: string;
  /** Allure feature grouping. */
  feature: string;
  /** Allure story. */
  story: string;
}

/**
 * Tracked connection cases, keyed by a stable symbol used in the specs. Each maps
 * to a Client=UI row in the Test Plan's connection block (Test # 1012-1089); the
 * `[matrix Rn]` id from the sheet is noted for cross-reference. `componentUnderTest`
 * is copied from the sheet's column D ("Component Under Test") verbatim.
 */
export const CONN_CASES = {
  // matrix R1
  kubeconfigConnect: {
    testId: 'TC-1012',
    sheetRow: 1012,
    componentUnderTest: 'UI/Connection Wizard',
    feature: 'Kubeconfig import',
    story: 'Register reachable single-context kubeconfig',
  },
  // matrix R3 (this test covers the multi-context discover + review-listing portion)
  wizardDiscoverContexts: {
    testId: 'TC-1014',
    sheetRow: 1014,
    componentUnderTest: 'UI/Connection Wizard',
    feature: 'Connection wizard',
    story: 'Register multi-context kubeconfig (discover + review contexts)',
  },
  // matrix X1
  stateTransition: {
    testId: 'TC-1061',
    sheetRow: 1061,
    componentUnderTest: 'UI/Connections Table',
    feature: 'Connection state transitions',
    story: 'Disconnect a connected cluster',
  },
  // matrix X3
  delete: {
    testId: 'TC-1063',
    sheetRow: 1063,
    componentUnderTest: 'UI/Connections Table',
    feature: 'Connection lifecycle',
    story: 'FSM delete (graceful)',
  },
} as const satisfies Record<string, ConnCase>;

export type ConnCaseKey = keyof typeof CONN_CASES;

/** Turn a Component value into a stable, grep-safe tag slug. */
const slug = (value: string): string =>
  value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');

/**
 * Playwright tags for a tracked case. These are grep-able (`--grep @TC-98`) and
 * are surfaced by allure-playwright as Allure `tag` labels and by the repo's
 * custom reporter, so the Test # and Component travel with the result without a
 * runtime Allure dependency.
 */
export function connTags(key: ConnCaseKey): string[] {
  const testCase = CONN_CASES[key];
  return [
    '@connections',
    `@${testCase.testId}`,
    `@cut:${slug(testCase.componentUnderTest)}`,
    `@client:${slug(CONN_CLIENT)}`,
  ];
}

/**
 * Coarse tags for a connection behaviour that has no dedicated sheet Test # yet
 * (e.g. hard delete, table filtering). Grouped by Component only so it still
 * lands in the Connections report; graduate it into {@link CONN_CASES} once a
 * Test # is assigned.
 */
export function connTagsUntracked(componentUnderTest: string): string[] {
  return ['@connections', `@cut:${slug(componentUnderTest)}`, `@client:${slug(CONN_CLIENT)}`];
}

/**
 * Emit the shared Allure label contract via Playwright annotations:
 * `testId`, `componentUnderTest`, `client`, `epic`, `feature`, `story`.
 * allure-playwright maps epic/feature/story to Allure labels and records the
 * rest as parameters, so no extra dependency is introduced. (The repo's custom
 * Playwright reporter only processes `relationship` annotations - it does not
 * read these; the tags from {@link connTags} are what it surfaces.)
 */
export async function annotateConnCase(testInfo: TestInfo, key: ConnCaseKey): Promise<void> {
  const testCase = CONN_CASES[key];
  testInfo.annotations.push(
    { type: 'testId', description: testCase.testId },
    { type: 'componentUnderTest', description: testCase.componentUnderTest },
    { type: 'client', description: CONN_CLIENT },
    { type: 'epic', description: CONN_EPIC },
    { type: 'feature', description: testCase.feature },
    { type: 'story', description: testCase.story },
  );

  // Emit the Test Plan row as a real Allure `tms` link (part A) - the UI-lane
  // counterpart to the CLI converter's `links` entry. This uses the
  // allure-js-commons runtime API for the same reason as the testGroup label:
  // allure-playwright maps only known annotation types (epic/feature/story) to
  // the report, so a custom link must be set through allure.link to appear. The
  // helper returns null for an out-of-block/malformed id (fail-safe, no wrong
  // link); every tracked case is in-block so this normally yields a link.
  const link = testPlanLink(testCase.testId);
  if (link) {
    await allure.link(link.url, link.name, 'tms');
  }
}

/**
 * Emit the shared Allure label contract for an untracked connection case - one
 * with no dedicated sheet Test # yet (see {@link connTagsUntracked}): `epic`,
 * `client`, `feature`, and `story` when provided. Mirrors
 * {@link annotateConnCase} so the sheet<->code label contract stays defined in
 * this file rather than inlined per spec.
 */
export function annotateConnCaseUntracked(
  testInfo: TestInfo,
  { feature, story }: { feature: string; story?: string },
): void {
  testInfo.annotations.push(
    { type: 'epic', description: CONN_EPIC },
    { type: 'client', description: CONN_CLIENT },
    { type: 'feature', description: feature },
  );
  if (story) {
    testInfo.annotations.push({ type: 'story', description: story });
  }
}
