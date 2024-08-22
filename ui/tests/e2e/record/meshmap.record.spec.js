import { expect, test } from '../fixtures/recorder';

test.describe('Record a videos', () => {
  test('Go to MeshMap extension and check the first design', async ({ page }) => {
    await page.goto('/extension/meshmap');

    await page.getByTestId('MuiDataTableBodyCell-2-0').first().click();

    await expect(page.locator('#custom-noti').getByRole('alert')).toBeVisible();
  });
});
