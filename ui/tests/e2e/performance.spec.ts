import { expect, test, Page } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';

/**
 * Waits for the performance dashboard to be visible, retrying via page reload
 * if the non-reactive CAN() permission gate renders <DefaultError/> before
 * capabilities have loaded.
 *
 * Known flake (meshery/meshery#20504): "performance-dashboard" only renders
 * when CAN(VIEW_PERFORMANCE_PROFILES) is true, and CAN() (ui/utils/can.ts)
 * reads a non-reactive module-level casl ability singleton. If the user's
 * capabilities load after the component mounts, it renders <DefaultError/>
 * permanently. A reload after capabilities have settled fixes it.
 */
async function waitForPerformanceDashboard(page: Page, maxRetries = 3): Promise<void> {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    const dashboard = page.getByTestId('performance-dashboard');
    try {
      await dashboard.waitFor({ state: 'visible', timeout: 30_000 });
      return; // element appeared — done
    } catch {
      if (attempt === maxRetries) {
        // Final attempt failed — surface the real assertion error.
        await expect(dashboard).toBeVisible();
      }
      // Capabilities may have settled by now; reload gives the component a
      // fresh mount where CAN() should return true.
      await page.reload({ waitUntil: 'domcontentloaded' });
    }
  }
}

test.describe('Performance Section Tests', () => {
  // Generous budget for the dashboard -> performance navigation chain on slow CI.
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page }: { page: Page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToPerformance();
    // Wait for the performance dashboard with retry-on-reload to handle the
    // CAN() race condition (see meshery/meshery#20504).
    await waitForPerformanceDashboard(page);
  });

  // Global chrome (navigation, notification, profile, header) is covered by the
  // stable indexui.spec.ts tests and by navigateToDashboard() above, so this
  // suite only asserts performance-specific controls.
  test('Performance dashboard controls', async ({ page }: { page: Page }) => {
    await expect(page.getByRole('button', { name: 'Run Test' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Manage Profiles' })).toBeVisible();
  });
});
