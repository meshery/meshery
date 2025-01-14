import { expect, test } from './fixtures/project';
import { ENV } from './env';
import os from 'os';

const verifyConnectionsResBody = (body, provider) => {
  expect(body).toEqual(
    expect.objectContaining({
      connections: expect.any(Array),
      total_count: expect.any(Number),
      page: expect.any(Number),
      page_size: expect.any(Number),
    }),
  );
  for (const connection of body.connections) {
    const mesheryKeys =
      provider === 'Meshery'
        ? {
            sub_type: expect.any(String),
            metadata: expect.anything(),
            kind: expect.any(String),
            name: expect.any(String),
            status: expect.any(String),
            type: expect.any(String),
          }
        : null;

    expect(connection).toEqual(
      expect.objectContaining({
        id: expect.any(String),
        credential_id: expect.any(String),
        user_id: expect.any(String),
        created_at: expect.any(String),
        updated_at: expect.any(String),
        deleted_at: expect.objectContaining({
          Time: expect.any(String),
          Valid: expect.any(Boolean),
        }),
        ...mesheryKeys,
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

test.describe.configure({ mode: 'serial' });
let connectionCount = 0;

test.beforeEach(async ({ page, provider }) => {
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
  await page.goto(`${ENV.MESHERY_SERVER_URL}`);
  await page.getByRole('button', { name: 'Lifecycle' }).click();
  await page.getByRole('button', { name: 'Connections' }).click();

  // Verify requests and responses expected on initial page load
  await connectionsReq;
  const res = await connectionsRes;
  const body = await res.json();
  verifyConnectionsResBody(body, provider);

  connectionCount = body.connections.length;
});

test('Verify that UI components are displayed', async ({ page }) => {
  // Verify that connections table is displayed (by checking for table headings)
  for (const heading of ['Name', 'Environments', 'Kind', 'Category', 'Status', 'Actions']) {
    await expect(page.getByRole('columnheader', { name: heading })).toBeVisible();
  }

  // Verify that all connections returned by server are displayed (by counting number of rows in the table)
  expect((await page.locator('tr').count()) - 2).toEqual(connectionCount); // -2 for not considering header and footer
});

test(
  'Add a cluster connection by uploading kubeconfig file',
  { tag: '@unstable' },
  async ({ page }) => {
    // Navigate to 'Connections' tab
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

    // Click Add Cluster button
    await page.getByTestId('connection-addCluster').click();

    // Verify "Add Kubernetes Cluster(s)" modal opens
    await expect(page.getByTestId('connection-addKubernetesModal')).toBeVisible();

    const fileChooserPromise = page.waitForEvent('filechooser');
    await page.getByTestId('connection-uploadKubeConfig').click();
    const fileChooser = await fileChooserPromise;

    // Attach existing kubeconfig file of the system, to test the upload without making any changes in configuration
    const kubeConfigPath = `${os.homedir()}/.kube/config`;
    await fileChooser.setFiles(kubeConfigPath);

    // Click "IMPORT" button
    await page.getByRole('button', { name: 'IMPORT', exact: true }).click();

    // Verify requests and responses
    await addConnectionReq;
    await addConnectionRes;

    // Verify displaying of success modal
    await expect(page.getByTestId('connection-discoveredModal')).toBeVisible();

    // Verify available contexts were connected
    await expect(page.getByRole('menuitem', { name: 'connected' })).toBeVisible();

    // Click "OK" button to close success modal
    await page.getByRole('button', { name: 'OK', exact: true }).click();
  },
);

transitionTests.forEach((t) => {
  test(t.name, { tag: '@unstable' }, async ({ page, provider }) => {
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
    verifyConnectionsResBody(body, provider);

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

test('Delete Kubernetes cluster connections', { tag: '@unstable' }, async ({ page }) => {
  // Navigate to 'Connections' tab
  await page.getByRole('tab', { name: 'Connections' }).click();
  // Find the row with the connection to be deleted
  const row = page.locator('tr').filter({ hasText: 'connected' }).first();

  // Fail the test if the connection is not found
  if ((await row.count()) === 0) {
    throw new Error(
      'No connected Kubernetes cluster found to delete. Ensure a connection exists before running this test.',
    );
  }

  //find the checkbox in the row
  const checkbox = row.locator('input[type="checkbox"]').first();
  await checkbox.check();

  // Click "Delete" button in the table
  await page.getByRole('button', { name: 'Delete', exact: true }).click();
  // Verify that Confirmation modal opened and delete
  await expect(page.getByText('Delete Connections')).toBeVisible();

  const responsePromise = page.waitForResponse(
    (response) =>
      response
        .url()
        .startsWith(`${ENV.MESHERY_SERVER_URL}/api/integrations/connections/kubernetes/status`) &&
      response.status() === 202,
  );

  await page.getByRole('button', { name: 'Delete', exact: true }).click();

  await responsePromise;
});
