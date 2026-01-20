import { expect, test } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';

test.describe('Logout Page Tests', () => {
  let dashboardPage;
  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
  });
  test('Logout from current user session', async ({ page }) => {
    await page.route('/user/logout', (route) => {
      route.fulfill({
        status: 302,
        headers: {
          Location: '/provider',
        },
      });
    });

    const waitForLogoutRequest = page.waitForRequest('/user/logout');

    await dashboardPage.navigateToLogout();

    const request = await waitForLogoutRequest;
    expect(request.url()).toContain('/user/logout');

    await page.waitForURL('/provider');
    expect(page.url()).toContain(`/provider`);
  });
});
