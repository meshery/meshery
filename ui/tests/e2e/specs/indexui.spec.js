import { expect, test } from '@playwright/test';
import { MesheryDashboardPage } from '../fixtures/pages';

test.describe('Index Page UI Tests', () => {
  test('Check Dashboard UI components are displayed', async ({ page }) => {
    const dashboardPage = new MesheryDashboardPage(page);

    await test.step('When I visit the dashboard page', async () => {
      await dashboardPage.goTo();
    });
    await test.step('Then I see the Left Navigation Panel is displayed', async () => {
      await expect(dashboardPage.page.locator('[data-test=navigation]')).toBeVisible();
    });

    await test.step('And the Settings button is displayed', async () => {
      await expect(dashboardPage.page.locator('[data-test=settings-button]')).toBeVisible();
    });

    await test.step('And the Notification button is displayed', async () => {
      await expect(dashboardPage.page.locator('[data-test=notification-button]')).toBeVisible();
    });

    await test.step('And the Profile button is displayed', async () => {
      await expect(dashboardPage.page.locator('[data-test=profile-button]')).toBeVisible();
    });
  });
});
