import { expect, test } from '@playwright/test';
import { ENV } from '../env';

test.describe('Lifecycle Service Mesh', () => {
  const selectServiceMeshType = async (page, adapterLocation) => {
    await page.locator('[data-cy="lifecycle-service-mesh-type"]').click();
    await page.locator('[role="listbox"]').getByText(adapterLocation).click();
  };

  const mesheryAdapters = [
    { adapterName: 'Istio', adapterPort: '10000', deploy: false },
    { adapterName: 'Consul', adapterPort: '10002', deploy: false },
  ];

  mesheryAdapters.forEach(({ adapterName, adapterPort, deploy }) => {
    const ADAPTER_LOCATION = `localhost:${adapterPort}`;
    test(`User can Configure Existing ${adapterName} adapter through Mesh Adapter URL from Management page`, async ({
      page,
    }) => {
      // Settings > Adapters Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/settings#service-mesh`);

      // "Mesh Adapter URL" Dropdown
      await page
        .locator('[data-cy="mesh-adapter-url"]')
        .locator('input')
        .fill(`${ADAPTER_LOCATION}`);
      await page.keyboard.press('Enter');
      // "Connect" Button
      await page.locator('[data-cy="btnSubmitMeshAdapter"]').click();
      await expect(page.getByText('Adapter was configured!')).toBeVisible();

      // Lifecycle > Service Mesh Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);
      await selectServiceMeshType(page, ADAPTER_LOCATION);

      // "Select Service Mesh Type" Dropdown
      await expect(page.locator('[data-cy="lifecycle-service-mesh-type"]')).toContainText(
        ADAPTER_LOCATION,
      );
    });

    test(`User can ping ${adapterName} Adapter`, async ({ page }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/`);

      // Lifecycle > Service Mesh Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);
      await selectServiceMeshType(page, ADAPTER_LOCATION);
      await expect(page.locator('[data-cy="lifecycle-service-mesh-type"]')).toContainText(
        ADAPTER_LOCATION,
      );

      // "Manage Service Mesh" Card's Ping adapter chip
      await page.locator('[data-cy="adapter-chip-ping"]').click();
      await expect(page.getByText('Adapter pinged!')).toBeVisible();
    });
  });
});
