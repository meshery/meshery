import { expect, test } from './fixtures/project';
import { mockPerfApis } from './mocks/mockPerfApi';
//
const performanceProfiles = [
  {
    profileName: 'Fortio-Perf-Test',
    serviceMesh: 'None',
    url: 'https://meshery.io/',
    loadGenerator: 'fortio',
    concurrentRequest: '2',
    qps: '2',
    duration: '15s',
  },
];

performanceProfiles.forEach((config) => {
  const { profileName, serviceMesh, url, loadGenerator, concurrentRequest, qps, duration } = config;

  test.describe(`Performance Management Tests with ${loadGenerator}`, () => {
    test.beforeEach(async ({ page }) => {
      await mockPerfApis(page, config);
      await page.goto('/');
      await page.getByRole('button', { name: 'meshery-button-2 Performance' }).click();
      await page.getByTestId('navigation').getByRole('button', { name: 'Profiles' }).click();
    });

    test(`Add performance profile with load generator ${loadGenerator}`, async ({ page }) => {
      await page.getByLabel('Add Performance Profile').click();
      await page.getByLabel('Profile Name').fill(profileName);
      await page.getByLabel('Technology').click();
      await page.getByRole('option', { name: serviceMesh }).click();
      await page.getByLabel('URL to test').fill(url);
      await page.getByLabel('Concurrent requests').fill(concurrentRequest);
      await page.getByLabel('Queries per second').fill(qps);
      await page.getByLabel('Duration').fill(duration);
      await page.getByLabel(loadGenerator).check();
      await page.getByTestId('run-performance-test').click();

      await expect(
        page.getByTestId('SnackbarContent-success').filter({ hasText: /submitted/i }),
      ).toHaveCount(1);
    });

    test(`View detailed result of a performance profile (Graph Visualiser) with load generator ${loadGenerator}`, async ({
      page,
    }) => {
      await expect(page.getByText(profileName)).toBeVisible();
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.getByTestId('TableChartIcon').first().click();
      await expect(page.getByText(url, { exact: true })).toBeVisible();
    });

    test(`Edit the configuration of a performance profile with load generator ${loadGenerator} and service mesh ${serviceMesh}`, async ({
      page,
    }) => {
      await page.getByText(profileName).click();
      await page.getByTestId('performanceProfileCard-edit').click();
      await page.getByLabel('Technology').click();
      await page.getByRole('option', { name: serviceMesh }).click();
      await page.getByLabel('Concurrent requests').fill('3');
      await page.getByLabel(loadGenerator).check();
      await page.getByTestId('run-performance-test').click();

      await expect(
        page.getByTestId('SnackbarContent-success').filter({ hasText: /submitted/i }),
      ).toHaveCount(1);
    });

    test(`Compare test of a performance profile with load generator ${loadGenerator}`, async ({
      page,
    }) => {
      await page.getByRole('button', { name: 'View Results', exact: true }).first().click();
      await page.getByTestId('MUIDataTableBodyRow-0').locator('input[type="checkbox"]').check();
      await page.getByTestId('CompareArrowsIcon').first().click();
      await expect(page.getByTestId('chart-dialog-title')).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Comparison' })).toBeVisible();
    });

    test(`Delete a performance profile with load generator ${loadGenerator}`, async ({ page }) => {
      await page.getByText(profileName).click();
      await page.getByTestId('performanceProfileCard-delete').first().click();
      await expect(
        page.getByTestId('SnackbarContent-success').filter({ hasText: /deleted/i }),
      ).toHaveCount(1);
    });
  });
});
