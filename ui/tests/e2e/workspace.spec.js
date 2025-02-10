import { expect, test } from './fixtures/project';
import { ENV } from './env';

const verifyWorkspaceResBody = (body, provider) => {
  expect(body).toEqual(
    expect.objectContaining({
      workspaces: expect.any(Array),
      total_count: expect.any(Number),
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
        : {
            name: expect.any(String),
            description: expect.any(String),
            org_id: expect.any(String),
            created_at: expect.any(String),
            updated_at: expect.any(String),
            owner: expect.any(String),
            deleted_at: expect.anything(),
        };

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
    if (!response.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/workspaces`)) return false;
    if (response.status() !== 200) return false;
    const body = await response.json();
    if (body.workspaces) return true;
    return false;
  });
  await page.goto(`${ENV.MESHERY_SERVER_URL}`);
  await page.getByRole('button', { name: 'Lifecycle' }).click();
  await page.getByRole('button', { name: 'Workspaces' }).click();

  await workspacesReq;
  const res = await workspacesRes;
  const body = await res.json();
  verifyWorkspaceResBody(body, provider);
  workspaceCount = body.workspaces.length;
});

test('Verify workspace page components', async ({ page }) => {
  // Verify Create button exists
  await expect(page.getByRole('button').filter({ hasText: 'Create' })).toBeVisible();

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

test('Switch between grid and table views', async ({ page }) => {
  // Verify initial grid view
  await expect(page.locator('.MuiGrid-container').first()).toBeVisible();

  // Switch to table view using aria-label
  await page.getByRole('button', { name: 'Switch View' }).click();

  // Verify table structure exists (simpler approach)
  await expect(page.locator('.MuiTable-root').first()).toBeVisible();
  await expect(page.locator('.MuiTableHead-root').first()).toBeVisible();
  await expect(page.locator('.MuiTableBody-root').first()).toBeVisible();

  const expectedHeaders = ['Name', 'Description', 'Owner'];
  for (const header of expectedHeaders) {
    await expect(
      page.locator('[class*="MUIDataTableHeadCell-data"]', { hasText: header }).first()
    ).toBeVisible();
  }

  // Switch back to grid view
  await page.getByRole('button', { name: 'Switch View' }).click();
  await expect(page.locator('.MuiGrid-container').first()).toBeVisible();
});

test('Create new workspace', async ({ page, provider }) => {
  const createWorkspaceRes = page.waitForResponse(
    (response) =>
      response.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces`) &&
      response.request().method() === 'POST',
  );

  // Click Create button
  await page.getByRole('button').filter({ hasText: 'Create' }).click();

  // Fill workspace form
  await page.getByLabel('Name').fill('Test Workspace');
  await page.getByLabel('Description').fill('Test workspace description');

  // Submit form
  await page.getByRole('button', { name: 'Save' }).click();

  // Wait for response
  const response = await createWorkspaceRes;
  const responseBody = await response.json();

  // Verify response payload
  expect(responseBody).toMatchObject({
    name: 'Test Workspace',
    description: 'Test workspace description',
  });

  // Verify success notification
  await expect(page.getByText('Workspace "Test Workspace" created')).toBeVisible();

  // Verify the new workspace is visible in the UI
  const workspaceCard = page.locator('.MuiGrid-container').filter({ 
    has: page.getByText('Test Workspace')
  });
  await expect(workspaceCard).toBeVisible();

  // Verify workspace details in the card
  await expect(workspaceCard.getByText('Test Workspace').first()).toBeVisible();
  await expect(workspaceCard.getByText('Test workspace description').first()).toBeVisible();
});

test('Update workspace', async ({ page }) => {
  const updateWorkspaceRes = page.waitForResponse(
    (response) =>
      response.url().includes(`${ENV.MESHERY_SERVER_URL}/api/workspaces`) &&
      response.request().method() === 'PUT',
  );

  // Click edit button on the first workspace card
  const workspaceCard = page.locator('.MuiGrid-container').first();
  await workspaceCard.getByTestId('edit-icon').click();

  // Fill workspace form
  await page.getByLabel('Name').fill('Updated Workspace');
  await page.getByLabel('Description').fill('Updated workspace description');

  // Submit form
  await page.getByRole('button', { name: 'Save' }).click();

  // Wait for response
  const response = await updateWorkspaceRes;
  const responseBody = await response.json();

  // Verify response payload
  expect(responseBody).toMatchObject({
    name: 'Updated Workspace',
    description: 'Updated workspace description',
  });

  // Verify success notification
  await expect(page.getByText('Workspace updated successfully')).toBeVisible();

  // Verify the updated workspace is visible in the UI
  const updatedCard = page.locator('.MuiGrid-container').filter({ 
    has: page.getByText('Updated Workspace')
  });
  await expect(updatedCard).toBeVisible();
  await expect(updatedCard.getByText('Updated workspace description')).toBeVisible();
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