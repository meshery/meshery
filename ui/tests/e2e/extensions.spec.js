import { test } from '@playwright/test';
import { DashboardPage } from './pages/DashboardPage';
import { ExtensionPage } from './pages/ExtensionPage';

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
  test.beforeEach(async ({ page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToExtensions();
  });

  test('Verify Kanvas Snapshot using data-testid', async ({ page }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyKanvasSnapshotVisible();
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyPerformanceAnalysisDetails();
  });

  test('Verify Kanvas Details', async ({ page, context }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyKanvasDetails(context, URLS.KANVAS.DOCS);
  });

  test('Verify Meshery Docker Extension Details', async ({ page, context }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyMesheryDockerExtensionDetails(context, URLS.DOCKER.EXTENSION);
  });

  test('Verify Meshery Design Embed Details', async ({ page, context }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyMesheryDesignEmbedDetails(context, URLS.KANVAS.DESIGNER_EMBED);
  });

  test('Verify Meshery Catalog Section Details', async ({ page, context }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyMesheryCatalogSectionDetails(context, URLS.MESHERY.CATALOG);
  });

  test('Verify Meshery Adapter for Istio Section', async ({ page, context }) => {
    const extensionPage = new ExtensionPage(page);
    await extensionPage.verifyMesheryAdapterIstioSection(context, URLS.MESHERY.ADATPER_DOCS);
  });
});
