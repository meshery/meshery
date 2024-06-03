import { expect, test } from '@playwright/test';
import { ENV } from '../env';

const verifyAdapterResBody = (body) => {
  expect(body).toBeInstanceOf(Array);
  body.forEach(({ adapter_location, name, version, git_commit_sha, ops }) => {
    expect(adapter_location).toMatch(/localhost:\d+/);
    expect(name).toEqual(expect.any(String));
    expect(version).toEqual(expect.any(String));
    expect(git_commit_sha).toEqual(expect.any(String));
    // ops can be null or array
    if (ops) {
      expect(ops).toBeInstanceOf(Array);
      ops.forEach(({ key, value, category }) => {
        expect(key).toEqual(expect.any(String));
        expect(value).toEqual(expect.any(String));
        // category may or may not be present, if present its an integer
        if (category) {
          expect(category).toEqual(expect.any(Number));
        } else {
          expect(category).toBeUndefined();
        }
      });
    } else {
      expect(ops).toBeNull();
    }
  });
};

test.describe('Settings', () => {
  test.describe('Service Meshes', () => {
    test.beforeEach(async ({ page }) => {
      const meshAdapterReq = page.waitForRequest(`${ENV.MESHERY_SERVER_URL}/api/system/adapters`);
      const meshAdapterRes = page.waitForResponse(`${ENV.MESHERY_SERVER_URL}/api/system/adapters`);

      await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);
      expect((await meshAdapterReq).method()).toBe('GET');
      const res = await meshAdapterRes;
      expect(res.status()).toBe(200);
      const body = await res.json();
      verifyAdapterResBody(body);

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
      const body = await res.json();
      verifyAdapterResBody(body);

      await expect(page.getByText('Adapter was configured!')).toBeVisible();
    });
  });
});
