import { expect } from '@playwright/test';

export class ModelsPage {
  constructor(page) {
    this.page = page;
    
    // Navigation elements
    this.createModelButton = page.getByTestId('TabBar-Button-CreateModel');
    this.importModelButton = page.getByTestId('TabBar-Button-ImportModel');
    this.searchIcon = page.getByTestId('search-icon');
    this.searchInput = page.locator('#searchClick');
    
    // Model creation form elements
    this.modelNameInput = page.locator('#model-name');
    this.modelDisplayNameInput = page.locator('#model-display-name');
    this.modelUrlInput = page.locator('#model-url');
    
    // Stepper navigation
    this.nextButton = page.getByTestId('UrlStepper-Button-Next');
    this.generateButton = page.getByTestId('UrlStepper-Button-Generate');
    this.finishButton = page.getByTestId('UrlStepper-Button-Finish');
    
    // Category and subcategory selects
    this.categorySelect = page.getByTestId('UrlStepper-Select-Category');
    this.subcategorySelect = page.getByTestId('UrlStepper-Select-Subcategory');
    
    // Logo and styling elements
    this.logoDarkThemeSelect = page.getByTestId('UrlStepper-Select-Logo-Dark-Theme');
    this.logoLightThemeSelect = page.getByTestId('UrlStepper-Select-Logo-Light-Theme');
    this.primaryColorSelect = page.getByTestId('UrlStepper-Select-Primary-Color');
    this.secondaryColorSelect = page.getByTestId('UrlStepper-Select-Secondary-Color');
    this.shapeSelect = page.getByTestId('UrlStepper-Select-Shape');
    
    // Source selection
    this.githubSourceRadio = page.getByTestId('UrlStepper-Select-Source-GitHub');
    this.visualAnnotationCheckbox = page.getByTestId('UrlStepper-Visual-Annotation-Checkbox');
    
    // Import sections
    this.fileImportHeading = page.getByRole('heading', { name: 'File Import' });
    this.urlImportHeading = page.getByRole('heading', { name: 'URL Import' });
    this.csvImportHeading = page.getByRole('heading', { name: 'CSV Import' });
    this.urlImportTextbox = page.getByRole('textbox', { name: 'URL' });
    this.fileInput = page.locator('input[type="file"]');
    
    // Results and messages
    this.modelImportMessagesWrapper = page.getByTestId('ModelImportMessages-Wrapper');
    this.exportModelButton = page.getByTestId('export-model-button');
    this.enabledCombobox = page.getByRole('combobox', { name: 'enabled' });
    this.ignoredOption = page.getByRole('option', { name: 'ignored' });
    this.ignoredCombobox = page.getByRole('combobox', { name: 'ignored' });
    
    // Generic buttons for import flows
    this.nextButtonGeneric = page.getByRole('button', { name: 'Next' });
    this.finishButtonGeneric = page.getByRole('button', { name: 'Finish' });
    this.generateButtonGeneric = page.getByRole('button', { name: 'Generate' });
  }

  async clickCreateModel() {
    await this.createModelButton.click();
  }

  async clickImportModel() {
    await this.importModelButton.click();
  }

  async fillModelDetails(modelName, modelDisplayName) {
    await this.modelNameInput.fill(modelName);
    await this.modelDisplayNameInput.fill(modelDisplayName);
  }

  async proceedThroughCategoryStep() {
    await this.nextButton.click();
    await expect(this.categorySelect).toBeVisible();
    await expect(this.subcategorySelect).toBeVisible();
  }

  async proceedThroughStylingStep() {
    await this.nextButton.click();
    await expect(this.logoDarkThemeSelect).toBeVisible();
    await expect(this.logoLightThemeSelect).toBeVisible();
    await expect(this.primaryColorSelect).toBeVisible();
    await expect(this.secondaryColorSelect).toBeVisible();
    await expect(this.shapeSelect).toBeVisible();
  }

