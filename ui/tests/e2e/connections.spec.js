import { expect, test } from './fixtures/project';
import { ENV } from './env';
import os from 'os';
import { waitForSnackBar } from './utils/waitForSnackBar';

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
    await page.getByTestId('connection-addCluster').waitFor();
  });
  test('Verify that UI components are displayed', async ({ page }) => {
    // Verify that connections table is displayed (by checking for table headings)
    for (const heading of ['Name', 'Environments', 'Kind', 'Category', 'Status', 'Actions']) {
      await expect(page.getByRole('columnheader', { name: heading })).toBeVisible();
    }
  });

  test('Add a cluster connection by uploading kubeconfig file', async ({
    page,
    clusterMetaData,
  }) => {
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

    await page.getByRole('button', { name: 'OK' }).click();

    // Search for the newly added cluster
    await page.getByTestId('ConnectionTable-search').getByRole('button').click();

    await page.getByRole('textbox', { name: 'Search Connections...' }).click();
    await page.getByRole('textbox', { name: 'Search Connections...' }).fill(clusterMetaData.name);

    const newConnectionRow = page.getByRole('menuitem', { hasText: clusterMetaData.name }).first();
    await expect(newConnectionRow).toContainText('connected');
  });

  transitionTests.forEach((t) => {
    test(t.name, async ({ page, clusterMetaData }) => {
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

      await page.getByTestId('ConnectionTable-search').getByRole('button').click();
      await page.getByRole('textbox', { name: 'Search Connections...' }).fill(clusterMetaData.name);

      const matchingRows = page.getByRole('menuitem', { hasText: clusterMetaData.name });

      const connectedRow = matchingRows
        .filter({
          has: page.locator('span', { hasText: 'connected' }),
        })
        .first();

      await expect(connectedRow).toBeVisible();

      // ===== TRANSITIONING TO A NEW STATE =====

      // open state transition options dropdown
      await connectedRow.locator('span', { hasText: 'connected' }).click();
      await page.getByRole('option', { name: t.transitionOption }).click();

      await expect(page.locator('#searchClick')).toBeVisible();
      await page.getByRole('button', { name: 'Confirm' }).click();

      await waitForSnackBar(page, 'Connection status updated');

      await stateTransitionReq;
      await stateTransitionRes;
      await getConnectionsReq;
      // expect new state to be shown as current state

      const updatedConnection = page
        .getByRole('menuitem', { hasText: clusterMetaData.name })
        .first();
      await expect(updatedConnection).toContainText(t.statusAfterTransition);

      // ===== TRANSITION BACK TO "connected" STATE =====
      // open state transition options dropdown again
      await updatedConnection.locator('span', { hasText: t.statusAfterTransition }).click();
      await page.getByRole('option', { name: t.restorationOption }).click();

      await expect(page.locator('#searchClick')).toBeVisible();
      await page.getByRole('button', { name: 'Confirm' }).click();

      await waitForSnackBar(page, 'Connection status updated');
    });
  });
  test('Delete Kubernetes cluster connections', async ({ page, clusterMetaData }) => {
    // Navigate to 'Connections' tab
    await page.getByRole('tab', { name: 'Connections' }).click();
    // Find the row with the connection to be deleted
    await page.getByTestId('ConnectionTable-search').getByRole('button').click();

    await page.getByRole('textbox', { name: 'Search Connections...' }).click();
    await page.getByRole('textbox', { name: 'Search Connections...' }).fill(clusterMetaData.name);

    const row = page.locator('tr').filter({ hasText: 'connected' }).first();

    // Fail the test if the connection is not found
    if ((await row.count()) === 0) {
      throw new Error(
        'No connected Kubernetes cluster found to delete. Ensure a connection exists before running this test.',
      );
    }

    //find the checkbox in the row
    const checkbox = row.locator('input[type="checkbox"]').first();
    await checkbox.click();

    // Click "Delete" button in the table
    await page.getByRole('button', { name: 'Delete', exact: true }).click();
    // Verify that Confirmation modal opened and delete
    await expect(page.getByText('Delete Connections')).toBeVisible();
    await page.getByRole('button', { name: 'DELETE', exact: true }).click();

    await waitForSnackBar(page, 'Connection status updated');
  });
});
