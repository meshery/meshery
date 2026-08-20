export class ProviderSelectionPage {
  constructor(page) {
    this.page = page;
    this.providerDropdown = this.page.getByLabel('Select Provider');
  }

  getProviderMenuItem(providerName) {
    // The MUI menu item renders as "${providerName}" in CI and as
    // "${providerName} More" locally (the suffix comes from an additional icon
    // button inside the menu item). Using a regex that matches either form
    // avoids a Playwright strict-mode violation in both environments.
    return this.page.getByRole('menuitem', {
      name: new RegExp(`^${providerName}(?:\\s+More)?\\s*$`, 'i'),
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