  async configureGitHubSource(modelUrl) {
    await this.nextButton.click();
    await this.githubSourceRadio.check();
    await this.modelUrlInput.fill(modelUrl);
  }

  async enableVisualAnnotationAndGenerate() {
    await this.nextButton.click();
    await this.visualAnnotationCheckbox.check();
    await this.nextButton.click();
    await this.generateButton.click();
  }

  async verifyModelImported(modelName) {
    // Fixed template literal syntax
    await expect(
      this.page.getByTestId(`ModelImportedSection-ModelHeader-${modelName}`)
    ).toBeVisible();
    await expect(this.modelImportMessagesWrapper).toBeVisible();
  }

  async finishModelCreation() {
    await this.finishButton.click();
  }

  async searchAndSelectModel(modelDisplayName) {
    await this.searchIcon.click();
    await this.searchInput.click();
    await this.searchInput.fill(modelDisplayName);
    await this.page.getByText(modelDisplayName).click();
  }

  async exportModel() {
    const downloadPromise = this.page.waitForEvent('download');
    await this.exportModelButton.click();
    const download = await downloadPromise;
    expect(download).toBeDefined();
    return download;
  }

  async changeModelStatus() {
    await this.enabledCombobox.click();
    await this.ignoredOption.click();
    await expect(this.ignoredCombobox).toBeVisible();
  }

  async selectFileImport() {
    await this.fileImportHeading.click();
  }

  async selectUrlImport() {
    await this.urlImportHeading.click();
  }

  async selectCsvImport() {
    await this.csvImportHeading.click();
  }

  async uploadFile(filePath) {
    await this.page.setInputFiles('input[type="file"]', filePath);
  }

  async fillUrlImport(url) {
    await this.urlImportTextbox.click();
    await this.urlImportTextbox.fill(url);
  }

  async proceedWithImport() {
    await this.nextButtonGeneric.click();
  }

  async finishImport() {
    await this.finishButtonGeneric.click();
  }

  async generateFromCsv() {
    await this.generateButtonGeneric.click();
  }

  // Complete workflow methods
  async createModelWorkflow(modelData) {
    await this.clickCreateModel();
    await this.fillModelDetails(modelData.MODEL_NAME, modelData.MODEL_DISPLAY_NAME);
    await this.proceedThroughCategoryStep();
    await this.proceedThroughStylingStep();
    await this.configureGitHubSource(modelData.MODEL_URL);
    await this.enableVisualAnnotationAndGenerate();
    await this.verifyModelImported(modelData.MODEL_NAME);
    await this.finishModelCreation();
  }

  async importModelViaFile(modelData) {
    await this.clickImportModel();
    await this.selectFileImport();
    await this.uploadFile(modelData.MODEL_FILE_IMPORT);
    await this.proceedWithImport();
    await this.verifyModelImported(modelData.MODEL_NAME);
    await this.finishImport();
  }

  async importModelViaUrl(modelData) {
    await this.clickImportModel();
    await this.selectUrlImport();
    await this.fillUrlImport(modelData.MODEL_URL_IMPORT);
    await this.proceedWithImport();
    await this.verifyModelImported(modelData.MODEL_NAME);
    await this.finishImport();
  }

  async importModelViaCsv(csvData) {
    await this.clickImportModel();
    await this.selectCsvImport();
    await this.proceedWithImport();
    
    // Upload models CSV
    await this.uploadFile(csvData.Models);
    await this.proceedWithImport();
    
    // Upload components CSV
    await this.uploadFile(csvData.Components);
    await this.proceedWithImport();
    
    // Upload relationships CSV
    await this.uploadFile(csvData.Relationships);
    await this.generateFromCsv();
    
    await this.verifyModelImported(csvData.Model_Name);
    await this.finishImport();
  }

  async searchExportAndChangeStatus(modelDisplayName) {
    await this.searchAndSelectModel(modelDisplayName);
    const download = await this.exportModel();
    await this.changeModelStatus();
    return download;
  }
}