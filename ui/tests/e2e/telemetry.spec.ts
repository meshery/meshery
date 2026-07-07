import { expect, test, Page } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';

// Telemetry is the section that replaced the legacy "Configure Metrics" flow
// (Grafana API key under Settings > Metrics, and the Performance page's
// configure-metrics button). It now lives under /telemetry with two sub-pages:
//   - /telemetry/charts  -> Grafana dashboards
//   - /telemetry/metrics -> Prometheus metrics
// Both render an empty state when no matching connection is registered, which
// is the expected state in a fresh CI environment.

test.describe('Telemetry Section Tests', () => {
  // navigateToDashboard() + navigateToSubMenuItem() chain their own 120s inner
  // waits; 180s gives the hook enough wall-clock on slow CI.
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page }: { page: Page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
  });

  test('Charts (Grafana) page loads', async ({ page }: { page: Page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToTelemetryCharts();

    await expect(page).toHaveURL(/telemetry\/charts/i);
    await expect(page.getByTestId('telemetry-charts')).toBeVisible();

    // Either the empty state (no Grafana connection registered, the CI default)
    // or the connection toolbar (an environment with a connection) is fine —
    // both prove the page mounted without crashing.
    await expect(
      page.getByTestId('telemetry-grafana-empty').or(page.getByTestId('telemetry-grafana-toolbar')),
    ).toBeVisible();
  });

  test('Metrics (Prometheus) page loads', async ({ page }: { page: Page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToTelemetryMetrics();

    await expect(page).toHaveURL(/telemetry\/metrics/i);
    await expect(page.getByTestId('telemetry-metrics')).toBeVisible();

    await expect(
      page
        .getByTestId('telemetry-prometheus-empty')
        .or(page.getByTestId('telemetry-prometheus-toolbar')),
    ).toBeVisible();
  });
});
