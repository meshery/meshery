import { expect, test } from '@playwright/test';
import { ENV } from './env';

const URLS = {
  MESHERY: {
    SETTINGS: `${ENV.MESHERY_SERVER_URL}/settings`,
  },
};

test.describe('Performance Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(ENV.MESHERY_SERVER_URL);
    const performanceNav = page.locator('[data-cy="Performance"]');
    await performanceNav.click();
  });

  test.unstable()('Verify Configure Metrics Navigation and Settings', async ({ page }) => {
    await expect(page.getByTestId('configure-metrics-button')).toBeVisible();
    await page.getByTestId('configure-metrics-button').click();

    await expect(page).toHaveURL(`${URLS.MESHERY.SETTINGS}#metrics`);

    // settings tabs
    await expect(page.getByTestId('settings-tab-adapters')).toBeVisible();
    await expect(page.getByTestId('settings-tab-metrics')).toBeVisible();
    await expect(page.getByTestId('settings-tab-registry')).toBeVisible();
    await expect(page.getByTestId('settings-tab-reset')).toBeVisible();

    // Mesh Adapter section
    await expect(page.getByTestId('adapters-available-label')).toBeVisible();

    // action buttons
    await expect(page.getByTestId('adapter-undeploy-button')).toBeVisible();
    await expect(page.getByTestId('adapter-connect-button')).toBeVisible();
    await expect(page.getByTestId('adapter-deploy-button')).toBeVisible();

    await expect(page.getByTestId('database-reset-button')).toBeVisible();

    // Grafana configuration elements
    await expect(page.getByTestId('grafana-base-url')).toBeVisible();
    await expect(page.getByTestId('grafana-api-key')).toBeVisible();
  });
});
