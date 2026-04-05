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
    ADAPTER_DOCS: 'https://docs.meshery.io/concepts/architecture/adapters',
  },
};

test.describe('Extensions Section Tests', () => {
  let extensionsPage;

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

  test('Verify Kanvas Details', async () => {
    await extensionsPage.verifyKanvasSignupUI();
    const hasAccess = await extensionsPage.hasKanvasAccess();
    if (hasAccess) {
      await expect(extensionsPage.kanvasSignupBtn).toBeDisabled();
    } else {
      await expect(extensionsPage.kanvasSignupBtn).toBeEnabled();
      await extensionsPage.verifyNewTab(extensionsPage.kanvasSignupBtn, URLS.KANVAS.DOCS);
    }
  });

  test('Verify Meshery Docker Extension Details', async () => {
    await expect(extensionsPage.dockerExtensionHeading).toBeVisible();
    await extensionsPage.verifyNewTab(
      extensionsPage.dockerExtensionDownloadBtn,
      URLS.DOCKER.EXTENSION,
    );
  });

  test('Verify Meshery Design Embed Details', async () => {
    await expect(extensionsPage.designEmbedLearnMoreBtn).toBeVisible();
    await extensionsPage.verifyNewTab(
      extensionsPage.designEmbedLearnMoreBtn,
      URLS.KANVAS.DESIGNER_EMBED,
    );
  });

  test('Verify Meshery Catalog Section Details', async () => {
    await expect(extensionsPage.catalogSectionHeading).toBeVisible();
    await extensionsPage.toggleCatalog();
    await extensionsPage.verifyNewTab(extensionsPage.catalogLink, URLS.MESHERY.CATALOG);
  });

  test('Verify Meshery Adapter for Istio Section', async () => {
    await extensionsPage.verifyNewTab(
      extensionsPage.adapterDocsIstioLink,
      URLS.MESHERY.ADAPTER_DOCS,
    );
  });
});
