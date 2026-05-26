import { expect, test, Page } from '@playwright/test';

const SETTINGS_TABS: string[] = [
  'settings-tab-adapters',
  'settings-tab-metrics',
  'settings-tab-registry',
  'settings-tab-reset',
];

const ACTION_BUTTONS: string[] = [
  'adapter-undeploy-button',
  'adapter-connect-button',
  'adapter-deploy-button',
];

const GRAFANA_ELEMENTS: string[] = ['grafana-api-key'];

const COMMON_UI_ELEMENTS: string[] = [
  'navigation',
  'notification-button',
  'profile-button',
  'header-menu',
];

test.describe('Performance Section Tests', () => {
  test.describe.configure({ timeout: 120000 });

  test.beforeEach(async ({ page }: { page: Page }) => {
    await page.goto('/performance', { waitUntil: 'domcontentloaded' });
    await expect(page).toHaveURL(/\/performance/i, { timeout: 120000 });
    await expect(page.getByTestId('configure-metrics-button')).toBeVisible({ timeout: 120000 });
  });

  test('Common UI elements', async ({ page }: { page: Page }) => {
    for (const elementId of COMMON_UI_ELEMENTS) {
      await expect(
        page.getByTestId(elementId),
        `UI element with ID ${elementId} should be visible`,
      ).toBeVisible();
    }
  });

  test.describe('Configure Metrics Navigation and Settings', () => {
    test.beforeEach(async ({ page }: { page: Page }) => {
      await expect(page.getByTestId('configure-metrics-button')).toBeVisible();
      await page.getByTestId('configure-metrics-button').click();
      await expect(page).toHaveURL(/metrics/i);
    });

    test('All settings tabs', async ({ page }: { page: Page }) => {
      for (const tabId of SETTINGS_TABS) {
        await expect(
          page.getByTestId(tabId),
          `Tab with ID ${tabId} should be visible`,
        ).toBeVisible();
      }
    });

    test('Action buttons on adapters tab', async ({ page }: { page: Page }) => {
      await page.getByTestId('settings-tab-adapters').click();

      for (const buttonId of ACTION_BUTTONS) {
        await expect(
          page.getByTestId(buttonId),
          `Button with ID ${buttonId} should be visible`,
        ).toBeVisible();
      }
    });

    test('Grafana elements on metrics tab', async ({ page }: { page: Page }) => {
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
