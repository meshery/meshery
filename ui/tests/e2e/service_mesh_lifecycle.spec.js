import { expect, test } from '@playwright/test';
import fs from 'fs';
import { ENV } from './env';

let adapterStatus = {}; // Global object to store adapter statuses

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
      // Visit Settings Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);

      // Navigate to 'Adapters' tab
      await page.getByRole('tab', { name: 'Adapters', exact: true }).click({ force: true });

      // Enter Mesh Adapter URL
      await page
        .locator('label')
        .filter({ hasText: /Mesh Adapter URL/ })
        .locator('..')
        .locator('input')
        .fill(`localhost:${adapterPort}`);
      await page.keyboard.press('Enter');

      // Click 'Connect' Button
      await page.getByRole('button', { name: 'Connect', exact: true }).click();

      // Verify success notification
      try {
        await expect(page.getByText('Adapter was configured!')).toBeVisible();
        adapterStatus[adapterName] = true; // Set status to true
      } catch (error) {
        adapterStatus[adapterName] = false; // Set status to false
        return;
      }

      // Visit Lifecycle > Service Mesh Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);

      // Open "Select Meshery Adapter" Dropdown
      const dropdown = page
        .locator('label')
        .filter({ hasText: /Select Meshery Adapter/ })
        .locator('..');
      await dropdown.click();

      // Select the adapter by URL
      await page.getByRole('option', { name: ADAPTER_LOCATION }).click();

      // Verify selection of the adapter by URL
      try {
        await expect(dropdown).toContainText(ADAPTER_LOCATION);
      } catch (error) {
        adapterStatus[adapterName] = false; // Set status to false
        return;
      }
    });

    test(`Ping ${adapterName} Adapter`, async ({ page }) => {
      // Visit Lifecycle > Service Mesh Page
      await page.goto(`${ENV.MESHERY_SERVER_URL}/management/service-mesh`);

      // Open "Select Meshery Adapter" Dropdown
      const dropdown = page
        .locator('label')
        .filter({ hasText: /Select Meshery Adapter/ })
        .locator('..');
      await dropdown.click();

      // Select the adapter by URL
      await page.getByRole('option', { name: ADAPTER_LOCATION }).click();

      // Verify selection of the adapter by URL
      await expect(dropdown).toContainText(ADAPTER_LOCATION);

      // Ping the adapter by clicking the chip containing adapter URL in "Manage Service Mesh" section
      await page.getByRole('button', { name: ADAPTER_LOCATION, exact: true }).click();

      // Verify that the adapter was pinged
      try {
        await expect(page.getByText('Adapter pinged!')).toBeVisible();
      } catch (error) {
        adapterStatus[adapterName] = false; // Set status to false
        return;
      }
    });
  });

  // Save adapterStatus to a JSON file after all tests
  test.afterAll(async () => {
    fs.writeFileSync('./adapterStatus.json', JSON.stringify(adapterStatus));
  });
});
