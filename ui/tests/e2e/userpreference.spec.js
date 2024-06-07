import { expect, test } from '@playwright/test';
import { ENV } from './env';

const userPreferenceTests = [
  {
    name: 'Toggle "Meshery Catalog Content"',
    apiURL: `${ENV.MESHERY_SERVER_URL}/api/user/prefs`,
    locator: '[data-cy="CatalogContentPreference"]',
    expectedMethod: 'POST',
    expectedStatus: 200,
  },
  {
    name: 'Toggle "Send Anonymous Usage Statistics"',
    apiURL: `${ENV.MESHERY_SERVER_URL}/api/user/prefs?contexts=all`,
    locator: '[data-cy="UsageStatsPreference"]',
    expectedMethod: 'POST',
    expectedStatus: 200,
  },
  {
    name: 'Toggle "Send Anonymous Performance Results"',
    apiURL: `${ENV.MESHERY_SERVER_URL}/api/user/prefs?contexts=all`,
    locator: '[data-cy="PerfResultPreference"]',
    expectedMethod: 'POST',
    expectedStatus: 200,
  },
];

test.describe('User Preferences Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    const userPrefReq = page.waitForRequest((request) =>
      request.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`),
    );
    const userPrefRes = page.waitForResponse((response) =>
      response.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`),
    );

    // Visit Settings Page
    await page.goto(`${ENV.MESHERY_SERVER_URL}/user/preferences`);

    // Verify requests and responses expected on initial page load
    expect((await userPrefReq).method()).toBe('GET');
    expect((await userPrefRes).status()).toBe(200);

    // Verify visibility of 'Extensions' Section
    await expect(
      page.locator(':nth-child(1) > .MuiFormControl-root > .MuiFormLabel-root'),
    ).toHaveText('Extensions');

    // Verify visibility of 'Analytics and Improvement Program' Section
    await expect(
      page.locator(':nth-child(2) > .MuiFormControl-root > .MuiFormLabel-root'),
    ).toHaveText('Analytics and Improvement Program');
  });

  for (const t of userPreferenceTests) {
    test(t.name, async ({ page }) => {
      const userPrefReq = page.waitForRequest(t.apiURL);
      const userPrefRes = page.waitForResponse(t.apiURL);

      // Check current state of switch (checked or unchecked)
      const prefSwitch = page.locator(t.locator);
      const wasChecked = (await prefSwitch.getAttribute('class')).includes('Mui-checked');

      // Toggle the state of switch
      await prefSwitch.click();

      // Verify requests and responses
      expect((await userPrefReq).method()).toBe(t.expectedMethod);
      expect((await userPrefRes).status()).toBe(t.expectedStatus);

      // Verify that state of switch changed
      await page.waitForTimeout(2000);
      if (wasChecked) await expect(prefSwitch).not.toHaveClass(/Mui-checked/);
      else await expect(prefSwitch).toHaveClass(/Mui-checked/);
    });
  }
});
