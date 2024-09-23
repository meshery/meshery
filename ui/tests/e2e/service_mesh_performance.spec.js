import { test as base, expect } from '@playwright/test';
import { PerformancePage } from './fixtures/performancePage';
import { ENV } from './env';

export const test = base.extend({
  performancePage: async ({ page }, use) => {
    const perfPage = new PerformancePage(page);
    await perfPage.navigate();
    await use(perfPage);
  },
});

// Disable this test until got fixed
test.describe('Service Mesh Performance Management Tests', () => {
  test.skip();

  const profileName = 'Sample-test';

  test.beforeEach(async ({ performancePage }) => {
    await performancePage.goToProfiles();
  });

  test('Run a performance test through profile', async ({ performancePage }) => {
    await performancePage.createPerformanceProfile(profileName);
    await expect(performancePage.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
  });

  test('View detailed result of a performance profile (Graph Visualiser)', async ({
    performancePage,
  }) => {
    await performancePage.goToProfiles();
    await performancePage.viewPerformanceProfileResult(profileName);
  });

  test('Run a performance test', async ({ performancePage }) => {
    await performancePage.runPerformanceTest(profileName);
  });

  test('View Results from a performance profile', async ({ performancePage }) => {
    await performancePage.goToProfiles();
    await performancePage.viewPerformanceProfileResult(profileName);
  });

  test('View/Edit the configuration of a performance profile', async ({ performancePage }) => {
    await performancePage.goToProfiles();
    await performancePage.viewPerformanceProfileConfiguration(profileName);
    await performancePage.page.locator('[aria-labelledby="meshName-label meshName"]').click();
    await performancePage.page.locator('[data-value="istio"]').click();
    await performancePage.fillInput(
      performancePage.page.getByRole('spinbutton', { name: 'Concurrent requests' }),
      '6',
    );
    await expect(
      performancePage.page.getByRole('button', { name: 'Run Test', exact: true }),
    ).toBeVisible();
    await performancePage.page.getByRole('button', { name: 'Run Test', exact: true }).click();
    const notification = await performancePage.page
      .locator('text=Initiating load test . . .')
      .first();
    await expect(notification).toBeVisible();
  });
});
