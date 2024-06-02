import { expect, test } from '@playwright/test';
import { ENV } from '../env';

test.describe('Test if UI components are displayed on Index Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(ENV.MESHERY_SERVER_URL);
  });

  test('Left Navigation Panel', async ({ page }) => {
    await expect(page.locator('[data-test=navigation]')).toBeVisible();
  });

  test('Settings button', async ({ page }) => {
    await expect(page.locator('[data-test=settings-button]')).toBeVisible();
  });

  test('Notification button', async ({ page }) => {
    await expect(page.locator('[data-test=notification-button]')).toBeVisible();
  });

  test('Profile button', async ({ page }) => {
    await expect(page.locator('[data-test=profile-button]')).toBeVisible();
  });

  // test('Service Mesh Section', async ({ page }) => {
  //   await expect(page.locator('[data-test=service-mesh]')).toBeVisible();
  // });
});
