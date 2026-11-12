export class DesignConfiguratorPage {
  constructor(page) {
    this.page = page;

    this.appBar = page.getByTestId('design-configurator-app-bar');
    this.codeEditor = page.getByTestId('design-configurator-code-editor');
    this.categorySelector = page.getByTestId('category-selector');
    this.databaseCategory = page.getByTestId('Database');
    this.modelSelector = page.getByTestId('model-selector');
    this.modelAzure = page.getByTestId('azure-db-for-mysql');
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
