import { expect, test } from '@playwright/test';
import { ENV } from '../env';

test.describe('Service Mesh Lifecycle Tests', () => {
  const mesheryAdapters = [
    { adapterName: 'Istio', adapterPort: '10000', deploy: false },
    { adapterName: 'Consul', adapterPort: '10002', deploy: false },
  ];

  mesheryAdapters.forEach(({ adapterName, adapterPort }) => {
    const ADAPTER_LOCATION = `localhost:${adapterPort}`;
    test(`Configure Existing ${adapterName} adapter through Mesh Adapter URL from Management page`, async ({
      page,
    }) => {
      await test.step('When I visit the settings page', async () => {
        await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);
      });

      await test.step('And I click the adapters tab', async () => {
        await page.getByRole('tab', { name: 'Adapters', exact: true }).click({ force: true });
      });

      await test.step('And I enter the Mesh adapter URL', async () => {
        await page
          .locator('label')
          .filter({ hasText: /Mesh Adapter URL/ })
          .locator('..')
          .locator('input')
          .fill(`localhost:${adapterPort}`);
        await page.keyboard.press('Enter');
      });

      await test.step('And I click the connect button', async () => {
        await page.getByRole('button', { name: 'Connect', exact: true }).click();
      });

      await test.step('And I see success notification', async () => {
        await expect(page.getByText('Adapter was configured!')).toBeVisible();
      });

      await test.step('When I visit the service mesh page', async () => {
        await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);
      });

      await test.step('And I Open `Select Meshery Adapter` Dropdown, select the adapter by URL, and verify the selection', async () => {
        const dropdown = page
          .locator('label')
          .filter({ hasText: /Select Meshery Adapter/ })
          .locator('..');
        await dropdown.click();
        await page.getByRole('option', { name: ADAPTER_LOCATION }).click();
        await expect(dropdown).toContainText(ADAPTER_LOCATION);
      });
    });

    test(`Ping ${adapterName} Adapter`, async ({ page }) => {
      await test.step('When I visit the service-mesh page', async () => {
        await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);
      });

      await test.step('And I Open `Select Meshery Adapter` Dropdown, select the adapter by URL, and verify the selection', async () => {
        const dropdown = page
          .locator('label')
          .filter({ hasText: /Select Meshery Adapter/ })
          .locator('..');
        await dropdown.click();
        await page.getByRole('option', { name: ADAPTER_LOCATION }).click();
        await expect(dropdown).toContainText(ADAPTER_LOCATION);
      });

      await test.step('And Ping the adapter by clicking the chip containing adapter URL in `Manage Service Mesh` section', async () => {
        await page.getByRole('button', { name: ADAPTER_LOCATION, exact: true }).click();
      });

      await test.step('And I Verify that the adapter was pinged', async () => {
        await expect(page.getByText('Adapter pinged!')).toBeVisible();
      });
    });
  });
});
