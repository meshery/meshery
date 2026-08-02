import type { TestInfo } from '@playwright/test';

/**
 * Traceability map for the Kubernetes Connection Playwright suite.
 *
 * Each tracked case mirrors a row in the Meshery Test Plan spreadsheet
 * ("Latest" tab): `testId` is column A ("Test #") and `componentUnderTest` is
 * column C ("Component"), verbatim. Keeping the same `testId` here that the CLI
 * (BATS) and reporting lanes use lets the Connections Allure report line up
 * results across clients on a single Test #.
 *
 * IMPORTANT: `testId`s are reused from EXISTING sheet rows, never invented. As
 * of this writing the sheet's connection rows are Test # 96-104 (the sheet has
 * no higher connection-lifecycle numbering yet). New lifecycle scenarios that
 * do not yet have a sheet row are grouped by feature only (see
 * {@link connTagsUntracked}) until the Test Plan is expanded and their Test #s
 * assigned - at which point they graduate into {@link CONN_CASES}. This file is
 * the single reconciliation point between the sheet and the specs.
 */

/** Allure epic shared by every connection case (groups the Connections report). */
export const CONN_EPIC = 'Kubernetes Connections';

/** Client label for the shared cross-lane contract (UI vs CLI). */
export const CONN_CLIENT = 'UI';

export interface ConnCase {
  /** Meshery Test Plan "Latest" tab, column A ("Test #"), e.g. `TC-98`. */
  testId: string;
  /** Numeric sheet row, for diffing this map against the sheet. */
  sheetRow: number;
  /** Meshery Test Plan "Latest" tab, column C ("Component"), verbatim. */
  componentUnderTest: string;
  /** Allure feature grouping. */
  feature: string;
  /** Allure story. */
  story: string;
}

/**
 * Tracked connection cases, keyed by a stable symbol used in the specs.
 *
 * The `componentUnderTest` values are copied from the sheet as-authored; note
 * that Test # 104 predates the consolidation of kubeconfig upload into the
 * Connection Wizard, so its component still reads "UI/Settings" even though the
 * test now drives the wizard. The traceability id is what must stay stable;
 * updating the sheet's column C is a Test-Plan edit owned separately.
 */
export const CONN_CASES = {
  wizardOpen: {
    testId: 'TC-96',
    sheetRow: 96,
    componentUnderTest: 'UI/Connection Wizard',
    feature: 'Connection wizard',
    story: 'Open wizard and show kubeconfig upload',
  },
  wizardDiscoverContexts: {
    testId: 'TC-97',
    sheetRow: 97,
    componentUnderTest: 'UI/Connection Wizard',
    feature: 'Connection wizard',
    story: 'Discover kubeconfig contexts',
  },
  stateTransition: {
    testId: 'TC-98',
    sheetRow: 98,
    componentUnderTest: 'UI/Connection Wizard',
    feature: 'Connection state transitions',
    story: 'Transition a connection between lifecycle states',
  },
  operatorStatus: {
    testId: 'TC-101',
    sheetRow: 101,
    componentUnderTest: 'UI/Settings',
    feature: 'Meshery Operator',
    story: 'Operator and controller status for a connected cluster',
  },
  kubeconfigConnect: {
    testId: 'TC-104',
    sheetRow: 104,
    componentUnderTest: 'UI/Settings',
    feature: 'Kubeconfig import',
    story: 'Register and connect a cluster via kubeconfig upload',
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
 * allure-playwright maps epic/feature/story to Allure labels and captures the
 * rest; the repo's custom reporter also reads annotations. No extra dependency
 * is introduced.
 */
export function annotateConnCase(testInfo: TestInfo, key: ConnCaseKey): void {
  const testCase = CONN_CASES[key];
  testInfo.annotations.push(
    { type: 'testId', description: testCase.testId },
    { type: 'componentUnderTest', description: testCase.componentUnderTest },
    { type: 'client', description: CONN_CLIENT },
    { type: 'epic', description: CONN_EPIC },
    { type: 'feature', description: testCase.feature },
    { type: 'story', description: testCase.story },
  );
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
