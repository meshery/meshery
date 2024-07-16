const { test, expect, MesheryDashboardPage } = require('../fixtures/pages');

test('Aggregation Charts are displayed', async ({ page }) => {
  const dashboardPage = new MesheryDashboardPage(page);

  await test.step('When I visit the Dashboard Page', async () => {
    await dashboardPage.goTo();
  });

  await test.step('Then I verify the presence of `Models by Category Chart`', async () => {
    await expect(
      dashboardPage.page.getByRole('heading', { name: /Models by Category/i }),
    ).toBeVisible();
    await expect(dashboardPage.page.locator('tspan').filter({ hasText: /Models/i })).toBeVisible();
    await expect(dashboardPage.page.getByText('by Category', { exact: true })).toBeVisible();
  });

  await test.step('Then I verify the presence of `Registry Chart`', async () => {
    await expect(dashboardPage.page.getByRole('heading', { name: 'Registry' })).toBeVisible();
    await expect(dashboardPage.page.getByText('Registered', { exact: true })).toBeVisible();
    await expect(dashboardPage.page.getByText(/Capabilities/i)).toBeVisible();
    await expect(dashboardPage.page.getByText(/by Type/i).first()).toBeVisible();
  });

  await test.step('Then I verify the presence of `Connections Chart`', async () => {
    await expect(dashboardPage.page.getByRole('heading', { name: 'Connections' })).toBeVisible();
    await expect(
      dashboardPage.page.locator('tspan').filter({ hasText: /Connections/i }),
    ).toBeVisible();
    await expect(dashboardPage.page.getByText(/by Status/i)).toBeVisible();
  });

  await test.step('Then I verify the presence of `Configuration Chart`', async () => {
    await expect(dashboardPage.page.getByRole('heading', { name: 'Configuration' })).toBeVisible();
    await expect(dashboardPage.page.getByText('Content', { exact: true })).toBeVisible();
    await expect(dashboardPage.page.getByText(/by Type/i).nth(1)).toBeVisible();
  });
});
