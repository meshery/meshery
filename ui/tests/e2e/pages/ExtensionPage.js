const { expect } = require('@playwright/test');

class ExtensionPage {
  constructor(page) {
    this.page = page;
    // Locators for Kanvas Snapshot section
    this.kanvasSnapshotHeading = page.getByTestId('kanvas-snapshot-heading');
    this.kanvasSnapshotDescription = page.getByTestId('kanvas-snapshot-description');
    this.kanvasSnapshotEnableButton = page.getByTestId('kanvas-snapshot-enable-btn');
    this.kanvasSnapshotImage = page.getByTestId('kanvas-snapshot-image');

    // Locators for Performance Analysis section
    this.performanceAnalysisHeading = page.getByTestId('performance-analysis-heading');
    this.performanceAnalysisEnableButton = page.getByTestId('performance-analysis-enable-btn');

    // Locators for Kanvas Details section
    this.kanvasSignupHeading = page.getByTestId('kanvas-signup-heading');
    this.kanvasSignupButton = page.getByTestId('kanvas-signup-btn');

    // Locators for Meshery Docker Extension section
    this.dockerExtensionHeading = page.getByTestId('docker-extension-heading');
    this.dockerExtensionDownloadButton = page.getByTestId('docker-extension-download-btn');

    // Locators for Meshery Design Embed section
    this.designEmbedLearnMoreButton = page.getByTestId('design-embed-learn-more-btn');

    // Locators for Meshery Catalog section
    this.catalogSectionHeading = page.getByTestId('catalog-section-heading');
    this.catalogToggleSwitch = page.getByTestId('catalog-toggle-switch');
    this.catalogLink = page.locator('a[href="https://meshery.io/catalog"]');

    // Locators for Meshery Adapter for Istio section
    this.adapterDocsIstioButton = page.getByTestId('adapter-docs-istio');
  }

  // Navigation method (assuming it's called after navigating to dashboard)
  async navigateToExtensions() {
    // If navigation to Extensions section is needed, add logic here
    // For now, assuming DashboardPage.navigateToExtensions() handles it
  }

  // Kanvas Snapshot section methods
  async verifyKanvasSnapshotVisible() {
    await expect(this.kanvasSnapshotHeading).toBeVisible();
    await expect(this.kanvasSnapshotDescription).toBeVisible();
    await expect(this.kanvasSnapshotEnableButton).toBeVisible();
    await expect(this.kanvasSnapshotEnableButton).toBeEnabled();
    await expect(this.kanvasSnapshotImage).toBeVisible();
  }

  // Performance Analysis section methods
  async verifyPerformanceAnalysisDetails() {
    await expect(this.performanceAnalysisHeading).toBeVisible();
    await expect(this.performanceAnalysisEnableButton).toBeVisible();
    await expect(this.performanceAnalysisEnableButton).toBeEnabled();
  }

  // Kanvas Details section methods
  async verifyKanvasDetails(context, expectedUrl = 'https://docs.layer5.io/kanvas/') {
    await expect(this.kanvasSignupHeading).toBeVisible();
    await expect(this.kanvasSignupButton).toBeVisible();
    if (await this.kanvasSignupButton.isEnabled()) {
      const [docsPage] = await Promise.all([
        context.waitForEvent('page'),
        this.kanvasSignupButton.click(),
      ]);
      await expect(docsPage).toHaveURL(expectedUrl);
      await docsPage.close();
    }
  }

  // Meshery Docker Extension section methods
  async verifyMesheryDockerExtensionDetails(context, expectedUrl = 'https://hub.docker.com/extensions/meshery/docker-extension-meshery') {
    await expect(this.dockerExtensionHeading).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.dockerExtensionDownloadButton.click(),
    ]);
    await expect(newPage).toHaveURL(expectedUrl);
    await newPage.close();
  }

  // Meshery Design Embed section methods
  async verifyMesheryDesignEmbedDetails(context, expectedUrl = 'https://docs.layer5.io/kanvas/designer/embedding-designs/') {
    await expect(this.designEmbedLearnMoreButton).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.designEmbedLearnMoreButton.click(),
    ]);
    await expect(newPage).toHaveURL(expectedUrl);
    await newPage.close();
  }

  // Meshery Catalog section methods
  async verifyMesheryCatalogSectionDetails(context, expectedUrl = 'https://meshery.io/catalog') {
    await expect(this.catalogSectionHeading).toBeVisible();
    await this.catalogToggleSwitch.click();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.catalogLink.click(),
    ]);
    await expect(newPage).toHaveURL(expectedUrl);
    await newPage.close();
  }

  // Meshery Adapter for Istio section methods
  async verifyMesheryAdapterIstioSection(context, expectedUrl = 'https://docs.meshery.io/concepts/architecture/adapters') {
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      this.adapterDocsIstioButton.click(),
    ]);
    await expect(docsPage).toHaveURL(expectedUrl);
    await docsPage.close();
  }
}

module.exports = { ExtensionPage };