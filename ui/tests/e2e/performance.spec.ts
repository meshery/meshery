import { expect, test, Page } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';

test.describe('Performance Section Tests', () => {
  // Generous budget for the dashboard -> performance navigation chain on slow CI.
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page }: { page: Page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToPerformance();
    // Readiness signal for the performance dashboard.
    //
    // Known flake (meshery/meshery#20504): "performance-dashboard" only renders
    // when CAN(VIEW_PERFORMANCE_PROFILES) is true, and CAN() (ui/utils/can.ts)
    // reads a non-reactive module-level casl ability singleton. If the user's
    // capabilities load after Dashboard mounts, it renders <DefaultError/> and
    // never re-renders, so this element can time out regardless of how long we
    // wait - a longer timeout does not fix it. The real fix is a reactive
    // permission gate, tracked in #20504.
    await expect(page.getByTestId('performance-dashboard')).toBeVisible();
  });

  // Global chrome (navigation, notification, profile, header) is covered by the
  // stable indexui.spec.ts tests and by navigateToDashboard() above, so this
  // suite only asserts performance-specific controls.
  test('Performance dashboard controls', async ({ page }: { page: Page }) => {
    await expect(page.getByRole('button', { name: 'Run Test' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Manage Profiles' })).toBeVisible();
  });
});
