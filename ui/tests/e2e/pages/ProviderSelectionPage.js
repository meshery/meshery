export class ProviderSelectionPage {
  constructor(page) {
    this.page = page;
    this.providerDropdown = this.page.getByLabel('Select Provider');
  }

  getProviderMenuItem(providerName) {
    return this.page.getByRole('menuitem', { name: providerName });
  }

  async navigateToProviderSelection() {
    await this.page.goto('/provider');
  }

  async selectProvider(providerName) {
    await this.providerDropdown.click();
    await this.getProviderMenuItem(providerName).click();
  }
}
