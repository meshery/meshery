import { expect } from '@playwright/test';
import { DashboardPage } from './DashboardPage';

export class ModelsPage {
  constructor(page) {
    this.page = page;
    this.dashboardPage = new DashboardPage(page);

    this.registryTab = page.getByTestId('settings-tab-registry');

    this.createModelButton = page.getByTestId('TabBar-Button-CreateModel');
    this.importModelButton = page.getByTestId('TabBar-Button-ImportModel');

    this.modelNameInput = page.locator('#model-name');
    this.modelDisplayNameInput = page.locator('#model-display-name');
    this.modelUrlInput = page.locator('#model-url');

    this.stepNextButton = page.getByTestId('UrlStepper-Button-Next');
    this.stepGenerateButton = page.getByTestId('UrlStepper-Button-Generate');
    this.stepFinishButton = page.getByTestId('UrlStepper-Button-Finish');

    this.categorySelect = page.getByTestId('UrlStepper-Select-Category');
    this.subcategorySelect = page.getByTestId('UrlStepper-Select-Subcategory');

    this.logoDarkSelect = page.getByTestId('UrlStepper-Select-Logo-Dark-Theme');
    this.logoLightSelect = page.getByTestId('UrlStepper-Select-Logo-Light-Theme');
    this.primaryColorSelect = page.getByTestId('UrlStepper-Select-Primary-Color');
    this.secondaryColorSelect = page.getByTestId('UrlStepper-Select-Secondary-Color');
    this.shapeSelect = page.getByTestId('UrlStepper-Select-Shape');

    this.sourceGithubRadio = page.getByTestId('UrlStepper-Select-Source-GitHub');
    this.visualAnnotationCheckbox = page.getByTestId('UrlStepper-Visual-Annotation-Checkbox');

    this.searchIcon = page.getByTestId('search-icon');
    this.searchInput = page.locator('#searchClick');
    this.exportModelButton = page.getByTestId('export-model-button');
    this.statusCombobox = page.locator('[aria-labelledby="entity-status-select-label"]');

    this.fileInput = page.locator('input[type="file"]');

    this.importFileHeading = page.getByRole('heading', { name: 'File Import' });
    this.importUrlHeading = page.getByRole('heading', { name: 'URL Import' });
    this.importCsvHeading = page.getByRole('heading', { name: 'CSV Import' });
    this.importUrlInput = page.getByRole('textbox', { name: 'URL' });
    this.importNextButton = page.getByRole('button', { name: 'Next' });
    this.importGenerateButton = page.getByRole('button', { name: 'Generate' });
    this.importFinishButton = page.getByRole('button', { name: 'Finish' });
  }

  async navigateToRegistrySettings() {
    await this.dashboardPage.navigateToDashboard();
    await this.dashboardPage.navigateToSettings();
    await this.registryTab.click();
  }

  async createModel(model) {
    await this.createModelButton.click();

    await this.modelNameInput.fill(model.MODEL_NAME);
    await this.modelDisplayNameInput.fill(model.MODEL_DISPLAY_NAME);

    await this.stepNextButton.click();

    await expect(this.categorySelect).toBeVisible();
    await expect(this.subcategorySelect).toBeVisible();

    await this.stepNextButton.click();

    await expect(this.logoDarkSelect).toBeVisible();
    await expect(this.logoLightSelect).toBeVisible();

    await expect(this.primaryColorSelect).toBeVisible();
    await expect(this.secondaryColorSelect).toBeVisible();
    await expect(this.shapeSelect).toBeVisible();

    await this.stepNextButton.click();

    await this.sourceGithubRadio.check();
    await this.modelUrlInput.fill(model.MODEL_URL);

    await this.stepNextButton.click();
    await this.visualAnnotationCheckbox.check();
    await this.stepNextButton.click();

    await this.stepGenerateButton.click();

    await this.expectModelImportSuccess(model.MODEL_NAME);

    await this.stepFinishButton.click();
  }

  async searchModel(displayName) {
    await this.searchIcon.click();
    await this.searchInput.click();
    await this.searchInput.fill(displayName);
    await this.page.getByText(displayName).click();
  }

  async exportModel() {
    const downloadPromise = this.page.waitForEvent('download');
    await this.exportModelButton.click();
    return downloadPromise;
  }

  async setModelStatus(statusLabel) {
    await this.statusCombobox.click();
    await this.page.getByRole('option', { name: statusLabel }).click();
    await expect(this.statusCombobox).toHaveText(new RegExp(statusLabel, 'i'));
  }

  async importModelViaFile(filePath, modelName) {
    await this.importModelButton.click();
    await this.importFileHeading.click();

    await this.fileInput.setInputFiles(filePath);
    await this.importNextButton.click();

    await this.expectModelImportSuccess(modelName);
    await this.importFinishButton.click();
  }

  async importModelViaUrl(modelUrl, modelName) {
    await this.importModelButton.click();
    await this.importUrlHeading.click();

    await this.importUrlInput.fill(modelUrl);

    await this.importNextButton.click();

    await this.expectModelImportSuccess(modelName);
    await this.importFinishButton.click();
  }

  async importModelViaCsv(csvImport) {
    await this.importModelButton.click();
    await this.importCsvHeading.click();

    await this.importNextButton.click();

    await this.fileInput.setInputFiles(csvImport.Models);
    await this.importNextButton.click();
    await this.fileInput.setInputFiles(csvImport.Components);
    await this.importNextButton.click();
    await this.fileInput.setInputFiles(csvImport.Relationships);

    await this.importGenerateButton.click();

    await this.expectModelImportSuccess(csvImport.Model_Name);
    await this.importFinishButton.click();
  }

  async expectModelImportSuccess(modelName) {
    await expect(
      this.page.getByTestId(`ModelImportedSection-ModelHeader-${modelName}`),
    ).toBeVisible();
    await expect(this.page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
  }
}
