import { expect } from '@playwright/test';

export class ExtensionsPage {
  constructor(page) {
    this.page = page;

    // Kanvas Snapshot Locators
    this.kanvasSnapshotHeading = page.getByTestId('kanvas-snapshot-heading');
    this.kanvasSnapshotDescription = page.getByTestId('kanvas-snapshot-description');
    this.kanvasSnapshotEnableBtn = page.getByTestId('kanvas-snapshot-enable-btn');
    this.kanvasSnapshotImage = page.getByTestId('kanvas-snapshot-image');

    // Performance Analysis Locators
    this.performanceAnalysisHeading = page.getByTestId('performance-analysis-heading');
    this.performanceAnalysisEnableBtn = page.getByTestId('performance-analysis-enable-btn');

    // Kanvas Details Locators
    this.kanvasSignupHeading = page.getByTestId('kanvas-signup-heading');
    this.kanvasSignupBtn = page.getByTestId('kanvas-signup-btn');

    // Docker Extension Locators
    this.dockerExtensionHeading = page.getByTestId('docker-extension-heading');
    this.dockerExtensionDownloadBtn = page.getByTestId('docker-extension-download-btn');

    // Design Embed Locators
    this.designEmbedLearnMoreBtn = page.getByTestId('design-embed-learn-more-btn');

    // Catalog Section Locators
    this.catalogSectionHeading = page.getByTestId('catalog-section-heading');
    this.catalogToggleSwitch = page.getByTestId('catalog-toggle-switch');
    this.catalogLink = page.locator('a[href="https://meshery.io/catalog"]');

    // Adapter Locators
    this.adapterDocsIstio = page.getByTestId('adapter-docs-istio');
  }

  // Kanvas Snapshot Methods
  async verifyKanvasSnapshotVisibility() {
    await expect(this.kanvasSnapshotHeading).toBeVisible();
    await expect(this.kanvasSnapshotDescription).toBeVisible();
    await expect(this.kanvasSnapshotEnableBtn).toBeVisible();
    await expect(this.kanvasSnapshotEnableBtn).toBeEnabled();
    await expect(this.kanvasSnapshotImage).toBeVisible();
  }

  async clickKanvasSnapshotEnableBtn() {
    await this.kanvasSnapshotEnableBtn.click();
  }

  // Performance Analysis Methods
  async verifyPerformanceAnalysisVisibility() {
    await expect(this.performanceAnalysisHeading).toBeVisible();
    await expect(this.performanceAnalysisEnableBtn).toBeVisible();
    await expect(this.performanceAnalysisEnableBtn).toBeEnabled();
  }

  async clickPerformanceAnalysisEnableBtn() {
    await this.performanceAnalysisEnableBtn.click();
  }

  // Kanvas Details Methods
  async verifyKanvasDetailsVisibility() {
    await expect(this.kanvasSignupHeading).toBeVisible();
    await expect(this.kanvasSignupBtn).toBeVisible();
  }

  async clickKanvasSignupBtn() {
    return await this.kanvasSignupBtn.click();
  }

  async isKanvasSignupBtnEnabled() {
    return await this.kanvasSignupBtn.isEnabled();
  }

  // Docker Extension Methods
  async verifyDockerExtensionVisibility() {
    await expect(this.dockerExtensionHeading).toBeVisible();
  }

  async clickDockerExtensionDownloadBtn() {
    return await this.dockerExtensionDownloadBtn.click();
  }

  // Design Embed Methods
  async verifyDesignEmbedVisibility() {
    await expect(this.designEmbedLearnMoreBtn).toBeVisible();
  }

  async clickDesignEmbedLearnMoreBtn() {
    return await this.designEmbedLearnMoreBtn.click();
  }

  // Catalog Section Methods
  async verifyCatalogSectionVisibility() {
    await expect(this.catalogSectionHeading).toBeVisible();
  }

  async clickCatalogToggleSwitch() {
    await this.catalogToggleSwitch.click();
  }

  async clickCatalogLink() {
    return await this.catalogLink.click();
  }

  // Adapter Methods
  async clickAdapterDocsIstio() {
    return await this.adapterDocsIstio.click();
  }

  // Helper method to handle new page opening
  async clickAndWaitForNewPage(context, clickAction) {
    const [newPage] = await Promise.all([context.waitForEvent('page'), clickAction()]);
    return newPage;
  }

  // Verification helper for URL checking
  async verifyNewPageURL(context, clickAction, expectedURL) {
    const newPage = await this.clickAndWaitForNewPage(context, clickAction);
    await expect(newPage).toHaveURL(expectedURL);
    await newPage.close();
  }
}
