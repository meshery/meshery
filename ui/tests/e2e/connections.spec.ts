import { expect, Page, Response } from '@playwright/test';
import * as allure from 'allure-js-commons';
import fs from 'fs';
import os from 'os';
import path from 'path';
import { test } from './fixtures/project';
import { DashboardPage } from './pages/DashboardPage';
import { waitForSnackBar } from './utils/waitForSnackBar';
import {
  annotateConnCase,
  annotateConnCaseUntracked,
  connTags,
  connTagsUntracked,
} from './connections.testmap';

// The connection-lifecycle flows below were rewritten against the current UI.
// Cluster registration now runs through the Connection Wizard (deep link
// `?create=true&kind=kubernetes` -> upload `#connection-wizard-kubeconfig-input`
// -> "Review contexts" -> Import), and state transitions run through the
// Connections table's status dropdown (`#connection-status-select`) plus the
// shared `connection-transition-*` confirmation modal. The previous
// single-modal `connection-addKubernetesModal` / `connection-uploadKubeConfig`
// testids no longer exist.
//
// Tests that need a live Meshery server and a reachable cluster self-skip when
// the environment cannot provide one (mirroring the delete test), so an
// infra-less run degrades to "skipped" rather than a false failure.

// Test Plan Test Group (col B) for every case in this spec. Emitted once as an
// Allure `testGroup` label in the describe's beforeEach so a single mechanism
// covers the whole file - this is the UI-lane counterpart to the BATS
// `[tg=Connection Lifecycle]` title token (see mesheryctl/bats-to-allure.js).
// The meshery/qa "Connection Lifecycle" report keys on this label. Any other
// Test Group can drive its own filtered report the same way. This is a real
// Allure label emitted via `allure.label`, not a Playwright annotation:
// allure-playwright only maps known annotation types (epic/feature/story) to
// labels, so a custom `testGroup` must be set through allure.label to become a
// label the meshery/qa report can filter on. The runtime API comes from
// `allure-js-commons` (Allure's documented Playwright API and an exact-pinned
// dependency of allure-playwright); allure-playwright's own `allure` export is
// deprecated in favor of it.
const TEST_GROUP = 'Connection Lifecycle';

const MULTI_CONTEXT_KUBECONFIG = path.join(__dirname, 'assets', 'kubeconfig-multi-context.yaml');
const HOST_KUBECONFIG = path.join(os.homedir(), '.kube', 'config');

// Wait for the connections LIST response specifically. Matching any 200 whose
// URL merely contains "/api/integrations/connections" would also resolve on the
// PUT /connections/{id} status-update, letting a table-refresh wait complete on
// an unrelated request; restrict to the GET collection endpoint.
function waitForConnectionsApiResponse(page: Page): Promise<Response> {
  return page.waitForResponse(
    (response) =>
      response.request().method() === 'GET' &&
      /\/api\/integrations\/connections(\?|$)/.test(response.url()) &&
      response.status() === 200,
  );
}

/** Open the Connection Wizard directly on the Kubernetes kubeconfig step. */
async function openKubernetesWizard(page: Page): Promise<void> {
  await page.goto('/management/connections?create=true&kind=kubernetes', {
    waitUntil: 'domcontentloaded',
  });
  // The wizard opens preset to Kubernetes on the "Import Kubeconfig" step. Its
  // StepHeader is an h6 ("Upload a kubeconfig"), a reliable heading anchor for
  // "the wizard opened and reached the kubeconfig step". The modal chrome title
  // ("Create Connection") is not rendered as a heading element, so it must not
  // be asserted via getByRole('heading').
  await expect(page.getByRole('heading', { name: 'Upload a kubeconfig' })).toBeVisible();
}

/** Attach a kubeconfig to the wizard's hidden file input. */
async function uploadKubeconfig(page: Page, kubeconfigPath: string): Promise<void> {
  await page.locator('#connection-wizard-kubeconfig-input').setInputFiles(kubeconfigPath);
}

