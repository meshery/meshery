import { test, expect } from '@playwright/test';
import path from 'path';
import { ENV } from './env';

const model = {
  MODEL_URL: 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm',
  MODEL_NAME: `test-model-${Date.now()}`,
  MODEL_DISPLAY_NAME: `Test Model ${Date.now()}`,
};
const TAR_FILE_PATH = path.resolve('tests/e2e/assets/test.tar');
const IMPORT_MODEL_URL =
  'https://raw.githubusercontent.com/vr-varad/meshery/refs/heads/feat/create_import_model_testing/ui/tests/e2e/assets/test.tar';

test.describe('Model Workflow Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(ENV.MESHERY_SERVER_URL);
    await page.getByTestId('settings-button').click();
    await page.getByTestId('settings-tab-registry').click();
  });

  test('Create a Model', async ({ page }) => {
    await page.getByTestId('create-model-button').click();

    await page.locator('#model-name').fill(model.MODEL_NAME);
    await page.locator('#model-display-name').fill(model.MODEL_DISPLAY_NAME);

    await page.getByTestId('next-button').click();

    await expect(page.getByTestId('category-select')).toBeVisible();
    await expect(page.getByTestId('subcategory-select')).toBeVisible();

    await page.getByTestId('next-button').click();

    await expect(page.getByTestId('logo-dark-theme')).toBeVisible();
    await expect(page.getByTestId('logo-light-theme')).toBeVisible();

    await expect(page.getByTestId('primary-color')).toBeVisible();
    await expect(page.getByTestId('secondary-color')).toBeVisible();
    await expect(page.getByTestId('shape-select')).toBeVisible();

    await page.getByTestId('next-button').click();

    await page.getByTestId('source-github').check();
    await page.locator('#model-url').fill(model.MODEL_URL);

    await page.getByTestId('next-button').click();
    await page.getByTestId('visual-annotation-checkbox').check();
    await page.getByTestId('next-button').click();

    await page.getByTestId('generate-button').click();

    await expect(page.getByTestId('model-import-section')).toBeVisible();
    await expect(page.getByTestId('model-import-messages')).toBeVisible();

    await page.getByTestId('finish-button').click();
  });

  test('Search a Model and Export it', async ({ page }) => {
    await page.getByTestId('search-icon').click();
    await page.locator('#searchClick').click();
    await page.locator('#searchClick').fill(model.MODEL_DISPLAY_NAME);
    await page.getByText(model.MODEL_DISPLAY_NAME).click();

    const downloadPromise = page.waitForEvent('download');
    await page.getByTestId('export-model-button').click();
    const download = await downloadPromise;
    expect(download).toBeDefined();
  });

  test('Import a Model via File Import', async ({ page }) => {
    await page.getByTestId('import-model-button').click();
    await page.getByRole('heading', { name: 'File Import' }).click();

    await page.setInputFiles('input[type="file"]', TAR_FILE_PATH);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(page.getByTestId('model-import-section')).toBeVisible();
    await expect(page.getByTestId('model-import-messages')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });

  test('Import a Model via Url Import', async ({ page }) => {
    await page.getByTestId('import-model-button').click();
    await page.getByRole('heading', { name: 'URL Import' }).click();

    await page.getByRole('textbox', { name: 'URL' }).click();
    await page.getByRole('textbox', { name: 'URL' }).fill(IMPORT_MODEL_URL);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(page.getByTestId('model-import-section')).toBeVisible();
    await expect(page.getByTestId('model-import-messages')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });
});
