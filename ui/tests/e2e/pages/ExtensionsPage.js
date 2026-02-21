import { expect } from '@playwright/test';
import { DashboardPage } from './DashboardPage';

export class ExtensionsPage {
  constructor(page) {
    this.page = page;

    this.kanvasSnapshotHeading = page.getByTestId('kanvas-snapshot-heading');
    this.kanvasSnapshotDescription = page.getByTestId('kanvas-snapshot-description');
    this.kanvasSnapshotEnableBtn = page.getByTestId('kanvas-snapshot-enable-btn');
    this.kanvasSnapshotImage = page.getByTestId('kanvas-snapshot-image');

    this.performanceHeading = page.getByTestId('performance-analysis-heading');
    this.performanceEnableBtn = page.getByTestId('performance-analysis-enable-btn');

    this.kanvasSignupHeading = page.getByTestId('kanvas-signup-heading');
    this.kanvasSignupBtn = page.getByTestId('kanvas-signup-btn');

    this.dockerExtensionHeading = page.getByTestId('docker-extension-heading');
    this.dockerExtensionDownloadBtn = page.getByTestId('docker-extension-download-btn');

    this.designEmbedLearnMoreBtn = page.getByTestId('design-embed-learn-more-btn');

    this.catalogSectionHeading = page.getByTestId('catalog-section-heading');
    this.catalogToggleSwitch = page.getByTestId('catalog-toggle-switch');
    this.catalogLink = page.getByTestId('catalog-link');

    this.adapterDocsIstioLink = page.getByTestId('adapter-docs-istio');
  }

  async goto() {
    const dashboardPage = new DashboardPage(this.page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToExtensions();
  }

  async verifyKanvasSnapshotDetails() {
    await expect(this.kanvasSnapshotHeading).toBeVisible();
    await expect(this.kanvasSnapshotDescription).toBeVisible();
    await expect(this.kanvasSnapshotEnableBtn).toBeVisible();
    await expect(this.kanvasSnapshotEnableBtn).toBeEnabled();
    await expect(this.kanvasSnapshotImage).toBeVisible();
  }

  async verifyPerformanceAnalysisDetails() {
    await expect(this.performanceHeading).toBeVisible();
    await expect(this.performanceEnableBtn).toBeVisible();
    await expect(this.performanceEnableBtn).toBeEnabled();
  }

  async verifyKanvasSignupUI() {
    await expect(this.kanvasSignupHeading).toBeVisible({ timeout: 10000 });
    await expect(this.kanvasSignupBtn).toBeVisible({ timeout: 10000 });
    await expect(this.kanvasSignupBtn).toBeEnabled({ timeout: 10000 });
  }

  async toggleCatalog() {
    await this.catalogToggleSwitch.click();
  }

  async verifyNewTab(context, locator, expectedUrl) {
    // Wait for the element to be visible and clickable before proceeding
    await expect(locator).toBeVisible({ timeout: 10000 });
    await expect(locator).toBeEnabled({ timeout: 10000 });

    const [newPage] = await Promise.all([
      context.waitForEvent('page', { timeout: 30000 }),
      locator.click(),
    ]);
    await newPage.waitForLoadState('domcontentloaded');
    await expect(newPage).toHaveURL(expectedUrl, { timeout: 15000 });
    await newPage.close();
  }
}
