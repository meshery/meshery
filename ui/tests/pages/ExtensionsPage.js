class ExtensionsPage {
  constructor(page) {
    this.page = page;
    this.addExtensionButton = page.locator('button#add-extension');
    this.extensionNameInput = page.locator('input#extension-name');
    this.submitButton = page.locator('button#submit-extension');
    this.extensionList = page.locator('.extension-list');
  }

  async goto() {
    await this.page.goto('/extensions');
  }

  async addExtension(name) {
    await this.addExtensionButton.click();
    await this.extensionNameInput.fill(name);
    await this.submitButton.click();
  }

  async getExtensionList() {
    return await this.extensionList.allTextContents();
  }
}

module.exports = { ExtensionsPage };
