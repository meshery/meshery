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

test.describe('Settings Page Tests', () => {
  test.beforeEach(async ({ page }) => {
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

    // Visit Settings Page
    await page.goto(`${ENV.MESHERY_SERVER_URL}/settings`);

    // Verify requests and responses expected on initial page load
    await meshAdapterReq;
    const res = await meshAdapterRes;
    const body = await res.json();
    verifyAdapterResBody(body);
  });

  test('Connect to Meshery Istio Adapter and configure it', async ({ page }) => {
    let meshManageReq;
    let meshManageRes;

    await test.step('When I click the `Adapters tab`', async () => {
      await page.getByRole('tab', { name: 'Adapters', exact: true }).click({ force: true });
    });

    await test.step('Initiate the request and response', async () => {
      meshManageReq = page.waitForRequest(
        (request) =>
          request.url() === `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage` &&
          request.method() === 'POST',
      );
      meshManageRes = page.waitForResponse(
        (response) =>
          response.url() === `${ENV.MESHERY_SERVER_URL}/api/system/adapter/manage` &&
          response.status() === 200,
      );
    });

    await test.step('And I enter the mesh adapter url', async () => {
      await page
        .locator('label')
        .filter({ hasText: /Mesh Adapter URL/ })
        .locator('..')
        .locator('input')
        .fill('localhost:10000');
      await page.keyboard.press('Enter');
    });

    await test.step('And I click the connect button', async () => {
      await page.getByRole('button', { name: 'Connect', exact: true }).click();
    });

    await test.step('And I verify the requests and reponses', async () => {
      await meshManageReq;
      const res = await meshManageRes;
      const body = await res.json();
      verifyAdapterResBody(body);
    });

    await test.step('And I see success notification', async () => {
      await expect(page.getByText('Adapter was configured!')).toBeVisible();
    });
  });
});
