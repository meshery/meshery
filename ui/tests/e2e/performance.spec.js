import { expect, test } from '@playwright/test';
import { ENV } from './env';

const SETTINGS_TABS = [
  'settings-tab-adapters',
  'settings-tab-metrics',
  'settings-tab-registry',
  'settings-tab-reset',
];

const ACTION_BUTTONS = [
  'adapter-undeploy-button',
  'adapter-connect-button',
  'adapter-deploy-button',
];

const GRAFANA_ELEMENTS = ['grafana-api-key'];

const COMMON_UI_ELEMENTS = [
  'navigation',
  'notification-button',
  'profile-button',
  'header-menu',
];

test.describe('Performance Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(ENV.MESHERY_SERVER_URL);
    const performanceNav = page.locator('[data-cy="performance"]');
    await performanceNav.click();
  });

  test('Common UI elements', async ({ page }) => {
    for (const elementId of COMMON_UI_ELEMENTS) {
      await expect(
        page.getByTestId(elementId),
        `UI element with ID ${elementId} should be visible`,
      ).toBeVisible();
    }
  });

  test.describe('Configure Metrics Navigation and Settings', () => {
    test.beforeEach(async ({ page }) => {
      await page.getByTestId('configure-metrics-button').click();
      await expect(page).toHaveURL(/metrics/i);
    });

    test('All settings tabs', async ({ page }) => {
      for (const tabId of SETTINGS_TABS) {
        await expect(
          page.getByTestId(tabId),
          `Tab with ID ${tabId} should be visible`,
        ).toBeVisible();
      }
    });

    test('Action buttons on adapters tab', async ({ page }) => {
      await page.getByTestId('settings-tab-adapters').click();

      for (const buttonId of ACTION_BUTTONS) {
        await expect(
          page.getByTestId(buttonId),
          `Button with ID ${buttonId} should be visible`,
        ).toBeVisible();
      }
    });

    test('Grafana elements on metrics tab', async ({ page }) => {
      await page.getByTestId('settings-tab-metrics').click();
      for (const grafanaId of GRAFANA_ELEMENTS) {
        await expect(
          page.getByTestId(grafanaId),
          `Grafana element with ID ${grafanaId} should be visible`,
        ).toBeVisible();
      }
    });
  });
});
