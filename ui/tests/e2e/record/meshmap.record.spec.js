import { expect, test } from '@playwright/test';
import { ENV } from '../env';

test.describe('Test meshmap and record', () => {
  test('Go to MeshMap extension and check skip on the tutorial modal', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/extension/meshmap`);

    await expect(page.getByText('Tutorial')).toBeVisible();

    await page.getByLabel('Do not show again').check();

    await page.locator('#tutorial-dialog-title svg').nth(2).click();
  });

  test('Go to MeshMap extension and go to the first design', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/extension/meshmap`);

    await page.getByTestId('MuiDataTableBodyCell-2-0').first().click();

    await expect(page.locator('#custom-noti').getByRole('alert')).toBeVisible();
  });

  test('Go to MeshMap extension settings and turn on again tutorial', async ({ page }) => {
    await page.goto(`${ENV.MESHERY_SERVER_URL}/extension/meshmap`);

    await page.getByLabel('CANVAS_CONTROLS').nth(1).hover();

    await page.getByRole('menuitem', { name: 'Options' }).click();

    await expect(page.getByText('Show Tutorial')).toBeVisible();

    // Sign of bad UX, should implemented optimistic updates so user not be confused
    await page
      .locator('li')
      .filter({ hasText: 'Show Tutorial' })
      .getByRole('checkbox')
      .check({ force: true });

    await page.getByRole('button').first().click();
  });
});
