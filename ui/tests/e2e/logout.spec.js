import { expect, test } from '@playwright/test';
import { ENV } from './env';

test.describe('Logout Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}`);
  });
  test('Logout from current user session', async ({ page }) => {
    await page.getByTestId('header-menu').click();
    await page.getByRole('menuitem', { name: 'Logout' }).click();
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/provider`);
  });
});
