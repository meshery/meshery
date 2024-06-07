import { expect, test } from '@playwright/test';
import { ENV } from './env';

test.describe('Service Mesh Lifecycle Tests', () => {
  const mesheryAdapters = [
    { adapterName: 'Istio', adapterPort: '10000', deploy: false },
    { adapterName: 'Consul', adapterPort: '10002', deploy: false },
  ];

  mesheryAdapters.forEach(({ adapterName, adapterPort, deploy }) => {
    const ADAPTER_LOCATION = `localhost:${adapterPort}`;
    test(`Configure Existing ${adapterName} adapter through Mesh Adapter URL from Management page`, async ({
      page,
    }) => {
      // Visit Settings Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);

      // Navigate to 'Adapters' tab
      await page.locator('[data-cy=tabServiceMeshes]').click({ force: true });

      // Enter Mesh Adapter URL
      await page.locator("[data-cy='mesh-adapter-url']").locator('input').fill('localhost:10000');
      await page.keyboard.press('Enter');

      // Click 'Connect' Button
      await page.locator('[data-cy=btnSubmitMeshAdapter]').click();

      // Verify success notification
      await expect(page.getByText('Adapter was configured!')).toBeVisible();

      // Visit Lifecycle > Service Mesh Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);

      // Open "Select Service Mesh Type" Dropdown
      await page.locator('[data-cy="lifecycle-service-mesh-type"]').click();

      // Select the adapter by URL
      await page.locator('[role="listbox"]').getByText(ADAPTER_LOCATION).click();

      // Verify selection of the adapter by URL
      await expect(page.locator('[data-cy="lifecycle-service-mesh-type"]')).toContainText(
        ADAPTER_LOCATION,
      );
    });

    test(`Ping ${adapterName} Adapter`, async ({ page }) => {
      // Visit Lifecycle > Service Mesh Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);

      // Open "Select Service Mesh Type" Dropdown
      await page.locator('[data-cy="lifecycle-service-mesh-type"]').click();

      // Select the adapter by URL
      await page.locator('[role="listbox"]').getByText(ADAPTER_LOCATION).click();

      // Verify selection of the adapter by URL
      await expect(page.locator('[data-cy="lifecycle-service-mesh-type"]')).toContainText(
        ADAPTER_LOCATION,
      );

      // Ping the adapter by clicking the chip containing adapter URL in "Manage Service Mesh" section
      await page.locator('[data-cy="adapter-chip-ping"]').click();

      // Verify that the adapter was pinged
      await expect(page.getByText('Adapter pinged!')).toBeVisible();
    });
  });
});
