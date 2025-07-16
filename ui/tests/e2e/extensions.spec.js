import { test } from '@playwright/test';
import { ExtensionsPage } from './pages/ExtensionsPage';

test.describe('Extensions Page Tests', () => {
  let extensionsPage;

  test.beforeEach(async ({ page }) => {
    extensionsPage = new ExtensionsPage(page);
    await page.goto('/extensions');
  });

  test('should verify Kanvas Snapshot visibility', async () => {
    await extensionsPage.verifyKanvasSnapshotVisibility();
  });

  test('should click Kanvas Snapshot enable button', async () => {
    await extensionsPage.clickKanvasSnapshotEnableBtn();
  });

  test('should verify Performance Analysis visibility', async () => {
    await extensionsPage.verifyPerformanceAnalysisVisibility();
  });

  test('should click Performance Analysis enable button', async () => {
    await extensionsPage.clickPerformanceAnalysisEnableBtn();
  });

  test('should verify Kanvas Details visibility', async () => {
    await extensionsPage.verifyKanvasDetailsVisibility();
  });

  test('should click Kanvas Signup button and open new page', async ({ context }) => {
    const newPage = await extensionsPage.clickAndWaitForNewPage(context, () =>
      extensionsPage.clickKanvasSignupBtn(),
    );
    await newPage.close();
  });

  test('should verify Docker Extension visibility', async () => {
    await extensionsPage.verifyDockerExtensionVisibility();
  });

  test('should click Docker Extension download button', async ({ context }) => {
    const newPage = await extensionsPage.clickAndWaitForNewPage(context, () =>
      extensionsPage.clickDockerExtensionDownloadBtn(),
    );
    await newPage.close();
  });

  test('should verify Design Embed visibility', async () => {
    await extensionsPage.verifyDesignEmbedVisibility();
  });

  test('should click Design Embed learn more button', async ({ context }) => {
    const newPage = await extensionsPage.clickAndWaitForNewPage(context, () =>
      extensionsPage.clickDesignEmbedLearnMoreBtn(),
    );
    await newPage.close();
  });

  test('should verify Catalog Section visibility', async () => {
    await extensionsPage.verifyCatalogSectionVisibility();
  });

  test('should click Catalog toggle switch', async () => {
    await extensionsPage.clickCatalogToggleSwitch();
  });

  test('should click Catalog link and verify URL', async ({ context }) => {
    await extensionsPage.verifyNewPageURL(
      context,
      () => extensionsPage.clickCatalogLink(),
      'https://meshery.io/catalog',
    );
  });

  test('should click Adapter docs Istio', async ({ context }) => {
    const newPage = await extensionsPage.clickAndWaitForNewPage(context, () =>
      extensionsPage.clickAdapterDocsIstio(),
    );
    await newPage.close();
  });
});
