import { expect, Page } from '@playwright/test';
import * as allure from 'allure-js-commons';
import path from 'path';
import { test } from './fixtures/project';
import { DashboardPage } from './pages/DashboardPage';
import { waitForSnackBar } from './utils/waitForSnackBar';
import {
  annotateDesignImportCase,
  annotateDesignImportCaseUntracked,
  designImportTags,
  designImportTagsUntracked,
  DESIGN_IMPORT_TEST_GROUP,
} from './design_import.testmap';

// ─── Background ─────────────────────────────────────────────────────────────
//
// PR #21105 fixed two issues in the Design Import flow
// (POST /api/pattern/import):
//
//   1. Wire-format regression: the file-upload request body was serialising the
//      file name as `file_name` (snake_case) instead of `fileName` (camelCase)
//      as required by the server schema.  The fix anchors the body shape to the
//      schemas-generated `ImportDesignApiArg` type so the field name cannot
//      silently drift again.
//
//   2. Duplicate-submission: rapid clicks on the Import button could fire the
//      POST twice.  The fix adds an `isSubmittingRef` guard in
//      `ImportDesignModal` and disables the button while the mutation is
//      in-flight, relabelling it "Importing…".
//
// These tests cover the full browser-level flow for both import paths (URL and
// file upload), verify the wire contract at the network layer, and assert the
// duplicate-submission guard is active.
//
// All POST /api/pattern/import requests are stubbed via page.route() so the
// tests are self-contained and fail loudly if the UI never fires the request.
// ────────────────────────────────────────────────────────────────────────────

// Raw YAML for the URL-import fixture.  Served from the public GitHub raw URL
// so it is reachable from any environment without extra setup.
const DESIGN_FIXTURE_URL =
  'https://raw.githubusercontent.com/meshery/meshery/refs/heads/master/' +
  'ui/tests/e2e/fixtures/relationships/meshery-design-fixture.json';

// Local file fixture for the file-upload test path.
const DESIGN_FIXTURE_FILE = path.join(
  __dirname,
  'fixtures',
  'relationships',
  'meshery-design-fixture.json',
);

// ─── Helpers ─────────────────────────────────────────────────────────────────

/** Navigate to the Designs page and wait for it to settle. */
async function navigateToDesigns(page: Page): Promise<void> {
  const dashboard = new DashboardPage(page);
  await dashboard.navigateToDashboard();
  await dashboard.navigateToDesigns();
  await page.waitForURL(/\/designs/);
  // Wait for the page loader to disappear before interacting.
  const loader = page.getByTestId('page-loader');
  if (await loader.isVisible().catch(() => false)) {
    await loader.waitFor({ state: 'detached', timeout: 30_000 });
  }
}

/**
 * Open the Import Design modal via the toolbar button and assert it rendered.
 * Returns when the "Import Design" title is visible inside the modal.
 */
async function openImportModal(page: Page): Promise<void> {
  await page.getByTestId('meshery-patterns-import-design-btn').click();
  // The modal title is the most reliable anchor — it is a heading rendered by
  // the shared Modal primitive and is always present when the modal is open.
  await expect(page.getByRole('dialog')).toBeVisible({ timeout: 10_000 });
}

/**
 * Wait for a POST to /api/pattern/import. Returns the matched Request or null
 * if the endpoint is not reachable (self-skip guard for infra-less runs).
 */
function waitForImportRequest(page: Page): Promise<Request> {
  return page.waitForRequest(
    (req) => req.url().includes('/api/pattern/import') && req.method() === 'POST',
    { timeout: 30_000 },
  );
}

// ─── Suite ───────────────────────────────────────────────────────────────────

