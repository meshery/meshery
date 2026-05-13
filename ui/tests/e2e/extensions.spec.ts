import { expect, test } from '@playwright/test';
import { ExtensionsPage } from './pages/ExtensionsPage';

const URLS = {
  DESIGNS: {
    DOCS: 'https://docs.meshery.io/',
    DESIGNER_EMBED: 'https://meshery.io/extensions/meshery-design-embed',
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
  let extensionsPage: ExtensionsPage;

  test.beforeEach(async ({ page }) => {
    extensionsPage = new ExtensionsPage(page);
    await extensionsPage.goto();
  });

  test('Verify Kanvas Snapshot UI elements', async () => {
    await extensionsPage.verifyKanvasSnapshotDetails();
  });

  test('Verify extension nav items use top-level layout', async () => {
    await extensionsPage.verifyExtensionNavItemsUseTopLevelLayout();
  });

  test('Verify Performance Analysis Details', async () => {
    await extensionsPage.verifyPerformanceAnalysisDetails();
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
      URLS.DESIGNS.DESIGNER_EMBED,
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
