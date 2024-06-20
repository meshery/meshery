import { expect, test } from '@playwright/test';
import { ENV } from './env';
import os from 'os';

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

test.describe('Connections Page Tests', () => {
  let connectionCount = 0;

  test.beforeEach(async ({ page }) => {
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
    verifyConnectionsResBody(body);

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

  test('Add a cluster connection by uploading kubeconfig file', async ({ page }) => {
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
    await page.getByRole('button', { name: 'Add Cluster' }).click();

    // Verify "Add Kubernetes Cluster(s)" modal opens
    await expect(page.getByRole('heading', { name: 'Add Kubernetes Cluster(s)' })).toBeVisible();

    // Attach existing kubeconfig file of the system, to test the upload without making any changes in configuration
    const kubeConfigPath = `${os.homedir()}/.kube/config`;
    await page.locator('input[type="file"]').setInputFiles(kubeConfigPath);

    // Click "IMPORT" button
    await page.getByRole('button', { name: 'IMPORT', exact: true }).click();

    // Verify requests and responses
    await addConnectionReq;
    await addConnectionRes;

    // Verify displaying of success modal
    await expect(page.getByText('Available contexts in')).toBeVisible();

    // Verify available contexts were connected
    await expect(page.getByRole('menuitem', { name: 'connected' })).toBeVisible();

    // Click "OK" button to close success modal
    await page.getByRole('button', { name: 'OK', exact: true }).click();
  });
});
