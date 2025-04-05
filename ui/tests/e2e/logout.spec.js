import { expect, test } from '@playwright/test';

test.describe('Logout Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
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

    await page.getByTestId('header-menu').click();
    await page.getByRole('menuitem', { name: 'Logout' }).click();

    const request = await waitForLogoutRequest;
    expect(request.url()).toContain('/user/logout');

    await page.waitForURL('/provider');
    expect(page.url()).toContain(`/provider`);
  });
});
