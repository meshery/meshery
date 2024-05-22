const { ENV } = require('./env');
const { test, expect } = require('./fixtures/pages');

test('Aggregation Charts are displayed', async ({ page }) => {
  await page.goto(ENV.MESHERY_SERVER_URL);
  await expect(
    page.getByRole('heading', {
      name: 'Models by Category',
    }),
  ).toBeVisible();

  await expect(
    page.getByRole('heading', {
      name: 'Registry',
    }),
  ).toBeVisible();

  await expect(
    page.getByRole('heading', {
      name: 'Connections',
    }),
  ).toBeVisible();

  await expect(
    page.getByRole('heading', {
      name: 'Configuration',
    }),
  ).toBeVisible();
});
