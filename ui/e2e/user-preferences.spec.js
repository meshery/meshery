// @ts-check
const { test, expect } = require('@playwright/test');

test.describe('User Preferences', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/user/preferences');
  });

  test('Extensions | Section is visible', async ({ page }) => {
    // Check if the Extensions section is visible
    await expect(page.getByRole('group', { name: 'Extensions' }).locator('legend')).toBeVisible();
  });

  test('Extensions | Deactivates Meshery Catalog Content', async ({ page }) => {
    let responsePromise = page.waitForResponse(response => response.url().includes('/api/user/prefs') && response.request().method() === 'POST');

    // Deactivate Meshery Catalog Content
    await page.getByLabel('Meshery Catalog Content').click();
    let response = await responsePromise;
    let responseBody = await response.json();

    expect(response.status()).toBe(200);
    expect(responseBody.usersExtensionPreferences.catalogContent).toBe(false);

    // Reactivate Meshery Catalog Content
    await page.getByLabel('Meshery Catalog Content').click();
  });

  test('Analytics and Improvement Program', async ({ page }) => {
    await expect(page.getByText('Analytics and Improvement Program')).toBeVisible();
  });
});
