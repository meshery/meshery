import { expect } from '@playwright/test';
import { DashboardPage } from '../e2e/pages/DashboardPage.js';
const URLS = {
  KANVAS: {
    DOCS: 'https://docs.layer5.io/kanvas/',
    DESIGNER_EMBED: 'https://docs.layer5.io/kanvas/designer/embedding-designs/',
  },
  DOCKER: {
    EXTENSION: 'https://hub.docker.com/extensions/meshery/docker-extension-meshery',
  },
  MESHERY: {
    CATALOG: 'https://meshery.io/catalog',
    ADATPER_DOCS: 'https://docs.meshery.io/concepts/architecture/adapters',
  },
};
export class ExtensionsPage {
  constructor(page) {
    this.page = page;
  }

  async navigate() {
    const dashboardPage = new DashboardPage(this.page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToExtensions();
  }

  async verifyKanvasSnapshot() {
    await expect(this.page.getByTestId('kanvas-snapshot-heading')).toBeVisible();
    await expect(this.page.getByTestId('kanvas-snapshot-description')).toBeVisible();

    const enableButton = this.page.getByTestId('kanvas-snapshot-enable-btn');
    await expect(enableButton).toBeVisible();
    await expect(enableButton).toBeEnabled();
    await expect(this.page.getByTestId('kanvas-snapshot-image')).toBeVisible();
  }

  async verifyPerformanceAnalysis() {
    await expect(this.page.getByTestId('performance-analysis-heading')).toBeVisible();
    const enableButton = this.page.getByTestId('performance-analysis-enable-btn');
    await expect(enableButton).toBeVisible();
    await expect(enableButton).toBeEnabled();
  }

  async verifyKanvasDetails(context) {
    await expect(this.page.getByTestId('kanvas-signup-heading')).toBeVisible();
    const button = this.page.getByTestId('kanvas-signup-btn');
    await expect(button).toBeVisible();
    if (await button.isEnabled()) {
      const [docsPage] = await Promise.all([context.waitForEvent('page'), button.click()]);
      await expect(docsPage).toHaveURL(URLS.KANVAS.DOCS);
      await docsPage.close();
    }
  }

  async verifyDockerExtension(context) {
    await expect(this.page.getByTestId('docker-extension-heading')).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId('docker-extension-download-btn').click(),
    ]);
    await expect(newPage).toHaveURL(URLS.DOCKER.EXTENSION);
    await newPage.close();
  }

  async verifyDesignEmbed(context) {
    await expect(this.page.getByTestId('design-embed-learn-more-btn')).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId('design-embed-learn-more-btn').click(),
    ]);
    await expect(newPage).toHaveURL(URLS.KANVAS.DESIGNER_EMBED);
    await newPage.close();
  }

  async verifyCatalogSection(context) {
    await expect(this.page.getByTestId('catalog-section-heading')).toBeVisible();
    const toggle = this.page.getByTestId('catalog-toggle-switch');
    await toggle.click();
    const catalogLink = this.page.locator('a[href="https://meshery.io/catalog"]');
    const [newPage] = await Promise.all([context.waitForEvent('page'), catalogLink.click()]);
    await expect(newPage).toHaveURL(URLS.MESHERY.CATALOG);
    await newPage.close();
  }

  async verifyAdapterIstio(context) {
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId('adapter-docs-istio').click(),
    ]);
    await expect(docsPage).toHaveURL(URLS.MESHERY.ADATPER_DOCS);
    await docsPage.close();
  }
}