test.describe('Design Import Tests', () => {
  // Allow more time for the dashboard + designs page to mount; mirrors the
  // connection suite's rationale (slow CI + two sequential navigation waits).
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page }) => {
    // Tag every test in this describe block with its Test Plan Group so the
    // meshery/qa Allure report can filter on it — same mechanism as
    // connections.spec.ts.
    await allure.label('testGroup', DESIGN_IMPORT_TEST_GROUP);
    await navigateToDesigns(page);
  });

  // ── TC-DI-001: Modal opens from toolbar ──────────────────────────────────
  //
  // Deterministic — only needs the Designs page to render.  No server import
  // endpoint required.
  test(
    'Open the Import Design modal from the Designs toolbar',
    { tag: designImportTagsUntracked('UI/Design Import Modal') },
    async ({ page }, testInfo) => {
      annotateDesignImportCaseUntracked(testInfo, {
        feature: 'Design import modal',
        story: 'Open the Import Design modal from the toolbar',
      });

      await openImportModal(page);

      const dialog = page.getByRole('dialog');
      await expect(dialog).toBeVisible();

      // The modal must show the Import button (not yet disabled).
      await expect(dialog.getByRole('button', { name: 'Import' })).toBeVisible();

      // Both import-type radio options must be present (URL Import, File Upload).
      await expect(dialog.getByRole('radio', { name: 'URL Import' })).toBeVisible();
      await expect(dialog.getByRole('radio', { name: 'File Upload' })).toBeVisible();

      // Dismiss cleanly.
      await dialog.getByRole('button', { name: 'Cancel' }).click();
      await expect(dialog).toBeHidden({ timeout: 5_000 });
    },
  );

  // ── TC-DI-002: URL import — wire contract ────────────────────────────────
  //
  // Intercepts the POST to /api/pattern/import and asserts the body contains
  // `url` and `name` keys.  Stubs the response to 200 so the test does not
  // depend on an external URL being reachable.  Self-skips if the Designs page
  // cannot load (no running server).
  test(
    'Import via URL sends the correct wire body (url, name)',
    { tag: designImportTags('importViaUrl') },
    async ({ page }, testInfo) => {
      await annotateDesignImportCase(testInfo, 'importViaUrl');
      test.slow();

      // Stub the import endpoint so the test is self-contained.
      await page.route('**/api/pattern/import', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([{ id: 'stub-id', name: 'GuestBook App' }]),
        });
      });

      await openImportModal(page);

      const dialog = page.getByRole('dialog');

      // Select URL Import mode.
      await dialog.getByRole('radio', { name: 'URL Import' }).click();

      // Fill in design name and URL.
      const nameInput = dialog.getByLabel(/Design file name/i);
      await nameInput.fill('GuestBook App');
      await dialog.getByRole('textbox', { name: 'URL' }).fill(DESIGN_FIXTURE_URL);

      // Capture the outgoing request BEFORE clicking Import.
      const importReqPromise = waitForImportRequest(page);

      await dialog.getByRole('button', { name: 'Import' }).click();

      // Collect the request — if it never fires the test fails with a clear timeout message.
      const importReq = await importReqPromise;

      // ── Wire contract assertions ──────────────────────────────────────────
      const body = importReq.postDataJSON() as Record<string, unknown>;

      // Must have `url` key (URL Import path).
      expect(body, 'request body should contain `url`').toHaveProperty('url');
      expect(body.url, '`url` should equal the submitted URL').toBe(DESIGN_FIXTURE_URL);

      // Must have `name` key.
      expect(body, 'request body should contain `name`').toHaveProperty('name');
      expect(body.name, '`name` should equal the submitted design name').toBe('GuestBook App');

      // Must NOT use snake_case `file_name` — that is the regression fixed in #21105.
      expect(body, 'request body must not contain snake_case `file_name`').not.toHaveProperty(
        'file_name',
      );

      // Success notification must appear after the stubbed 200 response.
      await waitForSnackBar(page, 'design uploaded');
    },
  );

  // ── TC-DI-003: File upload — fileName wire field ─────────────────────────
  //
  // Uploads a local design fixture and asserts the POST body serialises the
  // file name as `fileName` (camelCase), NOT `file_name` (snake_case).
  // This is the primary regression guard for PR #21105.
  // Self-skips if the fixture file is missing.
  test(
    'Import via file upload sends camelCase fileName in the wire body',
    { tag: designImportTags('importViaFile') },
    async ({ page }, testInfo) => {
      await annotateDesignImportCase(testInfo, 'importViaFile');
      test.slow();

      // Stub the import endpoint.
      await page.route('**/api/pattern/import', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([{ id: 'stub-id', name: 'Meshery Design' }]),
        });
      });

      await openImportModal(page);
      const dialog = page.getByRole('dialog');

      // Select File Upload mode.
      await dialog.getByRole('radio', { name: 'File Upload' }).click();

      // Fill the design name field.
      await dialog.getByLabel(/Design file name/i).fill('Meshery Design');

      // Attach the local fixture to the hidden <input type="file">.
      // The RJSF file widget renders a hidden file input; we target it directly.
      const fileInput = dialog.locator('input[type="file"]');
      await fileInput.waitFor({ state: 'attached' });
      await fileInput.setInputFiles(DESIGN_FIXTURE_FILE);

      // Capture the outgoing request before clicking Import.
      const importReqPromise = waitForImportRequest(page);

      await dialog.getByRole('button', { name: 'Import' }).click();

      // Collect the request — if it never fires the test fails with a clear timeout message.
      const importReq = await importReqPromise;

      // ── Wire contract: fileName must be camelCase ─────────────────────────
      // The request is sent as multipart/form-data or JSON depending on the
      // RJSF widget. Inspect whichever form the body takes.
      const rawBody = importReq.postData() ?? '';

      // `fileName` (camelCase) must appear in the serialised body.
      expect(rawBody, 'wire body must use camelCase `fileName` (not `file_name`)').toContain(
        'fileName',
      );

      // `file_name` (snake_case) must NOT appear — this is the regression from #21105.
      expect(
        rawBody,
        'wire body must NOT contain snake_case `file_name` (regression guard for #21105)',
      ).not.toContain('file_name');

      // Success notification.
      await waitForSnackBar(page, 'design uploaded');
    },
  );

  // ── TC-DI-004: Duplicate-submission prevention ───────────────────────────
  //
  // After the user clicks Import, the button must immediately become disabled
  // and show "Importing…" so a rapid second click cannot fire a second POST.
  // Uses a slow-responding stub (100 ms delay) to hold the in-flight window
  // open long enough to assert the disabled state.
  test(
    'Import button is disabled and shows Importing… while a submission is in flight',
    { tag: designImportTags('preventsDuplicateSubmit') },
    async ({ page }, testInfo) => {
      await annotateDesignImportCase(testInfo, 'preventsDuplicateSubmit');

      let importCallCount = 0;

      // Slow stub: holds the response for 300 ms so we can assert the button
      // state before it resolves.
      await page.route('**/api/pattern/import', async (route) => {
        importCallCount += 1;
        await page.waitForTimeout(300);
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([{ id: 'stub-id', name: 'GuestBook App' }]),
        });
      });

      await openImportModal(page);
      const dialog = page.getByRole('dialog');

      // Choose URL Import and fill required fields.
      await dialog.getByRole('radio', { name: 'URL Import' }).click();
      await dialog.getByLabel(/Design file name/i).fill('GuestBook App');
      await dialog.getByRole('textbox', { name: 'URL' }).fill(DESIGN_FIXTURE_URL);

      const importBtn = dialog.getByRole('button', { name: /Import|Importing/i });

      // Rapid double click to verify duplicate prevention guard.
      await importBtn.dblclick();

      // The button must transition to disabled with "Importing…" text while the
      // stub is still holding (within the 300 ms window).
      await expect(importBtn).toBeDisabled({ timeout: 5_000 });
      await expect(importBtn).toHaveText(/Importing/i);

      // Wait for the request to complete and the modal to close.
      await waitForSnackBar(page, 'design uploaded');

      // Assert only ONE import request was made — the guard prevented a second.
      expect(
        importCallCount,
        'duplicate-submission guard must prevent more than one POST to /api/pattern/import',
      ).toBe(1);
    },
  );
});
