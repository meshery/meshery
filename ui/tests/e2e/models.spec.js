import { test, expect } from '@playwright/test';
import path from 'path';

const model = {
  MODEL_URL: 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm',
  MODEL_NAME: `test-model-${Date.now()}`,
  MODEL_DISPLAY_NAME: `Test Model ${Date.now()}`,
};

const model_import = {
  MODEL_NAME: `test`,
  MODEL_URL_IMPORT:
    'https://raw.githubusercontent.com/vr-varad/meshery/refs/heads/feat/create_import_model_testing/ui/tests/e2e/assets/test.tar',
  MODEL_FILE_IMPORT: path.resolve('tests/e2e/assets/test.tar'),
  MODEL_CSV_IMPORT: {
    Model_Name: 'couchbase',
    Models: path.resolve('tests/e2e/assets/models.csv'),
    Components: path.resolve('tests/e2e/assets/components.csv'),
    Relationships: path.resolve('tests/e2e/assets/relationships.csv'),
  },
};

test.describe.serial('Model Workflow Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.getByTestId('settings-button').click();
    await page.getByTestId('settings-tab-registry').click();
  });

  test('Create a Model', async ({ page }) => {
    await page.getByTestId('TabBar-Button-CreateModel').click();

    await page.locator('#model-name').fill(model.MODEL_NAME);
    await page.locator('#model-display-name').fill(model.MODEL_DISPLAY_NAME);

    await page.getByTestId('UrlStepper-Button-Next').click();

    await expect(page.getByTestId('UrlStepper-Select-Category')).toBeVisible();
    await expect(page.getByTestId('UrlStepper-Select-Subcategory')).toBeVisible();

    await page.getByTestId('UrlStepper-Button-Next').click();

    await expect(page.getByTestId('UrlStepper-Select-Logo-Dark-Theme')).toBeVisible();
    await expect(page.getByTestId('UrlStepper-Select-Logo-Light-Theme')).toBeVisible();

    await expect(page.getByTestId('UrlStepper-Select-Primary-Color')).toBeVisible();
    await expect(page.getByTestId('UrlStepper-Select-Secondary-Color')).toBeVisible();
    await expect(page.getByTestId('UrlStepper-Select-Shape')).toBeVisible();

    await page.getByTestId('UrlStepper-Button-Next').click();

    await page.getByTestId('UrlStepper-Select-Source-GitHub').check();
    await page.locator('#model-url').fill(model.MODEL_URL);

    await page.getByTestId('UrlStepper-Button-Next').click();
    await page.getByTestId('UrlStepper-Visual-Annotation-Checkbox').check();
    await page.getByTestId('UrlStepper-Button-Next').click();

    await page.getByTestId('UrlStepper-Button-Generate').click();

    await expect(
      page.getByTestId(`ModelImportedSection-ModelHeader-${model.MODEL_NAME}`),
    ).toBeVisible();
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();

    await page.getByTestId('UrlStepper-Button-Finish').click();
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
    await page.getByRole('combobox', { name: 'enabled' }).click();
    await page.getByRole('option', { name: 'ignored' }).click();
    expect(page.getByRole('option', { name: 'ignored' }).isVisible()).toBeTruthy();
  });

  test('Import a Model via File Import', async ({ page }) => {
    await page.getByTestId('TabBar-Button-ImportModel').click();
    await page.getByRole('heading', { name: 'File Import' }).click();

    await page.setInputFiles('input[type="file"]', model_import.MODEL_FILE_IMPORT);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(
      page.getByTestId(`ModelImportedSection-ModelHeader-${model_import.MODEL_NAME}`),
    ).toBeVisible();
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });

  test('Import a Model via Url Import', async ({ page }) => {
    await page.getByTestId('TabBar-Button-ImportModel').click();
    await page.getByRole('heading', { name: 'URL Import' }).click();

    await page.getByRole('textbox', { name: 'URL' }).click();
    await page.getByRole('textbox', { name: 'URL' }).fill(model_import.MODEL_URL_IMPORT);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(
      page.getByTestId(`ModelImportedSection-ModelHeader-${model_import.MODEL_NAME}`),
    ).toBeVisible();
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });

  test('Import a Model via CSV Import', async ({ page }) => {
    await page.getByTestId('TabBar-Button-ImportModel').click();
    await page.getByRole('heading', { name: 'CSV Import' }).click();

    await page.getByRole('button', { name: 'Next' }).click();

    await page.setInputFiles('input[type="file"]', model_import.MODEL_CSV_IMPORT.Models);
    await page.getByRole('button', { name: 'Next' }).click();
    await page.setInputFiles('input[type="file"]', model_import.MODEL_CSV_IMPORT.Components);
    await page.getByRole('button', { name: 'Next' }).click();
    await page.setInputFiles('input[type="file"]', model_import.MODEL_CSV_IMPORT.Relationships);

    await page.getByRole('button', { name: 'Generate' }).click();

    await expect(
      page.getByTestId(
        `ModelImportedSection-ModelHeader-${model_import.MODEL_CSV_IMPORT.Model_Name}`,
      ),
    ).toBeVisible();
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });
});
