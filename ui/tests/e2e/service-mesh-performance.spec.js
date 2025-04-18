import { v4 as uuidv4 } from 'uuid';
import { expect, test } from './fixtures/project';

test.describe.serial('Service Mesh Performance Management Tests', () => {
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
      await page.goto('/performance/profiles');
      await page.getByLabel('Add Performance Profile').click();
      await page.getByLabel('Profile Name').fill(profileWithUUID);
      await page.getByLabel('Technology').click();
      await page.getByRole('option', { name: serviceMesh }).click();
      await page.getByLabel('URL to test').fill(url);
      await page.getByLabel('Concurrent requests').fill('2');
      await page.getByLabel('Queries per second').fill('2');
      await page.getByLabel('Duration').fill('15s');
      await page.getByLabel(loadGenerator).check();
      await page.getByTestId('run-performance-test').click();
      const successNofications = page.getByTestId('SnackbarContent-success');
      const infoNofications = page.getByTestId('SnackbarContent-info');

      await expect(successNofications.filter({ hasText: /submitted/i })).toHaveCount(1);
      await expect(infoNofications.filter({ hasText: /completed/i })).toHaveCount(1);
    });

    test(`View detailed result of a performance profile (Graph Visualiser) with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto('/performance/profiles');
      await expect(page.getByText(`${profileWithUUID}`)).toBeVisible();
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.getByTestId('TableChartIcon').first().click();
      await expect(page.getByText(`${url}`, { exact: true })).toBeTruthy();
    });

    test(`Edit the configuration of a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto('/performance/profiles');
      await page.getByText(profileWithUUID, { exact: true }).first().click();

      await page.getByTestId('performanceProfileCard-edit').click();
      await page.getByLabel('Technology').click();
      await page.getByRole('option', { name: serviceMesh }).click();
      await page.getByLabel('Concurrent requests').fill('3');
      await page.getByLabel(loadGenerator).check();

      await page.getByTestId('run-performance-test').click();
      const successNofications = page.getByTestId('SnackbarContent-success');
      const infoNofications = page.getByTestId('SnackbarContent-info');

      await expect(successNofications.filter({ hasText: /submitted/i })).toHaveCount(1);
      await expect(infoNofications.filter({ hasText: /completed/i })).toHaveCount(1);
    });

    test(`Compare test of a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto('/performance/profiles');
      await page.getByText(profileWithUUID, { exact: true }).first().click();
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.getByTestId('MUIDataTableBodyRow-0').locator('input[type="checkbox"]').check();

      await page.getByTestId('CompareArrowsIcon').first().click();
      await expect(page.getByTestId('chart-dialog-title')).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Comparison' })).toBeVisible();
    });

    test(`Delete a performance profile with load generator "${loadGenerator}" and service mesh "${serviceMesh}"`, async ({
      page,
    }) => {
      await page.goto('performance/profiles');
      await page.getByText(profileWithUUID, { exact: true }).first().click();
      await page.getByTestId('performanceProfileCard-delete').first().click();

      await expect(page.getByTestId('SnackbarContent-success').filter(/deleted/i)).toHaveCount(1);
    });
  });
});
