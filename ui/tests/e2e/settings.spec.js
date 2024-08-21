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

  test('Adapters, Metrics, Registry, Reset tabs are displayed', async ({ page }) => {
    // Navigate to "Adapters" tab
    await page.getByRole('tab', { name: 'Adapters' }).click();

    // Mesh Adapter URL and Available Adapters input elements are displayed
    await expect(page.locator('label').filter({ hasText: 'Mesh Adapter URL' })).toBeVisible();
    await expect(page.locator('label').filter({ hasText: 'Available Adapters' })).toBeVisible();

    // Navigate to "Metrics" tab
    await page.getByRole('tab', { name: 'Metrics' }).click();

    // Navigate to "Grafana" sub-tab
    await page.getByRole('tab', { name: 'Grafana' }).click();

    // Grafana Base URL input elements are displayed
    await expect(page.locator('label').filter({ hasText: 'Grafana Base URL' })).toBeVisible();

    // Navigate to "Prometheus" sub-tab
    await page.getByRole('tab', { name: 'Prometheus' }).click();

    // Prometheus Base URL subt-tab are displayed
    await expect(page.locator('label').filter({ hasText: 'Prometheus Base URL' })).toBeVisible();

    // Navigate to "Registry" tab
    await page.getByRole('tab', { name: 'Registry' }).click();

    // "Mesh Model" component is displayed
    await expect(page.locator('[data-test="workloads"]')).toBeVisible();

    /// Navigate to "DataBase summary" tab
    await page.getByRole('tab', { name: 'Reset' }).click();

    // "DataBase summary table" is displayed
    await expect(page.getByRole('button', { name: 'RESET DATABASE' })).toBeVisible();
    await expect(page.getByRole('columnheader', { name: 'Name' })).toBeVisible();
    await expect(page.getByRole('columnheader', { name: 'Count' })).toBeVisible();
  });

  test('Connect to Meshery Istio Adapter and configure it', async ({ page }) => {
    // Navigate to 'Adapters' tab
    await page.getByRole('tab', { name: 'Adapters', exact: true }).click({ force: true });

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
    await meshManageReq;
    const res = await meshManageRes;
    const body = await res.json();
    verifyAdapterResBody(body);

    // Verify success notification
    await expect(page.getByText('Adapter was configured!')).toBeVisible();
  });
});
