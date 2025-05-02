import { expect, test } from './fixtures/project';
import { ENV } from './env';
import os from 'os';

// name: Name of the test
// transitionOption: Option to be chosen from dropdown to transition to another state
// statusAfterTransition: Text shown in current state after transition
// restorationOption: Option to be chosen from dropdown to transition back to connected state
const transitionTests = [
  {
    name: 'Transition to disconnected state and then back to connected state',
    transitionOption: 'disconnected',
    statusAfterTransition: 'disconnected',
    restorationOption: 'connected',
  },
  {
    name: 'Transition to ignored state and then back to connected state',
    transitionOption: 'ignored',
    statusAfterTransition: 'ignored',
    restorationOption: 'registered',
  },
  {
    name: 'Transition to not found state and then back to connected state',
    transitionOption: 'not found',
    statusAfterTransition: 'not found',
    restorationOption: 'discovered',
  },
];

test.describe.serial('Connection Management Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: 'Lifecycle' }).click();
  });
  test('Verify that UI components are displayed', async ({ page }) => {
    // Verify that connections table is displayed (by checking for table headings)
    for (const heading of ['Name', 'Environments', 'Kind', 'Category', 'Status', 'Actions']) {
      await expect(page.getByRole('columnheader', { name: heading })).toBeVisible();
    }
  });

  test('Add a cluster connection by uploading kubeconfig file', async ({ page }) => {
    await page.getByRole('tab', { name: 'Connections' }).click();

    const addConnectionReq = page.waitForRequest(
      (request) =>
        request.url() === `${ENV.MESHERY_SERVER_URL}/api/system/kubernetes` &&
        request.method() === 'POST',
    );
    const addConnectionRes = page.waitForResponse(
      (response) =>
        response.url() === `${ENV.MESHERY_SERVER_URL}/api/system/kubernetes` &&
        response.status() === 200,
    );
    await page.getByTestId('connection-addCluster').click();

    // Verify "Add Kubernetes Cluster(s)" modal opens
    await expect(page.getByTestId('connection-addKubernetesModal')).toBeVisible();

    const fileChooserPromise = page.waitForEvent('filechooser');
    await page.getByTestId('connection-uploadKubeConfig').click();
    const fileChooser = await fileChooserPromise;
    // Attach existing kubeconfig file of the system, to test the upload without making any changes in configuration
    const kubeConfigPath = `${os.homedir()}/.kube/config`;
    await fileChooser.setFiles(kubeConfigPath);

    await page.getByRole('button', { name: 'IMPORT', exact: true }).click();

    await addConnectionReq;
    await addConnectionRes;

    await expect(page.getByTestId('connection-discoveredModal')).toBeVisible();

    await page.getByRole('button', { name: 'OK' }).click();

    const connectedItem = page.getByRole('menuitem', { name: 'connected' }).first();
    await connectedItem.scrollIntoViewIfNeeded();
    await expect(connectedItem).toBeVisible();
  });

  transitionTests.forEach((t) => {
    test(t.name, async ({ page }) => {
      const stateTransitionReq = page.waitForRequest(
        (request) =>
          request.url() ===
            `${ENV.MESHERY_SERVER_URL}/api/integrations/connections/kubernetes/status` &&
          request.method() === 'PUT',
      );

      const stateTransitionRes = page.waitForResponse(
        (response) =>
          response.url() ===
            `${ENV.MESHERY_SERVER_URL}/api/integrations/connections/kubernetes/status` &&
          response.status() === 202,
      );

      const getConnectionsReq = page.waitForRequest(
        (request) =>
          request.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/integrations/connections`) &&
          request.method() === 'GET',
      );

      const firstRow = page.getByRole('menuitem', { name: 'connected' }).first();
      await firstRow.scrollIntoViewIfNeeded();
      await expect(firstRow).toBeVisible();
      // ===== TRANSITIONING TO A NEW STATE =====

      // open state transition options dropdown
      await firstRow.locator('span', { hasText: 'connected' }).click();
      await page.getByRole('option', { name: t.transitionOption }).click();

      await expect(page.locator('#searchClick')).toBeVisible();
      await page.getByRole('button', { name: 'Confirm' }).click();

      await expect(page.getByText('Connection status updated')).toBeVisible();

      await stateTransitionReq;
      await stateTransitionRes;
      await getConnectionsReq;
      // expect new state to be shown as current state

      const updatedFirstRow = page.getByRole('menuitem', { name: t.statusAfterTransition }).first();
      await expect(updatedFirstRow).toBeVisible();

      // ===== TRANSITION BACK TO "connected" STATE =====
      // open state transition options dropdown again
      await updatedFirstRow.locator('span', { hasText: t.statusAfterTransition }).click();
      await page.getByRole('option', { name: t.restorationOption }).click();

      await expect(page.locator('#searchClick')).toBeVisible();
      await page.getByRole('button', { name: 'Confirm' }).click();
      // expect the state to be restored to "connected"

      const restoredFirstRow = page.getByRole('menuitem', { name: 'connected' }).first();
      await expect(restoredFirstRow).toBeVisible();

      await expect(page.getByText('Connection status updated')).toBeVisible();
    });
  });

  test('Delete Kubernetes cluster connections', async ({ page }) => {
    await page.getByRole('tab', { name: 'Connections' }).click();

    const firstRow = page
      .locator('tr', { has: page.getByRole('menuitem', { name: 'connected' }) })
      .first();
    await firstRow.scrollIntoViewIfNeeded();
    await expect(firstRow).toBeVisible();

    if ((await firstRow.count()) === 0) {
      throw new Error(
        'No connected Kubernetes cluster found to delete. Ensure a connection exists before running this test.',
      );
    }

    const checkbox = firstRow.locator('input[type="checkbox"]').first();
    await checkbox.check();

    await page.getByRole('button', { name: 'Delete', exact: true }).click();
    await expect(page.getByText('Delete Connections')).toBeVisible();

    const responsePromise = page.waitForResponse(
      (response) =>
        response
          .url()
          .startsWith(`${ENV.MESHERY_SERVER_URL}/api/integrations/connections/kubernetes/status`) &&
        response.status() === 202,
    );

    await page.getByRole('button', { name: 'DELETE', exact: true }).click();
    await expect(page.getByText('Connection status updated')).toBeVisible();

    await responsePromise;
  });
});
