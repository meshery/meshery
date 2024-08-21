import { expect, test } from '@playwright/test';
const { ENV } = require('./env');

test.describe('Navigator section', () => {
  test('should render and navigate to all dashboard and verify their visibilities', async ({
    page,
  }) => {
    // Navigate to "Dashboard" page
    await page.goto(ENV.MESHERY_SERVER_URL);

    // Verify Dashboard page elements are rendered and visible
    await expect(
      page.locator('div').filter({ hasText: 'DashboardLifecycle' }).nth(4),
    ).toBeVisible();
    await expect(page.getByText('DashboardLifecycle')).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Lifecycle' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Configuration' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Configuration' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'meshery-button-2 Performance' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Extensions' })).toBeVisible();

    // Navigate to "Connections" page by clicking on "Lifecycle" then "Connections" Tab
    await page.getByRole('button', { name: 'Lifecycle' }).click();
    await page.getByRole('button', { name: 'Connections' }).click();

    // Verify "Connections" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Lifecycle' })).toBeVisible();
    await expect(page.getByRole('tab', { name: 'Connections' })).toBeVisible();
    await expect(page.getByRole('tab', { name: 'MeshSync' })).toBeVisible();

    // Navigate to "Environments" page by clicking on "Environments" Tab
    await page.getByRole('button', { name: 'Environments' }).click();

    // Verify "Environments" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Environments' })).toBeVisible();

    // Navigate to "Workspaces" page by clicking on "Workspaces" Tab
    await page.getByRole('button', { name: 'Workspaces' }).click();

    // Verify "Workspaces" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Workspaces' })).toBeVisible();

    // Navigate to "Adapters" page by clicking on "Adapters" Tab
    await page.getByRole('button', { name: 'Adapters' }).click();

    // Verify "Adapters" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Adapters' })).toBeVisible();

    // Navigate to "Filters" page by clicking on "Configuration" then "Filters" Tab
    await page.getByRole('button', { name: 'Configuration' }).click();
    await page.getByRole('button', { name: 'Filters' }).click();

    // Verify "Filters" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'FiltersBETA' })).toBeVisible();

    // Navigate to "Design" page by clicking on "Design" Tab
    await page.getByRole('button', { name: 'pattern_trans Designs' }).click();

    // Verify "Design" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Configuration' })).toBeVisible();

    // Navigate to "Profile" page by clicking on "Performance" then "Profiles" Tab
    await page.getByRole('button', { name: 'meshery-button-2 Performance' }).click();
    await page.getByRole('button', { name: 'Profiles', exact: true }).click();

    // Verify "Profiles" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Profiles' })).toBeVisible();

    // Navigate to "Extension" page by clicking on the "Extensions" tab
    await page.getByRole('button', { name: 'Extensions' }).click();

    // Verify "Extensions" page elements are rendered and visible
    await expect(page.getByRole('heading', { name: 'Extensions' })).toBeVisible();
  });
});
