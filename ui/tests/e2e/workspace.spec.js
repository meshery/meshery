import { expect, test } from './fixtures/project';
import { ENV } from './env';

const verifyWorkspaceResBody = (body, provider) => {
  expect(body).toEqual(
    expect.objectContaining({
      workspaces: expect.any(Array),
      total_count: expect.any(Number),
      page: expect.any(Number),
      page_size: expect.any(Number),
    }),
  );

  for (const workspace of body.workspaces) {
    const mesheryKeys =
      provider === 'Meshery'
        ? {
            name: expect.any(String),
            description: expect.any(String),
            organization_id: expect.any(String),
            created_at: expect.any(String),
            updated_at: expect.any(String),
          }
        : null;

    expect(workspace).toEqual(
      expect.objectContaining({
        id: expect.any(String),
        ...mesheryKeys,
      }),
    );
  }
};

test.describe.configure({ mode: 'serial' });
let workspaceCount = 0;

test.beforeEach(async ({ page, provider }) => {
  const workspacesReq = page.waitForRequest(
    (request) =>
      request.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces`) &&
      request.method() === 'GET',
  );
  const workspacesRes = page.waitForResponse(async (response) => {
    if (!response.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces`)) return false;
    if (response.status() !== 200) return false;
    const body = await response.json();
    if (body.workspaces) return true;
    return false;
  });

  await page.goto(`${ENV.MESHERY_SERVER_URL}/workspace`);

  await workspacesReq;
  const res = await workspacesRes;
  const body = await res.json();
  verifyWorkspaceResBody(body, provider);

  workspaceCount = body.workspaces.length;
});

test('Verify workspace page components', async ({ page }) => {
  // Verify Create button exists
  await expect(page.getByRole('button').filter({ hasText: 'Create' })).toBeVisible();

  // Verify Search bar exists
  await expect(page.getByPlaceholder('Search Workspaces...')).toBeVisible();

  if (workspaceCount > 0) {
    // Verify workspace cards are displayed
    await expect(page.locator('.MuiGrid-container').first()).toBeVisible();
  } else {
    // Verify empty state
    await expect(page.getByText('No workspace available')).toBeVisible();
    await expect(
      page.getByText('Click "Create" to establish your first workspace.'),
    ).toBeVisible();
  }
});

test('Create new workspace', async ({ page, provider }) => {
  const createWorkspaceReq = page.waitForRequest(
    (request) =>
      request.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces`) &&
      request.method() === 'POST',
  );

  // Click Create button
  await page.getByRole('button').filter({ hasText: 'Create' }).click();

  // Fill workspace form
  await page.getByLabel('Name').fill('Test Workspace');
  await page.getByLabel('Description').fill('Test workspace description');

  // Submit form
  await page.getByRole('button', { name: 'Save' }).click();

  // Wait for request to complete
  const request = await createWorkspaceReq;
  const requestBody = JSON.parse(request.postData());

  // Verify request payload
  expect(requestBody).toEqual(
    expect.objectContaining({
      name: 'Test Workspace',
      description: 'Test workspace description',
    }),
  );

  // Verify success notification
  await expect(page.getByText('Workspace "Test Workspace" created')).toBeVisible();
});

test('Assign environments to workspace', async ({ page }) => {
  // Find and click assign environment button on workspace card
  const workspaceCard = page.locator('.MuiGrid-container').first();
  await workspaceCard.getByTestId('environment-icon').click();

  // Verify transfer list modal
  await expect(page.getByText(/Assign Environments to/)).toBeVisible();

  // Verify transfer list components
  await expect(page.getByText('Available Environments')).toBeVisible();
  await expect(page.getByText('Assigned Environments')).toBeVisible();

  // Close modal
  await page.getByRole('button', { name: 'Cancel' }).click();
});

test('Assign designs to workspace', async ({ page }) => {
  // Find and click assign design button on workspace card
  const workspaceCard = page.locator('.MuiGrid-container').first();
  await workspaceCard.getByTestId('design-icon').click();

  // Verify transfer list modal
  await expect(page.getByText(/Assign Designs to/)).toBeVisible();

  // Verify transfer list components
  await expect(page.getByText('Available Designs')).toBeVisible();
  await expect(page.getByText('Assigned Designs')).toBeVisible();

  // Close modal
  await page.getByRole('button', { name: 'Cancel' }).click();
});

test('Delete workspace', async ({ page }) => {
  // Find the workspace card
  const workspaceCard = page.locator('.MuiGrid-container').first();
  
  // Click checkbox to select workspace
  await workspaceCard.getByRole('checkbox').click();

  // Verify bulk action wrapper appears
  await expect(page.getByText(/workspace selected/)).toBeVisible();

  // Click delete icon
  await page.getByTestId('delete-icon').click();

  // Verify delete confirmation modal
  await expect(page.getByText(/Delete Workspace/)).toBeVisible();

  const deleteReq = page.waitForRequest(
    (request) =>
      request.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces`) &&
      request.method() === 'DELETE',
  );

  // Confirm deletion
  await page.getByRole('button', { name: 'DELETE' }).click();

  await deleteReq;

  // Verify success notification
  await expect(page.getByText(/Workspace .* deleted/)).toBeVisible();
});

test('Search workspaces', async ({ page }) => {
  const searchBar = page.getByPlaceholder('Search Workspaces...');
  
  // Enter search term
  await searchBar.fill('Test');

  // Wait for search request
  const searchReq = page.waitForRequest((request) =>
    request.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces?search=Test`),
  );

  await searchReq;

  // Verify filtered results
  const workspaceCards = page.locator('.MuiGrid-container').first();
  await expect(workspaceCards).toContainText('Test');
});

