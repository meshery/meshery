import { test, expect } from '@playwright/test';
import { ENV } from './env';

test.describe('Service Mesh Performance Management Tests', () => {
  const peformanceProfiles = [
    {
      profileName: 'Sample-Perf-Test',
      serviceMesh: 'None',
      url: 'https://layer5.io/',
      loadGenerator: 'fortio',
    },
  ];

  peformanceProfiles.forEach(({ profileName, serviceMesh, url, loadGenerator }) => {
    test(`Add performace profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByLabel('Add Performance Profile').click();
      await page.getByLabel('Profile Name').fill(profileName);
      await page.getByLabel('Service Mesh').click();
      await page.locator(`[data-value=${serviceMesh}]`).click();
      await page.getByRole('textbox', { name: 'url' }).fill(url);
      await page.getByRole('spinbutton', { name: 'Concurrent requests' }).fill('2');
      await page.getByRole('spinbutton', { name: 'Queries per second' }).fill('2');
      await page.getByRole('textbox', { name: 'Duration' }).fill('15s');
      await page.getByLabel(loadGenerator).check();

      await expect(page.getByRole('button', { name: 'Run Test', exact: true })).toBeVisible();
      await page.getByRole('button', { name: 'Run Test', exact: true }).click();

      await expect(await page.getByText('fetched the data.')).toBeVisible({
        timeout: 2 * 60 * 1000,
      });
    });

    test(`View detailed result of a performance profile (Graph Visualiser) with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.waitForSelector(`text=${profileName}`, {
        state: 'visible',
      });
      await page.getByRole('button', { name: 'View Results', exact: true }).click();

      expect(page.getByText('Sorry, no matching records')).toBeHidden();
      await page.getByLabel('more').first().click();

      await expect(await page.getByText(`URL: ${url}`, { exact: true })).toBeVisible();
    });

    test(`View/Edit the configuration of a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileName).first().click();

      await page.getByTestId('performanceProfieCard-edit').click();
      await page.getByLabel('Service Mesh').click();
      await page.locator(`[data-value=${serviceMesh}]`).click();
      await page.getByRole('spinbutton', { name: 'Concurrent requests' }).fill('3');
      await page.getByLabel(loadGenerator).check();

      await expect(page.getByRole('button', { name: 'Run Test', exact: true })).toBeVisible();
      await page.getByRole('button', { name: 'Run Test', exact: true }).click();

      await expect(page.getByText('fetched the data.')).toBeVisible({ timeout: 2 * 60 * 1000 });
    });

    test(`Delete a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileName, { exact: true }).click();
      await page.getByTestId('performanceProfieCard-delete').first().click();

      await expect(await page.getByText('Performance Profile Deleted!').first()).toBeHidden();
    });
  });
});
