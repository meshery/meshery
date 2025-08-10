import { DashboardPage } from './DashboardPage';
import { DesignPage } from './DesignPage/DesignPage';

export class DesignConfiguratorPage {
  constructor(page) {
    this.page = page;
    this.DashboardPage = new DashboardPage(page);
    this.DesignPage = new DesignPage(page);

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

  async navigateTo() {
    await this.DashboardPage.navigateToDashboard();
    await this.DashboardPage.navigateToDesigns();
    await this.DesignPage.navigateToDesignConfigurator();
  }

  async saveDesign() {
    await this.saveButton.click();
  }

  async updateDesign() {
    await this.updateButton.click();
  }
}
