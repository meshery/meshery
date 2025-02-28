import { expect, test } from '@playwright/test';
import { ENV } from './env';

test.describe('Logout Page Tests', { tag: '@unstable' }, () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}`);
  });
  test('Logout from current user session', async ({ page }) => {
    await page.locator('span:nth-child(5) > div > .MuiButtonBase-root').click();
    await page.getByRole('menuitem', { name: 'Logout' }).click();
    await expect(page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/provider` || `https://cloud.layer5.io/login`,
    );
  });
});
