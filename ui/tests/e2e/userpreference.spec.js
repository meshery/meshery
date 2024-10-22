import { expect, test } from '@playwright/test';
import { ENV } from './env';

const userPreferenceTests = [
  // {
  //   name: 'Toggle "Meshery Catalog Content"',
  //   apiURL: `${ENV.MESHERY_SERVER_URL}/api/user/prefs`,
  //   switchLabel: 'Meshery Catalog Content',
  //   expectedMethod: 'POST',
  //   expectedStatus: 200,
  // },
  {
    name: 'Toggle "Send Anonymous Usage Statistics"',
    apiURL: `${ENV.MESHERY_SERVER_URL}/api/user/prefs`,
    switchLabel: 'Send Anonymous Usage Statistics',
    expectedMethod: 'POST',
    expectedStatus: 200,
  },
  {
    name: 'Toggle "Send Anonymous Performance Results"',
    apiURL: `${ENV.MESHERY_SERVER_URL}/api/user/prefs`,
    switchLabel: 'Send Anonymous Performance Results',
    expectedMethod: 'POST',
    expectedStatus: 200,
  },
];

test.describe(
  'User Preferences Page Tests',
  {
    tag: '@unstable',
    annotation: [{ type: 'issue', description: 'https://github.com/meshery/meshery/issues/12329' }],
  },
  () => {
    test.beforeEach(async ({ page }) => {
      const userPrefReq = page.waitForRequest(
        (request) =>
          request.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`) &&
          request.method() === 'GET',
      );
      const userPrefRes = page.waitForResponse(
        (response) =>
          response.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`) &&
          response.status() === 200,
      );

      // Visit User Preferences Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/user/preferences`);

      // Verify requests and responses expected on initial page load
      await userPrefReq;
      await userPrefRes;

      // Verify visibility of 'Extensions' Section
      await expect(page.getByRole('group', { name: /Extensions.*/ })).toBeVisible();

      // Verify visibility of 'Analytics and Improvement Program' Section
      await expect(
        page.getByRole('group', { name: /Analytics and Improvement Program.*/ }),
      ).toBeVisible();

      // Verify visibility of 'Theme' Section
      await expect(page.getByRole('group', { name: /Theme.*/ })).toBeVisible();
    });

    for (const t of userPreferenceTests) {
      test(t.name, async ({ page }) => {
        const userPrefReq = page.waitForRequest(
          (request) => request.url() === t.apiURL && request.method() === t.expectedMethod,
        );
        const userPrefRes = page.waitForResponse(
          (response) => response.url() === t.apiURL && response.status() === t.expectedStatus,
        );

        // Check current state of switch (checked or unchecked)
        const prefSwitch = page.getByLabel(t.switchLabel);
        const wasChecked = await prefSwitch.isChecked();

        // Toggle the state of switch
        await prefSwitch.click();

        // Verify requests and responses
        await userPrefReq;
        await userPrefRes;

        // Verify that state of switch changed
        await page.waitForTimeout(2000);
        if (wasChecked) await expect(prefSwitch).not.toBeChecked();
        else await expect(prefSwitch).toBeChecked();
      });
    }
  },
);
