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
  let extensionsPage;

  test.beforeEach(async ({ page }) => {
    const dashboardPage = new DashboardPage(page);
    extensionsPage = new ExtensionsPage(page);
    await dashboardPage.navigateToDashboard();
    await dashboardPage.navigateToExtensions();
  });

  test('Verify Kanvas Snapshot using data-testid', async ({ page }) => {
    const section = await extensionsPage.verifyKanvasSnapshot();
    
    await expect(section.heading).toBeVisible();
    await expect(section.description).toBeVisible();
    await expect(section.enableButton).toBeVisible();
    await expect(section.enableButton).toBeEnabled();
    await expect(section.image).toBeVisible();
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    const section = await extensionsPage.getPerformanceAnalysisSection();
    
    await expect(section.heading).toBeVisible();
    await expect(section.enableButton).toBeVisible();
    await expect(section.enableButton).toBeEnabled();
  });

  test('Verify Kanvas Details', async ({ page, context }) => {
    const section = await extensionsPage.getKanvasSignupSection();
    
    await expect(section.heading).toBeVisible();
    await expect(section.button).toBeVisible();
    
    if (await section.button.isEnabled()) {
      const docsPage = await extensionsPage.clickKanvasSignupButton(context);
      await expect(docsPage).toHaveURL(URLS.KANVAS.DOCS);
      await docsPage.close();
    }
  });

  test('Verify Meshery Docker Extension Details', async ({ page, context }) => {
    const section = await extensionsPage.getDockerExtensionSection();
    await expect(section.heading).toBeVisible();
    
    const newPage = await extensionsPage.clickDockerExtensionDownload(context);
    await expect(newPage).toHaveURL(URLS.DOCKER.EXTENSION);
    await newPage.close();
  });

  test('Verify Meshery Design Embed Details', async ({ page, context }) => {
    await expect(extensionsPage.getDesignEmbedButton()).toBeVisible();
    
    const newPage = await extensionsPage.clickDesignEmbedLearnMore(context);
    await expect(newPage).toHaveURL(URLS.KANVAS.DESIGNER_EMBED);
    await newPage.close();
  });

  test('Verify Meshery Catalog Section Details', async ({ page, context }) => {
    const section = extensionsPage.getCatalogSection();
    await expect(section.heading).toBeVisible();
    
    await extensionsPage.toggleCatalog();
    const newPage = await extensionsPage.clickCatalogLink(context);
    await expect(newPage).toHaveURL(URLS.MESHERY.CATALOG);
    await newPage.close();
  });

  test('Verify Meshery Adapter for Istio Section', async ({ page, context }) => {
    const docsPage = await extensionsPage.clickAdapterDocs('istio', context);
    await expect(docsPage).toHaveURL(URLS.MESHERY.ADATPER_DOCS);
    await docsPage.close();
  });
});
