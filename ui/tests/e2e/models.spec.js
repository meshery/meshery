import { test, expect } from '@playwright/test';
import path from 'path';
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
    modelsPage = new ModelsPage(page);
    await modelsPage.navigateToRegistrySettings();
  });

  test('Create a Model', async () => {
    await modelsPage.createModel(model);
  });

  test('Search a Model and Export it', async () => {
    await modelsPage.searchModel(model.MODEL_DISPLAY_NAME);
    const download = await modelsPage.exportModel();
    expect(download).toBeDefined();
    await modelsPage.setModelStatus('ignored');
  });

  test('Import a Model via File Import', async () => {
    await modelsPage.importModelViaFile(model_import.MODEL_FILE_IMPORT, model_import.MODEL_NAME);
  });

  test('Import a Model via Url Import', async () => {
    await modelsPage.importModelViaUrl(model_import.MODEL_URL_IMPORT, model_import.MODEL_NAME);
  });

  test('Import a Model via CSV Import', async () => {
    await modelsPage.importModelViaCsv(model_import.MODEL_CSV_IMPORT);
  });
});
