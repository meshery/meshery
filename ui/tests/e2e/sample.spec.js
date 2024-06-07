import { expect, test } from '@playwright/test';
import { ENV } from './env';

test.describe('Sample Tests', () => {
  test('Visits meshery', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}`);
    expect(page.url()).toBe(`${ENV.MESHERY_SERVER_URL}/`);
  });

  test('Visits meshery settings page', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);
    expect(page.url()).toBe(`${ENV.MESHERY_SERVER_URL}/settings`);
  });
});
