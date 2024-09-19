import { test, expect } from '@playwright/test';
import { ENV } from './env';

export const test = base.extend({
  performancePage: async ({ page }, use) => {
    const perfPage = new PerformancePage(page);
    await perfPage.navigate();
    await use(perfPage);
  },
});

// Disable this test until got fixed
test.describe('Service Mesh Performance Management Tests', () => {
  test.skip();

  const profileName = 'Sample-test';

  test('View detailed result of a performance profile (Graph Visualiser)', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await page.waitForSelector(`text=${profileName}`, { state: 'visible' });
    await page.click(`button:has-text("View Results")`);
    await page.waitForSelector('button[aria-label="more"]', { state: 'visible' });
    await page.click('button[aria-label="more"]');

    await page.evaluate(() => {
      const sentinelStartDiv = document.querySelector('div[data-testid="sentinelStart"]');
      if (sentinelStartDiv) {
        sentinelStartDiv.setAttribute('data-testid', 'sentinel-graph');
      }
    });

    const graphVisible = await page.getByTestId('sentinel-graph').isVisible();
    expect(graphVisible);
  });

  test('Run a performance test', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/performance`);
    await page.getByRole('button', { name: 'Run Test', exact: true }).click();
    await page.getByLabel('Profile Name').fill(`${profileName}-2`);
    await page.locator('[aria-labelledby="meshName-label meshName"]').click();
    await page.locator('[data-value="istio"]').click();
    await page.getByRole('textbox', { name: 'url' }).fill('https://layer5.io/');
    await page.getByRole('spinbutton', { name: 'Concurrent requests' }).fill('5');
    await page.getByRole('spinbutton', { name: 'Queries per second' }).fill('5');
    await page.getByRole('textbox', { name: 'Duration' }).fill('15s');
    await expect(page.getByRole('button', { name: 'Run Test', exact: true })).toBeVisible();
    await page.getByRole('button', { name: 'Run Test', exact: true }).click();
    // Check for notification visibility
    await expect(page.getByText('Initiating load test . . .')).toBeVisible({ timeout: 90 * 1000 });
  });

  test('View Results from a performance profile', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    const perfTestProfile = await page.getByText(profileName).first();

    await expect(perfTestProfile).toBeVisible();

    await page.getByText('View Results').first().click();

    await page.getByLabel('more').first().click();

    await expect(page.getByText('Percentile Summary')).toBeVisible();
  });

  test('View/Edit the configuration of a performance profile', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await page.getByText(profileName).first().click();

    await page.getByTestId('performanceProfieCard-edit').click();

    await page.locator('[aria-labelledby="meshName-label meshName"]').click();
    await page.locator('[data-value="istio"]').click();
    await page.getByRole('spinbutton', { name: 'Concurrent requests' }).fill('6');
    await expect(page.getByRole('button', { name: 'Run Test', exact: true })).toBeVisible();
    await page.getByRole('button', { name: 'Run Test', exact: true }).click();

    await expect(page.getByText('Initiating load test . . .')).toBeVisible();
  });

  test('Delete a performance profile', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await page.getByText(profileName).first().click();
    await page.getByTestId('performanceProfieCard-delete').first().click();

    await expect(page.getByText('Performance Profile Deleted!')).toBeVisible();

    await page.getByText(`${profileName}-2`).first().click();
    await page.getByTestId('performanceProfieCard-delete').first().click();

    await expect(page.getByText('Performance Profile Deleted!')).toBeVisible();
  });
});
