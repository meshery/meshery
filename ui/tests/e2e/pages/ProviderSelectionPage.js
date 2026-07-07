export class ProviderSelectionPage {
  constructor(page) {
    this.page = page;
    this.providerDropdown = this.page.getByLabel('Select Provider');
  }

  getProviderMenuItem(providerName) {
    // Use a regex with word boundary to avoid "Local" matching "localhost:9876"
    return this.page.getByRole('menuitem', { name: new RegExp('^' + providerName + '\\b', 'i') });
  }

  async navigateToProviderSelection() {
    await this.page.goto('/provider');
  }

  async selectProvider(providerName) {
    await this.providerDropdown.click();
    await this.getProviderMenuItem(providerName).click();
  }
}
