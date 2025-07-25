import { expect } from '@playwright/test';

export class ModelsPage {
  constructor(page) {
    this.page = page;

    // Locators for common buttons and inputs
    this.createModelButton = page.getByTestId('TabBar-Button-CreateModel');
    this.importModelButton = page.getByTestId('TabBar-Button-ImportModel');
    this.nextButton = page.getByRole('button', { name: 'Next' });
    this.stepperNextButton = page.getByTestId('UrlStepper-Button-Next');
    this.generateButton = page.getByRole('button', { name: 'Generate' });
    this.finishButton = page.getByRole('button', { name: 'Finish' });
    this.fileInput = page.locator('input[type="file"]');

    // Locators for the 'Create Model' workflow
    this.modelNameInput = page.locator('#model-name');
    this.modelDisplayNameInput = page.locator('#model-display-name');
    this.githubSourceCheckbox = page.getByTestId('UrlStepper-Select-Source-GitHub');
    this.modelUrlInput = page.locator('#model-url');
    this.visualAnnotationCheckbox = page.getByTestId('UrlStepper-Visual-Annotation-Checkbox');

    // Locators for searching and viewing models
    this.searchIcon = page.getByTestId('search-icon');
    this.searchInput = page.locator('#searchClick');
    this.exportModelButton = page.getByTestId('export-model-button');

    // Locators for the 'Import Model' workflow
    this.fileImportTab = page.getByRole('heading', { name: 'File Import' });
    this.urlImportTab = page.getByRole('heading', { name: 'URL Import' });
    this.csvImportTab = page.getByRole('heading', { name: 'CSV Import' });
    this.urlImportInput = page.getByRole('textbox', { name: 'URL' });
  }

  /**
   * Complete create model workflow
   * @param {{name: string, displayName: string, url: string}} modelData
   */
  async createModel(modelData) {
    await this.createModelButton.click();

    await this.modelNameInput.fill(modelData.name);
    await this.modelDisplayNameInput.fill(modelData.displayName);
    await this.stepperNextButton.click();

    await expect(this.page.getByTestId('UrlStepper-Select-Category')).toBeVisible();
    await this.stepperNextButton.click();

    await expect(this.page.getByTestId('UrlStepper-Select-Logo-Dark-Theme')).toBeVisible();
    await this.stepperNextButton.click();

    await this.githubSourceCheckbox.check();
    await this.modelUrlInput.fill(modelData.url);
    await this.stepperNextButton.click();

    await this.visualAnnotationCheckbox.check();
    await this.stepperNextButton.click();

    await this.page.getByTestId('UrlStepper-Button-Generate').click();
    await this.verifyModelImported(modelData.name);
    await this.verifyImportMessagesVisible();
    await this.page.getByTestId('UrlStepper-Button-Finish').click();
  }

  /**
   * searches for a model by its display name, opens it, and exports it.
   * @param {string} modelDisplayName
   */
  async searchAndExportModel(modelDisplayName) {
    await this.searchIcon.click();
    await this.searchInput.click();
    await this.searchInput.fill(modelDisplayName);
    await this.page.getByText(modelDisplayName).click();

    const downloadPromise = this.page.waitForEvent('download');
    await this.exportModelButton.click();
    const download = await downloadPromise;
    expect(download).toBeDefined();

    const statusCombobox = this.page.getByRole('combobox', { name: 'enabled' });
    await statusCombobox.click();
    await this.page.getByRole('option', { name: 'ignored' }).click();
    await expect(this.page.getByRole('combobox', { name: 'ignored' })).toBeVisible();
  }

  /**
   * Imports a model from a local file.
   * @param {string} filePath - The absolute path to the model file.
   * @param {string} expectedModelName - The name of the model once imported.
   */
  async importModelByFile(filePath, expectedModelName) {
    await this.importModelButton.click();
    await this.fileImportTab.click();
    await this.fileInput.setInputFiles(filePath);
    await this.nextButton.click();
    await this.verifyModelImported(expectedModelName);
    await this.verifyImportMessagesVisible();
    await this.finishButton.click();
  }

  /**
   * Imports a model from a url.
   * @param {string} modelUrl - The URL of the model to import.
   * @param {string} expectedModelName - The name of the model once imported.
   */
  async importModelByUrl(modelUrl, expectedModelName) {
    await this.importModelButton.click();
    await this.urlImportTab.click();
    await this.urlImportInput.click();
    await this.urlImportInput.fill(modelUrl);
    await this.nextButton.click();
    await this.verifyModelImported(expectedModelName);
    await this.verifyImportMessagesVisible();
    await this.finishButton.click();
  }

  /**
   * Imports a model from CSV file e.g. tests/e2e/assets
   * @param {{models: string, components: string, relationships: string}} csvFiles
   * @param {string} expectedModelName
   */
  async importModelByCsv(csvFiles, expectedModelName) {
    await this.importModelButton.click();
    await this.csvImportTab.click();

    await this.nextButton.click();
    await this.fileInput.setInputFiles(csvFiles.models);
    await this.nextButton.click();
    await this.fileInput.setInputFiles(csvFiles.components);
    await this.nextButton.click();
    await this.fileInput.setInputFiles(csvFiles.relationships);

    await this.generateButton.click();
    await this.verifyModelImported(expectedModelName);
    await this.verifyImportMessagesVisible();
    await this.finishButton.click();
  }

  async verifyModelImported(modelName) {
    await expect(
      this.page.getByTestId(`ModelImportedSection-ModelHeader-${modelName}`),
    ).toBeVisible();
  }

  async verifyImportMessagesVisible() {
    await expect(this.page.getByTestId('ModelImportMessages-Wrapper')).toBeVisible();
  }
}
