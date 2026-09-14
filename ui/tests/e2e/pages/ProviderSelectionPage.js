export class ProviderSelectionPage {
  constructor(page) {
    this.page = page;
    this.providerDropdown = this.page.getByLabel('Select Provider');
  }

  getProviderMenuItem(providerName) {
    // The MUI menu item renders as "${providerName}" in some CI setups, and as
    // "${providerName} More information about ${providerName}" when the info
    // icon button is present. Using a regex that matches either form avoids a
    // Playwright strict-mode violation.
    return this.page.getByRole('menuitem', {
      name: new RegExp(`^${providerName}(?:\\s+More.*)?\\s*$`, 'i'),
    });
  }

  async navigateToProviderSelection() {
    await this.page.goto('/provider');
  }

  async selectProvider(providerName) {
    await this.providerDropdown.click();
    await this.getProviderMenuItem(providerName).click();
  }
}
