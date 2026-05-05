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
    await expect(this.kanvasSignupHeading).toBeVisible();
    await expect(this.kanvasSignupBtn).toBeVisible();
  }

  async hasKanvasAccess() {
    const btnText = await this.kanvasSignupBtn.textContent();
    return btnText?.trim() === 'Enabled';
  }

  async toggleCatalog() {
    await this.catalogToggleSwitch.click();
  }

  normalizeUrl(url) {
    const parsedUrl = new URL(url);
    const normalizedPath = parsedUrl.pathname.replace(/\/+$/, '') || '/';
    return `${parsedUrl.origin}${normalizedPath}${parsedUrl.search}${parsedUrl.hash}`;
  }

  async verifyNewTab(locator, expectedUrl) {
    const href = await locator.getAttribute('href');

    if (href) {
      expect(href).not.toBe('');
      expect(this.normalizeUrl(href)).toBe(this.normalizeUrl(expectedUrl));
      return;
    }

    await this.page.evaluate(() => {
      window.__mesheryOpenedUrl = null;
      // Save original window.open so it can be restored after the check.
      window.__mesheryOriginalOpen = window.open;
      window.open = (...args) => {
        window.__mesheryOpenedUrl = args[0] ?? null;
        return null;
      };
    });

    await locator.click();

    try {
      await expect
        .poll(async () => {
          const openedUrl = await this.page.evaluate(() => window.__mesheryOpenedUrl);
          return openedUrl ? this.normalizeUrl(openedUrl) : null;
        })
        .toBe(this.normalizeUrl(expectedUrl));
    } finally {
      await this.page.evaluate(() => {
        if (window.__mesheryOriginalOpen) {
          window.open = window.__mesheryOriginalOpen;
          delete window.__mesheryOriginalOpen;
        }
      });
    }
  }
}
