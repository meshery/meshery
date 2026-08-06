export class DesignConfiguratorPage {
  constructor(page) {
    this.page = page;

    this.appBar = page.getByTestId('design-configurator-app-bar');
    this.codeEditor = page.getByTestId('design-configurator-code-editor');
    this.categoryModelSelector = page.getByTestId('category-model-selector');
    this.databaseCategory = page.getByTestId('Database');
    this.backToCategories = page.getByTestId('back-to-categories');
    this.clearSelector = page.getByTestId('clear-category-model-selector');
    this.modelAzure = page.getByTestId('Database-azure-db-for-mysql');
    this.modelContainer = page.getByTestId('model-component-list');
    this.saveButton = page.getByTestId('design-configurator-save-design-btn');
    this.updateButton = page.getByTestId('design-configurator-update-design-btn');
  }

  async navigateTo(designId) {
    const query = designId ? `?design_id=${designId}` : '';
    await this.page.goto(`/configuration/designs/configurator${query}`);
  }

  async saveDesign() {
    await this.saveButton.click();
  }

  async updateDesign() {
    await this.updateButton.click();
  }
}
