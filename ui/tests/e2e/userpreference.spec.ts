import { expect, test } from '@playwright/test';
import { ENV } from './env';
import { DashboardPage } from './pages/DashboardPage';

const userPreferenceTests: {
  name: string;
  apiURL: string;
  switchLabel: string;
  expectedMethod: string;
  expectedStatus: number;
}[] = [
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

test.describe('User Preferences Page Tests', () => {
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

    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToPreferences();

    await userPrefReq;
    await userPrefRes;

    await expect(page.getByRole('group', { name: /Extensions.*/ })).toBeVisible();

    await expect(
      page.getByRole('group', { name: /Analytics and Improvement Program.*/ }),
    ).toBeVisible();

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

      const prefSwitch = page.getByLabel(t.switchLabel);
      const wasChecked = await prefSwitch.isChecked();

      await prefSwitch.click();

      await userPrefReq;
      await userPrefRes;

      await page.waitForTimeout(2000);
      if (wasChecked) {
        await expect(prefSwitch).not.toBeChecked();
      } else {
        await expect(prefSwitch).toBeChecked();
      }
    });
  }
});
