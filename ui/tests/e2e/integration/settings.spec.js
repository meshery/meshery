import { expect, test } from '@playwright/test';
import { ENV } from '../env';

test.describe('Settings', () => {
  test.describe('Service Meshes', () => {
    test.beforeEach(async ({ page }) => {
      const meshAdapterReq = page.waitForRequest(`${ENV.MESHERY_SERVER_URL}/api/system/adapters`);
      const meshAdapterRes = page.waitForResponse(`${ENV.MESHERY_SERVER_URL}/api/system/adapters`);

      await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);
      expect((await meshAdapterReq).method()).toBe('GET');
      const res = await meshAdapterRes;
      expect(res.status()).toBe(200);

      await page.locator('[data-cy="tabServiceMeshes"]').click({ force: true });
    });

    test('Adapter Connection Status', async ({ page }) => {
      await expect(page.locator("[data-cy='mesh-adapter-connections']")).toBeVisible();
    });

    test('select, submit, and confirm', async ({ page }) => {
      const meshManageReq = page.waitForRequest(
        `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage`,
      );
      const meshManageRes = page.waitForResponse(
        `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage`,
      );

      await page.locator("[data-cy='mesh-adapter-url']").locator('input').fill('localhost:10000');
      await page.keyboard.press('Enter');
      await page.locator('[data-cy=btnSubmitMeshAdapter]').click();

      expect((await meshManageReq).method()).toBe('POST');
      const res = await meshManageRes;
      expect(res.status()).toBe(200);

      await expect(page.getByText('Adapter was configured!')).toBeVisible();
    });
  });
});
