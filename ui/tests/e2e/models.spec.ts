import { test, expect } from '@playwright/test';
import path from 'path';
import { DashboardPage } from './pages/DashboardPage';

// Strongly typed inline to avoid JS linter false positives
const model: {
  MODEL_URL: string;
  MODEL_NAME: string;
  MODEL_DISPLAY_NAME: string;
} = {
  MODEL_URL: 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm',
  MODEL_NAME: `test-model-${Date.now()}`,
  MODEL_DISPLAY_NAME: `Test Model ${Date.now()}`,
};

// Strongly typed inline
const model_import: {
  MODEL_NAME: string;
  MODEL_URL_IMPORT: string;
  MODEL_FILE_IMPORT: string;
  MODEL_CSV_IMPORT: {
    Model_Name: string;
    Models: string;
    Components: string;
    Relationships: string;
  };
} = {
  MODEL_NAME: `test`,
  MODEL_URL_IMPORT:
    'https://raw.githubusercontent.com/meshery/meshery/master/ui/tests/e2e/assets/test.tar',
  MODEL_FILE_IMPORT: path.resolve('tests/e2e/assets/test.tar'),
  MODEL_CSV_IMPORT: {
    Model_Name: 'couchbase',
    Models: path.resolve('tests/e2e/assets/models.csv'),
    Components: path.resolve('tests/e2e/assets/components.csv'),
    Relationships: path.resolve('tests/e2e/assets/relationships.csv'),
  },
};

test.describe.serial('Model Workflow Tests', () => {
  test.describe.configure({ timeout: 240_000 });

  test.beforeEach(async ({ page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToSettings();
    await expect(page.getByTestId('settings-tab-registry')).toBeVisible();
    await page.getByTestId('settings-tab-registry').click();
  });

  test('Create a Model', async ({ page }) => {
    test.slow();
    test.setTimeout(600_000);

    await page.getByTestId('TabBar-Button-CreateModel').click();

    // UPGRADE 1: Replaced fragile #id locators with robust getByRole selectors
    await page.getByRole('textbox', { name: /model name/i }).fill(model.MODEL_NAME);
    await page.getByRole('textbox', { name: /model display name/i }).fill(model.MODEL_DISPLAY_NAME);

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
    await page.getByRole('textbox', { name: /url/i }).fill(model.MODEL_URL);

    await page.getByTestId('UrlStepper-Button-Next').click();
    await page.getByTestId('UrlStepper-Visual-Annotation-Checkbox').check();
    await page.getByTestId('UrlStepper-Button-Next').click();

    await page.getByTestId('UrlStepper-Button-Generate').click();

    await expect(
      page.getByTestId(`ModelImportedSection-ModelHeader-${model.MODEL_NAME}`),
    ).toBeVisible({ timeout: 540_000 });
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();

    await page.getByTestId('UrlStepper-Button-Finish').click();
  });

  test('Search a Model and Export it', async ({ page }) => {
    await page.getByTestId('search-icon').click();
    
    // UPGRADE 2: Replaced #searchClick with semantic getByRole searchbox
    await page.getByRole('searchbox', { name: /search/i }).fill(model.MODEL_DISPLAY_NAME);
    
    // UPGRADE 3: Added { exact: true } to prevent matching partial text elsewhere on the page
    await page.getByText(model.MODEL_DISPLAY_NAME, { exact: true }).click();

    const downloadPromise = page.waitForEvent('download');
    await page.getByTestId('export-model-button').click();
    const download = await downloadPromise;
    expect(download).toBeDefined();
    
    await page.getByRole('combobox', { name: /model status/i }).click();
    await page.getByRole('option', { name: /ignored/i }).click();
    await expect(page.getByRole('combobox', { name: /model status/i })).toHaveText(/ignored/i);
  });

  test('Import a Model via File Import', async ({ page }) => {
    await page.getByTestId('TabBar-Button-ImportModel').click();
    await page.getByRole('heading', { name: 'File Import', exact: true }).click();

    await page.locator('input[type="file"]').setInputFiles(model_import.MODEL_FILE_IMPORT);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(
      page.getByTestId(`ModelImportedSection-ModelHeader-${model_import.MODEL_NAME}`),
    ).toBeVisible();
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });

  test('Import a Model via Url Import', async ({ page }) => {
    await page.getByTestId('TabBar-Button-ImportModel').click();
    await page.getByRole('heading', { name: 'URL Import', exact: true }).click();

    await page.getByRole('textbox', { name: /url/i }).fill(model_import.MODEL_URL_IMPORT);

    await page.getByRole('button', { name: 'Next' }).click();

    await expect(
      page.getByTestId(`ModelImportedSection-ModelHeader-${model_import.MODEL_NAME}`),
    ).toBeVisible();
    await expect(page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
    await page.getByRole('button', { name: 'Finish' }).click();
  });

  test('Import a Model via CSV Import', async ({ page }) => {
    await page.getByTestId('TabBar-Button-ImportModel').click();
    await page.getByRole('heading', { name: 'CSV Import', exact: true }).click();

    await page.getByRole('button', { name: 'Next' }).click();

    await page.locator('input[type="file"]').setInputFiles(model_import.MODEL_CSV_IMPORT.Models);
    await page.getByRole('button', { name: 'Next' }).click();
    await page.locator('input[type="file"]').setInputFiles(model_import.MODEL_CSV_IMPORT.Components);
    await page.getByRole('button', { name: 'Next' }).click();
    await page.locator('input[type="file"]').setInputFiles(model_import.MODEL_CSV_IMPORT.Relationships);

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