export class ExtensionsPage {
  constructor(page) {
    this.page = page;
  }

  // Kanvas Snapshot
  async verifyKanvasSnapshot() {
    await this.page.getByTestId('kanvas-snapshot-heading').waitFor({ state: 'visible' });
    return {
      heading: this.page.getByTestId('kanvas-snapshot-heading'),
      description: this.page.getByTestId('kanvas-snapshot-description'),
      enableButton: this.page.getByTestId('kanvas-snapshot-enable-btn'),
      image: this.page.getByTestId('kanvas-snapshot-image'),
    };
  }

  // Performance Analysis
  async getPerformanceAnalysisSection() {
    return {
      heading: this.page.getByTestId('performance-analysis-heading'),
      enableButton: this.page.getByTestId('performance-analysis-enable-btn'),
    };
  }

  // Kanvas Details
  async getKanvasSignupSection() {
    return {
      heading: this.page.getByTestId('kanvas-signup-heading'),
      button: this.page.getByTestId('kanvas-signup-btn'),
    };
  }

  async clickKanvasSignupButton(context) {
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId('kanvas-signup-btn').click(),
    ]);
    return docsPage;
  }

  // Docker Extension
  async getDockerExtensionSection() {
    return {
      heading: this.page.getByTestId('docker-extension-heading'),
      downloadButton: this.page.getByTestId('docker-extension-download-btn'),
    };
  }

  async clickDockerExtensionDownload(context) {
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId('docker-extension-download-btn').click(),
    ]);
    return newPage;
  }

  // Design Embed
  async clickDesignEmbedLearnMore(context) {
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId('design-embed-learn-more-btn').click(),
    ]);
    return newPage;
  }

  getDesignEmbedButton() {
    return this.page.getByTestId('design-embed-learn-more-btn');
  }

  // Catalog Section
  getCatalogSection() {
    return {
      heading: this.page.getByTestId('catalog-section-heading'),
      toggleSwitch: this.page.getByTestId('catalog-toggle-switch'),
      link: this.page.locator('a[href="https://meshery.io/catalog"]'),
    };
  }

  async toggleCatalog() {
    await this.page.getByTestId('catalog-toggle-switch').click();
  }

  async clickCatalogLink(context) {
    const catalogLink = this.page.locator('a[href="https://meshery.io/catalog"]');
    const [newPage] = await Promise.all([context.waitForEvent('page'), catalogLink.click()]);
    return newPage;
  }

  // Adapters
  async clickAdapterDocs(adapterName, context) {
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId(`adapter-docs-${adapterName}`).click(),
    ]);
    return docsPage;
  }
}
