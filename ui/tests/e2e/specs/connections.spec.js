import { expect, test } from '@playwright/test';
import { ENV } from '../env';
import os from 'os';
import { ConnectionsPage } from '../fixtures/pages';

const verifyConnectionsResBody = (body) => {
  expect(body).toEqual(
    expect.objectContaining({
      connections: expect.any(Array),
      total_count: expect.any(Number),
      page: expect.any(Number),
      page_size: expect.any(Number),
    }),
  );
  for (const connection of body.connections) {
    expect(connection).toEqual(
      expect.objectContaining({
        id: expect.any(String),
        name: expect.any(String),
        credential_id: expect.any(String),
        type: expect.any(String),
        sub_type: expect.any(String),
        kind: expect.any(String),
        metadata: expect.anything(),
        status: expect.any(String),
        user_id: expect.any(String),
        created_at: expect.any(String),
        updated_at: expect.any(String),
        deleted_at: expect.objectContaining({
          Time: expect.any(String),
          Valid: expect.any(Boolean),
        }),
      }),
    );
  }
};

// name: Name of the test
// transitionOption: Option to be chosen from dropdown to transition to another state
// statusAfterTransition: Text shown in current state after transition
// restorationOption: Option to be chosen from dropdown to transition back to connected state
const transitionTests = [
  {
    name: 'Transition to disconnected state and then back to connected state',
    transitionOption: 'Disconnect',
    statusAfterTransition: 'disconnected',
    restorationOption: 'Connect',
  },
  {
    name: 'Transition to ignored state and then back to connected state',
    transitionOption: 'Ignore',
    statusAfterTransition: 'ignored',
    restorationOption: 'Register',
  },
  {
    name: 'Transition to not found state and then back to connected state',
    transitionOption: 'Not Found',
    statusAfterTransition: 'not found',
    restorationOption: 'Discover',
  },
];

test.describe.serial('Connections Page Tests', () => {
  let connectionCount = 0;

  test.beforeEach(async ({ page }) => {
    const connectionsPage = new ConnectionsPage(page);
    const connectionsReq = page.waitForRequest(
      (request) =>
        request.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/integrations/connections`) &&
        request.method() === 'GET',
    );
    const connectionsRes = page.waitForResponse(async (response) => {
      if (!response.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/integrations/connections`))
        return false;
      if (response.status() !== 200) return false;
      const body = await response.json();
      if (body.connections && body.connections.length > 0) return true;
      else return false;
    });

    // Visit Connections Page
    await connectionsPage.goTo();

    // Verify requests and responses expected on initial page load
    await connectionsReq;
    const res = await connectionsRes;
    const body = await res.json();
    verifyConnectionsResBody(body);

    connectionCount = body.connections.length;
  });

  test('Verify that UI components are displayed', async ({ page }) => {
    await test.step('I see the connections table is displayed', async () => {
      for (const heading of ['Name', 'Environments', 'Kind', 'Category', 'Status', 'Actions']) {
        await expect(page.getByRole('columnheader', { name: heading })).toBeVisible();
      }
    });

    await test.step('And I see the available connections are displayed in the table', async () => {
      const actualConnectionCount = (await page.locator('tr').count()) - 2; // -2 for not considering header and footer
      expect(actualConnectionCount).toEqual(connectionCount);
    });
  });

  test('Add a cluster connection by uploading kubeconfig file', async ({ page }) => {
    let addConnectionReq;
    let addConnectionRes;
    await test.step('When I navigate to connections tab', async () => {
      await page.getByRole('tab', { name: 'Connections' }).click();
    });

    await test.step('Initiate add connection request', async () => {
      addConnectionReq = page.waitForRequest(
        (request) =>
          request.url() === `${ENV.MESHERY_SERVER_URL}/api/system/kubernetes` &&
          request.method() === 'POST',
      );
      addConnectionRes = page.waitForResponse(
        (response) =>
          response.url() === `${ENV.MESHERY_SERVER_URL}/api/system/kubernetes` &&
          response.status() === 200,
      );
    });

    await test.step('And I click the Add Cluster', async () => {
      await page.getByRole('button', { name: 'Add Cluster' }).click();
    });

    await test.step('Then I verify `Add Kubernetes Cluster(s)` modal opens', async () => {
      await expect(page.getByText('Add Kubernetes Cluster(s)')).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Upload your kubeconfig' })).toBeVisible();
    });

    await test.step('And I attach the existing kubeconfig file of the system, to test the upload withou making changes in the configuration', async () => {
      const kubeConfigPath = `${os.homedir()}/.kube/config`;
      await page.locator('input[type="file"]').setInputFiles(kubeConfigPath);
    });

    await test.step('And I click `import` button', async () => {
      await page.getByRole('button', { name: 'IMPORT', exact: true }).click();
    });

    test.step('And I verify the requests and responses', async () => {
      await addConnectionReq;
      await addConnectionRes;
    });

    await test.step('And I see a success modal', async () => {
      await expect(page.getByText('Available contexts in')).toBeVisible();
    });

    await test.step('And I see that available contexts were connected', async () => {
      await expect(page.getByRole('menuitem', { name: 'connected' })).toBeVisible();
    });

    await test.step('And I click `OK` button to close the modal`', async () => {
      await page.getByRole('button', { name: 'OK', exact: true }).click();
    });
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

      const getConnectionsRes = page.waitForResponse(
        (response) =>
          response.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/integrations/connections`) &&
          response.status() === 200,
      );

      // since test run serially, and latest connection appears at topmost,
      // the connection created in the previous test will appear in the first row of the table
      const firstRow = page.locator('tbody').locator('tr').first();

      // expect connected state initially
      await expect(firstRow.locator('span', { hasText: 'connected' })).toBeVisible();

      // ===== TRANSITIONING TO A NEW STATE =====

      // open state transition options dropdown
      await firstRow.locator('span', { hasText: 'connected' }).click();

      // click required option
      await page.getByText(t.transitionOption, { exact: true }).click();

      // verify that Confirmation modal opened
      await expect(page.getByText('Connection Status Transition')).toBeVisible();

      // click "Confirm"
      await page.getByRole('button', { name: 'Confirm' }).click();

      // verify API requests and responses
      await stateTransitionReq;
      await stateTransitionRes;
      await getConnectionsReq;
      const res = await getConnectionsRes;
      const body = await res.json();
      verifyConnectionsResBody(body);

      // expect new state to be shown as current state
      await expect(firstRow.locator('span', { hasText: t.statusAfterTransition })).toBeVisible();

      // ===== TRANSITION BACK TO "connected" STATE =====

      // open state transition options dropdown again
      await firstRow.locator('span', { hasText: t.statusAfterTransition }).click();

      // click the option required to transition back to "connected" state
      await page.getByText(t.restorationOption, { exact: true }).click();

      // verify that Confirmation modal opened again
      await expect(page.getByText('Connection Status Transition')).toBeVisible();

      // click "Confirm"
      await page.getByRole('button', { name: 'Confirm' }).click();

      // expect the state to be restored to "connected"
      await expect(firstRow.locator('span', { hasText: 'connected' })).toBeVisible();
    });
  });
});
