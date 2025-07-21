import { DashboardPage } from '../e2e/pages/DashboardPage.js';

const SELECTORS = {
  kanvasSnapshot: {
    heading: 'kanvas-snapshot-heading',
    description: 'kanvas-snapshot-description',
    enableBtn: 'kanvas-snapshot-enable-btn',
    image: 'kanvas-snapshot-image',
  },
  performanceAnalysis: {
    heading: 'performance-analysis-heading',
    enableBtn: 'performance-analysis-enable-btn',
  },
  kanvasDetails: {
    heading: 'kanvas-signup-heading',
    button: 'kanvas-signup-btn',
  },
  dockerExtension: {
    heading: 'docker-extension-heading',
    downloadBtn: 'docker-extension-download-btn',
  },
  designEmbed: {
    learnMoreBtn: 'design-embed-learn-more-btn',
  },
  catalog: {
    heading: 'catalog-section-heading',
    toggle: 'catalog-toggle-switch',
    link: 'a[href="https://meshery.io/catalog"]',
  },
  adapterIstio: 'adapter-docs-istio',
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

  // Kanvas Snapshot
  getKanvasSnapshotElements() {
    return {
      heading: this.page.getByTestId(SELECTORS.kanvasSnapshot.heading),
      description: this.page.getByTestId(SELECTORS.kanvasSnapshot.description),
      enableBtn: this.page.getByTestId(SELECTORS.kanvasSnapshot.enableBtn),
      image: this.page.getByTestId(SELECTORS.kanvasSnapshot.image),
    };
  }

  // Performance Analysis
  getPerformanceAnalysisElements() {
    return {
      heading: this.page.getByTestId(SELECTORS.performanceAnalysis.heading),
      enableBtn: this.page.getByTestId(SELECTORS.performanceAnalysis.enableBtn),
    };
  }

  // Kanvas Details
  getKanvasDetailsElements() {
    return {
      heading: this.page.getByTestId(SELECTORS.kanvasDetails.heading),
      button: this.page.getByTestId(SELECTORS.kanvasDetails.button),
    };
  }

  async openKanvasDocs(context) {
    const button = this.page.getByTestId(SELECTORS.kanvasDetails.button);
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      button.click(),
    ]);
    return docsPage;
  }

  // Docker Extension
  async openDockerExtension(context) {
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId(SELECTORS.dockerExtension.downloadBtn).click(),
    ]);
    return newPage;
  }

  getDockerExtensionHeading() {
    return this.page.getByTestId(SELECTORS.dockerExtension.heading);
  }

  // Design Embed
  async openDesignEmbed(context) {
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId(SELECTORS.designEmbed.learnMoreBtn).click(),
    ]);
    return newPage;
  }

  getDesignEmbedButton() {
    return this.page.getByTestId(SELECTORS.designEmbed.learnMoreBtn);
  }

  // Catalog Section
  getCatalogHeading() {
    return this.page.getByTestId(SELECTORS.catalog.heading);
  }

  async toggleCatalogSection() {
    await this.page.getByTestId(SELECTORS.catalog.toggle).click();
  }

  async openCatalogLink(context) {
    const link = this.page.locator(SELECTORS.catalog.link);
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      link.click(),
    ]);
    return newPage;
  }

  // Adapter Istio
  async openAdapterIstioDocs(context) {
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      this.page.getByTestId(SELECTORS.adapterIstio).click(),
    ]);
    return docsPage;
  }
}
