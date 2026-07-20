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
