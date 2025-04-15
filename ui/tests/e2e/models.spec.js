import { test, expect } from '@playwright/test';
import path from 'path';
import { ENV } from './env';

const model = {
  MODEL_URL: 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm',
  MODEL_NAME: `test-model-${Date.now()}`,
  MODEL_DISPLAY_NAME: `Test Model ${Date.now()}`,
};
const TAR_FILE_PATH = path.resolve('tests/e2e/assets/test.tar');

test.describe('Model Workflow Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(ENV.MESHERY_SERVER_URL);
    await page.getByTestId('settings-button').getByRole('button').click();
    await page.getByTestId('settings-tab-registry').click();
  });

  test('Create a Model', async ({ page }) => {
    await page.getByTestId('create-model-button').click();

    await page.getByRole('textbox', { name: 'Model Name' }).fill(model.MODEL_NAME);
    await page.getByRole('textbox', { name: 'Model Display Name' }).fill(model.MODEL_DISPLAY_NAME);

    await page.getByRole('button', { name: 'Next' }).click();
    await page.getByRole('button', { name: 'Next' }).click();
    await page.getByRole('button', { name: 'Next' }).click();

    await page.getByRole('radio', { name: 'GitHub' }).check();
    await page.getByRole('textbox', { name: 'Model URL' }).fill(model.MODEL_URL);

    await page.getByRole('button', { name: 'Next' }).click();
    await page.getByRole('checkbox', { name: /components in this model/i }).check();
    await page.getByRole('button', { name: 'Next' }).click();

    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.getByTestId('model-import-section')).toBeVisible();
    await expect(page.getByTestId('model-import-messages')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });

  test('Search a Model and Export it', async ({ page }) => {
    await page.getByTestId('search-icon').click();
    await page.getByRole('textbox', { name: 'Search' }).click();
    await page.getByRole('textbox', { name: 'Search' }).fill(model.MODEL_DISPLAY_NAME);
    await page.getByText(model.MODEL_DISPLAY_NAME).click();

    const downloadPromise = page.waitForEvent('download');
    await page.getByRole('button', { name: 'Export Model' }).click();
    const download = await downloadPromise;
    expect(download).toBeDefined();
  });

  test('Import a Model', async ({ page }) => {
    await page.getByTestId('import-model-button').click();
    await page.getByRole('heading', { name: 'File Import' }).click();

    await page.setInputFiles('input[type="file"]', TAR_FILE_PATH);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(page.getByTestId('model-import-section')).toBeVisible();
    await expect(page.getByTestId('model-import-messages')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });
});
