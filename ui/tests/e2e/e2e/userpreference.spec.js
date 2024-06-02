import { expect, test } from '@playwright/test';
import { ENV } from '../env';

test.describe('User Preferences', () => {
  test.describe('Extensions | Analytics and Improvement Program', () => {
    test.beforeEach(async ({ page }) => {
      const userPrefReq = page.waitForRequest((request) =>
        request.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`),
      );
      const userPrefRes = page.waitForResponse((response) =>
        response.url().startsWith(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`),
      );

      await page.goto(`${ENV.MESHERY_SERVER_URL}/user/preferences`);
      expect((await userPrefReq).method()).toBe('GET');
      expect((await userPrefRes).status()).toBe(200);

      await expect(
        page.locator(':nth-child(1) > .MuiFormControl-root > .MuiFormLabel-root'),
      ).toHaveText('Extensions');
      await expect(
        page.locator(':nth-child(2) > .MuiFormControl-root > .MuiFormLabel-root'),
      ).toHaveText('Analytics and Improvement Program');
    });

    test.describe('Extensions', () => {
      test('toggle Meshery Catalog Content', async ({ page }) => {
        const userPrefReq = page.waitForRequest(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`);
        const userPrefRes = page.waitForResponse(`${ENV.MESHERY_SERVER_URL}/api/user/prefs`);

        const catalogContentSwitch = page.locator('[data-cy="CatalogContentPreference"]');
        const wasChecked = (await catalogContentSwitch.getAttribute('class')).includes(
          'Mui-checked',
        );
        await catalogContentSwitch.click();
        expect((await userPrefReq).method()).toBe('POST');
        expect((await userPrefRes).status()).toBe(200);

        await page.waitForTimeout(2000);
        if (wasChecked) await expect(catalogContentSwitch).not.toHaveClass(/Mui-checked/);
        else await expect(catalogContentSwitch).toHaveClass(/Mui-checked/);
      });
    });

    test.describe('Analytics and Improvement Program', () => {
      test.describe.configure({ retries: 3 });

      test('toggle "Send Anonymous Usage Statistics"', async ({ page }) => {
        const userPrefReq = page.waitForRequest(
          `${ENV.MESHERY_SERVER_URL}/api/user/prefs?contexts=all`,
        );
        const userPrefRes = page.waitForResponse(
          `${ENV.MESHERY_SERVER_URL}/api/user/prefs?contexts=all`,
        );

        const usageStatisticsSwitch = page.locator('[data-cy="UsageStatsPreference"]');
        const wasChecked = (await usageStatisticsSwitch.getAttribute('class')).includes(
          'Mui-checked',
        );
        await usageStatisticsSwitch.click();
        expect((await userPrefReq).method()).toBe('POST');
        expect((await userPrefRes).status()).toBe(200);

        await page.waitForTimeout(2000);
        if (wasChecked) await expect(usageStatisticsSwitch).not.toHaveClass(/Mui-checked/);
        else await expect(usageStatisticsSwitch).toHaveClass(/Mui-checked/);
      });

      test('toggle "Send Anonymous Performance Results"', async ({ page }) => {
        const userPrefReq = page.waitForRequest(
          `${ENV.MESHERY_SERVER_URL}/api/user/prefs?contexts=all`,
        );
        const userPrefRes = page.waitForResponse(
          `${ENV.MESHERY_SERVER_URL}/api/user/prefs?contexts=all`,
        );

        const performanceResultsSwitch = page.locator('[data-cy="PerfResultPreference"]');
        const wasChecked = (await performanceResultsSwitch.getAttribute('class')).includes(
          'Mui-checked',
        );
        await performanceResultsSwitch.click();
        expect((await userPrefReq).method()).toBe('POST');
        expect((await userPrefRes).status()).toBe(200);

        await page.waitForTimeout(2000);
        if (wasChecked) await expect(performanceResultsSwitch).not.toHaveClass(/Mui-checked/);
        else await expect(performanceResultsSwitch).toHaveClass(/Mui-checked/);
      });
    });
  });
});
