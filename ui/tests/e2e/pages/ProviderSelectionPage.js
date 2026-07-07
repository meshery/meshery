export class ProviderSelectionPage {
  constructor(page) {
    this.page = page;
    this.providerDropdown = this.page.getByLabel('Select Provider');
  }

  getProviderMenuItem(providerName) {
    // Escape providerName to avoid regex injection, then use word boundary.
    // We cannot use exact: true because the accessible name includes tooltip text.
    const escapedName = providerName.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    return this.page.getByRole('menuitem', { name: new RegExp('^' + escapedName + '\\b', 'i') });
  }

  async navigateToProviderSelection() {
    await this.page.goto('/provider');
  }

  async selectProvider(providerName) {
    await this.providerDropdown.click();
    await this.getProviderMenuItem(providerName).click();
  }
}
