import { test, expect } from '@playwright/test';
import path from 'path';
import { DashboardPage } from './pages/DashboardPage';
import { ModelsPage } from './pages/ModelsPage';

const model = {
  MODEL_URL: 'git://github.com/aws-controllers-k8s/apigatewayv2-controller/main/helm',
  MODEL_NAME: `test-model-${Date.now()}`, // Fixed: Added backticks for template literal
  MODEL_DISPLAY_NAME: `Test Model ${Date.now()}`, // Fixed: Added backticks for template literal
};

const model_import = {
  MODEL_NAME: 'test', // Fixed: Added quotes around string
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
  let dashboardPage;
  let modelsPage;

  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page);
    modelsPage = new ModelsPage(page);

    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToSettings();
    await page.getByTestId('settings-tab-registry').click();
  });

  test('Create a Model', async () => {
    await modelsPage.createModelWorkflow(model);
  });

  test('Search a Model and Export it', async () => {
    const download = await modelsPage.searchExportAndChangeStatus(model.MODEL_DISPLAY_NAME);
    expect(download).toBeDefined();
  });

  test('Import a Model via File Import', async () => {
    await modelsPage.importModelViaFile(model_import);
  });

  test('Import a Model via Url Import', async () => {
    await modelsPage.importModelViaUrl(model_import);
  });

  test('Import a Model via CSV Import', async () => {
    await modelsPage.importModelViaCsv(model_import.MODEL_CSV_IMPORT);
  });
});
