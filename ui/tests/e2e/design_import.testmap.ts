import type { TestInfo } from '@playwright/test';
import * as allure from 'allure-js-commons';

/**
 * Traceability map for the Design Import Playwright suite.
 *
 * Tracks the UI-lane cases that cover the import flow introduced/fixed in
 * meshery/meshery#21105:
 *   - `fileName` (camelCase) wire field on file-upload imports
 *   - Duplicate-submission prevention in `ImportDesignModal`
 *   - URL-import end-to-end flow
 *
 * The Test Group is `Design Import` — it will appear in the meshery/qa Allure
 * dashboard once the `testGroup` label is emitted (see `DESIGN_IMPORT_TEST_GROUP`
 * below and the `beforeEach` in `design_import.spec.ts`).
 */

/** Allure epic shared by every case in this spec. */
export const DESIGN_IMPORT_EPIC = 'Design Import';

/** Client label for the shared cross-lane contract. */
export const DESIGN_IMPORT_CLIENT = 'UI';

/** Allure testGroup label value — the qa dashboard keys on this. */
export const DESIGN_IMPORT_TEST_GROUP = 'Design Import';

export interface DesignImportCase {
  /** Meshery Test Plan column A — placeholder until rows are assigned. */
  testId: string;
  /** Allure feature grouping. */
  feature: string;
  /** Allure story. */
  story: string;
  /** Meshery Test Plan column D verbatim. */
  componentUnderTest: string;
}

/**
 * Tracked cases. testIds are placeholders (TC-DI-*) until the qa sheet assigns
 * real row numbers; graduate each to a real TC-NNNN once assigned.
 */
export const DESIGN_IMPORT_CASES = {
  /** Open the Import Design modal from the Designs toolbar. */
  openImportModal: {
    testId: 'TC-DI-001',
    componentUnderTest: 'UI/Design Import Modal',
    feature: 'Design import modal',
    story: 'Open the Import Design modal from the toolbar',
  },
  /**
   * Import via URL — verifies the wire body sends the `url` field (not
   * `file_name`), the success notification appears, and the modal closes.
   */
  importViaUrl: {
    testId: 'TC-DI-002',
    componentUnderTest: 'UI/Design Import Modal',
    feature: 'Design import wire contract',
    story: 'POST /api/pattern/import sends url for URL imports',
  },
  /**
   * Import via file upload — asserts the wire body uses `fileName`
   * (camelCase) NOT `file_name` (snake_case), fixing the regression from #21105.
   */
  importViaFile: {
    testId: 'TC-DI-003',
    componentUnderTest: 'UI/Design Import Modal',
    feature: 'Design import wire contract',
    story: 'POST /api/pattern/import sends camelCase fileName for file-upload imports',
  },
  /**
   * Duplicate-submission prevention — the Import button must be disabled while
   * an import is in-flight so rapid clicks cannot trigger a second request.
   */
  preventsDuplicateSubmit: {
    testId: 'TC-DI-004',
    componentUnderTest: 'UI/Design Import Modal',
    feature: 'Design import duplicate-submission prevention',
    story: 'Import button is disabled and shows Importing… while in flight',
  },
} as const satisfies Record<string, DesignImportCase>;

export type DesignImportCaseKey = keyof typeof DESIGN_IMPORT_CASES;

/** Playwright tags for a tracked case — grep-safe and surfaced as Allure tag labels. */
export function designImportTags(key: DesignImportCaseKey): string[] {
  const c = DESIGN_IMPORT_CASES[key];
  return ['@design-import', `@${c.testId}`, '@cut:ui-design-import-modal', '@client:ui'];
}

/** Untracked tag (no sheet row yet). */
export function designImportTagsUntracked(component: string): string[] {
  return [
    '@design-import',
    `@cut:${component.toLowerCase().replace(/[^a-z0-9]+/g, '-')}`,
    '@client:ui',
  ];
}

/** Emit the shared Allure annotation contract for a tracked case. */
export async function annotateDesignImportCase(
  testInfo: TestInfo,
  key: DesignImportCaseKey,
): Promise<void> {
  const c = DESIGN_IMPORT_CASES[key];
  testInfo.annotations.push(
    { type: 'testId', description: c.testId },
    { type: 'componentUnderTest', description: c.componentUnderTest },
    { type: 'client', description: DESIGN_IMPORT_CLIENT },
    { type: 'epic', description: DESIGN_IMPORT_EPIC },
    { type: 'feature', description: c.feature },
    { type: 'story', description: c.story },
  );
  // No testPlanLink yet — placeholder IDs are not in the connection block.
  void allure; // imported for parity; use allure.link() once TC-DI-* get real ids.
}

/** Emit annotations for an untracked case. */
export function annotateDesignImportCaseUntracked(
  testInfo: TestInfo,
  { feature, story }: { feature: string; story?: string },
): void {
  testInfo.annotations.push(
    { type: 'epic', description: DESIGN_IMPORT_EPIC },
    { type: 'client', description: DESIGN_IMPORT_CLIENT },
    { type: 'feature', description: feature },
  );
  if (story) {
    testInfo.annotations.push({ type: 'story', description: story });
  }
}
