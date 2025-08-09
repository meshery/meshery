import { DashboardPage } from './DashboardPage';
import { DesignPage } from './DesignPage/DesignPage';

export class DesignConfiguratorPage {
  constructor(page) {
    this.page = page;
    this.DashboardPage = new DashboardPage(page);
    this.DesignPage = new DesignPage(page);

    this.appBar = this.page.getByTestId('design-configurator-app-bar');
    this.codeEditor = this.page.getByTestId('design-configurator-code-editor');
    this.categorySelector = this.page.getByTestId('category-selector');
    this.databaseCategory = this.page.getByTestId('Database');
    this.modelSelector = this.page.getByTestId('model-selector');
    this.modelAzure = this.page.getByTestId('azure-db-for-mysql');
    this.modelContainer = this.page.getByTestId('model-component-list');
    this.saveButton = this.page.getByTestId('design-configurator-save-design-btn');
    this.updateButton = this.page.getByTestId('design-configurator-update-design-btn');
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