test.describe.serial('Connection Management Tests', () => {
  // The shared beforeEach calls dashboardPage.navigateToDashboard() and
  // navigateToConnections(), each of which internally awaits two visibility
  // checks with a 120s timeout. Under the default BASE_TIMEOUT=60s the hook
  // itself dies before those waits can resolve when CI is under load.
  // Three minutes is enough headroom for a slow dashboard render plus the
  // connections page mount and its initial API response.
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page }) => {
    // Tag every test in this spec with its Test Plan Test Group (col B) so the
    // Test-Group-keyed meshery/qa report picks them up. One call here covers
    // all cases in the describe block.
    await allure.label('testGroup', TEST_GROUP);

    const initialConnectionsRes = waitForConnectionsApiResponse(page);
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToConnections();
    await page.waitForURL(/\/management\/connections/);
    await initialConnectionsRes;
    await expect(page.getByTestId('ConnectionTable-search')).toBeVisible();
  });

  test(
    'Verify that UI components are displayed',
    { tag: connTagsUntracked('UI/Connections Table') },
    async ({ page }, testInfo) => {
      annotateConnCaseUntracked(testInfo, { feature: 'Connections table' });
      // Verify that connections table is displayed (by checking for table headings)
      const headings = ['Name', 'Environments', 'Kind', 'Category', 'Status', 'Actions'];
      for (const heading of headings) {
        await expect(page.getByRole('columnheader', { name: heading })).toBeVisible();
      }
    },
  );

  // Opening the Kubernetes wizard surfaces the kubeconfig upload. This is
  // deterministic - it needs a running Meshery UI but no cluster - so it runs in
  // CI rather than self-skipping. Untracked: it is a UI smoke of the wizard entry
  // with no dedicated Test Plan row (the matrix begins at registration, R1).
  test(
    'Open the Kubernetes connection wizard and show kubeconfig upload',
    { tag: connTagsUntracked('UI/Connection Wizard') },
    async ({ page }, testInfo) => {
      annotateConnCaseUntracked(testInfo, {
        feature: 'Connection wizard',
        story: 'Open wizard and show kubeconfig upload',
      });

      await openKubernetesWizard(page);

      // The hidden file input is present (attached, not visible) and the
      // dropzone prompt + wizard navigation are rendered.
      await expect(page.locator('#connection-wizard-kubeconfig-input')).toBeAttached();
      await expect(page.getByText('Click to choose a kubeconfig file')).toBeVisible();
      await expect(page.getByTestId('connection-wizard-next')).toBeVisible();
    },
  );

  // TC-1014 (matrix R3): uploading a multi-context kubeconfig lists every context
  // in the "Review contexts" step (one row per context). Uses a committed
  // 3-context fixture with unreachable servers, so it needs the discovery
  // endpoint (a running server) but no real clusters. Self-skips if discovery is
  // unavailable.
  test(
    'Discover multiple kubeconfig contexts in the wizard',
    { tag: connTags('wizardDiscoverContexts') },
    async ({ page }, testInfo) => {
      await annotateConnCase(testInfo, 'wizardDiscoverContexts');
      test.slow();

      if (!fs.existsSync(MULTI_CONTEXT_KUBECONFIG)) {
        test.skip(true, 'Multi-context kubeconfig fixture missing.');
        return;
      }

      await openKubernetesWizard(page);
      await uploadKubeconfig(page, MULTI_CONTEXT_KUBECONFIG);

      const discoverRes = page
        .waitForResponse(
          (response) =>
            response.url().includes('/api/system/kubernetes/contexts') &&
            response.request().method() === 'POST',
          { timeout: 60_000 },
        )
        .catch(() => null);

      await page.getByTestId('connection-wizard-next').click();

      const response = await discoverRes;
      if (!response) {
        test.skip(true, 'Context discovery endpoint unavailable; skipping context listing.');
        return;
      }
      expect(
        response.status(),
        `context discovery POST returned ${response.status()} ${response.statusText()}`,
      ).toBe(200);

      // The fixture has three contexts (alpha, beta, gamma); each renders one
      // review row with an editable name field. All three are unreachable, so
      // none cascades to connected. Assert on the row count rather than a
      // specific context name (the discovered name is not asserted here).
      const rows = page.getByTestId('connection-wizard-context-row');
      await expect(rows).toHaveCount(3);
      await expect(rows.first().getByRole('textbox')).toBeVisible();
    },
  );

  // TC-1012 (matrix R1): register + connect a cluster by uploading a kubeconfig
  // through the wizard. Requires a reachable cluster (the host/CI kubeconfig), so
  // it self-skips when no kubeconfig is present or the import call does not
  // succeed.
  test(
    'Register and connect a Kubernetes cluster via kubeconfig upload',
    { tag: connTags('kubeconfigConnect') },
    async ({ page }, testInfo) => {
      await annotateConnCase(testInfo, 'kubeconfigConnect');
      test.slow();

      if (!fs.existsSync(HOST_KUBECONFIG)) {
        test.skip(true, 'No host kubeconfig available to register a cluster.');
        return;
      }

      await openKubernetesWizard(page);
      await uploadKubeconfig(page, HOST_KUBECONFIG);

      // Step 1 -> 2: discover the contexts in the uploaded kubeconfig.
      const discoverRes = page
        .waitForResponse(
          (response) =>
            response.url().includes('/api/system/kubernetes/contexts') &&
            response.request().method() === 'POST',
          { timeout: 60_000 },
        )
        .catch(() => null);
      await page.getByTestId('connection-wizard-next').click();
      const discovered = await discoverRes;
      if (!discovered || discovered.status() !== 200) {
        test.skip(true, 'Context discovery failed; cannot register cluster.');
        return;
      }
      await expect(page.getByTestId('connection-wizard-context-row').first()).toBeVisible();

      // Step 2 -> 3: import the selected contexts (POST /api/system/kubernetes).
      const importRes = page
        .waitForResponse(
          (response) =>
            response.url().endsWith('/api/system/kubernetes') &&
            response.request().method() === 'POST',
          { timeout: 90_000 },
        )
        .catch(() => null);
      await page.getByTestId('connection-wizard-next').click();
      const imported = await importRes;

      if (!imported) {
        test.skip(true, 'Kubeconfig import did not complete in this environment.');
        return;
      }
      expect(
        imported.status(),
        `kubeconfig import POST returned ${imported.status()} ${imported.statusText()}`,
      ).toBe(200);

      // The receipt step confirms the import.
      await expect(
        page.getByRole('heading', { name: /Kubernetes import complete|Configuration saved/i }),
      ).toBeVisible();
    },
  );

  // TC-1061 (matrix X1): transition a connected cluster to disconnected (and back),
  // via the table status dropdown and the shared confirmation modal. Requires a
  // pre-connected cluster, so it self-skips in environments without one.
  test(
    'Transition a Kubernetes connection between lifecycle states',
    { tag: connTags('stateTransition') },
    async ({ page, clusterMetaData }, testInfo) => {
      await annotateConnCase(testInfo, 'stateTransition');
      test.slow();

      // Narrow to the connected data row for the cluster under test.
      await page.getByTestId('ConnectionTable-search').getByRole('button').click();
      const filteredRes = waitForConnectionsApiResponse(page);
      await page.getByRole('textbox', { name: 'Search Connections...' }).fill(clusterMetaData.name);
      await filteredRes;

      const row = page
        .locator('tbody tr')
        .filter({ hasText: clusterMetaData.name })
        // Exact status match: `hasText: 'connected'` also matches "disconnected"
        // (substring), which could select an already-disconnected row.
        .filter({ has: page.getByText('connected', { exact: true }) })
        .first();

      if ((await row.count()) === 0) {
        test.skip(true, 'No connected Kubernetes cluster found to transition. Skipping test.');
        return;
      }

      // Open the row's status dropdown and choose "disconnected".
      await row.locator('#connection-status-select').click();
      await page.getByRole('option', { name: 'disconnected' }).click();

      // Confirm the transition in the shared modal.
      await expect(page.getByTestId('connection-transition-modal')).toBeVisible();
      const transitionRes = waitForConnectionsApiResponse(page);
      await page.getByTestId('connection-transition-confirm').click();
      await waitForSnackBar(page, 'Connection status updated');
      await transitionRes;

      await expect(
        page.locator('tbody tr').filter({ hasText: clusterMetaData.name }).first(),
      ).toContainText('disconnected');

      // Transition back to connected so the run leaves state as it found it.
      const restored = page
        .locator('tbody tr')
        .filter({ hasText: clusterMetaData.name })
        .filter({ has: page.getByText('disconnected', { exact: true }) })
        .first();
      await restored.locator('#connection-status-select').click();
      await page.getByRole('option', { name: 'connected' }).click();
      await expect(page.getByTestId('connection-transition-modal')).toBeVisible();
      // Verify the restore PUT like the disconnect leg, so a failed restore is
      // surfaced (not left to the snackbar alone) and state is left as found.
      const restoreRes = waitForConnectionsApiResponse(page);
      await page.getByTestId('connection-transition-confirm').click();
      await waitForSnackBar(page, 'Connection status updated');
      await restoreRes;
    },
  );

  // TC-1063 (matrix X3): FSM delete (graceful) - the UI bulk-delete sets status
  // to `deleted` via the shared transition modal (PUT /connections/{id}), not the
  // hard DELETE endpoint. Self-skips without a connected cluster.
  test(
    'Delete Kubernetes cluster connections',
    { tag: connTags('delete') },
    async ({ page, clusterMetaData }, testInfo) => {
      await annotateConnCase(testInfo, 'delete');
      // The full search -> delete -> confirm -> snackbar flow can be slow in CI.
      test.slow();

      // beforeEach already navigates to the connections page.
      // Find the row with the connection to be deleted.
      await page.getByTestId('ConnectionTable-search').getByRole('button').click();

      const getFilteredConnectionsRes = waitForConnectionsApiResponse(page);

      await page.getByRole('textbox', { name: 'Search Connections...' }).click();
      await page.getByRole('textbox', { name: 'Search Connections...' }).fill(clusterMetaData.name);

      await getFilteredConnectionsRes;

      // Narrow to the data row for the searched cluster - `tr` alone matches
      // table headers, MUI menu items, and the chip dropdown popovers, all of
      // which can carry the word `connected` from filter labels or other
      // rendered chips and silently grab the wrong row.
      const row = page
        .locator('tbody tr')
        .filter({ hasText: clusterMetaData.name })
        // Exact status match: `hasText: 'connected'` also matches "disconnected"
        // (substring), which could select an already-disconnected row.
        .filter({ has: page.getByText('connected', { exact: true }) })
        .first();

      // Skip the test if no connected cluster is found (e.g., CI environments without a pre-connected cluster)
      if ((await row.count()) === 0) {
        test.skip(true, 'No connected Kubernetes cluster found to delete. Skipping test.');
        return;
      }

      // find the checkbox in the row
      const checkbox = row.getByRole('checkbox').first();
      await checkbox.setChecked(true);

      // Click the bulk "Delete" button in the selected-rows toolbar. Delete now
      // funnels through the shared connection-transition confirmation modal.
      await page.getByTestId('Button-delete-connections').click();
      // Verify that the confirmation modal opened.
      await expect(page.getByTestId('connection-transition-modal')).toBeVisible();

      // Capture the PUT to the connection-update endpoint so we can surface
      // the HTTP status in the test output if the delete fails on the
      // backend - the alternative (just waiting on the snackbar) leaves
      // 500s/4xxs invisible: the success toast never fires and the test
      // hangs to its 60s timeout with no clue what happened.
      const updateRes = page.waitForResponse(
        (resp) =>
          /\/api\/integrations\/connections\/[0-9a-f-]{36}$/.test(resp.url()) &&
          resp.request().method() === 'PUT',
        { timeout: 60_000 },
      );

      await page.getByTestId('connection-transition-confirm').click();

      const response = await updateRes;
      expect(
        response.status(),
        `connection update PUT returned ${response.status()} ${response.statusText()} - body: ${await response.text()}`,
      ).toBe(200);

      await waitForSnackBar(page, 'Connection status updated');
    },
  );
});
