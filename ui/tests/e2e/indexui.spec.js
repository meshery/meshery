import { expect, test } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';

test.describe('Index Page UI Tests', () => {
  let dashboardPage;
  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
  });

  test('Test if Left Navigation Panel is displayed', async () => {
    await expect(dashboardPage.navigationPanel).toBeVisible();
  });

  test('Test if Notification button is displayed', async () => {
    await expect(dashboardPage.notificationButton).toBeVisible();
  });

  test('Test if Profile button is displayed', async () => {
    await expect(dashboardPage.profileButton).toBeVisible();
  });
});
