import { expect, test } from '@playwright/test';
import { ExtensionsPage } from './pages/ExtensionsPage';

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

test.describe('Extensions Section Tests', () => {
  let extensionsPage: ExtensionsPage;

  test.beforeEach(async ({ page }) => {
    extensionsPage = new ExtensionsPage(page);
    await extensionsPage.goto();
  });

  test('Verify Kanvas Snapshot UI elements', async () => {
    await extensionsPage.verifyKanvasSnapshotDetails();
  });

  test('Verify Performance Analysis Details', async () => {
    await extensionsPage.verifyPerformanceAnalysisDetails();
  });

  test('Verify Kanvas Details', async ({ context }) => {
    await extensionsPage.verifyKanvasSignupUI();
    await extensionsPage.verifyNewTab(context, extensionsPage.kanvasSignupBtn, URLS.KANVAS.DOCS);
  });

  test('Verify Meshery Docker Extension Details', async ({ context }) => {
    await expect(extensionsPage.dockerExtensionHeading).toBeVisible();
    await extensionsPage.verifyNewTab(
      context,
      extensionsPage.dockerExtensionDownloadBtn,
      URLS.DOCKER.EXTENSION,
    );
  });

  test('Verify Meshery Design Embed Details', async ({ context }) => {
    await expect(extensionsPage.designEmbedLearnMoreBtn).toBeVisible();
    await extensionsPage.verifyNewTab(
      context,
      extensionsPage.designEmbedLearnMoreBtn,
      URLS.KANVAS.DESIGNER_EMBED,
    );
  });

  test('Verify Meshery Catalog Section Details', async ({ context }) => {
    await expect(extensionsPage.catalogSectionHeading).toBeVisible();
    await extensionsPage.toggleCatalog();
    await extensionsPage.verifyNewTab(context, extensionsPage.catalogLink, URLS.MESHERY.CATALOG);
  });

  test('Verify Meshery Adapter for Istio Section', async ({ context }) => {
    await extensionsPage.verifyNewTab(
      context,
      extensionsPage.adapterDocsIstioLink,
      URLS.MESHERY.ADATPER_DOCS,
    );
  });
});
