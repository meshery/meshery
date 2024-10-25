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
      await page.getByRole('option', { name: serviceMesh }).click();
      await page.getByLabel('URL to test').fill(url);
      await page.getByLabel('Concurrent requests').fill('2');
      await page.getByLabel('Queries per second').fill('2');
      await page.getByLabel('Duration').fill('15s');
      await page.getByLabel(loadGenerator).check();

      const runPerformanceTest = await page.getByTestId('run-performance-test');
      await expect(runPerformanceTest).toBeVisible();
      await runPerformanceTest.click();
      await expect(page.getByText('Load test has been submitted')).toBeVisible();
    });

    test(`View detailed result of a performance profile (Graph Visualiser) with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.waitForSelector(`text=${profileName}`, {
        state: 'visible',
      });
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();

      await page.getByTestId('open-performance-result-bar-chart').first().click();
      await expect(page.locator('[data-test="sentinelStart"]')).toBeHidden();
      await expect(await page.getByText(`URL: ${url}`, { exact: true })).toBeVisible();
    });

    test(`Edit the configuration of a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileName, { exact: true }).first().click();

      await page.getByTestId('performanceProfileCard-edit').click();
      await page.getByLabel('Service Mesh').click();
      await page.getByRole('option', { name: serviceMesh }).click();
      await page.getByLabel('Concurrent requests').fill('3');
      await page.getByLabel(loadGenerator).check();

      const runPerformanceTest = await page.getByTestId('run-performance-test');
      await expect(runPerformanceTest).toBeEnabled();
      await runPerformanceTest.click();
      await expect(page.getByText('Load test has been submitted')).toBeVisible();
    });

    test(`Compare test of a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileName, { exact: true }).first().click();
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.locator('[data-testid="MUIDataTableBodyRow-0"] input[type="checkbox"]').check();

      await page.locator('[title="Compare selected"]').click();
      await expect(page.locator('[id="chart-dialog-title"]')).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Comparison' })).toBeVisible();
    });

    test(`Delete a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileName, { exact: true }).first().click();
      await page.getByTestId('performanceProfileCard-delete').first().click();

      await expect(await page.getByText('Performance Profile Deleted!')).toBeVisible();
    });
  });
});
