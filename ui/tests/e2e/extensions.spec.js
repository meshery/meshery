import { test } from '@playwright/test';
import { ExtensionsPage } from '../pages/extensionsPage.js';

test.describe('Extensions Section Tests', () => {
  let extensionsPage;

  test.beforeEach(async ({ page }) => {
    extensionsPage = new ExtensionsPage(page);
    await extensionsPage.navigate();
  });

  test('Verify Kanvas Snapshot using data-testid', async () => {
    await extensionsPage.verifyKanvasSnapshot();
  });

  test('Verify Performance Analysis Details', async () => {
    await extensionsPage.verifyPerformanceAnalysis();
  });

  test('Verify Kanvas Details', async ({ context }) => {
    await extensionsPage.verifyKanvasDetails(context);
  });

  test('Verify Meshery Docker Extension Details', async ({ context }) => {
    await extensionsPage.verifyDockerExtension(context);
  });

  test('Verify Meshery Design Embed Details', async ({ context }) => {
    await extensionsPage.verifyDesignEmbed(context);
  });

  test('Verify Meshery Catalog Section Details', async ({ context }) => {
    await extensionsPage.verifyCatalogSection(context);
  });

  test('Verify Meshery Adapter for Istio Section', async ({ context }) => {
    await extensionsPage.verifyAdapterIstio(context);
  });
});
