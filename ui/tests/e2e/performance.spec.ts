import { expect, test, Page } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';

const COMMON_UI_ELEMENTS: string[] = [
  'navigation',
  'notification-button',
  'profile-button',
  'header-menu',
];

test.describe('Performance Section Tests', () => {
  // The shared beforeEach calls navigateToDashboard() (two 120s visibility
  // waits) and navigateToPerformance() (120s wait), then waits on the
  // performance dashboard. Under the default BASE_TIMEOUT=60s the hook dies
  // before those inner waits resolve when CI is slow. 180s gives the chain
  // enough wall-clock.
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page }: { page: Page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToPerformance();
    // Readiness signal for the performance dashboard. The legacy "Configure
    // Metrics" flow (meshery-metrics / grafana config) was removed when
    // telemetry moved to its own section under /telemetry — see telemetry.spec.ts.
    await expect(page.getByTestId('performance-dashboard')).toBeVisible();
  });

  test('Common UI elements', async ({ page }: { page: Page }) => {
    for (const elementId of COMMON_UI_ELEMENTS) {
      await expect(
        page.getByTestId(elementId),
        `UI element with ID ${elementId} should be visible`,
      ).toBeVisible();
    }
  });

  test('Performance dashboard controls', async ({ page }: { page: Page }) => {
    await expect(page.getByRole('button', { name: 'Run Test' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Manage Profiles' })).toBeVisible();
  });
});
