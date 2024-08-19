import { test } from '@playwright/test';
import { ENV } from './env';

test.describe('meshery settings performance', () => {
  test('Performance tab', async ({ page }) => {
    // Navigate to "User Preferences" page
    await page.goto(`${ENV.MESHERY_SERVER_URL}/user/preferences`);

    // Navigate to "Performance" tab
    await page.getByRole('tab', { name: 'Performance' }).click();

    // Fill concurrent requests, Queris per seconds, Time duration, and Generator
    await page.getByLabel('Concurrent requests *').click();
    await page.getByLabel('Concurrent requests *').fill('3');
    await page.getByLabel('Queries per second *').click();
    await page.getByLabel('Queries per second *').fill('3');
    await page.getByLabel('Duration*').click();
    await page.getByRole('option', { name: '30s' }).click();
    await page.getByLabel('nighthawk').check();

    // Click on "Save" button
    await page.getByRole('button', { name: 'Save' }).click();

    // And Preferences are saved successfully
    await page.getByText('Preferences saved').click();
  });
});
