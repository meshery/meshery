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
    await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);
  });

  test('Aggregation Charts are displayed', async ({ page }) => {
    await expect(
      page.getByRole('heading', {
        name: 'Models by Category',
      }),
    ).toBeVisible();

    await expect(
      page.getByRole('heading', {
        name: 'Registry',
      }),
    ).toBeVisible();

    await expect(
      page.getByRole('heading', {
        name: 'Connections',
        exact: true,
      }),
    ).toBeVisible();

    await expect(
      page.getByRole('heading', {
        name: 'Configuration',
      }),
    ).toBeVisible();
  });

  test('Connect to Meshery Istio Adapter and configure it', async ({ page }) => {
    // Navigate to 'Adapters' tab
    await page.getByRole('tab', { name: 'Adapters', exact: true }).click({ force: true });

    const meshAdapterReq = page.waitForRequest(
      (request) =>
        request.url() === `${ENV.MESHERY_SERVER_URL}/api/system/adapters` &&
        request.method() === 'GET',
    );
    const meshAdapterRes = page.waitForResponse(
      (response) =>
        response.url() === `${ENV.MESHERY_SERVER_URL}/api/system/adapters` &&
        response.status() === 200,
    );

    const meshManageReq = page.waitForRequest(
      (request) =>
        request.url() === `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage` &&
        request.method() === 'POST',
    );
    const meshManageRes = page.waitForResponse(
      (response) =>
        response.url() === `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage` &&
        response.status() === 200,
    );

    // Enter Mesh Adapter URL
    await page
      .locator('label')
      .filter({ hasText: /Mesh Adapter URL/ })
      .locator('..')
      .locator('input')
      .fill('localhost:10000');
    await page.keyboard.press('Enter');

    // Click 'Connect' Button
    await page.getByRole('button', { name: 'Connect', exact: true }).click();

    // Verify requests and responses
    await meshAdapterReq;
    const adapterRes = await meshAdapterRes;
    verifyAdapterResBody(await adapterRes.json());

    await meshManageReq;
    const manageRes = await meshManageRes;
    verifyAdapterResBody(await manageRes.json());

    // Verify success notification
    await expect(page.getByText('Adapter was configured!')).toBeVisible();
  });
});
