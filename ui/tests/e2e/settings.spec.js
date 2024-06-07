import { expect, test } from '@playwright/test';
import { ENV } from './env';

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

test.describe('Settings Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    const meshAdapterReq = page.waitForRequest(`${ENV.MESHERY_SERVER_URL}/api/system/adapters`);
    const meshAdapterRes = page.waitForResponse(`${ENV.MESHERY_SERVER_URL}/api/system/adapters`);

    // Visit Settings Page
    await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);

    // Verify requests and responses expected on initial page load
    expect((await meshAdapterReq).method()).toBe('GET');
    const res = await meshAdapterRes;
    expect(res.status()).toBe(200);
    const body = await res.json();
    verifyAdapterResBody(body);
  });

  test('Connect to Meshery Istio Adapter and configure it', async ({ page }) => {
    // Navigate to 'Adapters' tab
    await page.locator('[data-cy=tabServiceMeshes]').click({ force: true });

    // Verify visibility of 'Adapters' tab
    await expect(page.locator("[data-cy='mesh-adapter-connections']")).toBeVisible();

    const meshManageReq = page.waitForRequest(
      `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage`,
    );
    const meshManageRes = page.waitForResponse(
      `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage`,
    );

    // Enter Mesh Adapter URL
    await page.locator("[data-cy='mesh-adapter-url']").locator('input').fill('localhost:10000');
    await page.keyboard.press('Enter');

    // Click 'Connect' Button
    await page.locator('[data-cy=btnSubmitMeshAdapter]').click();

    // Verify requests and responses
    expect((await meshManageReq).method()).toBe('POST');
    const res = await meshManageRes;
    expect(res.status()).toBe(200);
    const body = await res.json();
    verifyAdapterResBody(body);

    // Verify success notification
    await expect(page.getByText('Adapter was configured!')).toBeVisible();
  });
});
