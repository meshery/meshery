import { expect, test } from '@playwright/test';
import { ENV } from './env';

const SETTINGS_TABS = [
  'settings-tab-adapters',
  'settings-tab-metrics',
  'settings-tab-registry',
  'settings-tab-reset',
];

const MESH_ADAPTER = ['adapters-available-label'];

const ACTION_BUTTONS = [
  'adapter-undeploy-button',
  'adapter-connect-button',
  'adapter-deploy-button',
];

const GARFANA_ELEMENTS = ['grafana-base-url', 'grafana-api-key'];

test.describe('Performance Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(ENV.MESHERY_SERVER_URL);
    const performanceNav = page.locator('[data-cy="Performance"]');
    await performanceNav.click();
  });

  test(
    'Verify Configure Metrics Navigation and Settings',
    // adding for testing
    { tag: '@unstable' },
    async ({ page }) => {
      await page.getByTestId('configure-metrics-button').click();
      await expect(page).toHaveURL(/.*#metrics/);

      for (const tabId of SETTINGS_TABS) {
        await expect(page.getByTestId(tabId)).toBeVisible();
      }

      for (const MeshId of MESH_ADAPTER) {
        await expect(page.getByTestId(MeshId)).toBeVisible();
      }

      for (const buttonId of ACTION_BUTTONS) {
        await expect(page.getByTestId(buttonId)).toBeVisible();
      }

      for (const garfanaId of GARFANA_ELEMENTS) {
        await expect(page.getByTestId(garfanaId)).toBeVisible();
      }
    },
  );
});
