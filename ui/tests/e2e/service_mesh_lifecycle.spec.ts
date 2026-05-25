import { expect, test, Page } from '@playwright/test';
import { ENV } from './env';
import { DashboardPage } from './pages/DashboardPage';

interface MesheryAdapter {
  adapterName: string;
  adapterPort: string;
  deploy: boolean;
}

test.describe.skip('Service Mesh Lifecycle Tests', { tag: '@unstable' }, () => {
  const mesheryAdapters: MesheryAdapter[] = [
    { adapterName: 'Istio', adapterPort: '10000', deploy: false },
  ];

  mesheryAdapters.forEach(({ adapterName, adapterPort }: MesheryAdapter) => {
    const ADAPTER_LOCATION = `localhost:${adapterPort}`;

    test(`Configure Existing ${adapterName} adapter through Mesh Adapter URL from Management page`, async ({
      page,
    }: {
      page: Page;
    }) => {
      const dashboardPage = new DashboardPage(page);
      await dashboardPage.navigateToDashboard();
      await dashboardPage.navigateToSettings();

      await page.getByRole('tab', { name: 'Adapters', exact: true }).click({ force: true });

      await page
        .locator('label')
        .filter({ hasText: /Mesh Adapter URL/ })
        .locator('..')
        .locator('input')
        .fill(`localhost:${adapterPort}`);
      await page.keyboard.press('Enter');

      await page.getByRole('button', { name: 'Connect', exact: true }).click();

      await expect(page.getByText('Adapter configured')).toBeVisible();

      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/adapter`);

      const dropdown = page
        .locator('label')
        .filter({ hasText: /Select Meshery Adapter/ })
        .locator('..');
      await dropdown.click();

      await page.getByRole('option', { name: ADAPTER_LOCATION }).click();

      await expect(dropdown).toContainText(ADAPTER_LOCATION);
    });

    test(`Ping ${adapterName} Adapter`, async ({ page }: { page: Page }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/adapter`);

      const dropdown = page
        .locator('label')
        .filter({ hasText: /Select Meshery Adapter/ })
        .locator('..');
      await dropdown.click();

      await page.getByRole('option', { name: ADAPTER_LOCATION }).click();

      await expect(dropdown).toContainText(ADAPTER_LOCATION);

      await page.getByRole('button', { name: ADAPTER_LOCATION, exact: true }).click();

      await expect(page.getByText('Adapter pinged!')).toBeVisible();
    });
  });
});
