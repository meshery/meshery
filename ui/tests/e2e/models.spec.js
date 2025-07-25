import { test, expect } from '@playwright/test';
import path from 'path';
import { DashboardPage } from './pages/DashboardPage';
import { ModelsPage } from './pages/ModelsPage';

const model = {
  MODEL_URL: 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm',
  MODEL_NAME: `test-model-${Date.now()}`,
  MODEL_DISPLAY_NAME: `Test Model ${Date.now()}`,
};

const model_import = {
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
  let modelsPage;

  test.beforeEach(async ({ page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToSettings();
    await page.getByTestId('settings-tab-registry').click();

    modelsPage = new ModelsPage(page);
  });

  test('Create a Model', async ({ page }) => {
    await modelsPage.createModel({
      name: model.MODEL_NAME,
      displayName: model.MODEL_DISPLAY_NAME,
      url: model.MODEL_URL,
    });

    await expect(page.getByText(model.MODEL_DISPLAY_NAME)).toBeVisible();
  });

  test('Search a Model and Export it', async ({ page }) => {
    await modelsPage.searchAndExportModel(model.MODEL_DISPLAY_NAME);
  });

  test('Import a Model via File Import', async ({ page }) => {
    await modelsPage.importModelByFile(model_import.MODEL_FILE_IMPORT, model_import.MODEL_NAME);
    await expect(page.getByText(model_import.MODEL_NAME, { exact: true })).toBeVisible();
  });

  test('Import a Model via Url Import', async ({ page }) => {
    await modelsPage.importModelByUrl(model_import.MODEL_URL_IMPORT, model_import.MODEL_NAME);
    await expect(page.getByText(model_import.MODEL_NAME, { exact: true })).toBeVisible();
  });

  test('Import a Model via CSV Import', async ({ page }) => {
    await modelsPage.importModelByCsv(
      {
        models: model_import.MODEL_CSV_IMPORT.Models,
        components: model_import.MODEL_CSV_IMPORT.Components,
        relationships: model_import.MODEL_CSV_IMPORT.Relationships,
      },
      model_import.MODEL_CSV_IMPORT.Model_Name,
    );
    await expect(page.getByText(model_import.MODEL_CSV_IMPORT.Model_Name)).toBeVisible();
  });
});
