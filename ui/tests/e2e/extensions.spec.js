import { expect, test } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';
import { ExtensionsPage } from './pages/ExtensionsPage';

// URLs used in tests
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

// Extensions Section Tests
test.describe('Extensions Section Tests', () => {
  let dashboardPage;
  let extensionsPage;

  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page);
    extensionsPage = new ExtensionsPage(page);
    
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToExtensions();
  });

  test('Verify Kanvas Snapshot using data-testid', async ({ page }) => {
    await extensionsPage.verifyKanvasSnapshotVisibility();
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    await extensionsPage.verifyPerformanceAnalysisVisibility();
  });

  test('Verify Kanvas Details', async ({ page, context }) => {
    await extensionsPage.verifyKanvasDetailsVisibility();
    
    if (await extensionsPage.isKanvasSignupBtnEnabled()) {
      await extensionsPage.verifyNewPageURL(
        context,
        () => extensionsPage.clickKanvasSignupBtn(),
        URLS.KANVAS.DOCS
      );
    }
  });

  test('Verify Meshery Docker Extension Details', async ({ page, context }) => {
    await extensionsPage.verifyDockerExtensionVisibility();
    
    await extensionsPage.verifyNewPageURL(
      context,
      () => extensionsPage.clickDockerExtensionDownloadBtn(),
      URLS.DOCKER.EXTENSION
    );
  });

  test('Verify Meshery Design Embed Details', async ({ page, context }) => {
    await extensionsPage.verifyDesignEmbedVisibility();
    
    await extensionsPage.verifyNewPageURL(
      context,
      () => extensionsPage.clickDesignEmbedLearnMoreBtn(),
      URLS.KANVAS.DESIGNER_EMBED
    );
  });

  test('Verify Meshery Catalog Section Details', async ({ page, context }) => {
    await extensionsPage.verifyCatalogSectionVisibility();
    await extensionsPage.clickCatalogToggleSwitch();
    
    await extensionsPage.verifyNewPageURL(
      context,
      () => extensionsPage.clickCatalogLink(),
      URLS.MESHERY.CATALOG
    );
  });

  test('Verify Meshery Adapter for Istio Section', async ({ page, context }) => {
    await extensionsPage.verifyNewPageURL(
      context,
      () => extensionsPage.clickAdapterDocsIstio(),
      URLS.MESHERY.ADATPER_DOCS
    );
  });
});