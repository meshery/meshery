import { test, expect } from './fixtures/project';
import { ENV } from './env';

const DESIGN_FILE_NAME = 'E2E Test Design - Design Spec';
const DESIGN_URI =
  'https://raw.githubusercontent.com/PR4NJ41/os-files/refs/heads/main/E2E%20Test%20Design%20-%20Test%20File%20k8s-manifest.yml';
const SUCCESS_NOTIFICATION = (design) => `Imported design '${design}' of type 'meshery-design'`;
const DELETE_NOTIFICATION = (design) => `"${design}" Design deleted`;

test.beforeEach(async ({ page }) => {
  await page.goto(`${ENV.MESHERY_SERVER_URL}/configuration/designs`);
  await expect(page.getByRole('heading', { name: 'Configuration' })).toBeVisible();
});

test('Import a design using URI', async ({ page }) => {
  await page.locator('button', { hasText: 'Import Design' }).click();
  await page.getByLabel('Design file name *').fill(DESIGN_FILE_NAME);
  await page.locator('input[name="root_uploadType"][value="1"]').click();
  await page.getByLabel('URL *').fill(DESIGN_URI);
  await page.getByRole('button', { name: 'Import' }).click();

  const notification = page.locator('div.notistack-Snackbar', {
    hasText: SUCCESS_NOTIFICATION(DESIGN_FILE_NAME),
  });
  await expect(notification).toBeVisible();
});

test('Delete design', async ({ page }) => {
  const designLocator = page.locator(`text=${DESIGN_FILE_NAME}`);
  await expect(designLocator).toBeVisible();
  await designLocator.click();

  await page.locator('svg[data-testid="DeleteIcon"]').click();
  await page.locator('button:has-text("DELETE")').click();

  const notification = page.locator('div.notistack-Snackbar', {
    hasText: DELETE_NOTIFICATION(DESIGN_FILE_NAME),
  });
  await expect(notification).toBeVisible();
});
