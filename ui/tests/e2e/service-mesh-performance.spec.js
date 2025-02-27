import { ENV } from './env';
import { v4 as uuidv4 } from 'uuid';
import { expect, test } from './fixtures/project';

test.describe('Service Mesh Performance Management Tests', { tag: '@unstable' }, () => {
  const peformanceProfiles = [
    {
      profileWithUUID: `Sample-Perf-Test-${uuidv4()}`,
      serviceMesh: 'None',
      url: 'https://layer5.io/',
      loadGenerator: 'fortio',
    },
  ];

  peformanceProfiles.forEach(({ profileWithUUID, serviceMesh, url, loadGenerator }) => {
    test(`Add performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByLabel('Add Performance Profile').click();
      await page.getByLabel('Profile Name').fill(profileWithUUID);
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
      await expect(await page.getByText(`${profileWithUUID}`)).toBeVisible();
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.getByTestId('open-performance-result-bar-chart').first().click();
      await expect(page.getByTestId('sentinelStart')).toBeHidden();
      await expect(await page.getByText(`URL: ${url}`, { exact: true })).toBeVisible();
    });

    test(`Edit the configuration of a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileWithUUID, { exact: true }).first().click();

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
      await page.getByText(profileWithUUID, { exact: true }).first().click();
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.getByTestId('MUIDataTableBodyRow-0').locator('input[type="checkbox"]').check();

      await page.getByTestId('compare-selected').click();
      await expect(page.getByTestId('chart-dialog-title')).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Comparison' })).toBeVisible();
    });

    test(`Delete a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
      await page.getByText(profileWithUUID, { exact: true }).first().click();
      await page.getByTestId('performanceProfileCard-delete').first().click();

      await expect(await page.getByText('Performance Profile Deleted!')).toBeVisible();
    });
  });
});
