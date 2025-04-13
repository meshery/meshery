import { test, expect } from '@playwright/test';
import path from 'path';

const MODEL_URL = 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm';
const MODEL_NAME = 'test-model';
const MODEL_DISPLAY_NAME = 'Test Model';
const TAR_FILE_PATH = path.resolve('tests/e2e/assets/test.tar');
const APP_URL = 'http://localhost:9081/';

test.beforeEach(async ({ page }) => {
  await page.goto(APP_URL);
  await page.getByTestId('settings-button').getByRole('button').click();
  await page.getByTestId('settings-tab-registry').click();
});

test('Create a Model', async ({ page }) => {
  await page.getByText('Create').click();
  await page.locator('div', { hasText: 'Create Model' }).nth(1).click();

  await page.getByRole('textbox', { name: 'Model Name' }).fill(MODEL_NAME);
  await page.getByRole('textbox', { name: 'Model Display Name' }).fill(MODEL_DISPLAY_NAME);

  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('button', { name: 'Next' }).click();

  await page.getByRole('radio', { name: 'GitHub' }).check();
  await page.getByRole('textbox', { name: 'Model URL' }).fill(MODEL_URL);

  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('checkbox', { name: /components in this model/i }).check();
  await page.getByRole('button', { name: 'Next' }).click();

  await page.getByRole('button', { name: 'Generate' }).click();
  await page.getByRole('button', { name: 'Finish' }).click();
});

test('Search a Model and Export it', async ({ page }) => {
  await page.getByRole('button', { name: 'Search' }).click();
  await page.getByRole('textbox', { name: 'Search' }).click();
  await page.getByRole('textbox', { name: 'Search' }).fill(MODEL_DISPLAY_NAME);
  await page.getByText(MODEL_DISPLAY_NAME).click();

  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Export Model' }).click();
  const download = await downloadPromise;
  expect(download).toBeDefined();
});

test('Import a Model', async ({ page }) => {
  await page.getByText('Import').click();
  await page.getByRole('heading', { name: 'File Import' }).click();

  console.log('TAR_FILE_PATH', TAR_FILE_PATH);
  await page.setInputFiles('input[type="file"]', TAR_FILE_PATH);

  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('button', { name: 'Finish' }).click();
});
